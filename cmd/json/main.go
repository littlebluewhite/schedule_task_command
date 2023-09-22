package main

import (
	"fmt"
	"github.com/goccy/go-json"
	"time"
)

func main() {
	t := Command{Message: MyError("asdf")}
	tb, e := json.Marshal(t)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(string(tb))
}

type Status int

const (
	Prepared Status = iota
	Process
	Success
	Failure
	Cancel
)

func (s *Status) String() string {
	return [...]string{"Prepared", "Process", "Success", "Failure", "Cancel"}[*s]
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

type Command struct {
	CommandId      string          `json:"command_id"`
	Token          string          `json:"token"`
	From           time.Time       `json:"from"`
	To             *time.Time      `json:"to"`
	TriggerFrom    []string        `json:"trigger_from"`
	TriggerAccount string          `json:"trigger_account"`
	StatusCode     int             `json:"status_code"`
	RespData       json.RawMessage `json:"resp_data"`
	Status         Status          `json:"status"`
	Message        error           `json:"message"`
	TemplateID     int             `json:"template_id"`
}

type MyError string

func (m MyError) Error() string {
	return string(m)
}

//func (m MyError) MarshalJSON() ([]byte, error) {
//	return json.Marshal(string(m))
//}
