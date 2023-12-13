package websocket_manager

import (
	"github.com/gofiber/contrib/websocket"
	"schedule_task_command/util/logFile"
)

type WebsocketManager struct {
	groups     map[Group]map[*websocket.Conn]struct{}
	l          logFile.LogFile
	register   chan groupConnect
	unregister chan groupConnect
	broadcast  chan groupMessage
}

func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{
		groups: map[Group]map[*websocket.Conn]struct{}{
			None:    make(map[*websocket.Conn]struct{}),
			Command: make(map[*websocket.Conn]struct{}),
			Task:    make(map[*websocket.Conn]struct{}),
		},
		l:          logFile.NewLogFile("app", "websocket_manager"),
		register:   make(chan groupConnect),
		unregister: make(chan groupConnect),
		broadcast:  make(chan groupMessage),
	}
}

func (wm *WebsocketManager) Run() {
	for {
		select {
		case gc := <-wm.register:
			wm.l.Info().Printf("client: %v, register %d", gc.client, gc.group)
			wm.groups[gc.group][gc.client] = struct{}{}
		case gc := <-wm.unregister:
			wm.l.Info().Printf("client: %v, unregister %d", gc.client, gc.group)
			delete(wm.groups[gc.group], gc.client)
			_ = gc.client.Close()
		case gm := <-wm.broadcast:
			wm.l.Info().Printf("send message to %d, message: %s", gm.group, string(gm.message))
			for client := range wm.groups[gm.group] {
				err := client.WriteMessage(websocket.TextMessage, gm.message)
				if err != nil {
					delete(wm.groups[gm.group], client)
					_ = client.Close()
				}
			}
		}
	}
}

func (wm *WebsocketManager) Register(d int, client *websocket.Conn) {
	wm.register <- groupConnect{
		group:  Group(d),
		client: client,
	}
}

func (wm *WebsocketManager) Unregister(d int, client *websocket.Conn) {
	wm.unregister <- groupConnect{
		group:  Group(d),
		client: client,
	}
}

func (wm *WebsocketManager) Broadcast(d int, message []byte) {
	wm.broadcast <- groupMessage{
		group:   Group(d),
		message: message,
	}
}
