package main

import (
	"fmt"
	"os"

	"github.com/mt-inside/go-openrgb/pkg/logging"
	"github.com/mt-inside/go-openrgb/pkg/model"
)

func main() {
	log := logging.GetLogger(false)

	m, err := model.NewModel(log, "localhost:6742", "mt is skill")
	if err != nil {
		log.Error(err, "Couldn't synchronise devices")
		os.Exit(1)
	}
	fmt.Println(m)

	// buf := getCommandLEDs([]colorful.Color{
	// 	colorful.Hsv(300, 1, 0.1),
	// 	colorful.Hsv(300, 1, 0.1),
	// 	colorful.Hsv(300, 1, 0.1),
	// 	colorful.Hsv(300, 1, 0.1),
	// 	colorful.Hsv(300, 1, 0.1),
	// 	colorful.Hsv(300, 1, 0.1),
	// 	colorful.Hsv(300, 1, 0.1),
	// 	colorful.Hsv(300, 1, 0.1),
	// })
	// err = c.sendCommand(1, cmdUpdateLEDs, buf)
	// if err != nil {
	// 	log.Error(err, "Can't set device-level LEDs")
	// }
	// err = c.sendCommand(2, cmdUpdateLEDs, buf)
	// if err != nil {
	// 	log.Error(err, "Can't set device-level LEDs")
	// }

	// bufZ := getCommandZoneLEDs(1, []colorful.Color{
	// 	colorful.Hsv(240, 1, 0.1),
	// 	colorful.Hsv(240, 1, 0.1),
	// 	colorful.Hsv(240, 1, 0.1),
	// 	colorful.Hsv(240, 1, 0.1),
	// 	colorful.Hsv(240, 1, 0.1),
	// 	colorful.Hsv(240, 1, 0.1),
	// 	colorful.Hsv(240, 1, 0.1),
	// 	colorful.Hsv(240, 1, 0.1),
	// 	colorful.Hsv(240, 1, 0.1),
	// 	colorful.Hsv(240, 1, 0.1),
	// 	colorful.Hsv(240, 1, 0.1),
	// 	colorful.Hsv(240, 1, 0.1),
	// })
	// err = c.sendCommand(3, cmdUpdateZoneLEDs, bufZ)
	// if err != nil {
	// 	log.Error(err, "Can't set zone-level LEDs")
	// }
}
