// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"encoding/json"
)

const TableNameHTTPSCommand = "https_command"

// HTTPSCommand mapped from table <https_command>
type HTTPSCommand struct {
	ID                int32           `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CommandTemplateID *int32          `gorm:"column:command_template_id" json:"command_template_id"`
	Method            string          `gorm:"column:method;not null" json:"method"`
	URL               string          `gorm:"column:url;not null" json:"url"`
	AuthorizationType *string         `gorm:"column:authorization_type" json:"authorization_type"`
	Params            json.RawMessage `gorm:"column:params;default:json_array()" json:"params"`
	Header            json.RawMessage `gorm:"column:header;default:json_array()" json:"header"`
	BodyType          *string         `gorm:"column:body_type" json:"body_type"`
	Body              json.RawMessage `gorm:"column:body" json:"body"`
}

// TableName HTTPSCommand's table name
func (*HTTPSCommand) TableName() string {
	return TableNameHTTPSCommand
}