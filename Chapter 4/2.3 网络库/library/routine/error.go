package routine

import "errors"

var (
	// ErrInvalidManager 不支持的协程控制器类型
	ErrInvalidManager = errors.New("unsupported routine manager type")
	// ErrInvalidPanicError 不支持的崩溃错误类型
	ErrInvalidPanicError = errors.New("unsupported panic error interface")
	// ErrManagerClosed 控制器关闭时，无法正常执行任务，请先调用 routine.Open()
	ErrManagerClosed = errors.New("cannot exec task when manager closed")
)
