package tempengine

import (
	"context"

	"github.com/matrixorigin/matrixone/pkg/container/batch"
	"github.com/matrixorigin/matrixone/pkg/container/vector"
	"github.com/matrixorigin/matrixone/pkg/pb/plan"
	"github.com/matrixorigin/matrixone/pkg/vm/engine"
)

func (tempRelation *TempRelation) Rows() int64 {
	return int64(tempRelation.rows)
}

func (tempRelation *TempRelation) Size(_ string) int64 {
	return 0
}

func (tempRelation *TempRelation) Ranges(context.Context) ([][]byte, error) {
	return nil, nil
}

func (tempRelation *TempRelation) TableDefs(context.Context) ([]engine.TableDef, error) {
	return nil, nil
}

func (tempRelation *TempRelation) GetPrimaryKeys(context.Context) ([]*engine.Attribute, error) {
	return nil, nil
}

func (tempRelation *TempRelation) GetHideKeys(context.Context) ([]*engine.Attribute, error) {
	return nil, nil
}

func (tempRelation *TempRelation) Write(context.Context, *batch.Batch) error {
	return nil
}

func (tempRelation *TempRelation) Update(context.Context, *batch.Batch) error {
	return nil
}

func (tempRelation *TempRelation) Delete(context.Context, *vector.Vector, string) error {
	return nil
}

func (tempRelation *TempRelation) Truncate(context.Context) (uint64, error) {
	return 0, nil
}

func (tempRelation *TempRelation) AddTableDef(context.Context, engine.TableDef) error {
	return nil
}

func (tempRelation *TempRelation) DelTableDef(context.Context, engine.TableDef) error {
	return nil
}

func (tempRelation *TempRelation) GetTableID(context.Context) string {
	return ""
}

// second argument is the number of reader, third argument is the filter extend, foruth parameter is the payload required by the engine
func (tempRelation *TempRelation) NewReader(context.Context, int, *plan.Expr, [][]byte) ([]engine.Reader, error) {
	return nil, nil
}
