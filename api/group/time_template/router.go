package time_template

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/app/dbs"
	"schedule_task_command/util/logFile"
)

func RegisterRouter(g group) {
	o := NewOperate(g.GetDbs())
	log := logFile.NewLogFile("router", "time_template.log")
	app := g.GetApp()

	tt := app.Group("/time_template")

	h := NewHandler(o, log)
	tt.Get("/", h.GetTimeTemplates)
	tt.Get("/:id", h.GetTimeTemplateById)
	tt.Post("/", h.AddTimeTemplate)
	tt.Patch("/", h.UpdateTimeTemplate)
	tt.Delete("/", h.DeleteTimeTemplate)

}

type group interface {
	GetApp() fiber.Router
	GetDbs() dbs.Dbs
}
