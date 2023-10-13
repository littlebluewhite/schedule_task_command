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
	tl := make([]e_command.Command, 0, len(commandIds))
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

func (o *Operate) Cancel(commandId, message string) error {
	if err := o.commandS.CancelCommand(commandId, message); err != nil {
		return err
	}
	return nil
}

func (o *Operate) GetHistory(templateId, start, stop, status string) ([]e_command.CommandPub, error) {
	s := e_command.S2Status(&status)
	if s != e_command.Success && s != e_command.Failure && s != e_command.Cancel && status != "" {
		return nil, HistoryStatusErr
	}
	ht, e := o.commandS.ReadFromHistory(templateId, start, stop, status)
	if e != nil {
		return nil, e
	}
	return ht, nil
}
