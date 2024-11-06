package schedule_server

import (
	"context"
	"github.com/littlebluewhite/schedule_task_command/entry/e_task"
	"github.com/littlebluewhite/schedule_task_command/entry/e_time"
	"time"
)

type taskServer interface {
	Start(ctx context.Context, removeTime time.Duration)
	ReadMap() map[uint64]e_task.Task
	GetList() []e_task.Task
	ExecuteWait(ctx context.Context, task e_task.Task) e_task.Task
	ExecuteReturnId(ctx context.Context, task e_task.Task) (id uint64, err error)
	ReadFromHistory(id, taskTemplateId, start, stop, status string) ([]e_task.TaskPub, error)
	Close()
}

type timeServer interface {
	Start(ctx context.Context)
	Execute(pt e_time.PublishTime) (bool, error)
	ReadFromHistory(id, templateId, start, stop, isTime string) ([]e_time.PublishTime, error)
}
