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

// Begin statement
type BeginCompound struct {
	statementImpl
}

type EndCompound struct {
	statementImpl
}

type CompoundStmt struct {
	statementImpl
	Stmts []Statement
}

func NewCompoundStmt(s []Statement, buf *buffer.Buffer) *CompoundStmt {
	c := buffer.Alloc[CompoundStmt](buf)
	c.Stmts = s
	return c
}

func NewBeginCompound(buf *buffer.Buffer) *BeginCompound {
	b := buffer.Alloc[BeginCompound](buf)
	return b
}

func NewEndCompound(buf *buffer.Buffer) *EndCompound {
	e := buffer.Alloc[EndCompound](buf)
	return e
}

func (node *CompoundStmt) Format(ctx *FmtCtx) {
	ctx.WriteString("begin ")
	for _, s := range node.Stmts {
		if s != nil {
			s.Format(ctx)
			ctx.WriteString("; ")
		}
	}
	ctx.WriteString("end")
}

func (node *CompoundStmt) GetStatementType() string { return "compound" }
func (node *CompoundStmt) GetQueryType() string     { return QueryTypeTCL }

func (node *BeginCompound) Format(ctx *FmtCtx) {
	ctx.WriteString("begin")
}

func (node *BeginCompound) GetStatementType() string { return "begin" }
func (node *BeginCompound) GetQueryType() string     { return QueryTypeTCL }

func (node *EndCompound) Format(ctx *FmtCtx) {
	ctx.WriteString("end")
}

func (node *EndCompound) GetStatementType() string { return "end" }
func (node *EndCompound) GetQueryType() string     { return QueryTypeTCL }
