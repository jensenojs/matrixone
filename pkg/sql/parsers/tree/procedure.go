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

type InOutArgType int

const (
	TYPE_IN InOutArgType = iota
	TYPE_OUT
	TYPE_INOUT
)

type ProcedureArgType struct {
	Type InOutArgType
}

type ProcedureArg interface {
	NodeFormatter
	GetName(ctx *FmtCtx) string
	GetType() int
}

type ProcedureArgImpl struct {
	ProcedureArg
}

// container holding list of arguments in udf
type ProcedureArgs []ProcedureArg

type ProcedureArgDecl struct {
	ProcedureArgImpl
	Name      *UnresolvedName
	Type      ResolvableTypeReference
	InOutType InOutArgType
}

type ProcedureArgForMarshal struct {
	Name      *UnresolvedName
	Type      ResolvableTypeReference
	InOutType InOutArgType
}

type ProcedureName struct {
	Name objName
}

func NewProcedureName(name Identifier, prefix ObjectNamePrefix, buf *buffer.Buffer) *ProcedureName {
	pn := buffer.Alloc[ProcedureName](buf)
	pn.Name.ObjectName = name
	pn.Name.ObjectNamePrefix = prefix
	return pn
}

func NewProcedureArgDecl(inOutType InOutArgType, name *UnresolvedName, typ ResolvableTypeReference, buf *buffer.Buffer) *ProcedureArgDecl {
	pad := buffer.Alloc[ProcedureArgDecl](buf)
	pad.Name = name
	pad.Type = typ
	pad.InOutType = inOutType
	return pad
}

func (node *ProcedureArgDecl) Format(ctx *FmtCtx) {
	// in out type
	switch node.InOutType {
	case TYPE_IN:
		ctx.WriteString("in ")
	case TYPE_OUT:
		ctx.WriteString("out ")
	case TYPE_INOUT:
		ctx.WriteString("inout ")
	}
	if node.Name != nil {
		node.Name.Format(ctx)
		ctx.WriteByte(' ')
	}
	node.Type.(*T).InternalType.Format(ctx)
}

func (node *ProcedureArgDecl) GetName(ctx *FmtCtx) string {
	node.Name.Format(ctx)
	return ctx.String()
}

func (node *ProcedureArgDecl) GetType() int {
	return int(node.InOutType)
}

func (node *ProcedureName) Format(ctx *FmtCtx) {
	if node.Name.ExplicitCatalog {
		ctx.WriteString(string(node.Name.CatalogName))
		ctx.WriteByte('.')
	}
	if node.Name.ExplicitSchema {
		ctx.WriteString(string(node.Name.SchemaName))
		ctx.WriteByte('.')
	}
	ctx.WriteString(string(node.Name.ObjectName))
}

func (node *ProcedureName) HasNoNameQualifier() bool {
	return !node.Name.ExplicitCatalog && !node.Name.ExplicitSchema
}

type CreateProcedure struct {
	statementImpl
	Name *ProcedureName
	Args ProcedureArgs
	Body string
}

func NewCreateProcedure(n *ProcedureName, args ProcedureArgs, body string, buf *buffer.Buffer) *CreateProcedure {
	drop := buffer.Alloc[CreateProcedure](buf)	
	drop.Name = n
	drop.Args = args
	drop.Body = body
	return drop
}

type DropProcedure struct {
	statementImpl
	Name     *ProcedureName
	IfExists bool
}

func NewDropProcedure(n *ProcedureName, ifs bool, buf *buffer.Buffer) *DropProcedure {
	drop := buffer.Alloc[DropProcedure](buf)	
	drop.Name = n
	drop.IfExists = ifs
	return drop
}

func (node *CreateProcedure) Format(ctx *FmtCtx) {
	ctx.WriteString("create procedure ")

	node.Name.Format(ctx)

	ctx.WriteString(" (")
	for i, def := range node.Args {
		if i != 0 {
			ctx.WriteString(",")
			ctx.WriteByte(' ')
		}
		def.Format(ctx)
	}
	ctx.WriteString(") '")

	ctx.WriteString(node.Body)
	ctx.WriteString("'")
}

func (node *DropProcedure) Format(ctx *FmtCtx) {
	ctx.WriteString("drop procedure ")
	if node.IfExists {
		ctx.WriteString("if exists ")
	}
	node.Name.Format(ctx)
}

func (node *CreateProcedure) GetStatementType() string { return "Create Procedure" }
func (node *CreateProcedure) GetQueryType() string     { return QueryTypeOth }
func (node *DropProcedure) GetStatementType() string   { return "Create Procedure" }
func (node *DropProcedure) GetQueryType() string       { return QueryTypeOth }

type CallStmt struct {
	statementImpl
	Name *ProcedureName
	Args Exprs
}

func NewCallStmt(n *ProcedureName, a Exprs, buf *buffer.Buffer) *CallStmt {
	c := buffer.Alloc[CallStmt](buf)
	c.Name = n
	c.Args = a
	return c
}

func (node *CallStmt) Format(ctx *FmtCtx) {
	ctx.WriteString("call ")
	node.Name.Format(ctx)
	ctx.WriteString("(")
	if len(node.Args) != 0 {
		node.Args.Format(ctx)
	}
	ctx.WriteString(")")
}

func (node *CallStmt) GetStatementType() string { return "Call" }
func (node *CallStmt) GetQueryType() string     { return QueryTypeOth }
