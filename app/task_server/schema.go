package task_server

import (
	"errors"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/entry/e_task_template"
	"sync"
)

type chs struct {
	rec chan e_task.Task
	mu  *sync.RWMutex
}

type executeParams struct {
	templateId     int
	triggerFrom    []string
	triggerAccount string
	token          string
}

type getStagesResult struct {
	sns      []int32
	stageMap map[int32]stageMapValue
}

type stageMapValue struct {
	monitor []e_task_template.TaskStage
	execute []e_task_template.TaskStage
}

var cannotFindTemplate = errors.New("can not find task template")
