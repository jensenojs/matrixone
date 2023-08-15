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
package checkpkdup

import (
	"bytes"

	"github.com/cespare/xxhash/v2"
	"github.com/matrixorigin/matrixone/pkg/common/moerr"
	"github.com/matrixorigin/matrixone/pkg/common/mpool"
	"github.com/matrixorigin/matrixone/pkg/container/vector"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/containers"
	"github.com/matrixorigin/matrixone/pkg/vm/process"
)

func String(_ any, buf *bytes.Buffer) {
	buf.WriteString(" check primarykey dup")
}

func Prepare(proc *process.Process, arg any) (err error) {
	ap := arg.(*Argument)
	ap.potentially_dup_keys = make(map[any]bool)
	ap.bitmp.InitWithSize(5 * mpool.MB)
	return nil
}

func hash(v []byte) uint64 {
	return xxhash.Sum64(v)
}

func Call(idx int, proc *process.Process, arg any, isFirst bool, isLast bool) (process.ExecStatus, error) {
	anal := proc.GetAnalyze(idx)
	anal.Start()
	defer anal.Stop()
	ap := arg.(*Argument)
	bat := proc.InputBatch()
	anal.Input(bat, isFirst)

	if bat == nil {
		if len(ap.potentially_dup_keys) == 0 {
			return process.ExecStop, nil
		} else {
			// need to run background SQL to double check here
		}
	}
	if bat.RowCount() == 0 {
		return process.ExecNext, nil
	}

	// build hash array from pk col values

	toCheckVec := bat.GetVector(0) // firstly assume that the index of the column to be checked is 0
	toCheckBytes := vector.MustBytesCol(toCheckVec)
	hashes := make([]uint64, toCheckVec.Length())
	for idx, bytes := range toCheckBytes {
		hashes[idx] = hash(bytes)
	}

	// insert into bitmap

	if ap.bitmp.EmptyByFlag() {
		ap.bitmp.AddMany(hashes)
	} else {
		for idx, val := range hashes {
			if ap.bitmp.Contains(val) {
				mayDupKey := containers.GetNonNullValue(toCheckVec, uint32(idx))
				if _, ok := ap.potentially_dup_keys[mayDupKey]; ok {
					dupentry, _ := mayDupKey.(string)
					return process.ExecStop, moerr.NewDuplicateEntry(proc.Ctx, dupentry, bat.Attrs[0])
				}
				ap.potentially_dup_keys[mayDupKey] = true
			} else {
				ap.bitmp.Add(val)
			}
		}
	}

	return process.ExecNext, nil
}
