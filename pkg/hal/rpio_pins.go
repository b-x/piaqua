package hal

import (
	"piaqua/pkg/config"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type RPIOPins struct {
	buttons []rpio.Pin
	relays  []rpio.Pin

	buttonsStates []rpio.State
}

const updateButtonsInterval = time.Millisecond * 100

func (p *RPIOPins) Init(hwConf *config.HardwareConf) error {
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
		//pin.Write(rpio.Low)
	}
	return nil
}

func (p *RPIOPins) Cleanup() {
	for _, pin := range p.buttons {
		pin.PullOff()
	}
	for _, pin := range p.relays {
		pin.Write(rpio.Low)
	}
	rpio.Close()
}

func (p *RPIOPins) Loop(quit <-chan struct{}, events chan<- Event) {
	ticker := time.NewTicker(updateButtonsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			p.updateButtonsState(events)
		}
	}
}

func (p *RPIOPins) updateButtonsState(events chan<- Event) {
	for i, pin := range p.buttons {
		// don't use EdgeDetected() - requires dtoverlay=gpio-no-irq
		// and doesn't solve button bouncing problem
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

func (p *RPIOPins) GetRelay(id int) bool {
	return p.relays[id].Read() != rpio.Low
}

func (p *RPIOPins) SetRelay(id int, on bool) {
	state := rpio.Low
	if on {
		state = rpio.High
	}
	p.relays[id].Write(state)
}
