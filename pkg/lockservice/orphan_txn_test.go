// Copyright 2024 Matrix Origin
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

package lockservice

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/matrixorigin/matrixone/pkg/common/morpc"
	pb "github.com/matrixorigin/matrixone/pkg/pb/lock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnlockOrphanTxn(t *testing.T) {
	runLockServiceTestsWithAdjustConfig(
		t,
		[]string{"s1", "s2"},
		time.Second*10,
		func(alloc *lockTableAllocator, s []*service) {
			l1 := s[0]
			l2 := s[1]
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			txn1 := []byte{1}
			txn2 := []byte{2}
			table1 := uint64(1)

			// table1 on l1
			mustAddTestLock(t, ctx, l1, table1, txn1, [][]byte{{1}}, pb.Granularity_Row)

			// txn2 lock row 2 on remote.
			mustAddTestLock(t, ctx, l2, table1, txn2, [][]byte{{2}}, pb.Granularity_Row)
			// l2 shutdown
			assert.NoError(t, l2.Close())

			// wait until txn2 unlocked
			for {
				txn := l1.activeTxnHolder.getActiveTxn(txn2, false, "")
				if txn == nil {
					return
				}
				time.Sleep(time.Millisecond * 100)
			}
		},
		func(c *Config) {
			c.RemoteLockTimeout.Duration = time.Second
			c.KeepRemoteLockDuration.Duration = time.Millisecond * 100
		},
	)
}

func TestCannotUnlockOrphanTxnWithCommunicationInterruption(t *testing.T) {
	var pause atomic.Bool
	remoteLockTimeout := time.Second
	runLockServiceTestsWithAdjustConfig(
		t,
		[]string{"s1", "s2"},
		time.Second*10,
		func(alloc *lockTableAllocator, s []*service) {
			l1 := s[0]
			l2 := s[1]
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			txn1 := []byte{1}
			txn2 := []byte{2}
			table1 := uint64(1)

			// table1 on l1
			mustAddTestLock(t, ctx, l1, table1, txn1, [][]byte{{1}}, pb.Granularity_Row)

			// txn2 lock row 2 on remote.
			mustAddTestLock(t, ctx, l2, table1, txn2, [][]byte{{2}}, pb.Granularity_Row)

			// skip all keep remote lock request and valid service request
			pause.Store(true)

			time.Sleep(remoteLockTimeout * 2)
			txn := l1.activeTxnHolder.getActiveTxn(txn2, false, "")
			assert.NotNil(t, txn)
		},
		func(c *Config) {
			c.RemoteLockTimeout.Duration = remoteLockTimeout
			c.KeepRemoteLockDuration.Duration = time.Millisecond * 100
			c.RPC.BackendOptions = append(c.RPC.BackendOptions,
				morpc.WithBackendFilter(
					func(req morpc.Message, _ string) bool {
						if m, ok := req.(*pb.Request); ok {
							skip := pause.Load() &&
								(m.Method == pb.Method_KeepRemoteLock ||
									m.Method == pb.Method_ValidateService)
							return !skip
						}
						return true
					}))
		},
	)
}

func TestUnlockOrphanTxnWithServiceRestart(t *testing.T) {
	runLockServiceTestsWithAdjustConfig(
		t,
		[]string{"s1", "s2"},
		time.Second*10,
		func(alloc *lockTableAllocator, s []*service) {
			l1 := s[0]
			l2 := s[1]
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			txn1 := []byte{1}
			txn2 := []byte{2}
			table1 := uint64(1)

			// table1 on l1
			mustAddTestLock(t, ctx, l1, table1, txn1, [][]byte{{1}}, pb.Granularity_Row)

			// txn2 lock row 2 on remote.
			mustAddTestLock(t, ctx, l2, table1, txn2, [][]byte{{2}}, pb.Granularity_Row)

			// l2 restart
			assert.NoError(t, l2.Close())
			cfg := l2.cfg
			l22 := NewLockService(cfg)
			defer l22.Close()

			// wait until txn2 unlocked
			for {
				txn := l1.activeTxnHolder.getActiveTxn(txn2, false, "")
				if txn == nil {
					return
				}
				time.Sleep(time.Millisecond * 100)
			}
		},
		func(c *Config) {
			c.RemoteLockTimeout.Duration = time.Second
			c.KeepRemoteLockDuration.Duration = time.Millisecond * 100
		},
	)
}

func TestGetTimeoutRemoveTxn(t *testing.T) {
	hold := newMapBasedTxnHandler(
		"s1",
		newFixedSlicePool(16),
		func(sid string) (bool, error) { return false, nil },
		func(ot []pb.OrphanTxn) error { return nil },
	).(*mapBasedTxnHolder)

	txnID1 := []byte("txn1")
	hold.getActiveTxn(txnID1, true, "s1")
	txnID2 := []byte("txn2")
	hold.getActiveTxn(txnID2, true, "s2")

	// s1(now-10s), s2(now-5s)
	now := time.Now()
	hold.mu.remoteServices["s1"].Value.time = now.Add(-time.Second * 10)
	hold.mu.remoteServices["s2"].Value.time = now.Add(-time.Second * 5)

	txns := hold.getTimeoutRemoveTxn(make(map[string]struct{}), nil, time.Second*20)
	assert.Equal(t, 0, len(txns))

	// s1 timeout
	txns = hold.getTimeoutRemoveTxn(make(map[string]struct{}), nil, time.Second*8)
	assert.Equal(t, 1, len(txns))
	hold.mu.RLock()
	assert.Equal(t, 1, hold.mu.dequeue.Len())
	assert.Equal(t, 1, len(hold.mu.remoteServices))
	hold.mu.RUnlock()

	// s2 timeout
	txns = hold.getTimeoutRemoveTxn(make(map[string]struct{}), nil, time.Second*2)
	assert.Equal(t, 1, len(txns))
	hold.mu.RLock()
	assert.Equal(t, 0, hold.mu.dequeue.Len())
	assert.Equal(t, 0, len(hold.mu.remoteServices))
	hold.mu.RUnlock()
}

func TestGetTimeoutRemoveTxnWithValid(t *testing.T) {
	hold := newMapBasedTxnHandler(
		"s1",
		newFixedSlicePool(16),
		func(sid string) (bool, error) { return sid == "s1", nil },
		func(ot []pb.OrphanTxn) error { return nil },
	).(*mapBasedTxnHolder)

	txnID1 := []byte("txn1")
	hold.getActiveTxn(txnID1, true, "s1")
	txnID2 := []byte("txn2")
	hold.getActiveTxn(txnID2, true, "s2")

	// s1(now-10s), s2(now-5s)
	now := time.Now()
	hold.mu.remoteServices["s1"].Value.time = now.Add(-time.Second * 10)
	hold.mu.remoteServices["s2"].Value.time = now.Add(-time.Second * 5)

	txns := hold.getTimeoutRemoveTxn(make(map[string]struct{}), nil, time.Second*20)
	assert.Equal(t, 0, len(txns))
	hold.mu.RLock()
	assert.Equal(t, 2, hold.mu.dequeue.Len())
	assert.Equal(t, 2, len(hold.mu.remoteServices))
	hold.mu.RUnlock()

	// s1 timeout
	txns = hold.getTimeoutRemoveTxn(make(map[string]struct{}), nil, time.Second*8)
	assert.Equal(t, 0, len(txns))
	hold.mu.RLock()
	assert.Equal(t, 2, hold.mu.dequeue.Len())
	assert.Equal(t, 2, len(hold.mu.remoteServices))
	hold.mu.RUnlock()

	// s2 timeout
	txns = hold.getTimeoutRemoveTxn(make(map[string]struct{}), nil, time.Second*2)
	assert.Equal(t, 1, len(txns))
	hold.mu.RLock()
	assert.Equal(t, 1, hold.mu.dequeue.Len())
	assert.Equal(t, 1, len(hold.mu.remoteServices))
	hold.mu.RUnlock()
}

func TestGetTimeoutRemoveTxnWithValidErrorAndNotifyOK(t *testing.T) {
	hold := newMapBasedTxnHandler(
		"s1",
		newFixedSlicePool(16),
		func(sid string) (bool, error) { return false, ErrTxnNotFound },
		func(ot []pb.OrphanTxn) error { return nil },
	).(*mapBasedTxnHolder)

	txnID1 := []byte("txn1")
	hold.getActiveTxn(txnID1, true, "s1")

	// s1(now-10s)
	now := time.Now()
	hold.mu.remoteServices["s1"].Value.time = now.Add(-time.Second * 10)

	// s1 timeout
	txns := hold.getTimeoutRemoveTxn(make(map[string]struct{}), nil, time.Second*8)
	assert.Equal(t, 1, len(txns))
	hold.mu.RLock()
	assert.Equal(t, 0, hold.mu.dequeue.Len())
	assert.Equal(t, 0, len(hold.mu.remoteServices))
	hold.mu.RUnlock()
}

func TestGetTimeoutRemoveTxnWithValidErrorAndNotifyFailed(t *testing.T) {
	hold := newMapBasedTxnHandler(
		"s1",
		newFixedSlicePool(16),
		func(sid string) (bool, error) { return false, ErrTxnNotFound },
		func(ot []pb.OrphanTxn) error { return ErrTxnNotFound },
	).(*mapBasedTxnHolder)

	txnID1 := []byte("txn1")
	hold.getActiveTxn(txnID1, true, "s1")

	// s1(now-10s)
	now := time.Now()
	hold.mu.remoteServices["s1"].Value.time = now.Add(-time.Second * 10)

	// s1 timeout
	txns := hold.getTimeoutRemoveTxn(make(map[string]struct{}), nil, time.Second*8)
	assert.Equal(t, 0, len(txns))
	hold.mu.RLock()
	assert.Equal(t, 1, hold.mu.dequeue.Len())
	assert.Equal(t, 1, len(hold.mu.remoteServices))
	hold.mu.RUnlock()
}

func TestCannotCommitTxnCanBeRemovedWithNotInActiveTxn(t *testing.T) {
	var mu sync.Mutex
	var actives [][]byte

	fn := func(f func(txn []byte) bool) {
		mu.Lock()
		defer mu.Unlock()
		for _, txn := range actives {
			if !f(txn) {
				return
			}
		}
	}

	runLockServiceTestsWithAdjustConfig(
		t,
		[]string{"s1", "s2"},
		time.Millisecond*20,
		func(alloc *lockTableAllocator, s []*service) {
			l1 := s[0]
			mu.Lock()
			actives = append(actives, []byte{1}, []byte{2})
			mu.Unlock()

			alloc.AddCannotCommit([]pb.OrphanTxn{
				{
					Service: l1.serviceID,
					Txn:     [][]byte{{1}, {3}},
				},
			})

			// wait txn 3 removed
			for {
				alloc.mu.RLock()
				v, ok := alloc.mu.cannotCommit[getUUIDFromServiceIdentifier(l1.serviceID)]
				require.True(t, ok)

				_, ok = v.txn[string([]byte{1})]
				if len(v.txn) == 1 && ok {
					alloc.mu.RUnlock()
					return
				}
				alloc.mu.RUnlock()
				time.Sleep(time.Millisecond * 10)
			}
		},
		func(c *Config) {
			c.TxnIterFunc = fn
		},
	)
}

func TestCannotCommitTxnCanBeRemovedWithRestart(t *testing.T) {
	var mu sync.Mutex
	var actives [][]byte

	fn := func(f func(txn []byte) bool) {
		mu.Lock()
		defer mu.Unlock()
		for _, txn := range actives {
			if !f(txn) {
				return
			}
		}
	}

	runLockServiceTestsWithAdjustConfig(
		t,
		[]string{"s1", "s2"},
		time.Millisecond*20,
		func(alloc *lockTableAllocator, s []*service) {
			l1 := s[0]
			mu.Lock()
			actives = append(actives, []byte{1}, []byte{3})
			mu.Unlock()

			alloc.AddCannotCommit([]pb.OrphanTxn{
				{
					Service: l1.serviceID,
					Txn:     [][]byte{{1}, {3}},
				},
			})

			cfg := l1.GetConfig()
			require.NoError(t, l1.Close())
			l1 = NewLockService(cfg).(*service)
			defer func() {
				l1.Close()
			}()

			// wait txn 3 removed
			for {
				alloc.mu.RLock()
				n := len(alloc.mu.cannotCommit)
				if n == 0 {
					alloc.mu.RUnlock()
					return
				}
				alloc.mu.RUnlock()
				time.Sleep(time.Millisecond * 10)
			}
		},
		func(c *Config) {
			c.TxnIterFunc = fn
		},
	)
}