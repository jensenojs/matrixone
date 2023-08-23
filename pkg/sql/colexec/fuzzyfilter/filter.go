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
	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/container/vector"
	"github.com/matrixorigin/matrixone/pkg/logutil"
	"github.com/matrixorigin/matrixone/pkg/vm/process"
)

// This operator is used to implement a way to ensure primary keys/unique keys are not duplicate in `INSERT` and `LOAD` statements
// There are two conditions needed to check
//    new pk/uk are not duplicate with each other
//    new pk/uk are not duplicate with existing data
//
// the BIG idea is to store
//    pk columns to be loaded
//    pk columns already exist
// both in a bitmap-like data structure, let's say bitmap below
//
// if the final bitmap claim that
//    case 1: have no duplicate keys,
// 	      passes duplicate constraint
//    case 2: Not sure if there are duplicate keys because of hash collision
//        start a background SQL to double check with func doubleCheckForFuzzyFilter
//
// However, backgroud SQL may slow, some optimizations could be applied
//    1. manually check whether collision keys duplicate or not,
//         if duplicate, then return error timely
//    2. shuffle inner operator and between operations

const (
	// One hundred million can be estimated as 2^26, 1MB can be estimated as 2^20.
	bitmpSize = 1 * mpool.MB
	// bitmap always assume that bitmap has been extended to at least row
	// but currently have no logic about Expand bitmap
	bitMask = 0xfffff
)

func generateHashes(bat *batch.Batch, rowCnt int, hashes []uint64) {
	pkCol := bat.GetVector(0)
	for idx := 0; idx < rowCnt; idx++ {
		hashes[idx] = xxhash.Sum64(pkCol.GetRawBytesAt(idx)) & bitMask
	}
}

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
		collisionCnt := ctr.rbat.GetVector(0).Length()

		if collisionCnt == 0 {
			// case 1: pass duplicate constraint
			proc.SetInputBatch(nil)
			return process.ExecStop, nil
		} else {
			if collisionCnt > 100 {
				logutil.Warnf("the hash collision cnt for fuzzy filter is too large : %d", collisionCnt)
			}
			// case 2: send collisionKeys to output operator to run background SQL
			proc.SetInputBatch(ctr.rbat)
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

	if ctr.rbat == nil {
		if err := generateRbat(proc, ap, bat); err != nil {
			return process.ExecStop, err
		}
	}

	// build hash array from pk col
	generateHashes(bat, rowCnt, hashes)

	// check hashes in to fuzzy filter
	for idx, hval := range hashes {
		if ctr.filter.Contains(hval) {
			appendCollisionKey(proc, ctr, idx, bat)

			// opts 1
			if err := checkIfDup(proc, ctr, idx, bat); err != nil {
				return process.ExecStop, err
			}

		} else {
			ctr.filter.Add(hval)
		}
	}

	proc.SetInputBatch(batch.EmptyBatch)
	return process.ExecNext, nil
}

// appendCollisionKey will append collision key into rbat
func appendCollisionKey(proc *process.Process, ctr *container, idx int, bat *batch.Batch) {
	if !ctr.isCpk {
		pkCol := bat.GetVector(0)
		ctr.rbat.GetVector(0).UnionOne(pkCol, int64(idx), proc.GetMPool())
	} else {
		for i := 1; i < len(bat.Attrs); i++ { // skip for __mo_cpkey_col
			partialPkCol := bat.GetVector(int32(i))
			ctr.rbat.GetVector(int32(i-1)).UnionOne(partialPkCol, int64(idx), proc.GetMPool())
		}
	}
}

// optimization 1
func checkIfDup(proc *process.Process, ctr *container, idx int, bat *batch.Batch) error {
	var buf bytes.Buffer
	if !ctr.isCpk {
		pkCol := bat.GetVector(0)
		collisionkey := getNonNullValue(pkCol, uint32(idx))
		buf.WriteString(fmt.Sprintf("%v", collisionkey))
	} else {
		buf.WriteString(fmt.Sprintf("("))

		var i int32
		var partialCollsionKey any

		for i = 1; i < int32(len(bat.Attrs)-1); i++ { // skip for __mo_cpkey_col
			partialCollsionKey = getNonNullValue(bat.GetVector(int32(i)), uint32(idx))
			buf.WriteString(fmt.Sprintf("%v,", partialCollsionKey))
		}

		partialCollsionKey = getNonNullValue(bat.GetVector(int32(i)), uint32(idx))
		buf.WriteString(fmt.Sprintf("%v)", partialCollsionKey))
	}
	s := buf.String()
	if _, ok := ctr.mayDuplicate[s]; ok {
		// key that has been inserted before, fail to pass dup constraint
		return moerr.NewDuplicateEntry(proc.Ctx, s, bat.Attrs[0])
	}
	ctr.mayDuplicate[s] = true
	return nil
}

// rbat will contain the keys that have hash collisions, it only store partial keys when isCpk is true
// it will also create two additional vectors at the end to store the database name and table name
//
// refer1 func doubleCheckForFuzzyFilter
// refer2 func compileFuzzyFilter
func generateRbat(proc *process.Process, arg *Argument, bat *batch.Batch) error {
	var rattrs []string
	baCnt := len(bat.Attrs)
	if baCnt > 1 {
		// handle with composite primary key
		arg.ctr.isCpk = true
		// skip for __mo_cpkey_col as we only store partial pk value
		rattrs = append(rattrs, bat.Attrs[1:]...)
	} else {
		rattrs = bat.Attrs
	}

	rattrs = append(rattrs, []string{"toCheckDb", "toCheckTbl"}...)
	rbat := batch.New(true, rattrs)

	if !arg.ctr.isCpk {
		rbat.SetVector(0, vector.NewVec(*bat.GetVector(0).GetType()))
	} else {
		for i := 1; i < baCnt; i++ {
			rbat.SetVector(int32(i-1), vector.NewVec(*bat.GetVector(int32(i)).GetType()))
		}
	}

	db := vector.NewConstBytes(types.T_binary.ToType(), []byte(arg.DbName), 1, proc.GetMPool())
	rbat.SetVector(int32(len(rattrs)-1-1), db)
	tbl := vector.NewConstBytes(types.T_binary.ToType(), []byte(arg.TblName), 1, proc.GetMPool())
	rbat.SetVector(int32(len(rattrs)-1), tbl)

	// no means, only for pass for output operater
	rbat.SetRowCount(1)
	arg.ctr.rbat = rbat
	return nil
}
