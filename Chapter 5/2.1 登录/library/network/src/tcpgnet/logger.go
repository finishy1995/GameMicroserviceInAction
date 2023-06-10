package tcpgnet

import (
	"ProjectX/library/log"
)

type logger struct {
}

// Debugf logs messages at DEBUG level.
func (l *logger) Debugf(format string, args ...interface{}) {
	log.Debug(format, args...)
}

// Infof logs messages at INFO level.
func (l *logger) Infof(format string, args ...interface{}) {
	log.Info(format, args...)
}

// Warnf logs messages at WARN level.
func (l *logger) Warnf(format string, args ...interface{}) {
	log.Warning(format, args...)
}

// Errorf logs messages at ERROR level.
func (l *logger) Errorf(format string, args ...interface{}) {
	log.Error(format, args...)
}

// Fatalf logs messages at FATAL level.
func (l *logger) Fatalf(format string, args ...interface{}) {
	log.Error(format, args...)
}
