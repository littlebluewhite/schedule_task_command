package websocket_hub

import (
	"context"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"schedule_task_command/entry/e_module"
	"schedule_task_command/util/logFile"
)

type Hub struct {
	clients        map[*client]struct{}
	registerChan   chan *client
	unregisterChan chan *client
	broadcast      chan []byte
	l              logFile.LogFile
}

func NewHub(module e_module.Module) *Hub {
	name := fmt.Sprintf("%s_hub.log", module)
	return &Hub{
		clients:        make(map[*client]struct{}),
		registerChan:   make(chan *client),
		unregisterChan: make(chan *client),
		broadcast:      make(chan []byte),
		l:              logFile.NewLogFile("websocket", name),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cli := <-h.registerChan:
			h.clients[cli] = struct{}{}
		case cli := <-h.unregisterChan:
			if _, ok := h.clients[cli]; ok {
				delete(h.clients, cli)
			}
		case msg := <-h.broadcast:
			h.l.Info().Printf("broadcast 2")
			for cli := range h.clients {
				go func(cli *client) {
					cli.send(msg)
				}(cli)
			}
		}
	}
}

func (h *Hub) register(cli *client) {
	h.registerChan <- cli
}

func (h *Hub) unRegister(cli *client) {
	h.unregisterChan <- cli
}

func (h *Hub) Broadcast(msg []byte) {
	h.l.Info().Printf("broadcast 1")
	h.broadcast <- msg
}

func (h *Hub) WsConnect(conn *websocket.Conn) {
	ctx, cancel := context.WithCancel(context.Background())
	cli := newClient(conn, h.l)
	defer func() {
		h.unRegister(cli)
		cancel()
		cli.close()
	}()
	h.register(cli)
	go cli.writePump(ctx)
	cli.readPump()
}
