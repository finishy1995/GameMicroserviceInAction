package codec

import "ProjectX/library/network/core"

var (
	zeroPacket = []byte("")
)

type BuiltInCodec struct {
}

// Encode ...
func (cc *BuiltInCodec) Encode(_ core.CodecConn, buf []byte) (_ []byte, _ error) {
	return buf, nil
}

// Decode ...
func (cc *BuiltInCodec) Decode(c core.CodecConn) ([]byte, error) {
	buf := c.Read()
	if len(buf) == 0 {
		return zeroPacket, nil
	}
	c.ResetBuffer()
	return buf, nil
}
