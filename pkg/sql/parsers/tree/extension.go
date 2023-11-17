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

type CreateExtension struct {
	statementImpl
	Language string
	Name     *BufIdentifier
	Filename *BufIdentifier
}

func NewCreateExtension(language string, name, filename Identifier, buf *buffer.Buffer) *CreateExtension {
	c := buffer.Alloc[CreateExtension](buf)
	c.Language = language
	n := NewBufIdentifier(name)
	f := NewBufIdentifier(filename)
	buf.Pin(n, f)

	c.Name = n
	c.Filename = f
	return c
}

type LoadExtension struct {
	statementImpl
	Name *BufIdentifier
}

func NewLoadExtension(name Identifier, buf *buffer.Buffer) *LoadExtension {
	l := buffer.Alloc[LoadExtension](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)
	l.Name = n
	return l
}

func (node *CreateExtension) Format(ctx *FmtCtx) {
	ctx.WriteString("create extension ")
	ctx.WriteString(node.Language)
	ctx.WriteString(" as ")
	node.Name.Format(ctx)
	ctx.WriteString(" file ")
	node.Filename.Format(ctx)
}

func (node *LoadExtension) Format(ctx *FmtCtx) {
	ctx.WriteString("load ")
	node.Name.Format(ctx)
}
