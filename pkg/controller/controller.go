package controller

import (
	"fmt"
	"os"
	"piaqua/pkg/config"
	"piaqua/pkg/hal"
	"strconv"
	"sync"
	"time"
)

// Controller aquarium controller
type Controller struct {
	configDir string
	hwConf    config.HardwareConf
	state     config.ControllerConf
	stop      chan struct{}
	events    chan hal.Event
	allDone   sync.WaitGroup
	sensors   hal.Sensors
	pins      hal.Pins
	lastID    int
	mutex     sync.Mutex
}

const updateStateInterval = time.Second

// NewController creates and runs a controller
func NewController(configDir string) (*Controller, error) {
	c := &Controller{configDir: configDir}
	err := c.hwConf.Read(configDir)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read hw config from %s: %s", configDir, err.Error())
	}

	err = c.state.Read(configDir)
	if os.IsNotExist(err) {
		c.state.Init(&c.hwConf)
		err = c.state.Write(configDir)
		if err != nil {
			return nil, fmt.Errorf("Couldn't write config to %s: %s", configDir, err.Error())
		}
		err = c.state.Read(configDir)
	}
	if err != nil {
		return nil, fmt.Errorf("Couldn't read config from %s: %s", configDir, err.Error())
	}
	err = c.state.CheckValid(&c.hwConf)
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
	for i := range c.state.Sensors {
		sensor := &c.state.Sensors[i]
		if sensor.Name == "" {
			sensor.Name = strconv.Itoa(i + 1)
		}
	}
	for i := range c.state.Relays {
		relay := &c.state.Relays[i]
		for j := range relay.Tasks {
			if j > c.lastID {
				c.lastID = j
			}
		}
		if relay.Name == "" {
			relay.Name = strconv.Itoa(i + 1)
		}
	}
	for i := range c.state.Actions {
		if i > c.lastID {
			c.lastID = i
		}
	}

	c.updateState(time.Now())
}

func (c *Controller) newID() int {
	c.lastID++
	return c.lastID
}

func (c *Controller) saveConfig() {
	if err := c.state.Write(c.configDir); err != nil {
		fmt.Println("couldn't save config:", err)
	}
}

func (c *Controller) processEvents(quit <-chan struct{}, events <-chan hal.Event) {
	ticker := time.Tick(updateStateInterval)
	for {
		select {
		case <-quit:
			return
		case <-ticker:
			c.onUpdateState()
		case event := <-events:
			switch e := event.(type) {
			case hal.ButtonPressed:
				c.onButtonPressed(e)
			case hal.TemperatureRead:
				c.onTempRead(e)
			case hal.TemperatureError:
				c.onTempError(e)
			}
		}
	}
}

func (c *Controller) onTempRead(e hal.TemperatureRead) {
	value := float32(e.Temp) / 1000

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.state.Sensors[e.ID].Value = &value
}

func (c *Controller) onTempError(e hal.TemperatureError) {
	fmt.Printf("Sensor %d error: %s\n", e.ID, e.Error)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.state.Sensors[e.ID].Value = nil
}

func (c *Controller) onButtonPressed(e hal.ButtonPressed) {
	fmt.Printf("Button %d pressed\n", e.ID)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()

	for _, action := range c.state.Actions {
		if action.Button != nil && *action.Button == e.ID {
			action.Toggle(now)
		}
	}

	c.applyChanges(now)
}

func (c *Controller) onUpdateState() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.updateState(time.Now())
}

func (c *Controller) applyChanges(t time.Time) {
	c.saveConfig()
	c.updateState(t)
}

func (c *Controller) updateState(t time.Time) {
	c.state.Update(t)
	c.updatePins()
}

func (c *Controller) updatePins() {
	for i := range c.state.Relays {
		relay := &c.state.Relays[i]
		if c.pins.GetRelay(i) != relay.On {
			c.pins.SetRelay(i, relay.On)

			if relay.On {
				fmt.Printf("Relay on:  '%s'\n", relay.Name)
			} else {
				fmt.Printf("Relay off: '%s'\n", relay.Name)
			}
		}
	}
}
