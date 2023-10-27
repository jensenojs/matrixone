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

package parsers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/matrixorigin/matrixone/pkg/common/buffer"
	"github.com/matrixorigin/matrixone/pkg/sql/parsers/dialect"
	"github.com/matrixorigin/matrixone/pkg/sql/parsers/dialect/mysql"
	"github.com/matrixorigin/matrixone/pkg/sql/parsers/dialect/postgresql"
	"github.com/matrixorigin/matrixone/pkg/sql/parsers/tree"
)

var (
	debugSQL = struct {
		input  string
		output string
	}{
		input: "use db1",
	}
)

func TestMysql(t *testing.T) {
	ctx := context.TODO()
	buf := buffer.New()
	defer buf.Free()

	if debugSQL.output == "" {
		debugSQL.output = debugSQL.input
	}
	ast, err := mysql.ParseOne(ctx, debugSQL.input, 1, buf)
	if err != nil {
		t.Errorf("Parse(%q) err: %v", debugSQL.input, err)
		return
	}
	out := tree.String(ast, dialect.MYSQL)
	if debugSQL.output != out {
		t.Errorf("Parsing failed. \nExpected/Got:\n%s\n%s", debugSQL.output, out)
	}
}

func TestPostgresql(t *testing.T) {
	ctx := context.TODO()
	buf := buffer.New()
	defer buf.Free()
	if debugSQL.output == "" {
		debugSQL.output = debugSQL.input
	}
	ast, err := postgresql.ParseOne(ctx, debugSQL.input, buf)
	if err != nil {
		t.Errorf("Parse(%q) err: %v", debugSQL.input, err)
		return
	}
	out := tree.String(ast, dialect.POSTGRESQL)
	if debugSQL.output != out {
		t.Errorf("Parsing failed. \nExpected/Got:\n%s\n%s", debugSQL.output, out)
	}
}

func TestSplitSqlBySemicolon(t *testing.T) {
	buf := buffer.New()
	defer buf.Free()

	ret := SplitSqlBySemicolon("select 1;select 2;select 3;", buf)
	require.Equal(t, 3, len(ret))
	require.Equal(t, "select 1", ret[0])
	require.Equal(t, "select 2", ret[1])
	require.Equal(t, "select 3", ret[2])

	ret = SplitSqlBySemicolon("select 1;select 2/*;;;*/;select 3;", buf)
	require.Equal(t, 3, len(ret))
	require.Equal(t, "select 1", ret[0])
	require.Equal(t, "select 2/*;;;*/", ret[1])
	require.Equal(t, "select 3", ret[2])

	ret = SplitSqlBySemicolon("select 1;select \"2;;\";select 3;", buf)
	require.Equal(t, 3, len(ret))
	require.Equal(t, "select 1", ret[0])
	require.Equal(t, "select \"2;;\"", ret[1])
	require.Equal(t, "select 3", ret[2])

	ret = SplitSqlBySemicolon("select 1;select '2;;';select 3;", buf)
	require.Equal(t, 3, len(ret))
	require.Equal(t, "select 1", ret[0])
	require.Equal(t, "select '2;;'", ret[1])
	require.Equal(t, "select 3", ret[2])

	ret = SplitSqlBySemicolon("select 1;select '2;;';select 3", buf)
	require.Equal(t, 3, len(ret))
	require.Equal(t, "select 1", ret[0])
	require.Equal(t, "select '2;;'", ret[1])
	require.Equal(t, "select 3", ret[2])

	ret = SplitSqlBySemicolon("select 1", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "select 1", ret[0])

	ret = SplitSqlBySemicolon(";;;", buf)
	require.Equal(t, 3, len(ret))
	require.Equal(t, "", ret[0])
	require.Equal(t, "", ret[1])
	require.Equal(t, "", ret[2])

	ret = SplitSqlBySemicolon(";;;  ", buf)
	require.Equal(t, 3, len(ret))
	require.Equal(t, "", ret[0])
	require.Equal(t, "", ret[1])
	require.Equal(t, "", ret[2])

	ret = SplitSqlBySemicolon(";", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "", ret[0])

	ret = SplitSqlBySemicolon("", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "", ret[0])

	ret = SplitSqlBySemicolon("   ;   ", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "", ret[0])

	ret = SplitSqlBySemicolon("   ", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "", ret[0])

	ret = SplitSqlBySemicolon("  ; /* abc */ ", buf)
	require.Equal(t, 2, len(ret))
	require.Equal(t, "", ret[0])
	require.Equal(t, "/* abc */", ret[1])

	ret = SplitSqlBySemicolon(" /* cde */  ; /* abc */ ", buf)
	require.Equal(t, 2, len(ret))
	require.Equal(t, "/* cde */", ret[0])
	require.Equal(t, "/* abc */", ret[1])

	ret = SplitSqlBySemicolon("   ;    ;  ", buf)
	require.Equal(t, 2, len(ret))
	require.Equal(t, "", ret[0])
	require.Equal(t, "", ret[1])

	ret = SplitSqlBySemicolon("   ;    ;", buf)
	require.Equal(t, 2, len(ret))
	require.Equal(t, "", ret[0])
	require.Equal(t, "", ret[1])

	ret = SplitSqlBySemicolon("   ;   ", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "", ret[0])
}

func TestHandleSqlForRecord(t *testing.T) {
	// Test remove /* cloud_user */ prefix
	var ret []string
	buf := buffer.New()
	defer buf.Free()

	ret = HandleSqlForRecord(" ;   ;  ", buf)
	require.Equal(t, 2, len(ret))
	require.Equal(t, "", ret[0])
	require.Equal(t, "", ret[1])

	ret = HandleSqlForRecord(" ; /* abc */  ", buf)
	require.Equal(t, 2, len(ret))
	require.Equal(t, "", ret[0])
	require.Equal(t, "/* abc */", ret[1])

	ret = HandleSqlForRecord(" /* cde */  ; /* abc */ ", buf)
	require.Equal(t, 2, len(ret))
	require.Equal(t, "/* cde */", ret[0])
	require.Equal(t, "/* abc */", ret[1])

	ret = HandleSqlForRecord(" /* cde */  ; /* abc */ ; "+stripCloudNonUser+" ; "+stripCloudUser, buf)
	require.Equal(t, 4, len(ret))
	require.Equal(t, "/* cde */", ret[0])
	require.Equal(t, "/* abc */", ret[1])
	require.Equal(t, "", ret[2])
	require.Equal(t, "", ret[3])

	ret = HandleSqlForRecord("  /* cloud_user */ select 1;   ", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "select 1", ret[0])

	ret = HandleSqlForRecord("  /* cloud_user */ select 1;  ", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "select 1", ret[0])

	ret = HandleSqlForRecord("  /* cloud_user */select * from t;/* cloud_user */select * from t;/* cloud_user */select * from t;", buf)
	require.Equal(t, 3, len(ret))
	require.Equal(t, "select * from t", ret[0])
	require.Equal(t, "select * from t", ret[1])
	require.Equal(t, "select * from t", ret[2])

	ret = HandleSqlForRecord("  /* cloud_user */  select * from t ;  /* cloud_user */  select * from t ; /* cloud_user */ select * from t ; ", buf)
	require.Equal(t, 3, len(ret))
	require.Equal(t, "select * from t", ret[0])
	require.Equal(t, "select * from t", ret[1])
	require.Equal(t, "select * from t", ret[2])

	ret = HandleSqlForRecord("  /* cloud_user */  select * from t ;  /* cloud_user */  select * from t ; /* cloud_user */ select * from t ; /* abc */ ", buf)
	require.Equal(t, 4, len(ret))
	require.Equal(t, "select * from t", ret[0])
	require.Equal(t, "select * from t", ret[1])
	require.Equal(t, "select * from t", ret[2])
	require.Equal(t, "/* abc */", ret[3])

	ret = HandleSqlForRecord("  /* cloud_user */  ", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "", ret[0])

	ret = HandleSqlForRecord("  /* cloud_user */   ", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "", ret[0])

	ret = HandleSqlForRecord("   "+stripCloudNonUser+"  select 1;   ", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "select 1", ret[0])

	ret = HandleSqlForRecord("  "+stripCloudNonUser+"  select * from t  ;  "+stripCloudNonUser+"   select * from t  ;   "+stripCloudNonUser+"   select * from t  ;   ", buf)
	require.Equal(t, 3, len(ret))
	require.Equal(t, "select * from t", ret[0])
	require.Equal(t, "select * from t", ret[1])
	require.Equal(t, "select * from t", ret[2])

	ret = HandleSqlForRecord("  "+stripCloudNonUser+"  select * from t  ;  "+stripCloudNonUser+"   select * from t  ;   "+stripCloudNonUser+"   select * from t  ; /* abc */  ", buf)
	require.Equal(t, 4, len(ret))
	require.Equal(t, "select * from t", ret[0])
	require.Equal(t, "select * from t", ret[1])
	require.Equal(t, "select * from t", ret[2])
	require.Equal(t, "/* abc */", ret[3])

	ret = HandleSqlForRecord("   "+stripCloudNonUser+"  ", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "", ret[0])

	ret = HandleSqlForRecord("   "+stripCloudUser+"  ", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "", ret[0])

	ret = HandleSqlForRecord("", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "", ret[0])

	// Test hide secret key

	ret = HandleSqlForRecord("create user u identified by '123456';", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "create user u identified by '******'", ret[0])

	ret = HandleSqlForRecord("create user u identified with '12345';", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "create user u identified with '******'", ret[0])

	ret = HandleSqlForRecord("create user u identified by random password;", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "create user u identified by random password", ret[0])

	ret = HandleSqlForRecord("create user if not exists abc1 identified by '123', abc2 identified by '234', abc3 identified with '111', abc3 identified by random password;", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "create user if not exists abc1 identified by '******', abc2 identified by '******', abc3 identified with '******', abc3 identified by random password", ret[0])

	ret = HandleSqlForRecord("create external table t (a int) URL s3option{'endpoint'='s3.us-west-2.amazonaws.com', 'access_key_id'='123', 'secret_access_key'='123', 'bucket'='test', 'filepath'='*.txt', 'region'='us-west-2'};", buf)
	require.Equal(t, 1, len(ret))
	require.Equal(t, "create external table t (a int) URL s3option{'endpoint'='s3.us-west-2.amazonaws.com', 'access_key_id'='******', 'secret_access_key'='******', 'bucket'='test', 'filepath'='*.txt', 'region'='us-west-2'}", ret[0])
}
