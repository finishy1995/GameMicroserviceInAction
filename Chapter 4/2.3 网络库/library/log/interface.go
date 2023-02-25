package log

import "context"

// Level 日志等级
type Level uint8

// NewLoggerFunc 创建日志实例的函数
type NewLoggerFunc func(service string, timeFormat string) Logger

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
)

const (
	DefaultLevel      = INFO
	DefaultService    = "Unknown"
	DefaultTimeFormat = "15:04:05.000"
)

// Logger 日志处理抽象
type Logger interface {
	// Debug 打印 Debug 信息
	Debug(string, ...interface{})
	// Info 打印 Info 信息
	Info(string, ...interface{})
	// Warning 打印 Warning 信息
	Warning(string, ...interface{})
	// Error 打印 Error 信息
	Error(string, ...interface{})
	// SetLevel 设置日志等级
	SetLevel(Level)
	// WithContext 设置日志上下文（通常是用作设置 trace span 信息，可以用作扩展功能的信息传递）
	WithContext(ctx context.Context)
}

// Config 日志配置
//
//	我们推荐用配置文件的方式设置日志
type Config struct {
	Type        string `json:",default=default,options=default|zero"`
	Level       string `json:",default=info,options=info|debug|warn|error"`
	ServiceName string `json:",optional"`
	TimeFormat  string `json:",optional"`
}

var (
	logInstance Logger
	logTypeMap  = map[string]NewLoggerFunc{
		"default": newGoLoggingLogger,
		"":        newGoLoggingLogger,
	}
	logLevelStringMap = map[string]Level{
		"info":  INFO,
		"debug": DEBUG,
		"warn":  WARNING,
		"error": ERROR,
	}
)

func init() {
	SetLogger(logTypeMap[""](DefaultService, DefaultTimeFormat))
	SetLevel(DefaultLevel)
}

// MustSetup 通过配置初始化日志，失败会使当前线程崩溃
//
//	在大多数场景下，推荐使用这个函数初始化/改动日志设置
func MustSetup(c *Config) {
	err := Setup(c)
	if err != nil {
		panic(err)
	}
}

// RegisterLoggerType 注册一个新的日志类型
//
//	通过这个函数，可以轻松扩展这个库
func RegisterLoggerType(name string, f NewLoggerFunc) {
	if _, ok := logTypeMap[name]; ok {
		Warning("replace exist logger type %s", name)
	}
	logTypeMap[name] = f
}

// Setup 通过配置初始化日志，失败会返回 error
func Setup(c *Config) error {
	if c == nil {
		return ErrInvalidConfig
	}

	newInstanceLogic, ok := logTypeMap[c.Type]
	if !ok {
		return ErrInvalidConfig
	}
	level, ok := logLevelStringMap[c.Level]
	if !ok {
		return ErrInvalidConfig
	}
	serviceName := c.ServiceName
	if serviceName == "" {
		serviceName = DefaultService
	}
	timeFormat := c.TimeFormat
	if timeFormat == "" {
		timeFormat = DefaultTimeFormat
	}
	SetLogger(newInstanceLogic(serviceName, timeFormat))
	SetLevel(level)

	return nil
}

// SetLevel 设置日志等级
func SetLevel(level Level) {
	logInstance.SetLevel(level)
}

// SetLogger 设置日志实例
//
//	我们更推荐使用 MustSetup
func SetLogger(logger Logger) {
	logInstance = logger
}

// Debug 打印 Debug 信息
func Debug(message string, param ...interface{}) {
	logInstance.Debug(message, param...)
}

// Info 打印 Info 信息
func Info(message string, param ...interface{}) {
	logInstance.Info(message, param...)
}

// Warning 打印 Warning 信息
func Warning(message string, param ...interface{}) {
	logInstance.Warning(message, param...)
}

// Error 打印 Error 信息
func Error(message string, param ...interface{}) {
	logInstance.Error(message, param...)
}

// WithContext 打印 WithContext 信息
func WithContext(ctx context.Context) {
	logInstance.WithContext(ctx)
}
