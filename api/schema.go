package api

import (
	"context"
	"github.com/gofiber/contrib/websocket"
	api2 "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"schedule_task_command/app/dbs/influxdb"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/entry/e_module"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/entry/e_time"
	"time"
)

type Dbs interface {
	GetSql() *gorm.DB
	GetCache() *cache.Cache
	GetRdb() redis.UniversalClient
	GetIdb() *influxdb.Influx
	Close()
}

type HistoryDB interface {
	Close()
	Writer() api2.WriteAPI
	Querier() api2.QueryAPI
}

type TaskServer interface {
	ReadMap() map[uint64]e_task.Task
	ReadOne(id uint64) (e_task.Task, error)
	DeleteOne(id uint64)
	GetList() []e_task.Task
	ExecuteReturnId(ctx context.Context, task e_task.Task) (id uint64, err error)
	ReadFromHistory(id, taskTemplateId, start, stop, status string) ([]e_task.TaskPub, error)
	GetCommandServer() CommandServer
	CancelTask(id uint64, message string) error
}

type TimeServer interface {
	Execute(pt e_time.PublishTime) (bool, error)
	ReadFromHistory(id, templateId, start, stop, isTime string) ([]e_time.PublishTime, error)
}

type CommandServer interface {
	ReadMap() map[uint64]e_command.Command
	ReadOne(id uint64) (e_command.Command, error)
	DeleteOne(id uint64)
	GetList() []e_command.Command
	ExecuteReturnId(ctx context.Context, command e_command.Command) (id uint64, err error)
	ReadFromHistory(id, commandTemplateId, start, stop, status string) ([]e_command.CommandPub, error)
	CancelCommand(id uint64, message string) error
	TestExecute(ctx context.Context, com e_command.Command) e_command.Command
}

type ScheduleSer interface {
	Start(ctx context.Context, interval, removeTime time.Duration)
	GetTimeServer() TimeServer
	GetTaskServer() TaskServer
}

type HubManager interface {
	RegisterHub(module e_module.Module)
	Broadcast(module e_module.Module, message []byte)
	WsConnect(module e_module.Module, conn *websocket.Conn) error
}

type Logger interface {
	Infoln(args ...interface{})
	Infof(s string, args ...interface{})
	Errorln(args ...interface{})
	Errorf(s string, args ...interface{})
	Warnln(args ...interface{})
	Warnf(s string, args ...interface{})
}
