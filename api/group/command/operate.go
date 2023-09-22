package command

import (
	"errors"
	"fmt"
	"schedule_task_command/api"
	"schedule_task_command/entry/e_command"
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

func (o *Operate) Find(commandIds []string) ([]e_command.Command, error) {
	tm := o.commandS.ReadMap()
	tl := make([]e_command.Command, len(commandIds))
	for _, commandId := range commandIds {
		t, ok := tm[commandId]
		if !ok {
			return nil, errors.New(fmt.Sprintf("cannot find command id: %s", commandId))
		} else {
			tl = append(tl, t)
		}
	}
	return tl, nil
}

func (o *Operate) Cancel(commandId string) error {
	tm := o.commandS.ReadMap()
	command, ok := tm[commandId]
	if !ok {
		return errors.New(fmt.Sprintf("cannot find command id: %s", commandId))
	}
	command.CancelFunc()
	return nil
}

func (o *Operate) GetHistory(templateId, status, start, stop string) ([]e_command.Command, error) {
	s := e_command.S2Status(&status)
	if s != e_command.Success && s != e_command.Failure && s != e_command.Cancel {
		return nil, HistoryStatusErr
	}
	ht, e := o.commandS.ReadFromHistory(templateId, status, start, stop)
	if e != nil {
		return nil, e
	}
	return ht, nil
}
