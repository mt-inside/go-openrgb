package model

import (
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mt-inside/go-openrgb/pkg/wire"
)

type DeviceList []*Device

func (list DeviceList) ByName(name string) []*Device {
	ds := []*Device{}

	for _, d := range list {
		if strings.EqualFold(d.GetName(), name) { // case insensitive
			ds = append(ds, d)
		}
	}

	return ds
}
func (list DeviceList) MustByName(name string) *Device {
	ds := list.ByName(name)
	if len(ds) != 1 {
		panic(fmt.Errorf("Not presicely one mode with name: %s", name))
	}
	return ds[0]
}

type Device struct {
	log              logr.Logger
	model            *Model
	index            uint32
	devType          wire.DeviceType
	name             string
	description      string
	version          string
	serial           string
	location         string
	Modes            ModeList
	serverActiveMode Mode
	newActiveMode    Mode
}

func (d *Device) GetActiveMode() Mode {
	return d.newActiveMode
}
func (d *Device) SetActiveMode(m Mode) {
	d.newActiveMode = m
}

func (d *Device) MustActiveDirectMode() *DirectMode {
	m := d.GetActiveMode()
	dm, ok := m.(*DirectMode)
	if !ok {
		em := m.(*EffectMode)
		panic(fmt.Errorf("Active Mode (%s) is not a direct one", em.GetName()))
	}
	return dm
}
func (d *Device) MustGetDirectModeAndActivate() *DirectMode {
	m := d.GetActiveMode()
	// The current mode might be direct, in which case we just return it and don't worry that there might be other direct modes
	if dm, ok := m.(*DirectMode); ok {
		return dm
	}

	// Search for a direct mode, return the first
	for _, m := range d.Modes {
		if dm, ok := m.(*DirectMode); ok {
			d.log.V(1).Info("Found an arbitrary Direct mode, but there may be more than one", "mode", dm.GetName())
			return dm
		}
	}

	panic(fmt.Errorf("Device has no Direct Modes"))
}

func (d *Device) GetName() string {
	for _, dev := range d.model.Devices {
		if dev != d && dev.name == d.name {
			return fmt.Sprintf("%s[%d]", d.name, d.index)
		}
	}
	return d.name
}

func (d *Device) Size() int {
	m := d.GetActiveMode()
	if _, ok := m.(*DirectMode); !ok {
		d.log.V(1).Info("Getting size of active mode, but mode is not Direct; this may not be what you want.", "mode", m.GetName())
	}
	return m.Size()
}

func (d *Device) SetColor(c colorful.Color) {
	d.GetActiveMode().SetColor(c)
}

func (d *Device) SetColors(cs []colorful.Color) {
	m := d.GetActiveMode()
	if _, ok := m.(*DirectMode); !ok {
		d.log.V(1).Info("Setting color(s) of active mode, but mode is not Direct; this may not be what you want.", "mode", m.GetName())
	}
	m.SetColors(cs)
}

func (d *Device) Diff() {
	var path [4]string
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

func (d *Device) render(indent int) []indentedString {
	ss := []indentedString{
		{indent, fmt.Sprintf("DEVICE [%s] %s", d.devType, d.GetName())},
		{indent + 1, d.description},
		{indent + 1, fmt.Sprintf("Active mode: %s", d.newActiveMode.GetName())},
	}
	for _, mode := range d.Modes {
		ss = append(ss, mode.render(indent+1)...)
	}

	return ss
}
func (d *Device) String() string {
	return renderIndents(d.render(0))
}
