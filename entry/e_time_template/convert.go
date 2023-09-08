package e_time_template

import "schedule_task_command/dal/model"

func Format(tt []model.TimeTemplate) []TimeTemplate {
	result := make([]TimeTemplate, 0, len(tt))
	for _, item := range tt {
		i := TimeTemplate{
			ID:        item.ID,
			Name:      item.Name,
			UpdatedAt: item.UpdatedAt,
			CreatedAt: item.CreatedAt,
			TimeData: TimeDatum{
				RepeatType:      item.TimeData.RepeatType,
				StartDate:       item.TimeData.StartDate,
				EndDate:         item.TimeData.EndDate,
				StartTime:       string(item.TimeData.StartTime),
				EndTime:         string(item.TimeData.EndTime),
				IntervalSeconds: item.TimeData.IntervalSeconds,
				ConditionType:   item.TimeData.ConditionType,
				TCondition:      item.TimeData.TCondition,
			},
		}
		result = append(result, i)
	}
	return result
}

func CreateConvert(c []*TimeTemplateCreate) []*model.TimeTemplate {
	result := make([]*model.TimeTemplate, 0, len(c))
	for _, item := range c {
		i := model.TimeTemplate{
			Name: item.Name,
			TimeData: model.TimeDatum{
				RepeatType:      item.TimeData.RepeatType,
				StartDate:       item.TimeData.StartDate,
				EndDate:         item.TimeData.EndDate,
				StartTime:       []byte(item.TimeData.StartTime.String()),
				EndTime:         []byte(item.TimeData.EndTime.String()),
				IntervalSeconds: item.TimeData.IntervalSeconds,
				ConditionType:   item.TimeData.ConditionType,
				TCondition:      item.TimeData.TCondition,
			},
		}
		result = append(result, &i)
	}
	return result
}

func UpdateConvert(ttMap map[int]model.TimeTemplate, utt []*TimeTemplateUpdate) (result []*model.TimeTemplate, err error) {
	for _, u := range utt {
		tt, ok := ttMap[int(u.ID)]
		if !ok {
			err = TimeTemplateNotFound(int(u.ID))
			return
		}
		if u.Name != nil {
			tt.Name = *u.Name
		}
		if u.TimeData != nil {
			tt.TimeData.RepeatType = u.TimeData.RepeatType
			tt.TimeData.StartDate = u.TimeData.StartDate
			tt.TimeData.EndDate = u.TimeData.EndDate
			tt.TimeData.StartTime = []byte(u.TimeData.StartTime.String())
			tt.TimeData.EndTime = []byte(u.TimeData.EndTime.String())
			tt.TimeData.IntervalSeconds = u.TimeData.IntervalSeconds
			tt.TimeData.ConditionType = u.TimeData.ConditionType
			tt.TimeData.TCondition = u.TimeData.TCondition
		}
		result = append(result, &tt)
	}
	return
}
