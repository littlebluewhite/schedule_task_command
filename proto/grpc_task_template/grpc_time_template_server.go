package grpc_task_template

import (
	"context"
	"fmt"
	"github.com/littlebluewhite/schedule_task_command/entry/e_time"
)

type TimeServer interface {
	Execute(pt e_time.PublishTime) (isTime bool, err error)
}

type TimeTemplateService struct {
	UnimplementedTimeTemplateServiceServer
	ts TimeServer
	o  *Operate
}

func NewTimeTemplateService(ts TimeServer, o *Operate) *TimeTemplateService {
	return &TimeTemplateService{
		ts: ts,
		o:  o,
	}
}

func (tts *TimeTemplateService) SendTimeTemplate(ctx context.Context, req *SendTimeTemplateRequest) (*SendTimeTemplateResponse, error) {
	fmt.Println("get send time template request")
	isTime, err := tts.ts.Execute(tts.o.generatePublishTime(req))
	if err != nil {
		return nil, err
	}
	return &SendTimeTemplateResponse{
		IsTime: isTime,
	}, nil
}
