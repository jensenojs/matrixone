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
	"strings"

	"github.com/matrixorigin/matrixone/pkg/common/buffer"
)

type SelectStatement interface {
	Statement
}

// Select represents a SelectStatement with an ORDER and/or LIMIT.
type Select struct {
	statementImpl
	Select         SelectStatement
	OrderBy        OrderBy
	Limit          *Limit
	With           *With
	Ep             *ExportParam
	SelectLockInfo *SelectLockInfo
}

func (node *Select) Format(ctx *FmtCtx) {
	if node.With != nil {
		node.With.Format(ctx)
		ctx.WriteByte(' ')
	}
	node.Select.Format(ctx)
	if len(node.OrderBy) > 0 {
		ctx.WriteByte(' ')
		node.OrderBy.Format(ctx)
	}
	if node.Limit != nil {
		ctx.WriteByte(' ')
		node.Limit.Format(ctx)
	}
	if node.Ep != nil {
		ctx.WriteByte(' ')
		node.Ep.Format(ctx)
	}
	if node.SelectLockInfo != nil {
		ctx.WriteByte(' ')
		node.SelectLockInfo.Format(ctx)
	}
}

func (node *Select) GetStatementType() string { return "Select" }
func (node *Select) GetQueryType() string     { return QueryTypeDQL }

func NewSelect(s SelectStatement, o OrderBy, l *Limit, ep *ExportParam, sinfo *SelectLockInfo, buf *buffer.Buffer) *Select {
	se := buffer.Alloc[Select](buf)
	se.Select = s
	se.OrderBy = o
	se.Limit = l
	se.Ep = ep
	se.SelectLockInfo = sinfo
	return se
}

// OrderBy represents an ORDER BY clause.
type OrderBy []*Order

func (node *OrderBy) Format(ctx *FmtCtx) {
	prefix := "order by "
	for _, n := range *node {
		ctx.WriteString(prefix)
		n.Format(ctx)
		prefix = ", "
	}
}

// the ordering expression.
type Order struct {
	Expr          Expr
	Direction     Direction
	NullsPosition NullsPosition
	//without order
	NullOrder bool
}

func (node *Order) Format(ctx *FmtCtx) {
	node.Expr.Format(ctx)
	if node.Direction != DefaultDirection {
		ctx.WriteByte(' ')
		ctx.WriteString(node.Direction.String())
	}
	if node.NullsPosition != DefaultNullsPosition {
		ctx.WriteByte(' ')
		ctx.WriteString(node.NullsPosition.String())
	}
}

func NewOrder(e Expr, d Direction, np NullsPosition, o bool, buf *buffer.Buffer) *Order {
	or := buffer.Alloc[Order](buf)
	or.Expr = e
	or.Direction = d
	or.NullsPosition = np
	or.NullOrder = o
	return or
}

// Direction for ordering results.
type Direction int8

// Direction values.
const (
	DefaultDirection Direction = iota
	Ascending
	Descending
)

var directionName = [...]string{
	DefaultDirection: "",
	Ascending:        "asc",
	Descending:       "desc",
}

func (d Direction) String() string {
	if d < 0 || d > Direction(len(directionName)-1) {
		return fmt.Sprintf("Direction(%d)", d)
	}
	return directionName[d]
}

type NullsPosition int8

const (
	DefaultNullsPosition NullsPosition = iota
	NullsFirst
	NullsLast
)

var nullsPositionName = [...]string{
	DefaultNullsPosition: "",
	NullsFirst:           "nulls first",
	NullsLast:            "nulls last",
}

func (np NullsPosition) String() string {
	if np < 0 || np >= NullsPosition(len(nullsPositionName)) {
		return fmt.Sprintf("NullsPosition(%d)", np)
	}
	return nullsPositionName[np]
}

// the LIMIT clause.
type Limit struct {
	Offset, Count Expr
}

func (node *Limit) Format(ctx *FmtCtx) {
	needSpace := false
	if node != nil && node.Count != nil {
		ctx.WriteString("limit ")
		node.Count.Format(ctx)
		needSpace = true
	}
	if node != nil && node.Offset != nil {
		if needSpace {
			ctx.WriteByte(' ')
		}
		ctx.WriteString("offset ")
		node.Offset.Format(ctx)
	}
}

func NewLimit(o, c Expr, buf *buffer.Buffer) *Limit {
	l := buffer.Alloc[Limit](buf)
	l.Offset = o
	l.Count = c
	return l
}

// the parenthesized SELECT/UNION/VALUES statement.
type ParenSelect struct {
	SelectStatement
	Select *Select
}

func NewParenSelect(s *Select, buf *buffer.Buffer) *ParenSelect {
	p := buffer.Alloc[ParenSelect](buf)
	p.Select = s
	return p
}

func (node *ParenSelect) Format(ctx *FmtCtx) {
	ctx.WriteByte('(')
	node.Select.Format(ctx)
	ctx.WriteByte(')')
}

// SelectClause represents a SELECT statement.
type SelectClause struct {
	SelectStatement
	Distinct bool
	Exprs    SelectExprs
	From     *From
	Where    *Where
	GroupBy  GroupBy
	Having   *Where
	Option   *BufString
}

func NewSelectClause(distinct bool, expr SelectExprs, from *From, where *Where, groupby GroupBy, having *Where, opt string, buf *buffer.Buffer) *SelectClause {
	sec := buffer.Alloc[SelectClause](buf)
	sec.Distinct = distinct
	sec.Exprs = expr
	sec.From = from
	sec.Where = where
	sec.GroupBy = groupby
	sec.Having = having
	bOption := NewBufString(opt)
	buf.Pin(bOption)
	sec.Option = bOption
	return sec
}

func (node *SelectClause) Format(ctx *FmtCtx) {
	ctx.WriteString("select ")
	if node.Distinct {
		ctx.WriteString("distinct ")
	}
	if node.Option.Get() != "" {
		ctx.WriteString(node.Option.Get())
		ctx.WriteByte(' ')
	}
	node.Exprs.Format(ctx)
	if len(node.From.Tables) > 0 {
		canFrom := true
		als, ok := node.From.Tables[0].(*AliasedTableExpr)
		if ok {
			tbl, ok := als.Expr.(*TableName)
			if ok {
				if string(tbl.ObjectName.Get()) == "" {
					canFrom = false
				}
			}
		}
		if canFrom {
			ctx.WriteByte(' ')
			node.From.Format(ctx)
		}
	}
	if node.Where != nil {
		ctx.WriteByte(' ')
		node.Where.Format(ctx)
	}
	if len(node.GroupBy) > 0 {
		ctx.WriteByte(' ')
		node.GroupBy.Format(ctx)
	}
	if node.Having != nil {
		ctx.WriteByte(' ')
		node.Having.Format(ctx)
	}
}

func (node *SelectClause) GetStatementType() string { return "Select" }
func (node *SelectClause) GetQueryType() string     { return QueryTypeDQL }

// WHERE or HAVING clause.
type Where struct {
	Type *BufString
	Expr Expr
}

func (node *Where) Format(ctx *FmtCtx) {
	ctx.WriteString(node.Type.Get())
	ctx.WriteByte(' ')
	node.Expr.Format(ctx)
}

const (
	AstWhere  = "where"
	AstHaving = "having"
)

func NewWhere(t string, e Expr, buf *buffer.Buffer) *Where {
	w := buffer.Alloc[Where](buf)
	bType := NewBufString(t)
	buf.Pin(bType)
	w.Type = bType
	w.Expr = e
	return w
}

// SELECT expressions.
type SelectExprs []SelectExpr

func (node *SelectExprs) Format(ctx *FmtCtx) {
	for i, n := range *node {
		if i > 0 {
			ctx.WriteString(", ")
		}
		n.Format(ctx)
	}
}

// a SELECT expression.
type SelectExpr struct {
	exprImpl
	Expr Expr
	As   *CStr
}

func (node *SelectExpr) Format(ctx *FmtCtx) {
	node.Expr.Format(ctx)
	if node.As != nil && !node.As.Empty() {
		ctx.WriteString(" as ")
		ctx.WriteString(node.As.Origin())
	}
}

// a GROUP BY clause.
type GroupBy []Expr

func NewGroupBy(es []Expr, buf *buffer.Buffer) GroupBy {
	g := buffer.MakeSlice[Expr](buf)
	g = es
	return g
}

func (node *GroupBy) Format(ctx *FmtCtx) {
	prefix := "group by "
	for _, n := range *node {
		ctx.WriteString(prefix)
		n.Format(ctx)
		prefix = ", "
	}
}

// func NewGroupBy(es []Expr, buf *buffer.Buffer) *GroupBy {
// 	g := buffer.Alloc[GroupBy](buf)

// 	return a
// }

const (
	JOIN_TYPE_FULL          = "FULL"
	JOIN_TYPE_LEFT          = "LEFT"
	JOIN_TYPE_RIGHT         = "RIGHT"
	JOIN_TYPE_CROSS         = "CROSS"
	JOIN_TYPE_INNER         = "INNER"
	JOIN_TYPE_STRAIGHT      = "STRAIGHT_JOIN"
	JOIN_TYPE_NATURAL       = "NATURAL"
	JOIN_TYPE_NATURAL_LEFT  = "NATURAL LEFT"
	JOIN_TYPE_NATURAL_RIGHT = "NATURAL RIGHT"
)

// the table expression
type TableExpr interface {
	NodeFormatter
}

var _ TableExpr = &Subquery{}

type JoinTableExpr struct {
	TableExpr
	JoinType *BufString
	Left     TableExpr
	Right    TableExpr
	Cond     JoinCond
}

func (node *JoinTableExpr) Format(ctx *FmtCtx) {
	if node.Left != nil {
		node.Left.Format(ctx)
	}
	if node.JoinType.Get() != "" && node.Right != nil {
		ctx.WriteByte(' ')
		ctx.WriteString(strings.ToLower(node.JoinType.Get()))
	}
	if node.JoinType.Get() != JOIN_TYPE_STRAIGHT && node.Right != nil {
		ctx.WriteByte(' ')
		ctx.WriteString("join")
	}
	if node.Right != nil {
		ctx.WriteByte(' ')
		node.Right.Format(ctx)
	}
	if node.Cond != nil {
		ctx.WriteByte(' ')
		node.Cond.Format(ctx)
	}
}

func NewJoinTableExpr(l TableExpr, jt string, r TableExpr, jc JoinCond, buf *buffer.Buffer) *JoinTableExpr {
	jte := buffer.Alloc[JoinTableExpr](buf)
	bJoinType := NewBufString(jt)
	buf.Pin(bJoinType)
	jte.JoinType = bJoinType
	jte.Left = l
	jte.Right = r
	jte.Cond = jc
	return jte
}

// the join condition.
type JoinCond interface {
	NodeFormatter
}

// the NATURAL join condition
type NaturalJoinCond struct {
	JoinCond
}

func (node *NaturalJoinCond) Format(ctx *FmtCtx) {
	ctx.WriteString("natural")
}

func NewNaturalJoinCond(buf *buffer.Buffer) *NaturalJoinCond {
	njc := buffer.Alloc[NaturalJoinCond](buf)
	return njc
}

// the ON condition for join
type OnJoinCond struct {
	JoinCond
	Expr Expr
}

func (node *OnJoinCond) Format(ctx *FmtCtx) {
	ctx.WriteString("on ")
	node.Expr.Format(ctx)
}

func NewOnJoinCond(e Expr, buf *buffer.Buffer) *OnJoinCond {
	ojc := buffer.Alloc[OnJoinCond](buf)
	ojc.Expr = e
	return ojc
}

// the USING condition
type UsingJoinCond struct {
	JoinCond
	Cols IdentifierList
}

func (node *UsingJoinCond) Format(ctx *FmtCtx) {
	ctx.WriteString("using (")
	node.Cols.Format(ctx)
	ctx.WriteByte(')')
}

func NewUsingJoinCond(c IdentifierList, buf *buffer.Buffer) *UsingJoinCond {
	ujc := buffer.Alloc[UsingJoinCond](buf)
	ujc.Cols = c
	return ujc
}

// the parenthesized TableExpr.
type ParenTableExpr struct {
	TableExpr
	Expr TableExpr
}

func (node *ParenTableExpr) Format(ctx *FmtCtx) {
	ctx.WriteByte('(')
	node.Expr.Format(ctx)
	ctx.WriteByte(')')
}

func NewParenTableExpr(e TableExpr, buf *buffer.Buffer) *ParenTableExpr {
	pte := buffer.Alloc[ParenTableExpr](buf)
	pte.Expr = e
	return pte
}

// The alias, optionally with a column list:
// "AS name" or "AS name(col1, col2)".
type AliasClause struct {
	NodeFormatter
	Alias *BufIdentifier
	Cols  IdentifierList
}

func NewAliasClause(alias Identifier, cols IdentifierList, buf *buffer.Buffer) *AliasClause {
	a := buffer.Alloc[AliasClause](buf)
	al := NewBufIdentifier(alias)
	buf.Pin(al)
	a.Alias = al
	a.Cols = cols
	return a
}

func (node *AliasClause) Format(ctx *FmtCtx) {
	if node.Alias.Get() != "" {
		ctx.WriteString(string(node.Alias.Get()))
	}
	if node.Cols != nil {
		ctx.WriteByte('(')
		node.Cols.Format(ctx)
		ctx.WriteByte(')')
	}
}

// the table expression coupled with an optional alias.
type AliasedTableExpr struct {
	TableExpr
	Expr       TableExpr
	As         AliasClause
	IndexHints []*IndexHint
}

func (node *AliasedTableExpr) Format(ctx *FmtCtx) {
	node.Expr.Format(ctx)
	if node.As.Alias.Get() != "" {
		ctx.WriteString(" as ")
		node.As.Format(ctx)
	}
	if node.IndexHints != nil {
		prefix := " "
		for _, hint := range node.IndexHints {
			ctx.WriteString(prefix)
			hint.Format(ctx)
			prefix = " "
		}
	}
}

func NewAliasedTableExpr(e TableExpr, a *AliasClause, idxs []*IndexHint, buf *buffer.Buffer) *AliasedTableExpr {
	ate := buffer.Alloc[AliasedTableExpr](buf)
	ate.Expr = e
	if a != nil {
		ate.As = *a
	}
	ate.IndexHints = idxs
	return ate
}

// the statements as a data source includes the select statement.
type StatementSource struct {
	TableExpr
	Statement Statement
}

func NewStatementSource(s Statement, buf *buffer.Buffer) *StatementSource {
	ss := buffer.Alloc[StatementSource](buf)
	ss.Statement = s
	return ss
}

// the list of table expressions.
type TableExprs []TableExpr

func (node *TableExprs) Format(ctx *FmtCtx) {
	prefix := ""
	for _, n := range *node {
		ctx.WriteString(prefix)
		n.Format(ctx)
		prefix = ", "
	}
}

// the FROM clause.
type From struct {
	Tables TableExprs
}

func (node *From) Format(ctx *FmtCtx) {
	ctx.WriteString("from ")
	node.Tables.Format(ctx)
}

func NewFrom(t TableExprs, buf *buffer.Buffer) *From {
	f := buffer.Alloc[From](buf)
	f.Tables = t
	return f
}

type IndexHintType int

const (
	HintUse IndexHintType = iota + 1
	HintIgnore
	HintForce
)

type IndexHintScope int

// Index hint scopes.
const (
	HintForScan IndexHintScope = iota + 1
	HintForJoin
	HintForOrderBy
	HintForGroupBy
)

type IndexHint struct {
	IndexNames []string // do NOT reassign after NewIndexHint
	HintType   IndexHintType
	HintScope  IndexHintScope
}

func NewIndexHint(indexNames []string, hintType IndexHintType, hintScope IndexHintScope, buf *buffer.Buffer) *IndexHint {
	i := buffer.Alloc[IndexHint](buf)
	i.HintType = hintType
	i.HintScope = hintScope

	if indexNames != nil {
		i.IndexNames = buffer.MakeSlice[string](buf)
		for _, v := range indexNames {
			i.IndexNames = buffer.AppendSlice[string](buf, i.IndexNames, buf.CopyString(v))
		}
	}
	return i
}

func (node *IndexHint) Format(ctx *FmtCtx) {
	indexHintType := ""
	switch node.HintType {
	case HintUse:
		indexHintType = "use index"
	case HintIgnore:
		indexHintType = "ignore index"
	case HintForce:
		indexHintType = "force index"
	}

	indexHintScope := ""
	switch node.HintScope {
	case HintForScan:
		indexHintScope = ""
	case HintForJoin:
		indexHintScope = " for join"
	case HintForOrderBy:
		indexHintScope = " for order by"
	case HintForGroupBy:
		indexHintScope = " for group by"
	}
	ctx.WriteString(indexHintType)
	ctx.WriteString(indexHintScope)
	ctx.WriteString("(")
	if node.IndexNames != nil {
		for i, value := range node.IndexNames {
			if i > 0 {
				ctx.WriteString(", ")
			}
			ctx.WriteString(value)
		}
	}
	ctx.WriteString(")")
}
