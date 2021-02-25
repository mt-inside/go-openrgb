package model

import (
	"fmt"
	"strings"

	"github.com/mt-inside/go-openrgb/pkg/wire"
)

type DeviceList []*Device

func (ds DeviceList) ByName(name string) (*Device, bool) {
	for _, d := range ds {
		if strings.EqualFold(d.name, name) { // case insensitive
			return d, true
		}
	}
	return nil, false
}
func (ds DeviceList) ByNameUnwrap(name string) *Device {
	d, ok := ds.ByName(name)
	if !ok {
		panic(fmt.Errorf("Device list doesn't contain: %s", name))
	}
	return d
}

type Device struct {
	index         uint32
	devType       wire.DeviceType
	name          string
	description   string
	version       string
	serial        string
	location      string
	activeModeIdx uint32 // TODO hide this, add functions for SetActiceMode(*Mode - got with ByName or whatever), GetActiceMode() *Mode
	Modes         ModeList
}

func (d *Device) render(indent int) []indentedString {
	ss := []indentedString{
		{indent, fmt.Sprintf("DEVICE [%s] %s", d.devType, d.name)},
		{indent + 1, d.description},
	}
	for _, mode := range d.Modes {
		ss = append(ss, mode.render(indent+1)...)
	}

	return ss
}
func (d *Device) String() string {
	return renderIndents(d.render(0))
}
