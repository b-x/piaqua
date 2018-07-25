package model

import "time"

type Sensors []Sensor

type Relays []Relay

type Actions map[int]*Action

type Sensor struct {
	Name string
}

type Relay struct {
	Name  string
	Tasks map[int]*RelayTask
}

type RelayTask struct {
	Start    time.Duration
	Stop     time.Duration
	Weekdays int
}

type Action struct {
	Name     string
	Button   int
	Relay    int
	Duration time.Duration
	Start    time.Time
}
