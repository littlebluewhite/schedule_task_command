package ping

import (
	"github.com/gofiber/fiber/v2"
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/util/my_log"
)

func RegisterRouter(g group) {
	o := NewOperate(g.GetDbs())
	log := my_log.NewLog("router/ping")
	app := g.GetApp()

	ping := app.Group("/ping")

	h := NewHandler(o, log)
	ping.Get("/test", h.GetPing)
	ping.Get("/list", h.GetListPing)
}

type group interface {
	GetApp() fiber.Router
	GetDbs() api.Dbs
}
