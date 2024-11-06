package task_server

import (
	"github.com/littlebluewhite/schedule_task_command/entry/e_command"
	"github.com/littlebluewhite/schedule_task_command/entry/e_module"
	"github.com/littlebluewhite/schedule_task_command/entry/e_task"
	"github.com/littlebluewhite/schedule_task_command/entry/e_task_template"
	"github.com/littlebluewhite/schedule_task_command/util"
	"sync"
)

type chs struct {
	mu *sync.RWMutex
	wg *sync.WaitGroup
}

type hubManager interface {
	Broadcast(module e_module.Module, message []byte)
}

type getStagesResult struct {
	sns      []int32
	stageMap map[int32]stageMap
}

type stageMap struct {
	monitor map[int32]stageItem
	execute map[int32]stageItem
}

type stageItem struct {
	templateStageItem e_task_template.StageItem
	taskStageItem     e_task.StageItem
}

type comBuilder struct {
	stageID int32
	parser  []e_task_template.ParserItem
	com     e_command.Command
}

type StreamCancel struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

var SendToRedisErr = util.MyErr("send task to redis cannot format")
var TaskCannotCancel = util.MyErr("task cannot be canceled")
