// Copyright 2021 - 2023 Matrix Origin
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

package bloomfilter

import (
	"github.com/matrixorigin/matrixone/pkg/container/vector"
)

func (bf *BloomFilter) Test() bool {

	return false
}

func (bf *BloomFilter) Add() bool {
	return false
}

func (bf *BloomFilter) TestAndAdd() bool {
	return false
}

func (bf *BloomFilter) TestAndAddForVector(v *vector.Vector) bool {

	return false
}

// func getHashValue() {
// 	keys := make([][]byte, hashmap.UnitLimit)
// 	states := make([][3]uint64, hashmap.UnitLimit)
// 	for i := 0; i < length; i += hashmap.UnitLimit {
// 		n := length - i
// 		if n > hashmap.UnitLimit {
// 			n = hashmap.UnitLimit
// 		}
// 		for j := 0; j < n; j++ {
// 			keys[j] = keys[j][:0]
// 		}
// 		encodeHashKeys(keys, parameters, i, n)

// 		hashtable.BytesBatchGenHashStates(&keys[0], &states[0], n)
// 		for j := 0; j < n; j++ {
// 			rs.AppendMustValue(int64(states[j][0]))
// 		}
// 	}
// }
