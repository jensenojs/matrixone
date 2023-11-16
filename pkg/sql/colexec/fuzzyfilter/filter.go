// Copyright 2023 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package fuzzyfilter

import (
	"bytes"
	// "fmt"
	"math"

	"github.com/bits-and-blooms/bloom"
	"github.com/matrixorigin/matrixone/pkg/common/moerr"
	"github.com/matrixorigin/matrixone/pkg/container/batch"
	"github.com/matrixorigin/matrixone/pkg/vm"
	"github.com/matrixorigin/matrixone/pkg/vm/process"
)

/*
This operator is used to implement a way to ensure primary keys/unique keys are not duplicate in `INSERT` and `LOAD` statements

the BIG idea is to store
    pk columns to be loaded
    pk columns already exist
both in a bitmap-like data structure, let's say bloom filter below

if the final bloom filter claim that
    case 1: have no duplicate keys
        pass duplicate constraint directly
    case 2: Not sure if there are duplicate keys because of hash collision
        start a background SQL to double check

Note:
1. backgroud SQL may slow, so some optimizations could be applied
    manually check whether collision keys duplicate or not,
        if duplicate, then return error timely
2. there is a corner case that no need to run background SQL
    on duplicate key update
*/

const (
	// Probability of false positives
	p float64 = 0.00001
	// Number of hash functions
	k uint = 3
)

// EstimateBitsNeed return the Number of bits should have in the filter
// by the formula: p = pow(1 - exp(-k / (m / n)), k)
//
//	==> m = - kn / ln(1 - p^(1/k)), use k * (1.001 * n) instead of kn to overcome floating point errors
func EstimateBitsNeed(n float64, k uint, p float64) float64 {
	return -float64(k) * math.Ceil(1.001*n) / math.Log(1-math.Pow(p, 1.0/float64(k)))
}

func (arg *Argument) String(buf *bytes.Buffer) {
	buf.WriteString(" fuzzy check duplicate constraint")
}

func (arg *Argument) Prepare(proc *process.Process) (err error) {
	e := EstimateBitsNeed(arg.N, k, p)
	m := uint(math.Ceil(e))
	if float64(m) < e {
		return moerr.NewInternalErrorNoCtx("Overflow occurred when estimating size of fuzzy filter")
	}
	arg.filter = bloom.New(m, k)
	return nil
}

func (arg *Argument) Call(proc *process.Process) (vm.CallResult, error) {
	anal := proc.GetAnalyze(arg.info.Idx)
	anal.Start()
	defer anal.Stop()

	result, err := arg.children[0].Call(proc)
	if err != nil {
		result.Status = vm.ExecStop
		return result, err
	}
	bat := result.Batch

	if bat == nil {
		// this will happen in such case:create unique index from a table that unique col have no data
		if arg.rbat == nil || arg.collisionCnt == 0 {
			result.Status = vm.ExecStop
			return result, nil
		}

		// fmt.Printf("Estimated row count is %f, collisionCnt is %d, fp is %f\n", ap.N, ap.collisionCnt, float64(ap.collisionCnt)/float64(ap.N))
		// send collisionKeys to output operator to run background SQL
		arg.rbat.SetRowCount(arg.rbat.Vecs[0].Length())
		result.Batch = arg.rbat
		arg.collisionCnt = 0
		result.Status = vm.ExecStop
		return result, nil
	}

	anal.Input(bat, arg.info.IsFirst)

	rowCnt := bat.RowCount()
	if rowCnt == 0 {
		return result, nil
	}

	if arg.rbat == nil {
		if err := generateRbat(proc, arg, bat); err != nil {
			result.Status = vm.ExecStop
			return result, err
		}
	}

	pkCol := bat.GetVector(0)
	for i := 0; i < rowCnt; i++ {
		var bytes = pkCol.GetRawBytesAt(i)
		if arg.filter.TestAndAdd(bytes) {
			appendCollisionKey(proc, arg, i, bat)
			arg.collisionCnt++
		}
	}
	result.Batch = batch.EmptyBatch
	return result, nil
}

// appendCollisionKey will append collision key into rbat
func appendCollisionKey(proc *process.Process, arg *Argument, idx int, bat *batch.Batch) {
	pkCol := bat.GetVector(0)
	arg.rbat.GetVector(0).UnionOne(pkCol, int64(idx), proc.GetMPool())
}

// rbat will contain the keys that have hash collisions
func generateRbat(proc *process.Process, arg *Argument, bat *batch.Batch) error {
	rbat := batch.NewWithSize(1)
	rbat.SetVector(0, proc.GetVector(*bat.GetVector(0).GetType()))
	arg.rbat = rbat
	return nil
}
