package dynamodb

import (
	"ProjectX/library/storage/core"
	"ProjectX/library/storage/src/tools"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"time"
)

type Storage struct {
	db     *dynamo.DB
	prefix string
}

const (
	DefaultAWSRegion = "us-east-1"
	DefaultTimeout   = 200 * time.Millisecond
)

func NewStorage(region string, mock string, prefix string, ak string, sk string) *Storage {
	st := new(Storage)
	if region == "" {
		region = DefaultAWSRegion
	}

	config := aws.NewConfig().WithRegion(region)
	if config == nil {
		return nil
	}
	if mock != "" {
		config = config.WithEndpoint(mock).WithCredentials(credentials.NewStaticCredentials("1", "1", "1"))
	} else {
		if ak != "" && sk != "" {
			config = config.WithCredentials(credentials.NewStaticCredentials(ak, sk, ""))
		}
	}
	mySession := session.Must(session.NewSession(config))

	st.db = dynamo.New(mySession, config)
	if prefix != "" {
		st.prefix = prefix + "-"
	}
	return st
}

func getContext() (aws.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}

// CreateTable 创建一个新的存储对象表
func (st *Storage) CreateTable(value interface{}) error {
	tableName := tools.GetStructOnlyName(value)
	if tableName == "" {
		return core.ErrUnsupportedValueType
	}
	hashKey, _ := tools.GetHashAndRangeKey(value, true)
	if hashKey == "" {
		return core.ErrUnsupportedValueType
	}
	tableName = st.prefix + tableName

	process := st.db.CreateTable(tableName, value)
	if process == nil {
		return core.ErrUnsupportedValueType
	}
	return process.Run()
}

// Create 创建一个新的存储对象（单主键时主键不相同，主键+排序键时有一个不相同）
// value 为符合 tag 定义的 struct
func (st *Storage) Create(value interface{}) error {
	tableName := tools.GetStructOnlyName(value)
	if tableName == "" {
		return core.ErrUnsupportedValueType
	}
	hashKey, _ := tools.GetHashAndRangeKey(value, true)
	if hashKey == "" {
		return core.ErrUnsupportedValueType
	}
	tableName = st.prefix + tableName

	table := st.db.Table(tableName)
	process := table.Put(value)
	if process == nil {
		return core.ErrUnsupportedValueType
	}

	ctx, cancel := getContext()
	defer cancel()
	err := process.If(fmt.Sprintf("attribute_not_exists(%s)", hashKey)).RunWithContext(ctx)
	_, ok := err.(*dynamodb.ConditionalCheckFailedException)
	if ok {
		return core.ErrDuplicateKey
	} else {
		return err
	}
}

// Delete 删除一个存储对象（单主键时不需要额外参数，主键+排序键时需要把排序键的值作为额外参数）
// value 为符合 tag 定义的 struct
func (st *Storage) Delete(value interface{}, hash interface{}, args ...interface{}) error {
	tableName := tools.GetStructOnlyName(value)
	if tableName == "" {
		return core.ErrUnsupportedValueType
	}
	hashKey, rangeKey := tools.GetHashAndRangeKey(value, true)
	if hashKey == "" {
		return core.ErrUnsupportedValueType
	}
	if rangeKey != "" && len(args) == 0 {
		return core.ErrMissingRangeValue
	}
	tableName = st.prefix + tableName

	table := st.db.Table(tableName)
	del := table.Delete(hashKey, hash)
	if del == nil {
		return core.ErrUnsupportedValueType
	}
	if rangeKey != "" {
		del = del.Range(rangeKey, args[0])
	}

	ctx, cancel := getContext()
	defer cancel()
	return del.RunWithContext(ctx)
}

// Save 保存一个存储对象（请勿用这个方法创建对象，可能会造成同步性问题）
// value 为符合 tag 定义的 struct ptr（注：一定要是 struct ptr）
func (st *Storage) Save(value interface{}) error {
	tableName := tools.GetStructName(value)
	if tableName == "" {
		return core.ErrUnsupportedValueType
	}
	hashKey, _ := tools.GetHashAndRangeKey(value, true)
	if hashKey == "" {
		return core.ErrUnsupportedValueType
	}
	version, err := tools.TrySetStructVersion(value)
	if err != nil {
		return err
	}
	tableName = st.prefix + tableName

	table := st.db.Table(tableName)
	process := table.Put(value).If(fmt.Sprintf("'%s' = ?", tools.VersionMark), version)
	if process == nil {
		return core.ErrUnsupportedValueType
	}

	ctx, cancel := getContext()
	defer cancel()
	return process.RunWithContext(ctx)
}

// First 获取符合要求的存储对象（单主键时不需要额外参数，主键+排序键时需要把排序键的值作为额外参数）
// value 为符合 tag 定义的 struct ptr
func (st *Storage) First(value interface{}, hash interface{}, args ...interface{}) error {
	tableName := tools.GetStructName(value)
	if tableName == "" {
		return core.ErrUnsupportedValueType
	}
	hashKey, rangeKey := tools.GetHashAndRangeKey(value, true)
	if hashKey == "" {
		return core.ErrUnsupportedValueType
	}
	if rangeKey != "" && len(args) == 0 {
		return core.ErrMissingRangeValue
	}
	tableName = st.prefix + tableName

	table := st.db.Table(tableName)
	query := table.Get(hashKey, hash)
	if query == nil {
		return core.ErrUnsupportedValueType
	}
	if rangeKey != "" {
		defaultOp := dynamo.GreaterOrEqual
		if len(args) >= 2 {
			defaultOp = dynamo.Operator(args[1].(string))
		}

		query = query.Range(rangeKey, defaultOp, args[0])
	}

	ctx, cancel := getContext()
	defer cancel()
	err := query.OneWithContext(ctx, value)
	if err == dynamo.ErrNotFound {
		return core.ErrNotFound
	}
	return err
}

// Find 获取所有符合要求的对象，性能远低于 First
// value 为符合 tag 定义的 struct slice ptr （注：&[]struct）
// limit 为限制数量， <= 0 即不限制数量
// expr 为表达式（空代表不使用表达式），参考 dynamodb 文档、或 https://github.com/guregu/dynamo
// 其他为补充表达式的具体值
func (st *Storage) Find(value interface{}, limit int64, expr string, args ...interface{}) error {
	tableName := tools.GetSliceStructName(value)
	if tableName == "" {
		return core.ErrUnsupportedValueType
	}
	tableName = st.prefix + tableName

	table := st.db.Table(tableName)
	process := table.Scan()
	if process == nil {
		return core.ErrUnsupportedValueType
	}
	if limit > 0 {
		process.Limit(limit)
	}
	if expr != "" {
		process.Filter(expr, args...)
	}

	ctx, cancel := getContext()
	defer cancel()
	return process.AllWithContext(ctx, value)
}
