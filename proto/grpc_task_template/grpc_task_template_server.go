package grpc_task_template

import (
	"context"
	"fmt"
	"github.com/littlebluewhite/schedule_task_command/entry/e_task"
)

type TaskServer interface {
	ExecuteReturnId(ctx context.Context, task e_task.Task) (id uint64, err error)
}

type TaskTemplateService struct {
	UnimplementedTaskTemplateServiceServer
	ts TaskServer
	o  *Operate
}

func NewTaskTemplateService(ts TaskServer, o *Operate) *TaskTemplateService {
	return &TaskTemplateService{
		ts: ts,
		o:  o,
	}
}

func (tts *TaskTemplateService) SendTaskTemplate(ctx context.Context, req *SendTaskTemplateRequest) (*SendTaskTemplateResponse, error) {
	fmt.Println("get send task template request")
	id, err := tts.ts.ExecuteReturnId(ctx, tts.o.GenerateTask(req))
	if err != nil {
		return nil, err
	}
	return &SendTaskTemplateResponse{
		TaskId: id,
	}, nil
}
