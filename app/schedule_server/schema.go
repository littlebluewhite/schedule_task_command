package schedule_server

import (
	"context"
	"time"
)

type taskServer interface {
	Start(ctx context.Context, removeTime time.Duration)
	Execute(ctx context.Context, templateId int, triggerFrom []string,
		triggerAccount string, token string) (taskId string)
}

type timeServer interface {
	Start(ctx context.Context)
}
