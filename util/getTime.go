package util

import "time"

func GetTimeInt(t time.Time) (d int) {
	d = t.Hour()*int(time.Hour) + t.Minute()*int(time.Minute) + t.Second()*int(time.Second)
	return
}

// CountWeek 計算此時刻是當月第幾個星期
func CountWeek(t time.Time) (c int) {
	monthFirst := t.Add(-time.Duration(t.Day()-1) * 24 * time.Hour)
	c = int(t.Unix()-monthFirst.Unix()) / (60 * 60 * 24) / 7
	return
}
