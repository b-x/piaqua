package model

import (
	"time"
)

func isValidDuration24h(duration time.Duration) bool {
	return duration > 0 && duration < time.Hour*24
}

func (wd Weekdays) IsValid() bool {
	return wd >= toWeekdays(time.Sunday) && wd <= toWeekdays(time.Saturday)
}

func (t *RelayTask) IsValid() bool {
	return isValidDuration24h(t.Start) && isValidDuration24h(t.Stop) && t.Weekdays.IsValid()
}

func (a *Action) IsValid(numRelays, numButtons int) bool {
	return a.Duration > 0 &&
		a.Relay >= 0 && a.Relay < numRelays &&
		a.Button >= -1 && a.Button < numButtons
}
