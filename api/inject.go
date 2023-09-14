package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"io"
	"os"
	"schedule_task_command/api/group"
	"schedule_task_command/api/group/command_template"
	"schedule_task_command/api/group/header_template"
	"schedule_task_command/api/group/ping"
	"schedule_task_command/api/group/schedule"
	"schedule_task_command/api/group/task_template"
	"schedule_task_command/api/group/time_template"
	"schedule_task_command/app/dbs"
	"schedule_task_command/app/schedule_server"
	"schedule_task_command/util/logFile"
)

func Inject(app *fiber.App, dbs dbs.Dbs, ss schedule_server.ScheduleSer) {
	// Middleware
	log := logFile.NewLogFile("api", "inject.log")
	fiberLog := getFiberLogFile(log)
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Output: fiberLog,
	}))

	//swagger routes
	app.Get("/swagger/*", swagger.HandlerDefault)

	// api group add cors middleware
	api := app.Group("/api", cors.New())

	// create new group
	g := group.NewAPIGroup(api, dbs, ss)

	// model registration
	ping.RegisterRouter(g)
	command_template.RegisterRouter(g)
	header_template.RegisterRouter(g)
	schedule.RegisterRouter(g)
	task_template.RegisterRouter(g)
	time_template.RegisterRouter(g)
}

func getFiberLogFile(log logFile.LogFile) io.Writer {
	fiberFile, err := os.OpenFile("./log/fiber.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error().Fatal("can not open log file: " + err.Error())
	}
	return io.MultiWriter(fiberFile, os.Stdout)
}
