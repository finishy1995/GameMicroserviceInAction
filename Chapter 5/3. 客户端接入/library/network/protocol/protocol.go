package protocol

type protocol struct {
	ReceiveMsg []byte
	ReplyMsg   []byte
	Flag       int
}

var ProtocolV001 = protocol{
	ReceiveMsg: []byte("//"),
	ReplyMsg:   []byte("Hi"),
	Flag:       2,
}
