package controller

import (
	"errors"
	"piaqua/pkg/config"
	"time"
)

var errID = errors.New("id out of bounds")
var errArg = errors.New("invalid argument")

func (c *Controller) SetSensorName(id int, name string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if id < 0 || id >= len(c.conf.Sensors) {
		return errID
	}
	sensor := &c.conf.Sensors[id]
	if sensor.Name == name {
		return nil
	}
	sensor.Name = name
	c.saveConfig()
	return nil
}

func (c *Controller) SetRelayName(id int, name string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if id < 0 || id >= len(c.conf.Relays) {
		return errID
	}
	relay := &c.conf.Relays[id]
	if relay.Name == name {
		return nil
	}
	relay.Name = name
	c.saveConfig()
	return nil
}

func (c *Controller) AddRelayTask(relayID int, task *config.RelayTask) (int, error) {
	if !task.IsValid() {
		return 0, errArg
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if relayID < 0 || relayID >= len(c.conf.Relays) {
		return 0, errID
	}

	relay := &c.conf.Relays[relayID]
	id := c.newID()
	relay.Tasks[id] = task
	c.saveConfig()
	return id, nil
}

func (c *Controller) UpdateRelayTask(relayID int, taskID int, task *config.RelayTask) error {
	if !task.IsValid() {
		return errArg
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if relayID < 0 || relayID >= len(c.conf.Relays) {
		return errID
	}

	relay := &c.conf.Relays[relayID]
	_, found := relay.Tasks[taskID]
	if !found {
		return errID
	}

	relay.Tasks[taskID] = task
	c.saveConfig()
	return nil
}

func (c *Controller) RemoveRelayTask(relayID int, taskID int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if relayID < 0 || relayID >= len(c.conf.Relays) {
		return errID
	}
	relay := &c.conf.Relays[relayID]
	_, found := relay.Tasks[taskID]
	if !found {
		return nil
	}

	delete(relay.Tasks, taskID)
	c.saveConfig()
	return nil
}

func (c *Controller) AddAction(action *config.Action) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !action.IsValid(&c.hwConf) {
		return 0, errArg
	}

	id := c.newID()
	c.conf.Actions[id] = action
	c.saveConfig()
	return id, nil
}

func (c *Controller) UpdateAction(id int, action *config.Action) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !action.IsValid(&c.hwConf) {
		return errArg
	}

	_, found := c.conf.Actions[id]
	if !found {
		return errID
	}

	c.conf.Actions[id] = action
	c.saveConfig()
	return nil
}

func (c *Controller) RemoveAction(id int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, found := c.conf.Actions[id]
	if !found {
		return nil
	}
	delete(c.conf.Actions, id)
	c.saveConfig()
	return nil
}

func (c *Controller) ToggleAction(id int) error {
	now := time.Now()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	action, found := c.conf.Actions[id]
	if !found {
		return errID
	}

	if action.IsActive(now) {
		action.Start = time.Time{}
	} else {
		action.Start = now
	}
	c.saveConfig()
	return nil
}
