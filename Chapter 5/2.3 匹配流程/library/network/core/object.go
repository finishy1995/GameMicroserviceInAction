package core

// Object 持续存在的对象，需要注意生命周期管理
type Object interface {
	// Run 持续执行逻辑，Run 函数主逻辑不需要使用Goroutine，而应该在外层调用时决定是否使用Goroutine
	Run()
	// Close 关闭对象푍
	Close()
}
