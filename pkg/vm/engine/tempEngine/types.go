package tempengine

import "github.com/matrixorigin/matrixone/pkg/vm/engine"

// note that,we don't need to think about Concurrency,
// as we know, in one Session, we hava one TempEngine,
// different sessions won't bother each other,
// and we can't execute multi-sqls one time

// TempEngine is used to store temporary tables
// For TempEngine, it just has only one database, itdoesn't need more
// when we find a tempRelation, just use "databaseName-TableName"
type TempEngine struct {
	// node engine.Node
	tempDatabase *TempDatabase
}

// once TempRelation starts to write a batch
// we will regard this batch as a block, and then we will
type TempRelation struct {
	tblSchema TableSchema
	blockNums uint64
	rows      uint64
	// string will be like tblName-blockId-attrName, the []byte is the data of the col in
	// this block (a block is a batch)
	bytesData map[string][]byte
}

type TempDatabase struct {
	nameToRelation map[string]*TempRelation // "string shoule be databaseName-tblName"
}

type TempReader struct {
	blockIdxs   []uint64 // indicate which blocks this reader should read
	ownRelation *TempRelation
}

type TableSchema struct {
	tblName string
	attrs   []engine.Attribute // record the
	// TODO: we need to support primaryKey, partition key and auto_increment_key so on
	// this needs to refer to the TAE engine
}
