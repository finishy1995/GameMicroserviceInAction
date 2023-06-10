// Package log 日志打印输出
//
// Example Usage
//
//	log.MustSetup(c) // 使用配置初始化日志
//	log.Info("hello %s", nickname) // Info 日志
//	log.Error("error %s", err.String()) // Error 日志
//	log.WithContext(ctx) // 设置 trace span 等附加信息
package log
