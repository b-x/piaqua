package controller

import (
	"log"
	"piaqua/pkg/config"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type pins struct {
	buttons []rpio.Pin
	relays  []rpio.Pin

	buttonsStates []rpio.State
}

const updateButtonsInterval = time.Millisecond * 100

func (p *pins) init(hwConf *config.HardwareConf) error {
	err := rpio.Open()
	if err != nil {
		return err
	}

	p.buttons = make([]rpio.Pin, len(hwConf.Buttons))
	p.relays = make([]rpio.Pin, len(hwConf.Relays))
	p.buttonsStates = make([]rpio.State, len(hwConf.Buttons))

	for i := range hwConf.Buttons {
		p.buttons[i] = rpio.Pin(hwConf.Buttons[i])
	}
	for i := range hwConf.Relays {
		p.relays[i] = rpio.Pin(hwConf.Relays[i])
	}
	for i := range p.buttonsStates {
		p.buttonsStates[i] = rpio.High
	}
	for _, pin := range p.buttons {
		pin.Input()
		pin.PullUp()
	}
	for _, pin := range p.relays {
		pin.Output()
		pin.Write(rpio.Low)
	}
	return nil
}

func (p *pins) cleanup() {
	for _, pin := range p.buttons {
		pin.PullOff()
	}
	for _, pin := range p.relays {
		pin.Write(rpio.Low)
	}
	rpio.Close()
}

func (p *pins) loop(quit <-chan struct{}) {
	ticker := time.Tick(updateButtonsInterval)
	for {
		select {
		case <-quit:
			return
		case <-ticker:
			p.updateButtonsState()
		}
	}
}

func (p *pins) updateButtonsState() {
	for i, pin := range p.buttons {
		// EdgeDetected() doesn't work...
		state := pin.Read()
		if state == p.buttonsStates[i] {
			continue
		}
		p.buttonsStates[i] = state
		if state == rpio.Low {
			//generate event
			log.Printf("button %d: %d\n", i, state)
		}

	}
}
