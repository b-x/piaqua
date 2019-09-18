package mock

import (
	"piaqua/pkg/config"
	"piaqua/pkg/hal"
)

type MockPins struct {
	relays []bool
}

func (p *MockPins) Init(hwConf *config.HardwareConf) error {
	p.relays = make([]bool, len(hwConf.Relays))
	return nil
}

func (p *MockPins) Cleanup() {
}

func (p *MockPins) GetRelay(id int) bool {
	return p.relays[id]
}

func (p *MockPins) SetRelay(id int, on bool) {
	p.relays[id] = on
}

func (p *MockPins) Loop(quit <-chan struct{}, events chan<- hal.Event) {
	<-quit
}
