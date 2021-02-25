package model

import (
	"fmt"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mt-inside/go-openrgb/pkg/wire"
)

type ZoneList []*Zone

func (zs ZoneList) ByName(name string) (*Zone, bool) {
	for _, z := range zs {
		if strings.EqualFold(z.name, name) {
			return z, true
		}
	}
	return nil, false
}
func (zs ZoneList) ByNameUnwrap(name string) *Zone {
	z, ok := zs.ByName(name)
	if !ok {
		panic(fmt.Errorf("Zone list doesn't contain: %s", name))
	}
	return z
}

type Zone struct {
	name     string
	zoneType wire.ZoneType
	minLEDs  uint32 // min!=max => user-resizable (depending on what's plugged in)
	maxLEDs  uint32
	Leds     []*LED
}

func (z *Zone) SetColor(c colorful.Color) {
	for _, l := range z.Leds {
		l.SetColor(c)
	}
}

func (z *Zone) render(indent int) []indentedString {
	// We skip a level; not rendering LED names.
	// Thus this is map()
	colors := []colorful.Color{}
	for _, led := range z.Leds {
		colors = append(colors, led.newColor)
	}

	ss := []indentedString{
		{indent, fmt.Sprintf("ZONE [%s] %s", z.zoneType, z.name)},
		{indent + 1, fmt.Sprintf("LEDs: (%d-%d) %d:%s", z.minLEDs, z.maxLEDs, len(z.Leds), renderColors(colors))},
	}

	return ss

}
func (z *Zone) String() string {
	return renderIndents(z.render(0))
}
