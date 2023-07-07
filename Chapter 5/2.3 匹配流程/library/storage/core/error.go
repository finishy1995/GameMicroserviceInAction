package core

import "errors"

var (
	// ErrUnsupportedValueType 不支持的参数类型
	ErrUnsupportedValueType = errors.New("unsupported value type")

	// ErrMissingRangeValue 排序值丢失
	ErrMissingRangeValue = errors.New("missing range value")

	// ErrDuplicateKey 主键重复
	ErrDuplicateKey = errors.New("duplicate key")

	// ErrUnsupportedExprType 不支持的筛选字符串类型
	ErrUnsupportedExprType = errors.New("unsupported expr type")

	// ErrNotFound 未查询到指定对象
	ErrNotFound = errors.New("no item found")

	// ErrExpiredValue 当前对象非最新
	ErrExpiredValue = errors.New("item has updated")
)
