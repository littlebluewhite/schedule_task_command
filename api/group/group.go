package group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/littlebluewhite/schedule_task_command/api"
)

type Group struct {
	app fiber.Router
	dbs api.Dbs
	ss  api.ScheduleSer
	hm  api.HubManager
}

func NewAPIGroup(app fiber.Router, dbs api.Dbs, ss api.ScheduleSer, hm api.HubManager) *Group {
	return &Group{
		app: app,
		dbs: dbs,
		ss:  ss,
		hm:  hm,
	}
}

func (g *Group) GetApp() fiber.Router {
	return g.app
}

func (g *Group) GetDbs() api.Dbs {
	return g.dbs
}

func (g *Group) GetScheduleServer() api.ScheduleSer {
	return g.ss
}

func (g *Group) GetWebsocketHub() api.HubManager {
	return g.hm
}
