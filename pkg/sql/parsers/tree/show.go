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

type Show interface {
	Explain
}

type showImpl struct {
	Show
}

// SHOW CREATE TABLE statement
type ShowCreateTable struct {
	showImpl
	Name *UnresolvedObjectName
}

func NewShowCreateTable(n *UnresolvedObjectName, buf *buffer.Buffer) *ShowCreateTable {
	s := buffer.Alloc[ShowCreateTable](buf)
	s.Name = n
	return s
}

func (node *ShowCreateTable) Format(ctx *FmtCtx) {
	ctx.WriteString("show create table ")
	node.Name.ToTableName().Format(ctx)
}

func (node *ShowCreateTable) GetStatementType() string { return "Show Create Table" }
func (node *ShowCreateTable) GetQueryType() string     { return QueryTypeOth }

func NewShowCreate(n *UnresolvedObjectName, buf *buffer.Buffer) *ShowCreateTable {
	sc := buffer.Alloc[ShowCreateTable](buf)
	sc.Name = n
	return sc
}

// SHOW CREATE VIEW statement
type ShowCreateView struct {
	showImpl
	Name *UnresolvedObjectName
}

func (node *ShowCreateView) Format(ctx *FmtCtx) {
	ctx.WriteString("show create view ")
	node.Name.ToTableName().Format(ctx)
}
func (node *ShowCreateView) GetStatementType() string { return "Show Create View" }
func (node *ShowCreateView) GetQueryType() string     { return QueryTypeOth }

func NewShowCreateView(n *UnresolvedObjectName, buf *buffer.Buffer) *ShowCreateView {
	var scv *ShowCreateView
	if buf != nil {
		scv = buffer.Alloc[ShowCreateView](buf)
	} else {
		scv = new(ShowCreateView)
	}
	scv.Name = n
	return scv
}

// SHOW CREATE DATABASE statement
type ShowCreateDatabase struct {
	showImpl
	IfNotExists bool
	Name        string
}

func (node *ShowCreateDatabase) Format(ctx *FmtCtx) {
	ctx.WriteString("show create database")
	if node.IfNotExists {
		ctx.WriteString(" if not exists")
	}
	ctx.WriteByte(' ')
	ctx.WriteString(string(node.Name))
}
func (node *ShowCreateDatabase) GetStatementType() string { return "Show Create View" }
func (node *ShowCreateDatabase) GetQueryType() string     { return QueryTypeOth }

func NewShowCreateDatabase(i bool, n string, buf *buffer.Buffer) *ShowCreateDatabase {
	scd := buffer.Alloc[ShowCreateDatabase](buf)
	scd.IfNotExists = i
	scd.Name = n
	return scd
}

// SHOW COLUMNS statement.
type ShowColumns struct {
	showImpl
	Ext     bool
	Full    bool
	Table   *UnresolvedObjectName
	ColName *UnresolvedName
	DBName  string
	Like    *ComparisonExpr
	Where   *Where
}

func (node *ShowColumns) Format(ctx *FmtCtx) {
	ctx.WriteString("show")
	if node.Ext {
		ctx.WriteString(" extended")
	}
	if node.Full {
		ctx.WriteString(" full")
	}
	ctx.WriteString(" columns")
	if node.Table != nil {
		ctx.WriteString(" from ")
		node.Table.Format(ctx)
	}
	if node.DBName != "" {
		ctx.WriteString(" from ")
		ctx.WriteString(node.DBName)
	}
	if node.Like != nil {
		ctx.WriteByte(' ')
		node.Like.Format(ctx)
	}
	if node.Where != nil {
		ctx.WriteByte(' ')
		node.Where.Format(ctx)
	}
}
func (node *ShowColumns) GetStatementType() string { return "Show Columns" }
func (node *ShowColumns) GetQueryType() string     { return QueryTypeOth }

func NewShowColumns1(t *UnresolvedObjectName, cn *UnresolvedName, buf *buffer.Buffer) *ShowColumns {
	sc := buffer.Alloc[ShowColumns](buf)
	sc.Table = t
	sc.ColName = cn
	return sc
}

func NewShowColumns2(e bool, f bool, t *UnresolvedObjectName, d string, l *ComparisonExpr, w *Where, buf *buffer.Buffer) *ShowColumns {
	sc := buffer.Alloc[ShowColumns](buf)
	sc.Ext = e
	sc.Full = f
	sc.Table = t
	sc.DBName = d
	sc.Like = l
	sc.Where = w
	return sc
}

// the SHOW DATABASES statement.
type ShowDatabases struct {
	showImpl
	Like  *ComparisonExpr
	Where *Where
}

func (node *ShowDatabases) Format(ctx *FmtCtx) {
	ctx.WriteString("show databases")
	if node.Like != nil {
		ctx.WriteByte(' ')
		node.Like.Format(ctx)
	}
	if node.Where != nil {
		ctx.WriteByte(' ')
		node.Where.Format(ctx)
	}
}
func (node *ShowDatabases) GetStatementType() string { return "Show Databases" }
func (node *ShowDatabases) GetQueryType() string     { return QueryTypeOth }

func NewShowDatabases(l *ComparisonExpr, w *Where, buf *buffer.Buffer) *ShowDatabases {
	sd := buffer.Alloc[ShowDatabases](buf)
	sd.Like = l
	sd.Where = w
	return sd
}

type ShowType int

const (
	ShowEngines = iota
	ShowCharset
	ShowCreateUser
	ShowTriggers
	ShowConfig
	ShowEvents
	ShowPlugins
	ShowProfile
	ShowProfiles
	ShowPrivileges
)

func (s ShowType) String() string {
	switch s {
	case ShowEngines:
		return "engines"
	case ShowCharset:
		return "charset"
	case ShowCreateUser:
		return "create user"
	case ShowTriggers:
		return "triggers"
	case ShowConfig:
		return "config"
	case ShowEvents:
		return "events"
	case ShowPlugins:
		return "plugins"
	case ShowProfile:
		return "profile"
	case ShowProfiles:
		return "profiles"
	case ShowPrivileges:
		return "privileges"
	default:
		return "not implemented"
	}
}

type ShowTarget struct {
	showImpl
	Global bool
	Type   ShowType
	DbName string
	Like   *ComparisonExpr
	Where  *Where
}

func NewShowTarget(n string, t ShowType, buf *buffer.Buffer) *ShowTarget {
	sh := buffer.Alloc[ShowTarget](buf)
	sh.Type = t
	sh.DbName = n
	return sh
}

func (node *ShowTarget) Format(ctx *FmtCtx) {
	ctx.WriteString("show ")
	if node.Global {
		ctx.WriteString("global ")
	}
	ctx.WriteString(node.Type.String())
	if node.DbName != "" {
		ctx.WriteString(" from ")
		ctx.WriteString(node.DbName)
	}
	if node.Like != nil {
		ctx.WriteByte(' ')
		node.Like.Format(ctx)
	}
	if node.Where != nil {
		ctx.WriteByte(' ')
		node.Where.Format(ctx)
	}
}
func (node *ShowTarget) GetStatementType() string { return "Show Target" }
func (node *ShowTarget) GetQueryType() string     { return QueryTypeOth }

type ShowTableStatus struct {
	showImpl
	DbName string
	Like   *ComparisonExpr
	Where  *Where
}

func NewShowTableStatus(dbName string, l *ComparisonExpr, w *Where, buf *buffer.Buffer) *ShowTableStatus {
	st := buffer.Alloc[ShowTableStatus](buf)
	st.DbName = dbName
	st.Like = l
	st.Where = w
	return st
}

func (node *ShowTableStatus) Format(ctx *FmtCtx) {
	ctx.WriteString("show table status")
	if node.DbName != "" {
		ctx.WriteString(" from ")
		ctx.WriteString(node.DbName)
	}
	if node.Like != nil {
		ctx.WriteByte(' ')
		node.Like.Format(ctx)
	}
	if node.Where != nil {
		ctx.WriteByte(' ')
		node.Where.Format(ctx)
	}
}
func (node *ShowTableStatus) GetStatementType() string { return "Show Table Status" }
func (node *ShowTableStatus) GetQueryType() string     { return QueryTypeOth }

type ShowGrants struct {
	showImpl
	Username      string
	Hostname      string
	Roles         []*Role
	ShowGrantType ShowGrantType
}

func NewShowGrants(u, h string, rs []*Role, s ShowGrantType, buf *buffer.Buffer) *ShowGrants {
	sh := buffer.Alloc[ShowGrants](buf)
	sh.Username = u
	sh.Hostname = h
	sh.ShowGrantType = s
	sh.Roles = rs
	return sh
}

type ShowGrantType int

const (
	GrantForUser = iota
	GrantForRole
)

func (node *ShowGrants) Format(ctx *FmtCtx) {
	if node.ShowGrantType == GrantForRole {
		ctx.WriteString("show grants")
		if node.Roles != nil {
			ctx.WriteString("for")
			ctx.WriteString(" ")
			ctx.WriteString(node.Roles[0].UserName)
		}
	} else {
		ctx.WriteString("show grants")
		if node.Username != "" {
			ctx.WriteString(" for ")
			ctx.WriteString(node.Username)
			if node.Hostname != "" {
				ctx.WriteString("@")
				ctx.WriteString(node.Hostname)
			}
		}
		if node.Roles != nil {
			prefix := ""
			for _, r := range node.Roles {
				ctx.WriteString(prefix)
				r.Format(ctx)
				prefix = ", "
			}
		}
	}
}
func (node *ShowGrants) GetStatementType() string { return "Show Grants" }
func (node *ShowGrants) GetQueryType() string     { return QueryTypeOth }

// SHOW SEQUENCES statement.
type ShowSequences struct {
	showImpl
	DBName string
	Where  *Where
}

func NewShowSequences(dbName string, where *Where, buf *buffer.Buffer) *ShowSequences {
	sh := buffer.Alloc[ShowSequences](buf)
	sh.DBName = dbName
	sh.Where = where
	return sh
}

func (node *ShowSequences) Format(ctx *FmtCtx) {
	ctx.WriteString("show sequences")
	if node.DBName != "" {
		ctx.WriteString(" from ")
		ctx.WriteString(node.DBName)
	}
	if node.Where != nil {
		ctx.WriteByte(' ')
		node.Where.Format(ctx)
	}
}
func (node *ShowSequences) GetStatementType() string { return "Show Sequences" }
func (node *ShowSequences) GetQueryType() string     { return QueryTypeOth }

// SHOW TABLES statement.
type ShowTables struct {
	showImpl
	Ext    bool
	Open   bool
	Full   bool
	DBName string
	Like   *ComparisonExpr
	Where  *Where
}

func (node *ShowTables) Format(ctx *FmtCtx) {
	ctx.WriteString("show")
	if node.Open {
		ctx.WriteString(" open")
	}
	if node.Full {
		ctx.WriteString(" full")
	}
	ctx.WriteString(" tables")
	if node.DBName != "" {
		ctx.WriteString(" from ")
		ctx.WriteString(node.DBName)
	}
	if node.Like != nil {
		ctx.WriteByte(' ')
		node.Like.Format(ctx)
	}
	if node.Where != nil {
		ctx.WriteByte(' ')
		node.Where.Format(ctx)
	}
}
func (node *ShowTables) GetStatementType() string { return "Show Tables" }
func (node *ShowTables) GetQueryType() string     { return QueryTypeOth }

func NewShowTables(o, f bool, n string, l *ComparisonExpr, w *Where, buf *buffer.Buffer) *ShowTables {
	st := buffer.Alloc[ShowTables](buf)
	st.Open = o
	st.Full = f
	st.DBName = n
	st.Like = l
	st.Where = w
	return st
}

// SHOW PROCESSLIST
type ShowProcessList struct {
	showImpl
	Full bool
}

func (node *ShowProcessList) Format(ctx *FmtCtx) {
	ctx.WriteString("show")
	if node.Full {
		ctx.WriteString(" full")
	}
	ctx.WriteString(" processlist")
}
func (node *ShowProcessList) GetStatementType() string { return "Show Processlist" }
func (node *ShowProcessList) GetQueryType() string     { return QueryTypeOth }

func NewShowProcessList(f bool, buf *buffer.Buffer) *ShowProcessList {
	spl := buffer.Alloc[ShowProcessList](buf)
	spl.Full = f
	return spl
}

type ShowErrors struct {
	showImpl
}

func (node *ShowErrors) Format(ctx *FmtCtx) {
	ctx.WriteString("show errors")
}
func (node *ShowErrors) GetStatementType() string { return "Show Errors" }
func (node *ShowErrors) GetQueryType() string     { return QueryTypeOth }

func NewShowErrors(buf *buffer.Buffer) *ShowErrors {
	se := buffer.Alloc[ShowErrors](buf)
	return se
}

type ShowWarnings struct {
	showImpl
}

func (node *ShowWarnings) Format(ctx *FmtCtx) {
	ctx.WriteString("show warnings")
}
func (node *ShowWarnings) GetStatementType() string { return "Show Warnings" }
func (node *ShowWarnings) GetQueryType() string     { return QueryTypeOth }

func NewShowWarnings(buf *buffer.Buffer) *ShowWarnings {
	sw := buffer.Alloc[ShowWarnings](buf)
	return sw
}

// SHOW collation statement
type ShowCollation struct {
	showImpl
	Like  *ComparisonExpr
	Where *Where
}

func NewShowCollation(buf *buffer.Buffer) *ShowCollation {
	s := buffer.Alloc[ShowCollation](buf)
	return s
}

func (node *ShowCollation) Format(ctx *FmtCtx) {
	ctx.WriteString("show collation")
	if node.Like != nil {
		ctx.WriteString(" like ")
		node.Like.Format(ctx)
	}
	if node.Where != nil {
		ctx.WriteByte(' ')
		node.Where.Format(ctx)
	}
}
func (node *ShowCollation) GetStatementType() string { return "Show Collation" }
func (node *ShowCollation) GetQueryType() string     { return QueryTypeOth }

// SHOW VARIABLES statement
// System Variables
type ShowVariables struct {
	showImpl
	Global bool
	Like   *ComparisonExpr
	Where  *Where
}

func (node *ShowVariables) Format(ctx *FmtCtx) {
	ctx.WriteString("show")
	if node.Global {
		ctx.WriteString(" global")
	}
	ctx.WriteString(" variables")
	if node.Like != nil {
		ctx.WriteByte(' ')
		node.Like.Format(ctx)
	}
	if node.Where != nil {
		ctx.WriteByte(' ')
		node.Where.Format(ctx)
	}
}
func (node *ShowVariables) GetStatementType() string { return "Show Variables" }
func (node *ShowVariables) GetQueryType() string     { return QueryTypeOth }

func NewShowVariables(g bool, l *ComparisonExpr, w *Where, buf *buffer.Buffer) *ShowVariables {
	sv := buffer.Alloc[ShowVariables](buf)
	sv.Global = g
	sv.Like = l
	sv.Where = w
	return sv
}

// SHOW STATUS statement
type ShowStatus struct {
	showImpl
	Global bool
	Like   *ComparisonExpr
	Where  *Where
}

func (node *ShowStatus) Format(ctx *FmtCtx) {
	ctx.WriteString("show")
	if node.Global {
		ctx.WriteString(" global")
	}
	ctx.WriteString(" status")
	if node.Like != nil {
		ctx.WriteString(" like ")
		node.Like.Format(ctx)
	}
	if node.Where != nil {
		ctx.WriteByte(' ')
		node.Where.Format(ctx)
	}
}
func (node *ShowStatus) GetStatementType() string { return "Show Status" }
func (node *ShowStatus) GetQueryType() string     { return QueryTypeOth }

func NewShowStatus(g bool, l *ComparisonExpr, w *Where, buf *buffer.Buffer) *ShowStatus {
	ss := buffer.Alloc[ShowStatus](buf)
	ss.Global = g
	ss.Like = l
	ss.Where = w
	return ss
}

// show index statement
type ShowIndex struct {
	showImpl
	TableName TableName
	Where     *Where
}

func (node *ShowIndex) Format(ctx *FmtCtx) {
	ctx.WriteString("show index from ")
	node.TableName.Format(ctx)
	if node.Where != nil {
		ctx.WriteByte(' ')
		node.Where.Format(ctx)
	}
}
func (node *ShowIndex) GetStatementType() string { return "Show Index" }
func (node *ShowIndex) GetQueryType() string     { return QueryTypeOth }

func NewShowIndex(t *TableName, w *Where, buf *buffer.Buffer) *ShowIndex {
	si := buffer.Alloc[ShowIndex](buf)
	if t != nil {
		si.TableName = *t
	}
	si.Where = w
	return si
}

// show Function or Procedure statement

type ShowFunctionOrProcedureStatus struct {
	showImpl
	Like       *ComparisonExpr
	Where      *Where
	IsFunction bool
}

func (node *ShowFunctionOrProcedureStatus) Format(ctx *FmtCtx) {
	if node.IsFunction {
		ctx.WriteString("show function status")
	} else {
		ctx.WriteString("show procedure status")
	}
	if node.Like != nil {
		ctx.WriteString(" like ")
		node.Like.Format(ctx)
	}
	if node.Where != nil {
		ctx.WriteByte(' ')
		node.Where.Format(ctx)
	}
}

func (node *ShowFunctionOrProcedureStatus) GetStatementType() string {
	return "Show Function Or Procedure Status"
}
func (node *ShowFunctionOrProcedureStatus) GetQueryType() string { return QueryTypeOth }

func NewShowFunctionOrProcedureStatus(l *ComparisonExpr, w *Where, i bool, buf *buffer.Buffer) *ShowFunctionOrProcedureStatus {
	sfps := buffer.Alloc[ShowFunctionOrProcedureStatus](buf)
	sfps.Like = l
	sfps.Where = w
	sfps.IsFunction = i
	return sfps
}

// show node list
type ShowNodeList struct {
	showImpl
}

func (node *ShowNodeList) Format(ctx *FmtCtx) {
	ctx.WriteString("show node list")
}

func (node *ShowNodeList) GetStatementType() string { return "Show Node List" }
func (node *ShowNodeList) GetQueryType() string     { return QueryTypeOth }

func NewShowNodeList(buf *buffer.Buffer) *ShowNodeList {
	snl := buffer.Alloc[ShowNodeList](buf)
	return snl
}

// show locks
type ShowLocks struct {
	showImpl
}

func (node *ShowLocks) Format(ctx *FmtCtx) {
	ctx.WriteString("show locks")
}

func (node *ShowLocks) GetStatementType() string { return "Show Locks" }
func (node *ShowLocks) GetQueryType() string     { return QueryTypeOth }

func NewShowLocks(buf *buffer.Buffer) *ShowLocks {
	sl := buffer.Alloc[ShowLocks](buf)
	return sl
}

// show table number
type ShowTableNumber struct {
	showImpl
	DbName string
}

func (node *ShowTableNumber) Format(ctx *FmtCtx) {
	ctx.WriteString("show table number")
	if node.DbName != "" {
		ctx.WriteString(" from ")
		ctx.WriteString(node.DbName)
	}
}
func (node *ShowTableNumber) GetStatementType() string { return "Show Table Number" }
func (node *ShowTableNumber) GetQueryType() string     { return QueryTypeOth }

func NewShowTableNumber(dbname string, buf *buffer.Buffer) *ShowTableNumber {
	stn := buffer.Alloc[ShowTableNumber](buf)
	stn.DbName = dbname
	return stn
}

// show column number
type ShowColumnNumber struct {
	showImpl
	Table  *UnresolvedObjectName
	DbName string
}

func (node *ShowColumnNumber) Format(ctx *FmtCtx) {
	ctx.WriteString("show column number")
	if node.Table != nil {
		ctx.WriteString(" from ")
		node.Table.Format(ctx)
	}
	if node.DbName != "" {
		ctx.WriteString(" from ")
		ctx.WriteString(node.DbName)
	}
}
func (node *ShowColumnNumber) GetStatementType() string { return "Show Column Number" }
func (node *ShowColumnNumber) GetQueryType() string     { return QueryTypeOth }

func NewShowColumnNumber(table *UnresolvedObjectName, dbname string, buf *buffer.Buffer) *ShowColumnNumber {
	scn := buffer.Alloc[ShowColumnNumber](buf)
	scn.Table = table
	scn.DbName = dbname
	return scn
}

// show table values
type ShowTableValues struct {
	showImpl
	Table  *UnresolvedObjectName
	DbName string
}

func (node *ShowTableValues) Format(ctx *FmtCtx) {
	ctx.WriteString("show table values")
	if node.Table != nil {
		ctx.WriteString(" from ")
		node.Table.Format(ctx)
	}
	if node.DbName != "" {
		ctx.WriteString(" from ")
		ctx.WriteString(node.DbName)
	}
}
func (node *ShowTableValues) GetStatementType() string { return "Show Table Values" }
func (node *ShowTableValues) GetQueryType() string     { return QueryTypeOth }

func NewShowTableValues(table *UnresolvedObjectName, dbname string, buf *buffer.Buffer) *ShowTableValues {
	stv := buffer.Alloc[ShowTableValues](buf)
	stv.Table = table
	stv.DbName = dbname
	return stv
}

type ShowAccounts struct {
	showImpl
	Like *ComparisonExpr
}

func NewShowAccounts(l *ComparisonExpr, buf *buffer.Buffer) *ShowAccounts {
	s := buffer.Alloc[ShowAccounts](buf)
	s.Like = l
	return s
}

func (node *ShowAccounts) Format(ctx *FmtCtx) {
	ctx.WriteString("show accounts")
	if node.Like != nil {
		ctx.WriteByte(' ')
		node.Like.Format(ctx)
	}
}

func (node *ShowAccounts) GetStatementType() string { return "Show Accounts" }
func (node *ShowAccounts) GetQueryType() string     { return QueryTypeOth }

type ShowPublications struct {
	showImpl
	Like *ComparisonExpr
}

func NewShowPublications(l *ComparisonExpr, buf *buffer.Buffer) *ShowPublications {
	s := buffer.Alloc[ShowPublications](buf)
	s.Like = l
	return s
}

func (node *ShowPublications) Format(ctx *FmtCtx) {
	ctx.WriteString("show publications")
	if node.Like != nil {
		ctx.WriteByte(' ')
		node.Like.Format(ctx)
	}
}

func (node *ShowPublications) GetStatementType() string { return "Show Publications" }
func (node *ShowPublications) GetQueryType() string     { return QueryTypeOth }

type ShowSubscriptions struct {
	showImpl
	Like *ComparisonExpr
}

func NewShowSubscriptions(l *ComparisonExpr, buf *buffer.Buffer) *ShowSubscriptions {
	s := buffer.Alloc[ShowSubscriptions](buf)
	s.Like = l
	return s
}

func (node *ShowSubscriptions) Format(ctx *FmtCtx) {
	ctx.WriteString("show subscriptions")
	if node.Like != nil {
		ctx.WriteByte(' ')
		node.Like.Format(ctx)
	}
}
func (node *ShowSubscriptions) GetStatementType() string { return "Show Subscriptions" }
func (node *ShowSubscriptions) GetQueryType() string     { return QueryTypeOth }

type ShowCreatePublications struct {
	showImpl
	Name string
}

func NewShowCreatePublications(n string, buf *buffer.Buffer) *ShowCreatePublications {
	s := buffer.Alloc[ShowCreatePublications](buf)
	s.Name = n
	return s
}

func (node *ShowCreatePublications) Format(ctx *FmtCtx) {
	ctx.WriteString("show create publication ")
	ctx.WriteString(node.Name)
}
func (node *ShowCreatePublications) GetStatementType() string { return "Show Create Publication" }
func (node *ShowCreatePublications) GetQueryType() string     { return QueryTypeOth }

type ShowTableSize struct {
	showImpl
	Table  *UnresolvedObjectName
	DbName string
}

func (node *ShowTableSize) Format(ctx *FmtCtx) {
	ctx.WriteString("show table size")
	if node.Table != nil {
		ctx.WriteString(" from ")
		node.Table.Format(ctx)
	}
	if node.DbName != "" {
		ctx.WriteString(" from ")
		ctx.WriteString(node.DbName)
	}
}
func (node *ShowTableSize) GetStatementType() string { return "Show Table Size" }
func (node *ShowTableSize) GetQueryType() string     { return QueryTypeDQL }

func NewShowTableSize(table *UnresolvedObjectName, dbname string, buf *buffer.Buffer) *ShowTableSize {
	sts := buffer.Alloc[ShowTableSize](buf)
	sts.Table = table
	sts.DbName = dbname
	return sts
}

// show Roles statement

type ShowRolesStmt struct {
	showImpl
	Like *ComparisonExpr
}

func NewShowRolesStmt(l *ComparisonExpr, buf *buffer.Buffer) *ShowRolesStmt {
	s := buffer.Alloc[ShowRolesStmt](buf)
	s.Like = l
	return s
}

func (node *ShowRolesStmt) Format(ctx *FmtCtx) {
	ctx.WriteString("show roles")
	if node.Like != nil {
		ctx.WriteString(" ")
		node.Like.Format(ctx)
	}
}

func (node *ShowRolesStmt) GetStatementType() string { return "Show Roles" }
func (node *ShowRolesStmt) GetQueryType() string     { return QueryTypeOth }

// ShowBackendServers indicates SHOW BACKEND SERVERS statement.
type ShowBackendServers struct {
	showImpl
}

func NewShowBackendServers(buf *buffer.Buffer) *ShowBackendServers {
	s := buffer.Alloc[ShowBackendServers](buf)
	return s
}

func (node *ShowBackendServers) Format(ctx *FmtCtx) {
	ctx.WriteString("show backend servers")
}

func (node *ShowBackendServers) GetStatementType() string { return "Show Backend Servers" }
func (node *ShowBackendServers) GetQueryType() string     { return QueryTypeOth }

type ShowConnectors struct {
	showImpl
}

func (node *ShowConnectors) Format(ctx *FmtCtx) {
	ctx.WriteString("show connectors")
}
func (node *ShowConnectors) GetStatementType() string { return "Show Connectors" }
func (node *ShowConnectors) GetQueryType() string     { return QueryTypeOth }

func NewShowConnectors(buf *buffer.Buffer) *ShowConnectors {
	sc := buffer.Alloc[ShowConnectors](buf)
	return sc
}

type EmptyStmt struct {
	statementImpl
}

func (e *EmptyStmt) String() string {
	return ""
}

func (e *EmptyStmt) Format(ctx *FmtCtx) {
	ctx.WriteString("")
}

func (e EmptyStmt) GetStatementType() string {
	return "InternalCmd"
}

func (e EmptyStmt) GetQueryType() string {
	return QueryTypeOth
}
