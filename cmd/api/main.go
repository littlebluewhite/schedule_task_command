package main

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"os/signal"
	"path/filepath"
	"runtime"
	"schedule_task_command/api"
	"schedule_task_command/api/group"
	"schedule_task_command/app/command_server"
	"schedule_task_command/app/dbs"
	"schedule_task_command/app/schedule_server"
	"schedule_task_command/app/task_server"
	"schedule_task_command/app/time_server"
	"schedule_task_command/app/websocket_hub"
	_ "schedule_task_command/docs"
	"schedule_task_command/util/config"
	"schedule_task_command/util/logFile"
	"strings"
	"syscall"
	"time"
)

var (
	mainLog  logFile.LogFile
	rootPath string
)

// 初始化配置
func init() {
	// log配置
	mainLog = logFile.NewLogFile("", "main.log")
	_, b, _, _ := runtime.Caller(0)
	rootPath = filepath.Dir(filepath.Dir(filepath.Dir(b)))
}

// @title           Schedule-Task-Command swagger API
// @version         2.13.12
// @description     This is a schedule-command server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Wilson
// @contact.url    https://github.com/littlebluewhite
// @contact.email  wwilson008@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      127.0.0.1:5487

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mainLog.Info().Println("command module start")

	// read config
	Config := config.NewConfig[config.Config](rootPath, "config", "config", config.Yaml)

	ServerConfig := Config.Server

	// DBs start includes SQL Cache
	DBS := dbs.NewDbs(mainLog, false, Config)
	defer func() {
		DBS.GetIdb().Close()
		mainLog.Info().Println("influxDB Disconnect")
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

	var sb strings.Builder
	sb.WriteString(":")
	sb.WriteString(ServerConfig.Port)
	//addr := sb.String()
	apiServer := fiber.New(
		fiber.Config{
			ReadTimeout:  ServerConfig.ReadTimeout * time.Minute,
			WriteTimeout: ServerConfig.WriteTimeout * time.Minute,
			AppName:      "schedule_task_command",
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		},
	)

	group.Inject(apiServer, DBS, scheduleServer, hm)

	// for api server shout down gracefully
	serverShutdown := make(chan struct{})
	go func() {
		_ = <-ctx.Done()
		mainLog.Info().Println("Gracefully shutting down api server")
		_ = apiServer.Shutdown()
		serverShutdown <- struct{}{}
	}()

	if err := apiServer.Listen(":5487"); err != nil {
		mainLog.Error().Fatalf("listen: %s\n", err)
	}

	// Listen for the interrupt signal.
	<-serverShutdown

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	time.Sleep(1 * time.Second)
	mainLog.Info().Println("Server exiting")
}
