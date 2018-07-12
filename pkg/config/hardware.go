package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// HardwareConf hardware config
type HardwareConf struct {
	Sensors []string
	Buttons []uint8
	Relays  []uint8
}

func (conf *HardwareConf) Read(dir string) error {
	content, err := ioutil.ReadFile(dir + "/hardware.yml")
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, conf)
}
