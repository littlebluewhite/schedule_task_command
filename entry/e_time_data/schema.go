package e_time_data

import (
	"github.com/goccy/go-json"
	"gorm.io/datatypes"
	"time"
)

type TimeDatum struct {
	RepeatType      RepeatType      `json:"repeat_type"`
	StartDate       time.Time       `json:"start_date"`
	EndDate         *time.Time      `json:"end_date"`
	StartTime       string          `json:"start_time"`
	EndTime         string          `json:"end_time"`
	IntervalSeconds *int32          `json:"interval_seconds"`
	ConditionType   ConditionType   `json:"condition_type"`
	TCondition      json.RawMessage `json:"t_condition"`
}

type TimeDatumCreate struct {
	RepeatType      RepeatType      `json:"repeat_type"`
	StartDate       time.Time       `json:"start_date" binding:"required"`
	EndDate         *time.Time      `json:"end_date"`
	StartTime       *datatypes.Time `json:"start_time" binding:"required"`
	EndTime         datatypes.Time  `json:"end_time" binding:"required"`
	IntervalSeconds *int32          `json:"interval_seconds"`
	ConditionType   ConditionType   `json:"condition_type"`
	TCondition      json.RawMessage `json:"t_condition"`
}

type TimeDatumUpdate struct {
	RepeatType      RepeatType      `json:"repeat_type" binding:"required"`
	StartDate       time.Time       `json:"start_date" binding:"required"`
	EndDate         *time.Time      `json:"end_date"`
	StartTime       *datatypes.Time `json:"start_time" binding:"required"`
	EndTime         datatypes.Time  `json:"end_time" binding:"required"`
	IntervalSeconds *int32          `json:"interval_seconds"`
	ConditionType   ConditionType   `json:"condition_type"`
	TCondition      json.RawMessage `json:"t_condition" binding:"required"`
}

type RepeatType int

const (
	repeatNone RepeatType = iota
	daily
	weekly
	monthly
)

func (r *RepeatType) String() string {
	return [...]string{"", "daily", "weekly", "monthly"}[*r]
}

func (r *RepeatType) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *RepeatType) UnmarshalJSON(data []byte) error {
	var repeatTypeStr string
	err := json.Unmarshal(data, &repeatTypeStr)
	if err != nil {
		return err
	}
	*r = S2RepeatType(&repeatTypeStr)
	return nil
}

func (r *RepeatType) ToModel() *string {
	s := r.String()
	if s == "" {
		return nil
	}
	return &s
}

type ConditionType int

const (
	conditionNone ConditionType = iota
	monthlyDay
	weeklyDay
	weeklyFirst
	weeklySecond
	weeklyThird
	weeklyFourth
)

func (c *ConditionType) String() string {
	return [...]string{
		"", "monthly_day", "weekly_day", "weekly_first",
		"weekly_second", "weekly_third", "weekly_fourth"}[*c]
}

func (c *ConditionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *ConditionType) UnmarshalJSON(data []byte) error {
	var repeatTypeStr string
	err := json.Unmarshal(data, &repeatTypeStr)
	if err != nil {
		return err
	}
	*c = S2ConditionType(&repeatTypeStr)
	return nil
}

func (c *ConditionType) ToModel() *string {
	s := c.String()
	if s == "" {
		return nil
	}
	return &s
}

var allWeekDay = [...]int{0, 1, 2, 3, 4, 5, 6}
var allMonthDay = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
