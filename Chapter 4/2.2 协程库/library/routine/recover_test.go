package routine

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCatch(t *testing.T) {
	assert.NotPanics(t, func() {
		defer Catch()
		panic("111")
	}, "should not panic")
}

func TestTrace(t *testing.T) {
	bytes := Trace(0)
	assert.Contains(t, string(bytes), "TestTrace")
}

func TestHandle(t *testing.T) {
	r := require.New(t)
	r.NotPanics(func() {
		Handle(errors.New("test"))
	})
	r.NotPanics(func() {
		Handle("test")
	})
	r.NotPanics(func() {
		Handle(111)
	})
}
