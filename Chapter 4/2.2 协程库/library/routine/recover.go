package routine

import (
	"ProjectX/library/log"
	"bytes"
	"errors"
	"fmt"
	"runtime"
)

// Catch 协程恢复处理
func Catch() {
	if r := recover(); r != nil {
		Handle(r)
	}
}

// Trace 获取堆栈信息
// 详见 runtime/debug.Stack()
func Trace(skip int) []byte {
	si := make([]uintptr, 0, 5)
	pc := make([]uintptr, 10)
	skip += 2
	for {
		n := runtime.Callers(skip, pc)
		if n == 0 {
			break
		}
		skip += n
		si = append(si, pc[0:n]...)
	}

	var buf = new(bytes.Buffer)
	for i := 0; i < len(si); i++ {
		buf.WriteString("\t")
		name, file, line := function(si[i])
		_, err := fmt.Fprintf(buf, "at %s() [%s:%d]\n", name, file, line)
		if err != nil {
			log.Error("bytes.Buffer IO error: %v", err)
		}
	}
	return buf.Bytes()
}

var (
	dunno     = []byte("???")
	dot       = []byte(".")
	centerDot = []byte("·")
	slash     = []byte("/")
)

// function 返回调用的函数名、文件名和行数
func function(pc uintptr) (name []byte, file string, line int) {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno, "???", 0
	}
	file, line = fn.FileLine(pc)
	name = []byte(fn.Name())

	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return
}

// Handle 恢复执行函数
func Handle(r interface{}) {
	var err error
	switch x := r.(type) {
	case string:
		err = errors.New(x)
	case error:
		err = x
	default:
		err = ErrInvalidPanicError
	}

	// 跳过 Catch(), Handle() 这两层必出现的堆栈，所以 skip = 2
	log.Error("goroutine panic: %v\n%s\n", err, Trace(2))
}
