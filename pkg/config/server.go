package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// ServerConf HTTP server config
type ServerConf struct {
	Address     string
	Path        string
	Credentials map[string]string
}

func (conf *ServerConf) Read(dir string) error {
	content, err := ioutil.ReadFile(dir + "/server.yml")
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, conf)
}
