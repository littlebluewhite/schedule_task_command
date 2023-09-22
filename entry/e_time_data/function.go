package e_time_data

import (
	"github.com/goccy/go-json"
	"gorm.io/datatypes"
	"schedule_task_command/util"
	"time"
)

func (td *TimeDatum) CheckTimeData(t time.Time) (result bool) {
	result = true
	ch := make(chan bool, 3)
	go func(t time.Time, ch chan bool) {
		ch <- td.checkTime(t)
	}(t, ch)
	go func(t time.Time, ch chan bool) {
		ch <- td.checkDate(t)
	}(t, ch)
	go func(t time.Time, ch chan bool) {
		ch <- td.checkCondition(t)
	}(t, ch)
	for i := 0; i < 3; i++ {
		select {
		case b := <-ch:
			if b == false {
				result = false
			}
		}
	}
	return
}

func (td *TimeDatum) checkTime(t time.Time) (result bool) {
	var startTime datatypes.Time
	var endTime datatypes.Time
	if err := startTime.UnmarshalJSON([]byte(td.StartTime)); err != nil {
		return
	}
	if err := endTime.UnmarshalJSON([]byte(td.EndTime)); err != nil {
		return
	}
	startInt := int(startTime)
	nowInt := util.GetTimeInt(t)
	endInt := int(endTime)
	if !(startInt <= nowInt && endInt+999999999 >= nowInt) {
		return
	}
	if td.IntervalSeconds == nil || *td.IntervalSeconds == 0 {
		result = true
		return
	}
	if duration := nowInt - startInt; (duration/int(time.Second))%int(*td.IntervalSeconds) == 0 {
		result = true
	}
	return
}

func (td *TimeDatum) checkDate(t time.Time) (result bool) {
	if td.EndDate == nil {
		if td.StartDate.Unix() <= t.Unix() {
			result = true
		}
	} else {
		if td.StartDate.Unix() <= t.Unix() && (*td.EndDate).Add(24*time.Hour).Unix() > t.Unix() {
			result = true
		}
	}
	return
}

func (td *TimeDatum) checkCondition(t time.Time) (result bool) {
	if td.RepeatType == repeatNone {
		result = true
	} else {
		switch td.RepeatType {
		case daily:
			result = true
		case weekly:
			result = td.checkWeekly(t)
		case monthly:
			result = td.checkMonthly(t)
		default:
			result = false
		}
	}
	return
}

func (td *TimeDatum) checkWeekly(t time.Time) (result bool) {
	if td.ConditionType == conditionNone {
		return
	}
	if td.ConditionType != weeklyDay {
		return
	}
	var conditions []int
	if err := json.Unmarshal(td.TCondition, &conditions); err != nil {
		return
	}
	result = util.Contains[int]([]int{int(t.Weekday())}, conditions)
	return
}

func (td *TimeDatum) checkMonthly(t time.Time) (result bool) {
	if td.ConditionType == conditionNone {
		return
	}
	var conditions []int
	if err := json.Unmarshal(td.TCondition, &conditions); err != nil {
		return
	}
	weekCount := util.CountWeek(t)
	switch td.ConditionType {
	case monthlyDay:
		result = util.Contains[int]([]int{t.Day()}, conditions)
	case weeklyFirst:
		if weekCount == 0 {
			result = util.Contains[int]([]int{int(t.Weekday())}, conditions)
		}
	case weeklySecond:
		if weekCount == 1 {
			result = util.Contains[int]([]int{int(t.Weekday())}, conditions)
		}
	case weeklyThird:
		if weekCount == 2 {
			result = util.Contains[int]([]int{int(t.Weekday())}, conditions)
		}
	case weeklyFourth:
		if weekCount == 3 {
			result = util.Contains[int]([]int{int(t.Weekday())}, conditions)
		}
	}
	return
}
