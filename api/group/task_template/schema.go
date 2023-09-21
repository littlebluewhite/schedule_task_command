package task_template

import "schedule_task_command/util"

type SendTask struct {
	TriggerFrom    []string `json:"trigger_from" example:"[task execute]"`
	TriggerAccount string   `json:"trigger_account" example:"Wilson"`
	Token          string   `json:"token"`
}

var CannotFindTemplate = util.MyErr("can not find task template")
