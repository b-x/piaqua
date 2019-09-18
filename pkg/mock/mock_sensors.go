package mock

import (
	"piaqua/pkg/config"
	"piaqua/pkg/hal"
)

type MockSensors struct {
	temp []int
}

func (s *MockSensors) Init(hwConf *config.HardwareConf) error {
	s.temp = make([]int, len(hwConf.Sensors))
	return nil
}

func (s *MockSensors) Loop(quit <-chan struct{}, events chan<- hal.Event) {
	for n, t := range s.temp {
		events <- hal.TemperatureRead{ID: n, Temp: t}
	}
	<-quit
}
