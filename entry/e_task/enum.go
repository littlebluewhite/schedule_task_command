package e_task

import (
	"github.com/goccy/go-json"
	"schedule_task_command/entry/e_task_template"
	"time"
)

type Status int

const (
	Prepared Status = iota
	Process
	Success
	Failure
	Cancel
)

func (s Status) String() string {
	return [...]string{"Prepared", "Process", "Success", "Failure", "Cancel"}[s]
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var tStatus string
	err := json.Unmarshal(data, &tStatus)
	if err != nil {
		return err
	}
	*s = S2Status(&tStatus)
	return nil
}

type Stage struct {
	Monitor []StageItem `json:"monitor"`
	Execute []StageItem `json:"execute"`
}

type TaskPub struct {
	ID             uint64                       `json:"id"`
	Token          string                       `json:"token"`
	From           time.Time                    `json:"from"`
	To             *time.Time                   `json:"to"`
	Variables      map[int]map[string]string    `json:"variables"`
	Source         string                       `json:"source"`
	TriggerFrom    []string                     `json:"trigger_from"`
	TriggerAccount string                       `json:"trigger_account"`
	Status         Status                       `json:"status"`
	StageNumber    int32                        `json:"stage_number"`
	Stages         map[int32]Stage              `json:"stages"`
	FailedCommands []FailedCommand              `json:"failed_command"`
	ClientMessage  string                       `json:"client_message"`
	Message        string                       `json:"message"`
	TemplateID     int                          `json:"template_id"`
	TaskData       e_task_template.TaskTemplate `json:"task_data"`
}
