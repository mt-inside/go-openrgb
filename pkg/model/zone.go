package model

import (
	"fmt"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mt-inside/go-openrgb/pkg/wire"
)

type ZoneList []*Zone

/* OpenRGB models every device as having at least one zone; making a zone with a generic name like "DRAM" for devices (like DIMMs) that don't really have them.
* This function just gets the one and only zone, where applicable */
func (zs ZoneList) MustSingle() *Zone {
	if len(zs) != 1 {
		panic("Not precisely one zone")
	}
	return zs[0]
}
func (zs ZoneList) ByName(name string) (*Zone, bool) {
	for _, z := range zs {
		if strings.EqualFold(z.GetName(), name) {
			return z, true
		}
	}
	return nil, false
}
func (zs ZoneList) MustByName(name string) *Zone {
	z, ok := zs.ByName(name)
	if !ok {
		panic(fmt.Errorf("Zone list doesn't contain: %s", name))
	}
	return z
}

type Zone struct {
	mode         Mode
	name         string
	zoneType     wire.ZoneType
	minLEDs      uint32 // min!=max => user-resizable (depending on what's plugged in)
	maxLEDs      uint32
	matrixWidth  uint32
	matrixHeight uint32
	Leds         []*LED
}

func (z *Zone) GetName() string {
	// never seen zones clash names, plus they have no wire index
	return z.name
}
func (z *Zone) Size() int {
	return len(z.Leds)
}

func (z *Zone) SetColor(c colorful.Color) {
	for _, l := range z.Leds {
		l.SetColor(c)
	}
}
func (z *Zone) SetColors(cs []colorful.Color) {
	if z.Size() != len(cs) {
		panic(fmt.Errorf("trying to set %d-led Zone with %d colors", len(z.Leds), len(cs)))
	}
	for i := range z.Leds {
		z.Leds[i].SetColor(cs[i])
	}
}

func (z *Zone) Diff() {
	var path [2]string
	path[0] = z.GetName()
	for _, l := range z.Leds {
		path[1] = l.GetName()
		if l.serverColor != l.newColor {
			fmt.Printf("%s: %s -> %s\n", strings.Join(path[:], "/"), l.serverColor.Hex(), l.newColor.Hex())
		}
	}
}

func (z *Zone) render(indent int) []indentedString {
	ss := []indentedString{
		{indent, fmt.Sprintf("ZONE [%s] %s", z.zoneType, z.GetName())},
	}
	if z.zoneType == wire.Planar {
		ss = append(ss, indentedString{indent + 1, fmt.Sprintf("Matrix %dx%d", z.matrixWidth, z.matrixHeight)})
	}
	ss = append(ss, indentedString{indent + 1, fmt.Sprintf("LEDs: (%d-%d) %d:%s", z.minLEDs, z.maxLEDs, len(z.Leds), renderLedColors(z.Leds))})

	return ss

}
func (z *Zone) String() string {
	return renderIndents(z.render(0))
}
