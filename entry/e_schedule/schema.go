package e_schedule

import (
	"fmt"
	"github.com/goccy/go-json"
	"schedule_task_command/entry/e_time_data"
	"schedule_task_command/util"
	"time"
)

type Schedule struct {
	ID             int32                 `json:"id"`
	Name           string                `json:"name"`
	Description    *string               `json:"description"`
	TimeDataID     int32                 `json:"time_data_id"`
	TaskTemplateID int32                 `json:"task_id"`
	Enabled        bool                  `json:"enabled"`
	UpdatedAt      *time.Time            `json:"updated_at"`
	CreatedAt      *time.Time            `json:"created_at"`
	TimeData       e_time_data.TimeDatum `json:"time_data"`
	Tags           json.RawMessage       `json:"tags"`
}

type ScheduleCreate struct {
	Name           string                      `json:"name" binding:"required"`
	Description    *string                     `json:"description"`
	TaskTemplateID int32                       `json:"task_id"`
	Enabled        bool                        `json:"enabled"`
	TimeData       e_time_data.TimeDatumCreate `json:"time_data" binding:"required"`
	Tags           json.RawMessage             `json:"tags"`
}

type ScheduleUpdate struct {
	ID             int32                        `json:"id" binding:"required"`
	Name           *string                      `json:"name"`
	Description    *string                      `json:"description"`
	TaskTemplateID *int32                       `json:"task_id"`
	Enabled        *bool                        `json:"enabled"`
	TimeData       *e_time_data.TimeDatumUpdate `json:"time_data"`
	Tags           *json.RawMessage             `json:"tags"`
}

func ScheduleNotFound(id int) util.MyErr {
	e := fmt.Sprintf("schedule id: %d not found", id)
	return util.MyErr(e)
}
