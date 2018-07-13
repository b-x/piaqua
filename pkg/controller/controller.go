package controller

import (
	"fmt"
	"piaqua/pkg/config"
	"sync"
)

// Controller aquarium controller
type Controller struct {
	hwConf  config.HardwareConf
	stop    chan struct{}
	allDone sync.WaitGroup
	sensors sensors
}

// NewController creates and runs a controller
func NewController(configDir string) (*Controller, error) {
	c := &Controller{}
	err := c.hwConf.Read(configDir)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read hw config from %s: %s", configDir, err.Error())
	}

	c.stop = make(chan struct{})
	c.sensors.init(&c.hwConf)

	go func() {
		c.allDone.Add(1)
		defer c.allDone.Done()

		c.sensors.loop(c.stop)
	}()

	return c, nil
}

// Stop stops controller
func (c *Controller) Stop() {
	close(c.stop)
	c.allDone.Wait()
}
