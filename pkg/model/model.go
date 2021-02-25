package model

import (
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mt-inside/go-openrgb/pkg/wire"
)

type Model struct {
	log     logr.Logger
	client  *wire.Client
	Devices DeviceList
}

func NewModel(log logr.Logger, addr, userAgent string) (*Model, error) {
	model := &Model{log: log.WithName("model").WithValues("server", addr)}

	c, err := wire.NewClient(model.log, "localhost:6742", "mt is skill")
	if err != nil {
		return nil, fmt.Errorf("Couldn't connect: %w", err)
	}
	model.client = c

	err = model.Thence(model.log)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (model *Model) Thence(log logr.Logger) error {
	wireDevs, err := wire.FetchDevices(model.client)
	if err != nil {
		return fmt.Errorf("Couldn't get devices: %w", err)
	}

	var deviceCount, modeCount, zoneCount, ledCount int

	for _, wireDev := range wireDevs {
		modelDev := &Device{
			index:         wireDev.Index,
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
				modelMode.colors = make([]colorful.Color, len(wireMode.Colors))
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
							name:        wireLeds[i].Name,
							serverColor: wireColors[i],
							newColor:    wireColors[i],
						}
						ledCount++
					}
					modelMode.Zones = append(modelMode.Zones, modelZone)
					ledOffset += wireZone.TotalLEDs
					zoneCount++
				}
				modelDev.Modes = append(modelDev.Modes, modelMode)
			}
			modeCount++
		}
		model.Devices = append(model.Devices, modelDev)
		deviceCount++
	}

	model.log.Info(
		"Synchronised from server",
		"devices", deviceCount,
		"modes", modeCount,
		"zones", zoneCount,
		"leds", ledCount,
	)

	return nil
}

// TODO: handle non-direct modes
// FIXME: assumes there's one direct mode and it's the one we're using. What we should do is:
// * switch on the (new) active mode and only write that - device.colors if it's a per-led mode, mode.colors if it's not. No point sync'ing anything else. While it wouldn't be an error to, what do you put in device->colors if there's >1 per-led mode?
// * update active mode (how?)
// TODO: optimise: walk the Diff object. If only one LED has changed, use cmdUpdateSingularLED, etc
func (m *Model) Thither() error {
	devs := []*wire.WriteDevice{}

	for _, d := range m.Devices {
		offset := 0

		ms := d.Modes.Directs()
		if len(ms) > 1 {
			m.log.Info("FIXME: Multiple direct modes not supported")
			continue
		}
		m := ms[0] // hack

		devColCount := 0
		for _, z := range m.Zones {
			devColCount += len(z.Leds)
		}
		colors := make([]colorful.Color, devColCount)

		for _, z := range m.Zones {
			for i, l := range z.Leds {
				colors[offset+i] = l.newColor
			}
			offset += len(z.Leds)
		}
		dev := &wire.WriteDevice{Index: d.index, Colors: colors}
		devs = append(devs, dev)
	}

	m.log.Info("Synchronising to server")

	return wire.WriteDevices(m.client, devs)
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

func (m *Model) Diff() { // TODO: return object
	var path [4]string
	for _, d := range m.Devices {
		path[0] = d.name
		for _, m := range d.Modes.Directs() {
			path[1] = m.name
			for _, z := range m.Zones {
				path[2] = z.name
				for _, l := range z.Leds {
					path[3] = l.name
					if l.serverColor != l.newColor {
						fmt.Printf("%s: %s -> %s\n", strings.Join(path[:], "/"), l.serverColor.Hex(), l.newColor.Hex())
					}
				}
			}
		}
	}
}
