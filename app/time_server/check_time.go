package time_server

import (
	"github.com/goccy/go-json"
	"gorm.io/datatypes"
	"schedule_task_command/util"
	"time"
)

func checkScheduleActive(s schedule, t time.Time) (result bool) {
	result = s.Enabled && checkTimeData(s.TimeData, t)
	return
}

func checkTimeData(td timeDatum, t time.Time) (result bool) {
	result = true
	ch := make(chan bool, 3)
	go func(td timeDatum, t time.Time, ch chan bool) {
		ch <- checkTime(td, t)
	}(td, t, ch)
	go func(td timeDatum, t time.Time, ch chan bool) {
		ch <- checkDate(td, t)
	}(td, t, ch)
	go func(td timeDatum, t time.Time, ch chan bool) {
		ch <- checkCondition(td, t)
	}(td, t, ch)
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

func checkTime(td timeDatum, t time.Time) (result bool) {
	var startTime datatypes.Time
	var endTime datatypes.Time
	if err := startTime.UnmarshalJSON(td.StartTime); err != nil {
		return
	}
	if err := endTime.UnmarshalJSON(td.EndTime); err != nil {
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

func checkDate(td timeDatum, t time.Time) (result bool) {
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

func checkCondition(td timeDatum, t time.Time) (result bool) {
	if td.RepeatType == nil {
		result = true
	} else {
		switch *td.RepeatType {
		case daily.String():
			result = true
		case weekly.String():
			result = checkWeekly(td, t)
		case monthly.String():
			result = checkMonthly(td, t)
		}
	}
	return
}

func checkWeekly(td timeDatum, t time.Time) (result bool) {
	if td.ConditionType == nil {
		return
	}
	if *td.ConditionType != weeklyDay.String() {
		return
	}
	var conditions []int
	if err := json.Unmarshal(td.TCondition, &conditions); err != nil {
		return
	}
	result = util.Contains[int]([]int{int(t.Weekday())}, conditions)
	return
}

func checkMonthly(td timeDatum, t time.Time) (result bool) {
	if td.ConditionType == nil {
		return
	}
	var conditions []int
	if err := json.Unmarshal(td.TCondition, &conditions); err != nil {
		return
	}
	weekCount := util.CountWeek(t)
	switch *td.ConditionType {
	case monthDay.String():
		result = util.Contains[int]([]int{t.Day()}, conditions)
	case weeklyFirst.String():
		if weekCount == 0 {
			result = util.Contains[int]([]int{int(t.Weekday())}, conditions)
		}
	case weeklySecond.String():
		if weekCount == 1 {
			result = util.Contains[int]([]int{int(t.Weekday())}, conditions)
		}
	case weeklyThird.String():
		if weekCount == 2 {
			result = util.Contains[int]([]int{int(t.Weekday())}, conditions)
		}
	case weeklyFourth.String():
		if weekCount == 3 {
			result = util.Contains[int]([]int{int(t.Weekday())}, conditions)
		}
	}
	return
}
