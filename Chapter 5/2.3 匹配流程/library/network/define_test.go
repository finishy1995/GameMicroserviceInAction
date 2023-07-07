package network

import (
	"ProjectX/library/network/core"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInfo_SupportClientServer(t *testing.T) {
	i := new(info)
	i.client = func() core.Client {
		return nil
	}
	i.server = nil
	r := require.New(t)
	r.True(i.SupportClient())
	r.False(i.SupportServer())
}
