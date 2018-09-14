package model

import "time"

type System struct {
	Sensors Sensors `json:"sensors"`
	Relays  Relays  `json:"relays"`
	Actions Actions `json:"actions"`
	Buttons int     `json:"buttons"`
}

type Sensors []Sensor

type Relays []Relay

type Actions map[int]*Action

type Sensor struct {
	Name  string   `json:"name"`
	Value *float32 `json:"value" yaml:"-"`
}

type Relay struct {
	Name  string             `json:"name"`
	On    bool               `json:"on" yaml:"-"`
	Tasks map[int]*RelayTask `json:"tasks"`
}

type RelayTask struct {
	Start    time.Duration `json:"start"`
	Stop     time.Duration `json:"stop"`
	Weekdays Weekdays      `json:"weekdays"`
}

type Action struct {
	Name     string        `json:"name"`
	On       bool          `json:"on" yaml:"-"`
	Button   *int          `json:"button"`
	Relay    int           `json:"relay"`
	Duration time.Duration `json:"duration"`
	Start    time.Time     `json:"start"`
}

type Weekdays int
