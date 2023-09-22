package e_time

import (
	"github.com/goccy/go-json"
	"schedule_task_command/entry/e_time_template"
	"schedule_task_command/util"
	"time"
)

type PublishTime struct {
	TemplateId     int                          `json:"template_id"`
	TriggerFrom    []string                     `json:"trigger_from"`
	TriggerAccount string                       `json:"trigger_account"`
	Token          string                       `json:"token"`
	Time           time.Time                    `json:"time"`
	IsTime         bool                         `json:"is_time"`
	Status         Status                       `json:"status"`
	Message        *util.MyErr                  `json:"message"`
	Template       e_time_template.TimeTemplate `json:"template"`
}

type Status int

const (
	Prepared Status = iota
	Success
	Failure
)

func (s *Status) String() string {
	return [...]string{
		"Prepared", "Success", "Failure"}[*s]
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var statusStr string
	err := json.Unmarshal(data, &statusStr)
	if err != nil {
		return err
	}
	*s = S2Status(&statusStr)
	return nil
}
