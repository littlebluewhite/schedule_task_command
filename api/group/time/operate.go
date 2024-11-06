package time

import (
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/entry/e_time"
)

type Operate struct {
	timeS api.TimeServer
}

func NewOperate(timeS api.TimeServer) *Operate {
	o := &Operate{
		timeS: timeS,
	}
	return o
}

func (o *Operate) GetHistory(id, templateId, start, stop, isTime string) ([]e_time.PublishTime, error) {
	pt, e := o.timeS.ReadFromHistory(id, templateId, start, stop, isTime)
	if e != nil {
		return nil, e
	}
	return pt, nil
}
