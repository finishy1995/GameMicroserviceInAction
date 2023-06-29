package network

import "errors"

var (
	// ErrUnsupportedNetType 不支持的网络类型
	ErrUnsupportedNetType = errors.New("unsupported net type")
)
