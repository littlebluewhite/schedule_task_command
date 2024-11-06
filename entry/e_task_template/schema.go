package e_task_template

import (
	"github.com/goccy/go-json"
	"github.com/littlebluewhite/schedule_task_command/entry/e_command_template"
	"time"
)

type TaskTemplate struct {
	ID         int32           `json:"id"`
	Name       string          `json:"name"`
	Visible    bool            `json:"visible"`
	UpdatedAt  *time.Time      `json:"updated_at"`
	CreatedAt  *time.Time      `json:"created_at"`
	StageItems []StageItem     `json:"stage_items"`
	Tags       json.RawMessage `json:"tags"`
}

type StageItem struct {
	ID                int32                              `json:"id"`
	Name              string                             `json:"name"`
	StageNumber       int32                              `json:"stage_number"`
	Mode              Mode                               `json:"mode"`
	CommandTemplateID int32                              `json:"command_template_id,omitempty"`
	Tags              json.RawMessage                    `json:"tags"`
	Variable          json.RawMessage                    `json:"variable"`
	Parser            []ParserItem                       `json:"parser"`
	CommandTemplate   e_command_template.CommandTemplate `json:"command_template,omitempty"`
}

type ParserItem struct {
	FromKey string `json:"from_key"`
	To      []To   `json:"to"`
}

type To struct {
	Key string `json:"key"`
	ID  int    `json:"id"`
}

type TaskTemplateCreate struct {
	Name       string            `json:"name" binding:"required"`
	Visible    bool              `json:"visible"`
	StageItems []StageItemCreate `json:"stage_items"`
	Tags       json.RawMessage   `json:"tags"`
}

type StageItemCreate struct {
	Name              string          `json:"name" binding:"required"`
	StageNumber       int32           `json:"stage_number" binding:"required"`
	Mode              Mode            `json:"mode" binding:"required"`
	CommandTemplateID int32           `json:"command_template_id"`
	Tags              json.RawMessage `json:"tags"`
	Variable          json.RawMessage `json:"variable"`
}

type TaskTemplateUpdate struct {
	ID         int32             `json:"id" binding:"required"`
	Name       *string           `json:"name"`
	Visible    *bool             `json:"visible"`
	StageItems []StageItemUpdate `json:"stage_items"`
	Tags       json.RawMessage   `json:"tags"`
}

type StageItemUpdate struct {
	ID                int32           `json:"id"`
	Name              *string         `json:"name" binding:"required"`
	StageNumber       *int32          `json:"stage_number" binding:"required"`
	Mode              Mode            `json:"mode" binding:"required"`
	CommandTemplateID *int32          `json:"command_template_id"`
	Tags              json.RawMessage `json:"tags"`
	Variable          json.RawMessage `json:"variable"`
	Parser            json.RawMessage `json:"parser"`
}

type SendTaskTemplate struct {
	TemplateId     int                       `json:"template_id"`
	Source         string                    `json:"source"`
	TriggerFrom    []map[string]string       `json:"trigger_from"`
	TriggerAccount string                    `json:"trigger_account"`
	Token          string                    `json:"token"`
	Variables      map[int]map[string]string `json:"variables"`
}
