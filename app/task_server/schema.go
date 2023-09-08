package task_server

import (
	"schedule_task_command/entry/e_task"
	"schedule_task_command/entry/e_task_template"
	"schedule_task_command/util"
	"sync"
)

type chs struct {
	rec chan e_task.Task
	mu  *sync.RWMutex
}

type getStagesResult struct {
	sns      []int32
	stageMap map[int32]stageMapValue
}

type stageMapValue struct {
	monitor []e_task_template.TaskStage
	execute []e_task_template.TaskStage
}

type SendTask struct {
	TemplateId     int      `json:"template_id"`
	TriggerFrom    []string `json:"trigger_from"`
	TriggerAccount string   `json:"trigger_account"`
	Token          string   `json:"token"`
}

var CannotFindTemplate = util.MyErr("can not find task template")
var TaskCanceled = util.MyErr("Task has been canceled")
