package model

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mt-inside/go-openrgb/pkg/wire"
)

type Zone struct {
	name     string
	zoneType wire.ZoneType
	minLEDs  uint32 // min!=max => user-resizable (depending on what's plugged in)
	maxLEDs  uint32
	Leds     []*LED
}

func (z *Zone) render(indent int) []indentedString {
	// We skip a level; not rendering LED names.
	// Thus this is map()
	colors := []*colorful.Color{}
	for _, led := range z.Leds {
		colors = append(colors, led.color)
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
