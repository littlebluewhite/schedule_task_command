package e_schedule

import (
	"github.com/littlebluewhite/schedule_task_command/dal/model"
	"github.com/littlebluewhite/schedule_task_command/entry/e_time_data"
)

func Format(sd []model.Schedule) []Schedule {
	result := make([]Schedule, 0, len(sd))
	for _, item := range sd {
		i := Schedule{
			ID:             item.ID,
			Name:           item.Name,
			Description:    item.Description,
			TaskTemplateID: item.TaskTemplateID,
			Enabled:        item.Enabled,
			UpdatedAt:      item.UpdatedAt,
			CreatedAt:      item.CreatedAt,
			Tags:           item.Tags,
			TimeData: e_time_data.TimeDatum{
				RepeatType:      e_time_data.S2RepeatType(item.TimeData.RepeatType),
				StartDate:       item.TimeData.StartDate,
				EndDate:         item.TimeData.EndDate,
				StartTime:       string(item.TimeData.StartTime),
				EndTime:         string(item.TimeData.EndTime),
				IntervalSeconds: item.TimeData.IntervalSeconds,
				ConditionType:   e_time_data.S2ConditionType(item.TimeData.ConditionType),
				TCondition:      item.TimeData.TCondition,
			},
		}
		result = append(result, i)
	}
	return result
}

func Model2Entry(s model.Schedule) Schedule {
	se := Schedule{
		ID:             s.ID,
		Name:           s.Name,
		Description:    s.Description,
		TaskTemplateID: s.TaskTemplateID,
		Enabled:        s.Enabled,
		UpdatedAt:      s.UpdatedAt,
		CreatedAt:      s.CreatedAt,
		Tags:           s.Tags,
		TimeData: e_time_data.TimeDatum{
			RepeatType:      e_time_data.S2RepeatType(s.TimeData.RepeatType),
			StartDate:       s.TimeData.StartDate,
			EndDate:         s.TimeData.EndDate,
			StartTime:       string(s.TimeData.StartTime),
			EndTime:         string(s.TimeData.EndTime),
			IntervalSeconds: s.TimeData.IntervalSeconds,
			ConditionType:   e_time_data.S2ConditionType(s.TimeData.ConditionType),
			TCondition:      s.TimeData.TCondition,
		},
	}
	return se
}

func CreateConvert(c []*ScheduleCreate) []*model.Schedule {
	result := make([]*model.Schedule, 0, len(c))
	for _, item := range c {
		i := model.Schedule{
			Name:           item.Name,
			Description:    item.Description,
			TaskTemplateID: item.TaskTemplateID,
			Enabled:        item.Enabled,
			Tags:           item.Tags,
			TimeData: model.TimeDatum{
				RepeatType:      item.TimeData.RepeatType.ToModel(),
				StartDate:       item.TimeData.StartDate,
				EndDate:         item.TimeData.EndDate,
				StartTime:       []byte(item.TimeData.StartTime.String()),
				EndTime:         []byte(item.TimeData.EndTime.String()),
				IntervalSeconds: item.TimeData.IntervalSeconds,
				ConditionType:   item.TimeData.ConditionType.ToModel(),
				TCondition:      item.TimeData.TCondition,
			},
		}
		result = append(result, &i)
	}
	return result
}

func UpdateConvert(sMap map[int]model.Schedule, us []*ScheduleUpdate) (result []*model.Schedule, err error) {
	for _, u := range us {
		s, ok := sMap[int(u.ID)]
		if !ok {
			err = ScheduleNotFound(int(u.ID))
			return
		}
		if u.Name != nil {
			s.Name = *u.Name
		}
		s.Description = u.Description
		if u.TaskTemplateID != nil {
			s.TaskTemplateID = *u.TaskTemplateID
		}
		if u.Enabled != nil {
			s.Enabled = *u.Enabled
		}
		if u.Tags != nil {
			s.Tags = *u.Tags
		}
		if u.TimeData != nil {
			s.TimeData.RepeatType = u.TimeData.RepeatType.ToModel()
			s.TimeData.StartDate = u.TimeData.StartDate
			s.TimeData.EndDate = u.TimeData.EndDate
			s.TimeData.StartTime = []byte(u.TimeData.StartTime.String())
			s.TimeData.EndTime = []byte(u.TimeData.EndTime.String())
			s.TimeData.IntervalSeconds = u.TimeData.IntervalSeconds
			s.TimeData.ConditionType = u.TimeData.ConditionType.ToModel()
			s.TimeData.TCondition = u.TimeData.TCondition
		}
		result = append(result, &s)
	}
	return
}
