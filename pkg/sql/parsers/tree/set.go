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

type SetVar struct {
	statementImpl
	Assignments []*VarAssignmentExpr
}

func (node *SetVar) Format(ctx *FmtCtx) {
	ctx.WriteString("set")
	if node.Assignments != nil {
		prefix := " "
		for _, a := range node.Assignments {
			ctx.WriteString(prefix)
			a.Format(ctx)
			prefix = ", "
		}
	}
}

// Accept implements NodeChecker interface.
func (node *SetVar) Accept(v Visitor) (Expr, bool) {
	//TODO: unimplement Accept interface
	panic("tree.SetVar Unimplement Accept")
}

func (node *SetVar) GetStatementType() string { return "Set Var" }
func (node *SetVar) GetQueryType() string     { return QueryTypeOth }

func NewSetVar(a []*VarAssignmentExpr, buf *buffer.Buffer) *SetVar {
	sv := buffer.Alloc[SetVar](buf)
	sv.Assignments = a
	return sv
}

// for variable = expr
type VarAssignmentExpr struct {
	NodeFormatter
	System   bool
	Global   bool
	Name     *BufString
	Value    Expr
	Reserved Expr
}

func (node *VarAssignmentExpr) Format(ctx *FmtCtx) {
	if node.Global {
		ctx.WriteString("global ")
	}
	ctx.WriteString(node.Name.Get())
	ctx.WriteString(" =")
	if node.Value != nil {
		ctx.WriteByte(' ')
		node.Value.Format(ctx)
	}
	if node.Reserved != nil {
		ctx.WriteByte(' ')
		node.Reserved.Format(ctx)
	}
}

func NewVarAssignmentExpr(s bool, g bool, n string, v Expr, r Expr, buf *buffer.Buffer) *VarAssignmentExpr {
	va := buffer.Alloc[VarAssignmentExpr](buf)
	va.System = s
	va.Global = g
	bn := NewBufString(n)
	buf.Pin(bn)
	va.Name = bn
	va.Value = v
	va.Reserved = r
	return va
}

type SetDefaultRoleType int

const (
	SET_DEFAULT_ROLE_TYPE_NONE SetDefaultRoleType = iota
	SET_DEFAULT_ROLE_TYPE_ALL
	SET_DEFAULT_ROLE_TYPE_NORMAL
)

type SetDefaultRole struct {
	statementImpl
	Type  SetDefaultRoleType
	Roles []*Role
	Users []*User
}

func (node *SetDefaultRole) Format(ctx *FmtCtx) {
	ctx.WriteString("set default role")
	switch node.Type {
	case SET_DEFAULT_ROLE_TYPE_NONE:
		ctx.WriteString(" none")
	case SET_DEFAULT_ROLE_TYPE_ALL:
		ctx.WriteString(" all")
	case SET_DEFAULT_ROLE_TYPE_NORMAL:
		prefix := " "
		for _, r := range node.Roles {
			ctx.WriteString(prefix)
			r.Format(ctx)
			prefix = ", "
		}
	}
	ctx.WriteString(" to")
	prefix := " "
	for _, u := range node.Users {
		ctx.WriteString(prefix)
		u.Format(ctx)
		prefix = ", "
	}
}

func (node *SetDefaultRole) GetStatementType() string { return "Set Role" }
func (node *SetDefaultRole) GetQueryType() string     { return QueryTypeOth }

func NewSetDefaultRole(t SetDefaultRoleType, r []*Role, u []*User, buf *buffer.Buffer) *SetDefaultRole {
	sdr := buffer.Alloc[SetDefaultRole](buf)
	sdr.Type = t
	sdr.Roles = r
	sdr.Users = u
	return sdr
}

type SetRoleType int

const (
	SET_ROLE_TYPE_NORMAL SetRoleType = iota
	SET_ROLE_TYPE_DEFAULT
	SET_ROLE_TYPE_NONE
	SET_ROLE_TYPE_ALL
	SET_ROLE_TYPE_ALL_EXCEPT
)

type SetRole struct {
	statementImpl
	SecondaryRole     bool
	SecondaryRoleType SecondaryRoleType
	Role              *Role
}

func NewSetRole(secondaryRole bool, secondaryRoleTyp SecondaryRoleType, role *Role, buf *buffer.Buffer) *SetRole {
	sr := buffer.Alloc[SetRole](buf)
	sr.SecondaryRole = secondaryRole
	sr.SecondaryRoleType = secondaryRoleTyp
	sr.Role = role
	return sr
}

func (node *SetRole) Format(ctx *FmtCtx) {
	ctx.WriteString("set")
	if !node.SecondaryRole {
		if node.Role != nil {
			ctx.WriteString(" role ")
			node.Role.Format(ctx)
		}
	} else {
		ctx.WriteString(" secondary role ")
		switch node.SecondaryRoleType {
		case SecondaryRoleTypeAll:
			ctx.WriteString("all")
		case SecondaryRoleTypeNone:
			ctx.WriteString("none")
		}
	}
}

func (node *SetRole) GetStatementType() string { return "Set Role" }
func (node *SetRole) GetQueryType() string     { return QueryTypeOth }

type SetPassword struct {
	statementImpl
	User     *User
	Password *BufString
}

func (node *SetPassword) Format(ctx *FmtCtx) {
	ctx.WriteString("set password")
	if node.User != nil {
		ctx.WriteString(" for ")
		node.User.Format(ctx)
	}
	ctx.WriteString(" = ")
	ctx.WriteString(node.Password.Get())
}
func NewSetPassword(u *User, p string, buf *buffer.Buffer) *SetPassword {
	sp := buffer.Alloc[SetPassword](buf)
	sp.User = u
	bp := NewBufString(p)
	buf.Pin(bp)
	sp.Password = bp
	return sp
}

func (node *SetPassword) GetStatementType() string { return "Set Password" }
func (node *SetPassword) GetQueryType() string     { return QueryTypeOth }

type IsolationLevelType int

const (
	ISOLATION_LEVEL_NONE IsolationLevelType = iota
	ISOLATION_LEVEL_REPEATABLE_READ
	ISOLATION_LEVEL_READ_COMMITTED
	ISOLATION_LEVEL_READ_UNCOMMITTED
	ISOLATION_LEVEL_SERIALIZABLE
)

func (ilt IsolationLevelType) String() string {
	switch ilt {
	case ISOLATION_LEVEL_NONE:
		return "isolation level none"
	case ISOLATION_LEVEL_REPEATABLE_READ:
		return "isolation level repeatable read"
	case ISOLATION_LEVEL_READ_COMMITTED:
		return "isolation level read committed"
	case ISOLATION_LEVEL_READ_UNCOMMITTED:
		return "isolation level read uncommitted"
	case ISOLATION_LEVEL_SERIALIZABLE:
		return "isolation level serializable"
	default:
		return "isolation level unknown"
	}
}

type AccessModeType int

const (
	ACCESS_MODE_NONE AccessModeType = iota
	ACCESS_MODE_READ_WRITE
	ACCESS_MODE_READ_ONLY
)

func (amt AccessModeType) String() string {
	switch amt {
	case ACCESS_MODE_NONE:
		return "none"
	case ACCESS_MODE_READ_WRITE:
		return "read write"
	case ACCESS_MODE_READ_ONLY:
		return "read only"
	default:
		return "unknown"
	}
}

type TransactionCharacteristic struct {
	IsLevel   bool
	Isolation IsolationLevelType
	Access    AccessModeType
}

func NewTransactionCharacteristic(islevel bool, isolation IsolationLevelType, access AccessModeType, buf *buffer.Buffer) *TransactionCharacteristic {
	t := buffer.Alloc[TransactionCharacteristic](buf)
	t.IsLevel = islevel
	t.Isolation = isolation
	t.Access = access
	return t
}

func (tc *TransactionCharacteristic) Format(ctx *FmtCtx) {
	if tc.IsLevel {
		ctx.WriteString(tc.Isolation.String())
	} else {
		ctx.WriteString(tc.Access.String())
	}
}

type SetTransaction struct {
	statementImpl
	Global        bool
	CharacterList []*TransactionCharacteristic
}

func NewSetTransaction(glo bool, charlist []*TransactionCharacteristic, buf *buffer.Buffer) *SetTransaction {
	se := buffer.Alloc[SetTransaction](buf)
	se.Global = glo
	se.CharacterList = charlist
	return se
}

func (node *SetTransaction) Format(ctx *FmtCtx) {
	ctx.WriteString("set")
	if node.Global {
		ctx.WriteString(" global")
	}
	ctx.WriteString(" transaction ")

	for i, c := range node.CharacterList {
		if i > 0 {
			ctx.WriteString(" , ")
		}
		c.Format(ctx)
	}
}

func (node *SetTransaction) GetStatementType() string { return "Set Transaction" }
func (node *SetTransaction) GetQueryType() string     { return QueryTypeTCL }
