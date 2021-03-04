package model

import (
	"fmt"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mt-inside/go-openrgb/pkg/wire"
)

type ModeList []Mode

func (ms ModeList) Directs() []*DirectMode {
	dms := []*DirectMode{}

	for _, m := range ms {
		if dm, ok := m.(*DirectMode); ok {
			dms = append(dms, dm)
		}
	}

	return dms
}

type Mode interface {
	Settable
	getIndex() uint32
	getWireName() string
	getWireMode() *wire.Mode
	render(indent int) []indentedString // TODO remove and see if anyone misses it
}

type EffectMode struct {
	device   *Device
	name     string
	Colors   []*LED // To make patches work sensibly, might have to inject a fake "all" zone in here. The server does this for DirectModes when the device has no zones
	wireMode *wire.Mode
}

func (m *EffectMode) getIndex() uint32 {
	return m.wireMode.Index
}
func (m *EffectMode) getWireName() string {
	return m.wireMode.Name
}
func (m *EffectMode) getWireMode() *wire.Mode {
	return m.wireMode
}

func (m *EffectMode) GetName() string {
	for _, mode := range m.device.Modes {
		if mode != m && mode.getWireName() == m.name {
			return fmt.Sprintf("%s[%d]", m.name, m.wireMode.Index)
		}
	}
	return m.name
}

func (m *EffectMode) Size() int {
	return len(m.Colors)
}

func (m *EffectMode) SetColor(c colorful.Color) {
	for i := 0; i < len(m.Colors); i++ {
		m.Colors[i].SetColor(c)
	}
}

func (m *EffectMode) SetColors(cs []colorful.Color) {
	if m.Size() != len(cs) {
		panic(fmt.Errorf("Trying to set %d-color Mode with %d colors.", m.Size(), len(cs)))
	}
	for i := range m.Colors {
		m.Colors[i].SetColor(cs[i])
	}
}

func (m *EffectMode) Diff() {
	var path [2]string
	path[0] = m.GetName()
	for _, l := range m.Colors {
		path[1] = l.GetName()
		if l.serverColor != l.newColor {
			fmt.Printf("%s: %s -> %s\n", strings.Join(path[:], "/"), l.serverColor.Hex(), l.newColor.Hex())
		}
	}
}

func (m *EffectMode) render(indent int) []indentedString {
	attrs := []string{}
	if !wire.ModeDirectionFlagsClear(m.wireMode.Flags) {
		attrs = append(attrs, fmt.Sprintf("direction: %s", m.wireMode.Direction))
	}
	if m.wireMode.Flags&wire.HasSpeed == wire.HasSpeed {
		attrs = append(attrs, fmt.Sprintf("speed: %d (%d-%d)", m.wireMode.Speed, m.wireMode.MinSpeed, m.wireMode.MaxSpeed))
	}

	activeFlag := ""
	if m.device.newActiveMode == m {
		activeFlag = " !ACTIVE!"
	}
	ss := []indentedString{
		{indent, fmt.Sprintf("MODE [Effect] %s%s (%s)", m.GetName(), activeFlag, m.wireMode.Flags)},
		{indent + 1, strings.Join(attrs, ",")},
	}
	if m.wireMode.ColorMode == wire.ColorModeNone {
		ss = append(ss, indentedString{indent + 1, fmt.Sprintf("Colors: %s", m.wireMode.ColorMode)})
	} else {
		ss = append(ss, indentedString{indent + 1, fmt.Sprintf("Colors: %s, (%d-%d) %s", m.wireMode.ColorMode, m.wireMode.MinColors, m.wireMode.MaxColors, renderLedColors(m.Colors))})
	}

	return ss
}
func (m *EffectMode) String() string {
	return renderIndents(m.render(0))
}

type DirectMode struct {
	device   *Device
	name     string
	Zones    ZoneList
	wireMode *wire.Mode
}

func (m *DirectMode) getIndex() uint32 {
	return m.wireMode.Index
}
func (m *DirectMode) getWireName() string {
	return m.wireMode.Name
}
func (m *DirectMode) getWireMode() *wire.Mode {
	return m.wireMode
}

func (m *DirectMode) GetName() string {
	for _, mode := range m.device.Modes {
		if mode != m && mode.getWireName() == m.name {
			return fmt.Sprintf("%s[%d]", m.name, m.wireMode.Index)
		}
	}
	return m.name
}

func (m *DirectMode) Size() int {
	size := 0
	for _, z := range m.Zones {
		size += z.Size()
	}
	return size
}

func (m *DirectMode) SetColor(c colorful.Color) {
	for _, z := range m.Zones {
		z.SetColor(c)
	}
}

func (m *DirectMode) SetColors(cs []colorful.Color) {
	if m.Size() != len(cs) {
		panic(fmt.Errorf("Trying to set %d-led Mode with %d colors.", m.Size(), len(cs)))
	}
	i := 0
	for _, z := range m.Zones {
		z.SetColors(cs[i : i+z.Size()])
		i += z.Size()
	}
}

func (m *DirectMode) Diff() {
	var path [3]string
	path[0] = m.GetName()
	for _, z := range m.Zones {
		path[1] = z.GetName()
		for _, l := range z.Leds {
			path[2] = l.GetName()
			if l.serverColor != l.newColor {
				fmt.Printf("%s: %s -> %s\n", strings.Join(path[:], "/"), l.serverColor.Hex(), l.newColor.Hex())
			}
		}
	}
}

func (m *DirectMode) render(indent int) []indentedString {
	activeFlag := ""
	if m.device.newActiveMode == m {
		activeFlag = " !ACTIVE!"
	}
	ss := []indentedString{
		{indent, fmt.Sprintf("MODE [Direct] %s%s", m.GetName(), activeFlag)},
	}

	for _, zone := range m.Zones {
		ss = append(ss, zone.render(indent+1)...)
	}

	return ss
}
func (m *DirectMode) String() string {
	return renderIndents(m.render(0))
}
