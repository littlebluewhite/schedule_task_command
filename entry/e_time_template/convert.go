package e_time_template

import (
	"schedule_task_command/dal/model"
	"schedule_task_command/entry/e_time_data"
)

func Format(tt []model.TimeTemplate) []TimeTemplate {
	result := make([]TimeTemplate, 0, len(tt))
	for _, item := range tt {
		i := TimeTemplate{
			ID:        item.ID,
			Name:      item.Name,
			Visible:   item.Visible,
			UpdatedAt: item.UpdatedAt,
			CreatedAt: item.CreatedAt,
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

func CreateConvert(c []*TimeTemplateCreate) []*model.TimeTemplate {
	result := make([]*model.TimeTemplate, 0, len(c))
	for _, item := range c {
		//fmt.Printf("%+v\n", item)
		//fmt.Printf("%[1]T, %+[1]v\n", item.TimeData.RepeatType)
		i := model.TimeTemplate{
			Name:    item.Name,
			Visible: item.Visible,
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
		if u.Visible != nil {
			tt.Visible = *u.Visible
		}
		if u.TimeData != nil {
			tt.TimeData.RepeatType = u.TimeData.RepeatType.ToModel()
			tt.TimeData.StartDate = u.TimeData.StartDate
			tt.TimeData.EndDate = u.TimeData.EndDate
			tt.TimeData.StartTime = []byte(u.TimeData.StartTime.String())
			tt.TimeData.EndTime = []byte(u.TimeData.EndTime.String())
			tt.TimeData.IntervalSeconds = u.TimeData.IntervalSeconds
			tt.TimeData.ConditionType = u.TimeData.ConditionType.ToModel()
			tt.TimeData.TCondition = u.TimeData.TCondition
		}
		result = append(result, &tt)
	}
	return
}
