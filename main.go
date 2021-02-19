package main

import (
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
		colorful.Hsv(300, 1, 0.1),
		colorful.Hsv(300, 1, 0.1),
		colorful.Hsv(300, 1, 0.1),
		colorful.Hsv(300, 1, 0.1),
		colorful.Hsv(300, 1, 0.1),
		colorful.Hsv(300, 1, 0.1),
		colorful.Hsv(300, 1, 0.1),
		colorful.Hsv(300, 1, 0.1),
	})
	err = c.sendCommand(1, cmdUpdateLEDs, buf)
	if err != nil {
		log.Error(err, "Can't set device-level LEDs")
	}
	err = c.sendCommand(2, cmdUpdateLEDs, buf)
	if err != nil {
		log.Error(err, "Can't set device-level LEDs")
	}

	bufZ := getCommandZoneLEDs(1, []colorful.Color{
		colorful.Hsv(240, 1, 0.1),
		colorful.Hsv(240, 1, 0.1),
		colorful.Hsv(240, 1, 0.1),
		colorful.Hsv(240, 1, 0.1),
		colorful.Hsv(240, 1, 0.1),
		colorful.Hsv(240, 1, 0.1),
		colorful.Hsv(240, 1, 0.1),
		colorful.Hsv(240, 1, 0.1),
		colorful.Hsv(240, 1, 0.1),
		colorful.Hsv(240, 1, 0.1),
		colorful.Hsv(240, 1, 0.1),
		colorful.Hsv(240, 1, 0.1),
	})
	err = c.sendCommand(3, cmdUpdateZoneLEDs, bufZ)
	if err != nil {
		log.Error(err, "Can't set zone-level LEDs")
	}
}
