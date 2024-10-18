package header_template

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/util/my_log"
)

func RegisterRouter(g group) {
	o := NewOperate(g.GetDbs())
	log := my_log.NewLog("router/header_template")
	app := g.GetApp()

	ht := app.Group("/header_template")

	ht.Use(func(c *fiber.Ctx) error {
		c.Locals("Module", "schedule_module-command")
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
	GetDbs() api.Dbs
}
