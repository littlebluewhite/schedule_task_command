package api

import (
	"context"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/entry/e_time"
	"time"
)

type TaskServer interface {
	ReadMap() map[string]e_task.Task
	GetList() []e_task.Task
	ExecuteReturnId(ctx context.Context, task e_task.Task) (taskId string, err error)
}

type TimeServer interface {
	Execute(pt e_time.PublishTime) (bool, error)
	ReadFromHistory(templateId, start, stop string) ([]e_time.PublishTime, error)
}

type ScheduleSer interface {
	Start(ctx context.Context, interval, removeTime time.Duration)
	GetTimeServer() TimeServer
	GetTaskServer() TaskServer
}
