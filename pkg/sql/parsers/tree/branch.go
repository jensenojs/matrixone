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

type ElseIfStmt struct {
	statementImpl
	Cond Expr
	Body []Statement
}

func NewElseIfStmt(cond Expr, body []Statement, buf *buffer.Buffer) *ElseIfStmt {
	a := buffer.Alloc[ElseIfStmt](buf)
	a.Cond = cond
	a.Body = body
	return a
}

func (node *ElseIfStmt) Format(ctx *FmtCtx) {
	ctx.WriteString("elseif ")
	node.Cond.Format(ctx)
	ctx.WriteString(" then ")
	for _, s := range node.Body {
		if s != nil {
			s.Format(ctx)
			ctx.WriteString("; ")
		}
	}
}
func (node *ElseIfStmt) GetStatementType() string { return "ElseIf Statement" }
func (node *ElseIfStmt) GetQueryType() string     { return QueryTypeTCL }

type IfStmt struct {
	statementImpl
	Cond  Expr
	Body  []Statement
	Elifs []*ElseIfStmt
	Else  []Statement
}

func NewIfStmt(cond Expr, body []Statement, elifs []*ElseIfStmt, elseBody []Statement, buf *buffer.Buffer) *IfStmt {
	a := buffer.Alloc[IfStmt](buf)
	a.Cond = cond
	a.Body = body
	a.Elifs = elifs
	a.Else = elseBody
	return a
}

func (node *IfStmt) Format(ctx *FmtCtx) {
	ctx.WriteString("if ")
	node.Cond.Format(ctx)
	ctx.WriteString(" then ")
	for _, s := range node.Body {
		if s != nil {
			s.Format(ctx)
			ctx.WriteString("; ")
		}
	}
	if node.Elifs != nil {
		for _, elif := range node.Elifs {
			elif.Format(ctx)
		}
	}
	if len(node.Else) != 0 {
		ctx.WriteString("else ")
		for _, s := range node.Else {
			if s != nil {
				s.Format(ctx)
				ctx.WriteString("; ")
			}
		}
	}
	ctx.WriteString("end if")
}

func (node *IfStmt) GetStatementType() string { return "If Statement" }
func (node *IfStmt) GetQueryType() string     { return QueryTypeTCL }

type WhenStmt struct {
	statementImpl
	Cond Expr
	Body []Statement
}

func NewWhenStmt(cond Expr, body []Statement, buf *buffer.Buffer) *WhenStmt {
	a := buffer.Alloc[WhenStmt](buf)
	a.Cond = cond
	a.Body = body
	return a
}

func (node *WhenStmt) Format(ctx *FmtCtx) {
	ctx.WriteString("when ")
	node.Cond.Format(ctx)
	ctx.WriteString(" then ")
	for _, s := range node.Body {
		if s != nil {
			s.Format(ctx)
			ctx.WriteString("; ")
		}
	}
}
func (node *WhenStmt) GetStatementType() string { return "When Statement" }
func (node *WhenStmt) GetQueryType() string     { return QueryTypeTCL }

type CaseStmt struct {
	statementImpl
	Expr  Expr
	Whens []*WhenStmt
	Else  []Statement
}

func NewCaseStmt(expr Expr, whens []*WhenStmt, elseBody []Statement, buf *buffer.Buffer) *CaseStmt {
	a := buffer.Alloc[CaseStmt](buf)
	a.Expr = expr
	a.Whens = whens 
	a.Else = elseBody 
	return a
}

func (node *CaseStmt) Format(ctx *FmtCtx) {
	ctx.WriteString("case ")
	node.Expr.Format(ctx)
	ctx.WriteByte(' ')
	for _, w := range node.Whens {
		w.Format(ctx)
	}
	if node.Else != nil {
		ctx.WriteString("else ")
		for _, s := range node.Else {
			if s != nil {
				s.Format(ctx)
				ctx.WriteString("; ")
			}
		}
	}
	ctx.WriteString("end case")
}

func (node *CaseStmt) GetStatementType() string { return "Case Statement" }
func (node *CaseStmt) GetQueryType() string     { return QueryTypeTCL }
