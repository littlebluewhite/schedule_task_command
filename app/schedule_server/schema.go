package schedule_server

import (
	"context"
	"schedule_task_command/entry/e_time_data"
	"time"
)

type taskServer interface {
	Start(ctx context.Context, removeTime time.Duration)
	Execute(ctx context.Context, templateId int, triggerFrom []string,
		triggerAccount string, token string) (taskId string)
}

type timeServer interface {
	Start(ctx context.Context)
	ReadFromHistory(templateId, start, stop string) ([]e_time_data.PublishTime, error)
}
