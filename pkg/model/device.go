package model

import (
	"fmt"

	"github.com/mt-inside/go-openrgb/pkg/wire"
)

type Device struct {
	//index         uint32
	devType       wire.DeviceType
	name          string
	description   string
	version       string
	serial        string
	location      string
	activeModeIdx uint32 // TODO hide this, add functions for SetActiceMode(*Mode - got with ByName or whatever), GetActiceMode() *Mode
	Modes         []Mode
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
