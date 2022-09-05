package tempengine

import (
	"context"
	"fmt"

	"github.com/matrixorigin/matrixone/pkg/vm/engine"
)

func (db *TempDatabase) Relations(context.Context) ([]string, error) {
	return nil, nil
}

func (db *TempDatabase) Relation(ctx context.Context, db_tblNme string) (engine.Relation, error) {
	// now default we can't find
	if val := db.nameToRelation[db_tblNme]; val != nil {
		return val, nil
	}
	return nil, fmt.Errorf("not found table %s", getTblName(db_tblNme))
}

func (db *TempDatabase) Delete(_ context.Context, db_tblNme string) error {
	delete(db.nameToRelation, db_tblNme)
	return nil
}

// remember that the db_tblname is really "databaseName-tblName"
func (db *TempDatabase) Create(ctx context.Context, db_tblName string, tblDefs []engine.TableDef) error {
	// Create Table - (name, table define)
	schema, err := DefsToSchema(db_tblName, tblDefs)
	if err != nil {
		return err
	}
	_, err = db.Relation(ctx, db_tblName)
	if err == nil {
		return fmt.Errorf("table '%s' already exists", getTblName(db_tblName))
	}
	db.nameToRelation[db_tblName] = &TempRelation{
		tblSchema: *schema,
		blockNums: 0,
		bytesData: make(map[string][]byte),
	}
	return nil
}

func DefsToSchema(db_tblName string, tblDefs []engine.TableDef) (*TableSchema, error) {
	tblSchema := &TableSchema{
		tblName: db_tblName,
	}
	for _, def := range tblDefs {
		switch v := def.(type) {
		case *engine.AttributeDef:
			tblSchema.attrs = append(tblSchema.attrs, v.Attr)
		default:
			return nil, fmt.Errorf("only support normal attribute now")
		}
	}
	return tblSchema, nil
}
