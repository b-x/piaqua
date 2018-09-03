package config

import (
	"fmt"
	"io/ioutil"
	"piaqua/pkg/model"

	"gopkg.in/yaml.v2"
)

// ControllerConf controller config
type ControllerConf model.System

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
	conf.Sensors = make(model.Sensors, len(hwConf.Sensors))
	conf.Relays = make(model.Relays, len(hwConf.Relays))
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
		if !action.IsValid(len(hwConf.Relays), len(hwConf.Buttons)) {
			return fmt.Errorf("Invalid action %d", i)
		}
	}
	return nil
}
