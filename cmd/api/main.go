package main

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	version "github.com/littlebluewhite/schedule_task_command"
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/api/group"
	"github.com/littlebluewhite/schedule_task_command/app/command_server"
	"github.com/littlebluewhite/schedule_task_command/app/dbs"
	"github.com/littlebluewhite/schedule_task_command/app/schedule_server"
	"github.com/littlebluewhite/schedule_task_command/app/task_server"
	"github.com/littlebluewhite/schedule_task_command/app/time_server"
	"github.com/littlebluewhite/schedule_task_command/app/websocket_hub"
	"github.com/littlebluewhite/schedule_task_command/docs"
	_ "github.com/littlebluewhite/schedule_task_command/docs"
	"github.com/littlebluewhite/schedule_task_command/proto/grpc_task_template"
	"github.com/littlebluewhite/schedule_task_command/util/config"
	"github.com/littlebluewhite/schedule_task_command/util/my_log"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// @title           Schedule-Task-Command swagger API
// @version         0
// @description     This is a schedule-command server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Wilson
// @contact.url    https://github.com/littlebluewhite
// @contact.email  wwilson008@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      127.0.0.1:5487

func main() {
	mainLog := my_log.NewLog("main")
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mainLog.Infoln("command module start")

	// read config
	Config := config.NewConfig[config.Config]("./config", "config", config.Yaml)

	ServerConfig := Config.Server

	// swagger docs host
	docsHost := fmt.Sprintf("%s:%s", ServerConfig.SwaggerHost, ServerConfig.Port)
	docs.SwaggerInfo.Host = docsHost
	docs.SwaggerInfo.Version = version.Version

	// DBs start includes SQL Cache
	DBS := dbs.NewDbs(mainLog, false, Config)
	defer func() {
		DBS.GetIdb().Close()
		mainLog.Infoln("influxDB Disconnect")
	}()

	// create websocket manager
	hm := websocket_hub.NewHubManager()

	// create servers
	commandServer := command_server.NewCommandServer(DBS, hm)
	// task server need commandServer
	taskServer := task_server.NewTaskServer[api.CommandServer](DBS, commandServer, hm)
	timeServer := time_server.NewTimeServer(DBS)
	// schedule server need task server and time server
	scheduleServer := schedule_server.NewScheduleServer[api.TaskServer, api.TimeServer](DBS, taskServer, timeServer)

	// start schedule server
	go func() {
		scheduleServer.Start(ctx,
			ServerConfig.Interval*time.Second,
			ServerConfig.CleanTime*time.Hour)
	}()

	// start gRPC server
	go grpc_task_template.StartGRPCServer(ctx, mainLog, Config, taskServer, DBS)

	var sb strings.Builder
	sb.WriteString(":")
	sb.WriteString(ServerConfig.Port)
	//addr := sb.String()
	apiServer := fiber.New(
		fiber.Config{
			ReadTimeout:  ServerConfig.ReadTimeout * time.Minute,
			WriteTimeout: ServerConfig.WriteTimeout * time.Minute,
			AppName:      "github.com/littlebluewhite/schedule_task_command",
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		},
	)

	group.Inject(apiServer, DBS, scheduleServer, hm)

	// for api server shout down gracefully
	serverShutdown := make(chan struct{})
	go func() {
		_ = <-ctx.Done()
		mainLog.Infoln("Gracefully shutting down api server")
		_ = apiServer.Shutdown()
		serverShutdown <- struct{}{}
	}()

	if err := apiServer.Listen(":5487"); err != nil {
		mainLog.Errorf("listen: %s\n", err)
	}

	// Listen for the interrupt signal.
	<-serverShutdown

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	DBS.Close()
	scheduleServer.Close()
	mainLog.Infoln("Server exiting")
}
