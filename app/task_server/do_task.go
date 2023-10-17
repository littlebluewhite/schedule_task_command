package task_server

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"schedule_task_command/dal/model"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/entry/e_command_template"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/entry/e_task_template"
	"schedule_task_command/util"
	"sort"
	"time"
)

func (t *TaskServer[T]) doTask(ctx context.Context, task e_task.Task) e_task.Task {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	task.Status.TStatus = e_task.Process
	task.CancelFunc = cancel
	// write task
	t.writeTask(task)

	stages := task.TaskData.Stages
	gsr := getStages(stages)

	for _, sn := range gsr.sns {
		task.Status.Stages = int(sn)
		// write task
		t.writeTask(task)

		s := gsr.stageMap[sn]
		task = t.doOneStage(ctx, s, task)
		if task.Status.FailedMessage != nil {
			e := util.MyErr(fmt.Sprintf("task id: %d failed at stage %d\n", task.ID, sn))
			task.Message = &e
			// cancel task
			break
		}
	}
	// no wrong, is success
	if task.Message == nil {
		task.Status.TStatus = e_task.Success
	}

	now := time.Now()
	task.To = &now

	// write client message
	task.ClientMessage = t.ReadMap()[task.ID].ClientMessage

	// write task
	t.writeTask(task)

	// write to history in influxdb
	t.writeToHistory(task)

	//send to redis channel
	if e := t.rdbPub(task); e != nil {
		panic(e)
	}
	return task
}

// getStages return stage number array without duplicates and return the map (monitor and ExecuteReturnId commands slice)
func getStages(stages []e_task_template.TaskStage) (gsr getStagesResult) {
	snSet := make(map[int32]struct{})
	gsr.stageMap = make(map[int32]stageMapValue)
	for i := 0; i < len(stages); i++ {
		sn := stages[i].StageNumber
		if _, ok := snSet[sn]; !ok {
			gsr.sns = append(gsr.sns, sn)
			snSet[sn] = struct{}{}
		}
		monitor := gsr.stageMap[sn].monitor
		execute := gsr.stageMap[sn].execute
		switch stages[i].Mode {
		case e_task_template.Monitor:
			monitor = append(monitor, stages[i])
		case e_task_template.Execute:
			execute = append(execute, stages[i])
		default:
		}
		gsr.stageMap[sn] = stageMapValue{monitor: monitor, execute: execute}
	}
	sort.Slice(gsr.sns, func(i, j int) bool {
		return gsr.sns[i] < gsr.sns[j]
	})
	return
}

func (t *TaskServer[T]) doOneStage(ctx context.Context, sv stageMapValue, task e_task.Task) e_task.Task {
	comNumber := len(sv.monitor) + len(sv.execute)
	ch := make(chan comBuilder, comNumber)

	triggerFrom := append(task.TriggerFrom, fmt.Sprintf("task: %d", task.ID))
	for _, stage := range sv.monitor {
		go func(stage e_task_template.TaskStage) {
			com := t.ts2Com(stage, triggerFrom, task)
			com = t.cs.ExecuteWait(ctx, com)
			ch <- comBuilder{mode: e_task_template.Monitor, name: stage.Name, com: com,
				tags: stage.Tags}
		}(stage)
	}
	// wait 500 milliseconds to Execute executed command
	time.Sleep(500 * time.Millisecond)
	for _, stage := range sv.execute {
		go func(stage e_task_template.TaskStage) {
			com := t.ts2Com(stage, triggerFrom, task)
			com = t.cs.ExecuteWait(ctx, com)
			ch <- comBuilder{mode: e_task_template.Execute, name: stage.Name, com: com,
				tags: stage.Tags}
		}(stage)
	}
	mts := make([]e_task.TaskStage, 0, len(sv.monitor))
	ets := make([]e_task.TaskStage, 0, len(sv.execute))
Loop:
	for i := 0; i < comNumber; i++ {
		select {
		case comB := <-ch:
			com := comB.com
			ts := e_task.TaskStage{
				Name:       comB.name,
				ID:         com.ID,
				From:       com.From,
				To:         com.To,
				Status:     com.Status,
				Message:    com.Message,
				Tags:       comB.tags,
				Variables:  com.Variables,
				TemplateId: com.TemplateId,
			}
			switch comB.mode {
			case e_task_template.Monitor:
				mts = append(mts, ts)
			case e_task_template.Execute:
				ets = append(ets, ts)
			}
			if com.Status != e_command.Success {
				task.Status.TStatus = e_task.TStatus(com.Status)
				task.Status.FailedCommandId = com.ID
				task.Status.FailedCommandTemplateId = com.TemplateId
				task.Status.FailedMessage = com.Message
				break Loop
			}
		}
	}
	task.Stages[task.Status.Stages] = e_task.TaskStageC{Execute: ets, Monitor: mts}
	return task
}

func (t *TaskServer[T]) ts2Com(stage e_task_template.TaskStage, triggerFrom []string,
	task e_task.Task) (c e_command.Command) {
	// get variables
	var variables map[string]string
	if v, ok := task.Variables[int(stage.ID)]; ok {
		variables = v
	} else {
		_ = json.Unmarshal(stage.Variable, &variables)
	}
	c = e_command.Command{
		TemplateId:     int(stage.CommandTemplateID),
		TriggerFrom:    triggerFrom,
		TriggerAccount: task.TriggerAccount,
		Token:          task.Token,
		Variables:      variables,
	}
	// use command template id first
	if stage.CommandTemplateID != 0 {
		var cacheMap map[int]model.CommandTemplate
		if x, found := t.dbs.GetCache().Get("commandTemplates"); found {
			cacheMap = x.(map[int]model.CommandTemplate)
			c.CommandData = e_command_template.M2Entry(cacheMap[int(stage.CommandTemplateID)])
		} else {
			t.l.Info().Printf("Cannot find command template id %v, so use template to execute command",
				stage.CommandTemplateID)
			c.CommandData = stage.CommandTemplate
		}
	} else {
		c.CommandData = stage.CommandTemplate
	}

	return
}
