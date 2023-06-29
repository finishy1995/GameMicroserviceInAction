package log

import "errors"

var (
	// ErrInvalidConfig 不支持的日志设置
	ErrInvalidConfig = errors.New("invalid log config")
)
