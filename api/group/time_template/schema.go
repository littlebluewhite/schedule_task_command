package time_template

import (
	"schedule_task_command/util"
)

type CheckTime struct {
	TriggerFrom    []string `json:"trigger_from"`
	TriggerAccount string   `json:"trigger_account"`
	Token          string   `json:"token"`
}

var NoStartTime = util.MyErr("No start time input")

var CannotFindTemplate = util.MyErr("can not find time template")
