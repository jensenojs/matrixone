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

// DROP Database statement
type DropDatabase struct {
	statementImpl
	Name     *BufIdentifier
	IfExists bool
}

func (node *DropDatabase) Format(ctx *FmtCtx) {
	ctx.WriteString("drop database")
	if node.IfExists {
		ctx.WriteByte(' ')
		ctx.WriteString("if exists")
	}
	if node.Name.Get() != "" {
		ctx.WriteByte(' ')
		ctx.WriteString(string(node.Name.Get()))
	}
}

func (node *DropDatabase) GetStatementType() string { return "Drop Database" }
func (node *DropDatabase) GetQueryType() string     { return QueryTypeDDL }

func NewDropDatabase(name Identifier, ifExists bool, buf *buffer.Buffer) *DropDatabase {
	d := buffer.Alloc[DropDatabase](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)

	d.Name = n
	d.IfExists = ifExists
	return d
}

// DROP Table statement
type DropTable struct {
	statementImpl
	IfExists bool
	Names    TableNames
}

func (node *DropTable) Format(ctx *FmtCtx) {
	ctx.WriteString("drop table")
	if node.IfExists {
		ctx.WriteString(" if exists")
	}
	ctx.WriteByte(' ')
	node.Names.Format(ctx)
}

func (node *DropTable) GetStatementType() string { return "Drop Table" }
func (node *DropTable) GetQueryType() string     { return QueryTypeDDL }

func NewDropTable(ifExists bool, names TableNames, buf *buffer.Buffer) *DropTable {
	d := buffer.Alloc[DropTable](buf)
	d.IfExists = ifExists
	d.Names = names
	return d
}

// DropView DROP View statement
type DropView struct {
	statementImpl
	IfExists bool
	Names    TableNames
}

func (node *DropView) Format(ctx *FmtCtx) {
	ctx.WriteString("drop view")
	if node.IfExists {
		ctx.WriteString(" if exists")
	}
	ctx.WriteByte(' ')
	node.Names.Format(ctx)
}

func (node *DropView) GetStatementType() string { return "Drop View" }
func (node *DropView) GetQueryType() string     { return QueryTypeDDL }

func NewDropView(ifExists bool, names TableNames, buf *buffer.Buffer) *DropView {
	d := buffer.Alloc[DropView](buf)
	d.IfExists = ifExists
	d.Names = names
	return d
}

type DropIndex struct {
	statementImpl
	Name       *BufIdentifier
	TableName  *TableName
	IfExists   bool
	MiscOption []MiscOption
}

func (node *DropIndex) Format(ctx *FmtCtx) {
	ctx.WriteString("drop index")
	if node.IfExists {
		ctx.WriteString(" if exists")
	}
	ctx.WriteByte(' ')
	ctx.WriteString(string(node.Name.Get()))

	ctx.WriteString(" on ")
	node.TableName.Format(ctx)
}

func (node *DropIndex) GetStatementType() string { return "Drop Index" }
func (node *DropIndex) GetQueryType() string     { return QueryTypeDDL }

func NewDropIndex(name Identifier, tableName *TableName, ifExists bool, buf *buffer.Buffer) *DropIndex {
	d := buffer.Alloc[DropIndex](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)
	d.Name = n
	d.TableName = tableName
	d.IfExists = ifExists
	return d
}

type DropRole struct {
	statementImpl
	IfExists bool
	Roles    []*Role
}

func (node *DropRole) Format(ctx *FmtCtx) {
	ctx.WriteString("drop role")
	if node.IfExists {
		ctx.WriteString(" if exists")
	}
	prefix := " "
	for _, r := range node.Roles {
		ctx.WriteString(prefix)
		r.Format(ctx)
		prefix = ", "
	}
}

func (node *DropRole) GetStatementType() string { return "Drop Role" }
func (node *DropRole) GetQueryType() string     { return QueryTypeDCL }

func NewDropRole(ifExists bool, roles []*Role, buf *buffer.Buffer) *DropRole {
	d := buffer.Alloc[DropRole](buf)
	d.IfExists = ifExists
	d.Roles = roles
	return d
}

type DropUser struct {
	statementImpl
	IfExists bool
	Users    []*User
}

func (node *DropUser) Format(ctx *FmtCtx) {
	ctx.WriteString("drop user")
	if node.IfExists {
		ctx.WriteString(" if exists")
	}
	prefix := " "
	for _, u := range node.Users {
		ctx.WriteString(prefix)
		u.Format(ctx)
		prefix = ", "
	}
}

func (node *DropUser) GetStatementType() string { return "Drop User" }
func (node *DropUser) GetQueryType() string     { return QueryTypeDCL }

func NewDropUser(ifExists bool, users []*User, buf *buffer.Buffer) *DropUser {
	d := buffer.Alloc[DropUser](buf)
	d.IfExists = ifExists
	d.Users = users
	return d
}

type DropAccount struct {
	statementImpl
	IfExists bool
	Name     *BufString
}

func NewDropAccount(i bool, n string, buf *buffer.Buffer) *DropAccount {
	d := buffer.Alloc[DropAccount](buf)
	d.IfExists = i
	bn := NewBufString(n)
	buf.Pin(bn)
	d.Name = bn
	return d
}

func (node *DropAccount) Format(ctx *FmtCtx) {
	ctx.WriteString("drop account")
	if node.IfExists {
		ctx.WriteString(" if exists")
	}
	ctx.WriteString(" ")
	ctx.WriteString(node.Name.Get())
}

func (node *DropAccount) GetStatementType() string { return "Drop Account" }
func (node *DropAccount) GetQueryType() string     { return QueryTypeDCL }

type DropPublication struct {
	statementImpl
	Name     *BufIdentifier
	IfExists bool
}

func NewDropPublication(i bool, name Identifier, buf *buffer.Buffer) *DropPublication {
	d := buffer.Alloc[DropPublication](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)
	d.IfExists = i
	d.Name = n
	return d
}

func (node *DropPublication) Format(ctx *FmtCtx) {
	ctx.WriteString("drop publication")
	if node.IfExists {
		ctx.WriteString(" if exists")
	}
	ctx.WriteByte(' ')
	node.Name.Format(ctx)
}

func (node *DropPublication) GetStatementType() string { return "Drop Publication" }
func (node *DropPublication) GetQueryType() string     { return QueryTypeDCL }
