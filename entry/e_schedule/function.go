package e_schedule

import (
	"time"
)

func (s *Schedule) CheckTimeData(t time.Time) (result bool) {
	return s.TimeData.CheckTimeData(t)
}
