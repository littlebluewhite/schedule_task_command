package e_schedule

import "schedule_task_command/entry/e_time_data"

func (s *Schedule) GetTimeData() e_time_data.TimeDatum {
	return s.TimeData
}
