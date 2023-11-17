// Copyright 2022 Matrix Origin
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

type Reset struct {
	Statement
	Name *BufIdentifier
}

func (node *Reset) Format(ctx *FmtCtx) {
	ctx.WriteString("reset")
	ctx.WriteString(" prepare ")
	node.Name.Format(ctx)
}

func (node *Reset) GetStatementType() string { return "Reset" }
func (node *Reset) GetQueryType() string     { return QueryTypeDCL }

func NewReset(name Identifier, buf *buffer.Buffer) *Reset {
	rs := buffer.Alloc[Reset](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)

	rs.Name = n
	return rs
}
