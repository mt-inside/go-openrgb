# go-openrgb
Golang SDK for [OpenRGB](https://gitlab.com/CalcProgrammer1/OpenRGB)

[![Checks](https://github.com/mt-inside/go-openrgb/actions/workflows/checks.yaml/badge.svg)](https://github.com/mt-inside/go-openrgb/actions/workflows/checks.yaml)
[![GitHub Issues](https://img.shields.io/github/issues-raw/mt-inside/go-openrgb)](https://github.com/mt-inside/go-openrgb/issues)

[![Go Reference](https://pkg.go.dev/badge/github.com/mt-inside/go-openrgb.svg)](https://pkg.go.dev/github.com/mt-inside/go-openrgb)

Uses the [lm-sensors](https://github.com/lm-sensors/lm-sensors) (linux monitoring sensors) pacakge, on top of the [hwmon](https://hwmon.wiki.kernel.org) kernel feature.

## Setup
* Install and Configure _OpenRGB_
  * See: https://gitlab.com/CalcProgrammer1/OpenRGB
* Run OpenRGB and start its network server (from the GUI, or `openrgb --server`)
* `go get github.com/mt-inside/go-openrgb`

## How it works
This module calls OpenRGB over the network.

## Example

### Code
```go
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
```

## Contributing
PRs welcome!

### Reaching me
* Raise an issue or PR here
* Discord: mt_inside#8886
* Twitter: @mt165
