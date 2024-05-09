package task

import (
	"github.com/gofiber/fiber/v2"
	"schedule_task_command/api"
	"schedule_task_command/util/my_log"
)

func RegisterRouter(g group) {
	taskS := g.GetScheduleServer().GetTaskServer()
	o := NewOperate(taskS)
	log := my_log.NewLog("router/task")
	app := g.GetApp()

	tt := app.Group("/task")

	tt.Use(func(c *fiber.Ctx) error {
		c.Locals("Module", "schedule_module-task")
		return c.Next()
	})

	h := NewHandler(o, log)
	tt.Get("/", h.GetTasks)
	tt.Get("/simple/", h.GetSimpleTasks)
	tt.Get("/history", h.GetHistory)
	tt.Get("/stage_item/status/:id", h.GetStageItemStatus)
	tt.Get("/simple/:id", h.GetSimpleTasksById)
	tt.Get("/:id", h.GetTaskById)
	tt.Delete("/:id", h.CancelTask)
}

type group interface {
	GetApp() fiber.Router
	GetScheduleServer() api.ScheduleSer
}
