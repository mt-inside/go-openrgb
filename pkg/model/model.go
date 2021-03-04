package model

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-logr/logr"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mt-inside/go-openrgb/pkg/wire"
)

/* TODO
* - commit all this
* - add the fake zone in effectModes, remove so much complexity in there and in mode.go
* - setting colors for non-direct modes doesn't work?
* - deal with switchying modes - should be able to set and get active mode, and most of the *Direct stuff actually wants to work on the currently active mode, panicing if it's not direct
*   - do the PR to openrgb
* - deal with writing non-direct modes (see Thither())
 */

// TODO rename to System
type Model struct {
	log     logr.Logger
	client  *wire.Client
	Devices DeviceList
}

type Settable interface {
	GetName() string
	Size() int
	SetColor(c colorful.Color)
	SetColors(cs []colorful.Color)
	Diff() // TODO: return object
}

func NewModel(log logr.Logger, addr, userAgent string) (*Model, error) {
	model := &Model{log: log.WithValues("server", addr)}

	c, err := wire.NewClient(model.log, "localhost:6742", userAgent)
	if err != nil {
		return nil, fmt.Errorf("Couldn't connect: %w", err)
	}
	model.client = c

	err = model.Thence()
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (m *Model) Size() int {
	size := 0
	for _, d := range m.Devices {
		size += d.Size()
	}
	return size
}

func (m *Model) SetColor(c colorful.Color) {
	for _, d := range m.Devices {
		d.SetColor(c)
	}
}

func (m *Model) SetColors(cs []colorful.Color) {
	if m.Size() != len(cs) {
		panic(fmt.Errorf("Trying to set %d-led Model with %d colors.", m.Size(), len(cs)))
	}
	i := 0
	for _, d := range m.Devices {
		d.SetColors(cs[i : i+d.Size()])
		i += d.Size()
	}
}

func (model *Model) Thence() error {
	wireDevs, err := wire.FetchDevices(model.client)
	if err != nil {
		return fmt.Errorf("Couldn't get devices: %w", err)
	}

	var deviceCount, modeCount, zoneCount, ledCount int

	for _, wireDev := range wireDevs {
		modelDev := &Device{
			log:         model.log.WithName("device").WithValues("device", wireDev.Name),
			model:       model,
			index:       wireDev.Index,
			devType:     wireDev.Type,
			name:        wireDev.Name,
			description: wireDev.Description,
			version:     wireDev.Version,
			serial:      wireDev.Serial,
			location:    wireDev.Location,
		}
		for _, wireMode := range wireDev.Modes {
			if wireMode.ColorMode != wire.PerLED {
				modelMode := &EffectMode{
					device:   modelDev,
					name:     wireMode.Name,
					wireMode: wireMode,
				}
				modelMode.Colors = make([]*LED, len(wireMode.Colors))
				for i, wireC := range wireMode.Colors {
					modelMode.Colors[i] = &LED{
						mode:        modelMode,
						name:        fmt.Sprintf("Color %d", i),
						serverColor: wireC,
						newColor:    wireC,
					}
				}
				modelDev.Modes = append(modelDev.Modes, modelMode)
			} else {
				ledOffset := uint32(0)
				modelMode := &DirectMode{
					device:   modelDev,
					name:     wireMode.Name,
					wireMode: wireMode,
				}
				for _, wireZone := range wireDev.Zones {
					modelZone := &Zone{
						mode:         modelMode,
						name:         wireZone.Name,
						zoneType:     wireZone.Type,
						minLEDs:      wireZone.MinLEDs,
						maxLEDs:      wireZone.MaxLEDs,
						matrixWidth:  wireZone.MatrixWidth,
						matrixHeight: wireZone.MatrixHeight,
					}

					modelZone.Leds = make([]*LED, wireZone.TotalLEDs)
					wireLeds := wireDev.LEDs[ledOffset : ledOffset+wireZone.TotalLEDs]
					wireColors := wireDev.Colors[ledOffset : ledOffset+wireZone.TotalLEDs]
					for i := uint32(0); i < wireZone.TotalLEDs; i++ {
						modelZone.Leds[i] = &LED{
							zone:        modelZone,
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
		modelDev.serverActiveMode = modelDev.Modes[wireDev.ActiveModeIdx]
		modelDev.newActiveMode = modelDev.serverActiveMode

		model.Devices = append(model.Devices, modelDev)
		deviceCount++
	}

	model.log.V(1).Info(
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
	for _, d := range m.Devices {
		mode := d.GetActiveMode()

		// Update the active mode.
		// The API is clunky - to change the colors array, we have to re-assert everything else about it
		// The command we have to use to do this also sets this mode to be the active one
		// - Hence we only do this for the currently active mode, even if there's diffs to the others
		// Note: this is in fact the only way to set the active mode
		// TODO: assumes it's a direct!
		// do we blank out speed and colors and stuff if we read a direct?
		// TODO do we clear the diff of the other modes, and of the LEDs if we don't walk them?
		// TODO hack it unto building, then add the EM::Zone
		wireMode := *mode.getWireMode()

		if em, ok := mode.(*EffectMode); ok {
			colors := make([]colorful.Color, d.Size())
			for i, l := range em.Colors {
				colors[i] = l.newColor
				l.serverColor = l.newColor
			}
			wireMode.Colors = colors
		}

		spew.Dump(wireMode)
		err := wire.SendUpdateMode(m.client, d.index, &wireMode)
		if err != nil {
			return err
		}

		if dm, ok := mode.(*DirectMode); ok {
			colors := make([]colorful.Color, d.Size())
			offset := 0
			for _, z := range dm.Zones {
				for i, l := range z.Leds {
					colors[offset+i] = l.newColor
					l.serverColor = l.newColor
				}
				offset += len(z.Leds)
			}

			err = wire.SendUpdateLEDs(m.client, &wire.UpdateLEDs{DeviceID: d.index, Colors: colors})
			if err != nil {
				return err
			}
		}
	}

	m.log.V(1).Info("Synchronised to server")

	return nil
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

// TODO: lots of duplicated code. Each Diff should defer to the layer below it. Those things should have parent pointers so they can build their full path
func (m *Model) Diff() {
	var path [4]string
	for _, d := range m.Devices {
		path[0] = d.GetName()
		for _, m := range d.Modes.Directs() {
			path[1] = m.GetName()
			for _, z := range m.Zones {
				path[2] = z.GetName()
				for _, l := range z.Leds {
					path[3] = l.GetName()
					if l.serverColor != l.newColor {
						fmt.Printf("%s: %s -> %s\n", strings.Join(path[:], "/"), l.serverColor.Hex(), l.newColor.Hex())
					}
				}
			}
		}
	}
}
