package e_command

import (
	"github.com/goccy/go-json"
	"schedule_task_command/entry/e_command_template"
	"schedule_task_command/util"
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

type Command struct {
	CommandId      string                             `json:"command_id"`
	Token          string                             `json:"token"`
	From           time.Time                          `json:"from"`
	To             *time.Time                         `json:"to"`
	TriggerFrom    []string                           `json:"trigger_from"`
	TriggerAccount string                             `json:"trigger_account"`
	StatusCode     int                                `json:"status_code"`
	RespData       json.RawMessage                    `json:"resp_data"`
	Status         Status                             `json:"status"`
	Message        *util.MyErr                        `json:"message"`
	TemplateId     int                                `json:"template_id"`
	Template       e_command_template.CommandTemplate `json:"template"`
	CancelFunc     func()
}

type CommandPub struct {
	CommandId      string                             `json:"command_id"`
	Token          string                             `json:"token"`
	From           time.Time                          `json:"from"`
	To             *time.Time                         `json:"to"`
	TriggerFrom    []string                           `json:"trigger_from"`
	TriggerAccount string                             `json:"trigger_account"`
	StatusCode     int                                `json:"status_code"`
	RespData       json.RawMessage                    `json:"resp_data"`
	Status         Status                             `json:"status"`
	Message        string                             `json:"message"`
	TemplateID     int                                `json:"template_id"`
	Template       e_command_template.CommandTemplate `json:"template"`
}
