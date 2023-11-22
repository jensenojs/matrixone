// Copyright 2021 Matrix Origin
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

package tree

import "github.com/matrixorigin/matrixone/pkg/common/buffer"

// Declare statement
type Declare struct {
	statementImpl
	Variables  []string // do NOT reassign after NewDeclare
	ColumnType *T
	DefaultVal Expr
}

func NewDeclare(vs []string, ct *T, dv Expr, buf *buffer.Buffer) *Declare {
	d := buffer.Alloc[Declare](buf)
	d.ColumnType = ct
	d.DefaultVal = dv

	if vs != nil {
		d.Variables = buffer.MakeSlice[string](buf)
		for _, v := range vs {
			d.Variables = buffer.AppendSlice[string](buf, d.Variables, buf.CopyString(v))
		}
	}
	return d
}

func (node *Declare) Format(ctx *FmtCtx) {
	ctx.WriteString("declare ")
	for _, v := range node.Variables {
		ctx.WriteString(v + " ")
	}
	node.ColumnType.InternalType.Format(ctx)
	ctx.WriteString(" default ")
	node.DefaultVal.Format(ctx)
}

func (node *Declare) GetStatementType() string { return "Declare" }
func (node *Declare) GetQueryType() string     { return QueryTypeOth }
