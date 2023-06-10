package routine

import (
	"ProjectX/library/log"
	"errors"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync/atomic"
)

// Ants 由 panjf2000/ants 实现的线程池
type Ants struct {
	// commonPool 常规 Goroutine 池
	pool *ants.Pool
	// currentTask 当前任务
	currentTask int32
	// isOpen 是否打开
	isOpen bool
	// maxRoutine 最大协程数
	maxRoutine int
}

var (
	logger       = new(Logger)
	panicHandler = func(i interface{}) {
		panic(i)
	}
)

// Logger 日志处理，将 ants 自带的日志库覆盖为项目内的日志库
type Logger struct {
}

// Printf 日志打印
func (l *Logger) Printf(format string, args ...interface{}) {
	log.Info(format, args...)
}

func newAnts(opts ...Option) Manager {
	options := DefaultManagerOptions
	for _, o := range opts {
		o(&options)
	}

	m := &Ants{
		maxRoutine:  -1,
		currentTask: 0,
		isOpen:      false,
	}
	m.maxRoutine = -1
	if options.PoolSize > 0 {
		m.maxRoutine = int(options.PoolSize)
	}
	ants.Release()

	return m
}

func (a *Ants) Open() error {
	if a.isOpen {
		return nil
	}

	a.isOpen = true
	if a.pool != nil {
		a.pool.Release()
	}
	p, err := ants.NewPool(
		a.maxRoutine,
		ants.WithPanicHandler(panicHandler),
		ants.WithLogger(logger),
		ants.WithNonblocking(true))
	if err != nil {
		a.isOpen = false
		return errors.New(fmt.Sprintf("routine open Ants failed. error: %s", err.Error()))
	}

	a.pool = p
	a.currentTask = 0
	return nil
}

func (a *Ants) Close() error {
	if !a.isOpen {
		return nil
	}

	a.isOpen = false
	if a.pool != nil {
		a.pool.Release()
		a.pool = nil
	}

	return nil
}

func (a *Ants) GetRoutineNumber() int {
	if a.pool == nil {
		return 0
	}
	return int(a.currentTask)
}

func (a *Ants) Run(needRecover bool, task Task) error {
	if a.pool == nil {
		return ErrManagerClosed
	}
	var realTask Task
	if needRecover {
		realTask = func() {
			defer Catch()
			defer func() {
				atomic.AddInt32(&a.currentTask, -1)
			}()
			atomic.AddInt32(&a.currentTask, 1)
			task()
		}
	} else {
		realTask = func() {
			defer func() {
				atomic.AddInt32(&a.currentTask, -1)
			}()
			atomic.AddInt32(&a.currentTask, 1)
			task()
		}
	}

	return a.pool.Submit(realTask)
}
