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

type CreateStage struct {
	statementImpl
	IfNotExists bool
	Name        *BufIdentifier
	Url         string
	Credentials StageCredentials
	Status      StageStatus
	Comment     StageComment
}

func NewCreateStage(ifs bool, name Identifier, url string, cs StageCredentials, st StageStatus, co StageComment, buf *buffer.Buffer) *CreateStage {
	cr := buffer.Alloc[CreateStage](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)

	cr.IfNotExists = ifs
	cr.Name = n
	cr.Url = url
	cr.Credentials = cs
	cr.Status = st
	cr.Comment = co
	return cr
}

func (node *CreateStage) Format(ctx *FmtCtx) {
	ctx.WriteString("create stage ")
	if node.IfNotExists {
		ctx.WriteString("if not exists ")
	}
	node.Name.Format(ctx)

	ctx.WriteString(" url=")
	ctx.WriteString(fmt.Sprintf("'%s'", node.Url))

	node.Credentials.Format(ctx)
	node.Status.Format(ctx)
	node.Comment.Format(ctx)
}

func (node *CreateStage) GetStatementType() string { return "Create Stage" }
func (node *CreateStage) GetQueryType() string     { return QueryTypeOth }

type DropStage struct {
	statementImpl
	IfNotExists bool
	Name        *BufIdentifier
}

func NewDropStage(i bool, name Identifier, buf *buffer.Buffer) *DropStage {
	d := buffer.Alloc[DropStage](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)
	d.Name = n
	d.IfNotExists = i
	return d
}

func (node *DropStage) Format(ctx *FmtCtx) {
	ctx.WriteString("drop stage ")
	if node.IfNotExists {
		ctx.WriteString("if not exists ")
	}
	node.Name.Format(ctx)
}

func (node *DropStage) GetStatementType() string { return "Drop Stage" }
func (node *DropStage) GetQueryType() string     { return QueryTypeOth }

type AlterStage struct {
	statementImpl
	IfNotExists       bool
	Name              *BufIdentifier
	UrlOption         StageUrl
	CredentialsOption StageCredentials
	StatusOption      StageStatus
	Comment           StageComment
}

func NewAlterStage(ifn bool, name Identifier, url StageUrl, cop StageCredentials, sop StageStatus, co StageComment, buf *buffer.Buffer) *AlterStage {
	al := buffer.Alloc[AlterStage](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)
	al.IfNotExists = ifn
	al.Name = n
	al.UrlOption = url
	al.CredentialsOption = cop
	al.StatusOption = sop
	al.Comment = co
	return al
}

func (node *AlterStage) Format(ctx *FmtCtx) {
	ctx.WriteString("alter stage ")
	if node.IfNotExists {
		ctx.WriteString("if not exists ")
	}
	node.Name.Format(ctx)
	ctx.WriteString(" set ")
	node.UrlOption.Format(ctx)
	node.CredentialsOption.Format(ctx)
	node.StatusOption.Format(ctx)
	node.Comment.Format(ctx)
}

func (node *AlterStage) GetStatementType() string { return "Alter Stage" }
func (node *AlterStage) GetQueryType() string     { return QueryTypeOth }

type StageStatusOption int

const (
	StageStatusEnabled StageStatusOption = iota
	StageStatusDisabled
)

func (sso StageStatusOption) String() string {
	switch sso {
	case StageStatusEnabled:
		return "enabled"
	case StageStatusDisabled:
		return "disabled"
	default:
		return "disabled"
	}
}

type StageStatus struct {
	Exist  bool
	Option StageStatusOption
}

func NewStageStatus(e bool, o StageStatusOption, buf *buffer.Buffer) *StageStatus {
	s := buffer.Alloc[StageStatus](buf)
	s.Exist = e
	s.Option = o
	return s
}

func (node *StageStatus) Format(ctx *FmtCtx) {
	if node.Exist {
		switch node.Option {
		case StageStatusEnabled:
			ctx.WriteString(" enabled")
		case StageStatusDisabled:
			ctx.WriteString(" disabled")
		}
	}
}

type StageComment struct {
	Exist   bool
	Comment string
}

func NewStageComment(e bool, c string, buf *buffer.Buffer) *StageComment {
	s := buffer.Alloc[StageComment](buf)
	s.Exist = e
	s.Comment = c
	return s
}

func (node *StageComment) Format(ctx *FmtCtx) {
	if node.Exist {
		ctx.WriteString(" comment ")
		ctx.WriteString(fmt.Sprintf("'%s'", node.Comment))
	}
}

type StageCredentials struct {
	Exist       bool
	Credentials []string
}

func NewStageCredentials(e bool, cs []string, buf *buffer.Buffer) *StageCredentials {
	s := buffer.Alloc[StageCredentials](buf)
	s.Exist = e
	s.Credentials = cs
	return s
}

func (node *StageCredentials) Format(ctx *FmtCtx) {
	if node.Exist {

		ctx.WriteString(" crentiasl=")
		ctx.WriteString("{")
		for i := 0; i < len(node.Credentials)-1; i += 2 {
			ctx.WriteString(fmt.Sprintf("'%s'", node.Credentials[i]))
			ctx.WriteString("=")
			ctx.WriteString(fmt.Sprintf("'%s'", node.Credentials[i+1]))
			if i != len(node.Credentials)-2 {
				ctx.WriteString(",")
			}
		}
		ctx.WriteString("}")
	}
}

type StageUrl struct {
	Exist bool
	Url   string
}

func NewStageUrl(e bool, u string, buf *buffer.Buffer) *StageUrl {
	s := buffer.Alloc[StageUrl](buf)
	s.Exist = e
	s.Url = u
	return s
}

func (node *StageUrl) Format(ctx *FmtCtx) {
	if node.Exist {
		ctx.WriteString(" url=")
		ctx.WriteString(fmt.Sprintf("'%s'", node.Url))
	}
}

type ShowStages struct {
	showImpl
	Like *ComparisonExpr
}

func NewShowStages(l *ComparisonExpr, buf *buffer.Buffer) *ShowStages {
	s := buffer.Alloc[ShowStages](buf)
	s.Like = l
	return s
}

func (node *ShowStages) Format(ctx *FmtCtx) {
	ctx.WriteString("show stages ")
	if node.Like != nil {
		node.Like.Format(ctx)
	}
}
func (node *ShowStages) GetStatementType() string { return "Show Stages" }
func (node *ShowStages) GetQueryType() string     { return QueryTypeOth }
