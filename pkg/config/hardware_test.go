package config

import (
	"fmt"
	"log"
)

func ExampleHardwareConf_Read() {
	conf := HardwareConf{}
	err := conf.Read("../../configs")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(conf.Sensors))
	fmt.Println(len(conf.Buttons))
	fmt.Println(len(conf.Relays))
	// Output:
	// 2
	// 3
	// 6
}
