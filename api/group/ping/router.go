package ping

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/app/dbs"
	"schedule_task_command/util/logFile"
)

func RegisterRouter(g group) {
	o := NewOperate(g.GetDbs())
	log := logFile.NewLogFile("router", "ping.log")
	app := g.GetApp()

	ping := app.Group("/ping")

	h := NewHandler(o, log)
	ping.Get("/test", h.GetPing)
	ping.Get("/list", h.GetListPing)
}

type group interface {
	GetApp() fiber.Router
	GetDbs() dbs.Dbs
}
