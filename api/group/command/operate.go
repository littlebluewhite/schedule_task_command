package command

import (
	"errors"
	"fmt"
	"schedule_task_command/api"
	"schedule_task_command/entry/e_command"
	"strconv"
)

type Operate struct {
	commandS api.CommandServer
}

func NewOperate(commandS api.CommandServer) *Operate {
	o := &Operate{
		commandS: commandS,
	}
	return o
}

func (o *Operate) List() ([]e_command.Command, error) {
	tl := o.commandS.GetList()
	return tl, nil
}

func (o *Operate) Find(ids []uint64) ([]e_command.Command, error) {
	tm := o.commandS.ReadMap()
	tl := make([]e_command.Command, 0, len(ids))
	for _, id := range ids {
		t, ok := tm[id]
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot find command id: %d", id))
		} else {
			tl = append(tl, t)
		}
	}
	return tl, nil
}

func (o *Operate) Cancel(id uint64, message string) error {
	if err := o.commandS.CancelCommand(id, message); err != nil {
		return err
	}
	return nil
}

func (o *Operate) GetHistory(id, templateId, start, stop, status string) ([]e_command.CommandPub, error) {
	s := e_command.S2Status(&status)
	if s != e_command.Success && s != e_command.Failure && s != e_command.Cancel && status != "" {
		return nil, HistoryStatusErr
	}
	ht, e := o.commandS.ReadFromHistory(id, templateId, start, stop, status)
	if e != nil {
		return nil, e
	}
	return ht, nil
}

func (o *Operate) FindById(id uint64) (c e_command.CommandPub, err error) {
	cm := o.commandS.ReadMap()
	com, ok := cm[id]
	if ok {
		c = e_command.ToPub(com)
		return
	}
	hc, err := o.GetHistory(strconv.FormatUint(id, 10), "", "0", "", "")
	if len(hc) > 0 {
		c = hc[0]
		return
	}
	err = errors.New(fmt.Sprintf("cannot find id: %d", id))
	return
}
