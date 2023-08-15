// Copyright 2023 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package checkpkdup

import (
	"github.com/matrixorigin/matrixone/pkg/common/bitmap"
	"github.com/matrixorigin/matrixone/pkg/vm/process"
)

type Argument struct {
	potentially_dup_keys map[any]bool
	bitmp                bitmap.Bitmap
}

func (arg *Argument) Free(proc *process.Process, pipelineFailed bool) {
	// ctr := arg.ctr
	// bitmap currently have no Free method, maybe need to add it
}
