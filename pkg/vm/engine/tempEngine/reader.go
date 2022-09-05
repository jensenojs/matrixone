package tempengine

import (
	"github.com/matrixorigin/matrixone/pkg/container/batch"
	"github.com/matrixorigin/matrixone/pkg/pb/plan"
	"github.com/matrixorigin/matrixone/pkg/vm/mheap"
)

func (reader *TempReader) Close() error {
	return nil
}
func (reader *TempReader) Read([]string, *plan.Expr, *mheap.Mheap) (*batch.Batch, error) {
	return nil, nil
}
