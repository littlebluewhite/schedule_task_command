// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameMonitor = "monitor"

// Monitor mapped from table <monitor>
type Monitor struct {
	ID                int32        `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	StatusCode        int32        `gorm:"column:status_code;not null" json:"status_code"`
	Interval          int32        `gorm:"column:interval;not null" json:"interval"`
	CommandTemplateID int32        `gorm:"column:command_template_id;not null" json:"command_template_id"`
	MConditions       []MCondition `gorm:"foreignKey:monitor_id" json:"m_conditions"`
}

// TableName Monitor's table name
func (*Monitor) TableName() string {
	return TableNameMonitor
}