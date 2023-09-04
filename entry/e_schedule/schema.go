package e_schedule

import (
	"github.com/goccy/go-json"
	"gorm.io/datatypes"
	"time"
)

type Schedule struct {
	ID          int32           `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	TimeDataID  int32           `json:"time_data_id"`
	TaskID      *int32          `json:"task_id"`
	Enabled     bool            `json:"enabled"`
	UpdatedAt   *time.Time      `json:"updated_at"`
	CreatedAt   *time.Time      `json:"created_at"`
	TimeData    TimeDatum       `json:"time_data"`
	Tags        json.RawMessage `json:"tags"`
}

type TimeDatum struct {
	RepeatType      *string         `json:"repeat_type"`
	StartDate       time.Time       `json:"start_date"`
	EndDate         *time.Time      `json:"end_date"`
	StartTime       string          `json:"start_time"`
	EndTime         string          `json:"end_time"`
	IntervalSeconds *int32          `json:"interval_seconds"`
	ConditionType   *string         `json:"condition_type"`
	TCondition      json.RawMessage `json:"t_condition"`
}

type ScheduleCreate struct {
	Name        string          `json:"name" binding:"required"`
	Description *string         `json:"description"`
	TaskID      *int32          `json:"task_id"`
	Enabled     bool            `json:"enabled"`
	TimeData    TimeDatumCreate `json:"time_data" binding:"required"`
	Tags        json.RawMessage `json:"tags"`
}

type TimeDatumCreate struct {
	RepeatType      *string         `json:"repeat_type"`
	StartDate       time.Time       `json:"start_date" binding:"required"`
	EndDate         *time.Time      `json:"end_date"`
	StartTime       *datatypes.Time `json:"start_time" binding:"required"`
	EndTime         datatypes.Time  `json:"end_time" binding:"required"`
	IntervalSeconds *int32          `json:"interval_seconds"`
	ConditionType   *string         `json:"condition_type"`
	TCondition      json.RawMessage `json:"t_condition"`
}

type ScheduleUpdate struct {
	ID          int32            `json:"id" binding:"required"`
	Name        *string          `json:"name"`
	Description *string          `json:"description"`
	TaskID      *int32           `json:"task_id"`
	Enabled     *bool            `json:"enabled"`
	TimeData    *TimeDatumUpdate `json:"time_data"`
	Tags        *json.RawMessage `json:"tags"`
}

type TimeDatumUpdate struct {
	RepeatType      *string         `json:"repeat_type" binding:"required"`
	StartDate       time.Time       `json:"start_date" binding:"required"`
	EndDate         *time.Time      `json:"end_date"`
	StartTime       *datatypes.Time `json:"start_time" binding:"required"`
	EndTime         datatypes.Time  `json:"end_time" binding:"required"`
	IntervalSeconds *int32          `json:"interval_seconds"`
	ConditionType   *string         `json:"condition_type"`
	TCondition      json.RawMessage `json:"t_condition" binding:"required"`
}
