package main

import (
	"os"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mt-inside/go-openrgb/pkg/model"
	"github.com/mt-inside/go-usvc"
)

func main() {
	log := usvc.GetLogger(false, 10)

	m, err := model.NewModel(log, "localhost:6742", "go-openrgb basic example")
	if err != nil {
		log.Error(err, "Couldn't synchronise devices and colors from server")
		os.Exit(1)
	}

	m.SetColor(colorful.Color{R: 1, G: 0, B: 0})
	m.Devices.MustByName("B550 Vision D").MustGetDirectModeAndActivate().SetColor(colorful.Color{R: 0, G: 1, B: 0})
	m.Devices.MustByName("B550 Vision D").MustGetDirectModeAndActivate().Zones.MustByName("D_LED1 Bottom").SetColor(colorful.Color{R: 0, G: 0, B: 1})
	m.Devices.MustByName("B550 Vision D").MustGetDirectModeAndActivate().Zones.MustByName("D_LED1 Bottom").Leds[0].SetColor(colorful.Color{R: 1, G: 1, B: 1})

	m.Diff()

	err = m.Thither()
	if err != nil {
		log.Error(err, "Couldn't synchronise colors to server")
		os.Exit(1)
	}
}
