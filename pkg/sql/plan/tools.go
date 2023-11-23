// Copyright 2021 - 2022 Matrix Origin
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

package plan

import (
	"go/constant"

	"github.com/matrixorigin/matrixone/pkg/common/buffer"
	"github.com/matrixorigin/matrixone/pkg/sql/parsers/tree"
)

const (
	moRecursiveLevelCol            = "__mo_recursive_level_col"
	moDefaultRecursionMax          = 100
	moCheckRecursionLevelFun       = "mo_check_level"
	moEnumCastIndexToValueFun      = "cast_index_to_value"
	moEnumCastValueToIndexFun      = "cast_value_to_index"
	moEnumCastIndexValueToIndexFun = "cast_index_value_to_index"
)

func makeZeroRecursiveLevel(buf *buffer.Buffer) *tree.SelectExpr {
	sel := buffer.Alloc[tree.SelectExpr](buf)
	sel.Expr = tree.NewNumValWithType(constant.MakeInt64(0), "0", false, tree.P_int64, buf)
	sel.As = tree.NewCStr(moRecursiveLevelCol, 1, buf)
	return sel

}

func makePlusRecursiveLevel(name string, buf *buffer.Buffer) *tree.SelectExpr {
	sel := buffer.Alloc[tree.SelectExpr](buf)
	a := tree.SetUnresolvedName(buf, name, moRecursiveLevelCol)
	b := tree.NewNumValWithType(constant.MakeInt64(1), "1", false, tree.P_int64, buf)
	expr := tree.NewBinaryExpr(tree.PLUS, a, b, buf)
	sel.Expr = expr
	sel.As = tree.NewCStr("", 1, buf)
	return sel
}
