package core

import "errors"

var (
	// ErrInvalidAddress 不合法的地址输入
	ErrInvalidAddress 		= errors.New("invalid address input")
	// ErrInvalidGetAgentFunc 不合法的新建 agent 函数
	ErrInvalidGetAgentFunc 	= errors.New("invalid new agent function")
	// ErrInvalidCodec 不合法的 Codec
	ErrInvalidCodec 		= errors.New("codec invalid")
	// ErrInvalidAgent 不合法的 Agent
	ErrInvalidAgent 		= errors.New("agent invalid")
	// ErrTooLessLength 写入缓冲长度太低
	ErrTooLessLength		= errors.New("buffer too less length")
	// ErrTooMoreLength 写入缓冲长度太高
	ErrTooMoreLength		= errors.New("buffer too more length")
	// ErrPacketSplit 网络包传输不完全
	ErrPacketSplit 			= errors.New("network packet split")
)
