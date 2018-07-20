package singleinstance

import (
	"net"
)

// Lock given name
func Lock(name string) bool {
	_, err := net.ListenUnix("unix", &net.UnixAddr{Name: "@/run/" + name, Net: "unix"})
	return err == nil
}
