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

package tree

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/matrixorigin/matrixone/pkg/common/buffer"
)

type CreateOption interface {
	NodeFormatter
}

type createOptionImpl struct {
	CreateOption
}

type CreateOptionDefault struct {
	createOptionImpl
}

type CreateOptionCharset struct {
	createOptionImpl
	IsDefault bool
	Charset   *BufString
}

func (node *CreateOptionCharset) Format(ctx *FmtCtx) {
	if node.IsDefault {
		ctx.WriteString("default ")
	}
	ctx.WriteString("character set ")
	ctx.WriteString(node.Charset.Get())
}

func NewCreateOptionCharset(isf bool, charset string, buf *buffer.Buffer) *CreateOptionCharset {
	c := buffer.Alloc[CreateOptionCharset](buf)
	c.IsDefault = isf
	bc := NewBufString(charset)
	buf.Pin(bc)
	c.Charset = bc
	return c
}

type CreateOptionCollate struct {
	createOptionImpl
	IsDefault bool
	Collate   *BufString
}

func (node *CreateOptionCollate) Format(ctx *FmtCtx) {
	if node.IsDefault {
		ctx.WriteString("default ")
	}
	ctx.WriteString("collate ")
	ctx.WriteString(node.Collate.Get())
}

func NewCreateOptionCollate(isd bool, collate string, buf *buffer.Buffer) *CreateOptionCollate {
	c := buffer.Alloc[CreateOptionCollate](buf)
	c.IsDefault = isd
	bc := NewBufString(collate)
	buf.Pin(bc)
	c.Collate = bc
	return c
}

type CreateOptionEncryption struct {
	createOptionImpl
	Encrypt *BufString
}

func (node *CreateOptionEncryption) Format(ctx *FmtCtx) {
	ctx.WriteString("encryption ")
	ctx.WriteString(node.Encrypt.Get())
}

func NewCreateOptionEncryption(encrypt string, buf *buffer.Buffer) *CreateOptionEncryption {
	c := buffer.Alloc[CreateOptionEncryption](buf)
	ec := NewBufString(encrypt)
	buf.Pin(ec)
	c.Encrypt = ec
	return c
}

type SubscriptionOption struct {
	statementImpl
	From        *BufIdentifier
	Publication *BufIdentifier
}

func (node *SubscriptionOption) Format(ctx *FmtCtx) {
	ctx.WriteString(" from ")
	node.From.Format(ctx)
	ctx.WriteString(" publication ")
	node.Publication.Format(ctx)
}

func NewSubscriptionOption(from, publication Identifier, buf *buffer.Buffer) *SubscriptionOption {
	s := buffer.Alloc[SubscriptionOption](buf)

	f := NewBufIdentifier(from)
	p := NewBufIdentifier(publication)
	buf.Pin(f, p)

	s.From = f
	s.Publication = p
	return s
}

type CreateDatabase struct {
	statementImpl
	IfNotExists        bool
	Name               *BufIdentifier
	CreateOptions      []CreateOption
	SubscriptionOption *SubscriptionOption
}

func (node *CreateDatabase) Format(ctx *FmtCtx) {
	ctx.WriteString("create database ")
	if node.IfNotExists {
		ctx.WriteString("if not exists ")
	}
	node.Name.Format(ctx)

	if node.SubscriptionOption != nil {
		node.SubscriptionOption.Format(ctx)
	}

	if node.CreateOptions != nil {
		for _, opt := range node.CreateOptions {
			ctx.WriteByte(' ')
			opt.Format(ctx)
		}
	}
}

func (node *CreateDatabase) GetStatementType() string { return "Create Database" }
func (node *CreateDatabase) GetQueryType() string     { return QueryTypeDDL }

func NewCreateDatabase(ifNotExists bool, name Identifier, sub *SubscriptionOption, createOptions []CreateOption, buf *buffer.Buffer) *CreateDatabase {
	c := buffer.Alloc[CreateDatabase](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)

	c.Name = n
	c.IfNotExists = ifNotExists
	c.SubscriptionOption = sub
	c.CreateOptions = createOptions
	return c
}

type CreateTable struct {
	statementImpl
	/*
		it is impossible to be the temporary table, the cluster table,
		the normal table and the external table at the same time.
	*/
	Temporary       bool
	IsClusterTable  bool
	IfNotExists     bool
	Table           *TableName
	Defs            TableDefs
	Options         []TableOption
	PartitionOption *PartitionOption
	ClusterByOption *ClusterByOption
	Param           *ExternParam
}

func NewCreateTable(temporary, isClusterTable, ifNotExists bool, table *TableName, defs TableDefs, options []TableOption, partitionOption *PartitionOption, clusterByOption *ClusterByOption, param *ExternParam, buf *buffer.Buffer) *CreateTable {
	c := buffer.Alloc[CreateTable](buf)
	c.Temporary = temporary
	c.IsClusterTable = isClusterTable
	c.IfNotExists = ifNotExists
	c.Table = table
	c.Defs = defs
	c.Options = options
	c.PartitionOption = partitionOption
	c.ClusterByOption = clusterByOption
	c.Param = param
	return c
}

func (node *CreateTable) Format(ctx *FmtCtx) {
	ctx.WriteString("create")
	if node.Temporary {
		ctx.WriteString(" temporary")
	}
	if node.IsClusterTable {
		ctx.WriteString(" cluster")
	}
	if node.Param != nil {
		ctx.WriteString(" external")
	}

	ctx.WriteString(" table")

	if node.IfNotExists {
		ctx.WriteString(" if not exists")
	}

	ctx.WriteByte(' ')
	node.Table.Format(ctx)

	ctx.WriteString(" (")
	for i, def := range node.Defs {
		if i != 0 {
			ctx.WriteString(",")
			ctx.WriteByte(' ')
		}
		def.Format(ctx)
	}
	ctx.WriteByte(')')

	if node.Options != nil {
		prefix := " "
		for _, t := range node.Options {
			ctx.WriteString(prefix)
			t.Format(ctx)
		}
	}

	if node.PartitionOption != nil {
		ctx.WriteByte(' ')
		node.PartitionOption.Format(ctx)
	}

	if node.Param != nil {
		if len(node.Param.Option) == 0 {
			ctx.WriteString(" infile ")
			ctx.WriteString("'" + node.Param.Filepath.Get() + "'")
		} else {
			if node.Param.ScanType == S3 {
				ctx.WriteString(" url s3option ")
			} else {
				ctx.WriteString(" infile ")
			}
			ctx.WriteString("{")
			for i := 0; i < len(node.Param.Option); i += 2 {
				switch strings.ToLower(node.Param.Option[i]) {
				case "endpoint":
					ctx.WriteString("'endpoint'='" + node.Param.Option[i+1] + "'")
				case "region":
					ctx.WriteString("'region'='" + node.Param.Option[i+1] + "'")
				case "access_key_id":
					ctx.WriteString("'access_key_id'='******'")
				case "secret_access_key":
					ctx.WriteString("'secret_access_key'='******'")
				case "bucket":
					ctx.WriteString("'bucket'='" + node.Param.Option[i+1] + "'")
				case "filepath":
					ctx.WriteString("'filepath'='" + node.Param.Option[i+1] + "'")
				case "compression":
					ctx.WriteString("'compression'='" + node.Param.Option[i+1] + "'")
				case "format":
					ctx.WriteString("'format'='" + node.Param.Option[i+1] + "'")
				case "jsondata":
					ctx.WriteString("'jsondata'='" + node.Param.Option[i+1] + "'")
				}
				if i != len(node.Param.Option)-2 {
					ctx.WriteString(", ")
				}
			}
			ctx.WriteString("}")
		}
		if node.Param.Tail.Fields != nil {
			ctx.WriteByte(' ')
			node.Param.Tail.Fields.Format(ctx)
		}

		if node.Param.Tail.Lines != nil {
			ctx.WriteByte(' ')
			node.Param.Tail.Lines.Format(ctx)
		}

		if node.Param.Tail.IgnoredLines != 0 {
			ctx.WriteString(" ignore ")
			ctx.WriteString(strconv.FormatUint(node.Param.Tail.IgnoredLines, 10))
			ctx.WriteString(" lines")
		}
		if node.Param.Tail.ColumnList != nil {
			prefix := " ("
			for _, c := range node.Param.Tail.ColumnList {
				ctx.WriteString(prefix)
				c.Format(ctx)
				prefix = ", "
			}
			ctx.WriteByte(')')
		}
		if node.Param.Tail.Assignments != nil {
			ctx.WriteString(" set ")
			node.Param.Tail.Assignments.Format(ctx)
		}
	}
}

func (node *CreateTable) GetStatementType() string { return "Create Table" }
func (node *CreateTable) GetQueryType() string     { return QueryTypeDDL }

type TableDef interface {
	NodeFormatter
}

type tableDefImpl struct {
	TableDef
}

// the list of table definitions
type TableDefs []TableDef

type ColumnTableDef struct {
	tableDefImpl
	Name       *UnresolvedName
	Type       ResolvableTypeReference
	Attributes []ColumnAttribute
}

func (node *ColumnTableDef) Format(ctx *FmtCtx) {
	node.Name.Format(ctx)

	ctx.WriteByte(' ')
	node.Type.(*T).InternalType.Format(ctx)

	if node.Attributes != nil {
		prefix := " "
		for _, a := range node.Attributes {
			ctx.WriteString(prefix)
			if a != nil {
				a.Format(ctx)
			}
		}
	}
}

func NewColumnTableDef(name *UnresolvedName, typ ResolvableTypeReference, attributes []ColumnAttribute, buf *buffer.Buffer) *ColumnTableDef {
	var c *ColumnTableDef
	if buf != nil {
		c = buffer.Alloc[ColumnTableDef](buf)
	} else {
		c = new(ColumnTableDef)
	}
	c.Name = name
	c.Type = typ
	c.Attributes = attributes
	return c
}

// column attribute
type ColumnAttribute interface {
	NodeFormatter
}

type columnAttributeImpl struct {
	ColumnAttribute
}

type AttributeNull struct {
	columnAttributeImpl
	Is bool //true NULL (default); false NOT NULL
}

func (node *AttributeNull) Format(ctx *FmtCtx) {
	if node.Is {
		ctx.WriteString("null")
	} else {
		ctx.WriteString("not null")
	}
}

func NewAttributeNull(is bool, buf *buffer.Buffer) *AttributeNull {
	a := buffer.Alloc[AttributeNull](buf)
	a.Is = is
	return a
}

type AttributeDefault struct {
	columnAttributeImpl
	Expr Expr
}

func (node *AttributeDefault) Format(ctx *FmtCtx) {
	ctx.WriteString("default ")
	node.Expr.Format(ctx)
}

func NewAttributeDefault(expr Expr, buf *buffer.Buffer) *AttributeDefault {
	a := buffer.Alloc[AttributeDefault](buf)
	a.Expr = expr
	return a
}

type AttributeAutoIncrement struct {
	columnAttributeImpl
	IsAutoIncrement bool
}

func (node *AttributeAutoIncrement) Format(ctx *FmtCtx) {
	ctx.WriteString("auto_increment")
}

func NewAttributeAutoIncrement(buf *buffer.Buffer) *AttributeAutoIncrement {
	a := buffer.Alloc[AttributeAutoIncrement](buf)
	return a
}

type AttributeUniqueKey struct {
	columnAttributeImpl
}

func (node *AttributeUniqueKey) Format(ctx *FmtCtx) {
	ctx.WriteString("unique key")
}

func NewAttributeUniqueKey(buf *buffer.Buffer) *AttributeUniqueKey {
	a := buffer.Alloc[AttributeUniqueKey](buf)
	return a
}

type AttributeUnique struct {
	columnAttributeImpl
}

func (node *AttributeUnique) Format(ctx *FmtCtx) {
	ctx.WriteString("unique")
}

func NewAttributeUnique(buf *buffer.Buffer) *AttributeUnique {
	a := buffer.Alloc[AttributeUnique](buf)
	return a
}

type AttributeKey struct {
	columnAttributeImpl
}

func (node *AttributeKey) Format(ctx *FmtCtx) {
	ctx.WriteString("key")
}

func NewAttributeKey(buf *buffer.Buffer) *AttributeKey {
	a := buffer.Alloc[AttributeKey](buf)
	return a
}

type AttributePrimaryKey struct {
	columnAttributeImpl
}

func (node *AttributePrimaryKey) Format(ctx *FmtCtx) {
	ctx.WriteString("primary key")
}

func NewAttributePrimaryKey(buf *buffer.Buffer) *AttributePrimaryKey {
	a := buffer.Alloc[AttributePrimaryKey](buf)
	return a
}

type AttributeComment struct {
	columnAttributeImpl
	CMT Expr
}

func (node *AttributeComment) Format(ctx *FmtCtx) {
	ctx.WriteString("comment ")
	node.CMT.Format(ctx)
}

func NewAttributeComment(c Expr, buf *buffer.Buffer) *AttributeComment {
	a := buffer.Alloc[AttributeComment](buf)
	a.CMT = c
	return a
}

type AttributeCollate struct {
	columnAttributeImpl
	Collate *BufString
}

func (node *AttributeCollate) Format(ctx *FmtCtx) {
	ctx.WriteString("collate ")
	ctx.WriteString(node.Collate.Get())
}

func NewAttributeCollate(c string, buf *buffer.Buffer) *AttributeCollate {
	a := buffer.Alloc[AttributeCollate](buf)
	bc := NewBufString(c)
	buf.Pin(bc)
	a.Collate = bc
	return a
}

type AttributeColumnFormat struct {
	columnAttributeImpl
	ColumnFormat *BufString
}

func (node *AttributeColumnFormat) Format(ctx *FmtCtx) {
	ctx.WriteString("format ")
	ctx.WriteString(node.ColumnFormat.Get())
}

func NewAttributeColumnFormat(cf string, buf *buffer.Buffer) *AttributeColumnFormat {
	a := buffer.Alloc[AttributeColumnFormat](buf)
	bc := NewBufString(cf)
	buf.Pin(bc)
	a.ColumnFormat = bc
	return a
}

type AttributeStorage struct {
	columnAttributeImpl
	Storage *BufString
}

func (node *AttributeStorage) Format(ctx *FmtCtx) {
	ctx.WriteString("storage ")
	ctx.WriteString(node.Storage.Get())
}

func NewAttributeStorage(st string, buf *buffer.Buffer) *AttributeStorage {
	a := buffer.Alloc[AttributeStorage](buf)
	s := NewBufString(st)
	buf.Pin(s)
	a.Storage = s
	return a
}

type AttributeCheckConstraint struct {
	columnAttributeImpl
	Name     *BufString
	Expr     Expr
	Enforced bool
}

func (node *AttributeCheckConstraint) Format(ctx *FmtCtx) {
	ctx.WriteString("constraint")
	if node.Name.Get() != "" {
		ctx.WriteByte(' ')
		ctx.WriteString(node.Name.Get())
	}
	ctx.WriteString(" check")
	ctx.WriteString(" (")
	node.Expr.Format(ctx)
	ctx.WriteString(") ")

	if node.Enforced {
		ctx.WriteString("enforced")
	} else {
		ctx.WriteString("not enforced")
	}
}

func NewAttributeCheck(e Expr, f bool, name string, buf *buffer.Buffer) *AttributeCheckConstraint {
	a := buffer.Alloc[AttributeCheckConstraint](buf)
	n := NewBufString(name)
	buf.Pin(n)
	a.Name = n
	a.Expr = e
	a.Enforced = f
	return a
}

type AttributeGeneratedAlways struct {
	columnAttributeImpl
	Expr   Expr
	Stored bool
}

func (node *AttributeGeneratedAlways) Format(ctx *FmtCtx) {
	node.Expr.Format(ctx)
}

func NewAttributeGeneratedAlways(e Expr, s bool, buf *buffer.Buffer) *AttributeGeneratedAlways {
	a := buffer.Alloc[AttributeGeneratedAlways](buf)
	a.Expr = e
	a.Stored = s
	return a
}

type AttributeLowCardinality struct {
	columnAttributeImpl
}

func (node *AttributeLowCardinality) Format(ctx *FmtCtx) {
	ctx.WriteString("low_cardinality")
}

func NewAttributeLowCardinality(buf *buffer.Buffer) *AttributeLowCardinality {
	a := buffer.Alloc[AttributeLowCardinality](buf)
	return a
}

type KeyPart struct {
	columnAttributeImpl
	ColName   *UnresolvedName
	Length    int
	Direction Direction // asc or desc
	Expr      Expr
}

func (node *KeyPart) Format(ctx *FmtCtx) {
	if node.ColName != nil {
		node.ColName.Format(ctx)
	}
	if node.Length != 0 {
		ctx.WriteByte('(')
		if node.Length == -1 {
			ctx.WriteString("0")
		} else {
			ctx.WriteString(strconv.Itoa(node.Length))
		}
		ctx.WriteByte(')')
	}
	if node.Direction != DefaultDirection {
		ctx.WriteByte(' ')
		ctx.WriteString(node.Direction.String())
		return
	}
	if node.Expr != nil {
		ctx.WriteByte('(')
		node.Expr.Format(ctx)
		ctx.WriteByte(')')
		if node.Direction != DefaultDirection {
			ctx.WriteByte(' ')
			ctx.WriteString(node.Direction.String())
		}
	}
}

func NewKeyPart(c *UnresolvedName, l int, e Expr, d Direction, buf *buffer.Buffer) *KeyPart {
	k := buffer.Alloc[KeyPart](buf)
	k.ColName = c
	k.Length = l
	k.Expr = e
	k.Direction = d
	return k
}

// in reference definition
type MatchType int

func (node *MatchType) ToString() string {
	switch *node {
	case MATCH_FULL:
		return "full"
	case MATCH_PARTIAL:
		return "partial"
	case MATCH_SIMPLE:
		return "simple"
	default:
		return "Unknown MatchType"
	}
}

const (
	MATCH_INVALID MatchType = iota
	MATCH_FULL
	MATCH_PARTIAL
	MATCH_SIMPLE
)

type ReferenceOptionType int

func (node *ReferenceOptionType) ToString() string {
	switch *node {
	case REFERENCE_OPTION_RESTRICT:
		return "restrict"
	case REFERENCE_OPTION_CASCADE:
		return "cascade"
	case REFERENCE_OPTION_SET_NULL:
		return "set null"
	case REFERENCE_OPTION_NO_ACTION:
		return "no action"
	case REFERENCE_OPTION_SET_DEFAULT:
		return "set default"
	default:
		return "Unknown ReferenceOptionType"
	}
}

// Reference option
const (
	REFERENCE_OPTION_INVALID ReferenceOptionType = iota
	REFERENCE_OPTION_RESTRICT
	REFERENCE_OPTION_CASCADE
	REFERENCE_OPTION_SET_NULL
	REFERENCE_OPTION_NO_ACTION
	REFERENCE_OPTION_SET_DEFAULT
)

type AttributeReference struct {
	columnAttributeImpl
	TableName *TableName
	KeyParts  []*KeyPart
	Match     MatchType
	OnDelete  ReferenceOptionType
	OnUpdate  ReferenceOptionType
}

func (node *AttributeReference) Format(ctx *FmtCtx) {
	ctx.WriteString("references ")
	node.TableName.Format(ctx)
	if node.KeyParts != nil {
		ctx.WriteByte('(')
		prefix := ""
		for _, k := range node.KeyParts {
			ctx.WriteString(prefix)
			k.Format(ctx)
			prefix = ", "
		}
		ctx.WriteByte(')')
	}
	if node.Match != MATCH_INVALID {
		ctx.WriteString(" match ")
		ctx.WriteString(node.Match.ToString())
	}
	if node.OnDelete != REFERENCE_OPTION_INVALID {
		ctx.WriteString(" on delete ")
		ctx.WriteString(node.OnDelete.ToString())
	}
	if node.OnUpdate != REFERENCE_OPTION_INVALID {
		ctx.WriteString(" on update ")
		ctx.WriteString(node.OnUpdate.ToString())
	}
}

func NewAttributeReference(t *TableName, kps []*KeyPart, ma MatchType,
	od ReferenceOptionType, ou ReferenceOptionType, buf *buffer.Buffer) *AttributeReference {
	a := buffer.Alloc[AttributeReference](buf)
	a.TableName = t
	a.KeyParts = kps
	a.Match = ma
	a.OnDelete = od
	a.OnUpdate = ou
	return a
}

type ReferenceOnRecord struct {
	OnDelete ReferenceOptionType
	OnUpdate ReferenceOptionType
}

func NewReferenceOnRecord(dup, up ReferenceOptionType, buf *buffer.Buffer) *ReferenceOnRecord {
	re := buffer.Alloc[ReferenceOnRecord](buf)
	re.OnDelete = dup
	re.OnUpdate = up
	return re
}

type AttributeAutoRandom struct {
	columnAttributeImpl
	BitLength int
}

func NewAttributeAutoRandom(b int, buf *buffer.Buffer) *AttributeAutoRandom {
	a := buffer.Alloc[AttributeAutoRandom](buf)
	a.BitLength = b
	return a
}

type AttributeOnUpdate struct {
	columnAttributeImpl
	Expr Expr
}

func (node *AttributeOnUpdate) Format(ctx *FmtCtx) {
	ctx.WriteString("on update ")
	node.Expr.Format(ctx)
}

func NewAttributeOnUpdate(e Expr, buf *buffer.Buffer) *AttributeOnUpdate {
	a := buffer.Alloc[AttributeOnUpdate](buf)
	a.Expr = e
	return a
}

type IndexType int

func (it IndexType) ToString() string {
	switch it {
	case INDEX_TYPE_BTREE:
		return "btree"
	case INDEX_TYPE_HASH:
		return "hash"
	case INDEX_TYPE_RTREE:
		return "rtree"
	case INDEX_TYPE_BSI:
		return "bsi"
	case INDEX_TYPE_ZONEMAP:
		return "zonemap"
	default:
		return "Unknown IndexType"
	}
}

const (
	INDEX_TYPE_INVALID IndexType = iota
	INDEX_TYPE_BTREE
	INDEX_TYPE_HASH
	INDEX_TYPE_RTREE
	INDEX_TYPE_BSI
	INDEX_TYPE_ZONEMAP
)

type VisibleType int

const (
	VISIBLE_TYPE_INVALID VisibleType = iota
	VISIBLE_TYPE_VISIBLE
	VISIBLE_TYPE_INVISIBLE
)

func (vt VisibleType) ToString() string {
	switch vt {
	case VISIBLE_TYPE_VISIBLE:
		return "visible"
	case VISIBLE_TYPE_INVISIBLE:
		return "invisible"
	default:
		return "Unknown VisibleType"
	}
}

type IndexOption struct {
	NodeFormatter
	KeyBlockSize             uint64
	IType                    IndexType
	ParserName               *BufString
	Comment                  *BufString
	Visible                  VisibleType
	EngineAttribute          *BufString
	SecondaryEngineAttribute string
}

// Must follow the following sequence when test
func (node *IndexOption) Format(ctx *FmtCtx) {
	if node.KeyBlockSize != 0 {
		ctx.WriteString("KEY_BLOCK_SIZE ")
		ctx.WriteString(strconv.FormatUint(node.KeyBlockSize, 10))
		ctx.WriteByte(' ')
	}
	if node.ParserName.Get() != "" {
		ctx.WriteString("with parser ")
		ctx.WriteString(node.ParserName.Get())
		ctx.WriteByte(' ')
	}
	if node.Comment.Get() != "" {
		ctx.WriteString("comment ")
		ctx.WriteString(node.Comment.Get())
		ctx.WriteByte(' ')
	}
	if node.Visible != VISIBLE_TYPE_INVALID {
		ctx.WriteString(node.Visible.ToString())
	}
}

func NewIndexOption2(kb uint64, co, pa string, vi VisibleType, buf *buffer.Buffer) *IndexOption {
	i := buffer.Alloc[IndexOption](buf)
	i.KeyBlockSize = kb
	bComment := NewBufString(co)
	buf.Pin(bComment)
	i.Comment = bComment
	bParserName := NewBufString(pa)
	buf.Pin(bParserName)
	i.ParserName = bParserName
	i.Visible = vi
	return i
}

func NewIndexOption(buf *buffer.Buffer) *IndexOption {
	i := buffer.Alloc[IndexOption](buf)
	return i
}

type PrimaryKeyIndex struct {
	tableDefImpl
	KeyParts         []*KeyPart
	Name             *BufString
	ConstraintSymbol *BufString
	Empty            bool
	IndexOption      *IndexOption
}

func (node *PrimaryKeyIndex) Format(ctx *FmtCtx) {
	if node.ConstraintSymbol.Get() != "" {
		ctx.WriteString("constraint " + node.ConstraintSymbol.Get() + " ")
	}
	ctx.WriteString("primary key")
	if node.Name.Get() != "" {
		ctx.WriteByte(' ')
		ctx.WriteString(node.Name.Get())
	}
	if !node.Empty {
		ctx.WriteString(" using none")
	}
	if node.KeyParts != nil {
		prefix := " ("
		for _, k := range node.KeyParts {
			ctx.WriteString(prefix)
			k.Format(ctx)
			prefix = ", "
		}
		ctx.WriteByte(')')
	}
	if node.IndexOption != nil {
		ctx.WriteByte(' ')
		node.IndexOption.Format(ctx)
	}
}

func NewPrimaryKeyIndex(keyParts []*KeyPart, name string, empty bool, indexOption *IndexOption, buf *buffer.Buffer) *PrimaryKeyIndex {
	p := buffer.Alloc[PrimaryKeyIndex](buf)
	p.KeyParts = keyParts
	bName := NewBufString(name)
	buf.Pin(bName)
	p.Name = bName
	p.Empty = empty
	p.IndexOption = indexOption
	return p
}

type Index struct {
	tableDefImpl
	IfNotExists bool
	KeyParts    []*KeyPart
	Name        *BufString
	KeyType     IndexType
	IndexOption *IndexOption
}

func (node *Index) Format(ctx *FmtCtx) {
	ctx.WriteString("index")
	if node.IfNotExists {
		ctx.WriteString(" if not exists")
	}
	if node.Name.Get() != "" {
		ctx.WriteByte(' ')
		ctx.WriteString(node.Name.Get())
	}
	if node.KeyType != INDEX_TYPE_INVALID {
		ctx.WriteString(" using ")
		ctx.WriteString(node.KeyType.ToString())
	}
	if node.KeyParts != nil {
		prefix := " ("
		for _, k := range node.KeyParts {
			ctx.WriteString(prefix)
			k.Format(ctx)
			prefix = ", "
		}
		ctx.WriteByte(')')
	}
	if node.IndexOption != nil {
		ctx.WriteByte(' ')
		node.IndexOption.Format(ctx)
	}
}

func NewIndex(ifn bool, keyParts []*KeyPart, name string, keyType IndexType, indexOption *IndexOption, buf *buffer.Buffer) *Index {
	i := buffer.Alloc[Index](buf)
	i.IfNotExists = ifn
	i.KeyParts = keyParts
	bName := NewBufString(name)
	buf.Pin(bName)
	i.Name = bName
	i.KeyType = keyType
	i.IndexOption = indexOption
	return i
}

type UniqueIndex struct {
	tableDefImpl
	KeyParts         []*KeyPart
	Name             *BufString
	ConstraintSymbol *BufString
	Empty            bool
	IndexOption      *IndexOption
}

func (node *UniqueIndex) Format(ctx *FmtCtx) {
	if node.ConstraintSymbol.Get() != "" {
		ctx.WriteString("constraint " + node.ConstraintSymbol.Get() + " ")
	}
	ctx.WriteString("unique key")
	if node.Name.Get() != "" {
		ctx.WriteByte(' ')
		ctx.WriteString(node.Name.Get())
	}
	if !node.Empty {
		ctx.WriteString(" using none")
	}
	if node.KeyParts != nil {
		prefix := " ("
		for _, k := range node.KeyParts {
			ctx.WriteString(prefix)
			k.Format(ctx)
			prefix = ", "
		}
		ctx.WriteByte(')')
	}
	if node.IndexOption != nil {
		ctx.WriteByte(' ')
		node.IndexOption.Format(ctx)
	}
}

func (node *UniqueIndex) GetIndexName() string {
	if len(node.Name.Get()) != 0 {
		return node.Name.Get()
	} else {
		return node.ConstraintSymbol.Get()
	}
}

func NewUniqueIndex(keyParts []*KeyPart, name string, empty bool, indexOption *IndexOption, buf *buffer.Buffer) *UniqueIndex {
	u := buffer.Alloc[UniqueIndex](buf)
	u.KeyParts = keyParts
	bName := NewBufString(name)
	buf.Pin(bName)
	u.Name = bName
	u.Empty = empty
	u.IndexOption = indexOption
	return u
}

type ForeignKey struct {
	tableDefImpl
	IfNotExists      bool
	KeyParts         []*KeyPart
	Name             *BufString
	ConstraintSymbol *BufString
	Refer            *AttributeReference
	Empty            bool
}

func (node *ForeignKey) Format(ctx *FmtCtx) {
	if node.ConstraintSymbol.Get() != "" {
		ctx.WriteString("constraint " + node.ConstraintSymbol.Get() + " ")
	}
	ctx.WriteString("foreign key")
	if node.IfNotExists {
		ctx.WriteString(" if not exists")
	}
	if node.Name.Get() != "" {
		ctx.WriteByte(' ')
		ctx.WriteString(node.Name.Get())
	}
	if !node.Empty {
		ctx.WriteString(" using none")
	}
	if node.KeyParts != nil {
		prefix := " ("
		for _, k := range node.KeyParts {
			ctx.WriteString(prefix)
			k.Format(ctx)
			prefix = ", "
		}
		ctx.WriteByte(')')
	}
	if node.Refer != nil {
		ctx.WriteByte(' ')
		node.Refer.Format(ctx)
	}
}

func NewForeignKey(ifNotExists bool, keyParts []*KeyPart, name string, refer *AttributeReference, empty bool, buf *buffer.Buffer) *ForeignKey {
	f := buffer.Alloc[ForeignKey](buf)
	f.IfNotExists = ifNotExists
	f.KeyParts = keyParts
	bName := NewBufString(name)
	buf.Pin(bName)
	f.Name = bName
	f.Refer = refer
	f.Empty = empty
	return f
}

type FullTextIndex struct {
	tableDefImpl
	KeyParts    []*KeyPart
	Name        *BufString
	Empty       bool
	IndexOption *IndexOption
}

func (node *FullTextIndex) Format(ctx *FmtCtx) {
	ctx.WriteString("fulltext")
	if node.Name.Get() != "" {
		ctx.WriteByte(' ')
		ctx.WriteString(node.Name.Get())
	}
	if !node.Empty {
		ctx.WriteString(" using none")
	}
	if node.KeyParts != nil {
		prefix := " ("
		for _, k := range node.KeyParts {
			ctx.WriteString(prefix)
			k.Format(ctx)
			prefix = ", "
		}
		ctx.WriteByte(')')
	}
	if node.IndexOption != nil {
		ctx.WriteByte(' ')
		node.IndexOption.Format(ctx)
	}
}

func NewFullTextIndex(keyParts []*KeyPart, name string, empty bool, indexOption *IndexOption, buf *buffer.Buffer) *FullTextIndex {
	f := buffer.Alloc[FullTextIndex](buf)
	f.KeyParts = keyParts
	bName := NewBufString(name)
	buf.Pin(bName)
	f.Name = bName
	f.Empty = empty
	f.IndexOption = indexOption
	return f
}

type CheckIndex struct {
	tableDefImpl
	Expr     Expr
	Enforced bool
}

func (node *CheckIndex) Format(ctx *FmtCtx) {
	ctx.WriteString("check (")
	node.Expr.Format(ctx)
	ctx.WriteByte(')')
	if node.Enforced {
		ctx.WriteString(" enforced")
	}
}

func NewCheckIndex(expr Expr, enforced bool, buf *buffer.Buffer) *CheckIndex {
	c := buffer.Alloc[CheckIndex](buf)
	c.Expr = expr
	c.Enforced = enforced
	return c
}

type TableOption interface {
	AlterTableOption
}

type tableOptionImpl struct {
	TableOption
}

type TableOptionProperties struct {
	tableOptionImpl
	Preperties []*Property
}

func (node *TableOptionProperties) Format(ctx *FmtCtx) {
	ctx.WriteString("properties")
	if node.Preperties != nil {
		prefix := "("
		for _, p := range node.Preperties {
			ctx.WriteString(prefix)
			p.Format(ctx)
			prefix = ", "
		}
		ctx.WriteByte(')')
	}
}

func NewTableOptionProperties(preperties []*Property, buf *buffer.Buffer) *TableOptionProperties {
	t := buffer.Alloc[TableOptionProperties](buf)
	t.Preperties = preperties
	return t
}

type Property struct {
	Key   *BufString
	Value *BufString
}

func NewProperty(key string, value string, buf *buffer.Buffer) *Property {
	p := buffer.Alloc[Property](buf)
	bKey := NewBufString(key)
	bValue := NewBufString(value)
	buf.Pin(bKey, bValue)
	p.Key = bKey
	p.Value = bValue
	return p
}

func (node *Property) Format(ctx *FmtCtx) {
	ctx.WriteString(node.Key.Get())
	ctx.WriteString(" = ")
	ctx.WriteString(node.Value.Get())
}

type TableOptionEngine struct {
	tableOptionImpl
	Engine *BufString
}

func (node *TableOptionEngine) Format(ctx *FmtCtx) {
	ctx.WriteString("engine = ")
	ctx.WriteString(node.Engine.Get())
}

func NewTableOptionEngine(engine string, buf *buffer.Buffer) *TableOptionEngine {
	t := buffer.Alloc[TableOptionEngine](buf)
	bEngine := NewBufString(engine)
	buf.Pin(bEngine)
	t.Engine = bEngine
	return t
}

type TableOptionEngineAttr struct {
	tableOptionImpl
	Engine *BufString
}

func (node *TableOptionEngineAttr) Format(ctx *FmtCtx) {
	ctx.WriteString("ENGINE_ATTRIBUTE = ")
	ctx.WriteString(node.Engine.Get())
}

func NewTableOptionEngineAttr(engine string, buf *buffer.Buffer) *TableOptionEngineAttr {
	t := buffer.Alloc[TableOptionEngineAttr](buf)
	bEngine := NewBufString(engine)
	buf.Pin(bEngine)
	t.Engine = bEngine
	return t
}

type TableOptionInsertMethod struct {
	tableOptionImpl
	Method *BufString
}

func (node *TableOptionInsertMethod) Format(ctx *FmtCtx) {
	ctx.WriteString("INSERT_METHOD = ")
	ctx.WriteString(node.Method.Get())
}

func NewTableOptionInsertMethod(s string, buf *buffer.Buffer) *TableOptionInsertMethod {
	t := buffer.Alloc[TableOptionInsertMethod](buf)
	bMethod := NewBufString(s)
	buf.Pin(bMethod)
	t.Method = bMethod
	return t
}

type TableOptionSecondaryEngine struct {
	tableOptionImpl
	Engine *BufString
}

func (node *TableOptionSecondaryEngine) Format(ctx *FmtCtx) {
	ctx.WriteString("engine = ")
	ctx.WriteString(node.Engine.Get())
}

func NewTableOptionSecondaryEngine(s string, buf *buffer.Buffer) *TableOptionSecondaryEngine {
	t := buffer.Alloc[TableOptionSecondaryEngine](buf)
	bEngine := NewBufString(s)
	buf.Pin(bEngine)
	t.Engine = bEngine
	return t
}

type TableOptionSecondaryEngineNull struct {
	tableOptionImpl
}

func NewTableOptionSecondaryEngineNull(buf *buffer.Buffer) *TableOptionSecondaryEngineNull {
	t := buffer.Alloc[TableOptionSecondaryEngineNull](buf)
	return t
}

type TableOptionCharset struct {
	tableOptionImpl
	Charset *BufString
}

func (node *TableOptionCharset) Format(ctx *FmtCtx) {
	ctx.WriteString("charset = ")
	ctx.WriteString(node.Charset.Get())
}

func NewTableOptionCharset(s string, buf *buffer.Buffer) *TableOptionCharset {
	t := buffer.Alloc[TableOptionCharset](buf)
	bCharset := NewBufString(s)
	buf.Pin(bCharset)
	t.Charset = bCharset
	return t
}

type TableOptionCollate struct {
	tableOptionImpl
	Collate *BufString
}

func (node *TableOptionCollate) Format(ctx *FmtCtx) {
	ctx.WriteString("Collate = ")
	ctx.WriteString(node.Collate.Get())
}

func NewTableOptionCollate(s string, buf *buffer.Buffer) *TableOptionCollate {
	t := buffer.Alloc[TableOptionCollate](buf)
	bCollate := NewBufString(s)
	buf.Pin(bCollate)
	t.Collate = bCollate
	return t
}

type TableOptionAUTOEXTEND_SIZE struct {
	tableOptionImpl
	Value uint64
}

func (node *TableOptionAUTOEXTEND_SIZE) Format(ctx *FmtCtx) {
	ctx.WriteString("AUTOEXTEND_SIZE = ")
	ctx.WriteString(strconv.FormatUint(node.Value, 10))
}

func NewTableOptionAUTOEXTEND_SIZE(value uint64, buf *buffer.Buffer) *TableOptionAUTOEXTEND_SIZE {
	t := buffer.Alloc[TableOptionAUTOEXTEND_SIZE](buf)
	t.Value = value
	return t
}

type TableOptionAutoIncrement struct {
	tableOptionImpl
	Value uint64
}

func (node *TableOptionAutoIncrement) Format(ctx *FmtCtx) {
	ctx.WriteString("auto_increment = ")
	ctx.WriteString(strconv.FormatUint(node.Value, 10))
}

func NewTableOptionAutoIncrement(value uint64, buf *buffer.Buffer) *TableOptionAutoIncrement {
	t := buffer.Alloc[TableOptionAutoIncrement](buf)
	t.Value = value
	return t
}

type TableOptionComment struct {
	tableOptionImpl
	Comment *BufString
}

func (node *TableOptionComment) Format(ctx *FmtCtx) {
	ctx.WriteString("comment = ")
	ctx.WriteString(node.Comment.Get())
}

func NewTableOptionComment(comment string, buf *buffer.Buffer) *TableOptionComment {
	t := buffer.Alloc[TableOptionComment](buf)
	bComment := NewBufString(comment)
	buf.Pin(bComment)
	t.Comment = bComment
	return t
}

type TableOptionAvgRowLength struct {
	tableOptionImpl
	Length uint64
}

func (node *TableOptionAvgRowLength) Format(ctx *FmtCtx) {
	ctx.WriteString("avg_row_length = ")
	ctx.WriteString(strconv.FormatUint(node.Length, 10))
}

func NewTableOptionAvgRowLength(length uint64, buf *buffer.Buffer) *TableOptionAvgRowLength {
	t := buffer.Alloc[TableOptionAvgRowLength](buf)
	t.Length = length
	return t
}

type TableOptionChecksum struct {
	tableOptionImpl
	Value uint64
}

func (node *TableOptionChecksum) Format(ctx *FmtCtx) {
	ctx.WriteString("checksum = ")
	ctx.WriteString(strconv.FormatUint(node.Value, 10))
}

func NewTableOptionChecksum(value uint64, buf *buffer.Buffer) *TableOptionChecksum {
	t := buffer.Alloc[TableOptionChecksum](buf)
	t.Value = value
	return t
}

type TableOptionCompression struct {
	tableOptionImpl
	Compression *BufString
}

func (node *TableOptionCompression) Format(ctx *FmtCtx) {
	ctx.WriteString("compression = ")
	ctx.WriteString(node.Compression.Get())
}

func NewTableOptionCompression(compression string, buf *buffer.Buffer) *TableOptionCompression {
	t := buffer.Alloc[TableOptionCompression](buf)
	bCompression := NewBufString(compression)
	buf.Pin(bCompression)
	t.Compression = bCompression
	return t
}

type TableOptionConnection struct {
	tableOptionImpl
	Connection *BufString
}

func (node *TableOptionConnection) Format(ctx *FmtCtx) {
	ctx.WriteString("connection = ")
	ctx.WriteString(node.Connection.Get())
}

func NewTableOptionConnection(connection string, buf *buffer.Buffer) *TableOptionConnection {
	t := buffer.Alloc[TableOptionConnection](buf)
	bConnection := NewBufString(connection)
	buf.Pin(bConnection)
	t.Connection = bConnection
	return t
}

type TableOptionPassword struct {
	tableOptionImpl
	Password *BufString
}

func (node *TableOptionPassword) Format(ctx *FmtCtx) {
	ctx.WriteString("password = ")
	ctx.WriteString(node.Password.Get())
}

func NewTableOptionPassword(password string, buf *buffer.Buffer) *TableOptionPassword {
	t := buffer.Alloc[TableOptionPassword](buf)
	bPassword := NewBufString(password)
	buf.Pin(bPassword)
	t.Password = bPassword
	return t
}

type TableOptionKeyBlockSize struct {
	tableOptionImpl
	Value uint64
}

func (node *TableOptionKeyBlockSize) Format(ctx *FmtCtx) {
	ctx.WriteString("key_block_size = ")
	ctx.WriteString(strconv.FormatUint(node.Value, 10))
}

func NewTableOptionKeyBlockSize(value uint64, buf *buffer.Buffer) *TableOptionKeyBlockSize {
	t := buffer.Alloc[TableOptionKeyBlockSize](buf)
	t.Value = value
	return t
}

type TableOptionMaxRows struct {
	tableOptionImpl
	Value uint64
}

func (node *TableOptionMaxRows) Format(ctx *FmtCtx) {
	ctx.WriteString("max_rows = ")
	ctx.WriteString(strconv.FormatUint(node.Value, 10))
}

func NewTableOptionMaxRows(value uint64, buf *buffer.Buffer) *TableOptionMaxRows {
	t := buffer.Alloc[TableOptionMaxRows](buf)
	t.Value = value
	return t
}

type TableOptionMinRows struct {
	tableOptionImpl
	Value uint64
}

func (node *TableOptionMinRows) Format(ctx *FmtCtx) {
	ctx.WriteString("min_rows = ")
	ctx.WriteString(strconv.FormatUint(node.Value, 10))
}

func NewTableOptionMinRows(value uint64, buf *buffer.Buffer) *TableOptionMinRows {
	t := buffer.Alloc[TableOptionMinRows](buf)
	t.Value = value
	return t
}

type TableOptionDelayKeyWrite struct {
	tableOptionImpl
	Value uint64
}

func (node *TableOptionDelayKeyWrite) Format(ctx *FmtCtx) {
	ctx.WriteString("key_write = ")
	ctx.WriteString(strconv.FormatUint(node.Value, 10))
}

func NewTableOptionDelayKeyWrite(value uint64, buf *buffer.Buffer) *TableOptionDelayKeyWrite {
	t := buffer.Alloc[TableOptionDelayKeyWrite](buf)
	t.Value = value
	return t
}

type RowFormatType uint64

const (
	ROW_FORMAT_DEFAULT RowFormatType = iota
	ROW_FORMAT_DYNAMIC
	ROW_FORMAT_FIXED
	ROW_FORMAT_COMPRESSED
	ROW_FORMAT_REDUNDANT
	ROW_FORMAT_COMPACT
)

func (node *RowFormatType) ToString() string {
	switch *node {
	case ROW_FORMAT_DEFAULT:
		return "default"
	case ROW_FORMAT_DYNAMIC:
		return "dynamic"
	case ROW_FORMAT_FIXED:
		return "fixed"
	case ROW_FORMAT_COMPRESSED:
		return "compressed"
	case ROW_FORMAT_REDUNDANT:
		return "redundant"
	case ROW_FORMAT_COMPACT:
		return "compact"
	default:
		return "Unknown RowFormatType"
	}
}

type TableOptionRowFormat struct {
	tableOptionImpl
	Value RowFormatType
}

func (node *TableOptionRowFormat) Format(ctx *FmtCtx) {
	ctx.WriteString("row_format = ")
	ctx.WriteString(node.Value.ToString())
}

func NewTableOptionRowFormat(value RowFormatType, buf *buffer.Buffer) *TableOptionRowFormat {
	t := buffer.Alloc[TableOptionRowFormat](buf)
	t.Value = value
	return t
}

type TableOptionStartTrans struct {
	tableOptionImpl
	Value bool
}

func (node *TableOptionStartTrans) Format(ctx *FmtCtx) {
	ctx.WriteString("START TRANSACTION")
}

func NewTableOptionStartTrans(value bool, buf *buffer.Buffer) *TableOptionStartTrans {
	t := buffer.Alloc[TableOptionStartTrans](buf)
	t.Value = value
	return t
}

type TableOptionSecondaryEngineAttr struct {
	tableOptionImpl
	Attr *BufString
}

func (node *TableOptionSecondaryEngineAttr) Format(ctx *FmtCtx) {
	ctx.WriteString("SECONDARY_ENGINE_ATTRIBUTE = ")
	ctx.WriteString(node.Attr.Get())
}

func NewTableOptionSecondaryEngineAttr(attr string, buf *buffer.Buffer) *TableOptionSecondaryEngineAttr {
	t := buffer.Alloc[TableOptionSecondaryEngineAttr](buf)
	bAttr := NewBufString(attr)
	buf.Pin(bAttr)
	t.Attr = bAttr
	return t
}

type TableOptionStatsPersistent struct {
	tableOptionImpl
	Value   uint64
	Default bool
}

func (node *TableOptionStatsPersistent) Format(ctx *FmtCtx) {
	ctx.WriteString("stats_persistent = ")
	if node.Default {
		ctx.WriteString("default")
	} else {
		ctx.WriteString(strconv.FormatUint(node.Value, 10))
	}
}

func NewTableOptionStatsPersistent(value uint64, default_ bool, buf *buffer.Buffer) *TableOptionStatsPersistent {
	t := buffer.Alloc[TableOptionStatsPersistent](buf)
	t.Value = value
	t.Default = default_
	return t
}

type TableOptionStatsAutoRecalc struct {
	tableOptionImpl
	Value   uint64
	Default bool //false -- see Value; true -- Value is useless
}

func (node *TableOptionStatsAutoRecalc) Format(ctx *FmtCtx) {
	ctx.WriteString("stats_auto_recalc = ")
	if node.Default {
		ctx.WriteString("default")
	} else {
		ctx.WriteString(strconv.FormatUint(node.Value, 10))
	}
}

func NewTableOptionStatsAutoRecalc(value uint64, default_ bool, buf *buffer.Buffer) *TableOptionStatsAutoRecalc {
	t := buffer.Alloc[TableOptionStatsAutoRecalc](buf)
	t.Value = value
	t.Default = default_
	return t
}

type TableOptionPackKeys struct {
	tableOptionImpl
	Value   int64
	Default bool
}

func (node *TableOptionPackKeys) Format(ctx *FmtCtx) {
	ctx.WriteString("pack_keys = ")
	if node.Default {
		ctx.WriteString("default")
	} else {
		ctx.WriteString(strconv.FormatInt(node.Value, 10))
	}
}

func NewTableOptionPackKeys(value int64, defalut_ bool, buf *buffer.Buffer) *TableOptionPackKeys {
	t := buffer.Alloc[TableOptionPackKeys](buf)
	t.Value = value
	t.Default = defalut_
	return t
}

type TableOptionTablespace struct {
	tableOptionImpl
	Name       *BufString
	StorageOpt *BufString
}

func (node *TableOptionTablespace) Format(ctx *FmtCtx) {
	ctx.WriteString("tablespace = ")
	ctx.WriteString(node.Name.Get())
	ctx.WriteString(node.StorageOpt.Get())
}

func NewTableOptionTablespace(name string, storageOpt string, buf *buffer.Buffer) *TableOptionTablespace {
	t := buffer.Alloc[TableOptionTablespace](buf)
	bName := NewBufString(name)
	bStorageOpt := NewBufString(storageOpt)
	buf.Pin(bName, bStorageOpt)
	t.Name = bName
	t.StorageOpt = bStorageOpt
	return t
}

type TableOptionDataDirectory struct {
	tableOptionImpl
	Dir *BufString
}

func (node *TableOptionDataDirectory) Format(ctx *FmtCtx) {
	ctx.WriteString("data directory = ")
	ctx.WriteString(node.Dir.Get())
}

func NewTableOptionDataDirectory(dir string, buf *buffer.Buffer) *TableOptionDataDirectory {
	t := buffer.Alloc[TableOptionDataDirectory](buf)
	bDir := NewBufString(dir)
	buf.Pin(bDir)
	t.Dir = bDir
	return t
}

type TableOptionIndexDirectory struct {
	tableOptionImpl
	Dir *BufString
}

func (node *TableOptionIndexDirectory) Format(ctx *FmtCtx) {
	ctx.WriteString("index directory = ")
	ctx.WriteString(node.Dir.Get())
}

func NewTableOptionIndexDirectory(d string, buf *buffer.Buffer) *TableOptionIndexDirectory {
	t := buffer.Alloc[TableOptionIndexDirectory](buf)
	bDir := NewBufString(d)
	buf.Pin(bDir)
	t.Dir = bDir
	return t
}

type TableOptionStorageMedia struct {
	tableOptionImpl
	Media *BufString
}

func (node *TableOptionStorageMedia) Format(ctx *FmtCtx) {
	ctx.WriteString("storage media = ")
	ctx.WriteString(node.Media.Get())
}

func NewTableOptionStorageMedia(media string, buf *buffer.Buffer) *TableOptionStorageMedia {
	a := buffer.Alloc[TableOptionStorageMedia](buf)
	bMedia := NewBufString(media)
	buf.Pin(bMedia)
	a.Media = bMedia
	return a
}

type TableOptionStatsSamplePages struct {
	tableOptionImpl
	Value   uint64
	Default bool //false -- see Value; true -- Value is useless
}

func (node *TableOptionStatsSamplePages) Format(ctx *FmtCtx) {
	ctx.WriteString("stats_sample_pages = ")
	if node.Default {
		ctx.WriteString("default")
	} else {
		ctx.WriteString(strconv.FormatUint(node.Value, 10))
	}
}

func NewTableOptionStatsSamplePages(value uint64, default_ bool, buf *buffer.Buffer) *TableOptionStatsSamplePages {
	a := buffer.Alloc[TableOptionStatsSamplePages](buf)
	a.Value = value
	a.Default = default_
	return a
}

type TableOptionUnion struct {
	tableOptionImpl
	Names TableNames
}

func (node *TableOptionUnion) Format(ctx *FmtCtx) {
	ctx.WriteString("union (")
	node.Names.Format(ctx)
	ctx.WriteByte(')')
}

func NewTableOptionUnion(names TableNames, buf *buffer.Buffer) *TableOptionUnion {
	a := buffer.Alloc[TableOptionUnion](buf)
	a.Names = names
	return a
}

type TableOptionEncryption struct {
	tableOptionImpl
	Encryption *BufString
}

func (node *TableOptionEncryption) Format(ctx *FmtCtx) {
	ctx.WriteString("encryption = ")
	ctx.WriteString(node.Encryption.Get())
}

func NewTableOptionEncryption(encryption string, buf *buffer.Buffer) *TableOptionEncryption {
	a := buffer.Alloc[TableOptionEncryption](buf)
	bEncryption := NewBufString(encryption)
	buf.Pin(bEncryption)
	a.Encryption = bEncryption
	return a
}

type PartitionType interface {
	NodeFormatter
}

type partitionTypeImpl struct {
	PartitionType
}

type HashType struct {
	partitionTypeImpl
	Linear bool
	Expr   Expr
}

func (node *HashType) Format(ctx *FmtCtx) {
	if node.Linear {
		ctx.WriteString("linear ")
	}
	ctx.WriteString("hash")
	if node.Expr != nil {
		ctx.WriteString(" (")
		node.Expr.Format(ctx)
		ctx.WriteByte(')')
	}
}

func NewHashType(linear bool, expr Expr, buf *buffer.Buffer) *HashType {
	h := buffer.Alloc[HashType](buf)
	h.Linear = linear
	h.Expr = expr
	return h
}

type KeyType struct {
	partitionTypeImpl
	Linear     bool
	ColumnList []*UnresolvedName
	Algorithm  int64
}

func (node *KeyType) Format(ctx *FmtCtx) {
	if node.Linear {
		ctx.WriteString("linear ")
	}
	ctx.WriteString("key")
	if node.Algorithm != 0 {
		ctx.WriteString(" algorithm = ")
		ctx.WriteString(strconv.FormatInt(node.Algorithm, 10))
	}
	if node.ColumnList != nil {
		prefix := " ("
		for _, c := range node.ColumnList {
			ctx.WriteString(prefix)
			c.Format(ctx)
			prefix = ", "
		}
		ctx.WriteByte(')')
	}
}

func NewKeyType(linear bool, columnList []*UnresolvedName, a int64, buf *buffer.Buffer) *KeyType {
	k := buffer.Alloc[KeyType](buf)
	k.Linear = linear
	k.ColumnList = columnList
	k.Algorithm = a
	return k
}

type RangeType struct {
	partitionTypeImpl
	Expr       Expr
	ColumnList []*UnresolvedName
}

func (node *RangeType) Format(ctx *FmtCtx) {
	ctx.WriteString("range")
	if node.ColumnList != nil {
		prefix := " columns ("
		for _, c := range node.ColumnList {
			ctx.WriteString(prefix)
			c.Format(ctx)
			prefix = ", "
		}
		ctx.WriteByte(')')
	}
	if node.Expr != nil {
		ctx.WriteString("(")
		node.Expr.Format(ctx)
		ctx.WriteByte(')')
	}
}

func NewRangeType(expr Expr, columnList []*UnresolvedName, buf *buffer.Buffer) *RangeType {
	r := buffer.Alloc[RangeType](buf)
	r.Expr = expr
	r.ColumnList = columnList
	return r
}

type ListType struct {
	partitionTypeImpl
	Expr       Expr
	ColumnList []*UnresolvedName
}

func (node *ListType) Format(ctx *FmtCtx) {
	ctx.WriteString("list")
	if node.ColumnList != nil {
		prefix := " columns ("
		for _, c := range node.ColumnList {
			ctx.WriteString(prefix)
			c.Format(ctx)
			prefix = ", "
		}
		ctx.WriteByte(')')
	}
	if node.Expr != nil {
		ctx.WriteString("(")
		node.Expr.Format(ctx)
		ctx.WriteByte(')')
	}
}

func NewListType(expr Expr, columnList []*UnresolvedName, buf *buffer.Buffer) *ListType {
	l := buffer.Alloc[ListType](buf)
	l.Expr = expr
	l.ColumnList = columnList
	return l
}

type PartitionBy struct {
	IsSubPartition bool // for format
	PType          PartitionType
	Num            uint64
}

func (node *PartitionBy) Format(ctx *FmtCtx) {
	node.PType.Format(ctx)
	if node.Num != 0 {
		if node.IsSubPartition {
			ctx.WriteString(" subpartitions ")
		} else {
			ctx.WriteString(" partitions ")
		}
		ctx.WriteString(strconv.FormatUint(node.Num, 10))
	}
}

func NewPartitionBy(pType PartitionType, buf *buffer.Buffer) *PartitionBy {
	p := buffer.Alloc[PartitionBy](buf)
	p.PType = pType
	return p
}

func NewPartitionBy2(is bool, pType PartitionType, nu uint64, buf *buffer.Buffer) *PartitionBy {
	p := buffer.Alloc[PartitionBy](buf)
	p.IsSubPartition = is
	p.PType = pType
	p.Num = nu
	return p
}

//type SubpartitionBy struct {
//	SubPType PartitionType
//	Num uint64
//}

type Values interface {
	NodeFormatter
}

type valuesImpl struct {
	Values
}

type ValuesLessThan struct {
	valuesImpl
	ValueList Exprs
}

func (node *ValuesLessThan) Format(ctx *FmtCtx) {
	ctx.WriteString("values less than (")
	node.ValueList.Format(ctx)
	ctx.WriteByte(')')
}

func NewValuesLessThan(valueList Exprs, buf *buffer.Buffer) *ValuesLessThan {
	v := buffer.Alloc[ValuesLessThan](buf)
	v.ValueList = valueList
	return v
}

type ValuesIn struct {
	valuesImpl
	ValueList Exprs
}

func (node *ValuesIn) Format(ctx *FmtCtx) {
	ctx.WriteString("values in (")
	node.ValueList.Format(ctx)
	ctx.WriteByte(')')
}

func NewValuesIn(valueList Exprs, buf *buffer.Buffer) *ValuesIn {
	v := buffer.Alloc[ValuesIn](buf)
	v.ValueList = valueList
	return v
}

type Partition struct {
	Name    *BufIdentifier
	Values  Values
	Options []TableOption
	Subs    []*SubPartition
}

func (node *Partition) Format(ctx *FmtCtx) {
	ctx.WriteString("partition ")
	ctx.WriteString(string(node.Name.Get()))
	if node.Values != nil {
		ctx.WriteByte(' ')
		node.Values.Format(ctx)
	}
	if node.Options != nil {
		prefix := " "
		for _, t := range node.Options {
			ctx.WriteString(prefix)
			t.Format(ctx)
			prefix = " "
		}
	}
	if node.Subs != nil {
		prefix := " ("
		for _, s := range node.Subs {
			ctx.WriteString(prefix)
			s.Format(ctx)
			prefix = ", "
		}
		ctx.WriteByte(')')
	}
}

func NewPartition(name Identifier, values Values, options []TableOption, subs []*SubPartition, buf *buffer.Buffer) *Partition {
	p := buffer.Alloc[Partition](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)

	p.Name = n
	p.Values = values
	p.Options = options
	p.Subs = subs
	return p
}

type SubPartition struct {
	Name    *BufIdentifier
	Options []TableOption
}

func (node *SubPartition) Format(ctx *FmtCtx) {
	ctx.WriteString("subpartition ")
	ctx.WriteString(string(node.Name.Get()))

	if node.Options != nil {
		prefix := " "
		for _, t := range node.Options {
			ctx.WriteString(prefix)
			t.Format(ctx)
			prefix = " "
		}
	}
}

func NewSubPartition(name Identifier, options []TableOption, buf *buffer.Buffer) *SubPartition {
	s := buffer.Alloc[SubPartition](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)
	s.Name = n
	s.Options = options
	return s
}

type ClusterByOption struct {
	ColumnList []*UnresolvedName
}

func NewClusterByOption(columnList []*UnresolvedName, buf *buffer.Buffer) *ClusterByOption {
	c := buffer.Alloc[ClusterByOption](buf)
	c.ColumnList = columnList
	return c
}

type PartitionOption struct {
	PartBy     *PartitionBy
	SubPartBy  *PartitionBy
	Partitions []*Partition
}

func (node *PartitionOption) Format(ctx *FmtCtx) {
	ctx.WriteString("partition by ")
	node.PartBy.Format(ctx)
	if node.SubPartBy != nil {
		ctx.WriteString(" subpartition by ")
		node.SubPartBy.Format(ctx)
	}
	if node.Partitions != nil {
		prefix := " ("
		for _, p := range node.Partitions {
			ctx.WriteString(prefix)
			p.Format(ctx)
			prefix = ", "
		}
		ctx.WriteByte(')')
	}
}

func NewPartitionOption(partBy *PartitionBy, subPartBy *PartitionBy, partitions []*Partition, buf *buffer.Buffer) *PartitionOption {
	p := buffer.Alloc[PartitionOption](buf)
	p.PartBy = partBy
	p.SubPartBy = subPartBy
	p.Partitions = partitions
	return p
}

type IndexCategory int

func (ic IndexCategory) ToString() string {
	switch ic {
	case INDEX_CATEGORY_UNIQUE:
		return "unique"
	case INDEX_CATEGORY_FULLTEXT:
		return "fulltext"
	case INDEX_CATEGORY_SPATIAL:
		return "spatial"
	default:
		return "Unknown IndexCategory"
	}
}

const (
	INDEX_CATEGORY_NONE IndexCategory = iota
	INDEX_CATEGORY_UNIQUE
	INDEX_CATEGORY_FULLTEXT
	INDEX_CATEGORY_SPATIAL
)

type CreateIndex struct {
	statementImpl
	Name        *BufIdentifier
	Table       *TableName
	IndexCat    IndexCategory
	IfNotExists bool
	KeyParts    []*KeyPart
	IndexOption *IndexOption
	MiscOption  []MiscOption
}

func (node *CreateIndex) Format(ctx *FmtCtx) {
	ctx.WriteString("create ")
	if node.IndexCat != INDEX_CATEGORY_NONE {
		ctx.WriteString(node.IndexCat.ToString())
		ctx.WriteByte(' ')
	}

	ctx.WriteString("index ")
	node.Name.Format(ctx)
	ctx.WriteByte(' ')

	if node.IndexOption != nil && node.IndexOption.IType != INDEX_TYPE_INVALID {
		ctx.WriteString("using ")
		ctx.WriteString(node.IndexOption.IType.ToString())
		ctx.WriteByte(' ')
	}

	ctx.WriteString("on ")
	node.Table.Format(ctx)

	ctx.WriteString(" (")
	if node.KeyParts != nil {
		prefix := ""
		for _, kp := range node.KeyParts {
			ctx.WriteString(prefix)
			kp.Format(ctx)
			prefix = ", "
		}
	}
	ctx.WriteString(")")
	if node.IndexOption != nil {
		ctx.WriteByte(' ')
		node.IndexOption.Format(ctx)
	}
}

func (node *CreateIndex) GetStatementType() string { return "Create Index" }
func (node *CreateIndex) GetQueryType() string     { return QueryTypeDDL }

func NewCreateIndex(name Identifier, table *TableName, indexCat IndexCategory, keyParts []*KeyPart, indexOption *IndexOption, miscOption []MiscOption, buf *buffer.Buffer) *CreateIndex {
	t := buffer.Alloc[CreateIndex](buf)
	n := NewBufIdentifier(name)
	buf.Pin(n)

	t.Name = n
	t.Table = table
	t.IndexCat = indexCat
	t.KeyParts = keyParts
	t.IndexOption = indexOption
	t.MiscOption = miscOption
	return t
}

type MiscOption interface {
	NodeFormatter
}

type miscOption struct {
	MiscOption
}

type AlgorithmDefault struct {
	miscOption
}

type AlgorithmInplace struct {
	miscOption
}

type AlgorithmCopy struct {
	miscOption
}

type LockDefault struct {
	miscOption
}

type LockNone struct {
	miscOption
}

type LockShared struct {
	miscOption
}

type LockExclusive struct {
	miscOption
}

type CreateRole struct {
	statementImpl
	IfNotExists bool
	Roles       []*Role
}

func (node *CreateRole) Format(ctx *FmtCtx) {
	ctx.WriteString("create role ")
	if node.IfNotExists {
		ctx.WriteString("if not exists ")
	}
	prefix := ""
	for _, r := range node.Roles {
		ctx.WriteString(prefix)
		r.Format(ctx)
		prefix = ", "
	}
}

func (node *CreateRole) GetStatementType() string { return "Create Role" }
func (node *CreateRole) GetQueryType() string     { return QueryTypeDCL }

func NewCreateRole(ifNotExists bool, roles []*Role, buf *buffer.Buffer) *CreateRole {
	c := buffer.Alloc[CreateRole](buf)
	c.IfNotExists = ifNotExists
	c.Roles = roles
	return c
}

type Role struct {
	NodeFormatter
	UserName *BufString
}

func (node *Role) Format(ctx *FmtCtx) {
	ctx.WriteString(node.UserName.Get())
}

func NewRole(userName string, buf *buffer.Buffer) *Role {
	r := buffer.Alloc[Role](buf)
	bUserName := NewBufString(userName)
	buf.Pin(bUserName)
	r.UserName = bUserName
	return r
}

type User struct {
	NodeFormatter
	Username   *BufString
	Hostname   *BufString
	AuthOption *AccountIdentified
}

func (node *User) Format(ctx *FmtCtx) {
	ctx.WriteString(node.Username.Get())
	if node.Hostname.Get() != "%" {
		ctx.WriteByte('@')
		ctx.WriteString(node.Hostname.Get())
	}
	if node.AuthOption != nil {
		node.AuthOption.Format(ctx)
	}
}

func NewUser(username string, hostname string, authOption *AccountIdentified, buf *buffer.Buffer) *User {
	u := buffer.Alloc[User](buf)
	bUsername := NewBufString(username)
	bHostname := NewBufString(hostname)
	buf.Pin(bUsername, bHostname)
	u.Username = bUsername
	u.Hostname = bHostname
	u.AuthOption = authOption
	return u
}

type UsernameRecord struct {
	Username *BufString
	Hostname *BufString
}

func NewUsernameRecord(u, h string, buf *buffer.Buffer) *UsernameRecord {
	us := buffer.Alloc[UsernameRecord](buf)
	bUsername := NewBufString(u)
	bHostname := NewBufString(h)
	buf.Pin(bUsername, bHostname)
	us.Username = bUsername
	us.Hostname = bHostname
	return us
}

type AuthRecord struct {
	AuthPlugin *BufString
	AuthString *BufString
	HashString *BufString
	ByAuth     bool
}

type TlsOption interface {
	NodeFormatter
}

type tlsOptionImpl struct {
	TlsOption
}

type TlsOptionNone struct {
	tlsOptionImpl
}

func (node *TlsOptionNone) Format(ctx *FmtCtx) {
	ctx.WriteString("none")
}

type TlsOptionSSL struct {
	tlsOptionImpl
}

func (node *TlsOptionSSL) Format(ctx *FmtCtx) {
	ctx.WriteString("ssl")
}

type TlsOptionX509 struct {
	tlsOptionImpl
}

func (node *TlsOptionX509) Format(ctx *FmtCtx) {
	ctx.WriteString("x509")
}

type TlsOptionCipher struct {
	tlsOptionImpl
	Cipher *BufString
}

func (node *TlsOptionCipher) Format(ctx *FmtCtx) {
	ctx.WriteString("cipher ")
	ctx.WriteString(node.Cipher.Get())
}

type TlsOptionIssuer struct {
	tlsOptionImpl
	Issuer *BufString
}

func (node *TlsOptionIssuer) Format(ctx *FmtCtx) {
	ctx.WriteString("issuer ")
	ctx.WriteString(node.Issuer.Get())
}

type TlsOptionSubject struct {
	tlsOptionImpl
	Subject *BufString
}

func (node *TlsOptionSubject) Format(ctx *FmtCtx) {
	ctx.WriteString("subject ")
	ctx.WriteString(node.Subject.Get())
}

type TlsOptionSan struct {
	tlsOptionImpl
	San *BufString
}

func (node *TlsOptionSan) Format(ctx *FmtCtx) {
	ctx.WriteString("san ")
	ctx.WriteString(node.San.Get())
}

type ResourceOption interface {
	NodeFormatter
}

type resourceOptionImpl struct {
	ResourceOption
}

type ResourceOptionMaxQueriesPerHour struct {
	resourceOptionImpl
	Count int64
}

func (node *ResourceOptionMaxQueriesPerHour) Format(ctx *FmtCtx) {
	ctx.WriteString("max_queries_per_hour ")
	ctx.WriteString(strconv.FormatInt(node.Count, 10))
}

type ResourceOptionMaxUpdatesPerHour struct {
	resourceOptionImpl
	Count int64
}

func (node *ResourceOptionMaxUpdatesPerHour) Format(ctx *FmtCtx) {
	ctx.WriteString("max_updates_per_hour ")
	ctx.WriteString(strconv.FormatInt(node.Count, 10))
}

type ResourceOptionMaxConnectionPerHour struct {
	resourceOptionImpl
	Count int64
}

func (node *ResourceOptionMaxConnectionPerHour) Format(ctx *FmtCtx) {
	ctx.WriteString("max_connections_per_hour ")
	ctx.WriteString(strconv.FormatInt(node.Count, 10))
}

type ResourceOptionMaxUserConnections struct {
	resourceOptionImpl
	Count int64
}

func (node *ResourceOptionMaxUserConnections) Format(ctx *FmtCtx) {
	ctx.WriteString("max_user_connections ")
	ctx.WriteString(strconv.FormatInt(node.Count, 10))
}

type UserMiscOption interface {
	NodeFormatter
}

type userMiscOptionImpl struct {
	UserMiscOption
}

type UserMiscOptionPasswordExpireNone struct {
	userMiscOptionImpl
}

func NewUserMiscOptionPasswordExpireNone(buf *buffer.Buffer) *UserMiscOptionPasswordExpireNone {
	us := buffer.Alloc[UserMiscOptionPasswordExpireNone](buf)
	return us
}

func (node *UserMiscOptionPasswordExpireNone) Format(ctx *FmtCtx) {
	ctx.WriteString("password expire")
}

type UserMiscOptionPasswordExpireDefault struct {
	userMiscOptionImpl
}

func NewUserMiscOptionPasswordExpireDefault(buf *buffer.Buffer) *UserMiscOptionPasswordExpireDefault {
	us := buffer.Alloc[UserMiscOptionPasswordExpireDefault](buf)
	return us
}

func (node *UserMiscOptionPasswordExpireDefault) Format(ctx *FmtCtx) {
	ctx.WriteString("password expire default")
}

type UserMiscOptionPasswordExpireNever struct {
	userMiscOptionImpl
}

func NewUserMiscOptionPasswordExpireNever(buf *buffer.Buffer) *UserMiscOptionPasswordExpireNever {
	us := buffer.Alloc[UserMiscOptionPasswordExpireNever](buf)
	return us
}

func (node *UserMiscOptionPasswordExpireNever) Format(ctx *FmtCtx) {
	ctx.WriteString("password expire never")
}

type UserMiscOptionPasswordExpireInterval struct {
	userMiscOptionImpl
	Value int64
}

func NewUserMiscOptionPasswordExpireInterval(v int64, buf *buffer.Buffer) *UserMiscOptionPasswordExpireInterval {
	us := buffer.Alloc[UserMiscOptionPasswordExpireInterval](buf)
	return us
}

func (node *UserMiscOptionPasswordExpireInterval) Format(ctx *FmtCtx) {
	ctx.WriteString("password expire interval ")
	ctx.WriteString(strconv.FormatInt(node.Value, 10))
	ctx.WriteString(" day")
}

type UserMiscOptionPasswordHistoryDefault struct {
	userMiscOptionImpl
}

func NewUserMiscOptionPasswordHistoryDefault(buf *buffer.Buffer) *UserMiscOptionPasswordHistoryDefault {
	us := buffer.Alloc[UserMiscOptionPasswordHistoryDefault](buf)
	return us
}

func (node *UserMiscOptionPasswordHistoryDefault) Format(ctx *FmtCtx) {
	ctx.WriteString("password history default")
}

type UserMiscOptionPasswordHistoryCount struct {
	userMiscOptionImpl
	Value int64
}

func NewUserMiscOptionPasswordHistoryCount(v int64, buf *buffer.Buffer) *UserMiscOptionPasswordHistoryCount {
	us := buffer.Alloc[UserMiscOptionPasswordHistoryCount](buf)
	us.Value = v
	return us
}

func (node *UserMiscOptionPasswordHistoryCount) Format(ctx *FmtCtx) {
	ctx.WriteString(fmt.Sprintf("password history %d", node.Value))
}

type UserMiscOptionPasswordReuseIntervalDefault struct {
	userMiscOptionImpl
}

func NewUserMiscOptionPasswordReuseIntervalDefault(buf *buffer.Buffer) *UserMiscOptionPasswordReuseIntervalDefault {
	us := buffer.Alloc[UserMiscOptionPasswordReuseIntervalDefault](buf)
	return us
}

func (node *UserMiscOptionPasswordReuseIntervalDefault) Format(ctx *FmtCtx) {
	ctx.WriteString("password reuse interval default")
}

type UserMiscOptionPasswordReuseIntervalCount struct {
	userMiscOptionImpl
	Value int64
}

func NewUserMiscOptionPasswordReuseIntervalCount(v int64, buf *buffer.Buffer) *UserMiscOptionPasswordReuseIntervalCount {
	us := buffer.Alloc[UserMiscOptionPasswordReuseIntervalCount](buf)
	us.Value = v
	return us
}

func (node *UserMiscOptionPasswordReuseIntervalCount) Format(ctx *FmtCtx) {
	ctx.WriteString(fmt.Sprintf("password reuse interval %d day", node.Value))
}

type UserMiscOptionPasswordRequireCurrentNone struct {
	userMiscOptionImpl
}

func NewUserMiscOptionPasswordRequireCurrentNone(buf *buffer.Buffer) *UserMiscOptionPasswordRequireCurrentNone {
	us := buffer.Alloc[UserMiscOptionPasswordRequireCurrentNone](buf)
	return us
}

func (node *UserMiscOptionPasswordRequireCurrentNone) Format(ctx *FmtCtx) {
	ctx.WriteString("password require current")
}

type UserMiscOptionPasswordRequireCurrentDefault struct {
	userMiscOptionImpl
}

func NewUserMiscOptionPasswordRequireCurrentDefault(buf *buffer.Buffer) *UserMiscOptionPasswordRequireCurrentDefault {
	us := buffer.Alloc[UserMiscOptionPasswordRequireCurrentDefault](buf)
	return us
}

func (node *UserMiscOptionPasswordRequireCurrentDefault) Format(ctx *FmtCtx) {
	ctx.WriteString("password require current default")
}

type UserMiscOptionPasswordRequireCurrentOptional struct {
	userMiscOptionImpl
}

func NewUserMiscOptionPasswordRequireCurrentOptional(buf *buffer.Buffer) *UserMiscOptionPasswordRequireCurrentOptional {
	us := buffer.Alloc[UserMiscOptionPasswordRequireCurrentOptional](buf)
	return us
}

func (node *UserMiscOptionPasswordRequireCurrentOptional) Format(ctx *FmtCtx) {
	ctx.WriteString("password require current optional")
}

type UserMiscOptionFailedLoginAttempts struct {
	userMiscOptionImpl
	Value int64
}

func NewUserMiscOptionFailedLoginAttempts(v int64, buf *buffer.Buffer) *UserMiscOptionFailedLoginAttempts {
	us := buffer.Alloc[UserMiscOptionFailedLoginAttempts](buf)
	us.Value = v
	return us
}

func (node *UserMiscOptionFailedLoginAttempts) Format(ctx *FmtCtx) {
	ctx.WriteString(fmt.Sprintf("failed_login_attempts %d", node.Value))
}

type UserMiscOptionPasswordLockTimeCount struct {
	userMiscOptionImpl
	Value int64
}

func NewUserMiscOptionPasswordLockTimeCount(v int64, buf *buffer.Buffer) *UserMiscOptionPasswordLockTimeCount {
	us := buffer.Alloc[UserMiscOptionPasswordLockTimeCount](buf)
	us.Value = v
	return us
}

func (node *UserMiscOptionPasswordLockTimeCount) Format(ctx *FmtCtx) {
	ctx.WriteString(fmt.Sprintf("password_lock_time %d", node.Value))
}

type UserMiscOptionPasswordLockTimeUnbounded struct {
	userMiscOptionImpl
}

func NewUserMiscOptionPasswordLockTimeUnbounded(buf *buffer.Buffer) *UserMiscOptionPasswordLockTimeUnbounded {
	us := buffer.Alloc[UserMiscOptionPasswordLockTimeUnbounded](buf)
	return us
}

func (node *UserMiscOptionPasswordLockTimeUnbounded) Format(ctx *FmtCtx) {
	ctx.WriteString("password_lock_time unbounded")
}

type UserMiscOptionAccountLock struct {
	userMiscOptionImpl
}

func NewUserMiscOptionAccountLock(buf *buffer.Buffer) *UserMiscOptionAccountLock {
	us := buffer.Alloc[UserMiscOptionAccountLock](buf)
	return us
}

func (node *UserMiscOptionAccountLock) Format(ctx *FmtCtx) {
	ctx.WriteString("lock")
}

type UserMiscOptionAccountUnlock struct {
	userMiscOptionImpl
}

func NewUserMiscOptionAccountUnlock(buf *buffer.Buffer) *UserMiscOptionAccountUnlock {
	us := buffer.Alloc[UserMiscOptionAccountUnlock](buf)
	return us
}

func (node *UserMiscOptionAccountUnlock) Format(ctx *FmtCtx) {
	ctx.WriteString("unlock")
}

type CreateUser struct {
	statementImpl
	IfNotExists bool
	Users       []*User
	Role        *Role
	MiscOpt     UserMiscOption
	// comment or attribute
	CommentOrAttribute *AccountCommentOrAttribute
}

func (node *CreateUser) Format(ctx *FmtCtx) {
	ctx.WriteString("create user ")
	if node.IfNotExists {
		ctx.WriteString("if not exists ")
	}
	if node.Users != nil {
		prefix := ""
		for _, u := range node.Users {
			ctx.WriteString(prefix)
			u.Format(ctx)
			prefix = ", "
		}
	}

	if node.Role != nil {
		ctx.WriteString(" default role")
		ctx.WriteString(" ")
		node.Role.Format(ctx)
	}

	if node.MiscOpt != nil {
		ctx.WriteString(" ")
		node.MiscOpt.Format(ctx)
	}

	node.CommentOrAttribute.Format(ctx)
}

func (node *CreateUser) GetStatementType() string { return "Create User" }
func (node *CreateUser) GetQueryType() string     { return QueryTypeDCL }

func NewCreateUser(ife bool, u []*User, r *Role, misc UserMiscOption, co *AccountCommentOrAttribute, buf *buffer.Buffer) *CreateUser {
	c := buffer.Alloc[CreateUser](buf)
	c.IfNotExists = ife
	c.Users = u
	c.Role = r
	c.MiscOpt = misc
	c.CommentOrAttribute = co
	return c
}

type CreateAccount struct {
	statementImpl
	IfNotExists bool
	Name        *BufString
	AuthOption  *AccountAuthOption
	//status_option or not
	StatusOption *AccountStatus
	//comment or not
	Comment *AccountComment
}

func (ca *CreateAccount) Format(ctx *FmtCtx) {
	ctx.WriteString("create account ")
	if ca.IfNotExists {
		ctx.WriteString("if not exists ")
	}
	ctx.WriteString(ca.Name.Get())
	ca.AuthOption.Format(ctx)
	ca.StatusOption.Format(ctx)
	ca.Comment.Format(ctx)
}

func (ca *CreateAccount) GetStatementType() string { return "Create Account" }
func (ca *CreateAccount) GetQueryType() string     { return QueryTypeDCL }

func NewCreateAccount(ifNotExists bool, name string, authOption *AccountAuthOption, statusOption *AccountStatus, comment *AccountComment, buf *buffer.Buffer) *CreateAccount {
	c := buffer.Alloc[CreateAccount](buf)
	c.IfNotExists = ifNotExists
	bName := NewBufString(name)
	buf.Pin(bName)
	c.Name = bName
	c.AuthOption = authOption
	c.StatusOption = statusOption
	c.Comment = comment
	return c
}

type AccountAuthOption struct {
	Equal          *BufString
	AdminName      *BufString
	IdentifiedType *AccountIdentified
}

func NewAccountAuthOption(equal string, adminName string, identifiedType *AccountIdentified, buf *buffer.Buffer) *AccountAuthOption {
	a := buffer.Alloc[AccountAuthOption](buf)
	bEqual := NewBufString(equal)
	bAdminName := NewBufString(adminName)
	buf.Pin(bEqual, bAdminName)
	a.Equal = bEqual
	a.AdminName = bAdminName
	a.IdentifiedType = identifiedType
	return a
}

func (node *AccountAuthOption) Format(ctx *FmtCtx) {
	ctx.WriteString(" admin_name")
	if len(node.Equal.Get()) != 0 {
		ctx.WriteString(" ")
		ctx.WriteString(node.Equal.Get())
	}

	ctx.WriteString(fmt.Sprintf(" '%s'", node.AdminName.Get()))
	node.IdentifiedType.Format(ctx)

}

type AccountIdentifiedOption int

const (
	AccountIdentifiedByPassword AccountIdentifiedOption = iota
	AccountIdentifiedByRandomPassword
	AccountIdentifiedWithSSL
)

type AccountIdentified struct {
	Typ AccountIdentifiedOption
	Str *BufString
}

func NewAccountIdentified(typ AccountIdentifiedOption, str string, buf *buffer.Buffer) *AccountIdentified {
	a := buffer.Alloc[AccountIdentified](buf)
	a.Typ = typ
	bStr := NewBufString(str)
	buf.Pin(bStr)
	a.Str = bStr
	return a
}

func (node *AccountIdentified) Format(ctx *FmtCtx) {
	switch node.Typ {
	case AccountIdentifiedByPassword:
		ctx.WriteString(" identified by '******'")
	case AccountIdentifiedByRandomPassword:
		ctx.WriteString(" identified by random password")
	case AccountIdentifiedWithSSL:
		ctx.WriteString(" identified with '******'")
	}
}

type AccountStatusOption int

const (
	AccountStatusOpen AccountStatusOption = iota
	AccountStatusSuspend
	AccountStatusRestricted
)

func (aso AccountStatusOption) String() string {
	switch aso {
	case AccountStatusOpen:
		return "open"
	case AccountStatusSuspend:
		return "suspend"
	case AccountStatusRestricted:
		return "restricted"
	default:
		return "open"
	}
}

type AccountStatus struct {
	Exist  bool
	Option AccountStatusOption
}

func NewAccountStatus(e bool, o AccountStatusOption, buf *buffer.Buffer) *AccountStatus {
	a := buffer.Alloc[AccountStatus](buf)
	a.Exist = e
	a.Option = o
	return a
}

func (node *AccountStatus) Format(ctx *FmtCtx) {
	if node.Exist {
		switch node.Option {
		case AccountStatusOpen:
			ctx.WriteString(" open")
		case AccountStatusSuspend:
			ctx.WriteString(" suspend")
		case AccountStatusRestricted:
			ctx.WriteString(" restricted")
		}
	}
}

type AccountComment struct {
	Exist   bool
	Comment *BufString
}

func NewAccountComment(e bool, c string, buf *buffer.Buffer) *AccountComment {
	a := buffer.Alloc[AccountComment](buf)
	a.Exist = e
	bComment := NewBufString(c)
	buf.Pin(bComment)
	a.Comment = bComment
	return a
}

func (node *AccountComment) Format(ctx *FmtCtx) {
	if node.Exist {
		ctx.WriteString(" comment ")
		ctx.WriteString(fmt.Sprintf("'%s'", node.Comment.Get()))
	}
}

type AccountCommentOrAttribute struct {
	Exist     bool
	IsComment bool
	Str       *BufString
}

func NewAccountCommentOrAttribute(e, i bool, s string, buf *buffer.Buffer) *AccountCommentOrAttribute {
	a := buffer.Alloc[AccountCommentOrAttribute](buf)
	a.Exist = e
	a.IsComment = i
	bStr := NewBufString(s)
	buf.Pin(bStr)
	a.Str = bStr
	return a
}

func (node *AccountCommentOrAttribute) Format(ctx *FmtCtx) {
	if node.Exist {
		if node.IsComment {
			ctx.WriteString(" comment ")
		} else {
			ctx.WriteString(" attribute ")
		}
		ctx.WriteString(fmt.Sprintf("'%s'", node.Str.Get()))
	}
}

type CreatePublication struct {
	statementImpl
	IfNotExists bool
	Name        *BufIdentifier
	Database    *BufIdentifier
	AccountsSet *AccountsSetOption
	Comment     *BufString
}

func NewCreatePublication(ifNotExists bool, name Identifier, database Identifier, accountsSet *AccountsSetOption, comment string, buf *buffer.Buffer) *CreatePublication {
	c := buffer.Alloc[CreatePublication](buf)
	n := NewBufIdentifier(name)
	d := NewBufIdentifier(database)
	buf.Pin(n, d)

	c.IfNotExists = ifNotExists
	c.Name = n
	c.Database = d
	c.AccountsSet = accountsSet
	bComment := NewBufString(comment)
	buf.Pin(bComment)
	c.Comment = bComment
	return c
}

func (node *CreatePublication) Format(ctx *FmtCtx) {
	ctx.WriteString("create publication ")
	if node.IfNotExists {
		ctx.WriteString(" if not exists")
	}
	node.Name.Format(ctx)
	ctx.WriteString(" database ")
	node.Database.Format(ctx)
	if node.AccountsSet != nil && len(node.AccountsSet.SetAccounts) > 0 {
		ctx.WriteString(" account ")
		prefix := ""
		for _, a := range node.AccountsSet.SetAccounts {
			ctx.WriteString(prefix)
			a.Format(ctx)
			prefix = ", "
		}
	}
	if len(node.Comment.Get()) > 0 {
		ctx.WriteString(" comment ")
		ctx.WriteString(fmt.Sprintf("'%s'", node.Comment.Get()))
	}
}

type AttributeVisable struct {
	columnAttributeImpl
	Is bool //true NULL (default); false NOT NULL
}

func (node *AttributeVisable) Format(ctx *FmtCtx) {
	if node.Is {
		ctx.WriteString("visible")
	} else {
		ctx.WriteString("not visible")
	}
}

func NewAttributeVisable(is bool, buf *buffer.Buffer) *AttributeVisable {
	a := buffer.Alloc[AttributeVisable](buf)
	a.Is = is
	return a
}

func (node *CreatePublication) GetStatementType() string { return "Create Publication" }
func (node *CreatePublication) GetQueryType() string     { return QueryTypeDCL }
