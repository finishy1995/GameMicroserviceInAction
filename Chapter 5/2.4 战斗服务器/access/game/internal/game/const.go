package game

const (
	InitAction int32 = iota
	PlayerJoinAction
	PlayerLeaveAction
	StartAction
	TerminateAction
)

type Status uint8

const (
	Idle Status = iota
	Lobby
	BeforeStart
	Playing
	BeforeTerminate
)
