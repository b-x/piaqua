package hal

import (
	"piaqua/pkg/config"
	"piaqua/pkg/w1therm"
	"time"
)

type Sensors struct {
	ids []string
}

const updateSensorsInterval = time.Second * 30

func (s *Sensors) Init(hwConf *config.HardwareConf) {
	s.ids = hwConf.Sensors
}

func (s *Sensors) Loop(quit <-chan struct{}, events chan<- Event) {
	s.readSensors(events) // immediate first read
	ticker := time.Tick(updateSensorsInterval)
	for {
		select {
		case <-quit:
			return
		case <-ticker:
			s.readSensors(events)
		}
	}
}

func (s *Sensors) readSensors(events chan<- Event) {
	for n, id := range s.ids {
		temp, err := w1therm.Temperature(id)
		if err != nil {
			events <- TemperatureError{n, err}
			continue
		}
		events <- TemperatureRead{n, temp}
	}
}
