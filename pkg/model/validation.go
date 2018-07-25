package model

import (
	"time"
)

func isValidDuration24h(duration time.Duration) bool {
	return duration > 0 && duration < time.Hour*24
}

func isValidWeekdays(weekdays int) bool {
	return weekdays > 0 && weekdays < (1<<7)
}

func (t *RelayTask) IsValid() bool {
	return isValidDuration24h(t.Start) && isValidDuration24h(t.Stop) && isValidWeekdays(t.Weekdays)
}

func (a *Action) IsValid(numRelays, numButtons int) bool {
	return a.Duration > 0 &&
		a.Relay >= 0 && a.Relay < numRelays &&
		a.Button >= -1 && a.Button < numButtons
}
