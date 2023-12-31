// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameTimeTemplate = "time_template"

// TimeTemplate mapped from table <time_template>
type TimeTemplate struct {
	ID         int32      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name       string     `gorm:"column:name;not null" json:"name"`
	Visible    bool       `gorm:"column:visible;not null" json:"visible"`
	TimeDataID int32      `gorm:"column:time_data_id;not null" json:"time_data_id"`
	UpdatedAt  *time.Time `gorm:"column:updated_at" json:"updated_at"`
	CreatedAt  *time.Time `gorm:"column:created_at" json:"created_at"`
	TimeData   TimeDatum  `gorm:"foreignKey:time_data_id" json:"time_data"`
}

// TableName TimeTemplate's table name
func (*TimeTemplate) TableName() string {
	return TableNameTimeTemplate
}
