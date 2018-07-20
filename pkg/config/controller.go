package config

import (
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// ControllerConf controller config
type ControllerConf struct {
	Sensors []Sensor
	Relays  []Relay
	Actions map[int]Action
}

// Sensor config
type Sensor struct {
	Name string
}

// Relay config
type Relay struct {
	Name  string
	Tasks map[int]RelayTask
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
}

const controllerConfigFilename = "/controller.yml"

func (conf *ControllerConf) Read(dir string) error {
	content, err := ioutil.ReadFile(dir + controllerConfigFilename)
	if os.IsNotExist(err) {
		return nil
	}
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

func (conf *ControllerConf) Validate(hwConf *HardwareConf) {
	if slen := len(hwConf.Sensors); slen != len(conf.Sensors) {
		conf.Sensors = make([]Sensor, slen)
	}
	if rlen := len(hwConf.Relays); rlen != len(conf.Relays) {
		conf.Relays = make([]Relay, rlen)
	}
}
