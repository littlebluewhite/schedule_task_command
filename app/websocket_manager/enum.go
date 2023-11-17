package websocket_manager

type Group int

const (
	None Group = iota
	Command
	Task
)
