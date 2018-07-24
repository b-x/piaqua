package controller

import (
	"fmt"
	"log"
	"piaqua/pkg/config"
	"piaqua/pkg/hal"
	"sync"
)

// Controller aquarium controller
type Controller struct {
	configDir string
	hwConf    config.HardwareConf
	conf      config.ControllerConf
	stop      chan struct{}
	events    chan hal.Event
	allDone   sync.WaitGroup
	sensors   hal.Sensors
	pins      hal.Pins
	state     ControllerState
	lastID    int
	mutex     sync.Mutex
}

// NewController creates and runs a controller
func NewController(configDir string) (*Controller, error) {
	c := &Controller{configDir: configDir}
	err := c.hwConf.Read(configDir)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read hw config from %s: %s", configDir, err.Error())
	}

	err = c.conf.Read(configDir)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read config from %s: %s", configDir, err.Error())
	}

	err = c.conf.Validate(&c.hwConf)
	if err != nil {
		return nil, fmt.Errorf("Invalid config %s: %s", configDir, err.Error())
	}

	c.stop = make(chan struct{})
	c.events = make(chan hal.Event, 16)
	c.sensors.Init(&c.hwConf)

	err = c.pins.Init(&c.hwConf)
	if err != nil {
		return nil, fmt.Errorf("Couldn't init pins: %s", err.Error())
	}

	c.init()

	eventSources := []hal.EventSource{&c.sensors, &c.pins}

	c.allDone.Add(len(eventSources) + 1)

	for _, source := range eventSources {
		go func(es hal.EventSource) {
			defer c.allDone.Done()
			es.Loop(c.stop, c.events)
		}(source)
	}

	go func() {
		defer c.allDone.Done()
		c.processEvents(c.stop, c.events)
	}()

	return c, nil
}

// Stop stops controller
func (c *Controller) Stop() {
	close(c.stop)
	c.allDone.Wait()

	c.pins.Cleanup()
}

func (c *Controller) init() {
	for i := range c.conf.Relays {
		for j := range c.conf.Relays[i].Tasks {
			if j > c.lastID {
				c.lastID = j
			}
		}
	}
	for i := range c.conf.Actions {
		if i > c.lastID {
			c.lastID = i
		}
	}
}

func (c *Controller) newID() int {
	c.lastID++
	return c.lastID
}

func (c *Controller) saveConfig() {
	if err := c.conf.Write(c.configDir); err != nil {
		log.Println("couldn't save config:", err)
	}
}

func (c *Controller) processEvents(quit <-chan struct{}, events <-chan hal.Event) {
	for {
		select {
		case <-quit:
			return
		case event := <-events:
			switch e := event.(type) {
			case hal.ButtonPressed:
				log.Printf("Button %d pressed\n", e.ID)
			case hal.TemperatureRead:
				log.Printf("Sensor %d temp: %2.1f\n", e.ID, float32(e.Temp/100)/10)
			case hal.TemperatureError:
				log.Printf("Sensor %d error: %s\n", e.ID, e.Error)
			}
		}
	}
}
