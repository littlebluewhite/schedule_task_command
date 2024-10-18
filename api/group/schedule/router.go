package schedule

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/util/my_log"
)

func RegisterRouter(g group) {
	o := NewOperate(g.GetDbs(), g.GetScheduleServer())
	log := my_log.NewLog("router/schedule")
	app := g.GetApp()

	s := app.Group("/schedule")

	s.Use(func(c *fiber.Ctx) error {
		c.Locals("Module", "schedule_module-schedule")
		return c.Next()
	})

	h := NewHandler(o, log)
	s.Get("/", h.GetSchedules)
	s.Get("/:id", h.GetScheduleById)
	s.Post("/", h.AddSchedule)
	s.Patch("/", h.UpdateSchedule)
	s.Delete("/", h.DeleteSchedule)
}

type group interface {
	GetApp() fiber.Router
	GetDbs() api.Dbs
	GetScheduleServer() api.ScheduleSer
}
