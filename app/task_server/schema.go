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

var TaskCanceled = util.MyErr("Task has been canceled")
var SendToRedisErr = util.MyErr("send task to redis cannot format")
var TaskNotFind = util.MyErr("can not find task")
var TaskCannotCancel = util.MyErr("task cannot be canceled")
