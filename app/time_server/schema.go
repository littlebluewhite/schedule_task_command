package time_server

import "schedule_task_command/util"

type SendTime struct {
	TemplateId     int      `json:"template_id"`
	TriggerFrom    []string `json:"trigger_from"`
	TriggerAccount string   `json:"trigger_account"`
	Token          string   `json:"token"`
}

var CannotFindTemplate = util.MyErr("can not find time template")
