package pool

type Status uint8

const (
	// StatusUnknown 未知状态
	StatusUnknown Status = iota
	// StatusPreparing 准备中
	StatusPreparing
	// StatusMatching 匹配中
	StatusMatching
	// StatusMatched 已匹配
	StatusMatched
	// StatusCanceled 已取消
	StatusCanceled
	// StatusTimeout 已超时
	StatusTimeout
)

type Action uint8

const (
	// ActionUnknown 未知动作
	ActionUnknown Action = iota
	// ActionAdd 增加一个匹配
	ActionAdd
	// ActionCancel 取消一个匹配
	ActionCancel
)
