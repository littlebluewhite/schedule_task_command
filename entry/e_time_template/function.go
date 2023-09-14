package e_time_template

import "schedule_task_command/entry/e_time_data"

func (tt *TimeTemplate) GetTimeData() e_time_data.TimeDatum {
	return tt.TimeData
}
