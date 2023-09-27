package api

import (
	"context"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/entry/e_time"
	"time"
)

type TaskServer interface {
	ReadMap() map[string]e_task.Task
	GetList() []e_task.Task
	ExecuteReturnId(ctx context.Context, task e_task.Task) (taskId string, err error)
	ReadFromHistory(taskTemplateId, status, start, stop string) ([]e_task.TaskPub, error)
	GetCommandServer() CommandServer
}

type TimeServer interface {
	Execute(pt e_time.PublishTime) (bool, error)
	ReadFromHistory(templateId, start, stop string) ([]e_time.PublishTime, error)
}

type CommandServer interface {
	ReadMap() map[string]e_command.Command
	GetList() []e_command.Command
	ExecuteReturnId(ctx context.Context, command e_command.Command) (commandId string, err error)
	ReadFromHistory(commandTemplateId, status, start, stop string) ([]e_command.CommandPub, error)
}

type ScheduleSer interface {
	Start(ctx context.Context, interval, removeTime time.Duration)
	GetTimeServer() TimeServer
	GetTaskServer() TaskServer
}
