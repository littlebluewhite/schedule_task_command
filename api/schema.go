package api

import (
	"context"
	"schedule_task_command/entry/e_time_data"
	"time"
)

type TaskServer interface {
	Execute(ctx context.Context, templateId int, triggerFrom []string,
		triggerAccount string, token string) (taskId string)
}

type TimeServer interface {
	Execute(templateId int, triggerFrom []string,
		triggerAccount string, token string) (bool, error)
	ReadFromHistory(templateId, start, stop string) ([]e_time_data.PublishTime, error)
}

type ScheduleSer interface {
	Start(ctx context.Context, interval, removeTime time.Duration)
	GetTimeServer() TimeServer
	GetTaskServer() TaskServer
}
