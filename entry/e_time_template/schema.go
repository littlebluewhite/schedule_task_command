package e_time_template

import (
	"fmt"
	"github.com/littlebluewhite/schedule_task_command/entry/e_time_data"
	"github.com/littlebluewhite/schedule_task_command/util"
	"time"
)

type TimeTemplate struct {
	ID        int32                 `json:"id"`
	Name      string                `json:"name"`
	Visible   bool                  `json:"visible"`
	UpdatedAt *time.Time            `json:"updated_at"`
	CreatedAt *time.Time            `json:"created_at"`
	TimeData  e_time_data.TimeDatum `json:"time_data"`
}

type TimeTemplateCreate struct {
	Name     string                      `json:"name" binding:"required"`
	Visible  bool                        `json:"visible"`
	TimeData e_time_data.TimeDatumCreate `json:"time_data" binding:"required"`
}

type TimeTemplateUpdate struct {
	ID       int32                        `json:"id" binding:"required"`
	Name     *string                      `json:"name"`
	Visible  *bool                        `json:"visible"`
	TimeData *e_time_data.TimeDatumUpdate `json:"time_data"`
}

func TimeTemplateNotFound(id int) util.MyErr {
	e := fmt.Sprintf("time template id: %d not found", id)
	return util.MyErr(e)
}
