package log

import (
	"context"
	"fmt"
	"github.com/op/go-logging"
	"os"
)

// goLoggingLogger 使用 go-logging 第三方库实现日志
type goLoggingLogger struct {
	log *logging.Logger
}

const (
	ModuleName = "Default"
)

var (
	levelMap = map[Level]logging.Level{
		DEBUG:   logging.DEBUG,
		INFO:    logging.INFO,
		WARNING: logging.WARNING,
		ERROR:   logging.ERROR,
	}
)

func newGoLoggingLogger(service string, timeFormat string) Logger {
	format := logging.MustStringFormatter(
		fmt.Sprintf("%%{color}%%{time:%s} [%%{level:.4s}] [%s]▶ %%{color:reset} %%{message}", timeFormat, service),
	)
	backend := logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format)
	log := &goLoggingLogger{
		log: logging.MustGetLogger(ModuleName),
	}
	logging.SetBackend(backend)

	return log
}

func (g *goLoggingLogger) Debug(s string, i ...interface{}) {
	g.log.Debugf(s, i...)
}

func (g *goLoggingLogger) Info(s string, i ...interface{}) {
	g.log.Infof(s, i...)
}

func (g *goLoggingLogger) Warning(s string, i ...interface{}) {
	g.log.Warningf(s, i...)
}

func (g *goLoggingLogger) Error(s string, i ...interface{}) {
	g.log.Errorf(s, i...)
}

func (g *goLoggingLogger) SetLevel(level Level) {
	if l, ok := levelMap[level]; ok {
		logging.SetLevel(l, ModuleName)
	} else {
		logging.SetLevel(levelMap[DefaultLevel], ModuleName)
	}
}

func (g *goLoggingLogger) WithContext(_ context.Context) {
	// TODO: 读取 context 中的 trace span 信息，并更新日志输出内容
}
