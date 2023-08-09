package routine

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLogger_Printf(t *testing.T) {
	assert.NotPanics(t, func() {
		logger.Printf("111 %s", "222")
	})
}

func TestNewAnts(t *testing.T) {
	r := require.New(t)
	a := newAnts()
	r.NotNil(a)
	a = newAnts(WithPoolSize(10))
	r.NotNil(a)
}

func TestAnts_Open(t *testing.T) {
	r := require.New(t)
	a := newAnts()
	r.NotNil(a)

	defer func() {
		r.Nil(a.Close())
	}()
	// 测试正常开启
	err := a.Open()
	r.Nil(err)

	// 测试重复开启不报错
	err = a.Open()
	r.Nil(err)

	// 测试关闭后再开启
	err = a.Close()
	r.Nil(err)
	err = a.Open()
	r.Nil(err)
}

func TestAnts_Close(t *testing.T) {
	r := require.New(t)
	a := newAnts()
	r.NotNil(a)

	// 测试正常关闭
	err := a.Close()
	r.Nil(err)

	// 测试开启后关闭
	err = a.Open()
	r.Nil(err)
	err = a.Close()
	r.Nil(err)

	// 测试重复关闭
	err = a.Close()
	r.Nil(err)
}

func TestAnts_Run(t *testing.T) {
	r := require.New(t)
	a := newAnts()
	r.NotNil(a)
	defer func() {
		r.Nil(a.Close())
	}()

	x := 0
	err := a.Run(false, func() {
		hh := 1
		hh++
	})
	r.Equal(ErrManagerClosed, err)

	err = a.Open()
	r.Nil(err)
	err = a.Run(false, func() {
		x++
	})
	r.Nil(err)
	err = a.Run(true, func() {
		panic("test")
	})
	r.Nil(err)
	time.Sleep(time.Millisecond * 20)
	r.Equal(1, x)
}

func TestAnts_GetRoutineNumber(t *testing.T) {
	r := require.New(t)
	a := newAnts()
	r.NotNil(a)
	defer func() {
		r.Nil(a.Close())
	}()

	r.Equal(0, a.GetRoutineNumber())

	err := a.Open()
	r.Nil(err)

	err = a.Run(false, func() {
		time.Sleep(time.Millisecond * 70)
	})
	time.Sleep(time.Millisecond * 20) // golang 支持的 Sleep 最小颗粒为 15-20ms，再小的间隔没有意义
	num := a.GetRoutineNumber()
	r.Nil(err)
	r.Equal(1, num)

	err = a.Run(false, func() {
		time.Sleep(time.Millisecond * 500)
	})
	time.Sleep(time.Millisecond * 20)
	num = a.GetRoutineNumber()
	r.Nil(err)
	r.Equal(2, num)

	time.Sleep(time.Millisecond * 50)
	r.Equal(1, a.GetRoutineNumber())
}
