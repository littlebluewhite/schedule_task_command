package task_server

import (
	"github.com/goccy/go-json"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/entry/e_task_template"
	"schedule_task_command/util"
	"sync"
)

type chs struct {
	mu *sync.RWMutex
}

type getStagesResult struct {
	sns      []int32
	stageMap map[int32]stageMapValue
}

type stageMapValue struct {
	monitor []e_task_template.TaskStage
	execute []e_task_template.TaskStage
}

type comBuilder struct {
	mode e_task_template.Mode
	name string
	com  e_command.Command
	tags json.RawMessage
}

var TaskCanceled = util.MyErr("Task has been canceled")
var SendToRedisErr = util.MyErr("send task to redis cannot format")
var TaskNotFind = util.MyErr("can not find task")
var TaskCannotCancel = util.MyErr("task cannot be canceled")
var TaskTemplateVariable = util.MyErr("task template variable failed to format")
