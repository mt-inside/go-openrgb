package main

import (
	"fmt"
	"os"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mt-inside/go-openrgb/pkg/logging"
	"github.com/mt-inside/go-openrgb/pkg/model"
)

func main() {
	log := logging.GetLogger(false)

	m, err := model.NewModel(log, "localhost:6742", "mt is skill")
	if err != nil {
		log.Error(err, "Couldn't synchronise devices and colors from server")
		os.Exit(1)
	}
	fmt.Println(m)

	m.Devices.ByNameUnwrap("B550 Vision D").Modes.DirectUnwrap().Zones.ByNameUnwrap("D_LED1 Bottom").Leds[0].SetColor(colorful.Hsv(120, 1, 1))
	m.Devices.ByNameUnwrap("B550 Vision D").Modes.DirectUnwrap().Zones.ByNameUnwrap("D_LED1 Bottom").Leds[18].SetColor(colorful.Hsv(120, 1, 1))
	//m.Devices.ByNameUnwrap("B550 Vision D").Modes.DirectUnwrap().Zones.ByNameUnwrap("D_LED1 Bottom").Leds[0].SetColor(colorful.Hsv(0, 0, 0))
	m.Devices.ByNameUnwrap("B550 Vision D").Modes.DirectUnwrap().Zones.ByNameUnwrap("D_LED1 Bottom").Leds[18].SetColor(colorful.Hsv(0, 0, 0))

	m.Diff()

	err = m.Thither()
	if err != nil {
		log.Error(err, "Couldn't synchronise colors to server")
		os.Exit(1)
	}
}
