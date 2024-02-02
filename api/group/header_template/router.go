package header_template

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/app/dbs"
	"schedule_task_command/util/logFile"
)

func RegisterRouter(g group) {
	o := NewOperate(g.GetDbs())
	log := logFile.NewLogFile("router", "header_template.log")
	app := g.GetApp()

	ht := app.Group("/header_template")

	ht.Use(func(c *fiber.Ctx) error {
		c.Locals("Module", "command_server-header_template")
		return c.Next()
	})

	h := NewHandler(o, log)
	ht.Get("/", h.GetheaderTemplates)
	ht.Get("/:id", h.GetHeaderTemplateById)
	ht.Post("/", h.AddHeaderTemplate)
	ht.Patch("/", h.UpdateHeaderTemplate)
	ht.Delete("/", h.DeleteHeaderTemplate)
}

type group interface {
	GetApp() fiber.Router
	GetDbs() dbs.Dbs
}
