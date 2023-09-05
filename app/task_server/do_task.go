package task_server

import (
	"context"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/entry/e_task_template"
	"sort"
	"time"
)

func (t *TaskServer) doTask(task e_task.Task) e_task.Task {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	task.Status.TStatus = e_task.Process
	task.CancelFunc = cancel
	// write task
	t.writeTask(task)

	stages := task.Template.Stages
	gsr := getStages(stages)

	for _, sn := range gsr.sns {
		task.Status.Stages = int(sn)
		// write task
		t.writeTask(task)

		s := gsr.stageMap[sn]
		task = t.doStages(ctx, s, task)
		if task.Status.FailedCommandId != "" {
			break
		}
	}
	// no wrong, is success
	if task.Status.FailedCommandId == "" {
		task.Status.TStatus = e_task.Success
	}

	now := time.Now()
	task.To = &now

	// write task
	t.writeTask(task)

	//send to redis channel
	if e := t.rdbPub(task); e != nil {
		panic(e)
	}
	return task
}

// getStages return stage number array without duplicates and return the map (stage number as key stages as value)
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
		case e_task_template.Mode(0).String():
			monitor = append(monitor, stages[i])
		case e_task_template.Mode(1).String():
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

func (t *TaskServer) doStages(ctx context.Context, sv stageMapValue, task e_task.Task) e_task.Task {
	comNumber := len(sv.monitor) + len(sv.execute)
	ch := make(chan e_command.Command, comNumber)
	defer close(ch)

	triggerFrom := append(task.TriggerFrom, "task")
	for _, stage := range sv.monitor {
		go func(stage e_task_template.TaskStage) {
			com, _ := t.cs.Execute(
				ctx, int(*stage.CommandTemplateID), triggerFrom, task.TriggerAccount, task.Token)
			ch <- com
		}(stage)
	}
	// wait 500 milliseconds to execute "execute command"
	time.Sleep(500 * time.Millisecond)
	for _, stage := range sv.execute {
		go func(stage e_task_template.TaskStage) {
			com, _ := t.cs.Execute(
				ctx, int(*stage.CommandTemplateID), triggerFrom, task.TriggerAccount, task.Token)
			ch <- com
		}(stage)
	}
	for i := 0; i < comNumber; i++ {
		select {
		case com := <-ch:
			if com.Status != e_command.Success {
				task.Status.TStatus = e_task.TStatus(com.Status)
				task.Status.FailedCommandId = com.CommandId
				task.Status.FailedCommandTemplateId = com.TemplateId
				task.Status.FailedMessage = com.Message
				return task
			}
		}
	}
	return task
}

func (t *TaskServer) writeTask(task e_task.Task) {
	t.chs.mu.Lock()
	defer t.chs.mu.Unlock()
	t.t[task.TaskId] = task
}
