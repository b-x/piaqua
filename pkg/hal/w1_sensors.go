package hal

import (
	"piaqua/pkg/config"
	"piaqua/pkg/w1therm"
	"time"
)

type W1Sensors struct {
	ids []string
}

const updateSensorsInterval = time.Second * 30

func (s *W1Sensors) Init(hwConf *config.HardwareConf) error {
	s.ids = hwConf.Sensors
	return nil
}

func (s *W1Sensors) Loop(quit <-chan struct{}, events chan<- Event) {
	s.readSensors(events) // immediate first read
	ticker := time.NewTicker(updateSensorsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			s.readSensors(events)
		}
	}
}

func (s *W1Sensors) readSensors(events chan<- Event) {
	for n, id := range s.ids {
		temp, err := w1therm.Temperature(id)
		if err != nil {
			events <- TemperatureError{n, err}
			continue
		}
		events <- TemperatureRead{n, temp}
	}
}
