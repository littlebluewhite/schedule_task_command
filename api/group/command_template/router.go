package command_template

import (
	"github.com/gofiber/fiber/v2"
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/api/group/task_template"
	"github.com/littlebluewhite/schedule_task_command/util/my_log"
)

func RegisterRouter(g group) {
	taskS := g.GetScheduleServer().GetTaskServer()
	commandS := taskS.GetCommandServer()
	taskTemplateO := task_template.NewOperate(g.GetDbs(), taskS)
	o := NewOperate(g.GetDbs(), commandS, taskTemplateO)
	log := my_log.NewLog("router/command_template")
	app := g.GetApp()

	ct := app.Group("/command_template")

	ct.Use(func(c *fiber.Ctx) error {
		c.Locals("Module", "schedule_module-command")
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
	ct.Post("/send/", h.SendCommandTemplate)
}

type group interface {
	GetApp() fiber.Router
	GetDbs() api.Dbs
	GetScheduleServer() api.ScheduleSer
}
