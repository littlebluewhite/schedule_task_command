package group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/api/group/command"
	"github.com/littlebluewhite/schedule_task_command/api/group/command_template"
	"github.com/littlebluewhite/schedule_task_command/api/group/header_template"
	"github.com/littlebluewhite/schedule_task_command/api/group/ping"
	"github.com/littlebluewhite/schedule_task_command/api/group/schedule"
	"github.com/littlebluewhite/schedule_task_command/api/group/task"
	"github.com/littlebluewhite/schedule_task_command/api/group/task_template"
	"github.com/littlebluewhite/schedule_task_command/api/group/time"
	"github.com/littlebluewhite/schedule_task_command/api/group/time_template"
	"github.com/littlebluewhite/schedule_task_command/api/group/ws"
	"github.com/littlebluewhite/schedule_task_command/util/my_log"
	"io"
	"os"
)

func Inject(app *fiber.App, dbs api.Dbs, ss api.ScheduleSer, hm api.HubManager) {
	// Middleware
	log := my_log.NewLog("api/group")
	fiberLog := getFiberLogFile(log)
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Output: fiberLog,
	}))

	//swagger routes
	app.Get("/swagger/*", swagger.HandlerDefault)

	// api group add cors middleware
	Api := app.Group("/api", cors.New())
	Api.Use(func(c *fiber.Ctx) error {
		c.Locals("Module", "logs")
		return c.Next()
	})

	// use middleware to write my_log
	o := NewOperate(dbs)
	h := NewHandler(o, log)
	Api.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		o.WriteLog(c)
		if err != nil {
			log.Errorln(err)
		}
		return err
	})
	Api.Get("/logs", h.GetHistory)

	// create new group
	g := NewAPIGroup(Api, dbs, ss, hm)

	// model registration
	ping.RegisterRouter(g)
	command_template.RegisterRouter(g)
	command.RegisterRouter(g)
	header_template.RegisterRouter(g)
	schedule.RegisterRouter(g)
	task_template.RegisterRouter(g)
	task.RegisterRouter(g)
	time_template.RegisterRouter(g)
	time.RegisterRouter(g)
	ws.RegisterRouter(g)
}

func getFiberLogFile(log api.Logger) io.Writer {
	fiberFile, err := os.OpenFile("./my_log/fiber.my_log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Errorf("can not open my_log file: " + err.Error())
	}
	return io.MultiWriter(fiberFile, os.Stdout)
}
