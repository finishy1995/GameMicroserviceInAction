package codec

import (
	"ProjectX/library/buffer/bytebuffer"
	"ProjectX/library/network/core"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLengthFieldBasedFrameCodec(t *testing.T) {
	codec := new(LengthFieldBasedFrameCodec)
	conn := newMockConn()
	r := require.New(t)

	test1 := []byte("hello world")
	bb, err := codec.Encode(conn, test1)
	r.Nil(err)
	conn.MockGetNetworkMsg(bb)
	bb, err = codec.Decode(conn)
	r.Nil(err)
	r.Equal(test1, bb)

	test2 := []byte("")
	bb, err = codec.Encode(conn, test2)
	r.Nil(err)
	conn.MockGetNetworkMsg(bb)
	bb, err = codec.Decode(conn)
	r.Nil(err)
	r.Equal(test2, bb)

	// 测试粘包拆包，测试发三个包，客户端先后收到两个小包、一个大包、一个小包
	test3 := []byte("hellooooo")
	test4 := []byte(" ")
	test5 := []byte("worlddddd")
	bytes := bytebuffer.Get()
	bb, err = codec.Encode(conn, test3)
	r.Nil(err)
	bytes.Write(bb)
	bb, err = codec.Encode(conn, test4)
	r.Nil(err)
	bytes.Write(bb)
	bb, err = codec.Encode(conn, test5)
	r.Nil(err)
	bytes.Write(bb)
	b := bytes.Bytes()
	lengthB := len(b)
	conn.MockGetNetworkMsg(b[:2])
	bb, err = codec.Decode(conn)
	r.Equal(core.ErrPacketSplit, err)
	conn.MockGetNetworkMsg(b[2:6])
	bb, err = codec.Decode(conn)
	r.Equal(core.ErrPacketSplit, err)
	conn.MockGetNetworkMsg(b[6 : lengthB-3])
	bb, err = codec.Decode(conn)
	r.Nil(err)
	r.Equal(test3, bb)
	bb, err = codec.Decode(conn)
	r.Nil(err)
	r.Equal(test4, bb)
	bb, err = codec.Decode(conn)
	r.Equal(core.ErrPacketSplit, err)
	conn.MockGetNetworkMsg(b[lengthB-3:])
	bb, err = codec.Decode(conn)
	r.Nil(err)
	r.Equal(test5, bb)
}
