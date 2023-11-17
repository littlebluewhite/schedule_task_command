package websocket_manager

import "github.com/gofiber/contrib/websocket"

type groupConnect struct {
	group  Group
	client *websocket.Conn
}

type groupMessage struct {
	group   Group
	message []byte
}
