package websocket_hub

import (
	"errors"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/entry/e_module"
	"github.com/littlebluewhite/schedule_task_command/util/my_log"
)

type hub interface {
	Run()
	Broadcast([]byte)
	WsConnect(client *websocket.Conn)
}

type HubManager struct {
	hubs map[e_module.Module]hub
	l    api.Logger
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[e_module.Module]hub),
		l:    my_log.NewLog("websocket/hub_manager"),
	}
}

func (hm *HubManager) RegisterHub(module e_module.Module) {
	h := NewHub(module)
	hm.hubs[module] = h
	go h.Run()
}

func (hm *HubManager) Broadcast(module e_module.Module, msg []byte) {
	hm.l.Infof("module: %s hub broadcast", module)
	hm.hubs[module].Broadcast(msg)
}

func (hm *HubManager) WsConnect(module e_module.Module, conn *websocket.Conn) error {
	if h, ok := hm.hubs[module]; !ok {
		eString := fmt.Sprintf("module %s not exist", module)
		hm.l.Errorln(eString)
		return errors.New(eString)
	} else {
		h.WsConnect(conn)
	}
	return nil
}
