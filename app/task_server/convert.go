package task_server

import (
	"github.com/littlebluewhite/schedule_task_command/entry/e_command"
	"github.com/littlebluewhite/schedule_task_command/entry/e_task"
)

func stageMap2taskStage(stages map[int32]stageMap) map[int32]e_task.Stage {
	taskStages := make(map[int32]e_task.Stage)
	for key, sm := range stages {
		monitor := make([]e_task.StageItem, 0, len(sm.monitor))
		for _, item := range sm.monitor {
			monitor = append(monitor, item.taskStageItem)
		}
		execute := make([]e_task.StageItem, 0, len(sm.execute))
		for _, item := range sm.execute {
			execute = append(execute, item.taskStageItem)
		}
		taskStages[key] = e_task.Stage{
			Monitor: monitor,
			Execute: execute,
		}
	}
	return taskStages
}

func com2stageItem(si stageItem, com e_command.Command) stageItem {
	si.taskStageItem.CommandID = com.ID
	si.taskStageItem.Status = com.Status
	si.taskStageItem.From = &com.From
	si.taskStageItem.To = com.To
	si.taskStageItem.Message = com.Message
	si.taskStageItem.Variables = com.Variables
	return si
}
