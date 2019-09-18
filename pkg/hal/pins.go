package hal

import "piaqua/pkg/config"

type Pins interface {
	EventSource

	Init(hwConf *config.HardwareConf) error
	Cleanup()

	GetRelay(id int) bool
	SetRelay(id int, on bool)
}
