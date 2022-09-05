package tempengine

import (
	"context"
	"time"

	"github.com/matrixorigin/matrixone/pkg/txn/client"
	"github.com/matrixorigin/matrixone/pkg/vm/engine"
)

func NewTempEngine() *TempEngine {
	return &TempEngine{
		tempDatabase: &TempDatabase{
			nameToRelation: make(map[string]*TempRelation),
		},
	}
}

// Delete deletes a database
func (tempEngine *TempEngine) Delete(ctx context.Context, databaseName string, op client.TxnOperator) error {
	return nil
}

// Create creates a database
func (tempEngine *TempEngine) Create(ctx context.Context, databaseName string, op client.TxnOperator) error {
	return nil
}

// Databases returns all database names
func (tempEngine *TempEngine) Databases(ctx context.Context, op client.TxnOperator) (databaseNames []string, err error) {
	return nil, nil
}

// Database creates a handle for a database
func (tempEngine *TempEngine) Database(ctx context.Context, databaseName string, op client.TxnOperator) (engine.Database, error) {
	return tempEngine.tempDatabase, nil
}

// Nodes returns all nodes for worker jobs
func (tempEngine *TempEngine) Nodes() (cnNodes engine.Nodes, err error) {
	return nil, nil
}

// Hints returns hints of engine features
// return value should not be cached
// since implementations may update hints after engine had initialized
func (tempEngine *TempEngine) Hints() (h engine.Hints) {
	h.CommitOrRollbackTimeout = time.Minute
	return
}
