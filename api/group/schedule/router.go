package schedule

import (
	"github.com/gofiber/fiber/v2"
	group2 "schedule_task_command/api/group"
	"schedule_task_command/app/dbs"
	"schedule_task_command/util/logFile"
)

func RegisterRouter(g group) {
	o := NewOperate(g.GetDbs(), g.GetScheduleServer())
	log := logFile.NewLogFile("router", "schedule.log")
	app := g.GetApp()

	s := app.Group("/schedule")

	h := NewHandler(o, log)
	s.Get("/", h.GetSchedules)
	s.Get("/:id", h.GetScheduleById)
	s.Post("/", h.AddSchedule)
	s.Patch("/", h.UpdateSchedule)
	s.Delete("/", h.DeleteSchedule)
}

type group interface {
	GetApp() fiber.Router
	GetDbs() dbs.Dbs
	GetScheduleServer() group2.ScheduleSer
}
