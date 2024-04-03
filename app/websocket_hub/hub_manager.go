package websocket_hub

import (
	"errors"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"schedule_task_command/entry/e_module"
	"schedule_task_command/util/logFile"
)

type hub interface {
	Run()
	Broadcast([]byte)
	WsConnect(client *websocket.Conn)
}

type HubManager struct {
	hubs map[e_module.Module]hub
	l    logFile.LogFile
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[e_module.Module]hub),
		l:    logFile.NewLogFile("websocket", "hub_manager.log"),
	}
}

func (hm *HubManager) RegisterHub(module e_module.Module) {
	h := NewHub(module)
	hm.hubs[module] = h
	go h.Run()
}

func (hm *HubManager) Broadcast(module e_module.Module, msg []byte) {
	hm.l.Info().Printf("module: %s hub broadcast", module)
	fmt.Println(hm.hubs)
	hm.hubs[module].Broadcast(msg)
}

func (hm *HubManager) WsConnect(module e_module.Module, conn *websocket.Conn) error {
	if h, ok := hm.hubs[module]; !ok {
		eString := fmt.Sprintf("module %s not exist", module)
		hm.l.Error().Println(eString)
		return errors.New(eString)
	} else {
		h.WsConnect(conn)
	}
	return nil
}
