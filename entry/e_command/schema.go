package e_command

import (
	"github.com/goccy/go-json"
	"github.com/littlebluewhite/schedule_task_command/entry/e_command_template"
	"github.com/littlebluewhite/schedule_task_command/util"
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
	ID             uint64                             `json:"id"`
	Token          string                             `json:"token"`
	From           time.Time                          `json:"from"`
	To             *time.Time                         `json:"to"`
	Variables      map[string]string                  `json:"variables"`
	Source         string                             `json:"source"`
	TriggerFrom    []map[string]string                `json:"trigger_from"`
	TriggerAccount string                             `json:"trigger_account"`
	StatusCode     int                                `json:"status_code"`
	RespData       json.RawMessage                    `json:"resp_data"`
	Status         Status                             `json:"status"`
	ClientMessage  string                             `json:"client_message"`
	Message        *util.MyErr                        `json:"message"`
	TemplateId     int32                              `json:"template_id"`
	Return         map[string]string                  `json:"return"`
	CommandData    e_command_template.CommandTemplate `json:"command_data"`
	CancelFunc     func()
}

type CommandPub struct {
	ID             uint64                             `json:"id"`
	Token          string                             `json:"token"`
	From           time.Time                          `json:"from"`
	To             *time.Time                         `json:"to"`
	Variables      map[string]string                  `json:"variables"`
	Source         string                             `json:"source"`
	TriggerFrom    []map[string]string                `json:"trigger_from"`
	TriggerAccount string                             `json:"trigger_account"`
	StatusCode     int                                `json:"status_code"`
	RespData       json.RawMessage                    `json:"resp_data"`
	Status         Status                             `json:"status"`
	ClientMessage  string                             `json:"client_message"`
	Message        string                             `json:"message"`
	TemplateID     int32                              `json:"template_id"`
	Return         map[string]string                  `json:"return"`
	CommandData    e_command_template.CommandTemplate `json:"command_data"`
}
