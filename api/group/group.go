package group

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/app/dbs"
	"schedule_task_command/app/time_server"
)

type Group struct {
	app fiber.Router
	dbs dbs.Dbs
	ts  time_server.TimeServer
}

func NewAPIGroup(app fiber.Router, dbs dbs.Dbs, ts time_server.TimeServer) *Group {
	return &Group{
		app: app,
		dbs: dbs,
		ts:  ts,
	}
}

func (g *Group) GetApp() fiber.Router {
	return g.app
}

func (g *Group) GetDbs() dbs.Dbs {
	return g.dbs
}

func (g *Group) GetTimeServer() time_server.TimeServer {
	return g.ts
}
