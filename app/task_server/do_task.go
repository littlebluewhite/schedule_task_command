package task_server

import (
	"context"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/entry/e_task_template"
	"sort"
)

func (t *TaskServer) doTask(task e_task.Task) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stages := task.Template.Stages
	gsr := getStages(stages)
	for _, sn := range gsr.sns {
		s := gsr.stageMap[sn]
		t.doStages(s)
	}
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

func (t *TaskServer) doStages(sv stageMapValue) {
	for _, stage := range sv.monitor {

	}
}
