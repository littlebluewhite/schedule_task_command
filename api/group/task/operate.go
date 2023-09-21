package task

import (
	"errors"
	"fmt"
	"schedule_task_command/api"
	"schedule_task_command/entry/e_task"
)

type Operate struct {
	taskS api.TaskServer
}

func NewOperate(taskS api.TaskServer) *Operate {
	o := &Operate{
		taskS: taskS,
	}
	return o
}

func (o *Operate) List() ([]e_task.Task, error) {
	tl := o.taskS.GetList()
	return tl, nil
}

func (o *Operate) Find(taskIds []string) ([]e_task.Task, error) {
	tm := o.taskS.ReadMap()
	tl := make([]e_task.Task, len(taskIds))
	for _, taskId := range taskIds {
		t, ok := tm[taskId]
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot find task id: %s", taskId))
		} else {
			tl = append(tl, t)
		}
	}
	return tl, nil
}

func (o *Operate) Cancel(taskId string) error {
	tm := o.taskS.ReadMap()
	task, ok := tm[taskId]
	if !ok {
		return errors.New(fmt.Sprintf("cannot find task id: %s", taskId))
	}
	task.CancelFunc()
	return nil
}
