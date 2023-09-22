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

	tt := app.Group("/command")

	h := NewHandler(o, log)
	tt.Get("/", h.GetCommands)
	tt.Get("/:commandId", h.GetCommandByCommandId)
	tt.Delete("/:commandId", h.CancelCommand)
	tt.Patch("/history", h.GetHistory)
}

type group interface {
	GetApp() fiber.Router
	GetScheduleServer() api.ScheduleSer
}
