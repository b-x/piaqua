package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// ControllerConf controller config
type ControllerConf struct {
	Sensors []Sensor
	Relays  []Relay
	Actions map[int]*Action
}

// Sensor config
type Sensor struct {
	Name string
}

// Relay config
type Relay struct {
	Name  string
	Tasks map[int]*RelayTask
}

// RelayTask config
type RelayTask struct {
	Start    time.Duration
	Stop     time.Duration
	Weekdays int
}

// Action config
type Action struct {
	Name     string
	Button   int
	Relay    int
	Duration time.Duration
	Start    time.Time
}

const controllerConfigFilename = "/controller.yml"

func (conf *ControllerConf) Read(dir string) error {
	content, err := ioutil.ReadFile(dir + controllerConfigFilename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, conf)
}

func (conf *ControllerConf) Write(dir string) error {
	content, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dir+controllerConfigFilename, content, 0644)
}

func (conf *ControllerConf) Init(hwConf *HardwareConf) {
	conf.Sensors = make([]Sensor, len(hwConf.Sensors))
	conf.Relays = make([]Relay, len(hwConf.Relays))
}

func (conf *ControllerConf) CheckValid(hwConf *HardwareConf) error {
	if slen := len(hwConf.Sensors); slen != len(conf.Sensors) {
		return fmt.Errorf("Invalid number of sensors")
	}
	if rlen := len(hwConf.Relays); rlen != len(conf.Relays) {
		return fmt.Errorf("Invalid number of relays")
	}

	for i := range conf.Relays {
		for j, task := range conf.Relays[i].Tasks {
			if !task.IsValid() {
				return fmt.Errorf("Invalid relay %d task %d", i, j)
			}
		}
	}
	for i, action := range conf.Actions {
		if !action.IsValid(hwConf) {
			return fmt.Errorf("Invalid action %d", i)
		}
	}
	return nil
}

func isValidDuration(duration time.Duration) bool {
	return duration > 0 && duration < time.Hour*24
}

func isValidWeekdays(weekdays int) bool {
	return weekdays > 0 && weekdays < (1<<7)
}

func (t *RelayTask) IsValid() bool {
	return isValidDuration(t.Start) && isValidDuration(t.Stop) && isValidWeekdays(t.Weekdays)
}

func (a *Action) IsValid(hwConf *HardwareConf) bool {
	return a.Duration > 0 &&
		a.Relay >= 0 && a.Relay < len(hwConf.Relays) &&
		a.Button >= -1 && a.Button < len(hwConf.Buttons)
}

func (a *Action) IsActive(t time.Time) bool {
	return a.Start.Add(a.Duration).After(t)
}
