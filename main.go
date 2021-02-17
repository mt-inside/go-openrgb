package main

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-logr/logr"
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

	time.Sleep(time.Hour)
}
