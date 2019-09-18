package hal

import "piaqua/pkg/config"

type Sensors interface {
	EventSource

	Init(hwConf *config.HardwareConf) error
}
