package command

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/util/logFile"
)

func RegisterRouter(g group) {
	commandS := g.GetScheduleServer().GetTaskServer().GetCommandServer()
	o := NewOperate(commandS)
	log := logFile.NewLogFile("router", "command.log")
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
