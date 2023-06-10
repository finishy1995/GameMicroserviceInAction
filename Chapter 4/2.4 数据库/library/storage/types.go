package storage

import (
	"ProjectX/library/storage/core"
	"time"
)

// Type 存储类型
type Type uint8

const (
	// InMemory 内存存储，在单个服务内线程安全，不保证集群内的数据安全，非生产环境的选择，以及少数测试场景
	InMemory Type = iota
	// DynamoDB 借助 NoSQL 数据库的数据存取，保证集群内的数据安全，保证1秒内的数据最终一致性（只有货币相关操作是强一致性）
	// 每个接口调用小于 10ms
	DynamoDB
	// 这里可以添加其他数据库类型，例如 MongoDB、MySQL 等，MySQL 可以借助 ORM 库实现
)

// Model 存储基本模型
type Model struct {
	core.Model
}

type Config struct {
	StorageType string        `json:",default=memory,options=memory|dynamo"`
	Host        string        `json:",optional"`
	Endpoint    string        `json:",optional"`
	Prefix      string        `json:",optional"`
	MaxLength   int           `json:",default=0"`
	Tick        time.Duration `json:",default=1s"`
	AK          string        `json:",optional"`
	SK          string        `json:",optional"`
	Prob        float64       `json:",default=0,optional"`
}

const (
	InMemoryStr = "memory"
	DynamoDBStr = "dynamo"
)
