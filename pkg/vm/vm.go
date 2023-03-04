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

package vm

import (
	"bytes"
	"strings"

	"github.com/matrixorigin/matrixone/pkg/common/moerr"
	"github.com/matrixorigin/matrixone/pkg/logutil"
	"github.com/matrixorigin/matrixone/pkg/sql/colexec/agg"
	"github.com/matrixorigin/matrixone/pkg/vm/process"
)

func (ins Instruction) MarshalBinary() ([]byte, error) {
	return nil, nil
}

func (ins Instruction) UnmarshalBinary(_ []byte) error {
	return nil
}

// String range instructions and call each operator's string function to show a query plan
func String(ins Instructions, buf *bytes.Buffer) {
	for i, in := range ins {
		if i > 0 {
			buf.WriteString(" -> ")
		}
		stringFunc[in.Op](in.Arg, buf)
	}
}

// Prepare range instructions and do init work for each operator's argument by calling its prepare function
func Prepare(ins Instructions, proc *process.Process) error {
	for _, in := range ins {
		if err := prepareFunc[in.Op](proc, in.Arg); err != nil {
			return err
		}
	}
	return nil
}

// more check here
func Run(ins Instructions, proc *process.Process, sql string) (end bool, err error) {
	var ok bool

	defer func() {
		if e := recover(); e != nil {
			err = moerr.ConvertPanicError(proc.Ctx, e)
		}
	}()
	// 	Scope 1 (Magic: Merge, Receiver: [4]): [merge -> output]
	//   PreScopes: {
	//   Scope 1 (Magic: Merge, Receiver: [3]): [merge group -> projection -> projection -> connect to MergeReceiver 4]
	//     PreScopes: {
	//     Scope 1 (Magic: Merge, Receiver: [2]): [merge -> group -> connect to MergeReceiver 3]
	//       PreScopes: {
	//       Scope 1 (Magic: Merge, Receiver: [0, 1]): [merge group -> projection -> restrict -> projection -> projection -> connect to MergeReceiver 2]
	//         PreScopes: {
	//         Scope 1 (Magic: Remote, Receiver: []): [projection -> group -> connect to MergeReceiver 0]
	//         DataSource: t.t1[a b],
	//         Scope 2 (Magic: Remote, Receiver: []): [projection -> group -> connect to MergeReceiver 1]
	//         DataSource: t.t1[a b],
	//       }
	//     }
	//   }
	// }
	for _, in := range ins {
		if strings.Contains(sql, "SELECT DISTINCT a, AVG( b) FROM test_7748 GROUP BY a HAVING AVG( b) > 50") {
			// logutil.Errorf("+++: Target sql statement showed")
			if len(ins) == 3 && ins[0].Op == Projection && ins[1].Op == Group && in.Op == Connector && proc.Reg.InputBatch != nil && proc.Reg.InputBatch.Aggs != nil {
				a := proc.Reg.InputBatch.Aggs[0].(*agg.UnaryAgg[int32, float64])
				logutil.Errorf("+++: After Group : the count(b) is %v, and the sum(b) is %f", a.Priv, a.Vs[0])
				// after group operator, let check avg result
				if ok, err = execFunc[in.Op](in.Idx, proc, in.Arg, in.IsFirst, in.IsLast); err != nil {
					return ok || end, err
				}
				if ok { // ok is true shows that at least one operator has done its work
					end = true
				}
			} else if len(ins) == 6 && ins[0].Op == MergeGroup && ins[2].Op == Restrict && in.Op == Projection {
				bVec := proc.Reg.InputBatch.GetVector(1)
				logutil.Errorf("+++: After MergeGroup : the values is %f", bVec.Col.([]float64)[0])
				if ok, err = execFunc[in.Op](in.Idx, proc, in.Arg, in.IsFirst, in.IsLast); err != nil {
					return ok || end, err
				}
				if ok { // ok is true shows that at least one operator has done its work
					end = true
				}
			} else {
				// logutil.Errorf("+++: not catch")
				if ok, err = execFunc[in.Op](in.Idx, proc, in.Arg, in.IsFirst, in.IsLast); err != nil {
					return ok || end, err
				}
				if ok { // ok is true shows that at least one operator has done its work
					end = true
				}
			}
		} else {
			if ok, err = execFunc[in.Op](in.Idx, proc, in.Arg, in.IsFirst, in.IsLast); err != nil {
				return ok || end, err
			}
			if ok { // ok is true shows that at least one operator has done its work
				end = true
			}
		}
	}
	return end, err
}
