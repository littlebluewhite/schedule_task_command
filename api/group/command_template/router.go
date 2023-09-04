package command_template

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/app/dbs"
	"schedule_task_command/util/logFile"
)

func RegisterRouter(g group) {
	o := NewOperate(g.GetDbs())
	log := logFile.NewLogFile("router", "command_template.log")
	app := g.GetApp()

	ct := app.Group("/command_template")

	h := NewHandler(o, log)
	ct.Get("/", h.GetCommandTemplates)
	ct.Get("/:id", h.GetCommandTemplateById)
	ct.Post("/", h.AddCommandTemplate)
	ct.Delete("/", h.DeleteCommandTemplate)
}

type group interface {
	GetApp() fiber.Router
	GetDbs() dbs.Dbs
}
