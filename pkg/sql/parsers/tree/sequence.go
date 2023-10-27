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

import (
	"fmt"

	"github.com/matrixorigin/matrixone/pkg/common/buffer"
)

type CreateSequence struct {
	statementImpl

	Name        *TableName
	Type        ResolvableTypeReference
	IfNotExists bool
	IncrementBy *IncrementByOption
	MinValue    *MinValueOption
	MaxValue    *MaxValueOption
	StartWith   *StartWithOption
	Cycle       bool
}

func NewCreateSequence(ifs bool, name *TableName, typ ResolvableTypeReference, incrementBy *IncrementByOption, min *MinValueOption, max *MaxValueOption, start *StartWithOption, cycle bool, buf *buffer.Buffer) *CreateSequence {
	cr := buffer.Alloc[CreateSequence](buf)
	cr.IfNotExists = ifs
	cr.Name = name
	cr.Type = typ
	cr.IncrementBy = incrementBy
	cr.MinValue = min
	cr.MaxValue = max
	cr.StartWith = start
	cr.Cycle = cycle
	return cr
}

func (node *CreateSequence) Format(ctx *FmtCtx) {
	ctx.WriteString("create sequence ")

	if node.IfNotExists {
		ctx.WriteString("if not exists ")
	}

	node.Name.Format(ctx)

	ctx.WriteString(" as ")
	node.Type.(*T).InternalType.Format(ctx)
	ctx.WriteString(" ")
	if node.IncrementBy != nil {
		node.IncrementBy.Format(ctx)
	}
	if node.MinValue != nil {
		node.MinValue.Format(ctx)
	}
	if node.MaxValue != nil {
		node.MaxValue.Format(ctx)
	}
	if node.StartWith != nil {
		node.StartWith.Format(ctx)
	}
	if node.Cycle {
		ctx.WriteString("cycle")
	} else {
		ctx.WriteString("no cycle")
	}
}

func (node *CreateSequence) GetStatementType() string { return "Create Sequence" }
func (node *CreateSequence) GetQueryType() string     { return QueryTypeDDL }

type IncrementByOption struct {
	Minus bool
	Num   any
}

func NewIncrementByOption(m bool, n any, buf *buffer.Buffer) *IncrementByOption {
	s := buffer.Alloc[IncrementByOption](buf)
	s.Minus = m
	s.Num = n
	return s
}

func (node *IncrementByOption) Format(ctx *FmtCtx) {
	ctx.WriteString("increment by ")
	formatAny(node.Minus, node.Num, ctx)
}

type MinValueOption struct {
	Minus bool
	Num   any
}

func NewMinValueOption(m bool, n any, buf *buffer.Buffer) *MinValueOption {
	s := buffer.Alloc[MinValueOption](buf)
	s.Minus = m
	s.Num = n
	return s
}

func (node *MinValueOption) Format(ctx *FmtCtx) {
	ctx.WriteString("minvalue ")
	formatAny(node.Minus, node.Num, ctx)
}

type MaxValueOption struct {
	Minus bool
	Num   any
}

func NewMaxValueOption(m bool, n any, buf *buffer.Buffer) *MaxValueOption {
	s := buffer.Alloc[MaxValueOption](buf)
	s.Minus = m
	s.Num = n
	return s
}

func (node *MaxValueOption) Format(ctx *FmtCtx) {
	ctx.WriteString("maxvalue ")
	formatAny(node.Minus, node.Num, ctx)
}

type StartWithOption struct {
	Minus bool
	Num   any
}

func NewStartWithOption(m bool, n any, buf *buffer.Buffer) *StartWithOption {
	s := buffer.Alloc[StartWithOption](buf)
	s.Minus = m
	s.Num = n
	return s
}

func (node *StartWithOption) Format(ctx *FmtCtx) {
	ctx.WriteString("start with ")
	formatAny(node.Minus, node.Num, ctx)
}

type CycleOption struct {
	Cycle bool
}

func NewCycleOption(cycle bool, buf *buffer.Buffer) *CycleOption {
	c := buffer.Alloc[CycleOption](buf)
	c.Cycle = cycle
	return c
}

func (node *CycleOption) Format(ctx *FmtCtx) {
	if node.Cycle {
		ctx.WriteString("cycle")
	} else {
		ctx.WriteString("no cycle")
	}
}

type TypeOption struct {
	Type ResolvableTypeReference
}

func NewTypeOption(t ResolvableTypeReference, buf *buffer.Buffer) *TypeOption {
	typ := buffer.Alloc[TypeOption](buf)
	typ.Type = t
	return typ
}

func (node *TypeOption) Format(ctx *FmtCtx) {
	ctx.WriteString(" as ")
	node.Type.(*T).InternalType.Format(ctx)
}

func formatAny(minus bool, num any, ctx *FmtCtx) {
	switch num := num.(type) {
	case uint64:
		ctx.WriteString(fmt.Sprintf("%v ", num))
	case int64:
		var v int64
		if minus {
			v = -num
		} else {
			v = num
		}
		ctx.WriteString(fmt.Sprintf("%v ", v))
	}
}

type DropSequence struct {
	statementImpl
	IfExists bool
	Names    TableNames
}

func NewDropSequence(ife bool, n TableNames, buf *buffer.Buffer) *DropSequence {
	dr := buffer.Alloc[DropSequence](buf)
	dr.IfExists = ife
	dr.Names = n
	return dr
}

func (node *DropSequence) Format(ctx *FmtCtx) {
	ctx.WriteString("drop sequence")
	if node.IfExists {
		ctx.WriteString(" if exists")
	}
	ctx.WriteByte(' ')
	node.Names.Format(ctx)
}

func (node *DropSequence) GetStatementType() string { return "Drop Sequence" }
func (node *DropSequence) GetQueryType() string     { return QueryTypeDDL }

type AlterSequence struct {
	statementImpl

	Name        *TableName
	Type        *TypeOption
	IfExists    bool
	IncrementBy *IncrementByOption
	MinValue    *MinValueOption
	MaxValue    *MaxValueOption
	StartWith   *StartWithOption
	Cycle       *CycleOption
}

func NewAlterSequence(ifs bool, name *TableName, typ *TypeOption, incBy *IncrementByOption, min *MinValueOption, max *MaxValueOption, start *StartWithOption, cy *CycleOption, buf *buffer.Buffer) *AlterSequence {
	al := buffer.Alloc[AlterSequence](buf)
	al.IfExists = ifs
	al.Name = name
	al.Type = typ
	al.IncrementBy = incBy
	al.MinValue = min
	al.MaxValue = max
	al.StartWith = start
	al.Cycle = cy
	return al
}

func (node *AlterSequence) Format(ctx *FmtCtx) {
	ctx.WriteString("alter sequence ")

	if node.IfExists {
		ctx.WriteString("if exists ")
	}

	node.Name.Format(ctx)

	if node.Type != nil {
		node.Type.Format(ctx)
	}
	ctx.WriteString(" ")
	if node.IncrementBy != nil {
		node.IncrementBy.Format(ctx)
	}
	if node.MinValue != nil {
		node.MinValue.Format(ctx)
	}
	if node.MaxValue != nil {
		node.MaxValue.Format(ctx)
	}
	if node.StartWith != nil {
		node.StartWith.Format(ctx)
	}
	if node.Cycle != nil {
		node.Cycle.Format(ctx)
	}
}

func (node *AlterSequence) GetStatementType() string { return "Alter Sequence" }
func (node *AlterSequence) GetQueryType() string     { return QueryTypeDDL }
