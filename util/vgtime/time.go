package vgtime

import "time"

// RefreshHour 每日更新时间
const RefreshHour = 0

// Day 24h
const Day = time.Hour * 24

// LastDailyRefreshTime 计算上一次每日刷新时间
func LastDailyRefreshTime() time.Time {
	curTime := time.Now()
	var duration time.Duration = 0
	duration += time.Duration(-curTime.Nanosecond()) * time.Nanosecond
	duration += time.Duration(-curTime.Second()) * time.Second
	duration += time.Duration(-curTime.Minute()) * time.Minute
	duration += time.Duration(RefreshHour-curTime.Hour()) * time.Hour
	if curTime.Hour() < RefreshHour {
		duration -= Day
	}

	return curTime.Add(duration)
}

// LastDailyRefreshTs 计算上一次每日刷新时间戳
func LastDailyRefreshTs() int64 {
	return LastDailyRefreshTime().Unix() * 1000
}

// NextDailyRefreshTime 计算下一次每日刷新时间
func NextDailyRefreshTime() time.Time {
	return LastDailyRefreshTime().Add(Day)
}

// NextDailyRefreshTs 计算下一次每日刷新时间戳
func NextDailyRefreshTs() int64 {
	return NextDailyRefreshTime().Unix() * 1000
}

// NeedDailyRefreshByTime 计算是否需要刷新
func NeedDailyRefreshByTime(tm time.Time) bool {
	return tm.Before(LastDailyRefreshTime())
}

// NeedDailyRefreshByTs 计算是否需要刷新
func NeedDailyRefreshByTs(ts int64) bool {
	return ts < LastDailyRefreshTs()
}

// Now 获取当前毫秒时间戳
func Now() int64 {
	t := time.Now()
	return t.Unix()*1000 + int64(t.Nanosecond())/1000000
}

// GetDayZeroTs 获取当天0点时间戳
func GetDayZeroTs(ts int64) int64 {
	tm := time.Unix(ts/1000, 0)
	year, month, day := tm.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix() * 1000
}

// GetWeekZeroTs 获取周一0点时间戳
func GetWeekZeroTs(ts int64) int64 {
	tm := time.Unix(ts/1000, 0)
	var duration time.Duration = 0
	duration += time.Duration(-tm.Nanosecond()) * time.Nanosecond
	duration += time.Duration(-tm.Second()) * time.Second
	duration += time.Duration(-tm.Minute()) * time.Minute
	duration += time.Duration(-tm.Hour()) * time.Hour
	if tm.Weekday() == time.Sunday {
		duration += time.Duration(-time.Duration(6) * Day)
	} else {
		duration += time.Duration(-time.Duration(tm.Weekday()-1) * Day)
	}

	return tm.Add(duration).Unix() * 1000
}

// LastWeeklyRefreshTs 获取上一次每周刷新时间
func LastWeeklyRefreshTs() int64 {
	return GetWeekZeroTs(Now()) + RefreshHour*3600*1000
}
