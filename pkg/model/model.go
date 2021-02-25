package model

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mt-inside/go-openrgb/pkg/wire"
)

type Model struct {
	log     logr.Logger
	client  *wire.Client
	Devices []*Device
}

func NewModel(log logr.Logger, addr, userAgent string) (*Model, error) {
	model := &Model{log: log}

	c, err := wire.NewClient(log, "localhost:6742", "mt is skill")
	if err != nil {
		return nil, fmt.Errorf("Couldn't connect: %w", err)
	}
	model.client = c

	wireDevs, err := wire.FetchDevices(model.client)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get devices: %w", err)
	}
	log.Info("Got devices from the wire", "count", len(wireDevs))

	for _, wireDev := range wireDevs {
		modelDev := &Device{
			devType:       wireDev.Type,
			name:          wireDev.Name,
			description:   wireDev.Description,
			version:       wireDev.Version,
			serial:        wireDev.Serial,
			location:      wireDev.Location,
			activeModeIdx: wireDev.ActiveModeIdx,
		}
		for _, wireMode := range wireDev.Modes {
			if wireMode.ColorMode != wire.PerLED {
				modelMode := &EffectMode{
					name:      wireMode.Name,
					flags:     wireMode.Flags,
					minSpeed:  wireMode.MinSpeed,
					speed:     wireMode.Speed,
					maxSpeed:  wireMode.MaxSpeed,
					direction: wireMode.Direction,
					colorMode: wireMode.ColorMode,
					minColors: wireMode.MinColors,
					maxColors: wireMode.MaxColors,
				}
				modelMode.colors = make([]*colorful.Color, len(wireMode.Colors))
				copy(modelMode.colors, wireMode.Colors)
				modelDev.Modes = append(modelDev.Modes, modelMode)
			} else {
				ledOffset := uint32(0)
				modelMode := &DirectMode{
					name: wireMode.Name,
				}
				for _, wireZone := range wireDev.Zones {
					modelZone := &Zone{
						name:     wireZone.Name,
						zoneType: wireZone.Type,
						minLEDs:  wireZone.MinLEDs,
						maxLEDs:  wireZone.MaxLEDs,
					}

					modelZone.Leds = make([]*LED, wireZone.TotalLEDs)
					wireLeds := wireDev.LEDs[ledOffset : ledOffset+wireZone.TotalLEDs]
					wireColors := wireDev.Colors[ledOffset : ledOffset+wireZone.TotalLEDs]
					for i := uint32(0); i < wireZone.TotalLEDs; i++ {
						modelZone.Leds[i] = &LED{
							name:  wireLeds[i].Name,
							color: wireColors[i],
						}
					}
					modelMode.Zones = append(modelMode.Zones, modelZone)
					ledOffset += wireZone.TotalLEDs
				}
				modelDev.Modes = append(modelDev.Modes, modelMode)
			}
		}
		model.Devices = append(model.Devices, modelDev)
	}

	return model, nil
}

func (m *Model) render(indent int) []indentedString {
	ss := make([]indentedString, len(m.Devices))

	for _, device := range m.Devices {
		ss = append(ss, device.render(indent)...)
	}

	return ss
}
func (m *Model) String() string {
	return renderIndents(m.render(0))
}
