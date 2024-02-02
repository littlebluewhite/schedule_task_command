package command_template

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/api/group/task_template"
	"schedule_task_command/app/dbs"
	"schedule_task_command/util/logFile"
)

func RegisterRouter(g group) {
	taskS := g.GetScheduleServer().GetTaskServer()
	commandS := taskS.GetCommandServer()
	taskTemplateO := task_template.NewOperate(g.GetDbs(), taskS)
	o := NewOperate(g.GetDbs(), commandS, taskTemplateO)
	log := logFile.NewLogFile("router", "command_template.log")
	app := g.GetApp()

	ct := app.Group("/command_template")

	ct.Use(func(c *fiber.Ctx) error {
		c.Locals("Module", "command_server-command_template")
		return c.Next()
	})

	// subscribe to redis
	go func() {
		rdbSub(o, log)
	}()

	h := NewHandler(o, log)
	ct.Get("/", h.GetCommandTemplates)
	ct.Get("/:id", h.GetCommandTemplateById)
	ct.Post("/", h.AddCommandTemplate)
	ct.Patch("/", h.UpdateCommandTemplate)
	ct.Delete("/", h.DeleteCommandTemplate)
	ct.Post("/execute/:id", h.ExecuteCommand)
}

type group interface {
	GetApp() fiber.Router
	GetDbs() dbs.Dbs
	GetScheduleServer() api.ScheduleSer
}
