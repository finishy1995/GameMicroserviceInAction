package core

// Codec 网络连接的 Codec。
//
// Codec 可以理解为协议，例：websocket 是 tcp 的一种 codec；如果你需要 websocket，那你只需要用 tcp 的网络库 + websocket codec
type Codec interface {
	// Encode 加密传输
	Encode(c CodecConn, buf []byte) ([]byte, error)
	// Decode 解密传输
	Decode(c CodecConn) ([]byte, error)
}
