package core

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithClientOpt(t *testing.T) {
	context := map[string]interface{}{"a": 1}
	opt1 := WithReconnect(true)
	opt2 := WithClientContext(context)
	options := DefaultClientOptions
	opt1(&options)
	opt2(&options)
	r := require.New(t)

	r.Equal(true, options.Reconnect)
	r.Equal(context, options.Context)
}
