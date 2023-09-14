package group

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/app/dbs"
	"schedule_task_command/app/schedule_server"
)

type Group struct {
	app fiber.Router
	dbs dbs.Dbs
	ss  schedule_server.ScheduleSer
}

func NewAPIGroup(app fiber.Router, dbs dbs.Dbs, ss schedule_server.ScheduleSer) *Group {
	return &Group{
		app: app,
		dbs: dbs,
		ss:  ss,
	}
}

func (g *Group) GetApp() fiber.Router {
	return g.app
}

func (g *Group) GetDbs() dbs.Dbs {
	return g.dbs
}

func (g *Group) GetScheduleServer() schedule_server.ScheduleSer {
	return g.ss
}
