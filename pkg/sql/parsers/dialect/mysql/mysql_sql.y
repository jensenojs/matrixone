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

%{
package mysql

import (
    "fmt"
    "strings"
    "go/constant"

    "github.com/matrixorigin/matrixone/pkg/sql/parsers/tree"
    "github.com/matrixorigin/matrixone/pkg/sql/parsers/util"
    "github.com/matrixorigin/matrixone/pkg/common/buffer"
    "github.com/matrixorigin/matrixone/pkg/defines"
)
%}

%struct {
    id  int
    str string
    item interface{}
}

%union {
    statement tree.Statement
    statements []tree.Statement

    alterTable tree.AlterTable
    alterTableOptions tree.AlterTableOptions
    alterTableOption tree.AlterTableOption
    alterColPosition *tree.ColumnPosition
    alterColumnOrderBy []*tree.AlterColumnOrder
    alterColumnOrder *tree.AlterColumnOrder

    tableDef tree.TableDef
    tableDefs tree.TableDefs
    tableName *tree.TableName
    tableNames tree.TableNames
    columnTableDef *tree.ColumnTableDef
    tableOption tree.TableOption
    tableOptions []tree.TableOption
    tableExprs tree.TableExprs
    tableExpr tree.TableExpr
    rowFormatType tree.RowFormatType
    matchType tree.MatchType
    attributeReference *tree.AttributeReference
    loadParam *tree.ExternParam
    tailParam *tree.TailParameter
    connectorOption *tree.ConnectorOption
    connectorOptions []*tree.ConnectorOption

    functionName *tree.FunctionName
    funcArg tree.FunctionArg
    funcArgs tree.FunctionArgs
    funcArgDecl *tree.FunctionArgDecl
    funcReturn *tree.ReturnType

    procName *tree.ProcedureName
    procArg tree.ProcedureArg
    procArgs tree.ProcedureArgs
    procArgDecl *tree.ProcedureArgDecl
    procArgType tree.InOutArgType

    from *tree.From
    where *tree.Where
    groupBy tree.GroupBy
    aliasedTableExpr *tree.AliasedTableExpr
    direction tree.Direction
    nullsPosition tree.NullsPosition
    orderBy tree.OrderBy
    order *tree.Order
    limit *tree.Limit
    unionTypeRecord *tree.UnionTypeRecord
    parenTableExpr *tree.ParenTableExpr
    identifierList tree.IdentifierList
    joinCond tree.JoinCond
    selectLockInfo *tree.SelectLockInfo

    columnType *tree.T
    unresolvedName *tree.UnresolvedName
    lengthScaleOpt tree.LengthScaleOpt
    tuple *tree.Tuple
    funcType tree.FuncType

    columnAttribute tree.ColumnAttribute
    columnAttributes []tree.ColumnAttribute
    attributeNull tree.AttributeNull
    expr tree.Expr
    exprs tree.Exprs
    rowsExprs []tree.Exprs
    comparisonOp tree.ComparisonOp
    referenceOptionType tree.ReferenceOptionType
    referenceOnRecord *tree.ReferenceOnRecord

    select *tree.Select
    selectStatement tree.SelectStatement
    selectExprs tree.SelectExprs
    selectExpr tree.SelectExpr

    insert *tree.Insert
    replace *tree.Replace
    createOption tree.CreateOption
    createOptions []tree.CreateOption
    indexType tree.IndexType
    indexCategory tree.IndexCategory
    keyParts []*tree.KeyPart
    keyPart *tree.KeyPart
    indexOption *tree.IndexOption
    comparisionExpr *tree.ComparisonExpr

    userMiscOption tree.UserMiscOption
    userMiscOptions []tree.UserMiscOption
    updateExpr *tree.UpdateExpr
    updateExprs tree.UpdateExprs
    completionType tree.CompletionType
    varAssignmentExpr *tree.VarAssignmentExpr
    varAssignmentExprs []*tree.VarAssignmentExpr
    setRole *tree.SetRole
    setDefaultRole *tree.SetDefaultRole
    privilege *tree.Privilege
    privileges []*tree.Privilege
    objectType tree.ObjectType
    privilegeType tree.PrivilegeType
    privilegeLevel *tree.PrivilegeLevel
    unresolveNames []*tree.UnresolvedName

    partitionOption *tree.PartitionOption
    clusterByOption *tree.ClusterByOption
    partitionBy *tree.PartitionBy
    windowSpec *tree.WindowSpec
    frameClause *tree.FrameClause
    frameBound *tree.FrameBound
    frameType tree.FrameType
    partition *tree.Partition
    partitions []*tree.Partition
    values tree.Values
    numVal *tree.NumVal
    subPartition *tree.SubPartition
    subPartitions []*tree.SubPartition

    subquery *tree.Subquery
    funcExpr *tree.FuncExpr

    roles []*tree.Role
    role *tree.Role
    usernameRecord *tree.UsernameRecord
    authRecord *tree.AuthRecord
    user *tree.User
    users []*tree.User
    tlsOptions []tree.TlsOption
    tlsOption tree.TlsOption
    resourceOption tree.ResourceOption
    resourceOptions []tree.ResourceOption
    unresolvedObjectName *tree.UnresolvedObjectName
    lengthOpt int32
    unsignedOpt bool
    zeroFillOpt bool
    ifNotExists bool
    defaultOptional bool
    streamOptional bool
    connectorOptional bool
    fullOpt bool
    boolVal bool
    int64Val int64
    strs []string

    duplicateKey tree.DuplicateKey
    fields *tree.Fields
    fieldsList []*tree.Fields
    lines *tree.Lines
    varExpr *tree.VarExpr
    varExprs []*tree.VarExpr
    loadColumn tree.LoadColumn
    loadColumns []tree.LoadColumn
    assignments []*tree.Assignment
    assignment *tree.Assignment
    properties []tree.Property
    property tree.Property
    exportParm *tree.ExportParam

    epxlainOptions []tree.OptionElem
    epxlainOption tree.OptionElem
    whenClause *tree.When
    whenClauseList []*tree.When
    withClause *tree.With
    cte *tree.CTE
    cteList []*tree.CTE

    accountAuthOption tree.AccountAuthOption
    alterAccountAuthOption tree.AlterAccountAuthOption
    accountIdentified tree.AccountIdentified
    accountStatus tree.AccountStatus
    accountComment tree.AccountComment
    stageComment tree.StageComment
    stageStatus tree.StageStatus
    stageUrl tree.StageUrl
    stageCredentials tree.StageCredentials
    accountCommentOrAttribute tree.AccountCommentOrAttribute
    userIdentified *tree.AccountIdentified
    accountRole *tree.Role
    showType tree.ShowType
    joinTableExpr *tree.JoinTableExpr

    indexHintType tree.IndexHintType
    indexHintScope tree.IndexHintScope
    indexHint *tree.IndexHint
    indexHintList []*tree.IndexHint
    indexVisibility tree.VisibleType

    killOption tree.KillOption
    statementOption tree.StatementOption

    tableLock tree.TableLock
    tableLocks []tree.TableLock
    tableLockType tree.TableLockType
    cstr *tree.CStr
    incrementByOption *tree.IncrementByOption
    minValueOption  *tree.MinValueOption
    maxValueOption  *tree.MaxValueOption
    startWithOption *tree.StartWithOption
    cycleOption     *tree.CycleOption
    alterTypeOption *tree.TypeOption

    whenClause2 *tree.WhenStmt
    whenClauseList2 []*tree.WhenStmt

    elseIfClause *tree.ElseIfStmt
    elseIfClauseList []*tree.ElseIfStmt
    subscriptionOption *tree.SubscriptionOption
    accountsSetOption *tree.AccountsSetOption

    transactionCharacteristic *tree.TransactionCharacteristic
    transactionCharacteristicList []*tree.TransactionCharacteristic
    isolationLevel tree.IsolationLevelType
    accessMode tree.AccessModeType
}

%token LEX_ERROR
%nonassoc EMPTY
%left <str> UNION EXCEPT INTERSECT MINUS
%nonassoc LOWER_THAN_ORDER
%nonassoc ORDER
%token <str> SELECT INSERT UPDATE DELETE FROM WHERE GROUP HAVING BY LIMIT OFFSET FOR CONNECT MANAGE GRANTS OWNERSHIP REFERENCE
%nonassoc LOWER_THAN_SET
%nonassoc <str> SET
%token <str> ALL DISTINCT DISTINCTROW AS EXISTS ASC DESC INTO DUPLICATE DEFAULT LOCK KEYS NULLS FIRST LAST AFTER
%token <str> INSTANT INPLACE COPY DISABLE ENABLE UNDEFINED MERGE TEMPTABLE DEFINER INVOKER SQL SECURITY CASCADED
%token <str> VALUES
%token <str> NEXT VALUE SHARE MODE
%token <str> SQL_NO_CACHE SQL_CACHE
%left <str> JOIN STRAIGHT_JOIN LEFT RIGHT INNER OUTER CROSS NATURAL USE FORCE
%nonassoc LOWER_THAN_ON
%nonassoc <str> ON USING
%left <str> SUBQUERY_AS_EXPR
%right <str> '('
%left <str> ')'
%nonassoc LOWER_THAN_STRING
%nonassoc <str> ID AT_ID AT_AT_ID STRING VALUE_ARG LIST_ARG COMMENT COMMENT_KEYWORD QUOTE_ID STAGE CREDENTIALS STAGES
%token <item> INTEGRAL HEX BIT_LITERAL FLOAT
%token <str>  HEXNUM
%token <str> NULL TRUE FALSE
%nonassoc LOWER_THAN_CHARSET
%nonassoc <str> CHARSET
%right <str> UNIQUE KEY
%left <str> OR PIPE_CONCAT
%left <str> XOR
%left <str> AND
%right <str> NOT '!'
%left <str> BETWEEN CASE WHEN THEN ELSE END ELSEIF
%nonassoc LOWER_THAN_EQ
%left <str> '=' '<' '>' LE GE NE NULL_SAFE_EQUAL IS LIKE REGEXP IN ASSIGNMENT ILIKE
%left <str> '|'
%left <str> '&'
%left <str> SHIFT_LEFT SHIFT_RIGHT
%left <str> '+' '-'
%left <str> '*' '/' DIV '%' MOD
%left <str> '^'
%right <str> '~' UNARY
%left <str> COLLATE
%right <str> BINARY UNDERSCORE_BINARY
%right <str> INTERVAL
%nonassoc <str> '.' ','

%token <str> OUT INOUT

// Transaction
%token <str> BEGIN START TRANSACTION COMMIT ROLLBACK WORK CONSISTENT SNAPSHOT
%token <str> CHAIN NO RELEASE PRIORITY QUICK

// Type
%token <str> BIT TINYINT SMALLINT MEDIUMINT INT INTEGER BIGINT INTNUM
%token <str> REAL DOUBLE FLOAT_TYPE DECIMAL NUMERIC DECIMAL_VALUE
%token <str> TIME TIMESTAMP DATETIME YEAR
%token <str> CHAR VARCHAR BOOL CHARACTER VARBINARY NCHAR
%token <str> TEXT TINYTEXT MEDIUMTEXT LONGTEXT
%token <str> BLOB TINYBLOB MEDIUMBLOB LONGBLOB JSON ENUM UUID VECF32 VECF64
%token <str> GEOMETRY POINT LINESTRING POLYGON GEOMETRYCOLLECTION MULTIPOINT MULTILINESTRING MULTIPOLYGON
%token <str> INT1 INT2 INT3 INT4 INT8 S3OPTION

// Select option
%token <str> SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT
%token <str> LOW_PRIORITY HIGH_PRIORITY DELAYED

// Create Table
%token <str> CREATE ALTER DROP RENAME ANALYZE ADD RETURNS
%token <str> SCHEMA TABLE SEQUENCE INDEX VIEW TO IGNORE IF PRIMARY COLUMN CONSTRAINT SPATIAL FULLTEXT FOREIGN KEY_BLOCK_SIZE
%token <str> SHOW DESCRIBE EXPLAIN DATE ESCAPE REPAIR OPTIMIZE TRUNCATE
%token <str> MAXVALUE PARTITION REORGANIZE LESS THAN PROCEDURE TRIGGER
%token <str> STATUS VARIABLES ROLE PROXY AVG_ROW_LENGTH STORAGE DISK MEMORY
%token <str> CHECKSUM COMPRESSION DATA DIRECTORY DELAY_KEY_WRITE ENCRYPTION ENGINE
%token <str> MAX_ROWS MIN_ROWS PACK_KEYS ROW_FORMAT STATS_AUTO_RECALC STATS_PERSISTENT STATS_SAMPLE_PAGES
%token <str> DYNAMIC COMPRESSED REDUNDANT COMPACT FIXED COLUMN_FORMAT AUTO_RANDOM ENGINE_ATTRIBUTE SECONDARY_ENGINE_ATTRIBUTE INSERT_METHOD
%token <str> RESTRICT CASCADE ACTION PARTIAL SIMPLE CHECK ENFORCED
%token <str> RANGE LIST ALGORITHM LINEAR PARTITIONS SUBPARTITION SUBPARTITIONS CLUSTER
%token <str> TYPE ANY SOME EXTERNAL LOCALFILE URL
%token <str> PREPARE DEALLOCATE RESET
%token <str> EXTENSION

// Sequence
%token <str> INCREMENT CYCLE MINVALUE
// publication
%token <str> PUBLICATION SUBSCRIPTIONS PUBLICATIONS


// MO table option
%token <str> PROPERTIES

// Index
%token <str> PARSER VISIBLE INVISIBLE BTREE HASH RTREE BSI
%token <str> ZONEMAP LEADING BOTH TRAILING UNKNOWN

// Alter
%token <str> EXPIRE ACCOUNT ACCOUNTS UNLOCK DAY NEVER PUMP MYSQL_COMPATIBILITY_MODE
%token <str> MODIFY CHANGE

// Time
%token <str> SECOND ASCII COALESCE COLLATION HOUR MICROSECOND MINUTE MONTH QUARTER REPEAT
%token <str> REVERSE ROW_COUNT WEEK

// Revoke
%token <str> REVOKE FUNCTION PRIVILEGES TABLESPACE EXECUTE SUPER GRANT OPTION REFERENCES REPLICATION
%token <str> SLAVE CLIENT USAGE RELOAD FILE TEMPORARY ROUTINE EVENT SHUTDOWN

// Type Modifiers
%token <str> NULLX AUTO_INCREMENT APPROXNUM SIGNED UNSIGNED ZEROFILL ENGINES LOW_CARDINALITY AUTOEXTEND_SIZE

// Account
%token <str> ADMIN_NAME RANDOM SUSPEND ATTRIBUTE HISTORY REUSE CURRENT OPTIONAL FAILED_LOGIN_ATTEMPTS PASSWORD_LOCK_TIME UNBOUNDED SECONDARY RESTRICTED

// User
%token <str> USER IDENTIFIED CIPHER ISSUER X509 SUBJECT SAN REQUIRE SSL NONE PASSWORD SHARED EXCLUSIVE
%token <str> MAX_QUERIES_PER_HOUR MAX_UPDATES_PER_HOUR MAX_CONNECTIONS_PER_HOUR MAX_USER_CONNECTIONS

// Explain
%token <str> FORMAT VERBOSE CONNECTION TRIGGERS PROFILES

// Load
%token <str> LOAD INLINE INFILE TERMINATED OPTIONALLY ENCLOSED ESCAPED STARTING LINES ROWS IMPORT DISCARD

// MODump
%token <str> MODUMP

// Window function
%token <str> OVER PRECEDING FOLLOWING GROUPS

// Supported SHOW tokens
%token <str> DATABASES TABLES SEQUENCES EXTENDED FULL PROCESSLIST FIELDS COLUMNS OPEN ERRORS WARNINGS INDEXES SCHEMAS NODE LOCKS ROLES
%token <str> TABLE_NUMBER COLUMN_NUMBER TABLE_VALUES TABLE_SIZE

// SET tokens
%token <str> NAMES GLOBAL PERSIST SESSION ISOLATION LEVEL READ WRITE ONLY REPEATABLE COMMITTED UNCOMMITTED SERIALIZABLE
%token <str> LOCAL EVENTS PLUGINS

// Functions
%token <str> CURRENT_TIMESTAMP DATABASE
%token <str> CURRENT_TIME LOCALTIME LOCALTIMESTAMP
%token <str> UTC_DATE UTC_TIME UTC_TIMESTAMP
%token <str> REPLACE CONVERT
%token <str> SEPARATOR TIMESTAMPDIFF
%token <str> CURRENT_DATE CURRENT_USER CURRENT_ROLE

// Time unit
%token <str> SECOND_MICROSECOND MINUTE_MICROSECOND MINUTE_SECOND HOUR_MICROSECOND
%token <str> HOUR_SECOND HOUR_MINUTE DAY_MICROSECOND DAY_SECOND DAY_MINUTE DAY_HOUR YEAR_MONTH
%token <str> SQL_TSI_HOUR SQL_TSI_DAY SQL_TSI_WEEK SQL_TSI_MONTH SQL_TSI_QUARTER SQL_TSI_YEAR
%token <str> SQL_TSI_SECOND SQL_TSI_MINUTE

// With
%token <str> RECURSIVE CONFIG DRAINER

// Stream
%token <str> SOURCE STREAM HEADERS CONNECTOR CONNECTORS DAEMON PAUSE CANCEL TASK RESUME

// Match
%token <str> MATCH AGAINST BOOLEAN LANGUAGE WITH QUERY EXPANSION WITHOUT VALIDATION

// Built-in function
%token <str> ADDDATE BIT_AND BIT_OR BIT_XOR CAST COUNT APPROX_COUNT APPROX_COUNT_DISTINCT
%token <str> APPROX_PERCENTILE CURDATE CURTIME DATE_ADD DATE_SUB EXTRACT
%token <str> GROUP_CONCAT MAX MID MIN NOW POSITION SESSION_USER STD STDDEV MEDIAN
%token <str> STDDEV_POP STDDEV_SAMP SUBDATE SUBSTR SUBSTRING SUM SYSDATE
%token <str> SYSTEM_USER TRANSLATE TRIM VARIANCE VAR_POP VAR_SAMP AVG RANK ROW_NUMBER
%token <str> DENSE_RANK BIT_CAST

// Sequence function
%token <str> NEXTVAL SETVAL CURRVAL LASTVAL

//JSON function
%token <str> ARROW

// Insert
%token <str> ROW OUTFILE HEADER MAX_FILE_SIZE FORCE_QUOTE PARALLEL

%token <str> UNUSED BINDINGS

// Do
%token <str> DO

// Declare
%token <str> DECLARE

// Iteration
%token <str> LOOP WHILE LEAVE ITERATE UNTIL

// Call
%token <str> CALL

// sp_begin_sym
%token <str> SPBEGIN

%token <str> BACKEND SERVERS

%type <statement> stmt block_stmt block_type_stmt normal_stmt
%type <statements> stmt_list stmt_list_return
%type <statement> create_stmt insert_stmt delete_stmt drop_stmt alter_stmt truncate_table_stmt alter_sequence_stmt
%type <statement> delete_without_using_stmt delete_with_using_stmt
%type <statement> drop_ddl_stmt drop_database_stmt drop_table_stmt drop_index_stmt drop_prepare_stmt drop_view_stmt drop_connector_stmt drop_function_stmt drop_procedure_stmt drop_sequence_stmt
%type <statement> drop_account_stmt drop_role_stmt drop_user_stmt
%type <statement> create_account_stmt create_user_stmt create_role_stmt
%type <statement> create_ddl_stmt create_table_stmt create_database_stmt create_index_stmt create_view_stmt create_function_stmt create_extension_stmt create_procedure_stmt create_sequence_stmt
%type <statement> create_stream_stmt create_connector_stmt pause_daemon_task_stmt cancel_daemon_task_stmt resume_daemon_task_stmt
%type <statement> show_stmt show_create_stmt show_columns_stmt show_databases_stmt show_target_filter_stmt show_table_status_stmt show_grants_stmt show_collation_stmt show_accounts_stmt show_roles_stmt show_stages_stmt
%type <statement> show_tables_stmt show_sequences_stmt show_process_stmt show_errors_stmt show_warnings_stmt show_target
%type <statement> show_procedure_status_stmt show_function_status_stmt show_node_list_stmt show_locks_stmt
%type <statement> show_table_num_stmt show_column_num_stmt show_table_values_stmt show_table_size_stmt
%type <statement> show_variables_stmt show_status_stmt show_index_stmt
%type <statement> show_servers_stmt show_connectors_stmt
%type <statement> alter_account_stmt alter_user_stmt alter_view_stmt update_stmt use_stmt update_no_with_stmt alter_database_config_stmt alter_table_stmt
%type <statement> transaction_stmt begin_stmt commit_stmt rollback_stmt
%type <statement> explain_stmt explainable_stmt
%type <statement> set_stmt set_variable_stmt set_password_stmt set_role_stmt set_default_role_stmt set_transaction_stmt
%type <statement> lock_stmt lock_table_stmt unlock_table_stmt
%type <statement> revoke_stmt grant_stmt
%type <statement> load_data_stmt
%type <statement> analyze_stmt
%type <statement> prepare_stmt prepareable_stmt deallocate_stmt execute_stmt reset_stmt
%type <statement> replace_stmt
%type <statement> do_stmt
%type <statement> declare_stmt
%type <statement> values_stmt
%type <statement> call_stmt
%type <statement> mo_dump_stmt
%type <statement> load_extension_stmt
%type <statement> kill_stmt
%type <statement> backup_stmt
%type <rowsExprs> row_constructor_list
%type <exprs>  row_constructor
%type <exportParm> export_data_param_opt
%type <loadParam> load_param_opt load_param_opt_2
%type <tailParam> tail_param_opt

// case statement
%type <statement> case_stmt
%type <whenClause2> when_clause2
%type <whenClauseList2> when_clause_list2
%type <statements> else_clause_opt2

// if statement
%type <statement> if_stmt
%type <elseIfClause> elseif_clause
%type <elseIfClauseList> elseif_clause_list elseif_clause_list_opt

// iteration
%type <statement> loop_stmt iterate_stmt leave_stmt repeat_stmt while_stmt
%type <statement> create_publication_stmt drop_publication_stmt alter_publication_stmt show_publications_stmt show_subscriptions_stmt
%type <statement> create_stage_stmt drop_stage_stmt alter_stage_stmt
%type <str> urlparams
%type <str> comment_opt view_list_opt view_opt security_opt view_tail check_type
%type <subscriptionOption> subcription_opt
%type <accountsSetOption> alter_publication_accounts_opt

%type <select> select_stmt select_no_parens
%type <selectStatement> simple_select select_with_parens simple_select_clause
%type <selectExprs> select_expression_list
%type <selectExpr> select_expression
%type <tableExprs> table_name_wild_list
%type <joinTableExpr>  table_references join_table
%type <tableExpr> into_table_name table_function table_factor table_reference escaped_table_reference
%type <direction> asc_desc_opt
%type <nullsPosition> nulls_first_last_opt
%type <order> order
%type <orderBy> order_list order_by_clause order_by_opt
%type <limit> limit_opt limit_clause
%type <str> insert_column
%type <identifierList> column_list column_list_opt partition_clause_opt partition_id_list insert_column_list accounts_list
%type <joinCond> join_condition join_condition_opt on_expression_opt
%type <selectLockInfo> select_lock_opt

%type <functionName> func_name
%type <funcArgs> func_args_list_opt func_args_list
%type <funcArg> func_arg
%type <funcArgDecl> func_arg_decl
%type <funcReturn> func_return
%type <str> func_lang extension_lang extension_name

%type <procName> proc_name
%type <procArgs> proc_args_list_opt proc_args_list
%type <procArg> proc_arg
%type <procArgDecl> proc_arg_decl
%type <procArgType> proc_arg_in_out_type

%type <tableDefs> table_elem_list_opt table_elem_list
%type <tableDef> table_elem constaint_def constraint_elem index_def table_elem_2
%type <tableName> table_name table_name_opt_wild
%type <tableNames> table_name_list
%type <columnTableDef> column_def
%type <columnType> mo_cast_type mysql_cast_type
%type <columnType> column_type char_type spatial_type time_type numeric_type decimal_type int_type as_datatype_opt
%type <str> integer_opt
%type <columnAttribute> column_attribute_elem keys
%type <columnAttributes> column_attribute_list column_attribute_list_opt
%type <tableOptions> table_option_list_opt table_option_list stream_option_list_opt stream_option_list
%type <str> charset_name storage_opt collate_name column_format storage_media algorithm_type able_type space_type lock_type with_type rename_type algorithm_type_2
%type <rowFormatType> row_format_options
%type <int64Val> field_length_opt max_file_size_opt
%type <matchType> match match_opt
%type <referenceOptionType> ref_opt on_delete on_update
%type <referenceOnRecord> on_delete_update_opt
%type <attributeReference> references_def
%type <alterTableOptions> alter_option_list
%type <alterTableOption> alter_option alter_table_drop alter_table_alter alter_table_rename
%type <alterColPosition> column_position
%type <alterColumnOrder> alter_column_order
%type <alterColumnOrderBy> alter_column_order_list
%type <indexVisibility> visibility

%type <tableOption> table_option stream_option
%type <connectorOption> connector_option
%type <connectorOptions> connector_option_list
%type <from> from_clause from_opt
%type <where> where_expression_opt having_opt
%type <groupBy> group_by_opt
%type <aliasedTableExpr> aliased_table_name
%type <unionTypeRecord> union_op
%type <parenTableExpr> table_subquery
%type <str> inner_join straight_join outer_join natural_join
%type <funcType> func_type_opt
%type <funcExpr> function_call_generic
%type <funcExpr> function_call_keyword
%type <funcExpr> function_call_nonkeyword
%type <funcExpr> function_call_aggregate
%type <funcExpr> function_call_window

%type <unresolvedName> column_name column_name_unresolved
%type <strs> enum_values force_quote_opt force_quote_list infile_or_s3_param infile_or_s3_params credentialsparams credentialsparam
%type <str> charset_keyword db_name db_name_opt
%type <str> not_keyword func_not_keyword
%type <str> non_reserved_keyword
%type <str> equal_opt column_keyword_opt
%type <str> as_opt_id name_string
%type <cstr> ident as_name_opt
%type <str> table_alias explain_sym prepare_sym deallocate_sym stmt_name reset_sym
%type <unresolvedObjectName> unresolved_object_name table_column_name
%type <unresolvedObjectName> table_name_unresolved
%type <comparisionExpr> like_opt
%type <fullOpt> full_opt
%type <str> database_name_opt auth_string constraint_keyword_opt constraint_keyword
%type <userMiscOption> pwd_or_lck pwd_or_lck_opt
//%type <userMiscOptions> pwd_or_lck_list

%type <expr> literal num_literal
%type <expr> predicate
%type <expr> bit_expr interval_expr
%type <expr> simple_expr else_opt
%type <expr> expression like_escape_opt boolean_primary col_tuple expression_opt
%type <exprs> expression_list_opt
%type <exprs> expression_list row_value window_partition_by window_partition_by_opt
%type <expr> datetime_scale_opt datetime_scale
%type <tuple> tuple_expression
%type <comparisonOp> comparison_operator and_or_some
%type <createOption> create_option
%type <createOptions> create_option_list_opt create_option_list
%type <ifNotExists> not_exists_opt
%type <defaultOptional> default_opt
%type <streamOptional> replace_opt
%type <str> database_or_schema
%type <indexType> using_opt
%type <indexCategory> index_prefix
%type <keyParts> index_column_list index_column_list_opt
%type <keyPart> index_column
%type <indexOption> index_option_list index_option
%type <roles> role_spec_list using_roles_opt
%type <role> role_spec
%type <cstr> role_name
%type <usernameRecord> user_name
%type <user> user_spec drop_user_spec user_spec_with_identified
%type <users> user_spec_list drop_user_spec_list user_spec_list_of_create_user
//%type <tlsOptions> require_clause_opt require_clause require_list
//%type <tlsOption> require_elem
//%type <resourceOptions> conn_option_list conn_options
//%type <resourceOption> conn_option
%type <updateExpr> update_value
%type <updateExprs> update_list on_duplicate_key_update_opt
%type <completionType> completion_type
%type <str> password_opt
%type <boolVal> grant_option_opt enforce enforce_opt

%type <varAssignmentExpr> var_assignment
%type <varAssignmentExprs> var_assignment_list
%type <str> var_name equal_or_assignment
%type <strs> var_name_list
%type <expr> set_expr
//%type <setRole> set_role_opt
%type <setDefaultRole> set_default_role_opt
%type <privilege> priv_elem
%type <privileges> priv_list
%type <objectType> object_type
%type <privilegeType> priv_type
%type <privilegeLevel> priv_level
%type <unresolveNames> column_name_list
%type <partitionOption> partition_by_opt
%type <clusterByOption> cluster_by_opt
%type <partitionBy> partition_method sub_partition_method sub_partition_opt
%type <windowSpec> window_spec_opt window_spec
%type <frameClause> window_frame_clause window_frame_clause_opt
%type <frameBound> frame_bound frame_bound_start
%type <frameType> frame_type
%type <str> fields_or_columns
%type <int64Val> algorithm_opt partition_num_opt sub_partition_num_opt
%type <boolVal> linear_opt
%type <partition> partition
%type <partitions> partition_list_opt partition_list
%type <values> values_opt
%type <tableOptions> partition_option_list
%type <subPartition> sub_partition
%type <subPartitions> sub_partition_list sub_partition_list_opt
%type <subquery> subquery
%type <incrementByOption> increment_by_opt
%type <minValueOption> min_value_opt
%type <maxValueOption> max_value_opt
%type <startWithOption> start_with_opt
%type <cycleOption> alter_cycle_opt
%type <alterTypeOption> alter_as_datatype_opt

%type <lengthOpt> length_opt length_option_opt length timestamp_option_opt
%type <lengthScaleOpt> float_length_opt decimal_length_opt
%type <unsignedOpt> unsigned_opt header_opt parallel_opt
%type <zeroFillOpt> zero_fill_opt
%type <boolVal> global_scope exists_opt distinct_opt temporary_opt cycle_opt drop_table_opt
%type <item> pwd_expire clear_pwd_opt
%type <str> name_confict distinct_keyword separator_opt
%type <insert> insert_data
%type <replace> replace_data
%type <rowsExprs> values_list
%type <str> name_datetime_scale braces_opt name_braces
%type <str> std_dev_pop extended_opt
%type <expr> expr_or_default
%type <exprs> data_values data_opt row_value

%type <boolVal> local_opt
%type <duplicateKey> duplicate_opt
%type <fields> load_fields field_item export_fields
%type <fieldsList> field_item_list
%type <str> field_terminator starting_opt lines_terminated_opt starting lines_terminated
%type <lines> load_lines export_lines_opt
%type <int64Val> ignore_lines
%type <varExpr> user_variable variable system_variable
%type <varExprs> variable_list
%type <loadColumn> columns_or_variable
%type <loadColumns> columns_or_variable_list columns_or_variable_list_opt
%type <unresolvedName> normal_ident
%type <updateExpr> load_set_item
%type <updateExprs> load_set_list load_set_spec_opt
%type <strs> index_name_and_type_opt index_name_list
%type <str> index_name index_type key_or_index_opt key_or_index insert_method_options
// type <str> mo_keywords
%type <properties> properties_list
%type <property> property_elem
%type <assignments> set_value_list
%type <assignment> set_value
%type <str> row_opt substr_option
%type <str> time_unit time_stamp_unit
%type <whenClause> when_clause
%type <whenClauseList> when_clause_list
%type <withClause> with_clause
%type <cte> common_table_expr
%type <cteList> cte_list

%type <epxlainOptions> utility_option_list
%type <epxlainOption> utility_option_elem
%type <str> utility_option_name utility_option_arg
%type <str> explain_option_key select_option_opt
%type <str> explain_foramt_value trim_direction
%type <str> priority_opt priority quick_opt ignore_opt wild_opt

%type <str> account_name account_admin_name account_role_name
%type <accountAuthOption> account_auth_option
%type <alterAccountAuthOption> alter_account_auth_option
%type <accountIdentified> account_identified
%type <accountStatus> account_status_option
%type <accountComment> account_comment_opt
%type <accountCommentOrAttribute> user_comment_or_attribute_opt
%type <stageComment> stage_comment_opt
%type <stageStatus> stage_status_opt
%type <stageUrl> stage_url_opt
%type <stageCredentials> stage_credentials_opt
%type <userIdentified> user_identified user_identified_opt
%type <accountRole> default_role_opt

%type <indexHintType> index_hint_type
%type <indexHintScope> index_hint_scope
%type <indexHint> index_hint
%type <indexHintList> index_hint_list index_hint_list_opt

%token <str> KILL
%type <killOption> kill_opt
%token <str> BACKUP FILESYSTEM
%type <statementOption> statement_id_opt
%token <str> QUERY_RESULT
%start start_command

%type<tableLock> table_lock_elem
%type<tableLocks> table_lock_list
%type<tableLockType> table_lock_type

%type<transactionCharacteristic> transaction_characteristic
%type<transactionCharacteristicList> transaction_characteristic_list
%type<isolationLevel> isolation_level
%type<accessMode> access_mode
%%

start_command:
    stmt_type


stmt_type:
    block_stmt
    {
        yylex.(*Lexer).AppendStmt($1)
    }
|   stmt_list

stmt_list:
    stmt
    {
        if $1 != nil {
            yylex.(*Lexer).AppendStmt($1)
        }
    }
|   stmt_list ';' stmt
    {
        if $3 != nil {
            yylex.(*Lexer).AppendStmt($3)
        }
    }

block_stmt:
    SPBEGIN stmt_list_return END
    {
        // Stmts
        $$ = tree.NewCompoundStmt(
            $2, 
            yylex.(*Lexer).buf,
        )
    }

stmt_list_return:
    block_type_stmt
    {
        $$ = buffer.MakeSlice[tree.Statement](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Statement](yylex.(*Lexer).buf, $$, $1)
    }
|   stmt_list_return ';' block_type_stmt
    {
        $$ = buffer.MakeSlice[tree.Statement](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Statement](yylex.(*Lexer).buf, $1, $3)
    }

block_type_stmt:
    block_stmt
|   case_stmt
|   if_stmt
|   loop_stmt
|   repeat_stmt
|   while_stmt
|   iterate_stmt
|   leave_stmt
|   normal_stmt
|   declare_stmt
    {
        $$ = $1
    }
|   /* EMPTY */
    {
        n := *buffer.Alloc[tree.Statement](yylex.(*Lexer).buf)
        n = nil
        $$ = n
    }

stmt:
    normal_stmt
    {
        $$ = $1
    }
|   declare_stmt
|   transaction_stmt
    {
        $$ = $1
    }
|   /* EMPTY */
    {
        n := *buffer.Alloc[tree.Statement](yylex.(*Lexer).buf)
        n = nil
        $$ = n
    }

normal_stmt:
    create_stmt
|   call_stmt
|   mo_dump_stmt
|   insert_stmt
|   replace_stmt
|   delete_stmt
|   drop_stmt
|   truncate_table_stmt
|   explain_stmt
|   prepare_stmt
|   deallocate_stmt
|   reset_stmt
|   execute_stmt
|   show_stmt
|   alter_stmt
|   analyze_stmt
|   update_stmt
|   use_stmt
|   set_stmt
|   lock_stmt
|   revoke_stmt
|   grant_stmt
|   load_data_stmt
|   load_extension_stmt
|   do_stmt
|   values_stmt
|   select_stmt
    {
        $$ = $1
    }
|   kill_stmt
|   backup_stmt

backup_stmt:
    BACKUP STRING FILESYSTEM STRING
	{
        // Timestamp IsS3 Dir Option
		$$ = tree.NewBackupStart(
            $2, 
            false, 
            $4, 
            nil, 
            yylex.(*Lexer).buf,
        )
	}
    | BACKUP STRING S3OPTION '{' infile_or_s3_params '}'
    {
        // Timestamp IsS3 Dir Option
    	$$ = tree.NewBackupStart(
            $2, 
            true, 
            "", 
            $5, 
            yylex.(*Lexer).buf,
        )
    }

kill_stmt:
    KILL kill_opt INTEGRAL statement_id_opt
    {
        var connectionId uint64
        switch v := $3.(type) {
            case uint64:
	            connectionId = v
            case int64:
	            connectionId = uint64(v)
            default:
	            yylex.Error("parse integral fail")
	        return 1
        }

        // Option ConnectionId StmtOption
	    $$ = tree.NewKill(
                $2, 
                connectionId, 
                $4, 
                yylex.(*Lexer).buf,
            )
    }

kill_opt:
{
    // Exist Typ
    $$ = tree.NewKillOption(
        false,
        0, 
        yylex.(*Lexer).buf,
    )
}
| CONNECTION
{
    // Exist Typ
    $$ = tree.NewKillOption(
        true, 
        tree.KillTypeConnection, 
        yylex.(*Lexer).buf,
    )
}
| QUERY
{
    // Exist Typ
    $$ = tree.NewKillOption(
        true, 
        tree.KillTypeQuery, 
        yylex.(*Lexer).buf,
    )
}

statement_id_opt:
{
    // Exist StatementId
    $$ = tree.NewStatementOption(
        false, 
        "", 
        yylex.(*Lexer).buf,
    )
}
| STRING
{
    // Exist StatementId
    $$ = tree.NewStatementOption(
        true, 
        $1, 
        yylex.(*Lexer).buf,
    )
}

call_stmt:
    CALL proc_name '(' expression_list_opt ')'
    {
        // Name Args
        $$ = tree.NewCallStmt(
            $2, 
            $4, 
            yylex.(*Lexer).buf,
        )
    }

leave_stmt:
    LEAVE ident
    {
        // Name
        $$ = tree.NewLeaveStmt(
            tree.Identifier($2.ToLower()), 
            yylex.(*Lexer).buf,
        )
    }

iterate_stmt:
    ITERATE ident
    {
        // Name
        $$ = tree.NewIterateStmt(
            tree.Identifier($2.ToLower()), 
            yylex.(*Lexer).buf,
        )
    }

while_stmt:
    WHILE expression DO stmt_list_return END WHILE
    {
        // Name Cond Body
        $$ = tree.NewWhileStmt(
            "", 
            $2, 
            $4, 
            yylex.(*Lexer).buf,
        )
    }
|   ident ':' WHILE expression DO stmt_list_return END WHILE ident
    {
        // Name Cond Body
        $$ = tree.NewWhileStmt(
            tree.Identifier($1.ToLower()), 
            $4, 
            $6, 
            yylex.(*Lexer).buf,
        )
    }

repeat_stmt:
    REPEAT stmt_list_return UNTIL expression END REPEAT
    {
        // Name Body Cond
        $$ = tree.NewRepeatStmt(
            "",
            $2,
            $4, 
            yylex.(*Lexer).buf,
        )
    }
|    ident ':' REPEAT stmt_list_return UNTIL expression END REPEAT ident
    {
        // Name Body Cond
        $$ = tree.NewRepeatStmt(
            tree.Identifier($1.ToLower()),
            $4, 
            $6, yylex.(*Lexer).buf,
        )
    }

loop_stmt:
    LOOP stmt_list_return END LOOP
    {
        // Name Body
        $$ = tree.NewLoopStmt(
            "", 
            $2, 
            yylex.(*Lexer).buf,
        )
    }
|   ident ':' LOOP stmt_list_return END LOOP ident
    {
        // Name Body
        $$ = tree.NewLoopStmt(
            tree.Identifier($1.ToLower()), 
            $4, 
            yylex.(*Lexer).buf,
        )
    }

if_stmt:
    IF expression THEN stmt_list_return elseif_clause_list_opt else_clause_opt2 END IF
    {
        // cond body elifs else
        $$ = tree.NewIfStmt(
            $2, 
            $4, 
            $5, 
            $6, 
            yylex.(*Lexer).buf,
        )
    }

elseif_clause_list_opt:
    {
        $$ = nil
    }
|    elseif_clause_list
    {
        $$ = $1
    }

elseif_clause_list:
    elseif_clause
    {
        $$ = buffer.MakeSlice[*tree.ElseIfStmt](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.ElseIfStmt](yylex.(*Lexer).buf, $$, $1)
    }
|   elseif_clause_list elseif_clause
    {
        $$ = buffer.MakeSlice[*tree.ElseIfStmt](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.ElseIfStmt](yylex.(*Lexer).buf, $1, $2)
    }

elseif_clause:
    ELSEIF expression THEN stmt_list_return
    {
        // Cond Body
        $$ = tree.NewElseIfStmt(
            $2,
            $4, 
            yylex.(*Lexer).buf,
        )
    }

case_stmt:
    CASE expression when_clause_list2 else_clause_opt2 END CASE
    {
        // Expr Whens Else
        $$ = tree.NewCaseStmt(
            $2,
            $3, 
            $4, 
            yylex.(*Lexer).buf,
        )
    }

when_clause_list2:
    when_clause2
    {
        $$ = buffer.MakeSlice[*tree.WhenStmt](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.WhenStmt](yylex.(*Lexer).buf, $$, $1)
    }
|   when_clause_list2 when_clause2
    {
        $$ = buffer.MakeSlice[*tree.WhenStmt](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.WhenStmt](yylex.(*Lexer).buf, $1, $2)
    }

when_clause2:
    WHEN expression THEN stmt_list_return
    {
        // Cond Body
        $$ = tree.NewWhenStmt(
            $2,
            $4, 
            yylex.(*Lexer).buf,
        )
    }

else_clause_opt2:
    /* empty */
    {
        $$ = nil
    }
|   ELSE stmt_list_return
    {
        $$ = $2
    }

mo_dump_stmt:
   MODUMP QUERY_RESULT STRING INTO STRING export_fields export_lines_opt header_opt max_file_size_opt force_quote_opt
    {
        // Outfile QueryId FilePath Fields Lines Header MaxFileSize ForeQuote
        ep := tree.NewExportParam(
            true, 
            $3, 
            $5, 
            $6, 
            $7, 
            $8, 
            uint64($9)*1024, 
            $10, 
            yylex.(*Lexer).buf,
        )

        // ExportParams
        $$ = tree.NewMoDump(
            ep, 
            yylex.(*Lexer).buf,
        )
    }

load_data_stmt:
    LOAD DATA local_opt load_param_opt duplicate_opt INTO TABLE table_name tail_param_opt parallel_opt
    {
        // Local Param DuplicateHandling Table
        $$ = tree.NewLoad(
            $3, 
            $4, 
            $5, 
            $8, 
            yylex.(*Lexer).buf,
        )
        $$.(*tree.Load).Param.Tail = $9
        $$.(*tree.Load).Param.Parallel = $10
    }

load_extension_stmt:
    LOAD extension_name
    {
        // Name
        $$ = tree.NewLoadExtension(
            tree.Identifier($2), 
            yylex.(*Lexer).buf,
        )
    }

load_set_spec_opt:
    {
        $$ = nil
    }
|   SET load_set_list
    {
        $$ = $2
    }

load_set_list:
    load_set_item
    {
        $$ = buffer.MakeSlice[*tree.UpdateExpr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.UpdateExpr](yylex.(*Lexer).buf, $$, $1)
    }
|   load_set_list ',' load_set_item
    {
        $$ = buffer.MakeSlice[*tree.UpdateExpr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.UpdateExpr](yylex.(*Lexer).buf, $1, $3)
    }

load_set_item:
    normal_ident '=' DEFAULT
    {
        names := buffer.MakeSlice[*tree.UnresolvedName](yylex.(*Lexer).buf)
        names = buffer.AppendSlice[*tree.UnresolvedName](yylex.(*Lexer).buf, names, $1)
        // Names Expr
        $$ = tree.NewUpdateExpr(
            names,
            tree.NewDefaultVal(yylex.(*Lexer).buf),
            yylex.(*Lexer).buf,
        )
    }
|   normal_ident '=' expression
    {
        names := buffer.MakeSlice[*tree.UnresolvedName](yylex.(*Lexer).buf)
        names = buffer.AppendSlice[*tree.UnresolvedName](yylex.(*Lexer).buf, names, $1)

        // Names Expr
        $$ = tree.NewUpdateExpr( 
            names,
            $3,
            yylex.(*Lexer).buf,
        )
    }

parallel_opt:
    {
        $$ = false
    }
|   PARALLEL STRING
    {
        str := strings.ToLower($2)
        if str == "true" {
            $$ = true
        } else if str == "false" {
            $$ = false
        } else {
            yylex.Error("error parallel flag")
            return 1
        }
    }

normal_ident:
    ident
    {
        $$ = tree.SetUnresolvedName(yylex.(*Lexer).buf, $1.Compare())
    }
|   ident '.' ident
    {
        $$ = tree.SetUnresolvedName(yylex.(*Lexer).buf, $1.Compare(), $3.Compare())
    }
|   ident '.' ident '.' ident
    {
        $$ = tree.SetUnresolvedName(yylex.(*Lexer).buf, $1.Compare(), $3.Compare(), $5.Compare())
    }

columns_or_variable_list_opt:
    {
        $$ = nil
    }
|   '(' ')'
    {
        $$ = nil
    }
|   '(' columns_or_variable_list ')'
    {
        $$ = $2
    }

columns_or_variable_list:
    columns_or_variable
    {
        switch $1.(type) {
        case *tree.UnresolvedName:
            $$ = buffer.MakeSlice[tree.LoadColumn](yylex.(*Lexer).buf)
            $$ = buffer.AppendSlice[tree.LoadColumn](yylex.(*Lexer).buf, $$, $1.(*tree.UnresolvedName))
        case *tree.VarExpr:
            $$ = buffer.MakeSlice[tree.LoadColumn](yylex.(*Lexer).buf)
            $$ = buffer.AppendSlice[tree.LoadColumn](yylex.(*Lexer).buf, $$, $1.(*tree.VarExpr))
        }
    }
|   columns_or_variable_list ',' columns_or_variable
    {
        $$ = buffer.MakeSlice[tree.LoadColumn](yylex.(*Lexer).buf)
        switch $3.(type) {
        case *tree.UnresolvedName:
            $$ = buffer.MakeSlice[tree.LoadColumn](yylex.(*Lexer).buf)
            $$ = buffer.AppendSlice[tree.LoadColumn](yylex.(*Lexer).buf, $1, $3.(*tree.UnresolvedName))
        case *tree.VarExpr:
            $$ = buffer.MakeSlice[tree.LoadColumn](yylex.(*Lexer).buf)
            $$ = buffer.AppendSlice[tree.LoadColumn](yylex.(*Lexer).buf, $1, $3.(*tree.VarExpr))
        }
    }

columns_or_variable:
    column_name_unresolved
    {
        $$ = $1
    }
|   user_variable
    {
        $$ = $1
    }

variable_list:
    variable
    {
        $$ = buffer.MakeSlice[*tree.VarExpr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.VarExpr](yylex.(*Lexer).buf, $$, $1)
    }
|   variable_list ',' variable
    {
        $$ = buffer.MakeSlice[*tree.VarExpr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.VarExpr](yylex.(*Lexer).buf, $1, $3)
    }

variable:
    system_variable
    {
        $$ = $1
    }
|   user_variable
    {
        $$ = $1
    }

system_variable:
    AT_AT_ID
    {
        vs := strings.Split($1, ".")
        var isGlobal bool
        if strings.ToLower(vs[0]) == "global" {
            isGlobal = true
        }
        var r string
        if len(vs) == 2 {
           r = vs[1]
        } else if len(vs) == 1 {
           r = vs[0]
        } else {
            yylex.Error("variable syntax error")
            return 1
        }
        // Name System Global
        $$ = tree.NewVarExpr( 
            r,
            true,
            isGlobal,
            yylex.(*Lexer).buf,
        )
    }

user_variable:
    AT_ID
    {
//        vs := strings.Split($1, ".")
//        var r string
//        if len(vs) == 2 {
//           r = vs[1]
//        } else if len(vs) == 1 {
//           r = vs[0]
//        } else {
//            yylex.Error("variable syntax error")
//            return 1
//        }
        // Name System Global
        $$ = tree.NewVarExpr( 
            $1,
            false,
            false,
            yylex.(*Lexer).buf,
        )
    }

ignore_lines:
    {
        $$ = 0
    }
|   IGNORE INTEGRAL LINES
    {
        $$ = $2.(int64)
    }
|   IGNORE INTEGRAL ROWS
    {
        $$ = $2.(int64)
    }

load_lines:
    {
        $$ = nil
    }
|   LINES starting lines_terminated_opt
    {
        // StartingBy TerminatedBy
        $$ = tree.NewLines( 
            $2,
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   LINES lines_terminated starting_opt
    {
        // StartingBy TerminatedBy
        $$ = tree.NewLines( 
            $3,
            $2,
            yylex.(*Lexer).buf,
        )
    }

starting_opt:
    {
        $$ = ""
    }
|   starting

starting:
    STARTING BY STRING
    {
        $$ = $3
    }

lines_terminated_opt:
    {
        $$ = "\n"
    }
|   lines_terminated

lines_terminated:
    TERMINATED BY STRING
    {
        $$ = $3
    }

load_fields:
    {
        $$ = nil
    }
|   fields_or_columns field_item_list
    {
        // Terminated Optionally EnclosedBy EscapedBy
        res := tree.NewFields( 
            "\t",
            false,
            0,
            0,
            yylex.(*Lexer).buf,
        )
        for _, f := range $2 {
            if f.Terminated.Get() != "" {
                res.Terminated = f.Terminated
            }
            if f.Optionally {
                res.Optionally = f.Optionally
            }
            if f.EnclosedBy != 0 {
                res.EnclosedBy = f.EnclosedBy
            }
            if f.EscapedBy != 0 {
                res.EscapedBy = f.EscapedBy
            }
        }
        $$ = res
    }

field_item_list:
    field_item
    {
        $$ = buffer.MakeSlice[*tree.Fields](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.Fields](yylex.(*Lexer).buf, $$, $1)
    }
|   field_item_list field_item
    {
        $$ = buffer.MakeSlice[*tree.Fields](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.Fields](yylex.(*Lexer).buf, $1, $2)
    }

field_item:
    TERMINATED BY field_terminator
    {
        // Terminated Optionally EnclosedBy EscapedBy
        $$ = tree.NewFields( 
            $3,
            false,
            0,
            0,
            yylex.(*Lexer).buf,
        )
    }
|   OPTIONALLY ENCLOSED BY field_terminator
    {
        str := $4
        if str != "\\" && len(str) > 1 {
            yylex.Error("error field terminator")
            return 1
        }
        var b byte
        if len(str) != 0 {
            b = byte(str[0])
        } else {
            b = 0
        }
        // Terminated Optionally EnclosedBy EscapedBy
        $$ = tree.NewFields( 
            "",
            true,
            b,
            0,
            yylex.(*Lexer).buf,
        )
    }
|   ENCLOSED BY field_terminator
    {
        str := $3
        if str != "\\" && len(str) > 1 {
            yylex.Error("error field terminator")
            return 1
        }
        var b byte
        if len(str) != 0 {
           b = byte(str[0])
        } else {
           b = 0
        }
        // Terminated Optionally EnclosedBy EscapedBy
        $$ = tree.NewFields( 
            "",
            false,
            b,
            0,
            yylex.(*Lexer).buf,
        )
    }
|   ESCAPED BY field_terminator
    {
        str := $3
        if str != "\\" && len(str) > 1 {
            yylex.Error("error field terminator")
            return 1
        }
        var b byte
        if len(str) != 0 {
           b = byte(str[0])
        } else {
           b = 0
        }
        // Terminated Optionally EnclosedBy EscapedBy
        $$ = tree.NewFields( 
            "",
            false,
            0,
            b,
            yylex.(*Lexer).buf,
        )
    }

field_terminator:
    STRING
// |   HEXNUM
// |   BIT_LITERAL

duplicate_opt:
    {
        $$ = tree.NewDuplicateKeyError(yylex.(*Lexer).buf)
    }
|   IGNORE
    {
        $$ = tree.NewDuplicateKeyIgnore(yylex.(*Lexer).buf)
    }
|   REPLACE
    {
        $$ = tree.NewDuplicateKeyReplace(yylex.(*Lexer).buf)
    }

local_opt:
    {
        $$ = false
    }
|   LOCAL
    {
        $$ = true
    }

grant_stmt:
    GRANT priv_list ON object_type priv_level TO role_spec_list grant_option_opt
    {
        // GrantType GrantPrivilege GrantRole GrantProxy
        $$ = tree.NewGrant( 
            tree.GrantTypePrivilege,
            // Privilege ObjectType PrivilegeLevel Role globalOption
            tree.NewGrantPrivilege(
                $2,
                $4,
                $5,
                $7,
                $8,
                yylex.(*Lexer).buf,
            ),
            nil,
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   GRANT role_spec_list TO drop_user_spec_list grant_option_opt
    {
        // GrantType GrantPrivilege GrantRole GrantProxy
        $$ = tree.NewGrant( 
            tree.GrantTypeRole,
            nil,
            // Roles Users GrantOption
            tree.NewGrantRole(
                $2,
                $4,
                $5,
                yylex.(*Lexer).buf,
            ),
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   GRANT PROXY ON user_spec TO user_spec_list grant_option_opt
    {
        // GrantType GrantPrivilege GrantRole GrantProxy
        $$ =  tree.NewGrant( 
            tree.GrantTypeProxy,
            nil,
            nil,
            // ProxyUser Users GrantOption
            tree.NewGrantProxy(
                $4,
                $6,
                $7,
                yylex.(*Lexer).buf,
            ),
            yylex.(*Lexer).buf,
        )
    }

grant_option_opt:
    {
        $$ = false
    }
|   WITH GRANT OPTION
    {
        $$ = true
    }
// |    WITH MAX_QUERIES_PER_HOUR INTEGRAL
// |    WITH MAX_UPDATES_PER_HOUR INTEGRAL
// |    WITH MAX_CONNECTIONS_PER_HOUR INTEGRAL
// |    WITH MAX_USER_CONNECTIONS INTEGRAL

revoke_stmt:
    REVOKE exists_opt  priv_list ON object_type priv_level FROM role_spec_list
    {
        // RevokeType RevokePrivileges RevokeRole
        $$ = tree.NewRevoke( 
            tree.RevokeTypePrivilege,
            // IfExists Privileges ObjType Level Roles
            tree.NewRevokePrivilege(
                $2,
                $3,
                $5,
                $6,
                $8,
                yylex.(*Lexer).buf,
            ),
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   REVOKE exists_opt role_spec_list FROM user_spec_list
    {
        // RevokeType RevokePrivileges RevokeRole
        $$ = tree.NewRevoke(
            tree.RevokeTypeRole,
            nil,
            tree.NewRevokeRole(
                $2,
                $3,
                $5,
                yylex.(*Lexer).buf,
            ),
            yylex.(*Lexer).buf,
        )
    }

priv_level:
    '*'
    {
        // level, dbName, tblName
        $$ = tree.NewPrivilegeLevel( 
            tree.PRIVILEGE_LEVEL_TYPE_STAR,
            "",
            "",
            yylex.(*Lexer).buf,
        )
    }
|   '*' '.' '*'
    {
        // level, dbName, tblName
        $$ = tree.NewPrivilegeLevel(
            tree.PRIVILEGE_LEVEL_TYPE_STAR_STAR,
            "",
            "",
            yylex.(*Lexer).buf,
        )

    }
|   ident '.' '*'
    {
        // level, dbName, tblName
        $$ = tree.NewPrivilegeLevel(
            tree.PRIVILEGE_LEVEL_TYPE_DATABASE_STAR,
            $1.Compare(),
            "",
            yylex.(*Lexer).buf,
        )

    }
|   ident '.' ident
    {
        // level, dbName, tblName
        $$ = tree.NewPrivilegeLevel(
            tree.PRIVILEGE_LEVEL_TYPE_DATABASE_TABLE,
            $1.Compare(),
            $3.Compare(),
            yylex.(*Lexer).buf,
        )
    }
|   ident
    {
        // level, dbName, tblName
        $$ = tree.NewPrivilegeLevel(
            tree.PRIVILEGE_LEVEL_TYPE_TABLE,
            "",
            $1.Compare(),
            yylex.(*Lexer).buf,
        )
    }

object_type:
    TABLE
    {
        $$ = tree.OBJECT_TYPE_TABLE
    }
|   DATABASE
    {
        $$ = tree.OBJECT_TYPE_DATABASE
    }
|   FUNCTION
    {
        $$ = tree.OBJECT_TYPE_FUNCTION
    }
|   PROCEDURE
    {
        $$ = tree.OBJECT_TYPE_PROCEDURE
    }
|   VIEW
    {
        $$ = tree.OBJECT_TYPE_VIEW
    }
|   ACCOUNT
    {
        $$ = tree.OBJECT_TYPE_ACCOUNT
    }


priv_list:
    priv_elem
    {
        $$ = buffer.MakeSlice[*tree.Privilege](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.Privilege](yylex.(*Lexer).buf, $$, $1)
    }
|   priv_list ',' priv_elem
    {
        $$ = buffer.MakeSlice[*tree.Privilege](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.Privilege](yylex.(*Lexer).buf, $1, $3)
    }

priv_elem:
    priv_type
    {
        // PrivilegeType []*UnresolvedName
        $$ = tree.NewPrivilege( 
            $1,
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   priv_type '(' column_name_list ')'
    {
        // PrivilegeType []*UnresolvedName
        $$ = tree.NewPrivilege( 
            $1,
            $3,
            yylex.(*Lexer).buf,
        )
    }

column_name_list:
    column_name
    {
        $$ = buffer.MakeSlice[*tree.UnresolvedName](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.UnresolvedName](yylex.(*Lexer).buf, $$, $1)
    }
|   column_name_list ',' column_name
    {
        $$ = buffer.MakeSlice[*tree.UnresolvedName](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.UnresolvedName](yylex.(*Lexer).buf, $1, $3)
    }

priv_type:
    ALL
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_ALL
    }
|    CREATE ACCOUNT
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_CREATE_ACCOUNT
    }
|    DROP ACCOUNT
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_DROP_ACCOUNT
    }
|    ALTER ACCOUNT
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_ALTER_ACCOUNT
    }
|    ALL PRIVILEGES
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_ALL
    }
|    ALTER TABLE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_ALTER_TABLE
    }
|    ALTER VIEW
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_ALTER_VIEW
    }
|    CREATE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_CREATE
    }
|    CREATE USER
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_CREATE_USER
    }
|    DROP USER
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_DROP_USER
    }
|    ALTER USER
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_ALTER_USER
    }
|    CREATE TABLESPACE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_CREATE_TABLESPACE
    }
|    TRIGGER
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_TRIGGER
    }
|    DELETE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_DELETE
    }
|    DROP TABLE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_DROP_TABLE
    }
|    DROP VIEW
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_DROP_VIEW
    }
|    EXECUTE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_EXECUTE
    }
|    INDEX
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_INDEX
    }
|    INSERT
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_INSERT
    }
|    SELECT
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_SELECT
    }
|    SUPER
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_SUPER
    }
|    CREATE DATABASE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_CREATE_DATABASE
    }
|    DROP DATABASE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_DROP_DATABASE
    }
|    SHOW DATABASES
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_SHOW_DATABASES
    }
|    CONNECT
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_CONNECT
    }
|    MANAGE GRANTS
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_MANAGE_GRANTS
    }
|    OWNERSHIP
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_OWNERSHIP
    }
|    SHOW TABLES
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_SHOW_TABLES
    }
|    CREATE TABLE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_CREATE_TABLE
    }
|    UPDATE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_UPDATE
    }
|    GRANT OPTION
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_GRANT_OPTION
    }
|    REFERENCES
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_REFERENCES
    }
|    REFERENCE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_REFERENCE
    }
|    REPLICATION SLAVE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_REPLICATION_SLAVE
    }
|    REPLICATION CLIENT
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_REPLICATION_CLIENT
    }
|    USAGE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_USAGE
    }
|    RELOAD
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_RELOAD
    }
|    FILE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_FILE
    }
|    CREATE TEMPORARY TABLES
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_CREATE_TEMPORARY_TABLES
    }
|    LOCK TABLES
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_LOCK_TABLES
    }
|    CREATE VIEW
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_CREATE_VIEW
    }
|    SHOW VIEW
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_SHOW_VIEW
    }
|    CREATE ROLE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_CREATE_ROLE
    }
|    DROP ROLE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_DROP_ROLE
    }
|    ALTER ROLE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_ALTER_ROLE
    }
|      CREATE ROUTINE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_CREATE_ROUTINE
    }
|    ALTER ROUTINE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_ALTER_ROUTINE
    }
|    EVENT
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_EVENT
    }
|    SHUTDOWN
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_SHUTDOWN
    }
|    TRUNCATE
    {
        $$ = tree.PRIVILEGE_TYPE_STATIC_TRUNCATE
    }

set_stmt:
    set_variable_stmt
|   set_password_stmt
|   set_role_stmt
|   set_default_role_stmt
|   set_transaction_stmt

set_transaction_stmt:
    SET TRANSACTION transaction_characteristic_list
    {
        // Global CharacterList
    	$$ = tree.NewSetTransaction( 
	        false,
	        $3,
            yylex.(*Lexer).buf,
        )
    }
|   SET GLOBAL TRANSACTION transaction_characteristic_list
    {
        // Global CharacterList
        $$ = tree.NewSetTransaction( 
            true,
            $4,
            yylex.(*Lexer).buf,
        )
    }
|   SET SESSION TRANSACTION transaction_characteristic_list
    {
        // Global CharacterList
        $$ = tree.NewSetTransaction( 
            false,
            $4,
            yylex.(*Lexer).buf,
        )
    }


transaction_characteristic_list:
    transaction_characteristic
    {
        $$ = buffer.MakeSlice[*tree.TransactionCharacteristic](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.TransactionCharacteristic](yylex.(*Lexer).buf, $$, $1)
    }
|   transaction_characteristic_list ',' transaction_characteristic
    {
        $$ = buffer.MakeSlice[*tree.TransactionCharacteristic](yylex.(*Lexer).buf)
	    $$ = buffer.AppendSlice[*tree.TransactionCharacteristic](yylex.(*Lexer).buf, $1, $3)
    }

transaction_characteristic:
    ISOLATION LEVEL isolation_level
    {
        // Islevel Isolation Access
        $$ = tree.NewTransactionCharacteristic( 
            true,
            $3,
            0,
            yylex.(*Lexer).buf,
        )
    }
|   access_mode
    {
        // Islevel Isolation Access
        $$ = tree.NewTransactionCharacteristic( 
            false,
            0,
            $1,
            yylex.(*Lexer).buf,
        )
    }

isolation_level:
    REPEATABLE READ
    {
    	$$ = tree.ISOLATION_LEVEL_REPEATABLE_READ
    }
|   READ COMMITTED
    {
    	$$ = tree.ISOLATION_LEVEL_READ_COMMITTED
    }
|   READ UNCOMMITTED
    {
    	$$ = tree.ISOLATION_LEVEL_READ_UNCOMMITTED
    }
|   SERIALIZABLE
    {
    	$$ = tree.ISOLATION_LEVEL_SERIALIZABLE
    }

access_mode:
   READ WRITE
   {
    	$$ = tree.ACCESS_MODE_READ_WRITE
   }
| READ ONLY
   {
    	$$ = tree.ACCESS_MODE_READ_ONLY
   }

set_role_stmt:
    SET ROLE role_spec
    {
        // SecondaryRole SecondaryRoleType Role
        $$ = tree.NewSetRole( 
            false,
            0,
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   SET SECONDARY ROLE ALL
    {
        // SecondaryRole SecondaryRoleType Role
        $$ = tree.NewSetRole( 
            true,
            tree.SecondaryRoleTypeAll,
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   SET SECONDARY ROLE NONE
    {
        // SecondaryRole SecondaryRoleType Role
        $$ = tree.NewSetRole( 
            true,
            tree.SecondaryRoleTypeNone,
            nil,
            yylex.(*Lexer).buf,
        )
    }

set_default_role_stmt:
    SET DEFAULT ROLE set_default_role_opt TO user_spec_list
    {
        dr := $4
        dr.Users = $6
        $$ = dr
    }

//set_role_opt:
//    ALL EXCEPT role_spec_list
//    {
//        $$ = &tree.SetRole{Type: tree.SET_ROLE_TYPE_ALL_EXCEPT, Roles: $3}
//    }
//|   DEFAULT
//    {
//        $$ = &tree.SetRole{Type: tree.SET_ROLE_TYPE_DEFAULT, Roles: nil}
//    }
//|   NONE
//    {
//        $$ = &tree.SetRole{Type: tree.SET_ROLE_TYPE_NONE, Roles: nil}
//    }
//|   ALL
//    {
//        $$ = &tree.SetRole{Type: tree.SET_ROLE_TYPE_ALL, Roles: nil}
//    }
//|   role_spec_list
//    {
//        $$ = &tree.SetRole{Type: tree.SET_ROLE_TYPE_NORMAL, Roles: $1}
//    }

set_default_role_opt:
    NONE
    {
        // Type Roles Users
        $$ = tree.NewSetDefaultRole( 
            tree.SET_DEFAULT_ROLE_TYPE_NONE, 
            nil,
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   ALL
    {
        // Type Roles Users
        $$ = tree.NewSetDefaultRole( 
            tree.SET_DEFAULT_ROLE_TYPE_ALL, 
            nil,
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   role_spec_list
    {
        // Type Roles Users
        $$ = tree.NewSetDefaultRole( 
            tree.SET_DEFAULT_ROLE_TYPE_NORMAL,
            $1,
            nil,
            yylex.(*Lexer).buf,
        )
    }

set_variable_stmt:
    SET var_assignment_list
    {
        $$ = tree.NewSetVar( 
            $2,
            yylex.(*Lexer).buf,
        )
    }

set_password_stmt:
    SET PASSWORD '=' password_opt
    {
        // User Password
        $$ = tree.NewSetPassword(
            nil,
            $4,
            yylex.(*Lexer).buf,
        )
    }
|   SET PASSWORD FOR user_spec '=' password_opt
    {
        // User Password
        $$ = tree.NewSetPassword(
            $4,
            $6,
            yylex.(*Lexer).buf,
        )
    }

password_opt:
    STRING
|   PASSWORD '(' auth_string ')'
    {
        $$ = $3
    }

var_assignment_list:
    var_assignment
    {
        $$ = buffer.MakeSlice[*tree.VarAssignmentExpr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.VarAssignmentExpr](yylex.(*Lexer).buf, $$, $1)
    }
|   var_assignment_list ',' var_assignment
    {
        $$ = buffer.MakeSlice[*tree.VarAssignmentExpr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.VarAssignmentExpr](yylex.(*Lexer).buf, $1, $3)
    }

var_assignment:
    var_name equal_or_assignment set_expr
    {
        // System Global Name Value Reserved
        $$ = tree.NewVarAssignmentExpr(
            true,
            false,
            $1,
            $3,
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   GLOBAL var_name equal_or_assignment set_expr
    {
        // System Global Name Value Reserved
        $$ = tree.NewVarAssignmentExpr(
            true,
            true,
            $2,
            $4,
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   PERSIST var_name equal_or_assignment set_expr
    {
        // System Global Name Value Reserved
        $$ = tree.NewVarAssignmentExpr(
            true,
            true,
            $2,
            $4,
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   SESSION var_name equal_or_assignment set_expr
    {
        // System Global Name Value Reserved
        $$ = tree.NewVarAssignmentExpr(
            true,
            false,
            $2,
            $4,
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   LOCAL var_name equal_or_assignment set_expr
    {
        // System Global Name Value Reserved
        $$ = tree.NewVarAssignmentExpr(
            true,
            false,
            $2,
            $4,
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   AT_ID equal_or_assignment set_expr
    {
        vs := strings.Split($1, ".")
        var isGlobal bool
        if strings.ToLower(vs[0]) == "global" {
            isGlobal = true
        }
        var r string
        if len(vs) == 2 {
            r = vs[1]
        } else if len(vs) == 1{
            r = vs[0]
        } else {
            yylex.Error("variable syntax error")
            return 1
        }
        // System Global Name Value Reserved
        $$ = tree.NewVarAssignmentExpr(
            false, 
            isGlobal, 
            r, 
            $3, 
            nil, 
            yylex.(*Lexer).buf,
        )
    }
|   AT_AT_ID equal_or_assignment set_expr
    {
        vs := strings.Split($1, ".")
        var isGlobal bool
        if strings.ToLower(vs[0]) == "global" {
            isGlobal = true
        }
        var r string
        if len(vs) == 2 {
            r = vs[1]
        } else if len(vs) == 1{
            r = vs[0]
        } else {
            yylex.Error("variable syntax error")
            return 1
        }
        // System Global Name Value Reserved
        $$ = tree.NewVarAssignmentExpr(true,
            isGlobal, 
            r, 
            $3, 
            nil, 
            yylex.(*Lexer).buf,
        ) 
    }
|   NAMES charset_name
    {
        // System Global Name Value Reserved
        $$ = tree.NewVarAssignmentExpr(
            false,
            false,
            strings.ToLower($1),
            tree.NewNumValWithType(constant.MakeString($2), $2, false, tree.P_char, yylex.(*Lexer).buf),
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   NAMES charset_name COLLATE DEFAULT
    {
        // System Global Name Value Reserved
        $$ = tree.NewVarAssignmentExpr(
            false,
            false,
            strings.ToLower($1),
            tree.NewNumValWithType(constant.MakeString($2), $2, false, tree.P_char, yylex.(*Lexer).buf),
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   NAMES charset_name COLLATE name_string
    {
        // System Global Name Value Reserved
        $$ = tree.NewVarAssignmentExpr(
            false,
            false,
            strings.ToLower($1), 
            tree.NewNumValWithType(constant.MakeString($2), $2, false, tree.P_char, yylex.(*Lexer).buf),
            tree.NewNumValWithType(constant.MakeString($4), $4, false, tree.P_char, yylex.(*Lexer).buf),
            yylex.(*Lexer).buf,
        )
    }
|   NAMES DEFAULT
    {
        $$ = tree.NewVarAssignmentExpr(
            false,
            false,
            strings.ToLower($1),
            tree.NewDefaultVal(yylex.(*Lexer).buf),
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   charset_keyword charset_name
    {
        // System Global Name Value Reserved
        $$ = tree.NewVarAssignmentExpr(
            false,
            false,
            strings.ToLower($1),
            tree.NewNumValWithType(constant.MakeString($2), $2, false, tree.P_char, yylex.(*Lexer).buf),
            nil, 
            yylex.(*Lexer).buf,
        )
    }
|   charset_keyword DEFAULT
    {
        // System Global Name Value Reserved
        $$ = tree.NewVarAssignmentExpr(
            false,
            false,
            strings.ToLower($1),
            tree.NewDefaultVal(yylex.(*Lexer).buf),
            nil, 
            yylex.(*Lexer).buf,
        )
    }

set_expr:
    ON
    {
        $$ = tree.NewNumValWithType(constant.MakeString($1), $1, false, tree.P_char, yylex.(*Lexer).buf)
    }
|   BINARY
    {
        $$ = tree.NewNumValWithType(constant.MakeString($1), $1, false, tree.P_char, yylex.(*Lexer).buf)
    }
|   expr_or_default
    {
        $$ = $1
    }

equal_or_assignment:
    '='
    {
        $$ = string($1)
    }
|   ASSIGNMENT
    {
        $$ = $1
    }

var_name:
    ident
    {
    	$$ = $1.Compare()
    }
|   ident '.' ident
    {
        $$ = $1.Compare() + "." + $3.Compare()
    }

var_name_list:
    var_name
    {
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $$, yylex.(*Lexer).buf.CopyString($1))
    }
|   var_name_list ',' var_name
    {
        /* $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf) */
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $1, yylex.(*Lexer).buf.CopyString($3))
    }

transaction_stmt:
    begin_stmt
|   commit_stmt
|   rollback_stmt

rollback_stmt:
    ROLLBACK completion_type
    {
        $$ = tree.NewRollbackTransaction($2, yylex.(*Lexer).buf)
    }

commit_stmt:
    COMMIT completion_type
    {
        $$ = tree.NewCommitTransaction($2, yylex.(*Lexer).buf)
    }

completion_type:
    {
        $$ = tree.COMPLETION_TYPE_NO_CHAIN
    }
|    WORK
    {
        $$ = tree.COMPLETION_TYPE_NO_CHAIN
    }
|   AND CHAIN NO RELEASE
    {
        $$ = tree.COMPLETION_TYPE_CHAIN
    }
|   AND CHAIN
    {
        $$ = tree.COMPLETION_TYPE_CHAIN
    }
|   AND NO CHAIN RELEASE
    {
        $$ = tree.COMPLETION_TYPE_RELEASE
    }
|   RELEASE
    {
        $$ = tree.COMPLETION_TYPE_RELEASE
    }
|   AND NO CHAIN NO RELEASE
    {
        $$ = tree.COMPLETION_TYPE_NO_CHAIN
    }
|   AND NO CHAIN
    {
        $$ = tree.COMPLETION_TYPE_NO_CHAIN
    }
|   NO RELEASE
    {
        $$ = tree.COMPLETION_TYPE_NO_CHAIN
    }

begin_stmt:
    BEGIN
    {
        $$ = tree.NewBeginTransaction(yylex.(*Lexer).buf)
    }
|   BEGIN WORK
    {
        $$ = tree.NewBeginTransaction(yylex.(*Lexer).buf)
    }
|   START TRANSACTION
    {
        $$ = tree.NewBeginTransaction(yylex.(*Lexer).buf)
    }
|   START TRANSACTION READ WRITE
    {
        $$ = tree.NewBeginTransactionWithMode(tree.READ_WRITE_MODE_READ_WRITE, yylex.(*Lexer).buf)
    }
|   START TRANSACTION READ ONLY
    {
        $$ = tree.NewBeginTransactionWithMode(tree.READ_WRITE_MODE_READ_ONLY, yylex.(*Lexer).buf)
    }
|   START TRANSACTION WITH CONSISTENT SNAPSHOT
    {
        $$ = tree.NewBeginTransaction(yylex.(*Lexer).buf)
    }

use_stmt:
    USE ident
    {
        // SecondaryRole Name SecondaryRoleType 
        $$ = tree.NewUse(
            false,
            $2,
            0,
            yylex.(*Lexer).buf,
        )
    }
|   USE
    {
        // SecondaryRole Name SecondaryRoleType 
        $$ = tree.NewUse( 
            false,
            nil,
            0,
            yylex.(*Lexer).buf,
        )
    }
|   USE ROLE role_spec
    {
        // SecondaryRole Name SecondaryRoleType 
        $$ = tree.NewUse( 
            false,
            nil,
            0,
            yylex.(*Lexer).buf,
        )
        $$.(*tree.Use).Role = $3
    }
|   USE SECONDARY ROLE ALL
    {
        // SecondaryRole Name SecondaryRoleType 
        $$ = tree.NewUse( 
            true,
            nil,
            tree.SecondaryRoleTypeAll,
            yylex.(*Lexer).buf,
        )
    }
|   USE SECONDARY ROLE NONE
    {
        // SecondaryRole Name SecondaryRoleType 
        $$ = tree.NewUse( 
            true,
            nil,
            tree.SecondaryRoleTypeNone,
            yylex.(*Lexer).buf,
        )
    }

update_stmt:
    update_no_with_stmt
|    with_clause update_no_with_stmt
    {
        $2.(*tree.Update).With = $1
        $$ = $2
    }

update_no_with_stmt:
    UPDATE priority_opt ignore_opt table_reference SET update_list where_expression_opt order_by_opt limit_opt
    {
        // Single-table syntax
        // Tables UpdateExprs Where OrderBy Limit
        exprs := buffer.MakeSlice[tree.TableExpr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.TableExpr](yylex.(*Lexer).buf, exprs, $4)
        $$ = tree.NewUpdate( 
            exprs,
            $6,
            $7,
            $8,
            $9,
            yylex.(*Lexer).buf,
        )
    }
|    UPDATE priority_opt ignore_opt table_references SET update_list where_expression_opt
    {
        // Multiple-table syntax
        // Tables UpdateExprs Where OrderBy Limit
        exprs := buffer.MakeSlice[tree.TableExpr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.TableExpr](yylex.(*Lexer).buf, exprs, $4)
        $$ = tree.NewUpdate(
            exprs,
            $6,
            $7,
            nil,
            nil,
            yylex.(*Lexer).buf,
        )
    }

update_list:
    update_value
    {
        $$ = buffer.MakeSlice[*tree.UpdateExpr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.UpdateExpr](yylex.(*Lexer).buf, $$, $1)
    }
|   update_list ',' update_value
    {
        $$ = buffer.MakeSlice[*tree.UpdateExpr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.UpdateExpr](yylex.(*Lexer).buf, $1, $3)
    }

update_value:
    column_name '=' expr_or_default
    {
        names := buffer.MakeSlice[*tree.UnresolvedName](yylex.(*Lexer).buf)
        names = buffer.AppendSlice[*tree.UnresolvedName](yylex.(*Lexer).buf, names, $1)
        $$ = tree.NewUpdateExpr(names, $3, yylex.(*Lexer).buf)
    }

lock_stmt:
    lock_table_stmt
|   unlock_table_stmt

lock_table_stmt:
    LOCK TABLES table_lock_list
    {
        // TableLocks
        $$ = tree.NewLockTableStmt($3, yylex.(*Lexer).buf)
    }

table_lock_list:
    table_lock_elem
    {
        $$ = buffer.MakeSlice[tree.TableLock](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.TableLock](yylex.(*Lexer).buf, $$, $1)
    }
|   table_lock_list ',' table_lock_elem
    {
        $$ = buffer.MakeSlice[tree.TableLock](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.TableLock](yylex.(*Lexer).buf, $1, $3)
    }

table_lock_elem:
    table_name table_lock_type
    {
        $$ = *tree.NewTableLock(*$1, $2, yylex.(*Lexer).buf)
    }

table_lock_type:
    READ
    {
        $$ = tree.TableLockRead
    }
|   READ LOCAL
    {
        $$ = tree.TableLockReadLocal
    }
|   WRITE
    {
        $$ = tree.TableLockWrite
    }
|   LOW_PRIORITY WRITE
    {
        $$ = tree.TableLockLowPriorityWrite
    }

unlock_table_stmt:
    UNLOCK TABLES
    {
       $$ = tree.NewUnLockTableStmt(yylex.(*Lexer).buf)
    }

prepareable_stmt:
    create_stmt
|   insert_stmt
|   delete_stmt
|   drop_stmt
|   show_stmt
|   update_stmt
|   select_stmt
    {
        $$ = $1
    }

prepare_stmt:
    prepare_sym stmt_name FROM prepareable_stmt
    {
        // Name Stmt
        $$ = tree.NewPrepareStmt(
            tree.Identifier($2),
            $4, 
            yylex.(*Lexer).buf,
        )
    }
|   prepare_sym stmt_name FROM STRING
    {
        // Name Sql
        $$ = tree.NewPrepareString(
            tree.Identifier($2), 
            $4, 
            yylex.(*Lexer).buf,
        )
    }

execute_stmt:
    execute_sym stmt_name
    {
        // Name        
        $$ = tree.NewExecute(
            tree.Identifier($2), 
            yylex.(*Lexer).buf,
        )
    }
|   execute_sym stmt_name USING variable_list
    {
        // Name Variables
        $$ = tree.NewExecuteWithVariables(
            tree.Identifier($2), 
            $4, 
            yylex.(*Lexer).buf,
        )
    }

deallocate_stmt:
    deallocate_sym PREPARE stmt_name
    {
        $$ = tree.NewDeallocate(tree.Identifier($3), false, yylex.(*Lexer).buf)
    }

reset_stmt:
    reset_sym PREPARE stmt_name
    {
        $$ = tree.NewReset(tree.Identifier($3), yylex.(*Lexer).buf)
    }

explainable_stmt:
    delete_stmt
|   load_data_stmt
|   insert_stmt
|   replace_stmt
|   update_stmt
|   select_stmt
    {
        $$ = $1
    }

explain_stmt:
    explain_sym unresolved_object_name
    {
        // Table ColName
        $$ = tree.NewShowColumns1(
            $2,
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   explain_sym unresolved_object_name column_name
    {
        // Table ColName
        $$ = tree.NewShowColumns1(
            $2,
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   explain_sym FOR CONNECTION INTEGRAL
    {
        $$ = tree.NewExplainFor("", uint64($4.(int64)), yylex.(*Lexer).buf)
    }
|   explain_sym FORMAT '=' STRING FOR CONNECTION INTEGRAL
    {
        $$ = tree.NewExplainFor($4, uint64($7.(int64)), yylex.(*Lexer).buf)
    }
|   explain_sym explainable_stmt
    {
        $$ = tree.NewExplainStmt($2, "text", yylex.(*Lexer).buf)
    }
|   explain_sym VERBOSE explainable_stmt
    {
        explainStmt := tree.NewExplainStmt($3, "text", yylex.(*Lexer).buf)
        optionElem := tree.MakeOptionElem("verbose", "NULL", yylex.(*Lexer).buf)
        options := tree.MakeOptions(optionElem, yylex.(*Lexer).buf)
        explainStmt.Options = options
        $$ = explainStmt
    }
|   explain_sym ANALYZE explainable_stmt
    {
        explainStmt := tree.NewExplainAnalyze($3, "text", yylex.(*Lexer).buf)
        optionElem := tree.MakeOptionElem("analyze", "NULL", yylex.(*Lexer).buf)
        options := tree.MakeOptions(optionElem, yylex.(*Lexer).buf)
        explainStmt.Options = options
        $$ = explainStmt
    }
|   explain_sym ANALYZE VERBOSE explainable_stmt
    {
        explainStmt := tree.NewExplainAnalyze($4, "text", yylex.(*Lexer).buf)
        optionElem1 := tree.MakeOptionElem("analyze", "NULL", yylex.(*Lexer).buf)
        optionElem2 := tree.MakeOptionElem("verbose", "NULL", yylex.(*Lexer).buf)
        options := tree.MakeOptions(optionElem1, yylex.(*Lexer).buf)
        options = buffer.AppendSlice[tree.OptionElem](yylex.(*Lexer).buf, options, optionElem2)
        explainStmt.Options = options
        $$ = explainStmt
    }
|   explain_sym '(' utility_option_list ')' explainable_stmt
    {
        if tree.IsContainAnalyze($3) {
            explainStmt := tree.NewExplainAnalyze($5, "text", yylex.(*Lexer).buf)
            explainStmt.Options = $3
            $$ = explainStmt
        } else {
            explainStmt := tree.NewExplainStmt($5, "text", yylex.(*Lexer).buf)
            explainStmt.Options = $3
            $$ = explainStmt
        }
    }

explain_option_key:
    ANALYZE
|   VERBOSE
|   FORMAT

explain_foramt_value:
    JSON
|   TEXT


prepare_sym:
    PREPARE

deallocate_sym:
    DEALLOCATE

execute_sym:
    EXECUTE

reset_sym:
    RESET

explain_sym:
    EXPLAIN
|   DESCRIBE
|   DESC

utility_option_list:
    utility_option_elem
    {
        $$ = tree.MakeOptions($1, yylex.(*Lexer).buf)
    }
|     utility_option_list ',' utility_option_elem
    {
        $$ = buffer.MakeSlice[tree.OptionElem](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.OptionElem](yylex.(*Lexer).buf, $1, $3);
    }

utility_option_elem:
    utility_option_name utility_option_arg
    {
        $$ = tree.MakeOptionElem($1, $2, yylex.(*Lexer).buf)
    }

utility_option_name:
    explain_option_key
    {
        $$ = $1
    }

utility_option_arg:
    TRUE                    { $$ = "true" }
|   FALSE                        { $$ = "false" }
|   explain_foramt_value                    { $$ = $1 }


analyze_stmt:
    ANALYZE TABLE table_name '(' column_list ')'
    {
        $$ = tree.NewAnalyzeStmt($3, $5, yylex.(*Lexer).buf)
    }

alter_stmt:
    alter_user_stmt
|   alter_account_stmt
|   alter_database_config_stmt
|   alter_view_stmt
|   alter_table_stmt
|   alter_publication_stmt
|   alter_stage_stmt
|   alter_sequence_stmt
// |    alter_ddl_stmt

alter_sequence_stmt:
    ALTER SEQUENCE exists_opt table_name alter_as_datatype_opt increment_by_opt min_value_opt max_value_opt start_with_opt alter_cycle_opt
    {
        // IFExists Name Type IncrementBy MinValue MaxValue StartWith MinValue Cycle
        $$ = tree.NewAlterSequence(
            $3,
            $4,
            $5,
            $6,
            $7,
            $8,
            $9,
            $10,
            yylex.(*Lexer).buf,
        )
    }


alter_view_stmt:
    ALTER VIEW exists_opt table_name column_list_opt AS select_stmt
    {
        // IfExists Name ColNames AsSource
        $$ = tree.NewAlterView(
            $3,
            $4,
            $5,
            $7,
            yylex.(*Lexer).buf,
        )
    }

alter_table_stmt:
    ALTER TABLE table_name alter_option_list
    {
        // TableName AlterTableOptions
        $$ = tree.NewAlterTable(
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }

alter_option_list:
    alter_option
    {
        
        $$ = buffer.MakeSlice[tree.AlterTableOption](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.AlterTableOption](yylex.(*Lexer).buf, $$, $1)
    }
|   alter_option_list ',' alter_option
    {
        $$ = buffer.MakeSlice[tree.AlterTableOption](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.AlterTableOption](yylex.(*Lexer).buf,  $1, $3)
    }

alter_option:
    ADD table_elem_2
    {
        // TableDef
        opt := tree.NewAlterOptionAdd(
            $2,
            yylex.(*Lexer).buf,
        )
        aopt := *buffer.Alloc[tree.AlterTableOption](yylex.(*Lexer).buf)
        aopt = opt
        $$ = aopt
    }
|    MODIFY column_keyword_opt column_def column_position
    {
    	opt := tree.NewAlterTableModifyColumnClause(
                tree.AlterTableModifyColumn,
                $3,
                $4,
                yylex.(*Lexer).buf,
	        )
        aopt := *buffer.Alloc[tree.AlterTableOption](yylex.(*Lexer).buf)
        aopt = opt
        $$ = aopt
    }
|   CHANGE column_keyword_opt column_name column_def column_position
    {
        // Type OldColumnName NewColumn Position
    	opt := tree.NewAlterTableChangeColumnClause(
		    tree.AlterTableChangeColumn,
		    $3,
		    $4,
		    $5,
            yylex.(*Lexer).buf,
	    )
        aopt := *buffer.Alloc[tree.AlterTableOption](yylex.(*Lexer).buf)
        aopt = opt
        $$ = aopt
    }
|  RENAME COLUMN column_name TO column_name
    {
        // Typ OldColumnName NewColumnName
    	opt := tree.NewAlterTableRenameColumnClause(
    		tree.AlterTableRenameColumn,
    		$3,
		    $5,
            yylex.(*Lexer).buf,
    	)
        aopt := *buffer.Alloc[tree.AlterTableOption](yylex.(*Lexer).buf)
        aopt = opt
        $$ = aopt
    }
|  ALTER column_keyword_opt column_name SET DEFAULT bit_expr
    {
        // Typ ColumnName DefaultExpr Visibility OptionType
    	opt := tree.NewAlterTableAlterColumnClause(
		    tree.AlterTableAlterColumn,
		    $3,
		    tree.NewAttributeDefault($6, yylex.(*Lexer).buf),
            0,
		    tree.AlterColumnOptionSetDefault,
            yylex.(*Lexer).buf,
	    )
        aopt := *buffer.Alloc[tree.AlterTableOption](yylex.(*Lexer).buf)
        aopt = opt
        $$ = aopt
    }
|  ALTER column_keyword_opt column_name SET visibility
    {
        // Typ ColumnName DefaultExpr Visibility OptionType
    	opt := tree.NewAlterTableAlterColumnClause(
		    tree.AlterTableAlterColumn,
		    $3,
            nil,
		    $5,
		    tree.AlterColumnOptionSetVisibility,
            yylex.(*Lexer).buf,
	    )
        aopt := *buffer.Alloc[tree.AlterTableOption](yylex.(*Lexer).buf)
        aopt = opt
        $$ = aopt
    }
|  ALTER column_keyword_opt column_name DROP DEFAULT
    {
        // Typ ColumnName DefaultExpr Visibility OptionType
    	opt := tree.NewAlterTableAlterColumnClause(
		    tree.AlterTableAlterColumn,
		    $3,
            nil,
            0,
		    tree.AlterColumnOptionDropDefault,
            yylex.(*Lexer).buf,
        )
        aopt := *buffer.Alloc[tree.AlterTableOption](yylex.(*Lexer).buf)
        aopt = opt
        $$ = aopt
    }
|  ORDER BY alter_column_order_list %prec LOWER_THAN_ORDER
    {
    	opt := tree.NewAlterTableOrderByColumnClause(
                tree.AlterTableOrderByColumn,
        	    $3,
                yylex.(*Lexer).buf,
            )
        aopt := *buffer.Alloc[tree.AlterTableOption](yylex.(*Lexer).buf)
        aopt = opt
        $$ = aopt
    }
|   DROP alter_table_drop
    {
        opt := *buffer.Alloc[tree.AlterTableOption](yylex.(*Lexer).buf)
        opt = $2
        $$ = opt
    }
|   ALTER alter_table_alter
    {
        opt := *buffer.Alloc[tree.AlterTableOption](yylex.(*Lexer).buf)
        opt = $2
        $$ = opt
    }
|   table_option
    {
        opt := *buffer.Alloc[tree.AlterTableOption](yylex.(*Lexer).buf)
        opt = $1
        $$ = opt
    }
|   RENAME rename_type alter_table_rename
    {
        opt := *buffer.Alloc[tree.AlterTableOption](yylex.(*Lexer).buf)
        opt = $3
        $$ = opt
    }
|   ADD column_keyword_opt column_def column_position
    {
        $$ = tree.AlterTableOption(
            // Column Position
            tree.NewAlterAddCol(
                $3,
                $4,
                yylex.(*Lexer).buf,
            ),
        )
    }
|   ALGORITHM equal_opt algorithm_type
    {
        // Type Enforce
        $$ = tree.NewAlterOptionAlterCheck(
            $1,
            false,
            yylex.(*Lexer).buf,
        )
    }
|   default_opt charset_keyword equal_opt charset_name COLLATE equal_opt charset_name
    {
        $$ = tree.NewTableOptionCharset($4, yylex.(*Lexer).buf)
    }
|   CONVERT TO CHARACTER SET charset_name
    {
        $$ = tree.NewTableOptionCharset($5, yylex.(*Lexer).buf)
    }
|   CONVERT TO CHARACTER SET charset_name COLLATE equal_opt charset_name
    {
        $$ = tree.NewTableOptionCharset($5, yylex.(*Lexer).buf)
    }
|   able_type KEYS
    {
        $$ = tree.NewTableOptionCharset($1, yylex.(*Lexer).buf)
    }
|   space_type TABLESPACE
    {
        $$ = tree.NewTableOptionCharset($1, yylex.(*Lexer).buf)
    }
|   FORCE
    {
        $$ = tree.NewTableOptionCharset($1, yylex.(*Lexer).buf)
    }
|   LOCK equal_opt lock_type
    {
        $$ = tree.NewTableOptionCharset($1, yylex.(*Lexer).buf)
    }
|   with_type VALIDATION
    {
        $$ = tree.NewTableOptionCharset($1, yylex.(*Lexer).buf)
    }

rename_type:
    {
        $$ = ""
    }
|   TO
|   AS

algorithm_type:
    DEFAULT
|   INSTANT
|   INPLACE
|   COPY

able_type:
    DISABLE
|   ENABLE

space_type:
    DISCARD
|   IMPORT

lock_type:
    DEFAULT
|   NONE
|   SHARED
|   EXCLUSIVE

with_type:
    WITHOUT
|   WITH

column_keyword_opt:
    {
	$$ = ""
    }
|   COLUMN
    {
        $$ = string("COLUMN")
    }

column_position:
    {
	$$ = tree.NewColumnPosition(
            tree.ColumnPositionNone,
            nil,
            yylex.(*Lexer).buf,
	    )
    }
|   FIRST
    {
	$$ = tree.NewColumnPosition(
            tree.ColumnPositionFirst,
            nil,
            yylex.(*Lexer).buf,
	    )
    }
|   AFTER column_name
    {
	$$ = tree.NewColumnPosition(
            tree.ColumnPositionAfter,
            $2,
            yylex.(*Lexer).buf,
        )
    }

alter_column_order_list:
     alter_column_order
    {
        $$ = buffer.MakeSlice[*tree.AlterColumnOrder](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.AlterColumnOrder](yylex.(*Lexer).buf, $$, $1)
    }
|   alter_column_order_list ',' alter_column_order
    {
        $$ = buffer.MakeSlice[*tree.AlterColumnOrder](yylex.(*Lexer).buf)
	    $$ = buffer.AppendSlice[*tree.AlterColumnOrder](yylex.(*Lexer).buf, $1, $3)
    }

alter_column_order:
    column_name asc_desc_opt
    {
        // Column Direction
	    $$ = tree.NewAlterColumnOrder(
            $1, 
            $2,
            yylex.(*Lexer).buf,
        )
    }

alter_table_rename:
    table_name_unresolved
    {
        // Name
        $$ = tree.NewAlterTableName( 
            $1,
            yylex.(*Lexer).buf,
        )
    }

alter_table_drop:
    INDEX ident
    {
        // Type Name
        $$ = tree.NewAlterOptionDrop( 
            tree.AlterTableDropIndex,
            tree.Identifier($2.Compare()),
            yylex.(*Lexer).buf,
        )
    }
|   KEY ident
    {
        // Type Name
        $$ = tree.NewAlterOptionDrop( 
            tree.AlterTableDropKey,
            tree.Identifier($2.Compare()),
            yylex.(*Lexer).buf,
        )
    }
|   ident
    {
        // Type Name
        $$ = tree.NewAlterOptionDrop( 
            tree.AlterTableDropColumn,
            tree.Identifier($1.Compare()),
            yylex.(*Lexer).buf,
        )
    }
|   COLUMN ident
    {
        // Type Name
        $$ = tree.NewAlterOptionDrop( 
            tree.AlterTableDropColumn,
            tree.Identifier($2.Compare()),
            yylex.(*Lexer).buf,
        )
    }
|   FOREIGN KEY ident
    {
        // Type Name
        $$ = tree.NewAlterOptionDrop( 
            tree.AlterTableDropForeignKey,
            tree.Identifier($3.Compare()),
            yylex.(*Lexer).buf,
        )
    }
|   PRIMARY KEY
    {
        // Type Name
        $$ = tree.NewAlterOptionDrop( 
            tree.AlterTableDropPrimaryKey,
            "",
            yylex.(*Lexer).buf,
        )
    }

alter_table_alter:
    INDEX ident visibility
    {
        // Type Enforce
        $$ = tree.NewAlterOptionAlterIndex( 
            tree.Identifier($2.Compare()),
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   CHECK ident enforce
    {
        // Type Enforce
        $$ = tree.NewAlterOptionAlterCheck(
            $1,
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   CONSTRAINT ident enforce
    {
        // Type Enforce
        $$ = tree.NewAlterOptionAlterCheck(
            $1,
            $3,
            yylex.(*Lexer).buf,
        )
    }

visibility:
    VISIBLE
    {
        $$ = tree.VISIBLE_TYPE_VISIBLE
    }
|   INVISIBLE
    {
   	    $$ = tree.VISIBLE_TYPE_INVISIBLE
    }


alter_account_stmt:
    ALTER ACCOUNT exists_opt account_name alter_account_auth_option account_status_option account_comment_opt
    {
        // IfExists Name AuthOption StatusOption Comment
        $$ = tree.NewAlterAccount( 
            $3,
            $4,
            $5,
            $6,
            $7,
            yylex.(*Lexer).buf,
        )
    }

alter_database_config_stmt:
    ALTER DATABASE db_name SET MYSQL_COMPATIBILITY_MODE '=' STRING
    {
        $$ = tree.NewAlterDataBaseConfig(
            "",
            $3,
            false,
            $7,
            yylex.(*Lexer).buf,
        )
    }
|   ALTER ACCOUNT CONFIG account_name SET MYSQL_COMPATIBILITY_MODE '=' STRING
    {
        $$ = tree.NewAlterDataBaseConfig(
            $4,
            "",
            true,
            $8,
            yylex.(*Lexer).buf,
        )
    }
|   ALTER ACCOUNT CONFIG SET MYSQL_COMPATIBILITY_MODE  var_name equal_or_assignment set_expr
    {
        // System Global Name Value Reserved
        v := tree.NewVarAssignmentExpr(
            true, 
            true,
            $6,
            $8,
            nil,
            yylex.(*Lexer).buf,
        )
        assignments := buffer.MakeSlice[*tree.VarAssignmentExpr](yylex.(*Lexer).buf)
        assignments = buffer.AppendSlice[*tree.VarAssignmentExpr](yylex.(*Lexer).buf, assignments, v)

        // Assignments
        $$ = tree.NewSetVar(
            assignments,
            yylex.(*Lexer).buf,
        )
    }

alter_account_auth_option:
{
    // Exist Equal AdminName IdentifiedType
    $$ = *tree.NewAlterAccountAuthOption(
        false,
        "",
        "",
        yylex.(*Lexer).buf,
    )
}
| ADMIN_NAME equal_opt account_admin_name account_identified
{
    // Exist Equal AdminName IdentifiedType
    $$ = *tree.NewAlterAccountAuthOption(
        true,
        $2,
        $3,
        yylex.(*Lexer).buf,
    )
    $$.IdentifiedType = $4
}

alter_user_stmt:
    ALTER USER exists_opt user_spec_list_of_create_user default_role_opt pwd_or_lck_opt user_comment_or_attribute_opt
    {
        // IfExists Users Role MiscOpt CommentOrAttribute
        $$ = tree.NewAlterUser( 
            $3,
            $4,
            $5,
            $6,
            $7,
            yylex.(*Lexer).buf,
        )
    }

default_role_opt:
    {
        $$ = nil
    }
|   DEFAULT ROLE account_role_name
    {
        $$ = tree.NewRole(
            $3,
            yylex.(*Lexer).buf,
        )
    }

exists_opt:
    {
        $$ = false
    }
|   IF EXISTS
    {
        $$ = true
    }

pwd_or_lck_opt:
    {
        $$ = nil
    }
|   pwd_or_lck
    {
        $$ = $1
    }

//pwd_or_lck_list:
//    pwd_or_lck
//    {
//        $$ = []tree.UserMiscOption{$1}
//    }
//|   pwd_or_lck_list pwd_or_lck
//    {
//        $$ = buffer.AppendSlice[tree.Statement](yylex.(*Lexer).buf, $1, $2)
//    }

pwd_or_lck:
    UNLOCK
    {
        $$ = tree.NewUserMiscOptionAccountUnlock(yylex.(*Lexer).buf)
    }
|   LOCK
    {
        $$ = tree.NewUserMiscOptionAccountLock(yylex.(*Lexer).buf)
    }
|   pwd_expire
    {
        $$ = tree.NewUserMiscOptionPasswordExpireNone(yylex.(*Lexer).buf)
    }
|   pwd_expire INTERVAL INTEGRAL DAY
    {
        $$ = tree.NewUserMiscOptionPasswordExpireInterval($3.(int64), yylex.(*Lexer).buf)
    }
|   pwd_expire NEVER
    {
        $$ = tree.NewUserMiscOptionPasswordExpireNever(yylex.(*Lexer).buf)
    }
|   pwd_expire DEFAULT
    {
        $$ = tree.NewUserMiscOptionPasswordExpireDefault(yylex.(*Lexer).buf)
    }
|   PASSWORD HISTORY DEFAULT
    {
        $$ = tree.NewUserMiscOptionPasswordHistoryDefault(yylex.(*Lexer).buf)
    }
|   PASSWORD HISTORY INTEGRAL
    {
        $$ = tree.NewUserMiscOptionPasswordHistoryCount($3.(int64), yylex.(*Lexer).buf)
    }
|   PASSWORD REUSE INTERVAL DEFAULT
    {
        $$ = tree.NewUserMiscOptionPasswordReuseIntervalDefault(yylex.(*Lexer).buf)
    }
|   PASSWORD REUSE INTERVAL INTEGRAL DAY
    {
        $$ = tree.NewUserMiscOptionPasswordReuseIntervalCount($4.(int64), yylex.(*Lexer).buf)
    }
|   PASSWORD REQUIRE CURRENT
    {
        $$ = tree.NewUserMiscOptionPasswordRequireCurrentNone(yylex.(*Lexer).buf)
    }
|   PASSWORD REQUIRE CURRENT DEFAULT
    {
        $$ = tree.NewUserMiscOptionPasswordRequireCurrentDefault(yylex.(*Lexer).buf)
    }
|   PASSWORD REQUIRE CURRENT OPTIONAL
    {
        $$ = tree.NewUserMiscOptionPasswordRequireCurrentOptional(yylex.(*Lexer).buf)
    }
|   FAILED_LOGIN_ATTEMPTS INTEGRAL
    {
        $$ = tree.NewUserMiscOptionFailedLoginAttempts($2.(int64), yylex.(*Lexer).buf)
    }
|   PASSWORD_LOCK_TIME INTEGRAL
    {
        $$ = tree.NewUserMiscOptionPasswordLockTimeCount($2.(int64), yylex.(*Lexer).buf)
    }
|   PASSWORD_LOCK_TIME UNBOUNDED
    {
        $$ = tree.NewUserMiscOptionPasswordLockTimeUnbounded(yylex.(*Lexer).buf)
    }

pwd_expire:
    PASSWORD EXPIRE clear_pwd_opt
    {
        $$ = nil
    }

clear_pwd_opt:
    {
        $$ = nil
    }

auth_string:
    STRING

show_stmt:
    show_create_stmt
|   show_columns_stmt
|   show_databases_stmt
|   show_tables_stmt
|   show_sequences_stmt
|   show_process_stmt
|   show_errors_stmt
|   show_warnings_stmt
|   show_variables_stmt
|   show_status_stmt
|   show_index_stmt
|   show_target_filter_stmt
|   show_table_status_stmt
|   show_grants_stmt
|   show_roles_stmt
|   show_collation_stmt
|   show_function_status_stmt
|   show_procedure_status_stmt
|   show_node_list_stmt
|   show_locks_stmt
|   show_table_num_stmt
|   show_column_num_stmt
|   show_table_values_stmt
|   show_table_size_stmt
|   show_accounts_stmt
|   show_publications_stmt
|   show_subscriptions_stmt
|   show_servers_stmt
|   show_stages_stmt
|   show_connectors_stmt

show_collation_stmt:
    SHOW COLLATION like_opt where_expression_opt
    {
        $$ = tree.NewShowCollation(yylex.(*Lexer).buf)
    }

show_stages_stmt:
    SHOW STAGES like_opt
    {
        // like
        $$ = tree.NewShowStages(
            $3,
            yylex.(*Lexer).buf,
        )
    }

show_grants_stmt:
    SHOW GRANTS
    {
        // Username HostName Roles ShowGrantType
        $$ = tree.NewShowGrants(
            "",
            "",
            nil,
            tree.GrantForUser,
            yylex.(*Lexer).buf,
        )
    }
|   SHOW GRANTS FOR user_name using_roles_opt
    {
        // Username HostName Roles ShowGrantType
        $$ = tree.NewShowGrants(
            $4.Username.Get(),
            $4.Hostname.Get(),
            $5,
            tree.GrantForUser,
            yylex.(*Lexer).buf,
        )
    }
|   SHOW GRANTS FOR ROLE ident
    {
        roles := buffer.MakeSlice[*tree.Role](yylex.(*Lexer).buf)
        roles = buffer.AppendSlice[*tree.Role](yylex.(*Lexer).buf, roles, tree.NewRole($5.Compare(), yylex.(*Lexer).buf))
        // Username HostName Roles ShowGrantType
        $$ = tree.NewShowGrants(
            "",
            "",
            roles,
            tree.GrantForRole,
            yylex.(*Lexer).buf,
        )
    }

using_roles_opt:
    {
        $$ = nil
    }
|    USING role_spec_list
    {
        $$ = $2
    }

show_table_status_stmt:
    SHOW TABLE STATUS from_or_in_opt db_name_opt like_opt where_expression_opt
    {
        // DbName Like Where
        $$ = tree.NewShowTableStatus(
            $5, 
            $6, 
            $7,
            yylex.(*Lexer).buf,
        )
    }

from_or_in_opt:
    {}
|    from_or_in

db_name_opt:
    {}
|    db_name

show_function_status_stmt:
    SHOW FUNCTION STATUS like_opt where_expression_opt
    {
       // Like Where IsFunction
       $$ = tree.NewShowFunctionOrProcedureStatus(
            $4,
            $5,
            true,
            yylex.(*Lexer).buf,
        )
    }

show_procedure_status_stmt:
    SHOW PROCEDURE STATUS like_opt where_expression_opt
    {
       // Like Where IsFunction
        $$ = tree.NewShowFunctionOrProcedureStatus(
            $4,
            $5,
            false,
            yylex.(*Lexer).buf,
        )
    }

show_roles_stmt:
    SHOW ROLES like_opt
    {
        // Like
        $$ = tree.NewShowRolesStmt(
            $3,
            yylex.(*Lexer).buf,
        )
    }

show_node_list_stmt:
    SHOW NODE LIST
    {
       $$ = tree.NewShowNodeList(yylex.(*Lexer).buf)
    }

show_locks_stmt:
    SHOW LOCKS
    {
       $$ = tree.NewShowLocks(yylex.(*Lexer).buf)
    }

show_table_num_stmt:
    SHOW TABLE_NUMBER from_or_in_opt db_name_opt
    {
      // DbName
      $$ = tree.NewShowTableNumber($4 , yylex.(*Lexer).buf)
    }

show_column_num_stmt:
    SHOW COLUMN_NUMBER table_column_name database_name_opt
    {
       // Table DbName
       $$ = tree.NewShowColumnNumber(
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }

show_table_values_stmt:
    SHOW TABLE_VALUES table_column_name database_name_opt
    {
       // Table DbName
       $$ = tree.NewShowTableValues(
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }

show_table_size_stmt:
    SHOW TABLE_SIZE table_column_name database_name_opt
    {
       // Table DbName
       $$ = tree.NewShowTableSize(
            $3, 
            $4,
            yylex.(*Lexer).buf,
        )
    }

show_target_filter_stmt:
    SHOW show_target like_opt where_expression_opt
    {
        s := $2.(*tree.ShowTarget)
        s.Like = $3
        s.Where = $4
        $$ = s
    }

show_target:
    CONFIG
    {
        // DbName, ShowType
        $$ = tree.NewShowTarget(
            "",
            tree.ShowConfig,
            yylex.(*Lexer).buf,
        )
    }
|    charset_keyword
    {
        // DbName, ShowType
        $$ = tree.NewShowTarget(
            "",
            tree.ShowCharset,
            yylex.(*Lexer).buf,
        )
    }
|    ENGINES
    {
        // DbName, ShowType
        $$ = tree.NewShowTarget(
            "",
            tree.ShowEngines,
            yylex.(*Lexer).buf,
        )
    }
|    TRIGGERS from_or_in_opt db_name_opt
    {
        // DbName, ShowType
        $$ = tree.NewShowTarget(
            $3,
            tree.ShowTriggers,
            yylex.(*Lexer).buf,
        )
    }
|    EVENTS from_or_in_opt db_name_opt
    {
        // DbName, ShowType
        $$ = tree.NewShowTarget(
            $3,
            tree.ShowEvents,
            yylex.(*Lexer).buf,
         )
    }
|    PLUGINS
    {
        // DbName, ShowType
        $$ = tree.NewShowTarget(
            "",
            tree.ShowPlugins,
            yylex.(*Lexer).buf,
        )
    }
|    PRIVILEGES
    {
        // DbName, ShowType
        $$ = tree.NewShowTarget(
            "",
            tree.ShowPrivileges,
            yylex.(*Lexer).buf,
        )
    }
|    PROFILES
    {
        // DbName, ShowType
        $$ = tree.NewShowTarget(
            "",
            tree.ShowProfiles,
            yylex.(*Lexer).buf,
        )
    }

show_index_stmt:
    SHOW extended_opt index_kwd from_or_in table_name where_expression_opt
    {
        // TableName Where
        $$ = tree.NewShowIndex(
            $5,
            $6,
            yylex.(*Lexer).buf,
        )
    }
|	SHOW extended_opt index_kwd from_or_in ident from_or_in ident where_expression_opt
    {
        prefix := tree.NewObjectNamePrefix(yylex.(*Lexer).buf)
        bi := tree.NewBufIdentifier($7.Compare())
        yylex.(*Lexer).buf.Pin(bi)
        prefix.SchemaName = bi
        prefix.ExplicitSchema = true
        
        tbl := tree.NewTableName(tree.Identifier($5.Compare()), *prefix, yylex.(*Lexer).buf)
        $$ = tree.NewShowIndex(
            tbl,
            $8,
            yylex.(*Lexer).buf,
        )
    }

extended_opt:
    {}
|    EXTENDED
    {}

index_kwd:
    INDEX
|   INDEXES
|   KEYS

show_variables_stmt:
    SHOW global_scope VARIABLES like_opt where_expression_opt
    {
        // Global Like Where
        $$ = tree.NewShowVariables( 
            $2,
            $4,
            $5,
            yylex.(*Lexer).buf,
        )
    }

show_status_stmt:
    SHOW global_scope STATUS like_opt where_expression_opt
    {
        // Global Like Where
        $$ = tree.NewShowStatus( 
            $2,
            $4,
            $5,
            yylex.(*Lexer).buf,
        )
    }

global_scope:
    {
        $$ = false
    }
|   GLOBAL
    {
        $$ = true
    }
|   SESSION
    {
        $$ = false
    }

show_warnings_stmt:
    SHOW WARNINGS limit_opt
    {
        $$ = tree.NewShowWarnings(yylex.(*Lexer).buf)
    }

show_errors_stmt:
    SHOW ERRORS limit_opt
    {
        $$ = tree.NewShowErrors(yylex.(*Lexer).buf)
    }

show_process_stmt:
    SHOW full_opt PROCESSLIST
    {
        // Full
        $$ = tree.NewShowProcessList(
            $2,
            yylex.(*Lexer).buf,
        )
    }

show_sequences_stmt:
    SHOW SEQUENCES database_name_opt where_expression_opt
    {
        $$ = tree.NewShowSequences( 
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }

show_tables_stmt:
    SHOW full_opt TABLES database_name_opt like_opt where_expression_opt
    {
        // Open Full DBName Like Where
        $$ = tree.NewShowTables( 
            false,
            $2,
            $4,
            $5,
            $6,
            yylex.(*Lexer).buf,
        )
    }
|   SHOW OPEN full_opt TABLES database_name_opt like_opt where_expression_opt
    {
        // Open Full DBName Like Where
        $$ = tree.NewShowTables( 
            true,
            $3,
            $5,
            $6,
            $7,
            yylex.(*Lexer).buf,
        )
    }

show_databases_stmt:
    SHOW DATABASES like_opt where_expression_opt
    {
        // Like Where
        $$ = tree.NewShowDatabases(
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }
|   SHOW SCHEMAS like_opt where_expression_opt
    {
        // Like Where
        $$ = tree.NewShowDatabases(
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }

show_columns_stmt:
    SHOW full_opt fields_or_columns table_column_name database_name_opt like_opt where_expression_opt
    {
        // Ext Full Table DBName Like Where
        $$ = tree.NewShowColumns2(
            false,
            $2,
            $4,
            $5,
            $6,
            $7,
            yylex.(*Lexer).buf,
        )
    }
|   SHOW EXTENDED full_opt fields_or_columns table_column_name database_name_opt like_opt where_expression_opt
    {
        // Ext Full Table DBName Like Where
        $$ = tree.NewShowColumns2(
            true,
            $3,
            $5,
            $6,
            $7,
            $8,
            yylex.(*Lexer).buf,
        )
    }

show_accounts_stmt:
    SHOW ACCOUNTS like_opt
    {
        // Like
        $$ = tree.NewShowAccounts(
            $3,
            yylex.(*Lexer).buf,
        )
    }

show_publications_stmt:
    SHOW PUBLICATIONS like_opt
    {
        // like
	    $$ = tree.NewShowPublications(
            $3,
            yylex.(*Lexer).buf,
        )
    }


show_subscriptions_stmt:
    SHOW SUBSCRIPTIONS like_opt
    {
        // like
	    $$ = tree.NewShowSubscriptions(
            $3,
            yylex.(*Lexer).buf,
        )
    }

like_opt:
    {
        $$ = nil
    }
|   LIKE simple_expr
    {
        $$ = tree.NewComparisonExpr(tree.LIKE, nil, $2, yylex.(*Lexer).buf)
    }
|   ILIKE simple_expr
    {
        $$ = tree.NewComparisonExpr(tree.ILIKE, nil, $2, yylex.(*Lexer).buf)
    }

database_name_opt:
    {
        $$ = ""
    }
|   from_or_in ident
    {
        $$ = $2.Compare()
    }

table_column_name:
    from_or_in unresolved_object_name
    {
        $$ = $2
    }

from_or_in:
    FROM
|   IN

fields_or_columns:
    FIELDS
|   COLUMNS

full_opt:
    {
        $$ = false
    }
|   FULL
    {
        $$ = true
    }

show_create_stmt:
    SHOW CREATE TABLE table_name_unresolved
    {
        // Name
        $$ = tree.NewShowCreateTable(
            $4,
            yylex.(*Lexer).buf,
        )
    }
|
    SHOW CREATE VIEW table_name_unresolved
    {
        // Name
        $$ = tree.NewShowCreateView(
            $4,
            yylex.(*Lexer).buf,
        )
    }
|   SHOW CREATE DATABASE not_exists_opt db_name
    {
        // IfNotExists
        // Name
        $$ = tree.NewShowCreateDatabase(
            $4,
            $5,
            yylex.(*Lexer).buf,
        )
    }
|   SHOW CREATE PUBLICATION db_name
    {
	    $$ = tree.NewShowCreatePublications(
            $4,
            yylex.(*Lexer).buf,
        )
    }

show_servers_stmt:
    SHOW BACKEND SERVERS
    {
        $$ = tree.NewShowBackendServers(yylex.(*Lexer).buf)
    }

table_name_unresolved:
    ident
    {
        $$ = tree.SetUnresolvedObjectName(1, [3]string{$1.Compare()}, yylex.(*Lexer).buf)
    }
|   ident '.' ident
    {
        $$ = tree.SetUnresolvedObjectName(2, [3]string{$3.Compare(), $1.Compare()}, yylex.(*Lexer).buf)
    }

db_name:
    ident
    {
    	$$ = $1.Compare()
    }

unresolved_object_name:
    ident
    {
        $$ = tree.SetUnresolvedObjectName(1, [3]string{$1.Compare()}, yylex.(*Lexer).buf)
    }
|   ident '.' ident
    {
        $$ = tree.SetUnresolvedObjectName(2, [3]string{$3.Compare(), $1.Compare()}, yylex.(*Lexer).buf)
    }
|   ident '.' ident '.' ident
    {
        $$ = tree.SetUnresolvedObjectName(3, [3]string{$5.Compare(), $3.Compare(), $1.Compare()}, yylex.(*Lexer).buf)
    }

truncate_table_stmt:
    TRUNCATE table_name
    {
    	$$ = tree.NewTruncateTable($2, yylex.(*Lexer).buf)
    }
|   TRUNCATE TABLE table_name
    {
	    $$ = tree.NewTruncateTable($3, yylex.(*Lexer).buf)
    }

drop_stmt:
    drop_ddl_stmt

drop_ddl_stmt:
    drop_database_stmt
|   drop_prepare_stmt
|   drop_table_stmt
|   drop_view_stmt
|   drop_index_stmt
|   drop_role_stmt
|   drop_user_stmt
|   drop_account_stmt
|   drop_function_stmt
|   drop_sequence_stmt
|   drop_publication_stmt
|   drop_procedure_stmt
|   drop_stage_stmt
|   drop_connector_stmt

drop_sequence_stmt:
    DROP SEQUENCE exists_opt table_name_list
    {
        // IfExists Names
        $$ = tree.NewDropSequence(
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }

drop_account_stmt:
    DROP ACCOUNT exists_opt account_name
    {
        // IfExists Name
        $$ = tree.NewDropAccount(
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }

drop_user_stmt:
    DROP USER exists_opt drop_user_spec_list
    {
        // IfExists Users
        $$ = tree.NewDropUser(
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }

drop_user_spec_list:
    drop_user_spec
    {
        $$ = buffer.MakeSlice[*tree.User](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.User](yylex.(*Lexer).buf, $$, $1)
    }
|   drop_user_spec_list ',' drop_user_spec
    {
        $$ = buffer.MakeSlice[*tree.User](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.User](yylex.(*Lexer).buf, $1, $3)
    }

drop_user_spec:
    user_name
    {
        // UserName HostName AuthOption
        $$ = tree.NewUser(
            $1.Username.Get(),
            $1.Hostname.Get(),
            nil,
            yylex.(*Lexer).buf,
        )
    }

drop_role_stmt:
    DROP ROLE exists_opt role_spec_list
    {
        // IfExists Roles
        $$ = tree.NewDropRole(
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }

drop_index_stmt:
    DROP INDEX exists_opt ident ON table_name
    {
        // Name TableName IfExists
        $$ = tree.NewDropIndex(
            tree.Identifier($4.Compare()),
            *$6,
            $3,
            yylex.(*Lexer).buf,
        )
    }

drop_table_stmt:
    DROP TABLE temporary_opt exists_opt table_name_list drop_table_opt
    {
        // IfExists Names
        $$ = tree.NewDropTable(
            $4,
            $5,
            yylex.(*Lexer).buf,
        )
    }
|   DROP STREAM exists_opt table_name_list
    {
        // IfExists Names
        $$ = tree.NewDropTable(
            $3, 
            $4, 
            yylex.(*Lexer).buf,
        )
    }

drop_connector_stmt:
    DROP CONNECTOR exists_opt table_name_list
    {
        // IfExists Names
    	$$ = tree.NewDropConnector(
            $3, 
            $4,
            yylex.(*Lexer).buf,
        )
    }

drop_view_stmt:
    DROP VIEW exists_opt table_name_list
    {
        // ifExists Names
        $$ = tree.NewDropView(
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }

drop_database_stmt:
    DROP DATABASE exists_opt ident
    {
        // Name IfExists
        $$ = tree.NewDropDatabase(
            tree.Identifier($4.Compare()),
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   DROP SCHEMA exists_opt ident
    {
        // Name IfExists
        $$ = tree.NewDropDatabase(
            tree.Identifier($4.Compare()),
            $3,
            yylex.(*Lexer).buf,
        )
    }

drop_prepare_stmt:
    DROP PREPARE stmt_name
    {
        $$ = tree.NewDeallocate(tree.Identifier($3), true, yylex.(*Lexer).buf)
    }

drop_function_stmt:
    DROP FUNCTION func_name '(' func_args_list_opt ')'
    {
        // Name Args
        $$ = tree.NewDropFunction(
            $3,
            $5,
            yylex.(*Lexer).buf,
        )
    }

drop_procedure_stmt:
    DROP PROCEDURE proc_name
    {
        // ProcedureName IfExists
        $$ = tree.NewDropProcedure(
            $3,
            false,
            yylex.(*Lexer).buf,
        )
    }
|    DROP PROCEDURE IF EXISTS proc_name
    {
        // ProcedureName IfExists
        $$ = tree.NewDropProcedure(
            $5,
            true,
            yylex.(*Lexer).buf,
        )
    }

delete_stmt:
    delete_without_using_stmt
|    delete_with_using_stmt
|    with_clause delete_with_using_stmt
    {
        $2.(*tree.Delete).With = $1
        $$ = $2
    }
|    with_clause delete_without_using_stmt
    {
        $2.(*tree.Delete).With = $1
        $$ = $2
    }

delete_without_using_stmt:
    DELETE priority_opt quick_opt ignore_opt FROM table_name partition_clause_opt as_opt_id where_expression_opt order_by_opt limit_opt
    {
        // Single-Table Syntax
        t := tree.NewAliasedTableExpr(
            $6,
            tree.NewAliasClause(
                tree.Identifier($8),
                nil,
                yylex.(*Lexer).buf,
            ),
            nil,
            yylex.(*Lexer).buf,
        )

        ts := buffer.MakeSlice[tree.TableExpr](yylex.(*Lexer).buf)
        ts = buffer.AppendSlice[tree.TableExpr](yylex.(*Lexer).buf, ts, t)

        // Table Where
        $$ = tree.NewDelete(
            ts,
            $9,
            yylex.(*Lexer).buf,
        )
        $$.(*tree.Delete).OrderBy = $10
        $$.(*tree.Delete).Limit = $11
    }
|    DELETE priority_opt quick_opt ignore_opt table_name_wild_list FROM table_references where_expression_opt
    {
        // Multiple-Table Syntax
        // Table Where
        ts := buffer.MakeSlice[tree.TableExpr](yylex.(*Lexer).buf)
        ts = buffer.AppendSlice[tree.TableExpr](yylex.(*Lexer).buf, ts, $7)

        $$ = tree.NewDelete(
            $5,
            $8,
            yylex.(*Lexer).buf,
        )
        $$.(*tree.Delete).TableRefs = ts
    }

delete_with_using_stmt:
    DELETE priority_opt quick_opt ignore_opt FROM table_name_wild_list USING table_references where_expression_opt
    {
        // Multiple-Table Syntax
        // Table Where
        ts := buffer.MakeSlice[tree.TableExpr](yylex.(*Lexer).buf)
        ts = buffer.AppendSlice[tree.TableExpr](yylex.(*Lexer).buf, ts, $8)

        $$ = tree.NewDelete(
            $6,
            $9,
            yylex.(*Lexer).buf,
        )
        $$.(*tree.Delete).TableRefs = ts
    }

table_name_wild_list:
    table_name_opt_wild
    {
        $$ = buffer.MakeSlice[tree.TableExpr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.TableExpr](yylex.(*Lexer).buf, $$, $1)
    }
|    table_name_wild_list ',' table_name_opt_wild
    {
        $$ = buffer.MakeSlice[tree.TableExpr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.TableExpr](yylex.(*Lexer).buf, $1, $3)
    }

table_name_opt_wild:
    ident wild_opt
    {
        prefix := tree.NewObjectNamePrefix(yylex.(*Lexer).buf)
        prefix.ExplicitSchema = false
        $$ = tree.NewTableName(tree.Identifier($1.Compare()), *prefix, yylex.(*Lexer).buf)
    }
|    ident '.' ident wild_opt
    {
        prefix := tree.NewObjectNamePrefix(yylex.(*Lexer).buf)
        bi := tree.NewBufIdentifier($1.Compare())
        yylex.(*Lexer).buf.Pin(bi)
        prefix.SchemaName = bi
        prefix.ExplicitSchema = true
        $$ = tree.NewTableName(tree.Identifier($3.Compare()), *prefix, yylex.(*Lexer).buf)
    }

wild_opt:
    %prec EMPTY
    {}
|    '.' '*'
    {}

priority_opt:
    {}
|    priority

priority:
    LOW_PRIORITY
|    HIGH_PRIORITY
|    DELAYED

quick_opt:
    {}
|    QUICK

ignore_opt:
    {}
|    IGNORE

replace_stmt:
    REPLACE into_table_name partition_clause_opt replace_data
    {
    	rep := $4
    	rep.Table = $2
    	rep.PartitionNames = $3
    	$$ = rep
    }

replace_data:
    VALUES values_list
    {
        // Rows RowWord
        vc := tree.NewValuesClause($2, false, yylex.(*Lexer).buf)
        $$ = tree.NewReplace(nil, tree.NewSelect(vc, nil, nil, nil, nil, yylex.(*Lexer).buf), yylex.(*Lexer).buf)
    }
|   select_stmt
    {
        $$ = tree.NewReplace(nil, $1, yylex.(*Lexer).buf)
    }
|   '(' insert_column_list ')' VALUES values_list
    {
        // Rows RowWord
        vc := tree.NewValuesClause($5, false, yylex.(*Lexer).buf)
        $$ = tree.NewReplace($2, tree.NewSelect(vc, nil, nil, nil, nil, yylex.(*Lexer).buf), yylex.(*Lexer).buf)
    }
|   '(' ')' VALUES values_list
    {
        // Rows RowWord
        vc := tree.NewValuesClause($4, false, yylex.(*Lexer).buf)
        $$ = tree.NewReplace(nil, tree.NewSelect(vc, nil, nil, nil, nil, yylex.(*Lexer).buf), yylex.(*Lexer).buf)
    }
|   '(' insert_column_list ')' select_stmt
    {
        $$ = tree.NewReplace($2, $4, yylex.(*Lexer).buf)
    }
|	SET set_value_list
	{
		if $2 == nil {
			yylex.Error("the set list of replace can not be empty")
			return 1
		}
        identList := buffer.MakeSlice[tree.Identifier](yylex.(*Lexer).buf)
        valueList := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
		for _, a := range $2 {
            identList = buffer.AppendSlice[tree.Identifier](yylex.(*Lexer).buf, identList, a.Column.Get())
            valueList = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, valueList, a.Expr)
		}
        
        vexprs := buffer.MakeSlice[tree.Exprs](yylex.(*Lexer).buf)
        vexprs = buffer.AppendSlice[tree.Exprs](yylex.(*Lexer).buf, vexprs, valueList)
        
		vc := tree.NewValuesClause(vexprs, false, yylex.(*Lexer).buf)
		$$ = tree.NewReplace(identList, tree.NewSelect(vc, nil, nil, nil, nil, yylex.(*Lexer).buf), yylex.(*Lexer).buf)
	}

insert_stmt:
    INSERT into_table_name partition_clause_opt insert_data on_duplicate_key_update_opt
    {
        ins := $4
        ins.Table = $2
        ins.PartitionNames = $3
        ins.OnDuplicateUpdate = $5
        $$ = ins
    }

accounts_list:
    account_name
    {
        $$ = buffer.MakeSlice[tree.Identifier](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Identifier](yylex.(*Lexer).buf, $$, tree.Identifier($1))
    }
|   accounts_list ',' account_name
    {
        $$ = buffer.MakeSlice[tree.Identifier](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Identifier](yylex.(*Lexer).buf, $1, tree.Identifier($3))
    }

insert_data:
    VALUES values_list
    {
        // Rows RowWord
        vc := tree.NewValuesClause($2, false, yylex.(*Lexer).buf)
        $$ = tree.NewInsert(nil, tree.NewSelect(vc, nil, nil, nil, nil, yylex.(*Lexer).buf), yylex.(*Lexer).buf)
    }
|   select_stmt
    {
        $$ = tree.NewInsert(nil, $1, yylex.(*Lexer).buf)
    }
|   '(' insert_column_list ')' VALUES values_list
    {
        // Rows RowWord
        vc := tree.NewValuesClause($5, false, yylex.(*Lexer).buf)
        $$ = tree.NewInsert($2, tree.NewSelect(vc, nil, nil, nil, nil, yylex.(*Lexer).buf), yylex.(*Lexer).buf)
    }
|   '(' ')' VALUES values_list
    {
        // Rows RowWord
        vc := tree.NewValuesClause($4, false, yylex.(*Lexer).buf)
        $$ = tree.NewInsert(nil, tree.NewSelect(vc, nil, nil, nil, nil, yylex.(*Lexer).buf), yylex.(*Lexer).buf)
    }
|   '(' insert_column_list ')' select_stmt
    {
        $$ = tree.NewInsert($2, $4, yylex.(*Lexer).buf)
    }
|   SET set_value_list
    {
        if $2 == nil {
            yylex.Error("the set list of insert can not be empty")
            return 1
        }

        identList := buffer.MakeSlice[tree.Identifier](yylex.(*Lexer).buf)
        valueList := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        
		for _, a := range $2 {
            identList = buffer.AppendSlice[tree.Identifier](yylex.(*Lexer).buf, identList, a.Column.Get())
            valueList = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, valueList, a.Expr)
		}
        
        vexprs := buffer.MakeSlice[tree.Exprs](yylex.(*Lexer).buf)
        vexprs = buffer.AppendSlice[tree.Exprs](yylex.(*Lexer).buf, vexprs, valueList)
        
        // Rows RowWord
        vc := tree.NewValuesClause(vexprs, false, yylex.(*Lexer).buf)

        $$ = tree.NewInsert(identList, tree.NewSelect(vc, nil, nil, nil, nil, yylex.(*Lexer).buf), yylex.(*Lexer).buf)
    }

on_duplicate_key_update_opt:
    {
        $$ = buffer.MakeSlice[*tree.UpdateExpr](yylex.(*Lexer).buf)
    }
|   ON DUPLICATE KEY UPDATE update_list
    {
      	$$ = $5
    }

set_value_list:
    {
        $$ = nil
    }
|    set_value
    {
        $$ = buffer.MakeSlice[*tree.Assignment](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.Assignment](yylex.(*Lexer).buf, $$, $1)
    }
|    set_value_list ',' set_value
    {
        $$ = buffer.MakeSlice[*tree.Assignment](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.Assignment](yylex.(*Lexer).buf, $1, $3)
    }

set_value:
    insert_column '=' expr_or_default
    {
        $$ = tree.NewAssignment(
            tree.Identifier($1),
            $3,
            yylex.(*Lexer).buf,
        )
    }

insert_column_list:
    insert_column
    {
        $$ = buffer.MakeSlice[tree.Identifier](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Identifier](yylex.(*Lexer).buf, $$, tree.Identifier($1))
    }
|   insert_column_list ',' insert_column
    {
        $$ = buffer.MakeSlice[tree.Identifier](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Identifier](yylex.(*Lexer).buf, $1, tree.Identifier($3))
    }

insert_column:
    ident
    {
        $$ = $1.Compare()
    }
|   ident '.' ident
    {
        $$ = $3.Compare()
    }

values_list:
    row_value
    {
        $$ = buffer.MakeSlice[tree.Exprs](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Exprs](yylex.(*Lexer).buf, $$, $1)
    }
|   values_list ',' row_value
    {
        $$ = buffer.MakeSlice[tree.Exprs](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Exprs](yylex.(*Lexer).buf, $1, $3)
    }

row_value:
    row_opt '(' data_opt ')'
    {
        $$ = $3
    }

row_opt:
    {}
|    ROW

data_opt:
    {
        $$ = nil
    }
|   data_values

data_values:
    expr_or_default
    {
        $$ = buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, $$, $1)
    }
|   data_values ',' expr_or_default
    {
        $$ = buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, $1, $3)
    }

expr_or_default:
    expression
|   DEFAULT
    {
        $$ = tree.NewDefaultVal(yylex.(*Lexer).buf)
    }

partition_clause_opt:
    {
        $$ = nil
    }
|   PARTITION '(' partition_id_list ')'
    {
        $$ = $3
    }

partition_id_list:
    ident
    {
        $$ = buffer.MakeSlice[tree.Identifier](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Identifier](yylex.(*Lexer).buf, $$, tree.Identifier($1.Compare()))
    }
|   partition_id_list ',' ident
    {
        $$ = buffer.MakeSlice[tree.Identifier](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Identifier](yylex.(*Lexer).buf, $1, tree.Identifier($3.Compare()))
    }

into_table_name:
    INTO table_name
    {
        $$ = $2
    }
|   table_name
    {
        $$ = $1
    }

export_data_param_opt:
    {
        $$ = nil
    }
|   INTO OUTFILE STRING export_fields export_lines_opt header_opt max_file_size_opt force_quote_opt
    {
        // Outfile QueryId FilePath Fields Lines Header MaxFileSize ForeQuote
        $$ = tree.NewExportParam(
            true,
            "",
            $3,
            $4,
            $5,
            $6,
            uint64($7)*1024,
            $8,
            yylex.(*Lexer).buf,
        )
    }

export_fields:
    {
        // Terminated Optionally EnclosedBy EscapedBy
        $$ = tree.NewFields(
            ",",
            false,
            '"',
            0,
            yylex.(*Lexer).buf,
        )
    }
|   FIELDS TERMINATED BY STRING
    {
        // Terminated Optionally EnclosedBy EscapedBy
        $$ = tree.NewFields(
            $4,
            false,
            '"',
            0,
            yylex.(*Lexer).buf,
        )
    }
|   FIELDS TERMINATED BY STRING ENCLOSED BY field_terminator
    {
        str := $7
        if str != "\\" && len(str) > 1 {
            yylex.Error("export1 error field terminator")
            return 1
        }
        var b byte
        if len(str) != 0 {
           b = byte(str[0])
        } else {
           b = 0
        }
        // Terminated Optionally EnclosedBy EscapedBy
        $$ = tree.NewFields(
            $4,
            false,
            b,
            0,
            yylex.(*Lexer).buf,
        )
    }
|   FIELDS ENCLOSED BY field_terminator
    {
        str := $4
        if str != "\\" && len(str) > 1 {
            yylex.Error("export2 error field terminator")
            return 1
        }
        var b byte
        if len(str) != 0 {
           b = byte(str[0])
        } else {
           b = 0
        }
        // Terminated Optionally EnclosedBy EscapedBy
        $$ = tree.NewFields(
            ",",
            false,
            b,
            0,
            yylex.(*Lexer).buf,
        )
    }

export_lines_opt:
    {
        // StartingBy TerminatedBy
        $$ = tree.NewLines(
            "",
            "\n",
            yylex.(*Lexer).buf,
        )
    }
|   LINES lines_terminated_opt
    {
        // StartingBy TerminatedBy
        $$ = tree.NewLines(
            "",
            $2,
            yylex.(*Lexer).buf,
        )
    }

header_opt:
    {
        $$ = true
    }
|   HEADER STRING
    {
        str := strings.ToLower($2)
        if str == "true" {
            $$ = true
        } else if str == "false" {
            $$ = false
        } else {
            yylex.Error("error header flag")
            return 1
        }
    }

max_file_size_opt:
    {
        $$ = 0
    }
|   MAX_FILE_SIZE INTEGRAL
    {
        $$ = $2.(int64)
    }

force_quote_opt:
    {
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
    }
|   FORCE_QUOTE '(' force_quote_list ')'
    {
        $$ = $3
    }


force_quote_list:
    ident
    {
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $$, yylex.(*Lexer).buf.CopyString($1.Compare()))
    }
|   force_quote_list ',' ident
    {
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $1, yylex.(*Lexer).buf.CopyString($3.Compare()))
    }

select_stmt:
    select_no_parens
|   select_with_parens
    {
        // Select OrderBy Limit Ep SelectLockInfo
        $$ = tree.NewSelect($1, nil, nil, nil, nil, yylex.(*Lexer).buf)
    }

select_no_parens:
    simple_select order_by_opt limit_opt export_data_param_opt select_lock_opt
    {
        // Select OrderBy Limit Ep SelectLockInfo
        $$ = tree.NewSelect($1, $2, $3, $4, $5, yylex.(*Lexer).buf)
    }
|   select_with_parens order_by_clause export_data_param_opt
    {
        // Select OrderBy Limit Ep SelectLockInfo
        $$ = tree.NewSelect($1, $2, nil, $3, nil, yylex.(*Lexer).buf)
    }
|   select_with_parens order_by_opt limit_clause export_data_param_opt
    {
        // Select OrderBy Limit Ep SelectLockInfo
        $$ = tree.NewSelect($1, $2, $3, $4, nil, yylex.(*Lexer).buf)
    }
|   with_clause simple_select order_by_opt limit_opt export_data_param_opt select_lock_opt
    {
        // Select OrderBy Limit Ep SelectLockInfo
        se := tree.NewSelect($2, $3, $4, $5, $6, yylex.(*Lexer).buf)
        se.With = $1
        $$ = se
    }
|   with_clause select_with_parens order_by_clause export_data_param_opt
    {
        // Select OrderBy Limit Ep SelectLockInfo
        se := tree.NewSelect($2, $3, nil, $4, nil, yylex.(*Lexer).buf)
        se.With = $1
        $$ = se
    }
|   with_clause select_with_parens order_by_opt limit_clause export_data_param_opt
    {
        // Select OrderBy Limit Ep SelectLockInfo
        se := tree.NewSelect($2, $3, $4, $5, nil, yylex.(*Lexer).buf)
        se.With = $1
        $$ = se
    }

with_clause:
    WITH cte_list
    {
        // IsRecursive CTEs
        $$ = tree.NewWith(
            false,
            $2,
            yylex.(*Lexer).buf,
        )
    }
|    WITH RECURSIVE cte_list
    {
        // IsRecursive CTEs
        $$ = tree.NewWith(
            true,
            $3,
            yylex.(*Lexer).buf,
        )
    }

cte_list:
    common_table_expr
    {
        $$ = buffer.MakeSlice[*tree.CTE](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.CTE](yylex.(*Lexer).buf, $$, $1)
    }
|    cte_list ',' common_table_expr
    {
        $$ = buffer.MakeSlice[*tree.CTE](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.CTE](yylex.(*Lexer).buf, $1, $3)
    }

common_table_expr:
    ident column_list_opt AS '(' stmt ')'
    {
        // Name Stmt
        $$ = tree.NewCTE(
            tree.NewAliasClause(tree.Identifier($1.Compare()), $2, yylex.(*Lexer).buf),
            $5,
            yylex.(*Lexer).buf,
        )
    }

column_list_opt:
    {
        $$ = nil
    }
|    '(' column_list ')'
    {
        $$ = $2
    }

limit_opt:
    {
        $$ = nil
    }
|   limit_clause
    {
        $$ = $1
    }

limit_clause:
    LIMIT expression
    {
        // Offset Count
        $$ = tree.NewLimit(
            nil,
            $2,
            yylex.(*Lexer).buf,
        )
    }
|   LIMIT expression ',' expression
    {
        // Offset Count
        $$ = tree.NewLimit(
            $2,
            $4,
            yylex.(*Lexer).buf,
        )
    }
|   LIMIT expression OFFSET expression
    {
        // Offset Count
        $$ = tree.NewLimit(
            $4,
            $2,
            yylex.(*Lexer).buf,
        )
    }

order_by_opt:
    {
        $$ = nil
    }
|   order_by_clause
    {
        $$ = $1
    }

order_by_clause:
    ORDER BY order_list
    {
        $$ = $3
    }

order_list:
    order
    {
        $$ = buffer.MakeSlice[*tree.Order](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.Order](yylex.(*Lexer).buf, $$, $1)
    }
|   order_list ',' order
    {
        $$ = buffer.MakeSlice[*tree.Order](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.Order](yylex.(*Lexer).buf, $1, $3)
    }

order:
    expression asc_desc_opt nulls_first_last_opt
    {
        // Expr Direction NullsPosition NullOrder
        $$ = tree.NewOrder(
            $1,
            $2,
            $3,
            false,
            yylex.(*Lexer).buf,
        )
    }

asc_desc_opt:
    {
        $$ = tree.DefaultDirection
    }
|   ASC
    {
        $$ = tree.Ascending
    }
|   DESC
    {
        $$ = tree.Descending
    }

nulls_first_last_opt:
    {
        $$ = tree.DefaultNullsPosition
    }
|   NULLS FIRST
    {
        $$ = tree.NullsFirst
    }
|   NULLS LAST
    {
        $$ = tree.NullsLast
    }

select_lock_opt:
    {
        $$ = nil
    }
|   FOR UPDATE
    {
        // LockType
        $$ = tree.NewSelectLockInfo(
            tree.SelectLockForUpdate,
            yylex.(*Lexer).buf,
        )
    }

select_with_parens:
    '(' select_no_parens ')'
    {
        $$ = tree.NewParenSelect($2, yylex.(*Lexer).buf)
    }
|   '(' select_with_parens ')'
    {
        s := tree.NewSelect($2, nil, nil, nil, nil, yylex.(*Lexer).buf)
        $$ = tree.NewParenSelect(s, yylex.(*Lexer).buf)
    }
|   '(' values_stmt ')'
    {
        valuesStmt := $2.(*tree.ValuesStatement);
        
        // Select OrderBy Limit Ep SelectLockInfo
        s := tree.NewSelect(
            tree.NewValuesClause(valuesStmt.Rows, true, yylex.(*Lexer).buf),
            valuesStmt.OrderBy,
            valuesStmt.Limit,
            nil,
            nil,
            yylex.(*Lexer).buf,
        )
        $$ = tree.NewParenSelect(s, yylex.(*Lexer).buf)
    }

simple_select:
    simple_select_clause
    {
        $$ = $1
    }
|   simple_select union_op simple_select_clause
    {
        // Type Left Right All Distinct
        $$ = tree.NewUnionClause(
            $2.Type,
            $1,
            $3,
            $2.All,
            $2.Distinct,
            yylex.(*Lexer).buf,
        )
    }
|   select_with_parens union_op simple_select_clause
    {
        // Type Left Right All Distinct
        $$ = tree.NewUnionClause(
            $2.Type,
            $1,
            $3,
            $2.All,
            $2.Distinct,
            yylex.(*Lexer).buf,
        )
    }
|   simple_select union_op select_with_parens
    {
        // Type Left Right All Distinct
        $$ = tree.NewUnionClause(
            $2.Type,
            $1,
            $3,
            $2.All,
            $2.Distinct,
            yylex.(*Lexer).buf,
        )
    }
|   select_with_parens union_op select_with_parens
    {
        // Type Left Right All Distinct
        $$ = tree.NewUnionClause(
            $2.Type,
            $1,
            $3,
            $2.All,
            $2.Distinct,
            yylex.(*Lexer).buf,
        )
    }

union_op:
    UNION
    {
        $$ = tree.NewUnionTypeRecord( 
            tree.UNION,
            false,
            false,
            yylex.(*Lexer).buf,
        )
    }
|   UNION ALL
    {
        $$ = tree.NewUnionTypeRecord( 
            tree.UNION,
            true,
            false,
            yylex.(*Lexer).buf,
        )
    }
|   UNION DISTINCT
    {
        $$ = tree.NewUnionTypeRecord( 
            tree.UNION,
            false,
            true,
            yylex.(*Lexer).buf,
        )
    }
|
    EXCEPT
    {
        $$ = tree.NewUnionTypeRecord( 
            tree.EXCEPT,
            false,
            false,
            yylex.(*Lexer).buf,
        )
    }
|   EXCEPT ALL
    {
        $$ = tree.NewUnionTypeRecord( 
            tree.EXCEPT,
            true,
            false,
            yylex.(*Lexer).buf,
        )
    }
|   EXCEPT DISTINCT
    {
        $$ = tree.NewUnionTypeRecord( 
            tree.EXCEPT,
            false,
            true,
            yylex.(*Lexer).buf,
        )
    }
|    INTERSECT
    {
        $$ = tree.NewUnionTypeRecord( 
            tree.INTERSECT,
            false,
            false,
            yylex.(*Lexer).buf,
        )
    }
|   INTERSECT ALL
    {
        $$ = tree.NewUnionTypeRecord( 
            tree.INTERSECT,
            true,
            false,
            yylex.(*Lexer).buf,
        )
    }
|   INTERSECT DISTINCT
    {
        $$ = tree.NewUnionTypeRecord( 
            tree.INTERSECT,
            false,
            true,
            yylex.(*Lexer).buf,
        )
    }
|    MINUS
    {
        $$ = tree.NewUnionTypeRecord( 
            tree.UT_MINUS,
            false,
            false,
            yylex.(*Lexer).buf,
        )
    }
|   MINUS ALL
    {
        $$ = tree.NewUnionTypeRecord( 
            tree.UT_MINUS,
            true,
            false,
            yylex.(*Lexer).buf,
        )
    }
|    MINUS DISTINCT
    {
        $$ = tree.NewUnionTypeRecord( 
            tree.UT_MINUS,
            false,
            true,
            yylex.(*Lexer).buf,
        )
    }

simple_select_clause:
    SELECT distinct_opt select_expression_list from_opt where_expression_opt group_by_opt having_opt
    {
    $$ = tree.NewSelectClause( 
            $2,
            $3,
            $4,
            $5,
            $6,
            $7,
            "",
            yylex.(*Lexer).buf,
        )
    }
|    SELECT select_option_opt select_expression_list from_opt where_expression_opt group_by_opt having_opt
    {
        $$ = tree.NewSelectClause( 
            false,
            $3,
            $4,
            $5,
            $6,
            $7,
            $2,
            yylex.(*Lexer).buf,
        )
    }

select_option_opt:
    SQL_SMALL_RESULT
    {
    	$$ = strings.ToLower($1)
    }
|	SQL_BIG_RESULT
	{
       $$ = strings.ToLower($1)
    }
|   SQL_BUFFER_RESULT
	{
    	$$ = strings.ToLower($1)
    }

distinct_opt:
    {
        $$ = false
    }
|   ALL
    {
        $$ = false
    }
|   distinct_keyword
    {
        $$ = true
    }

distinct_keyword:
    DISTINCT
|   DISTINCTROW

having_opt:
    {
        $$ = nil
    }
|   HAVING expression
    {
        // type expr
        $$ = tree.NewWhere(
            tree.AstHaving,
            $2,
            yylex.(*Lexer).buf,
        )
    }

group_by_opt:
    {
        $$ = nil
    }
|   GROUP BY expression_list
    {
        $$ = tree.NewGroupBy($3, yylex.(*Lexer).buf)
    }

where_expression_opt:
    {
        $$ = nil
    }
|   WHERE expression
    {
        // type expr
        $$ = tree.NewWhere(
            tree.AstWhere,
            $2,
            yylex.(*Lexer).buf,
        )
    }

select_expression_list:
    select_expression
    {
        $$ = buffer.MakeSlice[tree.SelectExpr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.SelectExpr](yylex.(*Lexer).buf, $$, $1)
    }
|   select_expression_list ',' select_expression
    {
        $$ = buffer.MakeSlice[tree.SelectExpr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.SelectExpr](yylex.(*Lexer).buf, $1, $3)
    }

select_expression:
    '*' %prec '*'
    {
        s := buffer.Alloc[tree.SelectExpr](yylex.(*Lexer).buf)
        s.Expr = tree.StarExpr()
        $$ = *s
    }
|   expression as_name_opt
    {
    	$2.SetConfig(0)
        s := buffer.Alloc[tree.SelectExpr](yylex.(*Lexer).buf)
        s.Expr = $1
        s.As = $2
        $$ = *s
    }
|   ident '.' '*' %prec '*'
    {
        s := buffer.Alloc[tree.SelectExpr](yylex.(*Lexer).buf)
        s.Expr = tree.SetUnresolvedNameWithStar(yylex.(*Lexer).buf, $1.Compare())
        $$ = *s
    }
|   ident '.' ident '.' '*' %prec '*'
    {
        s := buffer.Alloc[tree.SelectExpr](yylex.(*Lexer).buf)
        s.Expr = tree.SetUnresolvedNameWithStar(yylex.(*Lexer).buf, $3.Compare(), $1.Compare())
        $$ = *s
    }

from_opt:
    {
        prefix := tree.NewObjectNamePrefix(yylex.(*Lexer).buf)
        prefix.ExplicitSchema = false
        tn := tree.NewTableName(tree.Identifier(""), *prefix, yylex.(*Lexer).buf)
        tns := tree.NewAliasedTableExpr(tn, nil, nil, yylex.(*Lexer).buf)

        tables := buffer.MakeSlice[tree.TableExpr](yylex.(*Lexer).buf)
        tables = buffer.AppendSlice[tree.TableExpr](yylex.(*Lexer).buf, tables, tns)
        $$ = tree.NewFrom(
            tables,
            yylex.(*Lexer).buf,
        )
    }
|   from_clause
    {
        $$ = $1
    }

from_clause:
    FROM table_references
    {
        tables := buffer.MakeSlice[tree.TableExpr](yylex.(*Lexer).buf)
        tables = buffer.AppendSlice[tree.TableExpr](yylex.(*Lexer).buf, tables, $2)
        $$ = tree.NewFrom(
            tables,
            yylex.(*Lexer).buf,
        )
    }

table_references:
    escaped_table_reference
   	{
   		if t, ok := $1.(*tree.JoinTableExpr); ok {
   			$$ = t
   		} else {
            // Left JoinType Right JoinCondition
   			$$ = tree.NewJoinTableExpr(
                $1,
                tree.JOIN_TYPE_CROSS,
                nil,
                nil,
                yylex.(*Lexer).buf,
            )
   		}
    }
|   table_references ',' escaped_table_reference
    {
        // Left JoinType Right JoinCondition
        $$ = tree.NewJoinTableExpr(
            $1,
            tree.JOIN_TYPE_CROSS,
            $3, 
            nil,
            yylex.(*Lexer).buf,
        )
    }

escaped_table_reference:
    table_reference %prec LOWER_THAN_SET

table_reference:
    table_factor
|   join_table
	{
		$$ = $1
	}

join_table:
    table_reference inner_join table_factor join_condition_opt
    {
        // JoinType Left Right Cond
        $$ = tree.NewJoinTableExpr(
            $1,
            $2,
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }
|   table_reference straight_join table_factor on_expression_opt
    {
        // JoinType Left Right Cond
        $$ = tree.NewJoinTableExpr(
            $1,
            $2,
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }
|   table_reference outer_join table_factor join_condition
    {
        // JoinType Left Right Cond
        $$ = tree.NewJoinTableExpr(
            $1,
            $2,
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }
|   table_reference natural_join table_factor
    {
        // JoinType Left Right Cond
        $$ = tree.NewJoinTableExpr(
            $1,
            $2,
            $3,
            nil,
            yylex.(*Lexer).buf,
        )
    }

natural_join:
    NATURAL JOIN
    {
        $$ = tree.JOIN_TYPE_NATURAL
    }
|   NATURAL outer_join
    {
        if $2 == tree.JOIN_TYPE_LEFT {
            $$ = tree.JOIN_TYPE_NATURAL_LEFT
        } else {
            $$ = tree.JOIN_TYPE_NATURAL_RIGHT
        }
    }

outer_join:
    LEFT JOIN
    {
        $$ = tree.JOIN_TYPE_LEFT
    }
|   LEFT OUTER JOIN
    {
        $$ = tree.JOIN_TYPE_LEFT
    }
|   RIGHT JOIN
    {
        $$ = tree.JOIN_TYPE_RIGHT
    }
|   RIGHT OUTER JOIN
    {
        $$ = tree.JOIN_TYPE_RIGHT
    }

values_stmt:
    VALUES row_constructor_list order_by_opt limit_opt
    {
        // Rows OrderBy Limit
        $$ = tree.NewValuesStatement(
            $2,
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }

row_constructor_list:
    row_constructor
    {
        $$ = buffer.MakeSlice[tree.Exprs](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Exprs](yylex.(*Lexer).buf, $$, $1)
    }
|   row_constructor_list ',' row_constructor
    {
        $$ = buffer.MakeSlice[tree.Exprs](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Exprs](yylex.(*Lexer).buf, $1, $3)
    }

row_constructor:
    ROW '(' data_values ')'
    {
        $$ = $3
    }

on_expression_opt:
    %prec JOIN
    {
        $$ = nil
    }
|   ON expression
    {
        $$ = tree.NewOnJoinCond($2, yylex.(*Lexer).buf)
    }

straight_join:
    STRAIGHT_JOIN
    {
        $$ = tree.JOIN_TYPE_STRAIGHT
    }

inner_join:
    JOIN
    {
        $$ = tree.JOIN_TYPE_INNER
    }
|   INNER JOIN
    {
        $$ = tree.JOIN_TYPE_INNER
    }
|   CROSS JOIN
    {
        $$ = tree.JOIN_TYPE_CROSS
    }

join_condition_opt:
    %prec JOIN
    {
        $$ = nil
    }
|   join_condition
    {
        $$ = $1
    }

join_condition:
    ON expression
    {
        $$ = tree.NewOnJoinCond($2, yylex.(*Lexer).buf)
    }
|   USING '(' column_list ')'
    {
        $$ = tree.NewUsingJoinCond($3, yylex.(*Lexer).buf)
    }

column_list:
    ident
    {
        $$ = buffer.MakeSlice[tree.Identifier](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Identifier](yylex.(*Lexer).buf, $$, tree.Identifier($1.Compare()))
    }
|   column_list ',' ident
    {
        $$ = buffer.MakeSlice[tree.Identifier](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Identifier](yylex.(*Lexer).buf, $1, tree.Identifier($3.Compare()))
    }

table_factor:
    aliased_table_name
    {
        $$ = $1
    }
|   table_subquery as_opt_id column_list_opt
    {
        // TableExpr AliasCluase IndexHints
        $$ = tree.NewAliasedTableExpr(
            $1,
            tree.NewAliasClause(
                tree.Identifier($2),
                $3,
                yylex.(*Lexer).buf,
            ),
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   table_function as_opt_id
    {
        if $2 != "" {
            // TableExpr AliasCluase IndexHints
            $$ = tree.NewAliasedTableExpr(
                $1,
                tree.NewAliasClause(
                    tree.Identifier($2),
                    nil,
                    yylex.(*Lexer).buf,
                ),
                nil,
                yylex.(*Lexer).buf,
            )
        } else {
            $$ = $1
        }
    }
|   '(' table_references ')'
	{
		$$ = $2
	}

table_subquery:
    select_with_parens %prec SUBQUERY_AS_EXPR
    {
    	$$ = tree.NewParenTableExpr($1.(*tree.ParenSelect).Select, yylex.(*Lexer).buf)
    }

table_function:
    ident '(' expression_list_opt ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1.Compare()))
        function := tree.NewFuncExpr(tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf), $3, yylex.(*Lexer).buf)
        function.Type = tree.FUNC_TYPE_TABLE
        
        $$ = tree.NewTableFunction(
       	    function,
            yylex.(*Lexer).buf,
        )
    }

aliased_table_name:
    table_name as_opt_id index_hint_list_opt
    {
        // Expr As IndexHints
        $$ = tree.NewAliasedTableExpr(
            $1,
            tree.NewAliasClause(
                tree.Identifier($2),
                nil,
                yylex.(*Lexer).buf,
            ),
            $3,
            yylex.(*Lexer).buf,
        )
    }

index_hint_list_opt:
	{
		$$ = nil
	}
|	index_hint_list

index_hint_list:
	index_hint
	{
        $$ = buffer.MakeSlice[*tree.IndexHint](yylex.(*Lexer).buf)
		$$ = buffer.AppendSlice[*tree.IndexHint](yylex.(*Lexer).buf, $$, $1)
	}
|	index_hint_list index_hint
	{
        $$ = buffer.MakeSlice[*tree.IndexHint](yylex.(*Lexer).buf)
		$$ = buffer.AppendSlice[*tree.IndexHint](yylex.(*Lexer).buf, $1, $2)
	}

index_hint:
	index_hint_type index_hint_scope '(' index_name_list ')'
	{
        // indexNames IndexHintType IndexHintScope
		$$ = tree.NewIndexHint(
			$4,
			$1,
			$2,
            yylex.(*Lexer).buf,
		)
	}

index_hint_type:
	USE key_or_index
	{
		$$ = tree.HintUse
	}
|	IGNORE key_or_index
	{
		$$ = tree.HintIgnore
	}
|	FORCE key_or_index
	{
		$$ = tree.HintForce
	}

index_hint_scope:
	{
		$$ = tree.HintForScan
	}
|	FOR JOIN
	{
		$$ = tree.HintForJoin
	}
|	FOR ORDER BY
	{
		$$ = tree.HintForOrderBy
	}
|	FOR GROUP BY
	{
		$$ = tree.HintForGroupBy
	}

index_name_list:
	{
		$$ = nil
	}
|	ident
	{
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $$, yylex.(*Lexer).buf.CopyString($1.Compare()))

	}
|	index_name_list ',' ident
	{
        /* $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf) */
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $1, yylex.(*Lexer).buf.CopyString($3.Compare()))
	}
|	PRIMARY
	{
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $$, yylex.(*Lexer).buf.CopyString($1))
	}
|	index_name_list ',' PRIMARY
	{
        /* $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf) */
		$$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $1, yylex.(*Lexer).buf.CopyString($3))
	}

as_opt_id:
    {
        $$ = ""
    }
|   table_alias
    {
        $$ = $1
    }
|   AS table_alias
    {
        $$ = $2
    }

table_alias:
    ident
    {
    	$$ = $1.Compare()
    }
|   STRING

as_name_opt:
    {
        $$ = tree.NewCStr("", yylex.(*Lexer).lower, yylex.(*Lexer).buf)
    }
|   ident
    {
        $$ = $1
    }
|   AS ident
    {
        $$ = $2
    }
|   STRING
    {
        $$ = tree.NewCStr($1, yylex.(*Lexer).lower, yylex.(*Lexer).buf)
    }
|   AS STRING
    {
        $$ = tree.NewCStr($2, yylex.(*Lexer).lower, yylex.(*Lexer).buf)
    }

stmt_name:
    ident
    {
    	$$ = $1.Compare()
    }

//table_id:
//    id_or_var
//|   non_reserved_keyword

//id_or_var:
//    ID
//|	QUOTE_ID
//|   AT_ID
//|   AT_AT_ID

create_stmt:
    create_ddl_stmt
|   create_role_stmt
|   create_user_stmt
|   create_account_stmt
|   create_publication_stmt
|   create_stage_stmt

create_ddl_stmt:
    create_table_stmt
|   create_database_stmt
|   create_index_stmt
|   create_view_stmt
|   create_function_stmt
|   create_extension_stmt
|   create_sequence_stmt
|   create_procedure_stmt
|   create_stream_stmt
|   create_connector_stmt
|   pause_daemon_task_stmt
|   cancel_daemon_task_stmt
|   resume_daemon_task_stmt

create_extension_stmt:
    CREATE EXTENSION extension_lang AS extension_name FILE STRING
    {
        // Language Name FileName
        $$ = tree.NewCreateExtension(
            $3,
            tree.Identifier($5),
            tree.Identifier($7),
            yylex.(*Lexer).buf,
        )
    }

extension_lang:
    ident
    {
        $$ = $1.Compare()
    }

extension_name:
    ident
    {
        $$ = $1.Compare()
    }

create_procedure_stmt:
    CREATE PROCEDURE proc_name '(' proc_args_list_opt ')' STRING
    {
        // Name Args Body
        $$ = tree.NewCreateProcedure(
            $3,
            $5,
            $7,
            yylex.(*Lexer).buf,
        )
    }

proc_name:
    ident
    {
        prefix := tree.NewObjectNamePrefix(yylex.(*Lexer).buf)
        prefix.ExplicitSchema = false
        $$ = tree.NewProcedureName(tree.Identifier($1.ToLower()), *prefix, yylex.(*Lexer).buf)
    }
|   ident '.' ident
    {
        prefix := tree.NewObjectNamePrefix(yylex.(*Lexer).buf)
        bi := tree.NewBufIdentifier($1.ToLower())
        yylex.(*Lexer).buf.Pin(bi)
        prefix.SchemaName = bi
        prefix.ExplicitSchema = true
        $$ = tree.NewProcedureName(tree.Identifier($3.ToLower()), *prefix, yylex.(*Lexer).buf)
    }

proc_args_list_opt:
    {
        $$ = buffer.MakeSlice[tree.ProcedureArg](yylex.(*Lexer).buf)
    }
|   proc_args_list

proc_args_list:
    proc_arg
    {
        $$ = buffer.MakeSlice[tree.ProcedureArg](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.ProcedureArg](yylex.(*Lexer).buf, $$, $1)
    }
|   proc_args_list ',' proc_arg
    {
        $$ = buffer.MakeSlice[tree.ProcedureArg](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.ProcedureArg](yylex.(*Lexer).buf, $1, $3)
    }

proc_arg:
    proc_arg_decl
    {
        p := *buffer.Alloc[tree.ProcedureArg](yylex.(*Lexer).buf)
        p = $1
        $$ = p
    }

proc_arg_decl:
    proc_arg_in_out_type column_name column_type
    {
        $$ = tree.NewProcedureArgDecl($1, $2, $3, yylex.(*Lexer).buf)
    }

proc_arg_in_out_type:
    {
        $$ = tree.TYPE_IN
    }
|   IN
    {
        $$ = tree.TYPE_IN
    }
|   OUT
    {
        $$ = tree.TYPE_OUT
    }
|   INOUT
    {
        $$ = tree.TYPE_INOUT
    }


create_function_stmt:
    CREATE FUNCTION func_name '(' func_args_list_opt ')' RETURNS func_return LANGUAGE func_lang AS STRING
    {
        // Name Args ReturnType Language Body
        $$ = tree.NewCreateFunction(
            $3,
            $5,
            $8,
            $10,
            $12,
            yylex.(*Lexer).buf,
        )
    }

func_name:
    ident
    {
        prefix := tree.NewObjectNamePrefix(yylex.(*Lexer).buf)
        prefix.ExplicitSchema = false
        $$ = tree.NewFuncName(tree.Identifier($1.Compare()), *prefix, yylex.(*Lexer).buf)
    }
|   ident '.' ident
    {
        prefix := tree.NewObjectNamePrefix(yylex.(*Lexer).buf)
        prefix.ExplicitSchema = true
        bi := tree.NewBufIdentifier($1.Compare())
        yylex.(*Lexer).buf.Pin(bi)
        prefix.SchemaName = bi
        $$ = tree.NewFuncName(tree.Identifier($3.Compare()), *prefix, yylex.(*Lexer).buf)
    }

func_args_list_opt:
    {
        $$ = buffer.MakeSlice[tree.FunctionArg](yylex.(*Lexer).buf)
    }
|   func_args_list

func_args_list:
    func_arg
    {
        $$ = buffer.MakeSlice[tree.FunctionArg](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.FunctionArg](yylex.(*Lexer).buf, $$, $1)
    }
|   func_args_list ',' func_arg
    {
        $$ = buffer.MakeSlice[tree.FunctionArg](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.FunctionArg](yylex.(*Lexer).buf, $1, $3)
    }

func_arg:
    func_arg_decl
    {
        f := *buffer.Alloc[tree.FunctionArg](yylex.(*Lexer).buf)
        f = $1
        $$ = f
    }

func_arg_decl:
    column_type
    {
        $$ = tree.NewFunctionArgDecl(nil, $1, nil, yylex.(*Lexer).buf)
    }
|   column_name column_type
    {
        $$ = tree.NewFunctionArgDecl($1, $2, nil, yylex.(*Lexer).buf)
    }
|   column_name column_type DEFAULT literal
    {
        $$ = tree.NewFunctionArgDecl($1, $2, $4, yylex.(*Lexer).buf)
    }

func_lang:
    ident
    {
        $$ = $1.Compare()
    }

func_return:
    column_type
    {
        $$ = tree.NewReturnType($1, yylex.(*Lexer).buf)
    }

create_view_stmt:
    CREATE view_list_opt VIEW not_exists_opt table_name column_list_opt AS select_stmt view_tail
    {
        // Replace Name ColNames AsSource IfNotExists
        $$ = tree.NewCreateView(
            false,
            $5,
            $6,
            $8,
            $4,
            yylex.(*Lexer).buf,
        )
    }
|   CREATE replace_opt VIEW not_exists_opt table_name column_list_opt AS select_stmt view_tail
    {
        // Replace Name ColNames AsSource IfNotExists
        $$ = tree.NewCreateView(
            $2,
            $5,
            $6,
            $8,
            $4,
            yylex.(*Lexer).buf,
        )
    }

create_account_stmt:
    CREATE ACCOUNT not_exists_opt account_name account_auth_option account_status_option account_comment_opt
    {
        // IfNotExists Name AuthOption StatusOption Comment
   		$$ = tree.NewCreateAccount(
            $3,
            $4,
            $5,
            $6,
            $7,
            yylex.(*Lexer).buf,
    	)
    }

view_list_opt:
    view_opt
    {
        $$ = $1
    }
|   view_list_opt view_opt
    {
        $$ = $$ + $2
    }

view_opt:
    ALGORITHM '=' algorithm_type_2
    {
        $$ = "ALGORITHM = " + $3
    }
|   DEFINER '=' user_name
    {
        $$ = "DEFINER = "
    }
|   SQL SECURITY security_opt
    {
        $$ = "SQL SECURITY " + $3
    }

view_tail:
    {
        $$ = ""
    }
|   WITH check_type CHECK OPTION
    {
        $$ = "WITH " + $2 + " CHECK OPTION"
    }

algorithm_type_2:
    UNDEFINED
|   MERGE
|   TEMPTABLE

security_opt:
    DEFINER
|   INVOKER

check_type:
    {
        $$ = ""
    }
|   CASCADED
|   LOCAL

account_name:
    ident
    {
    	$$ = $1.Compare()
    }

account_auth_option:
    ADMIN_NAME equal_opt account_admin_name account_identified
    {
        // Equal AdminName IdentifiedType
        $$ = *tree.NewAccountAuthOption(
            $2,
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }

account_admin_name:
    STRING
    {
    	$$ = $1
    }
|	ident
	{
		$$ = $1.Compare()
	}

account_identified:
    IDENTIFIED BY STRING
    {
        // Typ Str
        $$ = *tree.NewAccountIdentified(
            tree.AccountIdentifiedByPassword,
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   IDENTIFIED BY RANDOM PASSWORD
    {
        // Typ Str
        $$ = *tree.NewAccountIdentified(
            tree.AccountIdentifiedByRandomPassword,
            "",
            yylex.(*Lexer).buf,
        )
    }
|   IDENTIFIED WITH STRING
    {
        // Typ Str
        $$ = *tree.NewAccountIdentified(
            tree.AccountIdentifiedWithSSL,
            $3,
            yylex.(*Lexer).buf,
        )
    }

account_status_option:
    {
        // Exist Option
        $$ = *tree.NewAccountStatus(
            false,
            0,
            yylex.(*Lexer).buf,
        )
    }
|   OPEN
    {
        // Exist Option
        $$ = *tree.NewAccountStatus(
            true,
            tree.AccountStatusOpen,
            yylex.(*Lexer).buf,
        )
    }
|   SUSPEND
    {
        // Exist Option
        $$ = *tree.NewAccountStatus(
            true,
            tree.AccountStatusSuspend,
            yylex.(*Lexer).buf,
        )
    }
|   RESTRICTED
    {
        // Exist Option
        $$ = *tree.NewAccountStatus(
            true,
            tree.AccountStatusRestricted,
            yylex.(*Lexer).buf,
        )
    }

account_comment_opt:
    {
        // Exist Comment
        $$ = *tree.NewAccountComment(
            false,
            "",
            yylex.(*Lexer).buf,
        )
    }
|   COMMENT_KEYWORD STRING
    {
        // Exist Comment
        $$ = *tree.NewAccountComment(
            true,
            $2,
            yylex.(*Lexer).buf,
        )
    }

create_user_stmt:
    CREATE USER not_exists_opt user_spec_list_of_create_user default_role_opt pwd_or_lck_opt user_comment_or_attribute_opt
    {
        // IfNotExists User Role MiscOpt CommentOrAttribute
        $$ = tree.NewCreateUser(
            $3,
            $4,
            $5,
            $6,
            $7,
            yylex.(*Lexer).buf,
        )
    }

create_publication_stmt:
    CREATE PUBLICATION not_exists_opt ident DATABASE ident alter_publication_accounts_opt comment_opt
    {
        // IfNotExists Name Database AccountsSet Comment
	    $$ = tree.NewCreatePublication(
	        $3,
	        tree.Identifier($4.Compare()),
	        tree.Identifier($6.Compare()),
	        $7,
	        $8,
            yylex.(*Lexer).buf,
	    )
    }

create_stage_stmt:
    CREATE STAGE not_exists_opt ident urlparams stage_credentials_opt stage_status_opt stage_comment_opt
    {
        // IfNotExists Name Url Credentials Status Comment
        $$ = tree.NewCreateStage(
            $3,
            tree.Identifier($4.Compare()),
            $5,
            $6,
            $7,
            $8,
            yylex.(*Lexer).buf,
        )
    }

stage_status_opt:
    {
        // Exist Option
        $$ = *tree.NewStageStatus(
            false,
            0,
            yylex.(*Lexer).buf,
        )
    }
|   ENABLE '=' TRUE
    {
        // Exist Option
        $$ = *tree.NewStageStatus(
            true,
            tree.StageStatusEnabled,
            yylex.(*Lexer).buf,
        )
    }
|   ENABLE '=' FALSE
    {
        // Exist Option
        $$ = *tree.NewStageStatus(
            true,
            tree.StageStatusDisabled,
            yylex.(*Lexer).buf,
        )
    }

stage_comment_opt:
    {
        // Exist Comment
        $$ = *tree.NewStageComment(
            false,
            "",
            yylex.(*Lexer).buf,
        )
    }
|   COMMENT_KEYWORD STRING
    {
        // Exist Comment
        $$ = *tree.NewStageComment(
            true,
            $2,
            yylex.(*Lexer).buf,
        )
    }

stage_url_opt:
    {
        // Exist URL
        $$ = *tree.NewStageUrl(
            false,
            "",
            yylex.(*Lexer).buf,
        )
    }
|   URL '=' STRING
    {
        // Exist URL
        $$ = *tree.NewStageUrl(
            true,
            $3,
            yylex.(*Lexer).buf,
        )
    }

stage_credentials_opt:
    {
        // Exist Credentials
        $$ = *tree.NewStageCredentials(
            false,
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   CREDENTIALS '=' '{' credentialsparams '}'
    {
        // Exist Credentials
        $$ = *tree.NewStageCredentials(
            true,
            $4,
            yylex.(*Lexer).buf,
        )
    }

credentialsparams:
    credentialsparam
    {
        $$ = $1
    }
|   credentialsparams ',' credentialsparam
    {
        /* $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf) */
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $1, $3...)
    }

credentialsparam:
    {
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
    }
|   STRING '=' STRING
    {
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $$, yylex.(*Lexer).buf.CopyString($1), yylex.(*Lexer).buf.CopyString($3))
    }

urlparams:
    URL '=' STRING
    {
        $$ = $3
    }

comment_opt:
    {
        $$ = ""
    }
|   COMMENT_KEYWORD STRING
    {
        $$ = $2
    }

alter_stage_stmt:
    ALTER STAGE exists_opt ident SET stage_url_opt stage_credentials_opt stage_status_opt stage_comment_opt
    {
        // IfNotExists Name UrlOption CredentialsOption StatusOption Comment
        $$ = tree.NewAlterStage(
            	$3,
	            tree.Identifier($4.Compare()),
	            $6,
	            $7,
	            $8,
	            $9,
                yylex.(*Lexer).buf,
            )
    }


alter_publication_stmt:
    ALTER PUBLICATION exists_opt ident alter_publication_accounts_opt comment_opt
    {
        // IfExists Name AccountsSet Comment
	    $$ = tree.NewAlterPublication(
	        $3,
	        tree.Identifier($4.Compare()),
	        $5,
	        $6,
            yylex.(*Lexer).buf,
	    )
    }

alter_publication_accounts_opt:
    {
	    $$ = nil
    }
    | ACCOUNT ALL
    {
        // All SetAccounts AddAccounts DropAccounts
	    $$ = tree.NewAccountsSetOption( 
            true,
            nil,
            nil,
            nil,
            yylex.(*Lexer).buf,
	    )
    }
    | ACCOUNT accounts_list
    {
        // All SetAccounts AddAccounts DropAccounts
	    $$ = tree.NewAccountsSetOption( 
            false,
            $2,
            nil,
            nil,
            yylex.(*Lexer).buf,
	    )
    }
    | ACCOUNT ADD accounts_list
    {
        // All SetAccounts AddAccounts DropAccounts
	    $$ = tree.NewAccountsSetOption( 
            false,
            nil,
            $3,
            nil,
            yylex.(*Lexer).buf,
	    )
    }
    | ACCOUNT DROP accounts_list
    {
        // All SetAccounts AddAccounts DropAccounts
	    $$ = tree.NewAccountsSetOption( 
            false,
            nil,
            nil,
            $3,
            yylex.(*Lexer).buf,
	    )
    }


drop_publication_stmt:
DROP PUBLICATION exists_opt ident
    {
        // ifExists Name
	    $$ = tree.NewDropPublication(
	        $3,
	        tree.Identifier($4.Compare()),
            yylex.(*Lexer).buf,
	    )
    }

drop_stage_stmt:
DROP STAGE exists_opt ident
    {
        // ifNotExists Name
        $$ = tree.NewDropStage(
            $3,
	        tree.Identifier($4.Compare()),
            yylex.(*Lexer).buf,
        )
    }

account_role_name:
    ident
    {
   		$$ = $1.Compare()
    }

user_comment_or_attribute_opt:
    {
        // Exist IsComment Str
        $$ = tree.NewAccountCommentOrAttribute( 
            false,
            false,
            "",
            yylex.(*Lexer).buf,
        )
    }
|   COMMENT_KEYWORD STRING
    {
        // Exist IsComment Str
        $$ = tree.NewAccountCommentOrAttribute( 
            true,
            true,
            $2,
            yylex.(*Lexer).buf,
        )
    }
|   ATTRIBUTE STRING
    {
        // Exist IsComment Str
        $$ = tree.NewAccountCommentOrAttribute( 
            true,
            false,
            $2,
            yylex.(*Lexer).buf,
        )
    }

//conn_options:
//    {
//        $$ = nil
//    }
//|   WITH conn_option_list
//    {
//        $$ = $2
//    }
//
//conn_option_list:
//    conn_option
//    {
//        $$ = []tree.ResourceOption{$1}
//    }
//|   conn_option_list conn_option
//    {
//        $$ = buffer.AppendSlice[tree.Statement](yylex.(*Lexer).buf, $1, $2)
//    }
//
//conn_option:
//    MAX_QUERIES_PER_HOUR INTEGRAL
//    {
//        $$ = &tree.ResourceOptionMaxQueriesPerHour{Count: $2.(int64)}
//    }
//|   MAX_UPDATES_PER_HOUR INTEGRAL
//    {
//        $$ = &tree.ResourceOptionMaxUpdatesPerHour{Count: $2.(int64)}
//    }
//|   MAX_CONNECTIONS_PER_HOUR INTEGRAL
//    {
//        $$ = &tree.ResourceOptionMaxConnectionPerHour{Count: $2.(int64)}
//    }
//|   MAX_USER_CONNECTIONS INTEGRAL
//    {
//        $$ = &tree.ResourceOptionMaxUserConnections{Count: $2.(int64)}
//    }


//require_clause_opt:
//    {
//        $$ = nil
//    }
//|   require_clause
//
//require_clause:
//    REQUIRE NONE
//    {
//        t := &tree.TlsOptionNone{}
//        $$ = []tree.TlsOption{t}
//    }
//|   REQUIRE SSL
//    {
//        t := &tree.TlsOptionSSL{}
//        $$ = []tree.TlsOption{t}
//    }
//|   REQUIRE X509
//    {
//        t := &tree.TlsOptionX509{}
//        $$ = []tree.TlsOption{t}
//    }
//|   REQUIRE require_list
//    {
//        $$ = $2
//    }
//
//require_list:
//    require_elem
//    {
//        $$ = []tree.TlsOption{$1}
//    }
//|   require_list AND require_elem
//    {
//        $$ = buffer.AppendSlice[tree.Statement](yylex.(*Lexer).buf, $1, $3)
//    }
//|   require_list require_elem
//    {
//        $$ = buffer.AppendSlice[tree.Statement](yylex.(*Lexer).buf, $1, $2)
//    }
//
//require_elem:
//    ISSUER STRING
//    {
//        $$ = &tree.TlsOptionIssuer{Issuer: $2}
//    }
//|   SUBJECT STRING
//    {
//        $$ = &tree.TlsOptionSubject{Subject: $2}
//    }
//|   CIPHER STRING
//    {
//        $$ = &tree.TlsOptionCipher{Cipher: $2}
//    }
//|   SAN STRING
//    {
//        $$ = &tree.TlsOptionSan{San: $2}
//    }
user_spec_list_of_create_user:
    user_spec_with_identified
    {
        $$ = buffer.MakeSlice[*tree.User](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.User](yylex.(*Lexer).buf, $$, $1)
    }
|   user_spec_list_of_create_user ',' user_spec_with_identified
    {
        $$ = buffer.MakeSlice[*tree.User](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.User](yylex.(*Lexer).buf, $1, $3)
    }

user_spec_with_identified:
    user_name user_identified
    {
        // Username Hostname AuthOption
        $$ = tree.NewUser(
            $1.Username.Get(),
            $1.Hostname.Get(),
            $2,
            yylex.(*Lexer).buf,
        )
    }

user_spec_list:
    user_spec
    {
        $$ = buffer.MakeSlice[*tree.User](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.User](yylex.(*Lexer).buf, $$, $1)
    }
|   user_spec_list ',' user_spec
    {
        $$ = buffer.MakeSlice[*tree.User](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.User](yylex.(*Lexer).buf, $1, $3)
    }

user_spec:
    user_name user_identified_opt
    {
        // Username Hostname AuthOption
        $$ = tree.NewUser(
            $1.Username.Get(),
            $1.Hostname.Get(),
            $2,
            yylex.(*Lexer).buf,
        )
    }

user_name:
    name_string
    {
        // Username Hostname
        $$ = tree.NewUsernameRecord(
            $1, 
            "%",
            yylex.(*Lexer).buf,
        )
    }
|   name_string '@' name_string
    {
        // Username Hostname
        $$ = tree.NewUsernameRecord(
            $1, 
            $3, 
            yylex.(*Lexer).buf,
        )
    }
|   name_string AT_ID
    {
        // Username Hostname
        $$ = tree.NewUsernameRecord(
            $1, 
            $2, 
            yylex.(*Lexer).buf,
        )
    }

user_identified_opt:
    {
        $$ = nil
    }
|   user_identified
    {
        $$ = $1
    }

user_identified:
    IDENTIFIED BY STRING
    {
        // Typ Str 
        $$ = tree.NewAccountIdentified(
            tree.AccountIdentifiedByPassword,
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   IDENTIFIED BY RANDOM PASSWORD
    {
        // Typ Str 
        $$ = tree.NewAccountIdentified(
            tree.AccountIdentifiedByRandomPassword,
            "",
            yylex.(*Lexer).buf,
        )
    }
|   IDENTIFIED WITH STRING
    {
        // Typ Str 
        $$ = tree.NewAccountIdentified(
            tree.AccountIdentifiedWithSSL,
            $3,
            yylex.(*Lexer).buf,
        )
    }

name_string:
    ident
    {
    	$$ = $1.Compare()
    }
|   STRING

create_role_stmt:
    CREATE ROLE not_exists_opt role_spec_list
    {
        // IfNotExists Roles
        $$ = tree.NewCreateRole(
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }

role_spec_list:
    role_spec
    {
        $$ = buffer.MakeSlice[*tree.Role](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.Role](yylex.(*Lexer).buf, $$, $1)
    }
|   role_spec_list ',' role_spec
    {
        $$ = buffer.MakeSlice[*tree.Role](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.Role](yylex.(*Lexer).buf, $1, $3)
    }

role_spec:
    role_name
    {
        $$ = tree.NewRole($1.Compare(), yylex.(*Lexer).buf)
    }
//|   name_string '@' name_string
//    {
//        $$ = &tree.Role{UserName: $1, HostName: $3}
//    }
//|   name_string AT_ID
//    {
//        $$ = &tree.Role{UserName: $1, HostName: $2}
//    }

role_name:
    ID
	{
		$$ = tree.NewCStr($1, yylex.(*Lexer).lower, yylex.(*Lexer).buf)
    }
|   QUOTE_ID
	{
		$$ = tree.NewCStr($1, yylex.(*Lexer).lower, yylex.(*Lexer).buf)
    }
|   STRING
	{
    	$$ = tree.NewCStr($1, yylex.(*Lexer).lower, yylex.(*Lexer).buf)
    }

index_prefix:
    {
        $$ = tree.INDEX_CATEGORY_NONE
    }
|   FULLTEXT
    {
        $$ = tree.INDEX_CATEGORY_FULLTEXT
    }
|   SPATIAL
    {
        $$ = tree.INDEX_CATEGORY_SPATIAL
    }
|   UNIQUE
    {
        $$ = tree.INDEX_CATEGORY_UNIQUE
    }

create_index_stmt:
    CREATE index_prefix INDEX ident using_opt ON table_name '(' index_column_list ')' index_option_list
    {
        var io *tree.IndexOption = nil
        if $11 == nil && $5 != tree.INDEX_TYPE_INVALID {
            io = tree.NewIndexOption(yylex.(*Lexer).buf)
            io.IType = $5
        } else if $11 != nil{
            io = tree.NewIndexOption(yylex.(*Lexer).buf)
            io = $11
            io.IType = $5
        }
        
        // Name Table IndexCat KeyParts IndexOption MiscOption
        $$ = tree.NewCreateIndex(
            tree.Identifier($4.Compare()),
            *$7,
            $2,
            $9,
            io,
            nil,
            yylex.(*Lexer).buf,
        )
    }

index_option_list:
    {
        $$ = nil
    }
|   index_option_list index_option
    {
        // Merge the options
        if $1 == nil {
            $$ = $2
        } else {
            opt1 := $1
            opt2 := $2
            if len(opt2.Comment.Get()) > 0 {
                opt1.Comment = opt2.Comment
            } else if opt2.KeyBlockSize > 0 {
                opt1.KeyBlockSize = opt2.KeyBlockSize
            } else if len(opt2.ParserName.Get()) > 0 {
                opt1.ParserName = opt2.ParserName
            } else if opt2.Visible != tree.VISIBLE_TYPE_INVALID {
                opt1.Visible = opt2.Visible
            }
            $$ = opt1
        }
    }

index_option:
    KEY_BLOCK_SIZE equal_opt INTEGRAL
    {
        // KeyBlockSize Comment ParserName Visible
        $$ = tree.NewIndexOption2(
            uint64($3.(int64)),
            "",
            "",
            0,
            yylex.(*Lexer).buf,
        )
    }
|   COMMENT_KEYWORD STRING
    {
        // KeyBlockSize Comment ParserName Visible
        $$ = tree.NewIndexOption2(
            0,
            $2,
            "",
            0,
            yylex.(*Lexer).buf,
        )
    }
|   WITH PARSER ident
    {
        // KeyBlockSize Comment ParserName Visible
        $$ = tree.NewIndexOption2(
            0,
            "",
            $3.Compare(),
            0,
            yylex.(*Lexer).buf,
        )
    }
|   VISIBLE
    {
        // KeyBlockSize Comment ParserName Visible
        $$ = tree.NewIndexOption2(
            0,
            "",
            "",
            tree.VISIBLE_TYPE_VISIBLE,
            yylex.(*Lexer).buf,
        )
    }
|   INVISIBLE
    {
        // KeyBlockSize Comment ParserName Visible
        $$ = tree.NewIndexOption2(
            0,
            "",
            "",
            tree.VISIBLE_TYPE_INVISIBLE,
            yylex.(*Lexer).buf,
        )
    }

index_column_list:
    index_column
    {
        $$ = buffer.MakeSlice[*tree.KeyPart](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.KeyPart](yylex.(*Lexer).buf, $$, $1)
    }
|   index_column_list ',' index_column
    {
        $$ = buffer.MakeSlice[*tree.KeyPart](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.KeyPart](yylex.(*Lexer).buf, $1, $3)
    }

index_column:
    column_name length_opt asc_desc_opt
    {
        // Order is parsed but just ignored as MySQL dtree.
        // ColName Length Expr Direction
        $$ = tree.NewKeyPart(
            $1,
            int($2),
            nil,
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   '(' expression ')' asc_desc_opt
    {
        // ColName Length Expr Direction
        $$ = tree.NewKeyPart(
            nil,
            0,
            $2,
            $4,
            yylex.(*Lexer).buf,
        )
    }

using_opt:
    {
        $$ = tree.INDEX_TYPE_INVALID
    }
|   USING BTREE
    {
        $$ = tree.INDEX_TYPE_BTREE
    }
|   USING HASH
    {
        $$ = tree.INDEX_TYPE_HASH
    }
|   USING RTREE
    {
        $$ = tree.INDEX_TYPE_RTREE
    }
|   USING BSI
    {
        $$ = tree.INDEX_TYPE_BSI
    }

create_database_stmt:
    CREATE database_or_schema not_exists_opt ident subcription_opt create_option_list_opt
    {
        // IfNotExists Name SubcriptionOption CreateOptions
        $$ = tree.NewCreateDatabase(
            $3,
            tree.Identifier($4.Compare()),
            $5,
            $6,
            yylex.(*Lexer).buf,
        )
    }
// CREATE comment_opt database_or_schema comment_opt not_exists_opt ident

subcription_opt:
    {
	    $$ = nil
    }
|   FROM account_name PUBLICATION ident
    {
        // From Publication
   	    $$ = tree.NewSubscriptionOption(
            tree.Identifier($2), 
            tree.Identifier($4.Compare()),
            yylex.(*Lexer).buf,
        )
    }


database_or_schema:
    DATABASE
|   SCHEMA

not_exists_opt:
    {
        $$ = false
    }
|   IF NOT EXISTS
    {
        $$ = true
    }

create_option_list_opt:
    {
        $$ = nil
    }
|   create_option_list
    {
        $$ = $1
    }

create_option_list:
    create_option
    {
        $$ = buffer.MakeSlice[tree.CreateOption](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.CreateOption](yylex.(*Lexer).buf, $$, $1)
    }
|   create_option_list create_option
    {
        $$ = buffer.MakeSlice[tree.CreateOption](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.CreateOption](yylex.(*Lexer).buf, $1, $2)
    }

create_option:
    default_opt charset_keyword equal_opt charset_name
    {
        // IsDefault Charset
        $$ = tree.NewCreateOptionCharset(
            $1,
            $4,
            yylex.(*Lexer).buf,
        )
    }
|   default_opt COLLATE equal_opt collate_name
    {
        // IsDefault Collate
        $$ = tree.NewCreateOptionCollate(
            $1,
            $4,
            yylex.(*Lexer).buf,
        )
    }
|   default_opt ENCRYPTION equal_opt STRING
    {
        // Encrypt
        $$ = tree.NewCreateOptionEncryption(
            $4, 
            yylex.(*Lexer).buf,
        )
    }

default_opt:
    {
        $$ = false
    }
|   DEFAULT
    {
        $$ = true
    }

create_connector_stmt:
    CREATE CONNECTOR FOR table_name WITH '(' connector_option_list ')'
    {
        // TableName Options
        $$ = tree.NewCreateConnector(
            $4,
            $7,
            yylex.(*Lexer).buf,
        )
    }

show_connectors_stmt:
    SHOW CONNECTORS
    {
	    $$ = tree.NewShowConnectors(yylex.(*Lexer).buf)
    }

pause_daemon_task_stmt:
    PAUSE DAEMON TASK INTEGRAL
    {
        var taskID uint64
        switch v := $4.(type) {
        case uint64:
    	    taskID = v
        case int64:
    	    taskID = uint64(v)
        default:
    	    yylex.Error("parse integral fail")
    	    return 1
        }
        // taskID
        $$ = tree.NewPauseDaemonTask(
            taskID,
            yylex.(*Lexer).buf,
        )
    }

cancel_daemon_task_stmt:
    CANCEL DAEMON TASK INTEGRAL
    {
        var taskID uint64
        switch v := $4.(type) {
        case uint64:
    	    taskID = v
        case int64:
    	    taskID = uint64(v)
        default:
    	    yylex.Error("parse integral fail")
    	    return 1
        }
        // TaskID
        $$ = tree.NewCancelDaemonTask(
            taskID,
            yylex.(*Lexer).buf,
        )
    }

resume_daemon_task_stmt:
    RESUME DAEMON TASK INTEGRAL
    {
        var taskID uint64
        switch v := $4.(type) {
        case uint64:
    	    taskID = v
        case int64:
    	    taskID = uint64(v)
        default:
    	    yylex.Error("parse integral fail")
    	    return 1
        }
        // TaskID
        $$ = tree.NewResumeDaemonTask(
            taskID,
            yylex.(*Lexer).buf,
        )
    }

create_stream_stmt:
    CREATE replace_opt STREAM not_exists_opt table_name '(' table_elem_list_opt ')' stream_option_list_opt
    {
        // Replace Source ifNotExists StreamName TableDefs ColNames AsSource Options
        $$ = tree.NewCreateStream(
            $2,
            false,
            $4,
            $5,
            $7,
            nil,
            nil,
            $9,
            yylex.(*Lexer).buf,
        )
    }
|   CREATE replace_opt SOURCE STREAM not_exists_opt table_name '(' table_elem_list_opt ')' stream_option_list_opt
    {
        // Replace Source ifNotExists StreamName TableDefs ColNames AsSource Options
        $$ = tree.NewCreateStream(
            $2,
            true,
            $5,
            $6,
            $8,
            nil,
            nil,
            $10,
            yylex.(*Lexer).buf,
        )
    }
|	CREATE replace_opt STREAM not_exists_opt table_name stream_option_list_opt AS select_stmt
    {
        // Replace Source ifNotExists StreamName TableDefs ColNames AsSource Options
        $$ = tree.NewCreateStream(
            $2,
            false,
            $4,
            $5,
            nil,
            nil,
            $8,
            $6,
            yylex.(*Lexer).buf,
        )
    }

replace_opt:
    {
        $$ = false
    }
|   OR REPLACE
    {
        $$ = true
    }

create_table_stmt:
    CREATE temporary_opt TABLE not_exists_opt table_name '(' table_elem_list_opt ')' table_option_list_opt partition_by_opt cluster_by_opt
    {
        // Temporary IsClusterTable IfNotExists TableDefs Options PartitionOption ClusterByOption Param
        $$ = tree.NewCreateTable(
            $2,
            false,
            $4,
            *$5,
            $7,
            $9,
            $10,
            $11,
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   CREATE EXTERNAL TABLE not_exists_opt table_name '(' table_elem_list_opt ')' load_param_opt_2
    {
        // Temporary IsClusterTable IfNotExists TableDefs Options PartitionOption ClusterByOption Param
        $$ = tree.NewCreateTable(
            false,
            false,
            $4,
            *$5,
            $7,
            nil,
            nil,
            nil,
            $9,
            yylex.(*Lexer).buf,
        )
    }
|   CREATE CLUSTER TABLE not_exists_opt table_name '(' table_elem_list_opt ')' table_option_list_opt partition_by_opt cluster_by_opt
    {
        // Temporary IsClusterTable IfNotExists TableDefs Options PartitionOption ClusterByOption Param
        $$ = tree.NewCreateTable(
            false,
            true, // isClusterTable
            $4,  //  ifNotExists
            *$5, // Table
            $7,  // Defs
            $9,  // Options
            $10, // PartitionOption
            $11, // ClusterByOption
            nil,
            yylex.(*Lexer).buf,
        )
    }
load_param_opt_2:
    load_param_opt tail_param_opt
    {
        $$ = $1
        $$.Tail = $2
    }

load_param_opt:
    INFILE STRING
    {
        // ExParamConst
        $$ = tree.NewExternParam(
            // ScanType FilePath CompressType Format Option Data
            tree.NewExParamConst(
                0,
                $2,
                tree.AUTO,
                tree.CSV,
                nil,
                "",
                yylex.(*Lexer).buf,
            ),
            yylex.(*Lexer).buf,
        )
    }
|   INLINE  FORMAT '=' STRING ','  DATA '=' STRING
    {
        // ExParamConst
        $$ = tree.NewExternParam(
            // ScanType FilePath CompressType Format Option Data
            tree.NewExParamConst(
                tree.INLINE,
                "",
                "",
                $4,
                nil,
                $8,
                yylex.(*Lexer).buf,
            ),
            yylex.(*Lexer).buf,
        )
    }
|   INFILE '{' infile_or_s3_params '}'
    {
        // ExParamConst
        $$ = tree.NewExternParam(
            // ScanType FilePath CompressType Format Option Data
            tree.NewExParamConst(
                0,
                "",
                "",
                "",
                $3,
                "",
                yylex.(*Lexer).buf,
            ),
            yylex.(*Lexer).buf,
        )
    }
|   URL S3OPTION '{' infile_or_s3_params '}'
    {
        // ExParamConst
        $$ = tree.NewExternParam(
            // ScanType FilePath CompressType Format Option Data
            tree.NewExParamConst(
                tree.S3,
                "",
                "",
                "",
                $4,
                "",
                yylex.(*Lexer).buf,
            ),
            yylex.(*Lexer).buf,
        )
    }

infile_or_s3_params:
    infile_or_s3_param
    {
        $$ = $1
    }
|   infile_or_s3_params ',' infile_or_s3_param
    {
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $1, $3...)
    }

infile_or_s3_param:
    {
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
    }
|   STRING '=' STRING
    {
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $$, yylex.(*Lexer).buf.CopyString($1), yylex.(*Lexer).buf.CopyString($3))
    }

tail_param_opt:
    load_fields load_lines ignore_lines columns_or_variable_list_opt load_set_spec_opt
    {
        // Fields Lines IgnoredLines ColumnList Assignments
        $$ = tree.NewTailParameter(
            $1,
            $2,
            uint64($3),
            $4,
            $5,
            yylex.(*Lexer).buf,
        )
    }
create_sequence_stmt:
    CREATE SEQUENCE not_exists_opt table_name as_datatype_opt increment_by_opt min_value_opt max_value_opt start_with_opt cycle_opt
    {
        // IfNotExists Name Type IncrementBy MinValue MaxValue StartWith Cycle
        $$ = tree.NewCreateSequence(
            $3,
            $4,
            $5,
            $6,
            $7,
            $8,
            $9,
            $10,
            yylex.(*Lexer).buf,
        )
    }
as_datatype_opt:
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, "bigint", 64, "", uint32(defines.MYSQL_TYPE_LONGLONG), yylex.(*Lexer).buf)
    }
|   AS column_type
    {
        $$ = $2
    }
alter_as_datatype_opt:
    {
        $$ = nil
    }
|   AS column_type
    {
        // Type
        $$ = tree.NewTypeOption(
            $2,
            yylex.(*Lexer).buf,
        )
    }
increment_by_opt:
    {
        $$ = nil
    }
|   INCREMENT BY INTEGRAL
    {
        $$ = tree.NewIncrementByOption(
            false,
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   INCREMENT INTEGRAL
    {
        // Minud Num
        $$ = tree.NewIncrementByOption(
            false,
            $2,
            yylex.(*Lexer).buf,
        )
    }
|   INCREMENT BY '-' INTEGRAL
    {
        // Minud Num
        $$ = tree.NewIncrementByOption(
            true,
            $4,
            yylex.(*Lexer).buf,
        )
    }
|   INCREMENT '-' INTEGRAL
    {
        // Minud Num
        $$ = tree.NewIncrementByOption(
            true,
            $3,
            yylex.(*Lexer).buf,
        )
    }
cycle_opt:
    {
        $$ = false
    }
|   NO CYCLE
    {
        $$ = false
    }
|   CYCLE
    {
        $$ = true
    }
min_value_opt:
    {
        $$ = nil
    }
|   MINVALUE INTEGRAL
    {
        // Minud Num
        $$ = tree.NewMinValueOption(
            false,
            $2,
            yylex.(*Lexer).buf,
        )
    }
|   MINVALUE '-' INTEGRAL
    {
        // Minud Num
        $$ = tree.NewMinValueOption(
            true,
            $3,
            yylex.(*Lexer).buf,
        )
    }
max_value_opt:
    {
        $$ = nil
    }
|   MAXVALUE INTEGRAL
    {
        // Minus Num
        $$ = tree.NewMaxValueOption(
            false,
            $2,
            yylex.(*Lexer).buf,
        )
    }
|   MAXVALUE '-' INTEGRAL
    {
        // Minus Num
        $$ = tree.NewMaxValueOption(
            true,
            $3,
            yylex.(*Lexer).buf,
        )
    }
alter_cycle_opt:
    {
        $$ = nil
    }
|   NO CYCLE
    {
        // Cycle
        $$ = tree.NewCycleOption(  
           false, 
           yylex.(*Lexer).buf,
        )
    }
|   CYCLE
    {
        // Cycle
        $$ = tree.NewCycleOption(  
           true, 
           yylex.(*Lexer).buf,
        )
    }
start_with_opt:
    {
        $$ = nil
    }
|   START WITH INTEGRAL
    {
        // Minus Num
        $$ = tree.NewStartWithOption(  
            false,
            $3, 
            yylex.(*Lexer).buf, 
        )  
    }
|   START INTEGRAL
    {
        // Minus Num
        $$ = tree.NewStartWithOption(  
            false,
            $2, 
            yylex.(*Lexer).buf, 
        )  
    }
|   START WITH '-' INTEGRAL
    {
        // Minus Num
        $$ = tree.NewStartWithOption(  
            true,
            $4, 
            yylex.(*Lexer).buf, 
        )  
    }
|   START '-' INTEGRAL
    {
        // Minus Num
        $$ = tree.NewStartWithOption(  
            true,
            $3, 
            yylex.(*Lexer).buf, 
        )  
    }
temporary_opt:
    {
        $$ = false
    }
|   TEMPORARY
    {
        $$ = true
    }

drop_table_opt:
    {
        $$ = true
    }
|   RESTRICT
    {
        $$ = true
    }
|   CASCADE
    {
        $$ = true
    }

partition_by_opt:
    {
        $$ = nil
    }
|   PARTITION BY partition_method partition_num_opt sub_partition_opt partition_list_opt
    {
        $3.Num = uint64($4)
        // PartBy SubPartBy Partitions
        $$ = tree.NewPartitionOption(
            *$3,
            $5,
            $6,
            yylex.(*Lexer).buf,
        )
    }

cluster_by_opt:
    {
        $$ = nil
    }
|   CLUSTER BY column_name
    {
        columnList := buffer.MakeSlice[*tree.UnresolvedName](yylex.(*Lexer).buf)
        columnList = buffer.AppendSlice[*tree.UnresolvedName](yylex.(*Lexer).buf, columnList, $3)
        // ColumnList
        $$ = tree.NewClusterByOption(columnList, yylex.(*Lexer).buf)

    }
    | CLUSTER BY '(' column_name_list ')'
    {
        // ColumnList
        $$ = tree.NewClusterByOption($4, yylex.(*Lexer).buf)
    }

sub_partition_opt:
    {
        $$ = nil
    }
|   SUBPARTITION BY sub_partition_method sub_partition_num_opt
    {
        // IsSubPartition PartitionType Num
        $$ = tree.NewPartitionBy2(
            true,
            $3,
            uint64($4),
            yylex.(*Lexer).buf,
        )
    }

partition_list_opt:
    {
        $$ = nil
    }
|   '(' partition_list ')'
    {
        $$ = $2
    }

partition_list:
    partition
    {
        $$ = buffer.MakeSlice[*tree.Partition](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.Partition](yylex.(*Lexer).buf, $$, $1)
    }
|   partition_list ',' partition
    {
        $$ = buffer.MakeSlice[*tree.Partition](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.Partition](yylex.(*Lexer).buf, $1, $3)
    }

partition:
    PARTITION ident values_opt sub_partition_list_opt
    {
        // Name Values Options Subs
        $$ = tree.NewPartition(
            tree.Identifier($2.Compare()),
            $3,
            nil,
            $4,
            yylex.(*Lexer).buf,
        )
    }
|   PARTITION ident values_opt partition_option_list sub_partition_list_opt
    {
        // Name Values Options Subs
        $$ = tree.NewPartition(
            tree.Identifier($2.Compare()),
            $3,
            $4,
            $5,
            yylex.(*Lexer).buf,
        )
    }

sub_partition_list_opt:
    {
        $$ = nil
    }
|   '(' sub_partition_list ')'
    {
        $$ = $2
    }

sub_partition_list:
    sub_partition
    {
        $$ = buffer.MakeSlice[*tree.SubPartition](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.SubPartition](yylex.(*Lexer).buf, $$, $1)
    }
|   sub_partition_list ',' sub_partition
    {
        $$ = buffer.MakeSlice[*tree.SubPartition](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.SubPartition](yylex.(*Lexer).buf, $1, $3)
    }

sub_partition:
    SUBPARTITION ident
    {
        // Name Options
        $$ =tree.NewSubPartition(
            tree.Identifier($2.Compare()),
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   SUBPARTITION ident partition_option_list
    {
        // Name Options
        $$ =tree.NewSubPartition(
            tree.Identifier($2.Compare()),
            $3,
            yylex.(*Lexer).buf,
        )
    }

partition_option_list:
    table_option
    {
        $$ = buffer.MakeSlice[tree.TableOption](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.TableOption](yylex.(*Lexer).buf, $$, $1)
    }
|   partition_option_list table_option
    {
        $$ = buffer.MakeSlice[tree.TableOption](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.TableOption](yylex.(*Lexer).buf, $1, $2)
    }

values_opt:
    {
        $$ = nil
    }
|   VALUES LESS THAN MAXVALUE
    {
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, tree.NewMaxValue(yylex.(*Lexer).buf))
        $$ = tree.NewValuesLessThan(
            exprs,
            yylex.(*Lexer).buf,
        )
    }
|   VALUES LESS THAN '(' expression_list ')'
    {
        $$ = tree.NewValuesLessThan(
            $5,
            yylex.(*Lexer).buf,
        )
    }
|   VALUES IN '(' expression_list ')'
    {
        $$ = tree.NewValuesIn($4, yylex.(*Lexer).buf)
    }

sub_partition_num_opt:
    {
        $$ = 0
    }
|   SUBPARTITIONS INTEGRAL
    {
        res := $2.(int64)
        if res == 0 {
            yylex.Error("partitions can not be 0")
            return 1
        }
        $$ = res
    }

partition_num_opt:
    {
        $$ = 0
    }
|   PARTITIONS INTEGRAL
    {
        res := $2.(int64)
        if res == 0 {
            yylex.Error("partitions can not be 0")
            return 1
        }
        $$ = res
    }

partition_method:
    RANGE '(' bit_expr ')'
    {
        // PartitionType
        $$ = tree.NewPartitionBy(
            tree.NewRangeType($3, nil, yylex.(*Lexer).buf),
            yylex.(*Lexer).buf,
        )
    }
|   RANGE fields_or_columns '(' column_name_list ')'
    {
        // PartitionType
        $$ = tree.NewPartitionBy(
            tree.NewRangeType(nil, $4, yylex.(*Lexer).buf),
            yylex.(*Lexer).buf,
        )
    }
|   LIST '(' bit_expr ')'
    {
        // PartitionType
        $$ = tree.NewPartitionBy(
            tree.NewListType($3, nil, yylex.(*Lexer).buf),
            yylex.(*Lexer).buf,
        )
    }
|   LIST fields_or_columns '(' column_name_list ')'
    {
        // PartitionType
        $$ = tree.NewPartitionBy(
            tree.NewListType(nil, $4, yylex.(*Lexer).buf),
            yylex.(*Lexer).buf,
        )
    }
|   sub_partition_method

sub_partition_method:
    linear_opt KEY algorithm_opt '(' ')'
    {
        // PartitionType
        $$ = tree.NewPartitionBy(
            tree.NewKeyType($1, nil, $3, yylex.(*Lexer).buf),
            yylex.(*Lexer).buf,
        )
    }
|   linear_opt KEY algorithm_opt '(' column_name_list ')'
    {
        // PartitionType
        $$ = tree.NewPartitionBy(
            tree.NewKeyType($1, $5, $3, yylex.(*Lexer).buf),
            yylex.(*Lexer).buf,
        )
    }
|   linear_opt HASH '(' bit_expr ')'
    {
        // PartitionType
        $$ = tree.NewPartitionBy(
            tree.NewHashType($1, $4, yylex.(*Lexer).buf),
            yylex.(*Lexer).buf,
        )
    }

algorithm_opt:
    {
        $$ = 2
    }
|   ALGORITHM '=' INTEGRAL
    {
        $$ = $3.(int64)
    }

linear_opt:
    {
        $$ = false
    }
|   LINEAR
    {
        $$ = true
    }

connector_option_list:
	connector_option
	{
        $$ = buffer.MakeSlice[*tree.ConnectorOption](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.ConnectorOption](yylex.(*Lexer).buf, $$, $1)
	}
|	connector_option_list ',' connector_option
	{
        $$ = buffer.MakeSlice[*tree.ConnectorOption](yylex.(*Lexer).buf)
		$$ = buffer.AppendSlice[*tree.ConnectorOption](yylex.(*Lexer).buf, $1, $3)
	}

connector_option:
	ident equal_opt literal
    {
        // key val
        $$ = tree.NewConnectorOption(
            tree.Identifier($1.Compare()), 
            $3,
            yylex.(*Lexer).buf,
        )
    }
    |   STRING equal_opt literal
        {
            // key val
             $$ = tree.NewConnectorOption(
                tree.Identifier($1), 
                $3,
                yylex.(*Lexer).buf,
            )
        }

stream_option_list_opt:
    {
        $$ = nil
    }
|	WITH '(' stream_option_list ')'
	{
		$$ = $3
	}

stream_option_list:
	stream_option
	{
        $$ = buffer.MakeSlice[tree.TableOption](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.TableOption](yylex.(*Lexer).buf, $$, $1)
	}
|	stream_option_list ',' stream_option
	{
        $$ = buffer.MakeSlice[tree.TableOption](yylex.(*Lexer).buf)
		$$ = buffer.AppendSlice[tree.TableOption](yylex.(*Lexer).buf, $1, $3)
	}

stream_option:
	ident equal_opt literal
    {
        // Identifier Expr 
        $$ = tree.NewCreateStreamWithOption(
            tree.Identifier($1.Compare()),
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   STRING equal_opt literal
    {
        // Identifier Expr 
        $$ = tree.NewCreateStreamWithOption(
            tree.Identifier($1),
            $3,
            yylex.(*Lexer).buf,
        )
    }

table_option_list_opt:
    {
        $$ = nil
    }
|   table_option_list
    {
        $$ = $1
    }

table_option_list:
    table_option
    {
        $$ = buffer.MakeSlice[tree.TableOption](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.TableOption](yylex.(*Lexer).buf, $$, $1)
    }
|   table_option_list ',' table_option
    {
        $$ = buffer.MakeSlice[tree.TableOption](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.TableOption](yylex.(*Lexer).buf, $1, $3)
    }
|   table_option_list table_option
    {
        $$ = buffer.MakeSlice[tree.TableOption](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.TableOption](yylex.(*Lexer).buf, $1, $2)
    }

table_option:
    AUTOEXTEND_SIZE equal_opt INTEGRAL
    {
        $$ = tree.NewTableOptionAUTOEXTEND_SIZE(uint64($3.(int64)), yylex.(*Lexer).buf)
    }
|   AUTO_INCREMENT equal_opt INTEGRAL
    {
        $$ = tree.NewTableOptionAutoIncrement(uint64($3.(int64)), yylex.(*Lexer).buf)
    }
|   AVG_ROW_LENGTH equal_opt INTEGRAL
    {
        $$ = tree.NewTableOptionAvgRowLength(uint64($3.(int64)), yylex.(*Lexer).buf)
    }
|   default_opt charset_keyword equal_opt charset_name
    {
        $$ = tree.NewTableOptionCharset($4, yylex.(*Lexer).buf)
    }
|   default_opt COLLATE equal_opt charset_name
    {
        $$ = tree.NewTableOptionCollate($4, yylex.(*Lexer).buf)
    }
|   CHECKSUM equal_opt INTEGRAL
    {
        $$ = tree.NewTableOptionChecksum(uint64($3.(int64)), yylex.(*Lexer).buf)
    }
|   COMMENT_KEYWORD equal_opt STRING
    {
    	str := util.DealCommentString($3)
        $$ = tree.NewTableOptionComment(str, yylex.(*Lexer).buf)
    }
|   COMPRESSION equal_opt STRING
    {
        $$ = tree.NewTableOptionCompression($3, yylex.(*Lexer).buf)
    }
|   CONNECTION equal_opt STRING
    {
        $$ = tree.NewTableOptionConnection($3, yylex.(*Lexer).buf)
    }
|   DATA DIRECTORY equal_opt STRING
    {
        $$ = tree.NewTableOptionDataDirectory($4, yylex.(*Lexer).buf)
    }
|   INDEX DIRECTORY equal_opt STRING
    {
        $$ = tree.NewTableOptionIndexDirectory($4, yylex.(*Lexer).buf)
    }
|   DELAY_KEY_WRITE equal_opt INTEGRAL
    {
        $$ = tree.NewTableOptionDelayKeyWrite(uint64($3.(int64)), yylex.(*Lexer).buf)
    }
|   ENCRYPTION equal_opt STRING
    {
        $$ = tree.NewTableOptionEncryption($3, yylex.(*Lexer).buf)
    }
|   ENGINE equal_opt table_alias
    {
        $$ = tree.NewTableOptionEngine($3, yylex.(*Lexer).buf)
    }
|   ENGINE_ATTRIBUTE equal_opt STRING
    {
        $$ = tree.NewTableOptionEngineAttr($3, yylex.(*Lexer).buf)
    }
|   INSERT_METHOD equal_opt insert_method_options
    {
        $$ = tree.NewTableOptionInsertMethod($3, yylex.(*Lexer).buf)
    }
|   KEY_BLOCK_SIZE equal_opt INTEGRAL
    {
        $$ = tree.NewTableOptionKeyBlockSize(uint64($3.(int64)), yylex.(*Lexer).buf)
    }
|   MAX_ROWS equal_opt INTEGRAL
    {
        $$ = tree.NewTableOptionMaxRows(uint64($3.(int64)), yylex.(*Lexer).buf)
    }
|   MIN_ROWS equal_opt INTEGRAL
    {
        $$ = tree.NewTableOptionMinRows(uint64($3.(int64)), yylex.(*Lexer).buf)
    }
|   PACK_KEYS equal_opt INTEGRAL
    {
        // Value Default
        $$ = tree.NewTableOptionPackKeys(
            $3.(int64),
            false,
            yylex.(*Lexer).buf,
        )
    }
|   PACK_KEYS equal_opt DEFAULT
    {
        // Value Default
        $$ = tree.NewTableOptionPackKeys(
            0,
            true,
            yylex.(*Lexer).buf,
        )
    }
|   PASSWORD equal_opt STRING
    {
        $$ = tree.NewTableOptionPassword($3, yylex.(*Lexer).buf)
    }
|   ROW_FORMAT equal_opt row_format_options
    {
        $$ = tree.NewTableOptionRowFormat($3, yylex.(*Lexer).buf)
    }
|   START TRANSACTION
    {
        $$ = tree.NewTableOptionStartTrans(true, yylex.(*Lexer).buf)
    }
|   SECONDARY_ENGINE_ATTRIBUTE equal_opt STRING
    {
        $$ = tree.NewTableOptionSecondaryEngineAttr($3, yylex.(*Lexer).buf)
    }
|   STATS_AUTO_RECALC equal_opt INTEGRAL
    {
        // Value Default
        $$ = tree.NewTableOptionStatsAutoRecalc(
            uint64($3.(int64)),
            false,
            yylex.(*Lexer).buf,
        )
    }
|   STATS_AUTO_RECALC equal_opt DEFAULT
    {
        // Value Default
        $$ = tree.NewTableOptionStatsAutoRecalc(
            0,
            true,
            yylex.(*Lexer).buf,
        )
    }
|   STATS_PERSISTENT equal_opt INTEGRAL
    {
        // Value Default
        $$ = tree.NewTableOptionStatsPersistent(
            uint64($3.(int64)),
            false,
            yylex.(*Lexer).buf,
        )
    }
|   STATS_PERSISTENT equal_opt DEFAULT
    {
        // Value Default
        $$ = tree.NewTableOptionStatsPersistent(
            0,
            true,
            yylex.(*Lexer).buf,
        )
    }
|   STATS_SAMPLE_PAGES equal_opt INTEGRAL
    {
        // Value Default
        $$ = tree.NewTableOptionStatsSamplePages(
            uint64($3.(int64)),
            false,
            yylex.(*Lexer).buf,
        )
    }
|   STATS_SAMPLE_PAGES equal_opt DEFAULT
    {
        // Value Default
        $$ = tree.NewTableOptionStatsSamplePages(
            0,
            true,
            yylex.(*Lexer).buf,
        )
    }
|   TABLESPACE equal_opt ident
    {
        $$= tree.NewTableOptionTablespace($3.Compare(), "", yylex.(*Lexer).buf)
    }
|   storage_opt
    {
        $$= tree.NewTableOptionTablespace("", $1, yylex.(*Lexer).buf)
    }
|   UNION equal_opt '(' table_name_list ')'
    {
        $$= tree.NewTableOptionUnion($4, yylex.(*Lexer).buf)
    }
|   PROPERTIES '(' properties_list ')'
    {
        $$ = tree.NewTableOptionProperties($3, yylex.(*Lexer).buf)
    }

properties_list:
    property_elem
    {
        $$ = buffer.MakeSlice[tree.Property](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Property](yylex.(*Lexer).buf, $$, $1)
    }
|    properties_list ',' property_elem
    {
        $$ = buffer.MakeSlice[tree.Property](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Property](yylex.(*Lexer).buf, $1, $3)
    }

property_elem:
    STRING '=' STRING
    {
        // key value
        $$ = *tree.NewProperty($1, $3, yylex.(*Lexer).buf)
    }

storage_opt:
    STORAGE DISK
    {
        $$ = " " + $1 + " " + $2
    }
|   STORAGE MEMORY
    {
        $$ = " " + $1 + " " + $2
    }

row_format_options:
    DEFAULT
    {
        $$ = tree.ROW_FORMAT_DEFAULT
    }
|   DYNAMIC
    {
        $$ = tree.ROW_FORMAT_DYNAMIC
    }
|   FIXED
    {
        $$ = tree.ROW_FORMAT_FIXED
    }
|   COMPRESSED
    {
        $$ = tree.ROW_FORMAT_COMPRESSED
    }
|   REDUNDANT
    {
        $$ = tree.ROW_FORMAT_REDUNDANT
    }
|   COMPACT
    {
        $$ = tree.ROW_FORMAT_COMPACT
    }

charset_name:
    name_string
|   BINARY

collate_name:
    name_string
|   BINARY

table_name_list:
    table_name
    {
        $$ = buffer.MakeSlice[*tree.TableName](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.TableName](yylex.(*Lexer).buf, $$, $1)
    }
|   table_name_list ',' table_name
    {
        $$ = buffer.MakeSlice[*tree.TableName](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.TableName](yylex.(*Lexer).buf, $1, $3)
    }

// Accepted patterns:
// <table>
// <schema>.<table>
table_name:
    ident
    {
        prefix := tree.NewObjectNamePrefix(yylex.(*Lexer).buf)
        prefix.ExplicitSchema = false
        $$ = tree.NewTableName(tree.Identifier($1.Compare()), *prefix, yylex.(*Lexer).buf)
    }
|   ident '.' ident
    {
        prefix := tree.NewObjectNamePrefix(yylex.(*Lexer).buf)
        bi := tree.NewBufIdentifier($1.Compare())
        yylex.(*Lexer).buf.Pin(bi)
        prefix.SchemaName = bi
        prefix.ExplicitSchema = true
        $$ = tree.NewTableName(tree.Identifier($3.Compare()), *prefix, yylex.(*Lexer).buf)
    }

table_elem_list_opt:
    {
        $$ = buffer.MakeSlice[tree.TableDef](yylex.(*Lexer).buf)
    }
|   table_elem_list

table_elem_list:
    table_elem
    {
        $$ = buffer.MakeSlice[tree.TableDef](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.TableDef](yylex.(*Lexer).buf, $$, $1)
    }
|   table_elem_list ',' table_elem
    {
        $$ = buffer.MakeSlice[tree.TableDef](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.TableDef](yylex.(*Lexer).buf, $1, $3)
    }

table_elem:
    column_def
    {
        tb := *buffer.Alloc[tree.TableDef](yylex.(*Lexer).buf)
        tb = $1
        $$ = tb
    }
|   constaint_def
    {
        $$ = $1
    }
|   index_def
    {
        $$ = $1
    }

table_elem_2:
    constaint_def
    {
        $$ = $1
    }
|   index_def
    {
        $$ = $1
    }

index_def:
    FULLTEXT key_or_index_opt index_name '(' index_column_list ')' index_option_list
    {
        // KeyParts Name Empty IndexOption
        $$ = tree.NewFullTextIndex(
            $5,
            $3,
            true,
            $7,
            yylex.(*Lexer).buf,
        )
    }
|   FULLTEXT key_or_index_opt index_name '(' index_column_list ')' USING index_type index_option_list
    {
        // KeyParts Name Empty IndexOption
        $$ = tree.NewFullTextIndex(
            $5,
            $3,
            true,
            $9,
            yylex.(*Lexer).buf,
        )
    }
|   key_or_index not_exists_opt index_name_and_type_opt '(' index_column_list ')' index_option_list
    {
        keyTyp := tree.INDEX_TYPE_INVALID
        if $3[1] != "" {
               t := strings.ToLower($3[1])
            switch t {
            case "zonemap":
                keyTyp = tree.INDEX_TYPE_ZONEMAP
            case "bsi":
                keyTyp = tree.INDEX_TYPE_BSI
            default:
                yylex.Error("Invail the type of index")
                return 1
            }
        }
        // IfNotExists KeyParts Name KeyType IndexOption
        $$ = tree.NewIndex(
            $2,
            $5,
            $3[0],
            keyTyp,
            $7,
            yylex.(*Lexer).buf,
        )
    }
|   key_or_index not_exists_opt index_name_and_type_opt '(' index_column_list ')' USING index_type index_option_list
    {
        keyTyp := tree.INDEX_TYPE_INVALID
        if $3[1] != "" {
               t := strings.ToLower($3[1])
            switch t {
            case "zonemap":
                keyTyp = tree.INDEX_TYPE_ZONEMAP
            case "bsi":
                keyTyp = tree.INDEX_TYPE_BSI
            default:
                yylex.Error("Invail the type of index")
                return 1
            }
        }
        // IfNotExists KeyParts Name KeyType IndexOption
        $$ = tree.NewIndex(
            $2,
            $5,
            $3[0],
            keyTyp,
            $9,
            yylex.(*Lexer).buf,
        )
    }

constaint_def:
    constraint_keyword constraint_elem
    {
        if $1 != "" {
            switch v := $2.(type) {
            case *tree.PrimaryKeyIndex:
                bc := tree.NewBufString($1)
                yylex.(*Lexer).buf.Pin(bc)
                v.ConstraintSymbol = bc
            case *tree.ForeignKey:
                bf := tree.NewBufString($1)
                yylex.(*Lexer).buf.Pin(bf)
                v.ConstraintSymbol = bf
            case *tree.UniqueIndex:
                bn := tree.NewBufString($1)
                yylex.(*Lexer).buf.Pin(bn)
                v.ConstraintSymbol = bn
            }
        }
        $$ = $2
    }
|   constraint_elem
    {
        $$ = $1
    }

constraint_elem:
    PRIMARY KEY index_name_and_type_opt '(' index_column_list ')' index_option_list
    {
        // KeyParts Name Empty IndexOption
        $$ = tree.NewPrimaryKeyIndex(
            $5,
            $3[0],
            $3[1] == "",
            $7,
            yylex.(*Lexer).buf,
        )
    }
|   PRIMARY KEY index_name_and_type_opt '(' index_column_list ')' USING index_type index_option_list
    {
        // KeyParts Name Empty IndexOption
        $$ = tree.NewPrimaryKeyIndex(
            $5,
            $3[0],
            $3[1] == "",
            $9,
            yylex.(*Lexer).buf,
        )
    }
|   UNIQUE key_or_index_opt index_name_and_type_opt '(' index_column_list ')' index_option_list
    {
        // KeyParts Name Empty IndexOption
        $$ = tree.NewUniqueIndex(
            $5,
            $3[0],
            $3[1] == "",
            $7,
            yylex.(*Lexer).buf,
        )
    }
|   UNIQUE key_or_index_opt index_name_and_type_opt '(' index_column_list ')' USING index_type index_option_list
    {
        // KeyParts Name Empty IndexOption
        $$ = tree.NewUniqueIndex(
            $5,
            $3[0],
            $3[1] == "",
            $9,
            yylex.(*Lexer).buf,
        )
    }
|   FOREIGN KEY not_exists_opt index_name '(' index_column_list ')' references_def
    {
        // ifNotExists keyParts Name Refer Empty
        $$ = tree.NewForeignKey(
            $3,
            $6,
            $4,
            $8,
            true,
            yylex.(*Lexer).buf,
        )
    }
|   CHECK '(' expression ')' enforce_opt
    {
        // Expr Enforced
        $$ = tree.NewCheckIndex(
            $3,
            $5,
            yylex.(*Lexer).buf,
        )
    }

enforce_opt:
    {
        $$ = false
    }
|    enforce

key_or_index_opt:
    {
        $$ = ""
    }
|   key_or_index
    {
        $$ = $1
    }

key_or_index:
    KEY
|   INDEX

index_name_and_type_opt:
    index_name
    {
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $$, yylex.(*Lexer).buf.CopyString($1), yylex.(*Lexer).buf.CopyString(""))
    }
|   index_name USING index_type
    {
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $$, yylex.(*Lexer).buf.CopyString($1), yylex.(*Lexer).buf.CopyString($3))
    }
|   ident TYPE index_type
    {
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $$, yylex.(*Lexer).buf.CopyString($1.Compare()), yylex.(*Lexer).buf.CopyString($3))
    }

index_type:
    BTREE
|   HASH
|   RTREE
|   ZONEMAP
|   BSI

insert_method_options:
    NO
|   FIRST
|   LAST

index_name:
    {
        $$ = ""
    }
|    ident
	{
		$$ = $1.Compare()
	}

column_def:
    column_name column_type column_attribute_list_opt
    {
        $$ = tree.NewColumnTableDef($1, $2, $3, yylex.(*Lexer).buf)
    }

column_name_unresolved:
    ident
    {
        $$ = tree.SetUnresolvedName(yylex.(*Lexer).buf, $1.Compare())
    }
|   ident '.' ident
    {
        $$ = tree.SetUnresolvedName(yylex.(*Lexer).buf, $1.Compare(), $3.Compare())
    }
|   ident '.' ident '.' ident
    {
        $$ = tree.SetUnresolvedName(yylex.(*Lexer).buf, $1.Compare(), $3.Compare(), $5.Compare())
    }

ident:
    ID
    {
		$$ = tree.NewCStr($1, yylex.(*Lexer).lower, yylex.(*Lexer).buf)
    }
|	QUOTE_ID
	{
    	$$ = tree.NewCStr($1, yylex.(*Lexer).lower, yylex.(*Lexer).buf)
    }
|   not_keyword
	{
    	$$ = tree.NewCStr($1, yylex.(*Lexer).lower, yylex.(*Lexer).buf)
    }
|   non_reserved_keyword
	{
    	$$ = tree.NewCStr($1, yylex.(*Lexer).lower, yylex.(*Lexer).buf)
    }

column_name:
    ident
    {
        $$ = tree.SetUnresolvedName(yylex.(*Lexer).buf, $1.Compare())
    }
|   ident '.' ident
    {
        $$ = tree.SetUnresolvedName(yylex.(*Lexer).buf, $1.Compare(), $3.Compare())
    }
|   ident '.' ident '.' ident
    {
        $$ = tree.SetUnresolvedName(yylex.(*Lexer).buf, $1.Compare(), $3.Compare(), $5.Compare())
    }

column_attribute_list_opt:
    {
        $$ = nil
    }
|   column_attribute_list
    {
        $$ = $1
    }

column_attribute_list:
    column_attribute_elem
    {
        $$ = buffer.MakeSlice[tree.ColumnAttribute](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.ColumnAttribute](yylex.(*Lexer).buf, $$, $1)
    }
|   column_attribute_list column_attribute_elem
    {
        $$ = buffer.MakeSlice[tree.ColumnAttribute](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.ColumnAttribute](yylex.(*Lexer).buf, $1, $2)
    }

column_attribute_elem:
    NULL
    {
        $$ = tree.NewAttributeNull(true, yylex.(*Lexer).buf)
    }
|   NOT NULL
    {
        $$ = tree.NewAttributeNull(false, yylex.(*Lexer).buf)
    }
|   DEFAULT bit_expr
    {
        $$ = tree.NewAttributeDefault($2, yylex.(*Lexer).buf)
    }
|   AUTO_INCREMENT
    {
        $$ = tree.NewAttributeAutoIncrement(yylex.(*Lexer).buf)
    }
|   keys
    {
        $$ = $1
    }
|   COMMENT_KEYWORD STRING
    {
    	str := util.DealCommentString($2)
        $$ = tree.NewAttributeComment(tree.NewNumValWithType(constant.MakeString(str), str, false, tree.P_char, yylex.(*Lexer).buf), yylex.(*Lexer).buf)
    }
|   COLLATE collate_name
    {
        $$ = tree.NewAttributeCollate($2, yylex.(*Lexer).buf)
    }
|   COLUMN_FORMAT column_format
    {
        $$ = tree.NewAttributeColumnFormat($2, yylex.(*Lexer).buf)
    }
|   SECONDARY_ENGINE_ATTRIBUTE '=' STRING
    {
        $$ = nil
    }
|   ENGINE_ATTRIBUTE '=' STRING
    {
        $$ = nil
    }
|   STORAGE storage_media
    {
        $$ = tree.NewAttributeStorage($2, yylex.(*Lexer).buf)
    }
|   AUTO_RANDOM field_length_opt
    {
        $$ = tree.NewAttributeAutoRandom(int($2), yylex.(*Lexer).buf)
   }
|   references_def
    {
        $$ = $1
    }
|   constraint_keyword_opt CHECK '(' expression ')'
    {
        $$ = tree.NewAttributeCheck($4, false, $1, yylex.(*Lexer).buf)
    }
|   constraint_keyword_opt CHECK '(' expression ')' enforce
    {
        $$ = tree.NewAttributeCheck($4, $6, $1, yylex.(*Lexer).buf)
    }
|   ON UPDATE name_datetime_scale datetime_scale_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($3))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        if $4 != nil {
            exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $4)
        }

        expr := tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
        $$ = tree.NewAttributeOnUpdate(expr, yylex.(*Lexer).buf)
    }
|   LOW_CARDINALITY
    {
	    $$ = tree.NewAttributeLowCardinality(yylex.(*Lexer).buf)
    }
|   VISIBLE
    {
        $$ = tree.NewAttributeVisable(true, yylex.(*Lexer).buf)
    }
|   INVISIBLE
    {
        $$ = tree.NewAttributeVisable(false, yylex.(*Lexer).buf)
    }
|   default_opt CHARACTER SET equal_opt ident
    {
        $$ = nil
    }
|	HEADER '(' STRING ')'
	{
		$$ = tree.NewAttributeHeader($3, yylex.(*Lexer).buf)
	}
|	HEADERS
	{
		$$ = tree.NewAttributeHeaders(yylex.(*Lexer).buf)
	}

enforce:
    ENFORCED
    {
        $$ = true
    }
|   NOT ENFORCED
    {
        $$ = false
    }

constraint_keyword_opt:
    {
        $$ = ""
    }
 |    constraint_keyword
     {
         $$ = $1
     }

constraint_keyword:
    CONSTRAINT
    {
        $$ = ""
    }
|   CONSTRAINT ident
    {
        $$ = $2.Compare()
    }

references_def:
    REFERENCES table_name index_column_list_opt match_opt on_delete_update_opt
    {
        // TableName KeyPart MatchType ReferenceOptionType ReferenceOptionType
        $$ = tree.NewAttributeReference(
            $2,
            $3,
            $4,
            $5.OnDelete,
            $5.OnUpdate,
            yylex.(*Lexer).buf,
        )
    }

on_delete_update_opt:
    %prec LOWER_THAN_ON
    {
        $$ = tree.NewReferenceOnRecord( 
            tree.REFERENCE_OPTION_INVALID,
            tree.REFERENCE_OPTION_INVALID,
            yylex.(*Lexer).buf,
        )
    }
|   on_delete %prec LOWER_THAN_ON
    {
        $$ = tree.NewReferenceOnRecord( 
            $1,
            tree.REFERENCE_OPTION_INVALID,
            yylex.(*Lexer).buf,
        )
    }
|   on_update %prec LOWER_THAN_ON
    {
        $$ = tree.NewReferenceOnRecord( 
            tree.REFERENCE_OPTION_INVALID,
            $1,
            yylex.(*Lexer).buf,
        )
    }
|   on_delete on_update
    {
        $$ = tree.NewReferenceOnRecord( 
            $1,
            $2,
            yylex.(*Lexer).buf,
        )
    }
|   on_update on_delete
    {
        $$ = tree.NewReferenceOnRecord( 
            $2,
            $1,
            yylex.(*Lexer).buf,
        )
    }

on_delete:
    ON DELETE ref_opt
    {
        $$ = $3
    }

on_update:
    ON UPDATE ref_opt
    {
        $$ = $3
    }

ref_opt:
    RESTRICT
    {
        $$ = tree.REFERENCE_OPTION_RESTRICT
    }
|   CASCADE
    {
        $$ = tree.REFERENCE_OPTION_CASCADE
    }
|   SET NULL
    {
        $$ = tree.REFERENCE_OPTION_SET_NULL
    }
|   NO ACTION
    {
        $$ = tree.REFERENCE_OPTION_NO_ACTION
    }
|   SET DEFAULT
    {
        $$ = tree.REFERENCE_OPTION_SET_DEFAULT
    }

match_opt:
    {
        $$ = tree.MATCH_INVALID
    }
|   match

match:
    MATCH FULL
    {
        $$ = tree.MATCH_FULL
    }
|   MATCH PARTIAL
    {
        $$ = tree.MATCH_PARTIAL
    }
|   MATCH SIMPLE
    {
        $$ = tree.MATCH_SIMPLE
    }

index_column_list_opt:
    {
        $$ = nil
    }
|   '(' index_column_list ')'
    {
        $$ = $2
    }

field_length_opt:
    {
        $$ = -1
    }
|   '(' INTEGRAL ')'
    {
        $$ = $2.(int64)
    }

storage_media:
    DEFAULT
|   DISK
|   MEMORY

column_format:
    DEFAULT
|   FIXED
|   DYNAMIC

subquery:
    select_with_parens %prec SUBQUERY_AS_EXPR
    {
        $$ = tree.NewSubquery(
            $1, 
            false,
            yylex.(*Lexer).buf,
        )
    }

bit_expr:
    bit_expr '&' bit_expr %prec '&'
    {
        $$ = tree.NewBinaryExpr(tree.BIT_AND, $1, $3, yylex.(*Lexer).buf)
    }
|   bit_expr '|' bit_expr %prec '|'
    {
        $$ = tree.NewBinaryExpr(tree.BIT_OR, $1, $3, yylex.(*Lexer).buf)
    }
|   bit_expr '^' bit_expr %prec '^'
    {
        $$ = tree.NewBinaryExpr(tree.BIT_XOR, $1, $3, yylex.(*Lexer).buf)
    }
|   bit_expr '+' bit_expr %prec '+'
    {
        $$ = tree.NewBinaryExpr(tree.PLUS, $1, $3, yylex.(*Lexer).buf)
    }
|   bit_expr '-' bit_expr %prec '-'
    {
        $$ = tree.NewBinaryExpr(tree.MINUS, $1, $3, yylex.(*Lexer).buf)
    }
|   bit_expr '*' bit_expr %prec '*'
    {
        $$ = tree.NewBinaryExpr(tree.MULTI, $1, $3, yylex.(*Lexer).buf)
    }
|   bit_expr '/' bit_expr %prec '/'
    {
        $$ = tree.NewBinaryExpr(tree.DIV, $1, $3, yylex.(*Lexer).buf)
    }
|   bit_expr DIV bit_expr %prec DIV
    {
        $$ = tree.NewBinaryExpr(tree.INTEGER_DIV, $1, $3, yylex.(*Lexer).buf)
    }
|   bit_expr '%' bit_expr %prec '%'
    {
        $$ = tree.NewBinaryExpr(tree.MOD, $1, $3, yylex.(*Lexer).buf)
    }
|   bit_expr MOD bit_expr %prec MOD
    {
        $$ = tree.NewBinaryExpr(tree.MOD, $1, $3, yylex.(*Lexer).buf)
    }
|   bit_expr SHIFT_LEFT bit_expr %prec SHIFT_LEFT
    {
        $$ = tree.NewBinaryExpr(tree.LEFT_SHIFT, $1, $3, yylex.(*Lexer).buf)
    }
|   bit_expr SHIFT_RIGHT bit_expr %prec SHIFT_RIGHT
    {
        $$ = tree.NewBinaryExpr(tree.RIGHT_SHIFT, $1, $3, yylex.(*Lexer).buf)
    }
|   simple_expr
    {
        $$ = $1
    }

simple_expr:
    normal_ident
    {
        $$ = $1
    }
|   variable
    {
        $$ = $1
    }
|   literal
    {
        $$ = $1
    }
|   '(' expression ')'
    {
        $$ = tree.NewParenExpr($2, yylex.(*Lexer).buf)
    }
|   '(' expression_list ',' expression ')'
    {
        exprs := buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, $2, $4)
        $$ = tree.NewTuple(exprs, yylex.(*Lexer).buf)
    }
|   '+'  simple_expr %prec UNARY
    {
        $$ = tree.NewUnaryExpr(tree.UNARY_PLUS, $2, yylex.(*Lexer).buf)
    }
|   '-'  simple_expr %prec UNARY
    {
        $$ = tree.NewUnaryExpr(tree.UNARY_MINUS, $2, yylex.(*Lexer).buf)
    }
|   '~'  simple_expr
    {
        $$ = tree.NewUnaryExpr(tree.UNARY_TILDE, $2, yylex.(*Lexer).buf)
    }
|   '!' simple_expr %prec UNARY
    {
        $$ = tree.NewUnaryExpr(tree.UNARY_MARK, $2, yylex.(*Lexer).buf)
    }
|   interval_expr
    {
        $$ = $1
    }
|   subquery
    {
        $$ = $1
    }
|   EXISTS subquery
    {
        $2.Exists = true
        $$ = $2
    }
|   CASE expression_opt when_clause_list else_opt END
    {
        // Exprs Whens Else
        $$ = tree.NewCaseExpr(
            $2,
            $3,
            $4,
            yylex.(*Lexer).buf,
        )
    }
|   CAST '(' expression AS mo_cast_type ')'
    {
        $$ = tree.NewCastExpr($3, $5, yylex.(*Lexer).buf)
    }
|   BIT_CAST '(' expression AS mo_cast_type ')'
    {
        $$ = tree.NewBitCastExpr($3, $5, yylex.(*Lexer).buf)
    }
|   CONVERT '(' expression ',' mysql_cast_type ')'
    {
        $$ = tree.NewCastExpr($3, $5, yylex.(*Lexer).buf)
    }
|   CONVERT '(' expression USING charset_name ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "convert")
        num := tree.NewNumValWithType(constant.MakeString($5), $5, false, tree.P_char, yylex.(*Lexer).buf)

        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $3, num)

        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }
|   function_call_generic
    {
        $$ = $1
    }
|   function_call_keyword
    {
        $$ = $1
    }
|   function_call_nonkeyword
    {
        $$ = $1
    }
|   function_call_aggregate
    {
        $$ = $1
    }
|   function_call_window
    {
        $$ = $1
    }

function_call_window:
	RANK '(' ')' window_spec
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        $$ = tree.NewFuncExprWithWinSpec(tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf), nil, 0, $4, yylex.(*Lexer).buf)
    }
|	ROW_NUMBER '(' ')' window_spec
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        $$ = tree.NewFuncExprWithWinSpec(tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf), nil, 0, $4, yylex.(*Lexer).buf)
    }
|	DENSE_RANK '(' ')' window_spec
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        $$ = tree.NewFuncExprWithWinSpec(tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf), nil, 0, $4, yylex.(*Lexer).buf)
    }

else_opt:
    {
        $$ = nil
    }
|    ELSE expression
    {
        $$ = $2
    }

expression_opt:
    {
        $$ = nil
    }
|    expression
    {
        $$ = $1
    }

when_clause_list:
    when_clause
    {
        $$ = buffer.MakeSlice[*tree.When](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.When](yylex.(*Lexer).buf, $$, $1)
    }
|    when_clause_list when_clause
    {
        $$ = buffer.MakeSlice[*tree.When](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[*tree.When](yylex.(*Lexer).buf, $1, $2)
    }

when_clause:
    WHEN expression THEN expression
    {
        // Cond Val
        $$ = tree.NewWhen(
            $2,
            $4,
            yylex.(*Lexer).buf,
        )
    }

mo_cast_type:
    column_type
{
   t := $$
   str := strings.ToLower(t.InternalType.FamilyString.Get())
   if str == "binary" {
        t.InternalType.Scale = -1
   } else if str == "char" {
   	if t.InternalType.DisplayWith == -1 {
   		bf := tree.NewBufString("varchar")
        yylex.(*Lexer).buf.Pin(bf)
   		t.InternalType.FamilyString = bf
   		t.InternalType.Oid = uint32(defines.MYSQL_TYPE_VARCHAR)
   	}
   }
}
|   SIGNED integer_opt
    {
        name := $1
        if $2 != "" {
            name = $2
        }

        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, name, 64, "", uint32(defines.MYSQL_TYPE_LONGLONG), yylex.(*Lexer).buf)
    }
|   UNSIGNED integer_opt
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, $2, 64, "", uint32(defines.MYSQL_TYPE_LONGLONG), yylex.(*Lexer).buf)
        $$.InternalType.Unsigned = true
    }

mysql_cast_type:
    decimal_type
|   BINARY length_option_opt
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.StringFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_VARCHAR), yylex.(*Lexer).buf)
        $$.InternalType.DisplayWith = $2
    }
|   CHAR length_option_opt
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.StringFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_VARCHAR), yylex.(*Lexer).buf)
        $$.InternalType.DisplayWith = $2
    }
|   DATE
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.DateFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_DATE), yylex.(*Lexer).buf)
    }
|   YEAR length_opt
    {
        // Family FamilyString Width Locale Oid DisplayWith Scale
        $$ = tree.NewSQLTypeWithDisplayScale(tree.IntFamily, $1, 16, "", uint32(defines.MYSQL_TYPE_YEAR), $2, 0, yylex.(*Lexer).buf)
    }
|   DATETIME timestamp_option_opt
    {
        // Family FamilyString Width Locale Oid DisplayWith Scale
        $$ = tree.NewSQLTypeWithDisplayScale(tree.TimestampFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_YEAR), $2, $2, yylex.(*Lexer).buf)
    }
|   TIME length_opt
    {
        // Family FamilyString Width Locale Oid DisplayWith Scale
        $$ = tree.NewSQLTypeWithDisplayScale(tree.TimeFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_TIME), $2, $2, yylex.(*Lexer).buf)
    }
|   SIGNED integer_opt
    {
        name := $1
        if $2 != "" {
            name = $2
        }

        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, name, 64, "", uint32(defines.MYSQL_TYPE_LONGLONG), yylex.(*Lexer).buf)
    }
|   UNSIGNED integer_opt
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, $2, 64, "", uint32(defines.MYSQL_TYPE_LONGLONG), yylex.(*Lexer).buf)
        $$.InternalType.Unsigned = true
    }

integer_opt:
    {}
|    INTEGER
|    INT

frame_bound:
	frame_bound_start
|   UNBOUNDED FOLLOWING
    {
        // Type UnBounded Expr 
        $$ = tree.NewFrameBound(
            tree.Following,
            true,
            nil, 
            yylex.(*Lexer).buf,
        )
    }
|   num_literal FOLLOWING
    {
        // Type UnBounded Expr 
        $$ = tree.NewFrameBound(
            tree.Following, 
            false,
            $1,
            yylex.(*Lexer).buf,
        )
    }
|	interval_expr FOLLOWING
	{
        // Type UnBounded Expr 
        $$ = tree.NewFrameBound(
            tree.Following, 
            false,
            $1,
            yylex.(*Lexer).buf,
        )
	}

frame_bound_start:
    CURRENT ROW
    {
        // Type UnBounded Expr 
        $$ = tree.NewFrameBound(
            tree.CurrentRow,
            false,
            nil, 
            yylex.(*Lexer).buf,
        )
    }
|   UNBOUNDED PRECEDING
    {
        // Type UnBounded Expr 
        $$ = tree.NewFrameBound(
            tree.Preceding, 
            true,
            nil, 
            yylex.(*Lexer).buf,
        )
    }
|   num_literal PRECEDING
    {
        // Type UnBounded Expr 
        $$ = tree.NewFrameBound(
            tree.Preceding, 
            false,
            $1,
            yylex.(*Lexer).buf,
        )
    }
|	interval_expr PRECEDING
	{
        // Type UnBounded Expr 
		$$ = tree.NewFrameBound(
            tree.Preceding, 
            false,
            $1,
            yylex.(*Lexer).buf,
        )
	}

frame_type:
    ROWS
    {
        $$ = tree.Rows
    }
|   RANGE
    {
        $$ = tree.Range
    }
|   GROUPS
    {
        $$ = tree.Groups
    }

window_frame_clause:
    frame_type frame_bound_start
    {
        // Type HasEnd Start End
        $$ = tree.NewFrameClause(
            $1,
            false,
            $2,
            tree.NewFrameBound(tree.CurrentRow, false, nil, yylex.(*Lexer).buf),
            yylex.(*Lexer).buf,
        )
    }
|   frame_type BETWEEN frame_bound AND frame_bound
    {
        // Type HasEnd Start End
        $$ = tree.NewFrameClause(
            $1,
            true,
            $3,
            $5,
            yylex.(*Lexer).buf,
        )
    }

window_frame_clause_opt:
    {
        $$ = nil
    }
|   window_frame_clause
    {
        $$ = $1
    }


window_partition_by:
   PARTITION BY expression_list
    {
        $$ = $3
    }

window_partition_by_opt:
    {
        $$ = nil
    }
|   window_partition_by
    {
        $$ = $1
    }

separator_opt:
    {
        $$ = ","
    }
|   SEPARATOR STRING
    {
       $$ = $2
    }

window_spec_opt:
    {
        $$ = nil
    }
|	window_spec

window_spec:
    OVER '(' window_partition_by_opt order_by_opt window_frame_clause_opt ')'
    {
    	hasFrame := true
    	f := buffer.Alloc[tree.FrameClause](yylex.(*Lexer).buf)
    	if $5 != nil {
    		f = $5
    	} else {
    		hasFrame = false
            // Type HasEnd Start End
    		f = tree.NewFrameClause(
                    tree.Range,
                    false,
                    nil,
                    nil,
                    yylex.(*Lexer).buf,
                )
    		if $4 == nil {
                // BoundType UnBounded
				f.Start = tree.NewFrameBound(
                    tree.Preceding,
                    true,
                    nil,
                    yylex.(*Lexer).buf,
                )
                // BoundType UnBounded
                f.End = tree.NewFrameBound(
                    tree.Following, 
                    true,
                    nil,
                    yylex.(*Lexer).buf,
                )
    		} else {
                // BoundType UnBounded
    			f.Start = tree.NewFrameBound(
                    tree.Preceding, 
                    true,
                    nil,
                    yylex.(*Lexer).buf,
                )
                // BoundType UnBounded
            	f.End = tree.NewFrameBound(
                    tree.CurrentRow,
                    false,
                    nil,
                    yylex.(*Lexer).buf,
                )
    		}
    	}
        // PartitionBy OrderBy FlameClause HasFrame
        $$ = tree.NewWindowSpec(
            $3,
            $4,
            f,
            hasFrame,
            yylex.(*Lexer).buf,
        )
    }

function_call_aggregate:
    GROUP_CONCAT '(' func_type_opt expression_list order_by_opt separator_opt ')' window_spec_opt
    {
	    name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        
        exprs := buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, $4, tree.NewNumValWithType(constant.MakeString($6), $6, false, tree.P_char, yylex.(*Lexer).buf))
        
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec(
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            $3,
            $8,
           yylex.(*Lexer).buf,
        )
        $$.OrderBy = $5
    }
|   AVG '(' func_type_opt expression  ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $4)

        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            $3,
            $6,
            yylex.(*Lexer).buf,
        )
    }
|   APPROX_COUNT '(' func_type_opt expression_list ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec(
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            $4,
            $3,
            $6,
            yylex.(*Lexer).buf,
        )
    }
|   APPROX_COUNT '(' '*' ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        es := tree.NewNumValWithType(constant.MakeString("*"), "*", false, tree.P_char, yylex.(*Lexer).buf)
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, es)
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            0,
            $5,
            yylex.(*Lexer).buf,
        )
    }
|   APPROX_COUNT_DISTINCT '(' expression_list ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            $3,
            0,
            $5,
            yylex.(*Lexer).buf,
        )
    }
|   APPROX_PERCENTILE '(' expression_list ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            $3,
            0,
            $5,
            yylex.(*Lexer).buf,
        )
    }
|   BIT_AND '(' func_type_opt expression ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $4)
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            $3,
            $6,
            yylex.(*Lexer).buf,
        )
    }
|   BIT_OR '(' func_type_opt expression ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $4)
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            $3,
            $6,
            yylex.(*Lexer).buf,
        )
    }
|   BIT_XOR '(' func_type_opt expression ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $4)
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            $3,
            $6,
            yylex.(*Lexer).buf,
        )
    }
|   COUNT '(' func_type_opt expression_list ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            $4,
            $3,
            $6,
            yylex.(*Lexer).buf,
        )
    }
|   COUNT '(' '*' ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        es := tree.NewNumValWithType(constant.MakeString("*"), "*", false, tree.P_char, yylex.(*Lexer).buf)
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, es)
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            0,
            $5,
            yylex.(*Lexer).buf,
        )
    }
|   MAX '(' func_type_opt expression ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $4)
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            $3,
            $6,
            yylex.(*Lexer).buf,
        )
    }
|   MIN '(' func_type_opt expression ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $4)
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            $3,
            $6,
            yylex.(*Lexer).buf,
        )
    }
|   SUM '(' func_type_opt expression ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $4)
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            $3,
            $6,
            yylex.(*Lexer).buf,
        )
    }
|   std_dev_pop '(' func_type_opt expression ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $4)
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            $3,
            $6,
            yylex.(*Lexer).buf,
        )
    }
|   STDDEV_SAMP '(' func_type_opt expression ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $4)
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            $3,
            $6,
            yylex.(*Lexer).buf,
        )
    }
|   VAR_POP '(' func_type_opt expression ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $4)
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            $3,
            $6,
            yylex.(*Lexer).buf,
        )
    }
|   VAR_SAMP '(' func_type_opt expression ')' window_spec_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $4)
        // Func Exprs Type WindowsSpec
        $$ = tree.NewFuncExprWithWinSpec( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            $3,
            $6,
            yylex.(*Lexer).buf,
        )
    }
|   MEDIAN '(' func_type_opt expression ')' window_spec_opt
    {
    	name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $4)
        // Func Exprs Type WindowsSpec
	    $$ = tree.NewFuncExprWithWinSpec( 
	        tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
	        exprs,
	        $3,
	        $6,
            yylex.(*Lexer).buf,
    	)
    }

std_dev_pop:
    STD
|   STDDEV
|   STDDEV_POP

function_call_generic:
    ID '(' expression_list_opt ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   substr_option '(' expression_list_opt ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   substr_option '(' expression FROM expression ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $3, $5)

        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }
|   substr_option '(' expression FROM expression FOR expression ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $3, $5, $7)
        
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }
|   EXTRACT '(' time_unit FROM expression ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        str := strings.ToLower($3)
        timeUinit := tree.NewNumValWithType(constant.MakeString(str), str, false, tree.P_char, yylex.(*Lexer).buf)

        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, timeUinit, $5)

        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }
|   func_not_keyword '(' expression_list_opt ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   VARIANCE '(' func_type_opt expression ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $4)
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
        $$.Type = $3
    }
|   NEXTVAL '(' expression_list ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "nextval")
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   SETVAL '(' expression_list  ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "setval")
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   CURRVAL '(' expression_list  ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "currval")
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   LASTVAL '('')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "lastval")
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   TRIM '(' expression ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        arg0 := tree.NewNumValWithType(constant.MakeInt64(0), "0", false, tree.P_int64, yylex.(*Lexer).buf)
        arg1 := tree.NewNumValWithType(constant.MakeString("both"), "both", false, tree.P_char, yylex.(*Lexer).buf)
        arg2 := tree.NewNumValWithType(constant.MakeString(" "), " ", false, tree.P_char, yylex.(*Lexer).buf)

        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, arg0, arg1, arg2, $3)
                
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }
|   TRIM '(' expression FROM expression ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        arg0 := tree.NewNumValWithType(constant.MakeInt64(1), "1", false, tree.P_int64, yylex.(*Lexer).buf)
        arg1 := tree.NewNumValWithType(constant.MakeString("both"), "both", false, tree.P_char, yylex.(*Lexer).buf)
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, arg0, arg1, $3, $5)

        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }
|   TRIM '(' trim_direction FROM expression ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        arg0 := tree.NewNumValWithType(constant.MakeInt64(2), "2", false, tree.P_int64, yylex.(*Lexer).buf)
        str := strings.ToLower($3)
        arg1 := tree.NewNumValWithType(constant.MakeString(str), str, false, tree.P_char, yylex.(*Lexer).buf)
        arg2 := tree.NewNumValWithType(constant.MakeString(" "), " ", false, tree.P_char, yylex.(*Lexer).buf)
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, arg0, arg1, arg2, $5)
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }
|   TRIM '(' trim_direction expression FROM expression ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        arg0 := tree.NewNumValWithType(constant.MakeInt64(3), "3", false, tree.P_int64, yylex.(*Lexer).buf)
        str := strings.ToLower($3)
        arg1 := tree.NewNumValWithType(constant.MakeString(str), str, false, tree.P_char, yylex.(*Lexer).buf)
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, arg0, arg1, $4, $6)
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }
|   VALUES '(' insert_column ')'
    {
        column := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($3))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, column)
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
    	$$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }


trim_direction:
    BOTH
|   LEADING
|   TRAILING

substr_option:
    SUBSTRING
|   SUBSTR
|   MID

time_unit:
    time_stamp_unit
    {
        $$ = $1
    }
|    SECOND_MICROSECOND
|    MINUTE_MICROSECOND
|    MINUTE_SECOND
|    HOUR_MICROSECOND
|    HOUR_SECOND
|    HOUR_MINUTE
|    DAY_MICROSECOND
|    DAY_SECOND
|    DAY_MINUTE
|    DAY_HOUR
|    YEAR_MONTH

time_stamp_unit:
    MICROSECOND
|    SECOND
|    MINUTE
|    HOUR
|    DAY
|    WEEK
|    MONTH
|    QUARTER
|    YEAR
|    SQL_TSI_SECOND
|    SQL_TSI_MINUTE
|    SQL_TSI_HOUR
|    SQL_TSI_DAY
|    SQL_TSI_WEEK
|    SQL_TSI_MONTH
|    SQL_TSI_QUARTER
|    SQL_TSI_YEAR

function_call_nonkeyword:
    CURTIME datetime_scale
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        es := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        if $2 != nil {
            es = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, es, $2)
        }

        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            es,
            yylex.(*Lexer).buf,
        )
    }
|   SYSDATE datetime_scale
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        es := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        if $2 != nil {
            es = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, es, $2)
        }

        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            es,
            yylex.(*Lexer).buf,
        )
    }
|	TIMESTAMPDIFF '(' time_stamp_unit ',' expression ',' expression ')'
	{
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        str := strings.ToLower($3)
        arg1 := tree.NewNumValWithType(constant.MakeString(str), str, false, tree.P_char, yylex.(*Lexer).buf)
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, arg1, $5, $7)
		$$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
	}
function_call_keyword:
    name_confict '(' expression_list_opt ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   name_braces braces_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            nil,
            yylex.(*Lexer).buf,
        )
    }
|    SCHEMA '('')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            nil,
            yylex.(*Lexer).buf,
        )
    }
|   name_datetime_scale datetime_scale_opt
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower($1))
        es := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        if $2 != nil {
            es = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, es, $2)
        }

        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            es,
            yylex.(*Lexer).buf,
        )
    }
|   BINARY '(' expression_list ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "binary")
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   BINARY literal
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "binary")
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $2)
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }
|   BINARY column_name
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "binary")
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $2)
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }
|   CHAR '(' expression_list ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "char")
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   CHAR '(' expression_list USING charset_name ')'
    {
        cn := tree.NewNumValWithType(constant.MakeString($5), $5, false, tree.P_char, yylex.(*Lexer).buf)
        es := $3
        es = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, es, cn)
        
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "char")
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            es,
            yylex.(*Lexer).buf,
        )
    }
|   DATE STRING
    {
        val := tree.NewNumValWithType(constant.MakeString($2), $2, false, tree.P_char, yylex.(*Lexer).buf)
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "date")
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, val)
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }
|   TIME STRING
    {
        val := tree.NewNumValWithType(constant.MakeString($2), $2, false, tree.P_char, yylex.(*Lexer).buf)
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "time")
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, val)
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }
|   INSERT '(' expression_list_opt ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "insert")
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   MOD '(' bit_expr ',' bit_expr ')'
    {
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $3, $5)
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "mod")
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }
|   PASSWORD '(' expression_list_opt ')'
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "password")
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            $3,
            yylex.(*Lexer).buf,
        )
    }
|   TIMESTAMP STRING
    {
        val := tree.NewNumValWithType(constant.MakeString($2), $2, false, tree.P_char, yylex.(*Lexer).buf)
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, val)
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "timestamp")
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }

datetime_scale_opt:
    {
        $$ = nil
    }
|   datetime_scale
    {
        $$ = $1
    }

datetime_scale:
   '(' ')'
    {
        $$ = nil
    }
|   '(' INTEGRAL ')'
    {
        ival, errStr := util.GetInt64($2)
        if errStr != "" {
            yylex.Error(errStr)
            return 1
        }
        str := fmt.Sprintf("%v", $2)
        $$ = tree.NewNumValWithType(constant.MakeInt64(ival), str, false, tree.P_int64, yylex.(*Lexer).buf)
    }

name_datetime_scale:
    CURRENT_TIME
|   CURRENT_TIMESTAMP
|   LOCALTIME
|   LOCALTIMESTAMP
|   UTC_TIME
|   UTC_TIMESTAMP

braces_opt:
    {}
|   '(' ')'
    {}

name_braces:
    CURRENT_USER
|   CURRENT_DATE
|   CURRENT_ROLE
|   UTC_DATE

name_confict:
    ASCII
|   CHARSET
|   COALESCE
|   COLLATION
|   DATE
|   DATABASE
|   DAY
|   HOUR
|   IF
|   INTERVAL
|   FORMAT
|   LEFT
|   MICROSECOND
|   MINUTE
|   MONTH
|   QUARTER
|   REPEAT
|   REPLACE
|   REVERSE
|   RIGHT
|   ROW_COUNT
|   SECOND
|   TIME
|   TIMESTAMP
|   TRUNCATE
|   USER
|   WEEK
|   YEAR
|   UUID

interval_expr:
    INTERVAL expression time_unit
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, "interval")
        str := strings.ToLower($3)
        arg2 := tree.NewNumValWithType(constant.MakeString(str), str, false, tree.P_char, yylex.(*Lexer).buf)
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $2, arg2)
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }

func_type_opt:
    {
        $$ = tree.FUNC_TYPE_DEFAULT
    }
|   DISTINCT
    {
        $$ = tree.FUNC_TYPE_DISTINCT
    }
|   ALL
    {
        $$ = tree.FUNC_TYPE_ALL
    }

tuple_expression:
    '(' expression_list ')'
    {
        $$ = tree.NewTuple($2, yylex.(*Lexer).buf)
    }

expression_list_opt:
    {
        $$ = nil
    }
|   expression_list
    {
        $$ = $1
    }

expression_list:
    expression
    {
        $$ = buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, $$, $1)
    }
|   expression_list ',' expression
    {
        $$ = buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, $1, $3)
    }

// See https://dev.mysql.com/doc/refman/8.0/en/expressions.html
expression:
    expression AND expression %prec AND
    {
        $$ = tree.NewAndExpr($1, $3, yylex.(*Lexer).buf)
    }
|   expression OR expression %prec OR
    {
        $$ = tree.NewOrExpr($1, $3, yylex.(*Lexer).buf)
    }
|   expression PIPE_CONCAT expression %prec PIPE_CONCAT
    {
        name := tree.SetUnresolvedName(yylex.(*Lexer).buf, strings.ToLower("concat"))
        exprs := buffer.MakeSlice[tree.Expr](yylex.(*Lexer).buf)
        exprs = buffer.AppendSlice[tree.Expr](yylex.(*Lexer).buf, exprs, $1, $3)
        
        $$ = tree.NewFuncExpr( 
            tree.FuncName2ResolvableFunctionReference(name, yylex.(*Lexer).buf),
            exprs,
            yylex.(*Lexer).buf,
        )
    }
|   expression XOR expression %prec XOR
    {
        $$ = tree.NewXorExpr($1, $3, yylex.(*Lexer).buf)
    }
|   NOT expression %prec NOT
    {
        $$ = tree.NewNotExpr($2, yylex.(*Lexer).buf)
    }
|   MAXVALUE
    {
        $$ = tree.NewMaxValue(yylex.(*Lexer).buf)
    }
|   boolean_primary
    {
        $$ = $1
    }

boolean_primary:
    boolean_primary IS NULL %prec IS
    {
        $$ = tree.NewIsNullExpr($1, yylex.(*Lexer).buf)
    }
|   boolean_primary IS NOT NULL %prec IS
    {
        $$ = tree.NewIsNotNullExpr($1, yylex.(*Lexer).buf)
    }
|   boolean_primary IS UNKNOWN %prec IS
    {
        $$ = tree.NewIsUnknownExpr($1, yylex.(*Lexer).buf)
    }
|   boolean_primary IS NOT UNKNOWN %prec IS
    {
        $$ = tree.NewIsNotUnknownExpr($1, yylex.(*Lexer).buf)
    }
|   boolean_primary IS TRUE %prec IS
    {
        $$ = tree.NewIsTrueExpr($1, yylex.(*Lexer).buf)
    }
|   boolean_primary IS NOT TRUE %prec IS
    {
        $$ = tree.NewIsNotTrueExpr($1, yylex.(*Lexer).buf)
    }
|   boolean_primary IS FALSE %prec IS
    {
        $$ = tree.NewIsFalseExpr($1, yylex.(*Lexer).buf)
    }
|   boolean_primary IS NOT FALSE %prec IS
    {
        $$ = tree.NewIsNotFalseExpr($1, yylex.(*Lexer).buf)
    }
|   boolean_primary comparison_operator predicate %prec '='
    {
        $$ = tree.NewComparisonExpr($2, $1, $3, yylex.(*Lexer).buf)
    }
|   boolean_primary comparison_operator and_or_some subquery %prec '='
    {
        $$ = tree.NewSubqueryComparisonExpr($2, $3, $1, $4, yylex.(*Lexer).buf)
        $$ = tree.NewSubqueryComparisonExpr($2, $3, $1, $4, yylex.(*Lexer).buf)
    }
|   predicate

predicate:
    bit_expr IN col_tuple
    {
        $$ = tree.NewComparisonExpr(tree.IN, $1, $3, yylex.(*Lexer).buf)
    }
|   bit_expr NOT IN col_tuple
    {
        $$ = tree.NewComparisonExpr(tree.NOT_IN, $1, $4, yylex.(*Lexer).buf)
    }
|   bit_expr LIKE simple_expr like_escape_opt
    {
        $$ = tree.NewComparisonExprWithEscape(tree.LIKE, $1, $3, $4, yylex.(*Lexer).buf)
    }
|   bit_expr NOT LIKE simple_expr like_escape_opt
    {
        $$ = tree.NewComparisonExprWithEscape(tree.NOT_LIKE, $1, $4, $5, yylex.(*Lexer).buf)
    }
|   bit_expr ILIKE simple_expr like_escape_opt
    {
        $$ = tree.NewComparisonExprWithEscape(tree.ILIKE, $1, $3, $4, yylex.(*Lexer).buf)
    }
|   bit_expr NOT ILIKE simple_expr like_escape_opt
    {
        $$ = tree.NewComparisonExprWithEscape(tree.NOT_ILIKE, $1, $4, $5, yylex.(*Lexer).buf)
    }
|   bit_expr REGEXP bit_expr
    {
        $$ = tree.NewComparisonExpr(tree.REG_MATCH, $1, $3, yylex.(*Lexer).buf)
    }
|   bit_expr NOT REGEXP bit_expr
    {
        $$ = tree.NewComparisonExpr(tree.NOT_REG_MATCH, $1, $4, yylex.(*Lexer).buf)
    }
|   bit_expr BETWEEN bit_expr AND predicate
    {
        $$ = tree.NewRangeCond(false, $1, $3, $5, yylex.(*Lexer).buf)
    }
|   bit_expr NOT BETWEEN bit_expr AND predicate
    {
        $$ = tree.NewRangeCond(true, $1, $4, $6, yylex.(*Lexer).buf)
    }
|   bit_expr

like_escape_opt:
    {
        $$ = nil
    }
|   ESCAPE simple_expr
    {
        $$ = $2
    }

col_tuple:
    tuple_expression
    {
        $$ = $1
    }
|   subquery
    {
        $$ = $1
    }
// |   LIST_ARG

and_or_some:
    ALL
    {
        $$ = tree.ALL
    }
|    ANY
    {
        $$ = tree.ANY
    }
|    SOME
    {
        $$ = tree.SOME
    }

comparison_operator:
    '='
    {
        $$ = tree.EQUAL
    }
|   '<'
    {
        $$ = tree.LESS_THAN
    }
|   '>'
    {
        $$ = tree.GREAT_THAN
    }
|   LE
    {
        $$ = tree.LESS_THAN_EQUAL
    }
|   GE
    {
        $$ = tree.GREAT_THAN_EQUAL
    }
|   NE
    {
        $$ = tree.NOT_EQUAL
    }
|   NULL_SAFE_EQUAL
    {
        $$ = tree.NULL_SAFE_EQUAL
    }

keys:
    PRIMARY KEY
    {
        $$ = tree.NewAttributePrimaryKey(yylex.(*Lexer).buf)
    }
|   UNIQUE KEY
    {
        $$ = tree.NewAttributeUniqueKey(yylex.(*Lexer).buf)
    }
|   UNIQUE
    {
        $$ = tree.NewAttributeUnique(yylex.(*Lexer).buf)
    }
|   KEY
    {
        $$ = tree.NewAttributeKey(yylex.(*Lexer).buf)
    }

num_literal:
    INTEGRAL
    {
        str := fmt.Sprintf("%v", $1)
        switch v := $1.(type) {
        case uint64:
            $$ = tree.NewNumValWithType(constant.MakeUint64(v), str, false, tree.P_uint64, yylex.(*Lexer).buf)
        case int64:
            $$ = tree.NewNumValWithType(constant.MakeInt64(v), str, false, tree.P_int64, yylex.(*Lexer).buf)
        default:
            yylex.Error("parse integral fail")
            return 1
        }
    }
|   FLOAT
    {
        fval := $1.(float64)
        $$ = tree.NewNumValWithType(constant.MakeFloat64(fval), yylex.(*Lexer).scanner.LastToken, false, tree.P_float64, yylex.(*Lexer).buf)
    }
|   DECIMAL_VALUE
    {
        $$ = tree.NewNumValWithType(constant.MakeString($1), $1, false, tree.P_decimal, yylex.(*Lexer).buf)
    }

literal:
    STRING
    {
        $$ = tree.NewNumValWithType(constant.MakeString($1), $1, false, tree.P_char, yylex.(*Lexer).buf)
    }
|   INTEGRAL
    {
        str := fmt.Sprintf("%v", $1)
        switch v := $1.(type) {
        case uint64:
            $$ = tree.NewNumValWithType(constant.MakeUint64(v), str, false, tree.P_uint64, yylex.(*Lexer).buf)
        case int64:
            $$ = tree.NewNumValWithType(constant.MakeInt64(v), str, false, tree.P_int64, yylex.(*Lexer).buf)
        default:
            yylex.Error("parse integral fail")
            return 1
        }
    }
|   FLOAT
    {
        fval := $1.(float64)
        $$ = tree.NewNumValWithType(constant.MakeFloat64(fval), yylex.(*Lexer).scanner.LastToken, false, tree.P_float64, yylex.(*Lexer).buf)
    }
|   TRUE
    {
        $$ = tree.NewNumValWithType(constant.MakeBool(true), "true", false, tree.P_bool, yylex.(*Lexer).buf)
    }
|   FALSE
    {
        $$ = tree.NewNumValWithType(constant.MakeBool(false), "false", false, tree.P_bool, yylex.(*Lexer).buf)
    }
|   NULL
    {
        $$ = tree.NewNumValWithType(constant.MakeUnknown(), "null", false, tree.P_null, yylex.(*Lexer).buf)
    }
|   HEXNUM
    {
        $$ = tree.NewNumValWithType(constant.MakeString($1), $1, false, tree.P_hexnum, yylex.(*Lexer).buf)
    }
|   UNDERSCORE_BINARY HEXNUM
    {
        if strings.HasPrefix($2, "0x") {
            $2 = $2[2:]
        }
        $$ = tree.NewNumValWithType(constant.MakeString($2), $2, false, tree.P_bit, yylex.(*Lexer).buf)
    }
|   DECIMAL_VALUE
    {
        $$ = tree.NewNumValWithType(constant.MakeString($1), $1, false, tree.P_decimal, yylex.(*Lexer).buf)
    }
|   BIT_LITERAL
    {
        switch v := $1.(type) {
        case uint64:
            $$ = tree.NewNumValWithType(constant.MakeUint64(v), yylex.(*Lexer).scanner.LastToken, false, tree.P_uint64, yylex.(*Lexer).buf)
        case int64:
            $$ = tree.NewNumValWithType(constant.MakeInt64(v), yylex.(*Lexer).scanner.LastToken, false, tree.P_int64, yylex.(*Lexer).buf)
        case string:
            $$ = tree.NewNumValWithType(constant.MakeString(v), v, false, tree.P_bit, yylex.(*Lexer).buf)
        default:
            yylex.Error("parse integral fail")
            return 1
        }
    }
|   VALUE_ARG
    {
        $$ = tree.NewParamExpr(yylex.(*Lexer).GetParamIndex(), yylex.(*Lexer).buf)
    }
|   UNDERSCORE_BINARY STRING
    {
        $$ = tree.NewNumValWithType(constant.MakeString($2), $2, false, tree.P_ScoreBinary, yylex.(*Lexer).buf)
    }


column_type:
    numeric_type unsigned_opt zero_fill_opt
    {
        $$ = $1
        $$.InternalType.Unsigned = $2
        $$.InternalType.Zerofill = $3
    }
|   char_type
|   time_type
|   spatial_type

numeric_type:
    int_type length_opt
    {
        $$ = $1
        $$.InternalType.DisplayWith = $2
    }
|   decimal_type
    {
        $$ = $1
    }

int_type:
    BIT
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.BitFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_BIT), yylex.(*Lexer).buf)
    }
|   BOOL
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.BoolFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_BOOL), yylex.(*Lexer).buf)
    }
|   BOOLEAN
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.BoolFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_BOOL), yylex.(*Lexer).buf)
    }
|   INT1
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, $1, 8, "", uint32(defines.MYSQL_TYPE_TINY), yylex.(*Lexer).buf)
    }
|   TINYINT
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, $1, 8, "", uint32(defines.MYSQL_TYPE_TINY), yylex.(*Lexer).buf)
    }
|   INT2
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, $1, 16, "", uint32(defines.MYSQL_TYPE_SHORT), yylex.(*Lexer).buf)
    }
|   SMALLINT
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, $1, 16, "", uint32(defines.MYSQL_TYPE_SHORT), yylex.(*Lexer).buf)
    }
|   INT3
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, $1, 24, "", uint32(defines.MYSQL_TYPE_INT24), yylex.(*Lexer).buf)
    }
|   MEDIUMINT
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, $1, 24, "", uint32(defines.MYSQL_TYPE_INT24), yylex.(*Lexer).buf)
    }
|   INT4
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, $1, 32, "", uint32(defines.MYSQL_TYPE_LONG), yylex.(*Lexer).buf)
    }
|   INT
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, $1, 32, "", uint32(defines.MYSQL_TYPE_LONG), yylex.(*Lexer).buf)
    }
|   INTEGER
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, $1, 32, "", uint32(defines.MYSQL_TYPE_LONG), yylex.(*Lexer).buf)
    }
|   INT8
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, $1, 64, "", uint32(defines.MYSQL_TYPE_LONGLONG), yylex.(*Lexer).buf)
    }
|   BIGINT
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.IntFamily, $1, 64, "", uint32(defines.MYSQL_TYPE_LONGLONG), yylex.(*Lexer).buf)
    }

decimal_type:
    DOUBLE float_length_opt
    {
        if $2.DisplayWith > 255 {
            yylex.Error("Display width for double out of range (max = 255)")
            return 1
        }
        if $2.Scale > 30 {
            yylex.Error("Display scale for double out of range (max = 30)")
            return 1
        }
        if $2.Scale != tree.NotDefineDec && $2.Scale > $2.DisplayWith {
            yylex.Error("For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a'))")
                return 1
        }
        // Family FamilyString Width Locale Oid DisplayWith Scale
        $$ = tree.NewSQLTypeWithDisplayScale(tree.FloatFamily, $1, 64, "", uint32(defines.MYSQL_TYPE_DOUBLE), $2.DisplayWith, $2.Scale, yylex.(*Lexer).buf)
    }
|   FLOAT_TYPE float_length_opt
    {
        if $2.DisplayWith > 255 {
            yylex.Error("Display width for float out of range (max = 255)")
            return 1
        }
        if $2.Scale > 30 {
            yylex.Error("Display scale for float out of range (max = 30)")
            return 1
        }
        if $2.Scale != tree.NotDefineDec && $2.Scale > $2.DisplayWith {
        	yylex.Error("For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a'))")
        	return 1
        }
        if $2.DisplayWith >= 24 {
            // Family FamilyString Width Locale Oid DisplayWith Scale
            $$ = tree.NewSQLTypeWithDisplayScale(tree.FloatFamily, $1, 64, "", uint32(defines.MYSQL_TYPE_DOUBLE), $2.DisplayWith, $2.Scale, yylex.(*Lexer).buf)
        } else {
            // Family FamilyString Width Locale Oid DisplayWith Scale
            $$ = tree.NewSQLTypeWithDisplayScale(tree.FloatFamily, $1, 32, "", uint32(defines.MYSQL_TYPE_FLOAT), $2.DisplayWith, $2.Scale, yylex.(*Lexer).buf)
        }
    }

|   DECIMAL decimal_length_opt
    {
        if $2.Scale != tree.NotDefineDec && $2.Scale > $2.DisplayWith {
        yylex.Error("For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a'))")
        return 1
        }
        if $2.DisplayWith > 38 || $2.DisplayWith < 0 {
            yylex.Error("For decimal(M), M must between 0 and 38.")
                return 1
        } else if $2.DisplayWith <= 16 {
            // Family FamilyString Width Locale Oid DisplayWith Scale
            $$ = tree.NewSQLTypeWithDisplayScale(tree.FloatFamily, $1, 64, "", uint32(defines.MYSQL_TYPE_DECIMAL), $2.DisplayWith, $2.Scale, yylex.(*Lexer).buf)
        } else {
            // Family FamilyString Width Locale Oid DisplayWith Scale
            $$ = tree.NewSQLTypeWithDisplayScale(tree.FloatFamily, $1, 128, "", uint32(defines.MYSQL_TYPE_DECIMAL), $2.DisplayWith, $2.Scale, yylex.(*Lexer).buf)
        }
    }
|   NUMERIC decimal_length_opt
    {
        if $2.Scale != tree.NotDefineDec && $2.Scale > $2.DisplayWith {
        yylex.Error("For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a'))")
        return 1
        }
        if $2.DisplayWith > 38 || $2.DisplayWith < 0 {
            yylex.Error("For decimal(M), M must between 0 and 38.")
                return 1
        } else if $2.DisplayWith <= 16 {
            // Family FamilyString Width Locale Oid DisplayWith Scale
            $$ = tree.NewSQLTypeWithDisplayScale(tree.FloatFamily, $1, 64, "", uint32(defines.MYSQL_TYPE_DECIMAL), $2.DisplayWith, $2.Scale, yylex.(*Lexer).buf)
        } else {
            // Family FamilyString Width Locale Oid DisplayWith Scale
            $$ = tree.NewSQLTypeWithDisplayScale(tree.FloatFamily, $1, 128, "", uint32(defines.MYSQL_TYPE_DECIMAL), $2.DisplayWith, $2.Scale, yylex.(*Lexer).buf)
        }
    }
|   REAL float_length_opt
    {
        // Family FamilyString Width Locale Oid DisplayWith Scale
        $$ = tree.NewSQLTypeWithDisplayScale(tree.FloatFamily, $1, 64, "", uint32(defines.MYSQL_TYPE_DOUBLE), $2.DisplayWith, $2.Scale, yylex.(*Lexer).buf)
    }

time_type:
    DATE
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.DateFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_DATE), yylex.(*Lexer).buf)
    }
|   TIME timestamp_option_opt
    {
        if $2 < 0 || $2 > 6 {
            yylex.Error("For Time(fsp), fsp must in [0, 6]")
            return 1
        } else {
            // Family FamilyString Width Locale Oid DisplayWith Scale
            $$ = tree.NewSQLTypeWithDisplayScale(tree.TimeFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_TIME), $2, $2, yylex.(*Lexer).buf)
            $$.InternalType.TimePrecisionIsSet = true
        }
    }
|   TIMESTAMP timestamp_option_opt
    {
        if $2 < 0 || $2 > 6 {
            yylex.Error("For Timestamp(fsp), fsp must in [0, 6]")
            return 1
        } else {
            // Family FamilyString Width Locale Oid DisplayWith Scale
            $$ = tree.NewSQLTypeWithDisplayScale(tree.TimestampFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_TIMESTAMP), $2, $2, yylex.(*Lexer).buf)
            $$.InternalType.TimePrecisionIsSet = true
        }
    }
|   DATETIME timestamp_option_opt
    {
        if $2 < 0 || $2 > 6 {
            yylex.Error("For Datetime(fsp), fsp must in [0, 6]")
            return 1
        } else {
            // Family FamilyString Width Locale Oid DisplayWith Scale
            $$ = tree.NewSQLTypeWithDisplayScale(tree.TimestampFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_DATETIME), $2, $2, yylex.(*Lexer).buf)
            $$.InternalType.TimePrecisionIsSet = true
        }
    }
|   YEAR length_opt
    {
        // Family FamilyString Width Locale Oid DisplayWith Scale
        $$ = tree.NewSQLTypeWithDisplayScale(tree.IntFamily, $1, 16, "", uint32(defines.MYSQL_TYPE_YEAR), $2, 0, yylex.(*Lexer).buf)
        $$.InternalType.TimePrecisionIsSet = true
    }

char_type:
    CHAR length_option_opt
    {
        // Family FamilyString Width Locale Oid DisplayWith Scale
        $$ = tree.NewSQLTypeWithDisplayScale(tree.StringFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_STRING), $2, 0, yylex.(*Lexer).buf)
    }
|   VARCHAR length_option_opt
    {
        // Family FamilyString Width Locale Oid DisplayWith Scale
        $$ = tree.NewSQLTypeWithDisplayScale(tree.StringFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_VARCHAR), $2, 0, yylex.(*Lexer).buf)
    }
|   BINARY length_option_opt
    {
        // Family FamilyString Width Locale Oid DisplayWith Scale
        $$ = tree.NewSQLTypeWithDisplayScale(tree.StringFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_VARCHAR), $2, 0, yylex.(*Lexer).buf)
    }
|   VARBINARY length_option_opt
    {
        // Family FamilyString Width Locale Oid DisplayWith Scale
        $$ = tree.NewSQLTypeWithDisplayScale(tree.StringFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_VARCHAR), $2, 0, yylex.(*Lexer).buf)
    }
|   TEXT
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.BlobFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_TEXT), yylex.(*Lexer).buf)
    }
|   TINYTEXT
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.BlobFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_TEXT), yylex.(*Lexer).buf)
    }
|   MEDIUMTEXT
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.BlobFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_TEXT), yylex.(*Lexer).buf)
    }
|   LONGTEXT
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.BlobFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_TEXT), yylex.(*Lexer).buf)
    }
|   BLOB
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.BlobFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_BLOB), yylex.(*Lexer).buf)
    }
|   TINYBLOB
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.BlobFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_TINY_BLOB), yylex.(*Lexer).buf)
    }
|   MEDIUMBLOB
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.BlobFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_MEDIUM_BLOB), yylex.(*Lexer).buf)
    }
|   LONGBLOB
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.BlobFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_LONG_BLOB), yylex.(*Lexer).buf)
    }
|   JSON
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.JsonFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_JSON), yylex.(*Lexer).buf)
    }
|   VECF32 length_option_opt
    {
        // Family FamilyString Width Locale Oid DisplayWith Scale
        $$ = tree.NewSQLTypeWithDisplayScale(tree.ArrayFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_VARCHAR), $2, 0, yylex.(*Lexer).buf)
    }
|   VECF64 length_option_opt
    {
        // Family FamilyString Width Locale Oid DisplayWith Scale
        $$ = tree.NewSQLTypeWithDisplayScale(tree.ArrayFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_VARCHAR), $2, 0, yylex.(*Lexer).buf)
    }
|   ENUM '(' enum_values ')'
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.EnumFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_ENUM), yylex.(*Lexer).buf)
        $$.InternalType.EnumValues = $3
    }
|   SET '(' enum_values ')'
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.SetFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_SET), yylex.(*Lexer).buf)
        $$.InternalType.EnumValues = $3
    }
|  UUID
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.UuidFamily, $1, 128, "", uint32(defines.MYSQL_TYPE_UUID), yylex.(*Lexer).buf)
    }

do_stmt:
    DO expression_list
    {
        $$ = tree.NewDo(
            $2,
            yylex.(*Lexer).buf,
        )
    }

declare_stmt:
    DECLARE var_name_list column_type
    {
        // Variables ColumnType DefaultVal
        $$ = tree.NewDeclare(
            $2,
            $3,
            tree.NewNumValWithType(constant.MakeUnknown(), "null", false, tree.P_null, yylex.(*Lexer).buf),
            yylex.(*Lexer).buf,
        )
    }
    |
    DECLARE var_name_list column_type DEFAULT expression
    {
        // Variables ColumnType DefaultVal
        $$ = tree.NewDeclare(
            $2,
            $3,
            $5,
            yylex.(*Lexer).buf,
        )
    }

spatial_type:
    GEOMETRY
    {
        // Family FamilyString Width Locale Oid
        $$ = tree.NewSQLType(tree.GeometryFamily, $1, 0, "", uint32(defines.MYSQL_TYPE_GEOMETRY), yylex.(*Lexer).buf)
    }
// |   POINT
// |   LINESTRING
// |   POLYGON
// |   GEOMETRYCOLLECTION
// |   MULTIPOINT
// |   MULTILINESTRING
// |   MULTIPOLYGON

// TODO:
// need to encode SQL string
enum_values:
    STRING
    {
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $$, yylex.(*Lexer).buf.CopyString($1))
    }
|   enum_values ',' STRING
    {
        $$ = buffer.MakeSlice[string](yylex.(*Lexer).buf)
        $$ = buffer.AppendSlice[string](yylex.(*Lexer).buf, $1, yylex.(*Lexer).buf.CopyString($3))
    }

length_opt:
    /* EMPTY */
    {
        $$ = 0
    }
|    length

timestamp_option_opt:
    /* EMPTY */
        {
            $$ = 0
        }
|    '(' INTEGRAL ')'
    {
        $$ = int32($2.(int64))
    }

length_option_opt:
    {
        $$ = int32(-1)
    }
|    '(' INTEGRAL ')'
    {
        $$ = int32($2.(int64))
    }

length:
   '(' INTEGRAL ')'
    {
        $$ = tree.GetDisplayWith(int32($2.(int64)))
    }

float_length_opt:
    /* EMPTY */
    {
        // DisplayWith Scale
        $$ = *tree.NewLengthScaleOpt(
            tree.NotDefineDisplayWidth,
            tree.NotDefineDec,
            yylex.(*Lexer).buf,
        )
    }
|   '(' INTEGRAL ')'
    {
        // DisplayWith Scale
        $$ = *tree.NewLengthScaleOpt(
            tree.GetDisplayWith(int32($2.(int64))),
            tree.NotDefineDec,
            yylex.(*Lexer).buf,
        )
    }
|   '(' INTEGRAL ',' INTEGRAL ')'
    {
        // DisplayWith Scale
        $$ = *tree.NewLengthScaleOpt(
            tree.GetDisplayWith(int32($2.(int64))),
            int32($4.(int64)),
            yylex.(*Lexer).buf,
        )
    }

decimal_length_opt:
    /* EMPTY */
    {
        // DisplayWith Scale
        $$ = *tree.NewLengthScaleOpt(
            38,           // this is the default precision for decimal
            0,
            yylex.(*Lexer).buf,
        )
    }
|   '(' INTEGRAL ')'
    {
        // DisplayWith Scale
        $$ = *tree.NewLengthScaleOpt(
            tree.GetDisplayWith(int32($2.(int64))),
            0,
            yylex.(*Lexer).buf,
        )
    }
|   '(' INTEGRAL ',' INTEGRAL ')'
    {
        // DisplayWith Scale
        $$ = *tree.NewLengthScaleOpt(
            tree.GetDisplayWith(int32($2.(int64))),
            int32($4.(int64)),
            yylex.(*Lexer).buf,
        )
    }

unsigned_opt:
    /* EMPTY */
    {
        $$ = false
    }
|   UNSIGNED
    {
        $$ = true
    }
|   SIGNED
    {
        $$ = false
    }

zero_fill_opt:
    /* EMPTY */
    {

    }
|   ZEROFILL
    {
        $$ = true
    }

charset_keyword:
    CHARSET
|   CHARACTER SET
|   CHAR SET

equal_opt:
    {
        $$ = ""
    }
|   '='
    {
        $$ = string($1)
    }

//sql_id:
//    id_or_var
//|   non_reserved_keyword

//reserved_sql_id:
//    sql_id
//|   reserved_keyword

//reserved_table_id:
//    table_id
//|   reserved_keyword

//reserved_keyword:
//    ADD
//|   ALL
//|   AND
//|   AS
//|   ASC
//|   ASCII
//|   AUTO_INCREMENT
//|   BETWEEN
//|   BINARY
//|   BY
//|   CASE
//|   CHAR
//|   COLLATE
//|   COLLATION
//|   CONVERT
//|   COALESCE
//|   CREATE
//|   CROSS
//|   CURRENT_DATE
//|   CURRENT_ROLE
//|   CURRENT_USER
//|   CURRENT_TIME
//|   CURRENT_TIMESTAMP
//|   CIPHER
//|   SAN
//|   SSL
//|   SUBJECT
//|   DATABASE
//|   DATABASES
//|   DEFAULT
//|   DELETE
//|   DESC
//|   DESCRIBE
//|   DISTINCT
//|   DISTINCTROW
//|   DIV
//|   DROP
//|   ELSE
//|   END
//|   ESCAPE
//|   EXISTS
//|   EXPLAIN
//|   FALSE
//|   FIRST
//|   AFTER
//|   FOR
//|   FORCE
//|   FROM
//|   GROUP
//|   HAVING
//|   HOUR
//|   IDENTIFIED
//|   IF
//|   IGNORE
//|   IN
//|   INFILE
//|   INDEX
//|   INNER
//|   INSERT
//|   INTERVAL
//|   INTO
//|   IS
//|   ISSUER
//|   JOIN
//|   KEY
//|   LAST
//|   LEFT
//|   LIKE
//|	ILIKE
//|   LIMIT
//|   LOCALTIME
//|   LOCALTIMESTAMP
//|   LOCK
//|   LOAD
//|   IMPORT
//|   MATCH
//|   MAXVALUE
//|   MOD
//|   MICROSECOND
//|   MINUTE
//|   NATURAL
//|   NOT
//|   NONE
//|   NULL
//|   NULLS
//|   ON
//|   OR
//|   ORDER
//|   OUTER
//|   REGEXP
//|   RENAME
//|   REPLACE
//|   RIGHT
//|   REQUIRE
//|   REPEAT
//|   ROW_COUNT
//|    REFERENCES
//|   RECURSIVE
//|   REVERSE
//|   SCHEMA
//|   SCHEMAS
//|   SELECT
//|   SECOND
//|   SEPARATOR
//|   SET
//|   SHOW
//|   STRAIGHT_JOIN
//|   TABLE
//|   THEN
//|   TO
//|   TRUE
//|   TRUNCATE
//|   UNION
//|   UNIQUE
//|   UPDATE
//|   USE
//|   USING
//|   UTC_DATE
//|   UTC_TIME
//|   UTC_TIMESTAMP
//|   VALUES
//|   WHEN
//|   WHERE
//|   WEEK
//|   WITH
//|   TERMINATED
//|   OPTIONALLY
//|   ENCLOSED
//|   ESCAPED
//|   STARTING
//|   LINES
//|   ROWS
//|   INT1
//|   INT2
//|   INT3
//|   INT4
//|   INT8
//|   CHECK
//|    CONSTRAINT
//|   PRIMARY
//|   FULLTEXT
//|   FOREIGN
//|    ROW
//|   OUTFILE
//|    SQL_SMALL_RESULT
//|    SQL_BIG_RESULT
//|    LEADING
//|    TRAILING
//|   CHARACTER
//|    LOW_PRIORITY
//|    HIGH_PRIORITY
//|    DELAYED
//|   PARTITION
//|    QUICK
//|   EXCEPT
//|   ADMIN_NAME
//|   RANDOM
//|   SUSPEND
//|   REUSE
//|   CURRENT
//|   OPTIONAL
//|   FAILED_LOGIN_ATTEMPTS
//|   PASSWORD_LOCK_TIME
//|   UNBOUNDED
//|   SECONDARY
//|   DECLARE
//|   MODUMP
//|   OVER
//|   PRECEDING
//|   FOLLOWING
//|   GROUPS
//|   LOCKS
//|   TABLE_NUMBER
//|   COLUMN_NUMBER
//|   TABLE_VALUES
//|   RETURNS
//|   MYSQL_COMPATIBILITY_MODE

non_reserved_keyword:
    ACCOUNT
|   ACCOUNTS
|   AGAINST
|   AVG_ROW_LENGTH
|   AUTO_RANDOM
|   ATTRIBUTE
|   ACTION
|   ALGORITHM
|   BEGIN
|   BIGINT
|   BIT
|   BLOB
|   BOOL
|   CANCEL
|   CHAIN
|   CHECKSUM
|   CLUSTER
|   COMPRESSION
|   COMMENT_KEYWORD
|   COMMIT
|   COMMITTED
|   CHARSET
|   COLUMNS
|   CONNECTION
|   CONSISTENT
|   COMPRESSED
|   COMPACT
|   COLUMN_FORMAT
|   CONNECTOR
|   CONNECTORS
|   SECONDARY_ENGINE_ATTRIBUTE
|   ENGINE_ATTRIBUTE
|   INSERT_METHOD
|   CASCADE
|   DAEMON
|   DATA
|	DAY
|   DATETIME
|   DECIMAL
|   DYNAMIC
|   DISK
|   DO
|   DOUBLE
|   DIRECTORY
|   DUPLICATE
|   DELAY_KEY_WRITE
|   ENUM
|   ENCRYPTION
|   ENGINE
|   EXPANSION
|   EXTENDED
|   EXPIRE
|   ERRORS
|   ENFORCED
|   FORMAT
|   FLOAT_TYPE
|   FULL
|   FIXED
|   FIELDS
|   GEOMETRY
|   GEOMETRYCOLLECTION
|   GLOBAL
|   PERSIST
|   GRANT
|   INT
|   INTEGER
|   INDEXES
|   ISOLATION
|   JSON
|   VECF32
|   VECF64
|   KEY_BLOCK_SIZE
|   KEYS
|   LANGUAGE
|   LESS
|   LEVEL
|   LINESTRING
|   LONGBLOB
|   LONGTEXT
|   LOCAL
|   LINEAR
|   LIST
|   MEDIUMBLOB
|   MEDIUMINT
|   MEDIUMTEXT
|   MEMORY
|   MODE
|   MULTILINESTRING
|   MULTIPOINT
|   MULTIPOLYGON
|   MAX_QUERIES_PER_HOUR
|   MAX_UPDATES_PER_HOUR
|   MAX_CONNECTIONS_PER_HOUR
|   MAX_USER_CONNECTIONS
|   MAX_ROWS
|   MIN_ROWS
|   MONTH
|   NAMES
|   NCHAR
|   NUMERIC
|   NEVER
|   NO
|   OFFSET
|   ONLY
|   OPTIMIZE
|   OPEN
|   OPTION
|   PACK_KEYS
|   PARTIAL
|   PARTITIONS
|   POINT
|   POLYGON
|   PROCEDURE
|   PROXY
|   QUERY
|   PAUSE
|   PROFILES
|   ROLE
|   RANGE
|   READ
|   REAL
|   REORGANIZE
|   REDUNDANT
|   REPAIR
|   REPEATABLE
|   RELEASE
|   RESUME
|   REVOKE
|   REPLICATION
|   ROW_FORMAT
|   ROLLBACK
|   RESTRICT
|   SESSION
|   SERIALIZABLE
|   SHARE
|   SIGNED
|   SMALLINT
|   SNAPSHOT
|   SPATIAL
|   START
|   STATUS
|   STORAGE
|   STATS_AUTO_RECALC
|   STATS_PERSISTENT
|   STATS_SAMPLE_PAGES
|	SOURCE
|   SUBPARTITIONS
|   SUBPARTITION
|   SIMPLE
|   TASK
|   TEXT
|   THAN
|   TINYBLOB
|   TIME %prec LOWER_THAN_STRING
|   TINYINT
|   TINYTEXT
|   TRANSACTION
|   TRIGGER
|   UNCOMMITTED
|   UNSIGNED
|   UNUSED
|   UNLOCK
|   USER
|   VARBINARY
|   VARCHAR
|   VARIABLES
|   VIEW
|   WRITE
|   WARNINGS
|   WORK
|   X509
|   ZEROFILL
|   YEAR
|   TYPE
|   HEADER
|   MAX_FILE_SIZE
|   FORCE_QUOTE
|   QUARTER
|   UNKNOWN
|   ANY
|   SOME
|   TIMESTAMP %prec LOWER_THAN_STRING
|   DATE %prec LOWER_THAN_STRING
|   TABLES
|   SEQUENCES
|   URL
|   PASSWORD %prec LOWER_THAN_EQ
|   HASH
|   ENGINES
|   TRIGGERS
|   HISTORY
|   LOW_CARDINALITY
|   S3OPTION
|   EXTENSION
|   NODE
|   ROLES
|   UUID
|   PARALLEL
|   INCREMENT
|   CYCLE
|   MINVALUE
|	PROCESSLIST
|   PUBLICATION
|   SUBSCRIPTIONS
|   PUBLICATIONS
|   PROPERTIES
|	WEEK
|   DEFINER
|   SQL
|   STAGE
|   STAGES
|   BACKUP
| FILESYSTEM

func_not_keyword:
    DATE_ADD
|    DATE_SUB
|   NOW
|    ADDDATE
|   CURDATE
|   POSITION
|   SESSION_USER
|   SUBDATE
|   SYSTEM_USER
|   TRANSLATE

not_keyword:
    ADDDATE
|   BIT_AND
|   BIT_OR
|   BIT_XOR
|   CAST
|   COUNT
|   APPROX_COUNT
|   APPROX_COUNT_DISTINCT
|   APPROX_PERCENTILE
|   CURDATE
|   CURTIME
|   DATE_ADD
|   DATE_SUB
|   EXTRACT
|   GROUP_CONCAT
|   MAX
|   MID
|   MIN
|   NOW
|   POSITION
|   SESSION_USER
|   STD
|   STDDEV
|   STDDEV_POP
|   STDDEV_SAMP
|   SUBDATE
|   SUBSTR
|   SUBSTRING
|   SUM
|   SYSDATE
|   SYSTEM_USER
|   TRANSLATE
|   TRIM
|   VARIANCE
|   VAR_POP
|   VAR_SAMP
|   AVG
|	TIMESTAMPDIFF
|   NEXTVAL
|   SETVAL
|   CURRVAL
|   LASTVAL
|   HEADERS
|   BIT_CAST

//mo_keywords:
//    PROPERTIES
//  BSI
//  ZONEMAP

%%
