package main

import (
	"piaqua/pkg/app"
	"piaqua/pkg/hal"
)

func main() {
	sensors := &hal.W1Sensors{}
	pins := &hal.RPIOPins{}

	app.Run(sensors, pins)
}
