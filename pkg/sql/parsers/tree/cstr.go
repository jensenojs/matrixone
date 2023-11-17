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
	"strings"

	"github.com/matrixorigin/matrixone/pkg/common/buffer"
)

type CStr struct {
	o *BufString
	c *BufString
	// quote bool

	buf *buffer.Buffer
}

func NewCStr(str string, lower int64, buf *buffer.Buffer) *CStr {
	var cs *CStr
	o := NewBufString(str)
	if buf != nil {
		cs = buffer.Alloc[CStr](buf)
		cs.buf = buf
		buf.Pin(o)
	} else {
		cs = new(CStr)
	}
	cs.o = o
	if lower == 0 {
		cs.c = cs.o
		return cs
	}
	c := NewBufString(strings.ToLower(cs.o.Get()))
	if buf != nil {
		buf.Pin(c)
	}
	cs.c = c
	return cs
}

func (cs *CStr) SetConfig(lower int64) {
	if lower == 0 {
		cs.c = cs.o
		return
	}
	c := NewBufString(strings.ToLower(cs.o.Get()))
	if cs.buf != nil {
		cs.buf.Pin(c)
	}
	cs.c = c
}

func (cs *CStr) ToLower() string {
	return strings.ToLower(cs.o.Get())
}

func (cs *CStr) Origin() string {
	return cs.o.Get()
}

func (cs *CStr) Compare() string {
	return cs.c.Get()
}

func (cs *CStr) Empty() bool {
	return len(cs.o.Get()) == 0
}
