package controller

import (
	"log"
	"piaqua/pkg/config"
	"piaqua/pkg/w1therm"
	"sync/atomic"
	"time"
)

type sensors struct {
	ids    []string
	values []int32
}

const updateInterval = time.Second * 30

func (s *sensors) init(hwConf *config.HardwareConf) {
	s.ids = hwConf.Sensors
	s.values = make([]int32, len(s.ids))
}

func (s *sensors) loop(quit <-chan struct{}) {
	s.readSensors() // immediate first read
	ticker := time.Tick(updateInterval)
	for {
		select {
		case <-quit:
			return
		case <-ticker:
			s.readSensors()
		}
	}
}

func (s *sensors) readSensors() {
	for n, id := range s.ids {
		temp, err := w1therm.Temperature(id)
		if err != nil {
			log.Printf("Could not read sensor %d: %s\n", n, err)
			continue
		}
		atomic.StoreInt32(&s.values[n], int32(temp))
		log.Printf("Sensor %d temp: %2.1f\n", n, float32(temp/100)/10)
	}
}

func (s *sensors) value(index int) int32 {
	return atomic.LoadInt32(&s.values[index])
}
