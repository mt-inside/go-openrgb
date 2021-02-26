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

	m.SetColor(colorful.Color{R: 1, G: 0, B: 0})
	m.Devices.MustByName("B550 Vision D").Modes.MustDirect().SetColor(colorful.Color{R: 0, G: 1, B: 0})
	m.Devices.MustByName("B550 Vision D").Modes.MustDirect().Zones.MustByName("D_LED1 Bottom").SetColor(colorful.Color{R: 0, G: 0, B: 1})
	m.Devices.MustByName("B550 Vision D").Modes.MustDirect().Zones.MustByName("D_LED1 Bottom").Leds[0].SetColor(colorful.Color{R: 1, G: 1, B: 1})

	m.Diff()

	err = m.Thither()
	if err != nil {
		log.Error(err, "Couldn't synchronise colors to server")
		os.Exit(1)
	}
}
