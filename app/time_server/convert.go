package time_server

import "schedule_task_command/dal/model"

func modelMap2scheduleMap(msm map[int]model.Schedule) map[int]schedule {
	result := make(map[int]schedule)
	for i, data := range msm {
		result[i] = model2schedule(data)
	}
	return result
}

func model2schedule(ms model.Schedule) schedule {
	return schedule{
		ID:      ms.ID,
		Name:    ms.Name,
		TaskID:  ms.TaskID,
		Enabled: ms.Enabled,
		TimeData: timeDatum{
			RepeatType:      ms.TimeData.RepeatType,
			StartDate:       ms.TimeData.StartDate,
			EndDate:         ms.TimeData.EndDate,
			StartTime:       ms.TimeData.StartTime,
			EndTime:         ms.TimeData.EndTime,
			IntervalSeconds: ms.TimeData.IntervalSeconds,
			ConditionType:   ms.TimeData.ConditionType,
			TCondition:      ms.TimeData.TCondition,
		},
	}
}
