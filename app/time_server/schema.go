package time_server

import (
	"github.com/goccy/go-json"
	"schedule_task_command/util"
	"time"
)

type SendTime struct {
	TemplateId     int      `json:"template_id"`
	TriggerFrom    []string `json:"trigger_from"`
	TriggerAccount string   `json:"trigger_account"`
	Token          string   `json:"token"`
}

type publishTime struct {
	TemplateId     int          `json:"template_id"`
	TriggerFrom    []string     `json:"trigger_from"`
	TriggerAccount string       `json:"trigger_account"`
	Token          string       `json:"token"`
	Time           time.Time    `json:"time"`
	IsTime         bool         `json:"is_time"`
	Status         Status       `json:"status"`
	Message        util.JsonErr `json:"message"`
}

type Status int

const (
	Prepared Status = iota
	Success
	Failure
)

func (s Status) String() string {
	return [...]string{
		"Prepared", "Success", "Failure"}[s]
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

var CannotFindTemplate = util.MyErr("can not find time template")
