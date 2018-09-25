package controller

import (
	"encoding/json"
	"errors"
	"piaqua/pkg/model"
	"time"
)

var errID = errors.New("id out of bounds")
var errArg = errors.New("invalid argument")

func (c *Controller) GetControllerState() ([]byte, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	return json.Marshal(c.state)
}

func (c *Controller) SetSensorName(id int, name string) error {
	if name == "" {
		return errArg
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if id < 0 || id >= len(c.state.Sensors) {
		return errID
	}
	sensor := &c.state.Sensors[id]
	if sensor.Name == name {
		return nil
	}
	sensor.Name = name
	c.saveConfig()
	return nil
}

func (c *Controller) SetRelayName(id int, name string) error {
	if name == "" {
		return errArg
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if id < 0 || id >= len(c.state.Relays) {
		return errID
	}
	relay := &c.state.Relays[id]
	if relay.Name == name {
		return nil
	}
	relay.Name = name
	c.saveConfig()
	return nil
}

func (c *Controller) AddRelayTask(relayID int, task *model.RelayTask) (int, error) {
	if !task.IsValid() {
		return 0, errArg
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if relayID < 0 || relayID >= len(c.state.Relays) {
		return 0, errID
	}

	relay := &c.state.Relays[relayID]
	id := c.newID()
	relay.Tasks[id] = task
	c.saveConfig()
	return id, nil
}

func (c *Controller) UpdateRelayTask(relayID int, taskID int, task *model.RelayTask) error {
	if !task.IsValid() {
		return errArg
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if relayID < 0 || relayID >= len(c.state.Relays) {
		return errID
	}

	relay := &c.state.Relays[relayID]
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

	if relayID < 0 || relayID >= len(c.state.Relays) {
		return errID
	}
	relay := &c.state.Relays[relayID]
	_, found := relay.Tasks[taskID]
	if !found {
		return nil
	}

	delete(relay.Tasks, taskID)
	c.saveConfig()
	return nil
}

func (c *Controller) AddAction(action *model.Action) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !action.IsValid(len(c.hwConf.Relays), len(c.hwConf.Buttons)) {
		return 0, errArg
	}

	id := c.newID()
	c.state.Actions[id] = action
	c.saveConfig()
	return id, nil
}

func (c *Controller) UpdateAction(id int, action *model.Action) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !action.IsValid(len(c.hwConf.Relays), len(c.hwConf.Buttons)) {
		return errArg
	}

	_, found := c.state.Actions[id]
	if !found {
		return errID
	}

	c.state.Actions[id] = action
	c.saveConfig()
	return nil
}

func (c *Controller) RemoveAction(id int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, found := c.state.Actions[id]
	if !found {
		return nil
	}
	delete(c.state.Actions, id)
	c.saveConfig()
	return nil
}

func (c *Controller) ToggleAction(id int) error {
	now := time.Now()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	action, found := c.state.Actions[id]
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
