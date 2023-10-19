package main

import (
	"fmt"
	"github.com/goccy/go-json"
	"time"
)

func main() {
	t := Command{
		Message: MyError("asdf"),
		Stages: map[int]string{
			1: "sss",
		},
		RespData:  []byte(`[{"id": 1,"value": "{{value}}"}]`),
		Variables: map[string]string{"value": "3"},
	}
	tb, e := json.Marshal(t)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(string(tb))

	h := map[string]interface{}{
		"1":   []string{"aaa"},
		"ccc": 3,
	}
	hb, e := json.Marshal(h)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(string(hb))
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
	CommandId      string            `json:"command_id"`
	Token          string            `json:"token"`
	From           time.Time         `json:"from"`
	To             *time.Time        `json:"to"`
	Variables      map[string]string `json:"variables"`
	TriggerFrom    []string          `json:"trigger_from"`
	TriggerAccount string            `json:"trigger_account"`
	StatusCode     int               `json:"status_code"`
	Stages         map[int]string    `json:"stages"`
	RespData       json.RawMessage   `json:"resp_data"`
	Status         Status            `json:"status"`
	Message        error             `json:"message"`
	TemplateID     int               `json:"template_id"`
}

type MyError string

func (m MyError) Error() string {
	return string(m)
}

//func (m MyError) MarshalJSON() ([]byte, error) {
//	return json.Marshal(string(m))
//}
