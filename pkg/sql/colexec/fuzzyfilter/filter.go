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
// 	 new pk/uk are not duplicate with each other
// 	 new pk/uk are not duplicate with existing data
//
// THE big idea is to store
// 	 pk columns to be loaded
//   pk columns already exist
// both in a bitmap-like data structure, let's say bitmap below
//
// if the final bitmap claim that
// 	 case 1: have no duplicate keys,
// 		passes duplicate constraint
//   case 2: may have duplicate keys because of hash collision
//      start a background SQL to double check
//
// However, backgroud SQL may slow, so we can do some optimizations
//   1. an record the keys that have hash collision, and manually check them whether duplicate or not,
//         if duplicate, then return error timely
//   2. shuffle inner operator and between operations

const (
	// One hundred million can be estimated as 2^26, 1MB can be estimated as 2^20.
	bitmpSize = 1 * mpool.MB
	// bitmap always assume that bitmap has been extended to at least row
	// but currently have no logic about Expand bitmap
	bitMask = 0xfffff

	// index for bat that fuzzy filter pass to next operator
	dbIdx   = 0 // where store dbName
	tblIdx  = 1 // where store tblName
	attrIdx = 2 // where store collsion keys for to check attribute
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
		collisionCnt := ctr.rbat.GetVector(attrIdx).Length()

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

	toCheck := bat.GetVector(0)
	if ctr.rbat == nil {
		if err := wrapUpCollisionBatch(proc, ap, bat.Attrs[0], *toCheck.GetType()); err != nil {
			return process.ExecStop, err
		}
	}

	// build hash array from pk col
	generateHashes(toCheck, rowCnt, hashes)

	// check hashes in to fuzzy filter
	for idx, hval := range hashes {
		if ctr.filter.Contains(hval) {
			ctr.rbat.GetVector(attrIdx).UnionOne(toCheck, int64(idx), proc.GetMPool())

			// opts 1
			collisionKey := fmt.Sprintf("%v", getNonNullValue(toCheck, uint32(idx)))
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

func wrapUpCollisionBatch(proc *process.Process, arg *Argument, attrName string, attrTyp types.Type) error {
	attrs := []string{"toCheckDb", "toCheckTbl", attrName}
	rbat := batch.New(true, attrs)

	db := vector.NewConstBytes(types.T_binary.ToType(), []byte(arg.DbName), 1, proc.GetMPool())
	tbl := vector.NewConstBytes(types.T_binary.ToType(), []byte(arg.TblName), 1, proc.GetMPool())

	rbat.SetVector(dbIdx, db)
	rbat.SetVector(tblIdx, tbl)
	rbat.SetVector(attrIdx, vector.NewVec(attrTyp))

	rbat.SetRowCount(1) // only for pass for output operater
	arg.ctr.rbat = rbat
	return nil
}
