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

type GrantType int

const (
	GrantTypePrivilege GrantType = iota
	GrantTypeRole
	GrantTypeProxy
)

type Grant struct {
	statementImpl
	Typ            GrantType
	GrantPrivilege *GrantPrivilege
	GrantRole      *GrantRole
	GrantProxy     *GrantProxy
}

func (node *Grant) Format(ctx *FmtCtx) {
	switch node.Typ {
	case GrantTypePrivilege:
		node.GrantPrivilege.Format(ctx)
	case GrantTypeRole:
		node.GrantRole.Format(ctx)
	case GrantTypeProxy:
		node.GrantProxy.Format(ctx)
	}
}

func (node *Grant) GetStatementType() string { return "Grant" }
func (node *Grant) GetQueryType() string     { return QueryTypeDCL }

func NewGrant(typ GrantType, grantPrivilege *GrantPrivilege, grantRole *GrantRole, grantProxy *GrantProxy, buf *buffer.Buffer) *Grant {
	g := buffer.Alloc[Grant](buf)
	g.Typ = typ
	g.GrantPrivilege = grantPrivilege
	g.GrantRole = grantRole
	g.GrantProxy = grantProxy
	return g
}

type GrantPrivilege struct {
	statementImpl
	Privileges []*Privilege
	//grant privileges
	ObjType ObjectType
	//grant privileges
	Level       *PrivilegeLevel
	Roles       []*Role
	GrantOption bool
}

func NewGrantPrivilege(pr []*Privilege, ob ObjectType, le *PrivilegeLevel, ro []*Role, gr bool, buf *buffer.Buffer) *GrantPrivilege {
	grant := buffer.Alloc[GrantPrivilege](buf)
	grant.Privileges = pr
	grant.ObjType = ob
	grant.Level = le
	grant.Roles = ro
	grant.GrantOption = gr
	return grant
}

func (node *GrantPrivilege) Format(ctx *FmtCtx) {
	ctx.WriteString("grant")
	if node.Privileges != nil {
		prefix := " "
		for _, p := range node.Privileges {
			ctx.WriteString(prefix)
			p.Format(ctx)
			prefix = ", "
		}
	}
	ctx.WriteString(" on")
	if node.ObjType != OBJECT_TYPE_NONE {
		ctx.WriteByte(' ')
		ctx.WriteString(node.ObjType.String())
	}
	if node.Level != nil {
		ctx.WriteByte(' ')
		node.Level.Format(ctx)
	}

	if node.Roles != nil {
		ctx.WriteString(" to")
		prefix := " "
		for _, r := range node.Roles {
			ctx.WriteString(prefix)
			r.Format(ctx)
			prefix = ", "
		}
	}
	if node.GrantOption {
		ctx.WriteString(" with grant option")
	}
}

func (node *GrantPrivilege) GetStatementType() string { return "Grant Privilege" }
func (node *GrantPrivilege) GetQueryType() string     { return QueryTypeDCL }

type GrantRole struct {
	statementImpl
	Roles       []*Role
	Users       []*User
	GrantOption bool
}

func NewGrantRole(roles []*Role, users []*User, grantoption bool, buf *buffer.Buffer) *GrantRole {
	g := buffer.Alloc[GrantRole](buf)
	g.Roles = roles
	g.Users = users
	g.GrantOption = grantoption
	return g
}

func (node *GrantRole) Format(ctx *FmtCtx) {
	ctx.WriteString("grant")
	if node.Roles != nil {
		prefix := " "
		for _, r := range node.Roles {
			ctx.WriteString(prefix)
			r.Format(ctx)
			prefix = ", "
		}
	}
	if node.Users != nil {
		ctx.WriteString(" to")
		prefix := " "
		for _, r := range node.Users {
			ctx.WriteString(prefix)
			r.Format(ctx)
			prefix = ", "
		}
	}
	if node.GrantOption {
		ctx.WriteString(" with grant option")
	}
}

func (node *GrantRole) GetStatementType() string { return "Grant Role" }
func (node *GrantRole) GetQueryType() string     { return QueryTypeDCL }

type GrantProxy struct {
	statementImpl
	ProxyUser   *User
	Users       []*User
	GrantOption bool
}

func NewGrantProxy(r *User, users []*User, grantoption bool, buf *buffer.Buffer) *GrantProxy {
	g := buffer.Alloc[GrantProxy](buf)
	g.ProxyUser = r
	g.Users = users
	g.GrantOption = grantoption
	return g
}

func (node *GrantProxy) Format(ctx *FmtCtx) {
	ctx.WriteString("grant")
	if node.ProxyUser != nil {
		ctx.WriteString(" proxy on ")
		node.ProxyUser.Format(ctx)
	}
	if node.Users != nil {
		ctx.WriteString(" to")
		prefix := " "
		for _, r := range node.Users {
			ctx.WriteString(prefix)
			r.Format(ctx)
			prefix = ", "
		}
	}
	if node.GrantOption {
		ctx.WriteString(" with grant option")
	}
}

func (node *GrantProxy) GetStatementType() string { return "Grant Proxy" }
func (node *GrantProxy) GetQueryType() string     { return QueryTypeDCL }
