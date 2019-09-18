package main

import "piaqua/pkg/mock"
import "piaqua/pkg/app"

func main() {
	sensors := &mock.MockSensors{}
	pins := &mock.MockPins{}

	app.Run(sensors, pins)
}
