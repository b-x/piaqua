package controller

import (
	"errors"
	"piaqua/pkg/config"
)

var errID = errors.New("id out of bounds")

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

func (c *Controller) AddRelayTask(relayID int, task config.RelayTask) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if relayID < 0 || relayID >= len(c.conf.Relays) {
		return 0, errID
	}
	relay := &c.conf.Relays[relayID]
	//TODO validate
	id := c.newID()
	relay.Tasks[id] = task
	c.saveConfig()
	return id, nil
}

func (c *Controller) UpdateRelayTask(relayID int, taskID int, task config.RelayTask) error {
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
	//TODO validate
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
