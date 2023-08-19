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
	"fmt"

	"github.com/cespare/xxhash/v2"
	"github.com/matrixorigin/matrixone/pkg/common/bitmap"
	"github.com/matrixorigin/matrixone/pkg/common/moerr"
	"github.com/matrixorigin/matrixone/pkg/common/mpool"
	"github.com/matrixorigin/matrixone/pkg/container/batch"
	"github.com/matrixorigin/matrixone/pkg/container/vector"
	"github.com/matrixorigin/matrixone/pkg/vm/process"
)

// This operator is used to implement a way to ensure primary keys/unique keys are not duplicate in `INSERT` and `LOAD` statements
// There are two conditions needed to check
// 	 new pk/uk are not duplicate with each other
// 	 new pk/uk are not duplicate with existing data
//
// THE big idea is to store
// 	 pk columns to be loaded
//   pk columns already exist
// both in a Bloom Filter-like data structure, let's say bloom filter below
//
// if the final Bloom Filter claim that
// 	 case 1: have no duplicate keys,
// 		passes duplicate constraint
//   case 2: may have duplicate keys
//      start a background SQL to double check, it should be a point select
//
// However, backgroudSQL may slow, so we can do some optimizations
//   1. an record the keys that have hash conflict, and manually check them whather duplicate or not
//   2. shuffle

const (
	// One hundred million can be estimated as 2^26, 1MB can be estimated as 2^20.
	bitmpSize = 1 * mpool.MB

	// bitmap always assume that bitmap has been extended to at least row
	// but currently have no logic about Expand bitmap
	bitMask = 0xfffff
)

func String(_ any, buf *bytes.Buffer) {
	buf.WriteString(" fuzzy check duplicate constraint")
}

func Prepare(proc *process.Process, arg any) (err error) {
	ap := arg.(*Argument)
	ap.ctr = new(container)
	ap.ctr.filter = new(bitmap.Bitmap)
	ap.ctr.filter.InitWithSize(bitmpSize)
	ap.ctr.mayDuplicate = make(map[string]bool)
	return nil
}

func Call(idx int, proc *process.Process, arg any, isFirst bool, isLast bool) (process.ExecStatus, error) {
	anal := proc.GetAnalyze(idx)
	anal.Start()
	defer anal.Stop()

	ap := arg.(*Argument)
	ctr := ap.ctr
	bat := proc.InputBatch()
	anal.Input(bat, isFirst)

	if bat == nil {
		if ctr.collisionKeys.Length() == 0 {
			// case 1:
			proc.SetInputBatch(nil)
			return process.ExecStop, nil
		} else {
			// case 2: send collisionKeys to output operator to run background SQL
			rbat := batch.New(true, ctr.toCheckAttr)
			rbat.SetVector(0, ctr.collisionKeys)
			rbat.SetRowCount(ctr.collisionKeys.Length())
			proc.SetInputBatch(rbat)
			return process.ExecNext, nil
		}
	}

	var hashes []uint64
	rowCnt := bat.RowCount()
	if rowCnt == 0 {
		proc.PutBatch(bat)
		proc.SetInputBatch(batch.EmptyBatch)
		return process.ExecNext, nil
	} else {
		hashes = make([]uint64, rowCnt)
	}

	toCheckVec := bat.GetVector(0)
	if ctr.collisionKeys == nil {
		ctr.collisionKeys = vector.NewVec(*toCheckVec.GetType())
		ctr.toCheckAttr = bat.Attrs
	}

	// build hash array from pk col
	generateHashes(toCheckVec, rowCnt, hashes)

	//  check hashes in to fuzzy filter
	for idx, hval := range hashes {
		if ctr.filter.Contains(hval) {
			ctr.collisionKeys.UnionOne(toCheckVec, int64(idx), proc.GetMPool())

			// FIXME: %v is the right way to go ?
			collisionKey := fmt.Sprintf("%v", getNonNullValue(toCheckVec, uint32(idx)))
			if _, ok := ctr.mayDuplicate[collisionKey]; ok {
				// key that has been inserted before, fail to pass constraint
				return process.ExecStop, moerr.NewDuplicateEntry(proc.Ctx, collisionKey, bat.Attrs[0])
			}
			ctr.mayDuplicate[collisionKey] = true
		} else {
			ctr.filter.Add(hval)
		}
	}

	proc.SetInputBatch(batch.EmptyBatch)
	return process.ExecNext, nil
}

func generateHashes(pkCol *vector.Vector, rowCnt int, hashes []uint64) {
	for idx := 0; idx < rowCnt; idx++ {
		hashes[idx] = xxhash.Sum64(pkCol.GetRawBytesAt(idx)) & bitMask
	}
}
