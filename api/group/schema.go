package group

import (
	"context"
	"time"
)

type TaskServer interface {
	Execute(ctx context.Context, templateId int, triggerFrom []string,
		triggerAccount string, token string) (taskId string)
}

type TimeServer interface {
	Execute(templateId int, triggerFrom []string,
		triggerAccount string, token string) bool
}

type ScheduleSer interface {
	Start(ctx context.Context, interval, removeTime time.Duration)
	GetTimeServer() TimeServer
	GetTaskServer() TaskServer
}
