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
	"strconv"
	"strings"

	"github.com/matrixorigin/matrixone/pkg/common/buffer"
)

type Explain interface {
	Statement
}

type explainImpl struct {
	Explain
	Statement Statement
	Format    string
	Options   []OptionElem
}

// EXPLAIN stmt statement
type ExplainStmt struct {
	explainImpl
}

func (node *ExplainStmt) Format(ctx *FmtCtx) {
	ctx.WriteString("explain")
	if node.Options != nil && len(node.Options) > 0 {
		ctx.WriteString(" (")
		var temp string
		for _, v := range node.Options {
			temp += v.Name
			if v.Value != "NULL" {
				temp += " " + v.Value
			}
			temp += ","
		}
		ctx.WriteString(temp[:len(temp)-1] + ")")
	}

	stmt := node.explainImpl.Statement
	switch st := stmt.(type) {
	case *ShowColumns:
		if st.Table != nil {
			ctx.WriteByte(' ')
			st.Table.ToTableName().Format(ctx)
		}
		if st.ColName != nil {
			ctx.WriteByte(' ')
			st.ColName.Format(ctx)
		}
	default:
		if stmt != nil {
			ctx.WriteByte(' ')
			stmt.Format(ctx)
		}
	}
}

func (node *ExplainStmt) GetStatementType() string { return "Explain" }
func (node *ExplainStmt) GetQueryType() string     { return QueryTypeOth }

func NewExplainStmt(stmt Statement, format string, buf *buffer.Buffer) *ExplainStmt {
	e := buffer.Alloc[ExplainStmt](buf)
	e.explainImpl.Statement = stmt
	e.explainImpl.Format = format
	return e
}

// EXPLAIN ANALYZE statement
type ExplainAnalyze struct {
	explainImpl
}

func (node *ExplainAnalyze) Format(ctx *FmtCtx) {
	ctx.WriteString("explain")
	if node.Options != nil && len(node.Options) > 0 {
		ctx.WriteString(" (")
		var temp string
		for _, v := range node.Options {
			temp += v.Name
			if v.Value != "NULL" {
				temp += " " + v.Value
			}
			temp += ","
		}
		ctx.WriteString(temp[:len(temp)-1] + ")")
	}

	stmt := node.explainImpl.Statement
	switch st := stmt.(type) {
	case *ShowColumns:
		if st.Table != nil {
			ctx.WriteByte(' ')
			st.Table.ToTableName().Format(ctx)
		}
		if st.ColName != nil {
			ctx.WriteByte(' ')
			st.ColName.Format(ctx)
		}
	default:
		if stmt != nil {
			ctx.WriteByte(' ')
			stmt.Format(ctx)
		}
	}
}

func (node *ExplainAnalyze) GetStatementType() string { return "Explain Analyze" }
func (node *ExplainAnalyze) GetQueryType() string     { return QueryTypeOth }

func NewExplainAnalyze(stmt Statement, format string, buf *buffer.Buffer) *ExplainAnalyze {
	e := buffer.Alloc[ExplainAnalyze](buf)
	e.explainImpl.Statement = stmt
	e.explainImpl.Format = format
	return e
}

// EXPLAIN FOR CONNECTION statement
type ExplainFor struct {
	explainImpl
	ID uint64
}

func (node *ExplainFor) Format(ctx *FmtCtx) {
	ctx.WriteString("explain format = ")
	ctx.WriteString(node.explainImpl.Format)
	ctx.WriteString(" for connection ")
	ctx.WriteString(strconv.FormatInt(int64(node.ID), 10))
}

func (node *ExplainFor) GetStatementType() string { return "Explain Format" }
func (node *ExplainFor) GetQueryType() string     { return QueryTypeOth }

func NewExplainFor(format string, id uint64, buf *buffer.Buffer) *ExplainFor {
	e := buffer.Alloc[ExplainFor](buf)
	e.explainImpl.Format = format
	e.ID = id
	return e
}

type OptionElem struct {
	Name  string
	Value string
}

func MakeOptionElem(name string, value string, buf *buffer.Buffer) OptionElem {
	o := buffer.Alloc[OptionElem](buf)
	o.Name = name
	o.Value = value
	return *o
}

func MakeOptions(elem OptionElem, buf *buffer.Buffer) []OptionElem {
	os := buffer.MakeSpecificSlice[OptionElem](buf, 0, 1)	
	os = buffer.AppendSlice[OptionElem](buf, os, elem)
	return os
}

func IsContainAnalyze(options []OptionElem) bool {
	if len(options) > 0 {
		for _, option := range options {
			if strings.EqualFold(option.Name, "analyze") && strings.EqualFold(option.Value, "true") {
				return true
			}
		}
	}
	return false
}
