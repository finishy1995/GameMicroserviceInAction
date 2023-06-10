package core

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithServerOpt(t *testing.T) {
	context := map[string]interface{}{"a": 1}
	opt1 := WithMaxConnNum(10)
	opt2 := WithServerContext(context)
	options := DefaultServerOptions
	opt1(&options)
	opt2(&options)
	r := require.New(t)

	r.Equal(10, options.MaxConnNum)
	r.Equal(context, options.Context)
}
