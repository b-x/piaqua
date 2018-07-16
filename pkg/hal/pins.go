package hal

import (
	"piaqua/pkg/config"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type Pins struct {
	buttons []rpio.Pin
	relays  []rpio.Pin

	buttonsStates []rpio.State
}

const updateButtonsInterval = time.Millisecond * 100

func (p *Pins) Init(hwConf *config.HardwareConf) error {
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

func (p *Pins) Cleanup() {
	for _, pin := range p.buttons {
		pin.PullOff()
	}
	for _, pin := range p.relays {
		pin.Write(rpio.Low)
	}
	rpio.Close()
}

func (p *Pins) Loop(quit <-chan struct{}, events chan<- Event) {
	ticker := time.Tick(updateButtonsInterval)
	for {
		select {
		case <-quit:
			return
		case <-ticker:
			p.updateButtonsState(events)
		}
	}
}

func (p *Pins) updateButtonsState(events chan<- Event) {
	for i, pin := range p.buttons {
		// EdgeDetected() doesn't work...
		state := pin.Read()
		if state == p.buttonsStates[i] {
			continue
		}
		p.buttonsStates[i] = state
		if state == rpio.Low {
			events <- ButtonPressed{i}
		}
	}
}

func (p *Pins) SetRelay(id int, on bool) {
	state := rpio.Low
	if on {
		state = rpio.High
	}
	p.relays[id].Write(state)
}
