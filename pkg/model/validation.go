package model

import (
	"time"
)

func isValidDuration24h(duration time.Duration) bool {
	return duration > 0 && duration < time.Hour*24
}

func (wd Weekdays) IsValid() bool {
	return wd > 0 && wd < 128
}

func (t *RelayTask) IsValid() bool {
	return isValidDuration24h(t.Start) && isValidDuration24h(t.Stop) && t.Weekdays.IsValid()
}

func (a *Action) IsValid(numRelays, numButtons int) bool {
	return a.Name != "" && a.Duration > 0 &&
		a.Relay >= 0 && a.Relay < numRelays &&
		(a.Button == nil || *a.Button < numButtons)
}
