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

type AlterUser struct {
	statementImpl
	IfExists bool
	Users    []*User
	Role     *Role
	MiscOpt  UserMiscOption
	// comment or attribute
	CommentOrAttribute AccountCommentOrAttribute
}

func (node *AlterUser) Format(ctx *FmtCtx) {
	ctx.WriteString("alter user")
	if node.IfExists {
		ctx.WriteString(" if exists")
	}
	if node.Users != nil {
		prefix := " "
		for _, u := range node.Users {
			ctx.WriteString(prefix)
			u.Format(ctx)
			prefix = ", "
		}
	}
	if node.Role != nil {
		ctx.WriteString(" default role ")
		node.Role.Format(ctx)
	}
	if node.MiscOpt != nil {
		prefix := " "
		ctx.WriteString(prefix)
		node.MiscOpt.Format(ctx)
	}
	node.CommentOrAttribute.Format(ctx)
}

func (node *AlterUser) GetStatementType() string { return "Alter User" }
func (node *AlterUser) GetQueryType() string     { return QueryTypeDCL }

func NewAlterUser(e bool, u []*User, r *Role, m UserMiscOption, c AccountCommentOrAttribute, buf *buffer.Buffer) *AlterUser {
	a := buffer.Alloc[AlterUser](buf)
	a.IfExists = e
	a.Role = r
	a.MiscOpt = m
	a.CommentOrAttribute = c
	a.Users = u
	return a
}

type AlterAccountAuthOption struct {
	Exist          bool
	Equal          *BufString
	AdminName      *BufString
	IdentifiedType AccountIdentified
}

func NewAlterAccountAuthOption(exist bool, equal, name string, buf *buffer.Buffer) *AlterAccountAuthOption {
	a := buffer.Alloc[AlterAccountAuthOption](buf)
	a.Exist = exist
	bEqual := NewBufString(equal)
	bAdminName := NewBufString(name)
	buf.Pin(bEqual, bAdminName)
	a.Equal = bEqual
	a.AdminName = bAdminName
	return a
}

func (node *AlterAccountAuthOption) Format(ctx *FmtCtx) {
	if node.Exist {
		ctx.WriteString(" admin_name")
		if len(node.Equal.Get()) != 0 {
			ctx.WriteString(" ")
			ctx.WriteString(node.Equal.Get())
		}

		ctx.WriteString(fmt.Sprintf(" '%s'", node.AdminName.Get()))
		node.IdentifiedType.Format(ctx)
	}
}

type AlterAccount struct {
	statementImpl
	IfExists   bool
	Name       *BufString
	AuthOption AlterAccountAuthOption
	//status_option or not
	StatusOption AccountStatus
	//comment or not
	Comment AccountComment
}

func NewAlterAccount(exist bool, name string, aopt AlterAccountAuthOption, sopt AccountStatus, c AccountComment, buf *buffer.Buffer) *AlterAccount {
	a := buffer.Alloc[AlterAccount](buf)
	a.IfExists = exist
	bName := NewBufString(name)
	buf.Pin(bName)
	a.Name = bName
	a.AuthOption = aopt
	a.StatusOption = sopt
	a.Comment = c
	return a
}

func (ca *AlterAccount) Format(ctx *FmtCtx) {
	ctx.WriteString("alter account ")
	if ca.IfExists {
		ctx.WriteString("if exists ")
	}
	ctx.WriteString(ca.Name.Get())
	ca.AuthOption.Format(ctx)
	ca.StatusOption.Format(ctx)
	ca.Comment.Format(ctx)
}

func (ca *AlterAccount) GetStatementType() string { return "Alter Account" }
func (ca *AlterAccount) GetQueryType() string     { return QueryTypeDCL }

type AlterView struct {
	statementImpl
	IfExists bool
	Name     *TableName
	ColNames IdentifierList
	AsSource *Select
}

func NewAlterView(exist bool, name *TableName, colNames IdentifierList, asSource *Select, buf *buffer.Buffer) *AlterView {
	a := buffer.Alloc[AlterView](buf)
	a.IfExists = exist
	a.Name = name
	a.ColNames = colNames
	a.AsSource = asSource
	return a
}

func (node *AlterView) Format(ctx *FmtCtx) {
	ctx.WriteString("alter ")

	ctx.WriteString("view ")

	if node.IfExists {
		ctx.WriteString("if exists ")
	}

	node.Name.Format(ctx)
	if len(node.ColNames) > 0 {
		ctx.WriteString(" (")
		node.ColNames.Format(ctx)
		ctx.WriteByte(')')
	}
	ctx.WriteString(" as ")
	node.AsSource.Format(ctx)
}

func (node *AlterView) GetStatementType() string { return "Alter View" }
func (node *AlterView) GetQueryType() string     { return QueryTypeDDL }

// alter configuration for mo_mysql_compatibility_mode
type AlterDataBaseConfig struct {
	statementImpl
	AccountName    *BufString
	DbName         *BufString
	IsAccountLevel bool
	UpdateConfig   *BufString
}

func NewAlterDataBaseConfig(accountName, dbName string, isAccountLevel bool, updateConfig string, buf *buffer.Buffer) *AlterDataBaseConfig {
	a := buffer.Alloc[AlterDataBaseConfig](buf)
	bAccountName := NewBufString(accountName)
	bDbName := NewBufString(dbName)
	bUpdateConfig := NewBufString(updateConfig)
	buf.Pin(bAccountName, bDbName, bUpdateConfig)
	a.AccountName = bAccountName
	a.DbName = bDbName
	a.UpdateConfig = bUpdateConfig
	a.IsAccountLevel = isAccountLevel
	return a
}

func (node *AlterDataBaseConfig) Format(ctx *FmtCtx) {

	if node.IsAccountLevel {
		ctx.WriteString("alter ")
		ctx.WriteString("account configuration ")

		ctx.WriteString("for ")
		ctx.WriteString(fmt.Sprintf("%s ", node.AccountName.Get()))
	} else {
		ctx.WriteString("alter ")
		ctx.WriteString("database configuration ")

		ctx.WriteString("for ")
		ctx.WriteString(fmt.Sprintf("%s ", node.DbName.Get()))
	}

	ctx.WriteString("as ")
	ctx.WriteString(fmt.Sprintf("%s ", node.UpdateConfig.Get()))
}

func (node *AlterDataBaseConfig) GetStatementType() string { return "Alter DataBase config" }
func (node *AlterDataBaseConfig) GetQueryType() string     { return QueryTypeDDL }

// AlterTable
// see https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
type AlterTable struct {
	statementImpl
	Table   *TableName
	Options AlterTableOptions
}

func NewAlterTable(table *TableName, options AlterTableOptions, buf *buffer.Buffer) *AlterTable {
	a := buffer.Alloc[AlterTable](buf)
	a.Table = table
	a.Options = options
	return a
}

func (node *AlterTable) Format(ctx *FmtCtx) {
	ctx.WriteString("alter table ")
	node.Table.Format(ctx)

	prefix := " "
	for _, t := range node.Options {
		ctx.WriteString(prefix)
		t.Format(ctx)
		prefix = ", "
	}
}

func (node *AlterTable) GetStatementType() string { return "Alter Table" }
func (node *AlterTable) GetQueryType() string     { return QueryTypeDDL }

type AlterTableOptions = []AlterTableOption

type AlterTableOption interface {
	NodeFormatter
}

type alterOptionImpl struct {
	AlterTableOption
}

type AlterOptionAlterIndex struct {
	alterOptionImpl
	Name       *BufIdentifier
	Visibility VisibleType
}

func NewAlterOptionAlterIndex(name Identifier, visibility VisibleType, buf *buffer.Buffer) *AlterOptionAlterIndex {
	a := buffer.Alloc[AlterOptionAlterIndex](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)

	a.Name = n
	a.Visibility = visibility
	return a
}

func (node *AlterOptionAlterIndex) Format(ctx *FmtCtx) {
	ctx.WriteString("alter index ")
	node.Name.Format(ctx)
	switch node.Visibility {
	case VISIBLE_TYPE_VISIBLE:
		ctx.WriteString(" visible")
	case VISIBLE_TYPE_INVISIBLE:
		ctx.WriteString(" invisible")
	}
}

type AlterOptionAlterCheck struct {
	alterOptionImpl
	Type    *BufString
	Enforce bool
}

func NewAlterOptionAlterCheck(t string, enforce bool, buf *buffer.Buffer) *AlterOptionAlterCheck {
	a := buffer.Alloc[AlterOptionAlterCheck](buf)
	bType := NewBufString(t)
	buf.Pin(bType)
	a.Type = bType
	a.Enforce = enforce
	return a
}

func (node *AlterOptionAlterCheck) Format(ctx *FmtCtx) {
	ctx.WriteString("alter ")
	ctx.WriteString(node.Type.Get() + " ")
	if node.Enforce {
		ctx.WriteString("enforce")
	} else {
		ctx.WriteString("not enforce")
	}
}

type AlterOptionAdd struct {
	alterOptionImpl
	Def TableDef
}

func NewAlterOptionAdd(def TableDef, buf *buffer.Buffer) *AlterOptionAdd {
	a := buffer.Alloc[AlterOptionAdd](buf)
	a.Def = def
	return a
}

func (node *AlterOptionAdd) Format(ctx *FmtCtx) {
	ctx.WriteString("add ")
	node.Def.Format(ctx)
}

type AlterTableDropType int

const (
	AlterTableDropColumn AlterTableDropType = iota
	AlterTableDropIndex
	AlterTableDropKey
	AlterTableDropPrimaryKey
	AlterTableDropForeignKey
)

type AlterOptionDrop struct {
	alterOptionImpl
	Typ  AlterTableDropType
	Name *BufIdentifier
}

func NewAlterOptionDrop(typ AlterTableDropType, name Identifier, buf *buffer.Buffer) *AlterOptionDrop {
	a := buffer.Alloc[AlterOptionDrop](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)

	a.Name = n
	a.Typ = typ
	return a
}

func (node *AlterOptionDrop) Format(ctx *FmtCtx) {
	ctx.WriteString("drop ")
	switch node.Typ {
	case AlterTableDropColumn:
		ctx.WriteString("column ")
		node.Name.Format(ctx)
	case AlterTableDropIndex:
		ctx.WriteString("index ")
		node.Name.Format(ctx)
	case AlterTableDropKey:
		ctx.WriteString("key ")
		node.Name.Format(ctx)
	case AlterTableDropPrimaryKey:
		ctx.WriteString("primary key")
	case AlterTableDropForeignKey:
		ctx.WriteString("foreign key ")
		node.Name.Format(ctx)
	}
}

type AlterTableName struct {
	Name *UnresolvedObjectName
}

func NewAlterTableName(name *UnresolvedObjectName, buf *buffer.Buffer) *AlterTableName {
	a := buffer.Alloc[AlterTableName](buf)
	a.Name = name
	return a
}

func (node *AlterTableName) Format(ctx *FmtCtx) {
	ctx.WriteString("rename to ")
	node.Name.ToTableName().Format(ctx)
}

type AlterColPos struct {
	PreColName *UnresolvedName
	Pos        int32
}

// suggest rename: AlterAddColumnPosition
type AlterAddCol struct {
	Column   *ColumnTableDef
	Position *ColumnPosition
}

func NewAlterAddCol(column *ColumnTableDef, position *ColumnPosition, buf *buffer.Buffer) *AlterAddCol {
	a := buffer.Alloc[AlterAddCol](buf)
	a.Column = column
	a.Position = position
	return a
}

func (node *AlterAddCol) Format(ctx *FmtCtx) {
	ctx.WriteString("add column ")
	node.Column.Format(ctx)
	node.Position.Format(ctx)
}

type AccountsSetOption struct {
	All          bool
	SetAccounts  IdentifierList
	AddAccounts  IdentifierList
	DropAccounts IdentifierList
}

func NewAccountsSetOption(al bool, se, ad, dr IdentifierList, buf *buffer.Buffer) *AccountsSetOption {
	a := buffer.Alloc[AccountsSetOption](buf)
	a.All = al
	a.SetAccounts = se
	a.AddAccounts = ad
	a.DropAccounts = dr
	return a
}

type AlterPublication struct {
	statementImpl
	IfExists    bool
	Name        *BufIdentifier
	AccountsSet *AccountsSetOption
	Comment     *BufString
}

func NewAlterPublication(exist bool, name Identifier, accountsSet *AccountsSetOption, comment string, buf *buffer.Buffer) *AlterPublication {
	a := buffer.Alloc[AlterPublication](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)

	a.IfExists = exist
	a.Name = n
	a.AccountsSet = accountsSet
bComment := NewBufString(comment)
buf.Pin(bComment)
	a.Comment = bComment
	return a
}

func (node *AlterPublication) Format(ctx *FmtCtx) {
	ctx.WriteString("alter publication ")
	if node.IfExists {
		ctx.WriteString("if exists ")
	}
	node.Name.Format(ctx)
	ctx.WriteString(" account ")
	if node.AccountsSet != nil {
		if node.AccountsSet.All {
			ctx.WriteString("all")
		} else {
			if len(node.AccountsSet.SetAccounts) > 0 {
				node.AccountsSet.SetAccounts.Format(ctx)
			}
			if len(node.AccountsSet.AddAccounts) > 0 {
				ctx.WriteString("add ")
				node.AccountsSet.AddAccounts.Format(ctx)
			}
			if len(node.AccountsSet.DropAccounts) > 0 {
				ctx.WriteString("drop ")
				node.AccountsSet.DropAccounts.Format(ctx)
			}
		}
	}
	if node.Comment.Get() != "" {
		ctx.WriteString(" comment ")
		ctx.WriteString(fmt.Sprintf("'%s'", node.Comment.Get()))
	}
}

func (node *AlterPublication) GetStatementType() string { return "Alter Publication" }
func (node *AlterPublication) GetQueryType() string     { return QueryTypeDCL }

type AlterTableModifyColumnClause struct {
	alterOptionImpl
	Typ       AlterTableOptionType
	NewColumn *ColumnTableDef
	Position  *ColumnPosition
}

func NewAlterTableModifyColumnClause(typ AlterTableOptionType, newColumn *ColumnTableDef, position *ColumnPosition, buf *buffer.Buffer) *AlterTableModifyColumnClause {
	a := buffer.Alloc[AlterTableModifyColumnClause](buf)
	a.Typ = typ
	a.NewColumn = newColumn
	a.Position = position
	return a
}

func (node *AlterTableModifyColumnClause) Format(ctx *FmtCtx) {
	ctx.WriteString("modify column ")
	node.NewColumn.Format(ctx)
	node.Position.Format(ctx)
}

type AlterTableChangeColumnClause struct {
	alterOptionImpl
	Typ           AlterTableOptionType
	OldColumnName *UnresolvedName
	NewColumn     *ColumnTableDef
	Position      *ColumnPosition
}

func NewAlterTableChangeColumnClause(typ AlterTableOptionType, oldColumnName *UnresolvedName, newColumn *ColumnTableDef, position *ColumnPosition, buf *buffer.Buffer) *AlterTableChangeColumnClause {
	a := buffer.Alloc[AlterTableChangeColumnClause](buf)
	a.Typ = typ
	a.OldColumnName = oldColumnName
	a.NewColumn = newColumn
	a.Position = position
	return a
}

func (node *AlterTableChangeColumnClause) Format(ctx *FmtCtx) {
	ctx.WriteString("change column")
	ctx.WriteString(" ")
	node.OldColumnName.Format(ctx)
	ctx.WriteString(" ")
	node.NewColumn.Format(ctx)
	node.Position.Format(ctx)
}

type AlterTableAddColumnClause struct {
	alterOptionImpl
	Typ        AlterTableOptionType
	NewColumns []*ColumnTableDef
	Position   *ColumnPosition
	// when Position is not none, the len(NewColumns) must be one
}

func NewAlterTableAddColumnClause(typ AlterTableOptionType, newColumns []*ColumnTableDef, position *ColumnPosition, buf *buffer.Buffer) *AlterTableAddColumnClause {
	a := buffer.Alloc[AlterTableAddColumnClause](buf)
	a.Typ = typ
	a.NewColumns = newColumns
	a.Position = position
	return a
}

func (node *AlterTableAddColumnClause) Format(ctx *FmtCtx) {
	ctx.WriteString("add column ")
	isFirst := true
	for _, column := range node.NewColumns {
		if isFirst {
			column.Format(ctx)
			isFirst = false
		}
		ctx.WriteString(", ")
		column.Format(ctx)
	}
	node.Position.Format(ctx)
}

type AlterTableRenameColumnClause struct {
	alterOptionImpl
	Typ           AlterTableOptionType
	OldColumnName *UnresolvedName
	NewColumnName *UnresolvedName
}

func NewAlterTableRenameColumnClause(typ AlterTableOptionType, oldColumnName *UnresolvedName, newColumnName *UnresolvedName, buf *buffer.Buffer) *AlterTableRenameColumnClause {
	a := buffer.Alloc[AlterTableRenameColumnClause](buf)
	a.Typ = typ
	a.OldColumnName = oldColumnName
	a.NewColumnName = newColumnName
	return a
}

func (node *AlterTableRenameColumnClause) Format(ctx *FmtCtx) {
	ctx.WriteString("rename column ")
	node.OldColumnName.Format(ctx)
	ctx.WriteString(" to ")
	node.NewColumnName.Format(ctx)
}

// AlterColumnOptionType is the type for AlterTableAlterColumn
type AlterColumnOptionType int

// AlterColumnOptionType types.
const (
	AlterColumnOptionSetDefault AlterColumnOptionType = iota
	AlterColumnOptionSetVisibility
	AlterColumnOptionDropDefault
)

type AlterTableAlterColumnClause struct {
	alterOptionImpl
	Typ         AlterTableOptionType
	ColumnName  *UnresolvedName
	DefalutExpr *AttributeDefault
	Visibility  VisibleType
	OptionType  AlterColumnOptionType
}

func NewAlterTableAlterColumnClause(typ AlterTableOptionType, columnName *UnresolvedName, defalutExpr *AttributeDefault, visibility VisibleType, optionType AlterColumnOptionType, buf *buffer.Buffer) *AlterTableAlterColumnClause {
	a := buffer.Alloc[AlterTableAlterColumnClause](buf)
	a.Typ = typ
	a.ColumnName = columnName
	a.DefalutExpr = defalutExpr
	a.Visibility = visibility
	a.OptionType = optionType
	return a
}

func (node *AlterTableAlterColumnClause) Format(ctx *FmtCtx) {
	ctx.WriteString("alter column ")
	node.ColumnName.Format(ctx)
	if node.OptionType == AlterColumnOptionSetDefault {
		ctx.WriteString(" set ")
		node.DefalutExpr.Format(ctx)
	} else if node.OptionType == AlterColumnOptionSetVisibility {
		ctx.WriteString(" set")
		switch node.Visibility {
		case VISIBLE_TYPE_VISIBLE:
			ctx.WriteString(" visible")
		case VISIBLE_TYPE_INVISIBLE:
			ctx.WriteString(" invisible")
		}
	} else {
		ctx.WriteString(" drop default")
	}
}

type AlterTableOrderByColumnClause struct {
	alterOptionImpl
	Typ              AlterTableOptionType
	AlterOrderByList []*AlterColumnOrder
}

func NewAlterTableOrderByColumnClause(typ AlterTableOptionType, alterOrderByList []*AlterColumnOrder, buf *buffer.Buffer) *AlterTableOrderByColumnClause {
	a := buffer.Alloc[AlterTableOrderByColumnClause](buf)
	a.Typ = typ
	a.AlterOrderByList = alterOrderByList
	return a
}

func (node *AlterTableOrderByColumnClause) Format(ctx *FmtCtx) {
	prefix := "order by "
	for _, columnOrder := range node.AlterOrderByList {
		ctx.WriteString(prefix)
		columnOrder.Format(ctx)
		prefix = ", "
	}
}

type AlterColumnOrder struct {
	Column    *UnresolvedName
	Direction Direction
}

func NewAlterColumnOrder(column *UnresolvedName, direction Direction, buf *buffer.Buffer) *AlterColumnOrder {
	a := buffer.Alloc[AlterColumnOrder](buf)
	a.Column = column
	a.Direction = direction
	return a
}

func (node *AlterColumnOrder) Format(ctx *FmtCtx) {
	node.Column.Format(ctx)
	if node.Direction != DefaultDirection {
		ctx.WriteByte(' ')
		ctx.WriteString(node.Direction.String())
	}
}

// AlterTableType is the type for AlterTableOptionType.
type AlterTableOptionType int

// AlterTable types.
const (
	AlterTableModifyColumn AlterTableOptionType = iota
	AlterTableChangeColumn
	AlterTableRenameColumn
	AlterTableAlterColumn
	AlterTableOrderByColumn
	AlterTableAddConstraint
	AlterTableAddColumn
)

// ColumnPositionType is the type for ColumnPosition.
type ColumnPositionType int

// ColumnPosition Types
// Do not change the value of a constant, as there are external dependencies
const (
	ColumnPositionNone  ColumnPositionType = -1
	ColumnPositionFirst ColumnPositionType = 0
	ColumnPositionAfter ColumnPositionType = 1
)

// ColumnPosition represent the position of the newly added column
type ColumnPosition struct {
	NodeFormatter
	// Tp is either ColumnPositionNone, ColumnPositionFirst or ColumnPositionAfter.
	Typ ColumnPositionType
	// RelativeColumn is the column the newly added column after if type is ColumnPositionAfter
	RelativeColumn *UnresolvedName
}

func NewColumnPosition(typ ColumnPositionType, relativeColumn *UnresolvedName, buf *buffer.Buffer) *ColumnPosition {
	a := buffer.Alloc[ColumnPosition](buf)
	a.Typ = typ
	a.RelativeColumn = relativeColumn
	return a
}

func (node *ColumnPosition) Format(ctx *FmtCtx) {
	switch node.Typ {
	case ColumnPositionNone:
		// do nothing
	case ColumnPositionFirst:
		ctx.WriteString(" first")
	case ColumnPositionAfter:
		ctx.WriteString(" after ")
		node.RelativeColumn.Format(ctx)
	}
}
