package task_template

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/app/dbs"
	"schedule_task_command/util/logFile"
)

func RegisterRouter(g group) {
	taskS := g.GetScheduleServer().GetTaskServer()
	o := NewOperate(g.GetDbs(), taskS)
	log := logFile.NewLogFile("router", "task_template.log")
	app := g.GetApp()

	tt := app.Group("/task_template")

	// subscribe to redis
	go func() {
		rdbSub(o, log)
	}()

	h := NewHandler(o, log)
	tt.Get("/", h.GetTaskTemplates)
	tt.Get("/:id", h.GetTaskTemplateById)
	tt.Post("/", h.AddTaskTemplate)
	tt.Patch("/", h.UpdateTaskTemplate)
	tt.Delete("/", h.DeleteTaskTemplate)
	tt.Post("/execute/:id", h.ExecuteTask)
}

type group interface {
	GetApp() fiber.Router
	GetDbs() dbs.Dbs
	GetScheduleServer() api.ScheduleSer
}