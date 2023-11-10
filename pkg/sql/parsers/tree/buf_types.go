// Copyright 2023 Matrix Origin
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

import "go/constant"

// TODO: need comment here
type BufString struct {
	s *string
}

func NewBufString(s string) *BufString {
	return &BufString{&s}
}

func (b *BufString) Get() string {
	if b != nil && b.s != nil {
		return *b.s
	}
	return ""
}

func (b *BufString) Set(s string) {
	if b != nil {
		b.s = &s
	}
}

type BufConstant struct {
	v *constant.Value
}

func NewBufConstant(value constant.Value) *BufConstant {
	return &BufConstant{&value}
}

func (b *BufConstant) Get() constant.Value {
	if b != nil && b.v != nil {
		return *b.v
	}
	return nil
}

type bufIdentifier interface {
	string | Identifier
}

func NewBufIdentifier[T string | Identifier](value T) *BufIdentifier {
	var id Identifier
	switch v := any(value).(type) {
	case string:
		id = Identifier(v)
	case Identifier:
		id = v
	}
	return &BufIdentifier{&id}
}

type BufIdentifier struct {
	i *Identifier
}

func (b *BufIdentifier) Get() Identifier {
	if b != nil && b.i != nil {
		return *b.i
	}
	return Identifier("")
}

func (b *BufIdentifier) Format(ctx *FmtCtx) {
	if b != nil && b.i != nil {
		ctx.WriteString(string(*b.i))
	}
}
