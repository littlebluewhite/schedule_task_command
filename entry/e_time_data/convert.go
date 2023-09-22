package e_time_data

func S2RepeatType(s *string) RepeatType {
	if s == nil {
		return repeatNone
	}
	switch *s {
	case "daily":
		return daily
	case "weekly":
		return weekly
	case "monthly":
		return monthly
	default:
		return repeatNone
	}
}

func S2ConditionType(s *string) ConditionType {
	if s == nil {
		return conditionNone
	}
	switch *s {
	case "monthly_day":
		return monthlyDay
	case "weekly_day":
		return weeklyDay
	case "weekly_first":
		return weeklyFirst
	case "weekly_second":
		return weeklySecond
	case "weekly_third":
		return weeklyThird
	case "weekly_fourth":
		return weeklyFourth
	default:
		return conditionNone
	}
}
