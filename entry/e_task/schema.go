package e_task

import (
	"github.com/goccy/go-json"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/entry/e_task_template"
	"schedule_task_command/util"
	"time"
)

type Task struct {
	ID             uint64                       `json:"id"`
	Token          string                       `json:"token"`
	From           time.Time                    `json:"from"`
	To             *time.Time                   `json:"to"`
	Variables      map[int]map[string]string    `json:"variables"`
	TriggerFrom    []string                     `json:"trigger_from"`
	TriggerAccount string                       `json:"trigger_account"`
	Status         Status                       `json:"status"`
	StageNumber    int32                        `json:"stage_number"`
	Stages         map[int32]Stage              `json:"stages"`
	FailedCommands []FailedCommand              `json:"failed_command"`
	ClientMessage  string                       `json:"client_message"`
	Message        *util.MyErr                  `json:"message"`
	TemplateId     int                          `json:"template_id"`
	TaskData       e_task_template.TaskTemplate `json:"task_data"`
	CancelFunc     func()
}

type FailedCommand struct {
	CommandID         uint64           `json:"command_id"`
	CommandTemplateID int32            `json:"command_template_id"`
	Message           *util.MyErr      `json:"message"`
	Status            e_command.Status `json:"status"`
}

type StageItem struct {
	Name              string               `json:"name"`
	StageID           int32                `json:"stage_id"`
	CommandID         uint64               `json:"command_id"`
	StageNumber       int32                `json:"stage_number"`
	Mode              e_task_template.Mode `json:"mode"`
	From              *time.Time           `json:"from"`
	To                *time.Time           `json:"to"`
	Status            e_command.Status     `json:"status"`
	Message           *util.MyErr          `json:"message"`
	Tags              json.RawMessage      `json:"tags"`
	Variables         map[string]string    `json:"variable"`
	CommandTemplateId int32                `json:"command_template_id"`
}

type SimpleTask struct {
	ID           uint64            `json:"id"`
	TemplateName string            `json:"template_name"`
	Status       int               `json:"status"`
	StageNumber  int32             `json:"stage_number"`
	StageItems   []SimpleStageItem `json:"stage_items"`
}

type SimpleStageItem struct {
	Name   string          `json:"name"`
	From   *time.Time      `json:"from"`
	To     *time.Time      `json:"to"`
	Status int             `json:"status"`
	Tags   json.RawMessage `json:"tags"`
}
