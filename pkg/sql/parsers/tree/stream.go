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

type CreateStream struct {
	statementImpl
	Replace     bool
	Source      bool
	IfNotExists bool
	StreamName  *TableName
	Defs        TableDefs
	ColNames    *BufIdentifierList
	AsSource    *Select
	Options     []TableOption
}

func NewCreateStream(replace, source, ifNotExists bool, streamName *TableName, defs TableDefs, colNames IdentifierList, asSource *Select, options []TableOption, buf *buffer.Buffer) *CreateStream {
	c := buffer.Alloc[CreateStream](buf)
	c.Replace = replace
	c.Source = source
	c.IfNotExists = ifNotExists
	c.StreamName = streamName
	c.Defs = defs
	bc := NewBufIdentifierList(colNames)
	buf.Pin(bc)
	c.ColNames = bc
	c.AsSource = asSource
	c.Options = options
	return c
}

func (node *CreateStream) Format(ctx *FmtCtx) {
	ctx.WriteString("create")
	if node.Replace {
		ctx.WriteString(" or replace")
	}
	if node.Defs != nil {
		if node.Source {
			ctx.WriteString(" source")
		}
		ctx.WriteString(" stream")
		if node.IfNotExists {
			ctx.WriteString(" if not exists")
		}
		ctx.WriteByte(' ')
		node.StreamName.Format(ctx)

		ctx.WriteString(" (")
		for i, def := range node.Defs {
			if i != 0 {
				ctx.WriteString(",")
				ctx.WriteByte(' ')
			}
			def.Format(ctx)
		}
		ctx.WriteByte(')')

		if node.Options != nil {
			prefix := " with ("
			for _, t := range node.Options {
				ctx.WriteString(prefix)
				t.Format(ctx)
				prefix = ", "
			}
			ctx.WriteByte(')')
		}
		return
	}
	ctx.WriteString(" stream")
	if node.IfNotExists {
		ctx.WriteString(" if not exists")
	}
	ctx.WriteByte(' ')
	node.StreamName.Format(ctx)
	if node.Options != nil {
		prefix := " with ("
		for _, t := range node.Options {
			ctx.WriteString(prefix)
			t.Format(ctx)
			prefix = ", "
		}
		ctx.WriteByte(')')
	}
	ctx.WriteString(" as ")
	node.AsSource.Format(ctx)
}

func (node *CreateStream) GetStatementType() string { return "Create Stream" }
func (node *CreateStream) GetQueryType() string     { return QueryTypeDDL }

type CreateStreamWithOption struct {
	createOptionImpl
	Key Identifier
	Val Expr
}

func (node *CreateStreamWithOption) Format(ctx *FmtCtx) {
	ctx.WriteString(string(node.Key))
	ctx.WriteString(" = ")
	node.Val.Format(ctx)
}

func NewCreateStreamWithOption(k Identifier, v Expr, buf *buffer.Buffer) *CreateStreamWithOption {
	c := buffer.Alloc[CreateStreamWithOption](buf)
	c.Key = k
	c.Val = v
	return c
}

type AttributeHeader struct {
	columnAttributeImpl
	Key *BufString
}

func (node *AttributeHeader) Format(ctx *FmtCtx) {
	ctx.WriteString("header(")
	ctx.WriteString(node.Key.Get())
	ctx.WriteByte(')')
}

func NewAttributeHeader(key string, buf *buffer.Buffer) *AttributeHeader {
	ah := buffer.Alloc[AttributeHeader](buf)
	bkey := NewBufString(key)
	buf.Pin(bkey)
	ah.Key = bkey
	return ah
}

type AttributeHeaders struct {
	columnAttributeImpl
}

func (node *AttributeHeaders) Format(ctx *FmtCtx) {
	ctx.WriteString("headers")
}

func NewAttributeHeaders(buf *buffer.Buffer) *AttributeHeaders {
	aHeads := buffer.Alloc[AttributeHeaders](buf)
	return aHeads
}
