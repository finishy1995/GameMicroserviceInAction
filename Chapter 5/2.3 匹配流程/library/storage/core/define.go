package core

// Model 数据标准模型，使用 Version 做版本管理 + 乐观锁
type Model struct {
	Version uint64
}
