// Copyright 2021 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plan

import (
	"github.com/matrixorigin/matrixone/pkg/catalog"
	"github.com/matrixorigin/matrixone/pkg/common/moerr"
	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/pb/plan"
	"github.com/matrixorigin/matrixone/pkg/sql/parsers/tree"
)

func (builder *QueryBuilder) bindInsert(stmt *tree.Insert, ctx *BindContext) (int32, error) {
	if len(stmt.OnDuplicateUpdate) > 1 || stmt.OnDuplicateUpdate[0] != nil {
		return 0, moerr.NewUnsupportedDML(builder.GetContext(), "on duplicate key update")
	}

	tbl := stmt.Table.(*tree.TableName)
	dbName := string(tbl.SchemaName)
	tblName := string(tbl.ObjectName)
	if len(dbName) == 0 {
		dbName = builder.compCtx.DefaultDatabase()
	}

	_, t := builder.compCtx.Resolve(dbName, tblName, nil)
	if t == nil {
		return 0, moerr.NewNoSuchTable(builder.compCtx.GetContext(), dbName, tblName)
	}
	if t.TableType == catalog.SystemSourceRel {
		return 0, moerr.NewNYIf(builder.compCtx.GetContext(), "insert stream %s", tblName)
	}

	tblInfo, err := getDmlTableInfo(builder.compCtx, tree.TableExprs{stmt.Table}, nil, nil, "insert")
	if err != nil {
		return 0, err
	}
	rewriteInfo := &dmlSelectInfo{
		typ:     "insert",
		rootId:  -1,
		tblInfo: tblInfo,
	}
	tableDef := tblInfo.tableDefs[0]
	// clusterTable, err := getAccountInfoOfClusterTable(ctx, stmt.Accounts, tableDef, tblInfo.isClusterTable[0])
	// if err != nil {
	// 	return 0, err
	// }
	// if len(stmt.OnDuplicateUpdate) > 0 && clusterTable.IsClusterTable {
	// 	return 0, moerr.NewNotSupported(builder.compCtx.GetContext(), "INSERT ... ON DUPLICATE KEY UPDATE ... for cluster table")
	// }

	builder.haveOnDuplicateKey = len(stmt.OnDuplicateUpdate) > 0
	if stmt.IsRestore {
		builder.isRestore = true
		oldSnapshot := builder.compCtx.GetSnapshot()
		builder.compCtx.SetSnapshot(&Snapshot{
			Tenant: &plan.SnapshotTenant{
				TenantName: "xxx",
				TenantID:   stmt.FromDataTenantID,
			},
		})
		defer func() {
			builder.compCtx.SetSnapshot(oldSnapshot)
		}()
	}

	_, _, colName2Idx, err := builder.initInsertStmt(ctx, stmt, rewriteInfo)
	if err != nil {
		return 0, err
	}
	replaceStmt := getRewriteToReplaceStmt(tableDef, stmt, rewriteInfo, builder.isPrepareStatement)
	if replaceStmt != nil {
		return 0, moerr.NewUnsupportedDML(builder.GetContext(), "rewrite to replace")
	}
	lastNodeID := rewriteInfo.rootId

	//colName2Idx := make(map[string]int)

	onDupAction := plan.Node_ERROR
	if len(stmt.OnDuplicateUpdate) == 1 && stmt.OnDuplicateUpdate[0] == nil {
		onDupAction = plan.Node_IGNORE
	}

	selectNode := builder.qry.Nodes[lastNodeID]
	idxScanNodes := make([][]*plan.Node, len(tblInfo.tableDefs))

	for _, tableDef := range tblInfo.tableDefs {
		for _, idxDef := range tableDef.Indexes {
			if !catalog.IsRegularIndexAlgo(idxDef.IndexAlgo) {
				return 0, moerr.NewUnsupportedDML(builder.GetContext(), "have vector index table")
			}

		}
	}

	for i, tableDef := range tblInfo.tableDefs {
		scanTag := builder.genNewTag()
		builder.addNameByColRef(scanTag, tableDef)

		scanNodeID := builder.appendNode(&plan.Node{
			NodeType:     plan.Node_TABLE_SCAN,
			TableDef:     tableDef,
			ObjRef:       tblInfo.objRef[i],
			BindingTags:  []int32{scanTag},
			ScanSnapshot: ctx.snapshot,
		}, ctx)

		pkName := tableDef.Pkey.PkeyColName
		pkPos := colName2Idx[tableDef.Name+"."+tableDef.Pkey.PkeyColName]
		pkTyp := tableDef.Cols[pkPos].Typ
		leftExpr := &plan.Expr{
			Typ: pkTyp,
			Expr: &plan.Expr_Col{
				Col: &plan.ColRef{
					RelPos: scanTag,
					ColPos: pkPos,
				},
			},
		}

		rightExpr := &plan.Expr{
			Typ: pkTyp,
			Expr: &plan.Expr_Col{
				Col: &plan.ColRef{
					RelPos: selectNode.BindingTags[0],
					ColPos: colName2Idx[tableDef.Name+"."+pkName],
				},
			},
		}

		joinCond, _ := BindFuncExprImplByPlanExpr(builder.GetContext(), "=", []*plan.Expr{
			leftExpr,
			rightExpr,
		})

		lastNodeID = builder.appendNode(&plan.Node{
			NodeType:          plan.Node_DEDUP_JOIN,
			Children:          []int32{scanNodeID, lastNodeID},
			OnList:            []*plan.Expr{joinCond},
			OnDuplicateAction: onDupAction,
		}, ctx)

		idxScanNodes[i] = make([]*plan.Node, len(tableDef.Indexes))

		for j, idxDef := range tableDef.Indexes {
			if !idxDef.TableExist {
				continue
			}

			idxObjRef, idxTableDef := builder.compCtx.Resolve(tblInfo.objRef[0].SchemaName, idxDef.IndexTableName, nil)
			colName2Idx[idxTableDef.Name+"."+catalog.IndexTablePrimaryColName] = pkPos

			argsLen := len(idxDef.Parts)
			if !idxDef.Unique {
				argsLen++
			}

			if argsLen == 1 {
				colName2Idx[idxTableDef.Name+"."+catalog.IndexTableIndexColName] = colName2Idx[tableDef.Name+"."+idxDef.Parts[0]]
			} else {
				args := make([]*plan.Expr, argsLen)

				for k, part := range idxDef.Parts {
					colPos := colName2Idx[tableDef.Name+"."+part]
					args[k] = DeepCopyExpr(selectNode.ProjectList[colPos])
				}

				if !idxDef.Unique {
					args[len(idxDef.Parts)] = DeepCopyExpr(selectNode.ProjectList[pkPos])
				}

				idxExpr, _ := BindFuncExprImplByPlanExpr(builder.GetContext(), "serial_full", args)
				colName2Idx[idxTableDef.Name+"."+catalog.IndexTableIndexColName] = int32(len(selectNode.ProjectList))
				selectNode.ProjectList = append(selectNode.ProjectList, idxExpr)
			}

			if !idxDef.Unique {
				continue
			}

			idxTag := builder.genNewTag()
			builder.addNameByColRef(idxTag, idxTableDef)

			idxScanNodes[i][j] = &plan.Node{
				NodeType:     plan.Node_TABLE_SCAN,
				TableDef:     idxTableDef,
				ObjRef:       idxObjRef,
				BindingTags:  []int32{idxTag},
				ScanSnapshot: ctx.snapshot,
			}
			idxTableNodeID := builder.appendNode(idxScanNodes[i][j], ctx)

			idxPkPos := idxTableDef.Name2ColIndex[catalog.IndexTableIndexColName]
			pkTyp := idxTableDef.Cols[idxPkPos].Typ

			leftExpr := &plan.Expr{
				Typ: pkTyp,
				Expr: &plan.Expr_Col{
					Col: &plan.ColRef{
						RelPos: idxTag,
						ColPos: idxPkPos,
					},
				},
			}

			rightExpr := &plan.Expr{
				Typ: pkTyp,
				Expr: &plan.Expr_Col{
					Col: &plan.ColRef{
						RelPos: selectNode.BindingTags[0],
						ColPos: colName2Idx[idxTableDef.Name+"."+catalog.IndexTableIndexColName],
					},
				},
			}

			joinCond, _ := BindFuncExprImplByPlanExpr(builder.GetContext(), "=", []*plan.Expr{
				leftExpr,
				rightExpr,
			})

			lastNodeID = builder.appendNode(&plan.Node{
				NodeType:          plan.Node_DEDUP_JOIN,
				Children:          []int32{idxTableNodeID, lastNodeID},
				OnList:            []*plan.Expr{joinCond},
				OnDuplicateAction: onDupAction,
			}, ctx)
		}
	}

	for i, tableDef := range tblInfo.tableDefs {
		if tableDef.Pkey.PkeyColName == catalog.FakePrimaryKeyColName {
			continue
		}

		pkPos, pkTyp := getPkPos(tableDef, false)

		lockTarget := &plan.LockTarget{
			TableId:            tableDef.TblId,
			PrimaryColIdxInBat: int32(pkPos),
			PrimaryColTyp:      pkTyp,
			RefreshTsIdxInBat:  -1, //unsupported now
		}

		//if tableDef.Partition != nil {
		//	lockTarget.IsPartitionTable = true
		//	lockTarget.FilterColIdxInBat = int32(partitionIdx)
		//	lockTarget.PartitionTableIds = partTableIDs
		//}

		lastNodeID = builder.appendNode(&plan.Node{
			NodeType:    plan.Node_LOCK_OP,
			Children:    []int32{lastNodeID},
			LockTargets: []*plan.LockTarget{lockTarget},
		}, ctx)

		for _, idxNode := range idxScanNodes[i] {
			if idxNode == nil {
				continue
			}

			pkPos, pkTyp := getPkPos(idxNode.TableDef, false)

			lockTarget := &plan.LockTarget{
				TableId:            idxNode.TableDef.TblId,
				PrimaryColIdxInBat: int32(pkPos),
				PrimaryColTyp:      pkTyp,
				RefreshTsIdxInBat:  -1, //unsupported now
			}

			//if tableDef.Partition != nil {
			//	lockTarget.IsPartitionTable = true
			//	lockTarget.FilterColIdxInBat = int32(partitionIdx)
			//	lockTarget.PartitionTableIds = partTableIDs
			//}

			lastNodeID = builder.appendNode(&plan.Node{
				NodeType:    plan.Node_LOCK_OP,
				Children:    []int32{lastNodeID},
				LockTargets: []*plan.LockTarget{lockTarget},
			}, ctx)
		}
	}

	for i, tableDef := range tblInfo.tableDefs {
		insertCols := make([]*plan.Expr, len(tableDef.Cols))
		for j, col := range tableDef.Cols {
			insertCols[j] = &plan.Expr{
				Typ: col.Typ,
				Expr: &plan.Expr_Col{
					Col: &plan.ColRef{
						RelPos: selectNode.BindingTags[0],
						ColPos: int32(colName2Idx[tableDef.Name+"."+col.Name]),
					},
				},
			}
		}

		lastNodeID = builder.appendNode(&plan.Node{
			NodeType: plan.Node_INSERT,
			Children: []int32{lastNodeID},
			ObjRef:   tblInfo.objRef[i],
			TableDef: tableDef,
			InsertCtx: &plan.InsertCtx{
				Ref:            tblInfo.objRef[i],
				IsClusterTable: tableDef.TableType == catalog.SystemClusterRel,
				TableDef:       tableDef,
				//PartitionTableIds:   paritionTableIds,
				//PartitionTableNames: paritionTableNames,
				//PartitionIdx:        int32(partitionIdx),
			},
			InsertDeleteCols: insertCols,
		}, ctx)

		for _, idxNode := range idxScanNodes[i] {
			if idxNode == nil {
				continue
			}

			insertCols := make([]*plan.Expr, len(idxNode.TableDef.Cols))
			for j, col := range idxNode.TableDef.Cols {
				insertCols[j] = &plan.Expr{
					Typ: col.Typ,
					Expr: &plan.Expr_Col{
						Col: &plan.ColRef{
							RelPos: selectNode.BindingTags[0],
							ColPos: int32(colName2Idx[tableDef.Name+"."+col.Name]),
						},
					},
				}
			}

			lastNodeID = builder.appendNode(&plan.Node{
				NodeType: plan.Node_INSERT,
				Children: []int32{lastNodeID},
				ObjRef:   idxNode.ObjRef,
				TableDef: idxNode.TableDef,
				InsertCtx: &plan.InsertCtx{
					Ref:            idxNode.ObjRef,
					IsClusterTable: idxNode.TableDef.TableType == catalog.SystemClusterRel,
					TableDef:       idxNode.TableDef,
					//PartitionTableIds:   paritionTableIds,
					//PartitionTableNames: paritionTableNames,
					//PartitionIdx:        int32(partitionIdx),
				},
				InsertDeleteCols: insertCols,
			}, ctx)
		}
	}

	reCheckifNeedLockWholeTable(builder)

	return lastNodeID, err
}

func (builder *QueryBuilder) initInsertStmt(bindCtx *BindContext, stmt *tree.Insert, info *dmlSelectInfo) (bool, map[string]bool, map[string]int32, error) {
	var err error
	colName2Idx := make(map[string]int32)
	// var uniqueCheckOnAutoIncr string
	var insertColumns []string
	tableDef := info.tblInfo.tableDefs[0]
	//tableObjRef := info.tblInfo.objRef[0]
	colToIdx := make(map[string]int)
	oldColPosMap := make(map[string]int)
	tableDef.Name2ColIndex = make(map[string]int32)
	for i, col := range tableDef.Cols {
		colToIdx[col.Name] = i
		oldColPosMap[col.Name] = i
		tableDef.Name2ColIndex[col.Name] = int32(i)
	}
	info.tblInfo.oldColPosMap = append(info.tblInfo.oldColPosMap, oldColPosMap)
	info.tblInfo.newColPosMap = append(info.tblInfo.newColPosMap, oldColPosMap)

	// dbName := string(stmt.Table.(*tree.TableName).SchemaName)
	// if dbName == "" {
	// 	dbName = builder.compCtx.DefaultDatabase()
	// }

	existAutoPkCol := false

	insertWithoutUniqueKeyMap := make(map[string]bool)
	//var ifInsertFromUniqueColMap map[string]bool
	if insertColumns, err = getInsertColsFromStmt(builder.GetContext(), stmt, tableDef); err != nil {
		return false, nil, nil, err
	}

	var astSelect *tree.Select
	switch selectImpl := stmt.Rows.Select.(type) {
	// rewrite 'insert into tbl values (1,1)' to 'insert into tbl select * from (values row(1,1))'
	case *tree.ValuesClause:
		isAllDefault := false
		if selectImpl.Rows[0] == nil {
			isAllDefault = true
		}
		if isAllDefault {
			for j, row := range selectImpl.Rows {
				if row != nil {
					return false, nil, nil, moerr.NewWrongValueCountOnRow(builder.GetContext(), j+1)
				}
			}
		} else {
			colCount := len(insertColumns)
			for j, row := range selectImpl.Rows {
				if len(row) != colCount {
					return false, nil, nil, moerr.NewWrongValueCountOnRow(builder.GetContext(), j+1)
				}
			}
		}

		// example1:insert into a values ();
		// but it does not work at the case:
		// insert into a(a) values (); insert into a values (0),();
		if isAllDefault && stmt.Columns != nil {
			return false, nil, nil, moerr.NewInvalidInput(builder.GetContext(), "insert values does not match the number of columns")
		}
		err = buildValueScan(isAllDefault, info, builder, bindCtx, tableDef, selectImpl, insertColumns, colToIdx, stmt.OnDuplicateUpdate)
		if err != nil {
			return false, nil, nil, err
		}

	case *tree.SelectClause:
		astSelect = stmt.Rows

		subCtx := NewBindContext(builder, bindCtx)
		info.rootId, err = builder.bindSelect(astSelect, subCtx, false)
		if err != nil {
			return false, nil, nil, err
		}
		//ifInsertFromUniqueColMap = make(map[string]bool)

	case *tree.ParenSelect:
		astSelect = selectImpl.Select

		subCtx := NewBindContext(builder, bindCtx)
		info.rootId, err = builder.bindSelect(astSelect, subCtx, false)
		if err != nil {
			return false, nil, nil, err
		}
		// ifInsertFromUniqueColMap = make(map[string]bool)

	default:
		return false, nil, nil, moerr.NewInvalidInput(builder.GetContext(), "insert has unknown select statement")
	}

	if err = builder.addBinding(info.rootId, tree.AliasClause{Alias: derivedTableName}, bindCtx); err != nil {
		return false, nil, nil, err
	}

	lastNode := builder.qry.Nodes[info.rootId]
	if len(insertColumns) != len(lastNode.ProjectList) {
		return false, nil, nil, moerr.NewInvalidInput(builder.GetContext(), "insert values does not match the number of columns")
	}

	var findTableDefFromSource func(node *plan.Node) *plan.TableDef
	findTableDefFromSource = func(node *plan.Node) *plan.TableDef {
		if node == nil {
			return nil
		}
		if node.NodeType == plan.Node_TABLE_SCAN {
			return node.TableDef
		} else {
			if len(node.Children) == 0 {
				return nil
			}
			return findTableDefFromSource(builder.qry.Nodes[node.Children[0]])
		}
	}

	//getUniqueColMap := func(tableDef *plan.TableDef) map[string]struct{} {
	//	result := make(map[string]struct{})
	//	if tableDef != nil {
	//		if tableDef.Pkey != nil {
	//			for _, name := range tableDef.Pkey.Names {
	//				if name != catalog.FakePrimaryKeyColName {
	//					result[name] = struct{}{}
	//				}
	//			}
	//		}
	//		for _, index := range tableDef.Indexes {
	//			if index.Unique {
	//				for _, name := range index.Parts {
	//					result[name] = struct{}{}
	//				}
	//			}
	//		}
	//	}
	//	return result
	//}
	//var fromUniqueCols map[string]struct{}
	//if ifInsertFromUniqueColMap != nil {
	//	tableDef := findTableDefFromSource(lastNode)
	//	fromUniqueCols = getUniqueColMap(tableDef)
	//}

	tag := builder.qry.Nodes[info.rootId].BindingTags[0]
	info.derivedTableId = info.rootId
	oldProject := append([]*Expr{}, lastNode.ProjectList...)

	insertColToExpr := make(map[string]*Expr)
	for i, column := range insertColumns {
		colIdx := colToIdx[column]
		projExpr := &plan.Expr{
			Typ: oldProject[i].Typ,
			Expr: &plan.Expr_Col{
				Col: &plan.ColRef{
					RelPos: tag,
					ColPos: int32(i),
				},
			},
		}
		if tableDef.Cols[colIdx].Typ.Id == int32(types.T_enum) {
			projExpr, err = funcCastForEnumType(builder.GetContext(), projExpr, tableDef.Cols[colIdx].Typ)
			if err != nil {
				return false, nil, nil, err
			}
		} else {
			projExpr, err = forceCastExpr(builder.GetContext(), projExpr, tableDef.Cols[colIdx].Typ)
			if err != nil {
				return false, nil, nil, err
			}
		}
		insertColToExpr[column] = projExpr
		//if ifInsertFromUniqueColMap != nil {
		//	col := lastNode.ProjectList[i].GetCol()
		//	if col != nil {
		//		if _, ok := fromUniqueCols[col.Name]; ok {
		//			ifInsertFromUniqueColMap[column] = true
		//		}
		//	}
		//}
	}

	// create table t(a int, b int unique key);
	// insert into t(a) values (1);  -> isInsertWithoutUniqueKey = true,  then we do not need a plan to insert unique_key_hidden_table;
	// create table t(a int, b int unique key auto_increment)	-> isInsertWithoutUniqueKey is allways false
	// create table t(a int, b int unique key default 10) 		-> isInsertWithoutUniqueKey is allways false
	for _, idx := range tableDef.Indexes {
		if idx.Unique {
			withoutUniqueCol := true
			for _, name := range idx.Parts {
				_, ok := insertColToExpr[name]
				if ok {
					withoutUniqueCol = false
					break
				} else {
					// insert without unique
					// then need check col is not auto_incr & default is not null
					col := tableDef.Cols[tableDef.Name2ColIndex[name]]
					if col.Typ.AutoIncr || (col.Default.Expr != nil && !isNullExpr(col.Default.Expr)) {
						withoutUniqueCol = false
						break
					}
				}
			}
			insertWithoutUniqueKeyMap[idx.IndexName] = withoutUniqueCol
		}
	}

	// have tables : t1(a default 0, b int, pk(a,b)) ,  t2(j int,k int)
	// rewrite 'insert into t1 select * from t2' to
	// select 'select _t.j, _t.k from (select * from t2) _t(j,k)
	// --------
	// rewrite 'insert into t1(b) values (1)' to
	// select 'select 0, _t.column_0 from (select * from values (1)) _t(column_0)
	projectList := make([]*Expr, 0, len(tableDef.Cols))
	isPkCol := make(map[string]bool)
	for _, name := range tableDef.Pkey.Names {
		isPkCol[name] = true
	}
	for i, col := range tableDef.Cols {
		if oldExpr, exists := insertColToExpr[col.Name]; exists {
			projectList = append(projectList, oldExpr)
			// if col.Typ.AutoIncr {
			// if _, ok := pkCols[col.Name]; ok {
			// 	uniqueCheckOnAutoIncr, err = builder.compCtx.GetDbLevelConfig(dbName, "unique_check_on_autoincr")
			// 	if err != nil {
			// 		return false, nil, err
			// 	}
			// 	if uniqueCheckOnAutoIncr == "Error" {
			// 		return false, nil, moerr.NewInvalidInput(builder.GetContext(), "When unique_check_on_autoincr is set to error, insertion of the specified value into auto-incr pk column is not allowed.")
			// 	}
			// }
			// }
		} else {
			defExpr, err := getDefaultExpr(builder.GetContext(), col)
			if err != nil {
				return false, nil, nil, err
			}

			if col.Typ.AutoIncr && isPkCol[col.Name] {
				existAutoPkCol = true
			}

			projectList = append(projectList, defExpr)
		}

		colName2Idx[tableDef.Name+"."+col.Name] = int32(i)
	}

	// append ProjectNode
	projectCtx := NewBindContext(builder, bindCtx)
	lastTag := builder.genNewTag()
	info.rootId = builder.appendNode(&plan.Node{
		NodeType:    plan.Node_PROJECT,
		ProjectList: projectList,
		Children:    []int32{info.rootId},
		BindingTags: []int32{lastTag},
	}, projectCtx)

	info.projectList = make([]*Expr, 0, len(projectList))
	info.derivedTableId = info.rootId
	for i, e := range projectList {
		info.projectList = append(info.projectList, &plan.Expr{
			Typ: e.Typ,
			Expr: &plan.Expr_Col{
				Col: &plan.ColRef{
					RelPos: lastTag,
					ColPos: int32(i),
				},
			},
		})
	}
	info.idx = int32(len(info.projectList))

	return existAutoPkCol, insertWithoutUniqueKeyMap, colName2Idx, nil
}