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

package client

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/lni/goutils/leaktest"
	"github.com/matrixorigin/matrixone/pkg/common/runtime"
	"github.com/matrixorigin/matrixone/pkg/logutil"
	"github.com/matrixorigin/matrixone/pkg/pb/metadata"
	"github.com/matrixorigin/matrixone/pkg/pb/timestamp"
	"github.com/matrixorigin/matrixone/pkg/pb/txn"
	"github.com/matrixorigin/matrixone/pkg/txn/clock"
	"github.com/matrixorigin/matrixone/pkg/txn/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdjustClient(t *testing.T) {
	runtime.SetupServiceBasedRuntime("", runtime.DefaultRuntime())
	c := &txnClient{}
	c.adjust()
	assert.NotNil(t, c.generator)
	assert.NotNil(t, c.generator)
}

func TestNewTxnAndReset(t *testing.T) {
	rt := runtime.NewRuntime(metadata.ServiceType_CN, "",
		logutil.GetPanicLogger(),
		runtime.WithClock(clock.NewHLCClock(func() int64 {
			return 1
		}, 0)))
	runtime.SetupServiceBasedRuntime("", rt)
	c := NewTxnClient("", newTestTxnSender())
	c.Resume()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	tx, err := c.New(ctx, newTestTimestamp(0))
	assert.Nil(t, err)
	txnMeta := tx.(*txnOperator).mu.txn
	assert.Equal(t, timestamp.Timestamp{PhysicalTime: 0}, txnMeta.SnapshotTS)
	assert.NotEmpty(t, txnMeta.ID)
	assert.Equal(t, txn.TxnStatus_Active, txnMeta.Status)

	require.NoError(t, tx.Rollback(ctx))

	tx, err = c.RestartTxn(ctx, tx, newTestTimestamp(0))
	assert.Nil(t, err)
	txnMeta = tx.(*txnOperator).mu.txn
	assert.Equal(t, timestamp.Timestamp{PhysicalTime: 0}, txnMeta.SnapshotTS)
	assert.NotEmpty(t, txnMeta.ID)
	assert.Equal(t, txn.TxnStatus_Active, txnMeta.Status)
}

func TestNewTxnWithNormalStateWait(t *testing.T) {
	defer leaktest.AfterTest(t)()
	rt := runtime.NewRuntime(metadata.ServiceType_CN, "",
		logutil.GetPanicLogger(),
		runtime.WithClock(clock.NewHLCClock(func() int64 {
			return 1
		}, 0)))
	runtime.SetupServiceBasedRuntime("", rt)
	c := NewTxnClient("", newTestTxnSender())
	defer func() {
		require.NoError(t, c.Close())
	}()

	// Do not resume the txn client for now.
	// c.Resume()
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()
			tx, err := c.New(ctx, newTestTimestamp(0))
			assert.Nil(t, err)
			txnMeta := tx.(*txnOperator).mu.txn
			assert.Equal(t, int64(0), txnMeta.SnapshotTS.PhysicalTime)
			assert.NotEmpty(t, txnMeta.ID)
			assert.Equal(t, txn.TxnStatus_Active, txnMeta.Status)
		}()
	}
	// Resume it now.
	c.Resume()
	wg.Wait()
}

func TestNewTxnWithNormalStateNoWait(t *testing.T) {
	defer leaktest.AfterTest(t)()
	rt := runtime.NewRuntime(metadata.ServiceType_CN, "",
		logutil.GetPanicLogger(),
		runtime.WithClock(clock.NewHLCClock(func() int64 {
			return 1
		}, 0)))
	runtime.SetupServiceBasedRuntime("", rt)
	c := NewTxnClient("", newTestTxnSender(), WithNormalStateNoWait(true))
	defer func() {
		require.NoError(t, c.Close())
	}()

	// Do not resume the txn client.
	// c.Resume()
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()
			tx, err := c.New(ctx, newTestTimestamp(0))
			assert.Error(t, err)
			assert.Nil(t, tx)
		}()
	}
	wg.Wait()
}

func TestNewTxnWithSnapshotTS(t *testing.T) {
	rt := runtime.NewRuntime(metadata.ServiceType_CN, "",
		logutil.GetPanicLogger(),
		runtime.WithClock(clock.NewHLCClock(func() int64 {
			return 1
		}, 0)))
	runtime.SetupServiceBasedRuntime("", rt)
	c := NewTxnClient("", newTestTxnSender())
	c.Resume()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	tx, err := c.New(ctx, newTestTimestamp(0), WithSnapshotTS(timestamp.Timestamp{PhysicalTime: 10}))
	assert.Nil(t, err)
	txnMeta := tx.(*txnOperator).mu.txn
	assert.Equal(t, timestamp.Timestamp{PhysicalTime: 10}, txnMeta.SnapshotTS)
	assert.NotEmpty(t, txnMeta.ID)
	assert.Equal(t, txn.TxnStatus_Active, txnMeta.Status)
}

type fakeRunningPipelinesManager struct{}

func (m *fakeRunningPipelinesManager) PauseService()            {}
func (m *fakeRunningPipelinesManager) KillAllQueriesWithError() {}
func (m *fakeRunningPipelinesManager) ResumeService()           {}

func TestTxnClientAbortAllRunningTxn(t *testing.T) {
	SetRunningPipelineManagement(&fakeRunningPipelinesManager{})
	rt := runtime.NewRuntime(metadata.ServiceType_CN, "",
		logutil.GetPanicLogger(),
		runtime.WithClock(clock.NewHLCClock(func() int64 {
			return 1
		}, 0)))
	runtime.SetupServiceBasedRuntime("", rt)

	c := NewTxnClient("", newTestTxnSender())
	c.Resume()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	var ops []TxnOperator
	for i := 0; i < 10; i++ {
		op, err := c.New(ctx, newTestTimestamp(0))
		assert.Nil(t, err)
		ops = append(ops, op)
	}
	require.Equal(t, 10, len(c.(*txnClient).mu.activeTxns))

	c.AbortAllRunningTxn()
	require.Equal(t, 0, len(c.(*txnClient).mu.activeTxns))
	for _, op := range ops {
		assert.Equal(t, txn.TxnStatus_Aborted, op.(*txnOperator).mu.txn.Status)
	}
}

func TestTxnClientPauseAndResume(t *testing.T) {
	rt := runtime.NewRuntime(metadata.ServiceType_CN, "",
		logutil.GetPanicLogger(),
		runtime.WithClock(clock.NewHLCClock(func() int64 {
			return 1
		}, 0)))
	runtime.SetupServiceBasedRuntime("", rt)
	c := NewTxnClient("", newTestTxnSender())

	c.Pause()
	require.Equal(t, paused, c.(*txnClient).mu.state)
	c.Resume()
	require.Equal(t, normal, c.(*txnClient).mu.state)
}

func TestLimit(t *testing.T) {
	RunTxnTests(
		func(tc TxnClient, ts rpc.TxnSender) {
			ctx := context.Background()

			c := make(chan struct{})
			c2 := make(chan struct{})
			n := 0
			go func() {
				defer close(c2)
				for {
					select {
					case <-c:
						return
					default:
						op, err := tc.New(ctx, newTestTimestamp(0))
						require.NoError(t, err)
						require.NoError(t, op.Rollback(ctx))
						n++
					}
				}
			}()
			time.Sleep(time.Millisecond * 200)
			close(c)
			<-c2
			require.True(t, n < 5)
		},
		WithTxnLimit(1))
}

func TestMaxActiveTxnWithWaitPrevClosed(t *testing.T) {
	RunTxnTests(
		func(tc TxnClient, ts rpc.TxnSender) {
			ctx := context.Background()
			op1, err := tc.New(ctx, newTestTimestamp(0), WithUserTxn())
			require.NoError(t, err)

			c := make(chan struct{})
			go func() {
				defer close(c)
				_, err = tc.New(ctx, newTestTimestamp(0), WithUserTxn())
				require.NoError(t, err)
			}()

			require.NoError(t, op1.Rollback(ctx))
			<-c
		},
		WithMaxActiveTxn(1))
}

func TestMaxActiveTxnWithWaitTimeout(t *testing.T) {
	RunTxnTests(
		func(tc TxnClient, ts rpc.TxnSender) {
			ctx := context.Background()
			op1, err := tc.New(ctx, newTestTimestamp(0), WithUserTxn())
			require.NoError(t, err)
			defer func() {
				require.NoError(t, op1.Rollback(ctx))
			}()

			ctx2, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err = tc.New(ctx2, newTestTimestamp(0), WithUserTxn())
			require.Error(t, err)
		},
		WithMaxActiveTxn(1))
}

func TestOpenTxnWithWaitPausedDisabled(t *testing.T) {
	c := &txnClient{}
	c.mu.state = paused

	op := &txnOperator{}
	op.opts.options = op.opts.options.WithDisableWaitPaused()

	require.Error(t, c.openTxn(op))
}
