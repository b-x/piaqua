package w1therm

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strconv"
)

var errRead = errors.New("failed to read temperature")
var errParse = errors.New("failed to parse temperature")

// Temperature reads temperature from sensor
func Temperature(id string) (int, error) {
	data, err := ioutil.ReadFile("/sys/bus/w1/devices/" + id + "/w1_slave")
	if err != nil {
		return 0, errRead
	}
	if bytes.Index(data, []byte("YES")) < 0 {
		return 0, errParse
	}
	idx := bytes.LastIndex(data, []byte("t="))
	if idx < 0 {
		return 0, errParse
	}
	val, err := strconv.ParseInt(string(data[idx+2:len(data)-1]), 10, 32)
	if err != nil {
		return 0, errParse
	}
	return int(val), nil
}
