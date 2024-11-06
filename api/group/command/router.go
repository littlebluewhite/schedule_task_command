package command

import (
	"github.com/gofiber/fiber/v2"
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/util/my_log"
)

func RegisterRouter(g group) {
	commandS := g.GetScheduleServer().GetTaskServer().GetCommandServer()
	o := NewOperate(commandS)
	log := my_log.NewLog("router/command")
	app := g.GetApp()

	c := app.Group("/command")

	c.Use(func(c *fiber.Ctx) error {
		c.Locals("Module", "schedule_module-command")
		return c.Next()
	})

	h := NewHandler(o, log)
	c.Get("/", h.GetCommands)
	c.Get("/history", h.GetHistory)
	c.Get("/:id", h.GetCommandById)
	c.Delete("/:id", h.CancelCommand)
}

type group interface {
	GetApp() fiber.Router
	GetScheduleServer() api.ScheduleSer
}
