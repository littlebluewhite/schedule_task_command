package time_template

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/app/dbs"
	"schedule_task_command/util/my_log"
)

func RegisterRouter(g group) {
	timeS := g.GetScheduleServer().GetTimeServer()
	o := NewOperate(g.GetDbs(), timeS)
	log := my_log.NewLog("router/time_template")
	app := g.GetApp()

	tt := app.Group("/time_template")

	tt.Use(func(c *fiber.Ctx) error {
		c.Locals("Module", "schedule_module-time")
		return c.Next()
	})

	// subscribe to redis
	go func() {
		rdbSub(o, log)
	}()

	go func() {
		receiveStream(o, log)
	}()

	h := NewHandler(o, log)
	tt.Get("/", h.GetTimeTemplates)
	tt.Get("/:id", h.GetTimeTemplateById)
	tt.Post("/", h.AddTimeTemplate)
	tt.Patch("/", h.UpdateTimeTemplate)
	tt.Delete("/", h.DeleteTimeTemplate)
	tt.Post("/checkTime/:id", h.CheckTime)
}

type group interface {
	GetApp() fiber.Router
	GetDbs() dbs.Dbs
	GetScheduleServer() api.ScheduleSer
}
