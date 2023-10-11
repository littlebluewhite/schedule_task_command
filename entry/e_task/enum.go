package e_task

import (
	"github.com/goccy/go-json"
	"schedule_task_command/entry/e_task_template"
	"time"
)

type TStatus int

const (
	Prepared TStatus = iota
	Process
	Success
	Failure
	Cancel
)

func (s TStatus) String() string {
	return [...]string{"Prepared", "Process", "Success", "Failure", "Cancel"}[s]
}

func (s TStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *TStatus) UnmarshalJSON(data []byte) error {
	var tStatus string
	err := json.Unmarshal(data, &tStatus)
	if err != nil {
		return err
	}
	*s = S2Status(&tStatus)
	return nil
}

type TaskStageC struct {
	Monitor []TaskStage `json:"monitor"`
	Execute []TaskStage `json:"execute"`
}

type TaskPub struct {
	TaskId         string                       `json:"task_id"`
	Token          string                       `json:"token"`
	From           time.Time                    `json:"from"`
	To             *time.Time                   `json:"to"`
	Variables      map[string]map[string]string `json:"variables"`
	TriggerFrom    []string                     `json:"trigger_from"`
	TriggerAccount string                       `json:"trigger_account"`
	Status         Status                       `json:"status"`
	Stages         map[int]TaskStageC           `json:"stages"`
	AccountMessage string                       `json:"account_message"`
	Message        string                       `json:"message"`
	TemplateID     int                          `json:"template_id"`
	Template       e_task_template.TaskTemplate `json:"template"`
}
