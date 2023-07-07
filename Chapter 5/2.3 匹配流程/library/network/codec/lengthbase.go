package codec

import (
	"ProjectX/library/network/core"
	"encoding/binary"
)

type LengthFieldBasedFrameCodec struct {
}

type innerBuffer []byte

func (in *innerBuffer) readN(n int) (buf []byte, err error) {
	if n == 0 {
		return nil, nil
	}

	if n < 0 {
		return nil, core.ErrTooLessLength
	} else if n > len(*in) {
		return nil, core.ErrTooMoreLength
	}
	buf = (*in)[:n]
	*in = (*in)[n:]
	return
}

// Encode ...
func (cc *LengthFieldBasedFrameCodec) Encode(_ core.CodecConn, buf []byte) (out []byte, err error) {
	length := len(buf)
	if length < 0 {
		return nil, core.ErrTooLessLength
	}

	out = make([]byte, 4)
	binary.BigEndian.PutUint32(out, uint32(length))

	out = append(out, buf...)
	return
}

// Decode ...
func (cc *LengthFieldBasedFrameCodec) Decode(c core.CodecConn) ([]byte, error) {
	var (
		in          innerBuffer
		err         error
		frameLength uint64
	)
	in = c.Read()

	lenBuf, err := in.readN(4)
	if err != nil {
		return nil, core.ErrPacketSplit
	} else {
		frameLength = uint64(binary.BigEndian.Uint32(lenBuf))
	}

	// real message length
	msgLength := int(frameLength)
	msg, err := in.readN(msgLength)
	if err != nil {
		return nil, core.ErrPacketSplit
	}

	fullMessage := make([]byte, msgLength)
	copy(fullMessage, msg)
	c.ShiftN(msgLength + 4)
	return fullMessage, nil
}
