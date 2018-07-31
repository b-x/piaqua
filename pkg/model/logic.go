package model

import "time"

func (a *Action) IsActive(t time.Time) bool {
	return a.Start.Add(a.Duration).After(t) && !a.Start.After(t)
}

func (rt *RelayTask) IsActive(t time.Time) bool {
	rel := toDuration24h(t)
	today := rt.Weekdays.contains(t.Weekday())
	if rt.Start < rt.Stop {
		return today && rel >= rt.Start && rel < rt.Stop
	}
	tomorrow := rt.Weekdays.contains(prevWeekday(t.Weekday()))
	return today && rel >= rt.Start || tomorrow && rel < rt.Stop
}

func (wd Weekdays) contains(d time.Weekday) bool {
	return wd&toWeekdays(d) != 0
}

func toWeekdays(d time.Weekday) Weekdays {
	return 1 << uint(d)
}

func prevWeekday(d time.Weekday) time.Weekday {
	return (d + 6) % 7
}

func toDuration24h(t time.Time) time.Duration {
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return t.Sub(date)
}
