package routine

import "github.com/finishy1995/codegenerator/library/log"

// Task Goroutine 需要执行的任务
type Task func()

// NewManagerFunc 创建协程控制器实例的函数
type NewManagerFunc func(opts ...Option) Manager

// Manager 协程控制器
type Manager interface {
	// Open 开启控制器
	Open() error
	// Close 关闭控制器
	Close() error
	// GetRoutineNumber 获取当前正在运行的协程数量
	GetRoutineNumber() int
	// Run 把任务交给协程处理
	Run(needRecover bool, task Task) error
}

// Options 选项
type Options struct {
	// PoolSize 协程池大小，负数代表无上限限制
	PoolSize int32
}

// Option 选项闭包
type Option func(*Options)

// WithPoolSize 设置协程池大小，负数代表无上限限制
func WithPoolSize(size int32) Option {
	return func(options *Options) {
		options.PoolSize = size
	}
}

var (
	managerInstance Manager
	managerTypeMap  = map[string]NewManagerFunc{
		"default": newAnts,
		"":        newAnts,
	}

	// DefaultManagerOptions 默认 Manager 选项
	DefaultManagerOptions = Options{
		PoolSize: -1,
	}
)

func init() {
	err := SetManager("")
	if err != nil {
		panic(err)
	}
	err = Open()
	if err != nil {
		panic(err)
	}
}

func SetManager(typ string, opts ...Option) error {
	f, ok := managerTypeMap[typ]
	if !ok {
		return ErrInvalidManager
	}
	newInstance := f(opts...)
	if newInstance == nil {
		return ErrInvalidManager
	}
	if managerInstance != nil {
		err := managerInstance.Close()
		if err != nil {
			log.Error("close old routine manager failed, error: %s", err.Error())
		}
	}
	managerInstance = newInstance
	return nil
}

// RegisterManagerType 注册一个新的协程控制器类型
//
//	通过这个函数，可以轻松扩展这个库
func RegisterManagerType(name string, f NewManagerFunc) {
	if _, ok := managerTypeMap[name]; ok {
		log.Warning("replace exist routine manager type %s", name)
	}
	managerTypeMap[name] = f
}

// Open 开启控制器
func Open() error {
	return managerInstance.Open()
}

// Close 关闭控制器
func Close() error {
	return managerInstance.Close()
}

// GetRoutineNumber 获取当前正在运行的协程数量
func GetRoutineNumber() int {
	return managerInstance.GetRoutineNumber()
}

// Run 把任务交给协程处理
func Run(needRecover bool, task Task) error {
	return managerInstance.Run(needRecover, task)
}
