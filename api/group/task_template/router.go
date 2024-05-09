package task_template

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/app/dbs"
	"schedule_task_command/util/my_log"
)

func RegisterRouter(g group) {
	taskS := g.GetScheduleServer().GetTaskServer()
	o := NewOperate(g.GetDbs(), taskS)
	log := my_log.NewLog("router/task_template")
	app := g.GetApp()

	tt := app.Group("/task_template")

	tt.Use(func(c *fiber.Ctx) error {
		c.Locals("Module", "schedule_module-task")
		return c.Next()
	})

	// subscribe to redis
	go func() {
		rdbSub(o, log)
	}()

	go func() {
		receiveStream(o, log)
	}()

	h := NewHandler(o, log)
	tt.Get("/", h.GetTaskTemplates)
	tt.Get("/:id", h.GetTaskTemplateById)
	tt.Post("/", h.AddTaskTemplate)
	tt.Patch("/", h.UpdateTaskTemplate)
	tt.Delete("/", h.DeleteTaskTemplate)
	tt.Post("/execute/:id", h.ExecuteTask)
}

type group interface {
	GetApp() fiber.Router
	GetDbs() dbs.Dbs
	GetScheduleServer() api.ScheduleSer
}
