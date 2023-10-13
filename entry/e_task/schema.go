package e_task

import (
	"github.com/goccy/go-json"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/entry/e_task_template"
	"schedule_task_command/util"
	"time"
)

type Task struct {
	TaskId         string                       `json:"task_id"`
	Token          string                       `json:"token"`
	From           time.Time                    `json:"from"`
	To             *time.Time                   `json:"to"`
	Variables      map[string]map[string]string `json:"variables"`
	TriggerFrom    []string                     `json:"trigger_from"`
	TriggerAccount string                       `json:"trigger_account"`
	Status         Status                       `json:"status"`
	Stages         map[int]TaskStageC           `json:"stages"`
	ClientMessage  string                       `json:"client_message"`
	Message        *util.MyErr                  `json:"message"`
	TemplateId     int                          `json:"template_id"`
	TaskData       e_task_template.TaskTemplate `json:"task_data"`
	CancelFunc     func()
}

type Status struct {
	TStatus                 TStatus     `json:"task_status"`
	Stages                  int         `json:"stages"`
	FailedCommandId         string      `json:"failed_command_id"`
	FailedCommandTemplateId int         `json:"failed_command_template_id"`
	FailedMessage           *util.MyErr `json:"failed_message"`
}

type TaskStage struct {
	Name       string           `json:"name"`
	CommandId  string           `json:"command_id"`
	From       time.Time        `json:"from"`
	To         *time.Time       `json:"to"`
	Status     e_command.Status `json:"status"`
	Message    *util.MyErr      `json:"message"`
	Tags       json.RawMessage  `json:"tags"`
	TemplateId int              `json:"template_id"`
}
