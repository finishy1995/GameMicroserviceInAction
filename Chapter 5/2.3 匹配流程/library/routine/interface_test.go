package routine

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	assert.Nil(t, Open())
}

func TestClose(t *testing.T) {
	assert.Nil(t, Close())
}

func TestGetRoutineNumber(t *testing.T) {
	assert.Equal(t, 0, GetRoutineNumber())
}

func TestRun(t *testing.T) {
	r := require.New(t)
	r.Nil(Open())

	x := 10
	err := Run(false, func() {
		x++
		time.Sleep(time.Millisecond * 50)
	})
	r.Nil(err)
	time.Sleep(time.Millisecond * 20)
	r.Equal(1, GetRoutineNumber())
	r.Equal(11, x)
}

func TestWithPoolSize(t *testing.T) {
	assert.NotNil(t, WithPoolSize(-1))
}

func TestRegisterManagerType(t *testing.T) {
	RegisterManagerType("TestRegisterManagerType", func(opts ...Option) Manager {
		return nil
	})
	// 测试重复注册同一个名称
	RegisterManagerType("TestRegisterManagerType", func(opts ...Option) Manager {
		return nil
	})
}

func TestSetManager(t *testing.T) {
	r := require.New(t)

	err := SetManager("wrong")
	r.Equal(ErrInvalidManager, err)
	RegisterManagerType("test", func(opts ...Option) Manager {
		return nil
	})
	err = SetManager("test")
	r.Equal(ErrInvalidManager, err)

	err = SetManager("")
	r.Nil(err)
}
