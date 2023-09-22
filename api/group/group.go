package group

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/app/dbs"
)

type Group struct {
	app fiber.Router
	dbs dbs.Dbs
	ss  api.ScheduleSer
}

func NewAPIGroup(app fiber.Router, dbs dbs.Dbs, ss api.ScheduleSer) *Group {
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

func (g *Group) GetScheduleServer() api.ScheduleSer {
	return g.ss
}
