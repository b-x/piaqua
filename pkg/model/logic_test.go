package model

import (
	"fmt"
	"testing"
	"time"
)

func TestAction_IsActive(t *testing.T) {
	var tab = []struct {
		daction time.Duration
		dnow    time.Duration
		active  bool
	}{
		{time.Hour, -time.Second, false},
		{time.Hour, 0, true},
		{time.Hour, time.Second, true},
		{time.Hour, time.Hour - time.Second, true},
		{time.Hour, time.Hour, false},
		{time.Hour, time.Hour + time.Second, false},
	}
	t0 := time.Now()
	for n, tt := range tab {
		t.Run(fmt.Sprintf("t%v", n), func(t *testing.T) {
			action := Action{Start: t0, Duration: tt.daction}
			if action.IsActive(t0.Add(tt.dnow)) != tt.active {
				t.Errorf("Action(Duration: %v).IsActive(%v) != %v", tt.daction, tt.dnow, tt.active)
			}
		})
	}
}

func TestRelayTask_IsActive(t *testing.T) {
	toDuration := func(str string) time.Duration {
		var h, m, s time.Duration
		_, err := fmt.Sscanf(str, "%d:%d:%d", &h, &m, &s)
		if err != nil {
			panic("invalid duration: " + str)
		}
		return h*time.Hour + m*time.Minute + s*time.Second
	}
	Tue := toWeekdays(time.Tuesday)
	TueWed := toWeekdays(time.Tuesday) | toWeekdays(time.Wednesday)
	var tab = []struct {
		start    string
		stop     string
		weekdays Weekdays
		dnow     string
		wd       time.Weekday
		active   bool
	}{
		// from < to
		{"07:15:30", "15:30:00", TueWed, "06:00:00", time.Monday, false},
		{"07:15:30", "15:30:00", TueWed, "12:00:00", time.Monday, false},
		{"07:15:30", "15:30:00", TueWed, "16:00:00", time.Monday, false},
		{"07:15:30", "15:30:00", TueWed, "06:00:00", time.Tuesday, false},
		{"07:15:30", "15:30:00", TueWed, "12:00:00", time.Tuesday, true},
		{"07:15:30", "15:30:00", TueWed, "16:00:00", time.Tuesday, false},
		{"07:15:30", "15:30:00", TueWed, "06:00:00", time.Wednesday, false},
		{"07:15:30", "15:30:00", TueWed, "12:00:00", time.Wednesday, true},
		{"07:15:30", "15:30:00", TueWed, "16:00:00", time.Wednesday, false},
		{"07:15:30", "15:30:00", TueWed, "06:00:00", time.Thursday, false},
		{"07:15:30", "15:30:00", TueWed, "12:00:00", time.Thursday, false},
		{"07:15:30", "15:30:00", TueWed, "16:00:00", time.Thursday, false},
		// bounds
		{"07:15:30", "15:30:00", Tue, "07:15:29", time.Tuesday, false},
		{"07:15:30", "15:30:00", Tue, "07:15:30", time.Tuesday, true},
		{"07:15:30", "15:30:00", Tue, "07:15:31", time.Tuesday, true},
		{"07:15:30", "15:30:00", Tue, "15:29:29", time.Tuesday, true},
		{"07:15:30", "15:30:00", Tue, "15:30:00", time.Tuesday, false},
		{"07:15:30", "15:30:00", Tue, "15:30:01", time.Tuesday, false},
		// from > to (ends in the next day)
		{"22:22:22", "05:05:05", TueWed, "04:00:00", time.Monday, false},
		{"22:22:22", "05:05:05", TueWed, "06:00:00", time.Monday, false},
		{"22:22:22", "05:05:05", TueWed, "23:00:00", time.Monday, false},
		{"22:22:22", "05:05:05", TueWed, "04:00:00", time.Tuesday, false},
		{"22:22:22", "05:05:05", TueWed, "06:00:00", time.Tuesday, false},
		{"22:22:22", "05:05:05", TueWed, "23:00:00", time.Tuesday, true},
		{"22:22:22", "05:05:05", TueWed, "04:00:00", time.Wednesday, true},
		{"22:22:22", "05:05:05", TueWed, "06:00:00", time.Wednesday, false},
		{"22:22:22", "05:05:05", TueWed, "23:00:00", time.Wednesday, true},
		{"22:22:22", "05:05:05", TueWed, "04:00:00", time.Thursday, true},
		{"22:22:22", "05:05:05", TueWed, "06:00:00", time.Thursday, false},
		{"22:22:22", "05:05:05", TueWed, "23:00:00", time.Thursday, false},
		{"22:22:22", "05:05:05", TueWed, "04:00:00", time.Friday, false},
		{"22:22:22", "05:05:05", TueWed, "06:00:00", time.Friday, false},
		{"22:22:22", "05:05:05", TueWed, "23:00:00", time.Friday, false},
		// bounds
		{"22:22:22", "05:05:05", Tue, "22:22:21", time.Tuesday, false},
		{"22:22:22", "05:05:05", Tue, "22:22:22", time.Tuesday, true},
		{"22:22:22", "05:05:05", Tue, "22:22:23", time.Tuesday, true},
		{"22:22:22", "05:05:05", Tue, "05:05:04", time.Wednesday, true},
		{"22:22:22", "05:05:05", Tue, "05:05:05", time.Wednesday, false},
		{"22:22:22", "05:05:05", Tue, "05:05:06", time.Wednesday, false},
		// 24h
		{"08:00:00", "08:00:00", Tue, "07:59:59", time.Tuesday, false},
		{"08:00:00", "08:00:00", Tue, "08:00:00", time.Tuesday, true},
		{"08:00:00", "08:00:00", Tue, "08:00:01", time.Tuesday, true},
		{"08:00:00", "08:00:00", Tue, "07:59:59", time.Wednesday, true},
		{"08:00:00", "08:00:00", Tue, "08:00:00", time.Wednesday, false},
		{"08:00:00", "08:00:00", Tue, "08:00:01", time.Wednesday, false},
	}
	_ = tab
	t_sunday := time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local)
	if t_sunday.Weekday() != time.Sunday {
		panic("invalid t_sunday")
	}
	for n, tt := range tab {
		t.Run(fmt.Sprintf("t%v", n), func(t *testing.T) {
			task := RelayTask{toDuration(tt.start), toDuration(tt.stop), tt.weekdays}
			t0 := t_sunday.Add(time.Hour * 24 * time.Duration(tt.wd))
			if task.IsActive(t0.Add(toDuration(tt.dnow))) != tt.active {
				t.Errorf("RelayTask(Start: %v, Stop: %v, Wd: %v).IsActive(T: %v, Wd: %v) != %v",
					tt.start, tt.stop, tt.weekdays, tt.dnow, tt.wd, tt.active)
			}
		})
	}
}
