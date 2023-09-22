package schedule_server

import (
	"context"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/entry/e_time"
	"time"
)

type taskServer interface {
	Start(ctx context.Context, removeTime time.Duration)
	ReadMap() map[string]e_task.Task
	GetList() []e_task.Task
	ExecuteWait(ctx context.Context, task e_task.Task) e_task.Task
	ExecuteReturnId(ctx context.Context, task e_task.Task) (taskId string, err error)
	ReadFromHistory(taskTemplateId, status, start, stop string) ([]e_task.Task, error)
}

type timeServer interface {
	Start(ctx context.Context)
	Execute(pt e_time.PublishTime) (bool, error)
	ReadFromHistory(templateId, start, stop string) ([]e_time.PublishTime, error)
}
