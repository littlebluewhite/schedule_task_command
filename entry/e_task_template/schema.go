package e_task_template

import (
	"fmt"
	"github.com/goccy/go-json"
	"schedule_task_command/entry/e_command_template"
	"schedule_task_command/util"
	"time"
)

type TaskTemplate struct {
	ID        int32           `json:"id"`
	Name      string          `json:"name"`
	Variable  json.RawMessage `json:"variable"`
	UpdatedAt *time.Time      `json:"updated_at"`
	CreatedAt *time.Time      `json:"created_at"`
	Stages    []TaskStage     `json:"stages"`
	Tags      json.RawMessage `json:"tags"`
}

type TaskStage struct {
	ID                int32                              `json:"id"`
	Name              string                             `json:"name"`
	StageNumber       int32                              `json:"stage_number"`
	Mode              Mode                               `json:"mode"`
	CommandTemplateID int32                              `json:"command_template_id,omitempty"`
	Tags              json.RawMessage                    `json:"tags"`
	CommandTemplate   e_command_template.CommandTemplate `json:"command_template,omitempty"`
}

type TaskTemplateCreate struct {
	Name     string            `json:"name" binding:"required"`
	Variable json.RawMessage   `json:"variable"`
	Stages   []TaskStageCreate `json:"stages"`
	Tags     json.RawMessage   `json:"tags"`
}

type TaskStageCreate struct {
	Name              string          `json:"name" binding:"required"`
	StageNumber       int32           `json:"stage_number" binding:"required"`
	Mode              Mode            `json:"mode" binding:"required"`
	CommandTemplateID int32           `json:"command_template_id"`
	Tags              json.RawMessage `json:"tags"`
}

type TaskTemplateUpdate struct {
	ID       int32             `json:"id" binding:"required"`
	Name     *string           `json:"name"`
	Variable json.RawMessage   `json:"variable"`
	Stages   []TaskStageUpdate `json:"stages"`
	Tags     json.RawMessage   `json:"tags"`
}

type TaskStageUpdate struct {
	ID                int32           `json:"id"`
	Name              *string         `json:"name" binding:"required"`
	StageNumber       *int32          `json:"stage_number" binding:"required"`
	Mode              Mode            `json:"mode" binding:"required"`
	CommandTemplateID *int32          `json:"command_template_id"`
	Tags              json.RawMessage `json:"tags"`
}

type Mode int

const (
	NoneMode Mode = iota
	Monitor
	Execute
)

func (m Mode) String() string {
	return [...]string{"", "monitor", "execute"}[m]
}

func (m Mode) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

func (m *Mode) UnmarshalJSON(data []byte) error {
	var tStatus string
	err := json.Unmarshal(data, &tStatus)
	if err != nil {
		return err
	}
	*m = S2Mode(&tStatus)
	return nil
}

func TaskTemplateNotFound(id int) util.MyErr {
	e := fmt.Sprintf("task template id: %d not found", id)
	return util.MyErr(e)
}

type SendTaskTemplate struct {
	TemplateId     int      `json:"template_id"`
	TriggerFrom    []string `json:"trigger_from"`
	TriggerAccount string   `json:"trigger_account"`
	Token          string   `json:"token"`
}
