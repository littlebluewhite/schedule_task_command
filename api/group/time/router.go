package time

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/util/logFile"
)

func RegisterRouter(g group) {
	timeS := g.GetScheduleServer().GetTimeServer()
	o := NewOperate(timeS)
	log := logFile.NewLogFile("router", "time.log")
	app := g.GetApp()

	tt := app.Group("/time")

	h := NewHandler(o, log)
	tt.Get("/history", h.GetHistory)
}

type group interface {
	GetApp() fiber.Router
	GetScheduleServer() api.ScheduleSer
}
