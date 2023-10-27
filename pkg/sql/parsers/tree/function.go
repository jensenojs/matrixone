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

type FunctionArg interface {
	NodeFormatter
	Expr
	GetName(ctx *FmtCtx) string
	GetType(ctx *FmtCtx) string
}

type FunctionArgImpl struct {
	FunctionArg
}

// container holding list of arguments in udf
type FunctionArgs []FunctionArg

type FunctionArgDecl struct {
	FunctionArgImpl
	Name       *UnresolvedName
	Type       ResolvableTypeReference
	DefaultVal Expr
}

type ReturnType struct {
	Type ResolvableTypeReference
}

func (node *FunctionArgDecl) Format(ctx *FmtCtx) {
	if node.Name != nil {
		node.Name.Format(ctx)
		ctx.WriteByte(' ')
	}
	node.Type.(*T).InternalType.Format(ctx)
	if node.DefaultVal != nil {
		ctx.WriteString(" default ")
		ctx.PrintExpr(node, node.DefaultVal, true)
	}
}

func (node *FunctionArgDecl) GetName(ctx *FmtCtx) string {
	node.Name.Format(ctx)
	return ctx.String()
}

func (node *FunctionArgDecl) GetType(ctx *FmtCtx) string {
	node.Type.(*T).InternalType.Format(ctx)
	return ctx.String()
}

func (node *ReturnType) Format(ctx *FmtCtx) {
	node.Type.(*T).InternalType.Format(ctx)
}

type FunctionName struct {
	Name objName
}

type CreateFunction struct {
	statementImpl
	Name       *FunctionName
	Args       FunctionArgs
	ReturnType *ReturnType
	Body       string
	Language   string
}

func NewCreateFunction(n *FunctionName, args FunctionArgs, returnTyp *ReturnType, language , body string, buf *buffer.Buffer) *CreateFunction {
	c := buffer.Alloc[CreateFunction](buf)	
	c.Name = n
	c.Args = args
	c.ReturnType = returnTyp
	c.Body = body
	c.Language = language
	return c
}



type DropFunction struct {
	statementImpl
	Name *FunctionName
	Args FunctionArgs
}

func NewDropFunction(n *FunctionName, args FunctionArgs, buf *buffer.Buffer) *DropFunction {
	drop := buffer.Alloc[DropFunction](buf)	
	drop.Name = n
	drop.Args = args
	return drop
}

func (node *FunctionName) Format(ctx *FmtCtx) {
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

func (node *FunctionName) HasNoNameQualifier() bool {
	return !node.Name.ExplicitCatalog && !node.Name.ExplicitSchema
}

func (node *CreateFunction) Format(ctx *FmtCtx) {
	ctx.WriteString("create function ")

	node.Name.Format(ctx)

	ctx.WriteString(" (")

	for i, def := range node.Args {
		if i != 0 {
			ctx.WriteString(",")
			ctx.WriteByte(' ')
		}
		def.Format(ctx)
	}

	ctx.WriteString(")")
	ctx.WriteString(" returns ")

	node.ReturnType.Format(ctx)

	ctx.WriteString(" language ")
	ctx.WriteString(node.Language)

	ctx.WriteString(" as '")

	ctx.WriteString(node.Body)
	ctx.WriteString("'")
}

func (node *DropFunction) Format(ctx *FmtCtx) {
	ctx.WriteString("drop function ")
	node.Name.Format(ctx)
	ctx.WriteString(" (")

	for i, def := range node.Args {
		if i != 0 {
			ctx.WriteString(",")
			ctx.WriteByte(' ')
		}
		def.Format(ctx)
	}

	ctx.WriteString(")")
}

func NewFunctionArgDecl(name *UnresolvedName, typ ResolvableTypeReference, defaultVal Expr, buf *buffer.Buffer) *FunctionArgDecl {
	fad := buffer.Alloc[FunctionArgDecl](buf)
	fad.Name = name
	fad.Type = typ
	fad.DefaultVal = defaultVal
	return fad
}

func NewFuncName(name Identifier, prefix ObjectNamePrefix, buf *buffer.Buffer) *FunctionName {
	fn := buffer.Alloc[FunctionName](buf)
	fn.Name.ObjectName = name
	fn.Name.ObjectNamePrefix = prefix
	return fn
}

func NewReturnType(typ ResolvableTypeReference, buf *buffer.Buffer) *ReturnType {
	rt := buffer.Alloc[ReturnType](buf)
	rt.Type = typ
	return rt
}

func (node *CreateFunction) GetStatementType() string { return "CreateFunction" }
func (node *CreateFunction) GetQueryType() string     { return QueryTypeDDL }

func (node *DropFunction) GetStatementType() string { return "DropFunction" }
func (node *DropFunction) GetQueryType() string     { return QueryTypeDDL }
