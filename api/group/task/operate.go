package task

import (
	"errors"
	"fmt"
	"schedule_task_command/api"
	"schedule_task_command/entry/e_task"
	"strconv"
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

func (o *Operate) Find(ids []uint64) ([]e_task.Task, error) {
	tl := make([]e_task.Task, 0, len(ids))
	for _, id := range ids {
		t, err := o.taskS.ReadOne(id)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("cannot find id: %d", id))
		} else {
			tl = append(tl, t)
		}
	}
	return tl, nil
}

func (o *Operate) Cancel(id uint64, message string) error {
	if err := o.taskS.CancelTask(id, message); err != nil {
		return err
	}
	return nil
}

func (o *Operate) GetHistory(id, templateId, start, stop, status string) ([]e_task.TaskPub, error) {
	s := e_task.S2Status(&status)
	if s != e_task.Success && s != e_task.Failure && s != e_task.Cancel && status != "" {
		return nil, HistoryStatusErr
	}
	ht, e := o.taskS.ReadFromHistory(id, templateId, start, stop, status)
	if e != nil {
		return nil, e
	}
	return ht, nil
}

func (o *Operate) FindById(id uint64) (t e_task.TaskPub, err error) {
	task, err := o.taskS.ReadOne(id)
	if err == nil {
		t = e_task.ToPub(task)
		return
	}
	ht, err := o.GetHistory(strconv.FormatUint(id, 10), "", "0", "", "")
	if len(ht) > 0 {
		t = ht[0]
		return
	}
	err = errors.New(fmt.Sprintf("cannot find id: %d", id))
	return
}
