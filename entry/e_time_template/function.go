package e_time_template

import (
	"time"
)

func (tt *TimeTemplate) CheckTimeData(t time.Time) (result bool) {
	return tt.TimeData.CheckTimeData(t)
}
