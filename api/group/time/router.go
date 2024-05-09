package time

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/util/my_log"
)

func RegisterRouter(g group) {
	timeS := g.GetScheduleServer().GetTimeServer()
	o := NewOperate(timeS)
	log := my_log.NewLog("router/time")
	app := g.GetApp()

	tt := app.Group("/time")

	tt.Use(func(c *fiber.Ctx) error {
		c.Locals("Module", "schedule_module-time")
		return c.Next()
	})

	h := NewHandler(o, log)
	tt.Get("/history", h.GetHistory)
}

type group interface {
	GetApp() fiber.Router
	GetScheduleServer() api.ScheduleSer
}
