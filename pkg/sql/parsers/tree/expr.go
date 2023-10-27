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

// AST for the expression
type Expr interface {
	fmt.Stringer
	NodeFormatter
	NodeChecker
}

type exprImpl struct {
	Expr
}

func (node *exprImpl) String() string {
	return ""
}

// Binary Operator
type BinaryOp int

const (
	PLUS BinaryOp = iota
	MINUS
	MULTI
	DIV         // /
	INTEGER_DIV //
	BIT_OR      // |
	BIT_AND     // &
	BIT_XOR     // ^
	LEFT_SHIFT  // <<
	RIGHT_SHIFT // >>
	MOD         // %
)

func (op BinaryOp) ToString() string {
	switch op {
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case MULTI:
		return "*"
	case DIV:
		return "/"
	case INTEGER_DIV:
		return "div"
	case BIT_OR:
		return "|"
	case BIT_AND:
		return "&"
	case BIT_XOR:
		return "^"
	case LEFT_SHIFT:
		return "<<"
	case RIGHT_SHIFT:
		return ">>"
	case MOD:
		return "%"
	default:
		return "Unknown BinaryExprOperator"
	}
}

// binary expression
type BinaryExpr struct {
	exprImpl

	//operator
	Op BinaryOp

	//left expression
	Left Expr

	//right expression
	Right Expr
}

func (node *BinaryExpr) Format(ctx *FmtCtx) {
	ctx.PrintExpr(node, node.Left, true)
	ctx.WriteByte(' ')
	ctx.WriteString(node.Op.ToString())
	ctx.WriteByte(' ')
	ctx.PrintExpr(node, node.Right, false)
}

// Accept implements NodeChecker interface
func (node *BinaryExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}

	node = newNode.(*BinaryExpr)
	tmpNode, ok := node.Left.Accept(v)
	if !ok {
		return node, false
	}
	node.Left = tmpNode

	tmpNode, ok = node.Right.Accept(v)
	if !ok {
		return node, false
	}
	node.Right = tmpNode
	return v.Exit(node)
}

func NewBinaryExpr(op BinaryOp, left Expr, right Expr, buf *buffer.Buffer) *BinaryExpr {
	var b *BinaryExpr
	if buf != nil {
		b = buffer.Alloc[BinaryExpr](buf)
	} else {
		b = new(BinaryExpr)
	}
	b.Op = op
	b.Left = left
	b.Right = right
	return b
}

// unary expression
type UnaryOp int

const (
	//-
	UNARY_MINUS UnaryOp = iota
	//+
	UNARY_PLUS
	//~
	UNARY_TILDE
	//!
	UNARY_MARK
)

func (op UnaryOp) ToString() string {
	switch op {
	case UNARY_MINUS:
		return "-"
	case UNARY_PLUS:
		return "+"
	case UNARY_TILDE:
		return "~"
	case UNARY_MARK:
		return "!"
	default:
		return "Unknown UnaryExprOperator"
	}
}

// unary expression
type UnaryExpr struct {
	exprImpl

	//operator
	Op UnaryOp

	//expression
	Expr Expr
}

func (e *UnaryExpr) Format(ctx *FmtCtx) {
	if _, unary := e.Expr.(*UnaryExpr); unary {
		ctx.WriteString(e.Op.ToString())
		ctx.WriteByte(' ')
		ctx.PrintExpr(e, e.Expr, true)
		return
	}
	ctx.WriteString(e.Op.ToString())
	ctx.PrintExpr(e, e.Expr, true)
}

func (e *UnaryExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(e)
	if skipChildren {
		return v.Exit(newNode)
	}
	e = newNode.(*UnaryExpr)
	node, ok := e.Expr.Accept(v)
	if !ok {
		return e, false
	}
	e.Expr = node
	return v.Exit(e)
}

func (e *UnaryExpr) String() string {
	return unaryOpName[e.Op] + e.Expr.String()
}

func NewUnaryExpr(op UnaryOp, expr Expr, buf *buffer.Buffer) *UnaryExpr {
	u := buffer.Alloc[UnaryExpr](buf)
	u.Op = op
	u.Expr = expr
	return u
}

var unaryOpName = []string{
	"-",
	"+",
	"~",
	"!",
}

// comparion operation
type ComparisonOp int

const (
	EQUAL            ComparisonOp = iota // =
	LESS_THAN                            // <
	LESS_THAN_EQUAL                      // <=
	GREAT_THAN                           // >
	GREAT_THAN_EQUAL                     // >=
	NOT_EQUAL                            // <>, !=
	IN                                   // IN
	NOT_IN                               // NOT IN
	LIKE                                 // LIKE
	NOT_LIKE                             // NOT LIKE
	ILIKE
	NOT_ILIKE
	REG_MATCH     // REG_MATCH
	NOT_REG_MATCH // NOT REG_MATCH
	IS_DISTINCT_FROM
	IS_NOT_DISTINCT_FROM
	NULL_SAFE_EQUAL // <=>
	//reference: https://dev.mysql.com/doc/refman/8.0/en/all-subqueries.html
	//subquery with ANY,SOME,ALL
	//operand comparison_operator [ANY | SOME | ALL] (subquery)
	ANY
	SOME
	ALL
)

func (op ComparisonOp) ToString() string {
	switch op {
	case EQUAL:
		return "="
	case LESS_THAN:
		return "<"
	case LESS_THAN_EQUAL:
		return "<="
	case GREAT_THAN:
		return ">"
	case GREAT_THAN_EQUAL:
		return ">="
	case NOT_EQUAL:
		return "!="
	case IN:
		return "in"
	case NOT_IN:
		return "not in"
	case LIKE:
		return "like"
	case NOT_LIKE:
		return "not like"
	case REG_MATCH:
		return "reg_match"
	case NOT_REG_MATCH:
		return "not reg_match"
	case IS_DISTINCT_FROM:
		return "is distinct from"
	case IS_NOT_DISTINCT_FROM:
		return "is not distinct from"
	case NULL_SAFE_EQUAL:
		return "<=>"
	case ANY:
		return "any"
	case SOME:
		return "some"
	case ALL:
		return "all"
	case ILIKE:
		return "ilike"
	case NOT_ILIKE:
		return "not ilike"
	default:
		return "Unknown ComparisonExprOperator"
	}
}

type ComparisonExpr struct {
	exprImpl
	Op ComparisonOp

	//ANY SOME ALL with subquery
	SubOp  ComparisonOp
	Left   Expr
	Right  Expr
	Escape Expr
}

func (node *ComparisonExpr) Format(ctx *FmtCtx) {
	if node.Left != nil {
		ctx.PrintExpr(node, node.Left, true)
		ctx.WriteByte(' ')
	}
	ctx.WriteString(node.Op.ToString())
	ctx.WriteByte(' ')

	if node.SubOp != ComparisonOp(0) {
		ctx.WriteString(node.SubOp.ToString())
		ctx.WriteByte(' ')
	}

	ctx.PrintExpr(node, node.Right, false)
	if node.Escape != nil {
		ctx.WriteString(" escape ")
		ctx.PrintExpr(node, node.Escape, true)
	}
}

// Accept implements NodeChecker interface
func (node *ComparisonExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}

	node = newNode.(*ComparisonExpr)
	tmpNode, ok := node.Left.Accept(v)
	if !ok {
		return node, false
	}
	node.Left = tmpNode

	tmpNode, ok = node.Right.Accept(v)
	if !ok {
		return node, false
	}
	node.Right = tmpNode
	return v.Exit(node)
}

func NewComparisonExpr(op ComparisonOp, left Expr, right Expr, buf *buffer.Buffer) *ComparisonExpr {
	var c *ComparisonExpr
	if buf != nil {
		c = buffer.Alloc[ComparisonExpr](buf)
	} else {
		c = new(ComparisonExpr)
	}
	c.Op = op
	c.SubOp = ComparisonOp(0)
	c.Left = left
	c.Right = right
	return c
}

func NewSubqueryComparisonExpr(op ComparisonOp, subOp ComparisonOp, left Expr, right Expr, buf *buffer.Buffer) *ComparisonExpr {
	c := buffer.Alloc[ComparisonExpr](buf)
	c.Op = op
	c.SubOp = subOp
	c.Left = left
	c.Right = right
	return c
}

func NewComparisonExprWithSubop(op ComparisonOp, subOp ComparisonOp, left Expr, right Expr, buf *buffer.Buffer) *ComparisonExpr {
	c := buffer.Alloc[ComparisonExpr](buf)
	c.Op = op
	c.SubOp = subOp
	c.Left = left
	c.Right = right
	return c
}

func NewComparisonExprWithEscape(op ComparisonOp, l, r, e Expr, buf *buffer.Buffer) *ComparisonExpr {
	c := buffer.Alloc[ComparisonExpr](buf)
	c.Op = op
	c.Left = l
	c.Right = r
	c.Escape = e
	return c

}

// and expression
type AndExpr struct {
	exprImpl
	Left, Right Expr
}

func (node *AndExpr) Format(ctx *FmtCtx) {
	ctx.PrintExpr(node, node.Left, true)
	ctx.WriteString(" and ")
	ctx.PrintExpr(node, node.Right, false)
}

func (node *AndExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}

	node = newNode.(*AndExpr)
	tmpNode, ok := node.Left.Accept(v)
	if !ok {
		return node, false
	}
	node.Left = tmpNode

	tmpNode, ok = node.Right.Accept(v)
	if !ok {
		return node, false
	}
	node.Right = tmpNode

	return v.Exit(node)
}

func NewAndExpr(left Expr, right Expr, buf *buffer.Buffer) *AndExpr {
	var a *AndExpr
	if buf != nil {
		a = buffer.Alloc[AndExpr](buf)
	} else {
		a = new(AndExpr)
	}
	a.Left = left
	a.Right = right
	return a
}

// xor expression
type XorExpr struct {
	exprImpl
	Left, Right Expr
}

func (node *XorExpr) Format(ctx *FmtCtx) {
	ctx.PrintExpr(node, node.Left, true)
	ctx.WriteString(" xor ")
	ctx.PrintExpr(node, node.Right, false)
}

func (node *XorExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}

	node = newNode.(*XorExpr)
	tmpNode, ok := node.Left.Accept(v)
	if !ok {
		return node, false
	}
	node.Left = tmpNode

	tmpNode, ok = node.Right.Accept(v)
	if !ok {
		return node, false
	}
	node.Right = tmpNode

	return v.Exit(node)
}

func NewXorExpr(left Expr, right Expr, buf *buffer.Buffer) *XorExpr {
	x := buffer.Alloc[XorExpr](buf)
	x.Left = left
	x.Right = right
	return x
}

// or expression
type OrExpr struct {
	exprImpl
	Left, Right Expr
}

func (node *OrExpr) Format(ctx *FmtCtx) {
	ctx.PrintExpr(node, node.Left, true)
	ctx.WriteString(" or ")
	ctx.PrintExpr(node, node.Right, false)
}

func (node *OrExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}

	node = newNode.(*OrExpr)
	tmpNode, ok := node.Left.Accept(v)
	if !ok {
		return node, false
	}
	node.Left = tmpNode

	tmpNode, ok = node.Right.Accept(v)
	if !ok {
		return node, false
	}
	node.Right = tmpNode

	return v.Exit(node)
}

func NewOrExpr(left Expr, right Expr, buf *buffer.Buffer) *OrExpr {
	var o *OrExpr
	if buf != nil {
		o = buffer.Alloc[OrExpr](buf)
	} else {
		o = new(OrExpr)
	}
	o.Left = left
	o.Right = right
	return o
}

// not expression
type NotExpr struct {
	exprImpl
	Expr Expr
}

func (node *NotExpr) Format(ctx *FmtCtx) {
	ctx.WriteString("not ")
	ctx.PrintExpr(node, node.Expr, true)
}

func (node *NotExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}

	node = newNode.(*NotExpr)
	tmpNode, ok := node.Expr.Accept(v)
	if !ok {
		return node, false
	}
	node.Expr = tmpNode
	return v.Exit(node)
}

func NewNotExpr(expr Expr, buf *buffer.Buffer) *NotExpr {
	n := buffer.Alloc[NotExpr](buf)
	n.Expr = expr
	return n
}

// is null expression
type IsNullExpr struct {
	exprImpl
	Expr Expr
}

func (node *IsNullExpr) Format(ctx *FmtCtx) {
	ctx.PrintExpr(node, node.Expr, true)
	ctx.WriteString(" is null")
}

// Accept implements NodeChecker interface
func (node *IsNullExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	node = newNode.(*IsNullExpr)
	tmpNode, ok := node.Expr.Accept(v)
	if !ok {
		return node, false
	}
	node.Expr = tmpNode
	return v.Exit(node)
}

func NewIsNullExpr(expr Expr, buf *buffer.Buffer) *IsNullExpr {
	var n *IsNullExpr
	if buf != nil {
		n = buffer.Alloc[IsNullExpr](buf)
	} else {
		n = new(IsNullExpr)
	}
	n.Expr = expr
	return n
}

// is not null expression
type IsNotNullExpr struct {
	exprImpl
	Expr Expr
}

func (node *IsNotNullExpr) Format(ctx *FmtCtx) {
	ctx.PrintExpr(node, node.Expr, true)
	ctx.WriteString(" is not null")
}

// Accept implements NodeChecker interface
func (node *IsNotNullExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	node = newNode.(*IsNotNullExpr)
	tmpNode, ok := node.Expr.Accept(v)
	if !ok {
		return node, false
	}
	node.Expr = tmpNode
	return v.Exit(node)
}

func NewIsNotNullExpr(expr Expr, buf *buffer.Buffer) *IsNotNullExpr {
	n := buffer.Alloc[IsNotNullExpr](buf)
	n.Expr = expr
	return n
}

// is unknown expression
type IsUnknownExpr struct {
	exprImpl
	Expr Expr
}

func (node *IsUnknownExpr) Format(ctx *FmtCtx) {
	ctx.PrintExpr(node, node.Expr, true)
	ctx.WriteString(" is unknown")
}

// Accept implements NodeChecker interface
func (node *IsUnknownExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	node = newNode.(*IsUnknownExpr)
	tmpNode, ok := node.Expr.Accept(v)
	if !ok {
		return node, false
	}
	node.Expr = tmpNode
	return v.Exit(node)
}

func NewIsUnknownExpr(expr Expr, buf *buffer.Buffer) *IsUnknownExpr {
	u := buffer.Alloc[IsUnknownExpr](buf)
	u.Expr = expr
	return u
}

// is not unknown expression
type IsNotUnknownExpr struct {
	exprImpl
	Expr Expr
}

func (node *IsNotUnknownExpr) Format(ctx *FmtCtx) {
	ctx.PrintExpr(node, node.Expr, true)
	ctx.WriteString(" is not unknown")
}

// Accept implements NodeChecker interface
func (node *IsNotUnknownExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	node = newNode.(*IsNotUnknownExpr)
	tmpNode, ok := node.Expr.Accept(v)
	if !ok {
		return node, false
	}
	node.Expr = tmpNode
	return v.Exit(node)
}

func NewIsNotUnknownExpr(expr Expr, buf *buffer.Buffer) *IsNotUnknownExpr {
	nu := buffer.Alloc[IsNotUnknownExpr](buf)
	nu.Expr = expr
	return nu
}

// is true expression
type IsTrueExpr struct {
	exprImpl
	Expr Expr
}

// Accept implements NodeChecker interface
func (node *IsTrueExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	node = newNode.(*IsTrueExpr)
	tmpNode, ok := node.Expr.Accept(v)
	if !ok {
		return node, false
	}
	node.Expr = tmpNode
	return v.Exit(node)
}

func (node *IsTrueExpr) Format(ctx *FmtCtx) {
	ctx.PrintExpr(node, node.Expr, true)
	ctx.WriteString(" is true")
}

func NewIsTrueExpr(e Expr, buf *buffer.Buffer) *IsTrueExpr {
	i := buffer.Alloc[IsTrueExpr](buf)
	i.Expr = e
	return i
}

// is not true expression
type IsNotTrueExpr struct {
	exprImpl
	Expr Expr
}

func (node *IsNotTrueExpr) Format(ctx *FmtCtx) {
	ctx.PrintExpr(node, node.Expr, true)
	ctx.WriteString(" is not true")
}

// Accept implements NodeChecker interface
func (node *IsNotTrueExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	node = newNode.(*IsNotTrueExpr)
	tmpNode, ok := node.Expr.Accept(v)
	if !ok {
		return node, false
	}
	node.Expr = tmpNode
	return v.Exit(node)
}

func NewIsNotTrueExpr(e Expr, buf *buffer.Buffer) *IsNotTrueExpr {
	in := buffer.Alloc[IsNotTrueExpr](buf)
	in.Expr = e
	return in
}

// is false expression
type IsFalseExpr struct {
	exprImpl
	Expr Expr
}

func (node *IsFalseExpr) Format(ctx *FmtCtx) {
	ctx.PrintExpr(node, node.Expr, true)
	ctx.WriteString(" is false")
}

// Accept implements NodeChecker interface
func (node *IsFalseExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	node = newNode.(*IsFalseExpr)
	tmpNode, ok := node.Expr.Accept(v)
	if !ok {
		return node, false
	}
	node.Expr = tmpNode
	return v.Exit(node)
}

func NewIsFalseExpr(e Expr, buf *buffer.Buffer) *IsFalseExpr {
	i := buffer.Alloc[IsFalseExpr](buf)
	i.Expr = e
	return i
}

// is not false expression
type IsNotFalseExpr struct {
	exprImpl
	Expr Expr
}

func (node *IsNotFalseExpr) Format(ctx *FmtCtx) {
	ctx.PrintExpr(node, node.Expr, true)
	ctx.WriteString(" is not false")
}

// Accept implements NodeChecker interface
func (node *IsNotFalseExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	node = newNode.(*IsNotFalseExpr)
	tmpNode, ok := node.Expr.Accept(v)
	if !ok {
		return node, false
	}
	node.Expr = tmpNode
	return v.Exit(node)
}

func NewIsNotFalseExpr(e Expr, buf *buffer.Buffer) *IsNotFalseExpr {
	inf := buffer.Alloc[IsNotFalseExpr](buf)
	inf.Expr = e
	return inf
}

// subquery interface
type SubqueryExpr interface {
	Expr
}

// subquery
type Subquery struct {
	SubqueryExpr

	Select SelectStatement
	Exists bool
}

func (node *Subquery) Format(ctx *FmtCtx) {
	if node.Exists {
		ctx.WriteString("exists ")
	}
	node.Select.Format(ctx)
}

func (node *Subquery) Accept(v Visitor) (Expr, bool) {
	panic("unimplement Subquery Accept")
}

func NewSubquery(s SelectStatement, e bool, buf *buffer.Buffer) *Subquery {
	var su *Subquery
	if buf != nil {
		su = buffer.Alloc[Subquery](buf)
	} else {
		su = new(Subquery)
	}
	su.Select = s
	su.Exists = e
	return su
}

// a list of expression.
type Exprs []Expr

func (node Exprs) Format(ctx *FmtCtx) {
	prefix := ""
	for _, n := range node {
		ctx.WriteString(prefix)
		n.Format(ctx)
		prefix = ", "
	}
}

// ast fir the list of expression
type ExprList struct {
	exprImpl
	Exprs Exprs
}

func (n *ExprList) Accept(v Visitor) (Expr, bool) {
	panic("unimplement ExprList Accept")
}

// the parenthesized expression.
type ParenExpr struct {
	exprImpl
	Expr Expr
}

func (node *ParenExpr) Format(ctx *FmtCtx) {
	ctx.WriteByte('(')
	node.Expr.Format(ctx)
	ctx.WriteByte(')')
}

func (node *ParenExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	node = newNode.(*ParenExpr)
	if node.Expr != nil {
		tmpNode, ok := node.Expr.Accept(v)
		if !ok {
			return node, false
		}
		node.Expr = tmpNode
	}
	return v.Exit(node)
}

func NewParenExpr(e Expr, buf *buffer.Buffer) *ParenExpr {
	var p *ParenExpr
	if buf != nil {
		p = buffer.Alloc[ParenExpr](buf)
	} else {
		p = new(ParenExpr)
	}
	p.Expr = e
	return p
}

type FuncType int

func (node *FuncType) ToString() string {
	switch *node {
	case FUNC_TYPE_DISTINCT:
		return "distinct"
	case FUNC_TYPE_ALL:
		return "all"
	case FUNC_TYPE_TABLE:
		return "table function"
	default:
		return "Unknown FuncType"
	}
}

const (
	FUNC_TYPE_DEFAULT FuncType = iota
	FUNC_TYPE_DISTINCT
	FUNC_TYPE_ALL
	FUNC_TYPE_TABLE
)

// AggType specifies the type of aggregation.
type AggType int

const (
	_ AggType = iota
	AGG_TYPE_GENERAL
)

// the common interface to UnresolvedName and QualifiedFunctionName.
type FunctionReference interface {
	fmt.Stringer
	NodeFormatter
}

var _ FunctionReference = &UnresolvedName{}

// function reference
type ResolvableFunctionReference struct {
	FunctionReference
}

func (node *ResolvableFunctionReference) Format(ctx *FmtCtx) {
	node.FunctionReference.(*UnresolvedName).Format(ctx)
}

func FuncName2ResolvableFunctionReference(funcName *UnresolvedName, buf *buffer.Buffer) ResolvableFunctionReference {
	if buf != nil {
		r := buffer.Alloc[ResolvableFunctionReference](buf)	
		r.FunctionReference = funcName
		return *r
	}
	return ResolvableFunctionReference{FunctionReference: funcName}
}

// function call expression
type FuncExpr struct {
	exprImpl
	Func  ResolvableFunctionReference
	Type  FuncType
	Exprs Exprs

	//specify the type of aggregation.
	AggType AggType

	WindowSpec *WindowSpec

	OrderBy OrderBy
}

func NewFuncExpr(fun ResolvableFunctionReference, expr Exprs, buf *buffer.Buffer) *FuncExpr {
	fx := buffer.Alloc[FuncExpr](buf)
	fx.Func = fun
	fx.Exprs = expr
	return fx
}

func NewFuncExprWithWinSpec(fun ResolvableFunctionReference, expr Exprs, typ FuncType, win *WindowSpec, buf *buffer.Buffer) *FuncExpr {
	fx := buffer.Alloc[FuncExpr](buf)
	fx.Func = fun
	fx.Exprs = expr
	fx.Type = typ
	fx.WindowSpec = win
	return fx
}

func (node *FuncExpr) Format(ctx *FmtCtx) {
	node.Func.Format(ctx)

	ctx.WriteString("(")
	if node.Type != FUNC_TYPE_DEFAULT && node.Type != FUNC_TYPE_TABLE {
		ctx.WriteString(node.Type.ToString())
		ctx.WriteByte(' ')
	}
	if node.Func.FunctionReference.(*UnresolvedName).Parts[0] == "trim" {
		trimExprsFormat(ctx, node.Exprs)
	} else {
		node.Exprs.Format(ctx)
	}

	if node.OrderBy != nil {
		node.OrderBy.Format(ctx)
	}

	ctx.WriteByte(')')

	if node.WindowSpec != nil {
		ctx.WriteString(" ")
		node.WindowSpec.Format(ctx)
	}
}

// Accept implements NodeChecker interface
func (node *FuncExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	node = newNode.(*FuncExpr)
	for i, val := range node.Exprs {
		tmpNode, ok := val.Accept(v)
		if !ok {
			return node, false
		}
		node.Exprs[i] = tmpNode
	}
	return v.Exit(node)
}

func trimExprsFormat(ctx *FmtCtx, exprs Exprs) {
	tp := exprs[0].(*NumVal).String()
	switch tp {
	case "0":
		exprs[3].Format(ctx)
	case "1":
		exprs[2].Format(ctx)
		ctx.WriteString(" from ")
		exprs[3].Format(ctx)
	case "2":
		exprs[1].Format(ctx)
		ctx.WriteString(" from ")
		exprs[3].Format(ctx)
	case "3":
		exprs[1].Format(ctx)
		ctx.WriteString(" ")
		exprs[2].Format(ctx)
		ctx.WriteString(" from ")
		exprs[3].Format(ctx)
	default:
		panic("unknown trim type")
	}
}

type WindowSpec struct {
	PartitionBy Exprs
	OrderBy     OrderBy
	HasFrame    bool
	Frame       *FrameClause
}

func NewWindowSpec(pb Exprs, o OrderBy, f *FrameClause, ha bool, buf *buffer.Buffer) *WindowSpec {
	wi := buffer.Alloc[WindowSpec](buf)
	wi.PartitionBy = pb
	wi.OrderBy = o
	wi.Frame = f
	wi.HasFrame = ha
	return wi
}

func (node *WindowSpec) Format(ctx *FmtCtx) {
	ctx.WriteString("over (")
	flag := false
	if len(node.PartitionBy) > 0 {
		ctx.WriteString("partition by ")
		node.PartitionBy.Format(ctx)
		flag = true
	}

	if len(node.OrderBy) > 0 {
		if flag {
			ctx.WriteString(" ")
		}
		node.OrderBy.Format(ctx)
		flag = true
	}

	if node.Frame != nil && node.HasFrame {
		if flag {
			ctx.WriteString(" ")
		}
		node.Frame.Format(ctx)
	}

	ctx.WriteByte(')')
}

type FrameType int

const (
	Rows FrameType = iota
	Range
	Groups
)

type FrameClause struct {
	Type   FrameType
	HasEnd bool
	Start  *FrameBound
	End    *FrameBound
}

func NewFrameClause(t FrameType, h bool, s, e *FrameBound, buf *buffer.Buffer) *FrameClause {
	fr := buffer.Alloc[FrameClause](buf)
	fr.Type = t
	fr.HasEnd = h
	fr.Start = s
	fr.End = e
	return fr
}

func (node *FrameClause) Format(ctx *FmtCtx) {
	switch node.Type {
	case Rows:
		ctx.WriteString("rows")
	case Range:
		ctx.WriteString("range")
	case Groups:
		ctx.WriteString("groups")
	}
	ctx.WriteString(" ")
	if !node.HasEnd {
		node.Start.Format(ctx)
		return
	}
	ctx.WriteString("between ")
	node.Start.Format(ctx)
	ctx.WriteString(" and ")
	node.End.Format(ctx)
}

type BoundType int

const (
	Following BoundType = iota
	Preceding
	CurrentRow
)

type FrameBound struct {
	Type      BoundType
	UnBounded bool
	Expr      Expr
}

func NewFrameBound(t BoundType, u bool, e Expr, buf *buffer.Buffer) *FrameBound {
	f := buffer.Alloc[FrameBound](buf)
	f.Type = t
	f.UnBounded = u
	f.Expr = e
	return f
}

func (node *FrameBound) Format(ctx *FmtCtx) {
	if node.UnBounded {
		ctx.WriteString("unbounded")
	}
	if node.Type == CurrentRow {
		ctx.WriteString("current row")
	} else {
		if node.Expr != nil {
			node.Expr.Format(ctx)
		}
		if node.Type == Preceding {
			ctx.WriteString(" preceding")
		} else {
			ctx.WriteString(" following")
		}
	}
}

// type reference
type ResolvableTypeReference interface {
}

var _ ResolvableTypeReference = &UnresolvedObjectName{}
var _ ResolvableTypeReference = &T{}

// the Cast expression
type CastExpr struct {
	exprImpl
	Expr Expr
	Type ResolvableTypeReference
}

func (node *CastExpr) Format(ctx *FmtCtx) {
	ctx.WriteString("cast(")
	node.Expr.Format(ctx)
	ctx.WriteString(" as ")
	node.Type.(*T).InternalType.Format(ctx)
	ctx.WriteByte(')')
}

// Accept implements NodeChecker interface
func (node *CastExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	node = newNode.(*CastExpr)
	tmpNode, ok := node.Expr.Accept(v)
	if !ok {
		return node, false
	}
	node.Expr = tmpNode
	return v.Exit(node)
}

func NewCastExpr(e Expr, t ResolvableTypeReference, buf *buffer.Buffer) *CastExpr {
	c := buffer.Alloc[CastExpr](buf)
	c.Expr = e
	c.Type = t
	return c
}

type BitCastExpr struct {
	exprImpl
	Expr Expr
	Type ResolvableTypeReference
}

func (node *BitCastExpr) Format(ctx *FmtCtx) {
	ctx.WriteString("bit_cast(")
	node.Expr.Format(ctx)
	ctx.WriteString(" as ")
	node.Type.(*T).InternalType.Format(ctx)
	ctx.WriteByte(')')
}

// Accept implements NodeChecker interface
func (node *BitCastExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	node = newNode.(*BitCastExpr)
	tmpNode, ok := node.Expr.Accept(v)
	if !ok {
		return node, false
	}
	node.Expr = tmpNode
	return v.Exit(node)
}

func NewBitCastExpr(e Expr, t ResolvableTypeReference, buf *buffer.Buffer) *BitCastExpr {
	bc := buffer.Alloc[BitCastExpr](buf)
	bc.Expr = e
	bc.Type = t
	return bc
}

// the parenthesized list of expressions.
type Tuple struct {
	exprImpl
	Exprs Exprs
}

func (node *Tuple) Format(ctx *FmtCtx) {
	if node.Exprs != nil {
		ctx.WriteByte('(')
		node.Exprs.Format(ctx)
		ctx.WriteByte(')')
	}
}

// Accept implements NodeChecker interface
func (node *Tuple) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	node = newNode.(*Tuple)
	for i, val := range node.Exprs {
		tmpNode, ok := val.Accept(v)
		if !ok {
			return node, false
		}
		node.Exprs[i] = tmpNode
	}
	return v.Exit(node)
}

func NewTuple(es Exprs, buf *buffer.Buffer) *Tuple {
	var t *Tuple
	if buf != nil {
		t = buffer.Alloc[Tuple](buf)
	} else {
		t = new(Tuple)
	}
	t.Exprs = es
	return t
}

// the BETWEEN or a NOT BETWEEN expression
type RangeCond struct {
	exprImpl
	Not      bool
	Left     Expr
	From, To Expr
}

func (node *RangeCond) Format(ctx *FmtCtx) {
	ctx.PrintExpr(node, node.Left, true)
	if node.Not {
		ctx.WriteString(" not")
	}
	ctx.WriteString(" between ")
	ctx.PrintExpr(node, node.From, true)
	ctx.WriteString(" and ")
	ctx.PrintExpr(node, node.To, false)
}

// Accept implements NodeChecker interface
func (node *RangeCond) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}

	node = newNode.(*RangeCond)
	tmpNode, ok := node.Left.Accept(v)
	if !ok {
		return node, false
	}
	node.Left = tmpNode

	tmpNode, ok = node.From.Accept(v)
	if !ok {
		return node, false
	}
	node.From = tmpNode

	tmpNode, ok = node.To.Accept(v)
	if !ok {
		return node, false
	}
	node.To = tmpNode

	return v.Exit(node)
}

func NewRangeCond(n bool, l, f, t Expr, buf *buffer.Buffer) *RangeCond {
	rc := buffer.Alloc[RangeCond](buf)
	rc.Not = n
	rc.Left = l
	rc.From = f
	rc.To = t
	return rc
}

// Case-When expression.
type CaseExpr struct {
	exprImpl
	Expr  Expr
	Whens []*When
	Else  Expr
}

func (node *CaseExpr) Format(ctx *FmtCtx) {
	ctx.WriteString("case")
	if node.Expr != nil {
		ctx.WriteByte(' ')
		node.Expr.Format(ctx)
	}
	ctx.WriteByte(' ')
	prefix := ""
	for _, w := range node.Whens {
		ctx.WriteString(prefix)
		w.Format(ctx)
		prefix = " "
	}
	if node.Else != nil {
		ctx.WriteString(" else ")
		node.Else.Format(ctx)
	}
	ctx.WriteString(" end")
}

// Accept implements NodeChecker interface
func (node *CaseExpr) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	node = newNode.(*CaseExpr)

	tmpNode, ok := node.Expr.Accept(v)
	if !ok {
		return node, false
	}
	node.Expr = tmpNode

	for _, when := range node.Whens {
		tmpNode, ok = when.Cond.Accept(v)
		if !ok {
			return node, false
		}
		when.Cond = tmpNode

		tmpNode, ok = when.Val.Accept(v)
		if !ok {
			return node, false
		}
		when.Val = tmpNode
	}

	tmpNode, ok = node.Else.Accept(v)
	if !ok {
		return node, false
	}
	node.Else = tmpNode

	return v.Exit(node)
}

func NewCaseExpr(e Expr, whens []*When, elsep Expr, buf *buffer.Buffer) *CaseExpr {
	c := buffer.Alloc[CaseExpr](buf)
	c.Expr = e
	c.Whens = whens
	c.Else = elsep
	return c
}

// When sub-expression.
type When struct {
	Cond Expr
	Val  Expr
}

func (node *When) Format(ctx *FmtCtx) {
	ctx.WriteString("when ")
	node.Cond.Format(ctx)
	ctx.WriteString(" then ")
	node.Val.Format(ctx)
}

func NewWhen(c Expr, v Expr, buf *buffer.Buffer) *When {
	w := buffer.Alloc[When](buf)
	w.Cond = c
	w.Val = v
	return w
}

// IntervalType is the type for time and timestamp units.
type IntervalType int

func (node *IntervalType) ToString() string {
	switch *node {
	case INTERVAL_TYPE_SECOND:
		return "second"
	default:
		return "Unknown IntervalType"
	}
}

const (
	//an invalid time or timestamp unit
	INTERVAL_TYPE_INVALID IntervalType = iota
	//the time or timestamp unit MICROSECOND.
	INTERVAL_TYPE_MICROSECOND
	//the time or timestamp unit SECOND.
	INTERVAL_TYPE_SECOND
	//the time or timestamp unit MINUTE.
	INTERVAL_TYPE_MINUTE
	//the time or timestamp unit HOUR.
	INTERVAL_TYPE_HOUR
	//the time or timestamp unit DAY.
	INTERVAL_TYPE_DAY
	//the time or timestamp unit WEEK.
	INTERVAL_TYPE_WEEK
	//the time or timestamp unit MONTH.
	INTERVAL_TYPE_MONTH
	//the time or timestamp unit QUARTER.
	INTERVAL_TYPE_QUARTER
	//the time or timestamp unit YEAR.
	INTERVAL_TYPE_YEAR
	//the time unit SECOND_MICROSECOND.
	INTERVAL_TYPE_SECOND_MICROSECOND
	//the time unit MINUTE_MICROSECOND.
	INTERVAL_TYPE_MINUTE_MICROSECOND
	//the time unit MINUTE_SECOND.
	INTERVAL_TYPE_MINUTE_SECOND
	//the time unit HOUR_MICROSECOND.
	INTERVAL_TYPE_HOUR_MICROSECOND
	//the time unit HOUR_SECOND.
	INTERVAL_TYPE_HOUR_SECOND
	//the time unit HOUR_MINUTE.
	INTERVAL_TYPE_HOUR_MINUTE
	//the time unit DAY_MICROSECOND.
	INTERVAL_TYPE_DAY_MICROSECOND
	//the time unit DAY_SECOND.
	INTERVAL_TYPE_DAY_SECOND
	//the time unit DAY_MINUTE.
	INTERVAL_TYPE_DAYMINUTE
	//the time unit DAY_HOUR.
	INTERVAL_TYPE_DAYHOUR
	//the time unit YEAR_MONTH.
	INTERVAL_TYPE_YEARMONTH
)

// INTERVAL / time unit
type IntervalExpr struct {
	exprImpl
	Expr Expr
	Type IntervalType
}

func (node *IntervalExpr) Format(ctx *FmtCtx) {
	ctx.WriteString("interval")
	if node.Expr != nil {
		ctx.WriteByte(' ')
		node.Expr.Format(ctx)
	}
	if node.Type != INTERVAL_TYPE_INVALID {
		ctx.WriteByte(' ')
		ctx.WriteString(node.Type.ToString())
	}
}

// Accept implements NodeChecker Accept interface.
func (node *IntervalExpr) Accept(v Visitor) (Expr, bool) {
	// TODO:
	panic("unimplement interval expr Accept")
}

func NewIntervalExpr(it IntervalType, buf *buffer.Buffer) *IntervalExpr {
	i := buffer.Alloc[IntervalExpr](buf)
	i.Type = it
	return i
}

// the DEFAULT expression.
type DefaultVal struct {
	exprImpl
	Expr Expr
}

func (node *DefaultVal) Format(ctx *FmtCtx) {
	ctx.WriteString("default")
	if node.Expr != nil {
		node.Expr.Format(ctx)
	}
}

// Accept implements NodeChecker interface
func (node *DefaultVal) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	node = newNode.(*DefaultVal)
	tmpNode, ok := node.Expr.Accept(v)
	if !ok {
		return node, false
	}
	node.Expr = tmpNode
	return v.Exit(node)
}

func NewDefaultVal(buf *buffer.Buffer) *DefaultVal {
	d := buffer.Alloc[DefaultVal](buf)
	return d
}

type UpdateVal struct {
	exprImpl
}

func (node *UpdateVal) Format(ctx *FmtCtx) {}

// Accept implements NodeChecker interface
func (node *UpdateVal) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	return v.Exit(node)
}

type TypeExpr interface {
	Expr
}

/*
Variable Expression Used in Set Statement,
Load Data statement, Show statement,etc.
Variable types:
User-Defined Variable
Local-Variable: DECLARE statement
System Variable: Global System Variable, Session System Variable
*/

type VarExpr struct {
	exprImpl
	Name   string
	System bool
	Global bool
	Expr   Expr
}

// incomplete
func (node *VarExpr) Format(ctx *FmtCtx) {
	if node.Name != "" {
		ctx.WriteByte('@')
		if node.System {
			ctx.WriteByte('@')
		}
		ctx.WriteString(node.Name)
	}
}

// Accept implements NodeChecker Accept interface.
func (node *VarExpr) Accept(v Visitor) (Expr, bool) {
	panic("unimplement VarExpr Accept")
}

func NewVarExpr(n string, sys, glob bool, buf *buffer.Buffer) *VarExpr {
	v := buffer.Alloc[VarExpr](buf)
	v.Name = n
	v.System = sys
	v.Global = glob
	return v
}

// select a from t1 where a > ?
type ParamExpr struct {
	exprImpl
	Offset int
}

func (node *ParamExpr) Format(ctx *FmtCtx) {
	ctx.WriteByte('?')
}

// Accept implements NodeChecker Accept interface.
func (node *ParamExpr) Accept(v Visitor) (Expr, bool) {
	panic("unimplement ParamExpr Accept")
}

func NewParamExpr(offset int, buf *buffer.Buffer) *ParamExpr {
	p := buffer.Alloc[ParamExpr](buf)
	p.Offset = offset
	return p
}

type MaxValue struct {
	exprImpl
}

func (node *MaxValue) Format(ctx *FmtCtx) {
	ctx.WriteString("MAXVALUE")
}

func NewMaxValue(buf *buffer.Buffer) *MaxValue {
	m := buffer.Alloc[MaxValue](buf)
	return m
}

// Accept implements NodeChecker interface
func (node *MaxValue) Accept(v Visitor) (Expr, bool) {
	newNode, skipChildren := v.Enter(node)
	if skipChildren {
		return v.Exit(newNode)
	}
	return v.Exit(node)
}
