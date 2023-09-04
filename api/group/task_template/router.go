package task_template

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/app/dbs"
	"schedule_task_command/util/logFile"
)

func RegisterRouter(g group) {
	o := NewOperate(g.GetDbs())
	log := logFile.NewLogFile("router", "task_template.json.log")
	app := g.GetApp()

	tt := app.Group("/task_template")

	h := NewHandler(o, log)
	tt.Get("/", h.GetTaskTemplates)
	tt.Get("/:id", h.GetTaskTemplateById)
	tt.Post("/", h.AddTaskTemplate)
	tt.Patch("/", h.UpdateTaskTemplate)
	tt.Delete("/", h.DeleteTaskTemplate)
}

type group interface {
	GetApp() fiber.Router
	GetDbs() dbs.Dbs
}
