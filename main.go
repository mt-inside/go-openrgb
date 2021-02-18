package main

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-logr/logr"
	"github.com/lucasb-eyer/go-colorful"
)

var (
	log logr.Logger
)

func init() {
	log = getLogger(false)
}

func main() {
	c, err := NewClient("localhost:6742", "mt is skill")
	if err != nil {
		log.Error(err, "Couldn't make client")
	}
	defer c.Close()

	devs, err := FetchDevices(c)
	if err != nil {
		log.Error(err, "Couldn't get devices")
	}
	spew.Dump(devs)

	buf := getCommandLEDs([]colorful.Color{
		colorful.Color{R: 1, G: 0, B: 0},
		colorful.Color{R: 1, G: 0, B: 0},
		colorful.Color{R: 1, G: 0, B: 0},
		colorful.Color{R: 1, G: 0, B: 0},
		colorful.Color{R: 1, G: 0, B: 0},
		colorful.Color{R: 1, G: 0, B: 0},
		colorful.Color{R: 1, G: 0, B: 0},
		colorful.Color{R: 1, G: 0, B: 0},
	})
	err = c.sendCommand(1, cmdUpdateLEDs, buf)
	if err != nil {
		log.Error(err, "Can't set device-level LEDs")
	}

	time.Sleep(time.Hour)
}
