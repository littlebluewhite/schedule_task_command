package task

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/util/logFile"
)

func RegisterRouter(g group) {
	taskS := g.GetScheduleServer().GetTaskServer()
	o := NewOperate(taskS)
	log := logFile.NewLogFile("router", "task.log")
	app := g.GetApp()

	tt := app.Group("/task")

	h := NewHandler(o, log)
	tt.Get("/", h.GetTasks)
	tt.Get("/:taskId", h.GetTaskByTaskId)
	tt.Delete("/:taskId", h.CancelTask)
	tt.Patch("/history", h.GetHistory)
}

type group interface {
	GetApp() fiber.Router
	GetScheduleServer() api.ScheduleSer
}
