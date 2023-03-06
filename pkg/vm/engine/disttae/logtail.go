// Copyright 2022 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package disttae

import (
	"context"

	"github.com/matrixorigin/matrixone/pkg/catalog"
	"github.com/matrixorigin/matrixone/pkg/container/batch"
	"github.com/matrixorigin/matrixone/pkg/pb/api"
)

func consumeEntry(
	ctx context.Context,
	primaryIdx int,
	engine *Engine,
	state *PartitionState,
	e *api.Entry,
) error {

	state.HandleLogtailEntry(ctx, e, primaryIdx, engine.mp)
	if isMetaTable(e.TableName) {
		return nil
	}

	if e.EntryType == api.Entry_Insert {
		switch e.TableId {
		case catalog.MO_TABLES_ID:
			bat, _ := batch.ProtoBatchToBatch(e.Bat)
			engine.catalog.InsertTable(bat)
		case catalog.MO_DATABASE_ID:
			bat, _ := batch.ProtoBatchToBatch(e.Bat)
			engine.catalog.InsertDatabase(bat)
		case catalog.MO_COLUMNS_ID:
			bat, _ := batch.ProtoBatchToBatch(e.Bat)
			engine.catalog.InsertColumns(bat)
		}
		return nil
	}

	switch e.TableId {
	case catalog.MO_TABLES_ID:
		bat, _ := batch.ProtoBatchToBatch(e.Bat)
		engine.catalog.DeleteTable(bat)
	case catalog.MO_DATABASE_ID:
		bat, _ := batch.ProtoBatchToBatch(e.Bat)
		engine.catalog.DeleteDatabase(bat)
	}
	return nil

	/*
		if isMetaTable(e.TableName) {
			return nil
		}

		if e.EntryType == api.Entry_Insert {
			if isMetaTable(e.TableName) {
				vec, err := vector.ProtoVectorToVector(e.Bat.Vecs[catalog.BLOCKMETA_ID_IDX+MO_PRIMARY_OFF])
				if err != nil {
					return err
				}
				timeVec, err := vector.ProtoVectorToVector(e.Bat.Vecs[catalog.BLOCKMETA_COMMITTS_IDX+MO_PRIMARY_OFF])
				if err != nil {
					return err
				}
				vec1, err := vector.ProtoVectorToVector(e.Bat.Vecs[catalog.BLOCKMETA_ENTRYSTATE_IDX+MO_PRIMARY_OFF])
				if err != nil {
					return err
				}
				vs := vector.MustTCols[uint64](vec)
				timestamps := vector.MustTCols[types.TS](timeVec)
				es := vector.MustTCols[bool](vec1)
				for i, v := range vs {
					if e.TableId == testTableId {
						logutil.Errorf("+++: deleted by block id %d, entrystate is %v, ts is %v", v, es[i], timestamps[i].ToTimestamp())
					}
					if err := tbl.parts[idx].DeleteByBlockID(ctx, timestamps[i].ToTimestamp(), v); err != nil {
						if !moerr.IsMoErrCode(err, moerr.ErrTxnWriteConflict) { // 1
							return err
						}
					}
				}
				return engine.getMetaPartitions(e.TableName)[idx].Insert(ctx, -1, e.Bat, false)
			}
			switch e.TableId {
			case catalog.MO_TABLES_ID:
				bat, _ := batch.ProtoBatchToBatch(e.Bat)
				engine.catalog.InsertTable(bat)
			case catalog.MO_DATABASE_ID:
				bat, _ := batch.ProtoBatchToBatch(e.Bat)
				engine.catalog.InsertDatabase(bat)
			case catalog.MO_COLUMNS_ID:
				bat, _ := batch.ProtoBatchToBatch(e.Bat)
				engine.catalog.InsertColumns(bat)
			}
			return nil
		}

		switch e.TableId {
		case catalog.MO_TABLES_ID:
			bat, _ := batch.ProtoBatchToBatch(e.Bat)
			engine.catalog.DeleteTable(bat)
		case catalog.MO_DATABASE_ID:
			bat, _ := batch.ProtoBatchToBatch(e.Bat)
			engine.catalog.DeleteDatabase(bat)
		}
	*/
}
