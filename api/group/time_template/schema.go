package time_template

import (
	"github.com/littlebluewhite/schedule_task_command/util"
	"time"
)

type SendTime struct {
	TemplateId     int                 `json:"template_id"`
	TriggerFrom    []map[string]string `json:"trigger_from"`
	TriggerAccount string              `json:"trigger_account"`
	Token          string              `json:"token"`
	Time           *time.Time          `json:"time,omitempty"`
}

type CheckTime struct {
	TriggerFrom    []map[string]string `json:"trigger_from"`
	TriggerAccount string              `json:"trigger_account"`
	Token          string              `json:"token"`
	Time           *time.Time          `json:"time,omitempty"`
}

var CannotFindTemplate = util.MyErr("can not find time template")
