package task_server

import (
	"context"
	"fmt"
	"github.com/littlebluewhite/schedule_task_command/dal/model"
	"github.com/littlebluewhite/schedule_task_command/entry/e_command"
	"github.com/littlebluewhite/schedule_task_command/entry/e_command_template"
	"github.com/littlebluewhite/schedule_task_command/entry/e_task"
	"github.com/littlebluewhite/schedule_task_command/entry/e_task_template"
	"github.com/littlebluewhite/schedule_task_command/util"
	"sort"
	"time"
)

func (t *TaskServer[T]) doTask(ctx context.Context, task e_task.Task) e_task.Task {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	task.Status = e_task.Process
	task.CancelFunc = cancel
	// write task
	t.writeTask(task)

	stages := task.TaskData.StageItems
	gsr := getTaskStage(stages)
	task.Stages = stageMap2taskStage(gsr.stageMap)
	// write task
	t.writeTask(task)
	t.publishContainer(context.Background(), task)

	for _, sn := range gsr.sns {
		task.StageNumber = sn
		// write task
		t.writeTask(task)
		// publish
		t.publishContainer(context.Background(), task)

		s := gsr.stageMap
		task = t.doOneStage(ctx, s, sn, task)
		if len(task.FailedCommands) != 0 {
			e := util.MyErr(fmt.Sprintf("task id: %d failed at stage %d\n", task.ID, sn))
			task.Message = &e
			// cancel task
			break
		}
	}
	// no wrong, is success
	if task.Message == nil {
		task.Status = e_task.Success
	}

	now := time.Now()
	task.To = &now

	// write task
	t.writeTask(task)
	// write to history in influxdb

	go func() {
		t.writeToHistory(task)
	}()

	go func() {
		time.Sleep(300 * time.Millisecond)
		// publish
		t.publishContainer(context.Background(), task)
	}()

	return task
}

func getTaskStage(stages []e_task_template.StageItem) (gsr getStagesResult) {
	stageNumbers := make(map[int32]struct{})
	gsr.stageMap = make(map[int32]stageMap)
	for i := 0; i < len(stages); i++ {
		sn := stages[i].StageNumber
		if _, ok := stageNumbers[sn]; !ok {
			gsr.sns = append(gsr.sns, sn)
			gsr.stageMap[sn] = stageMap{
				monitor: make(map[int32]stageItem),
				execute: make(map[int32]stageItem),
			}
			stageNumbers[sn] = struct{}{}
		}
		switch stages[i].Mode {
		case e_task_template.Monitor:
			gsr.stageMap[sn].monitor[stages[i].ID] = stageItem{
				templateStageItem: stages[i],
				taskStageItem: e_task.StageItem{
					Name:              stages[i].Name,
					StageID:           stages[i].ID,
					StageNumber:       stages[i].StageNumber,
					Mode:              stages[i].Mode,
					Status:            e_command.Prepared,
					Tags:              stages[i].Tags,
					CommandTemplateId: stages[i].CommandTemplateID,
				},
			}
		case e_task_template.Execute:
			gsr.stageMap[sn].execute[stages[i].ID] = stageItem{
				templateStageItem: stages[i],
				taskStageItem: e_task.StageItem{
					Name:              stages[i].Name,
					StageID:           stages[i].ID,
					StageNumber:       stages[i].StageNumber,
					Mode:              stages[i].Mode,
					Status:            e_command.Prepared,
					Tags:              stages[i].Tags,
					CommandTemplateId: stages[i].CommandTemplateID,
				},
			}
		case e_task_template.NoneMode:
		}
	}
	sort.Slice(gsr.sns, func(i, j int) bool {
		return gsr.sns[i] < gsr.sns[j]
	})
	return
}

func (t *TaskServer[T]) doOneStage(ctx context.Context, s map[int32]stageMap, stageNumber int32, task e_task.Task) e_task.Task {
	sm := s[stageNumber]

	// change stage status -> Process
	for key, si := range sm.monitor {
		si.taskStageItem.Status = e_command.Process
		s[stageNumber].monitor[key] = stageItem{
			templateStageItem: si.templateStageItem,
			taskStageItem:     si.taskStageItem,
		}
	}
	for key, si := range sm.execute {
		si.taskStageItem.Status = e_command.Process
		s[stageNumber].execute[key] = stageItem{
			templateStageItem: si.templateStageItem,
			taskStageItem:     si.taskStageItem,
		}
	}
	task.Stages = stageMap2taskStage(s)
	// write task
	t.writeTask(task)
	// publish
	t.publishContainer(context.Background(), task)

	comNumber := len(sm.monitor) + len(sm.execute)
	monitorCh := make(chan comBuilder, len(sm.monitor))
	defer close(monitorCh)
	executeCh := make(chan comBuilder, len(sm.execute))
	defer close(executeCh)

	triggerFrom := append(task.TriggerFrom, map[string]string{"task": fmt.Sprintf("%d", task.ID)})
	for _, stage := range sm.monitor {
		go func(stage e_task_template.StageItem) {
			com := t.ts2Com(stage, triggerFrom, task)
			com = t.cs.ExecuteWait(ctx, com)
			monitorCh <- comBuilder{com: com, stageID: stage.ID, parser: stage.Parser}
		}(stage.templateStageItem)
	}
	// wait 500 milliseconds to Execute executed command
	time.Sleep(500 * time.Millisecond)
	for _, stage := range sm.execute {
		go func(stage e_task_template.StageItem) {
			com := t.ts2Com(stage, triggerFrom, task)
			com = t.cs.ExecuteWait(ctx, com)
			executeCh <- comBuilder{com: com, stageID: stage.ID, parser: stage.Parser}
		}(stage.templateStageItem)
	}
	for i := 0; i < comNumber; i++ {
		var comB comBuilder
		select {
		case comB = <-monitorCh:
			s[stageNumber].monitor[comB.stageID] = com2stageItem(s[stageNumber].monitor[comB.stageID], comB.com)
		case comB = <-executeCh:
			s[stageNumber].execute[comB.stageID] = com2stageItem(s[stageNumber].execute[comB.stageID], comB.com)
		}
		task.Stages = stageMap2taskStage(s)
		// command return parser to variables
		task.Variables = commandReturn2Variables(task.Variables, comB)
		// command status is Cancel and task client message is empty
		if comB.com.Status == e_command.Cancel && task.ClientMessage == "" {
			// write client message
			task2, err := t.ReadOne(task.ID)
			if err != nil {
				t.l.Errorln(err)
			}
			task.ClientMessage = task2.ClientMessage
		}
		// write task
		t.writeTask(task)
		// publish
		t.publishContainer(context.Background(), task)
		if comB.com.Status != e_command.Success {
			task.Status = e_task.Status(comB.com.Status)
			task.FailedCommands = append(task.FailedCommands, e_task.FailedCommand{
				CommandID:         comB.com.ID,
				CommandTemplateID: comB.com.TemplateId,
				Message:           comB.com.Message,
				Status:            comB.com.Status,
			})
		}
	}
	return task
}

func (t *TaskServer[T]) ts2Com(stage e_task_template.StageItem, triggerFrom []map[string]string,
	task e_task.Task) (c e_command.Command) {
	// get variables
	variables, err := getStageVariables(stage, task)
	if err != nil {
		t.l.Errorf("task id: %d, task template id: %d, stage id: %d, error: %v",
			task.ID, task.TemplateId, stage.ID, err)
	}
	c = e_command.Command{
		TemplateId:     stage.CommandTemplateID,
		Source:         task.Source,
		TriggerFrom:    triggerFrom,
		TriggerAccount: task.TriggerAccount,
		Token:          task.Token,
		Variables:      variables,
	}
	// use command template as command data
	if stage.CommandTemplate.ID != 0 {
		c.CommandData = stage.CommandTemplate
		// fetch command template as command data
	} else if stage.CommandTemplateID == 0 {
		var cacheMap map[int]model.CommandTemplate
		if x, found := t.dbs.GetCache().Get("commandTemplates"); found {
			cacheMap = x.(map[int]model.CommandTemplate)
		}
		c.CommandData = e_command_template.M2Entry(cacheMap[int(stage.CommandTemplateID)])
	}
	trigger := map[string]string{"command_template": fmt.Sprintf("%d", c.CommandData.ID)}
	c.TriggerFrom = append(c.TriggerFrom, trigger)

	return
}
