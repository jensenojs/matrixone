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

package tree

import "github.com/matrixorigin/matrixone/pkg/common/buffer"

type Prepare interface {
	Statement
}

type prepareImpl struct {
	Prepare
	Format string
}

type PrepareStmt struct {
	prepareImpl
	Name Identifier
	Stmt Statement
}

type PrepareString struct {
	prepareImpl
	Name Identifier
	Sql  string
}

func (node *PrepareStmt) Format(ctx *FmtCtx) {
	ctx.WriteString("prepare ")
	node.Name.Format(ctx)
	ctx.WriteString(" from ")
	node.Stmt.Format(ctx)
}

func (node *PrepareString) Format(ctx *FmtCtx) {
	ctx.WriteString("prepare ")
	node.Name.Format(ctx)
	ctx.WriteString(" from ")
	ctx.WriteString(node.Sql)
}

func (node *PrepareStmt) GetStatementType() string   { return "Prepare" }
func (node *PrepareStmt) GetQueryType() string       { return QueryTypeOth }
func (node *PrepareString) GetStatementType() string { return "Prepare" }
func (node *PrepareString) GetQueryType() string     { return QueryTypeOth }

func NewPrepareStmt(name Identifier, stmt Statement, buf *buffer.Buffer) *PrepareStmt {
	ps := buffer.Alloc[PrepareStmt](buf)
	ps.Name = name
	ps.Stmt = stmt
	return ps
}

func NewPrepareString(name Identifier, sql string, buf *buffer.Buffer) *PrepareString {
	ps := buffer.Alloc[PrepareString](buf)
	ps.Name = name
	ps.Sql = sql
	return ps
}
