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

type RepeatStmt struct {
	statementImpl
	Name *BufIdentifier
	Body []Statement
	Cond Expr
}

func NewRepeatStmt(name Identifier, b []Statement, c Expr, buf *buffer.Buffer) *RepeatStmt {
	r := buffer.Alloc[RepeatStmt](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)

	r.Name = n
	r.Cond = c
	r.Body = b
	return r
}

func (node *RepeatStmt) Format(ctx *FmtCtx) {
	n := node.Name.Get()
	if n != "" {
		ctx.WriteString(string(n))
		ctx.WriteString(": ")
	}
	ctx.WriteString("repeat ")
	for _, s := range node.Body {
		if s != nil {
			s.Format(ctx)
			ctx.WriteString("; ")
		}
	}
	ctx.WriteString("until ")
	node.Cond.Format(ctx)
	ctx.WriteString(" end repeat")
	if n != "" {
		ctx.WriteString(string(n))
	}
}

func (node *RepeatStmt) GetStatementType() string { return "Repeat Statement" }
func (node *RepeatStmt) GetQueryType() string     { return QueryTypeTCL }

type WhileStmt struct {
	statementImpl
	Name *BufIdentifier
	Cond Expr
	Body []Statement
}

func NewWhileStmt(name Identifier, c Expr, b []Statement, buf *buffer.Buffer) *WhileStmt {
	l := buffer.Alloc[WhileStmt](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)
	l.Name = n

	l.Cond = c
	l.Body = b
	return l
}

func (node *WhileStmt) Format(ctx *FmtCtx) {
	n := node.Name.Get()
	if n != "" {
		ctx.WriteString(string(n))
		ctx.WriteString(": ")
	}
	ctx.WriteString("while ")
	node.Cond.Format(ctx)
	ctx.WriteString(" do ")
	for _, s := range node.Body {
		if s != nil {
			s.Format(ctx)
			ctx.WriteString("; ")
		}
	}
	ctx.WriteString("end while")
	if n != "" {
		ctx.WriteByte(' ')
		ctx.WriteString(string(n))
	}
}

func (node *WhileStmt) GetStatementType() string { return "While Statement" }
func (node *WhileStmt) GetQueryType() string     { return QueryTypeTCL }

type LoopStmt struct {
	statementImpl
	Name *BufIdentifier
	Body []Statement
}

func NewLoopStmt(name Identifier, b []Statement, buf *buffer.Buffer) *LoopStmt {
	l := buffer.Alloc[LoopStmt](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)

	l.Name = n
	l.Body = b
	return l
}

func (node *LoopStmt) Format(ctx *FmtCtx) {
	n := node.Name.Get()
	if n != "" {
		ctx.WriteString(string(n))
		ctx.WriteString(": ")
	}
	ctx.WriteString("loop ")
	for _, s := range node.Body {
		if s != nil {
			s.Format(ctx)
			ctx.WriteString("; ")
		}
	}
	ctx.WriteString("end loop")
	if n != "" {
		ctx.WriteByte(' ')
		ctx.WriteString(string(n))
	}
}

func (node *LoopStmt) GetStatementType() string { return "Loop Statement" }
func (node *LoopStmt) GetQueryType() string     { return QueryTypeTCL }

type IterateStmt struct {
	statementImpl
	Name *BufIdentifier
}

func NewIterateStmt(name Identifier, buf *buffer.Buffer) *IterateStmt {
	l := buffer.Alloc[IterateStmt](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)

	l.Name = n
	return l
}

func (node *IterateStmt) Format(ctx *FmtCtx) {
	ctx.WriteString("iterate ")
	ctx.WriteString(string(node.Name.Get()))
}

func (node *IterateStmt) GetStatementType() string { return "Iterate Statement" }
func (node *IterateStmt) GetQueryType() string     { return QueryTypeTCL }

type LeaveStmt struct {
	statementImpl
	Name *BufIdentifier
}

func NewLeaveStmt(name Identifier, buf *buffer.Buffer) *LeaveStmt {
	l := buffer.Alloc[LeaveStmt](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)

	l.Name = n
	return l
}

func (node *LeaveStmt) Format(ctx *FmtCtx) {
	ctx.WriteString("leave ")
	ctx.WriteString(string(node.Name.Get()))
}

func (node *LeaveStmt) GetStatementType() string { return "Leave Statement" }
func (node *LeaveStmt) GetQueryType() string     { return QueryTypeTCL }
