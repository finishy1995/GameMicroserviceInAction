package core

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVerifyAddress(t *testing.T) {
	r := require.New(t)

	// 检查格式
	r.False(VerifyAddress("127.0.0.1:50:"))
	r.False(VerifyAddress(":127.0.0.1:50"))
	r.True(VerifyAddress(":50"))
	r.True(VerifyAddress("0.0.0.0:50"))

	// 检查 ip
	r.False(VerifyAddress("0.0.0:60"))
	r.False(VerifyAddress("0.0.0.0.0:60"))
	r.True(VerifyAddress("[2001:db8::68]:60"))

	// 检查端口
	r.True(VerifyAddress(":660"))
	r.False(VerifyAddress(":100000"))
	r.False(VerifyAddress(":1a1"))
}

func TestGenerateID(t *testing.T) {
	n := 10000
	ids := make(map[ID]bool, n)
	for i := 0; i < n; i++ {
		id := GenerateID()
		if ids[id] {
			assert.FailNow(t, "duplicate ID", id)
		}
		ids[id] = true
	}
}
