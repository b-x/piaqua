package model

import "time"

func (a *Action) IsActive(t time.Time) bool {
	return a.Start.Add(a.Duration).After(t)
}

func (rt *RelayTask) IsActive(t time.Time) bool {
	//TODO
	return false
}
