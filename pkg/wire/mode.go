package wire

import (
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

//go:generate stringer -type=ColorMode
type ColorMode int

const (
	None         ColorMode = 0
	PerLED       ColorMode = 1
	ModeSpecific ColorMode = 2
	Random       ColorMode = 3
)

//go:generate stringer -type=ModeDirection
type ModeDirection int

const (
	Left       ModeDirection = 0
	Right      ModeDirection = 1
	Up         ModeDirection = 2
	Down       ModeDirection = 3
	Horizontal ModeDirection = 4
	Vertical   ModeDirection = 5
)

type ModeFlags int

const (
	HasSpeed             ModeFlags = 1 << 0
	HasDirectionLR       ModeFlags = 1 << 1
	HasDirectionUD       ModeFlags = 1 << 2
	HasDirectionHV       ModeFlags = 1 << 3
	HasBrightness        ModeFlags = 1 << 4
	HasPerLEDColor       ModeFlags = 1 << 5
	HasModeSpecificColor ModeFlags = 1 << 6
	HasRandomColor       ModeFlags = 1 << 7
)

var ModeFlag_names = []string{
	"HasSpeed",
	"HasDirectionLR",
	"HasDirectionUD",
	"HasDirectionHV",
	"HasBrightness",
	"HasPerLEDColor",
	"HasModeSpecificColor",
	"HasRandomColor",
}

// Can't find a library or generator for this
func (f ModeFlags) String() string {
	flagNames := []string{}

	for i := 0; i < len(ModeFlag_names); i++ {
		if (uint32(f) & (1 << i)) != 0 { // Cast is neccessary otherwise x&y is always false...
			flagNames = append(flagNames, ModeFlag_names[i])
		}
	}

	if len(flagNames) == 0 {
		return "0"
	}
	return strings.Join(flagNames, "|")
}

type Mode struct {
	Index uint16
	Name  string

	Value uint32 // driver-internal

	MinSpeed uint32
	Speed    uint32 // speed of the effect
	MaxSpeed uint32

	Direction ModeDirection // direction of the effect

	Flags     ModeFlags // available ColorModes, Directions, and other attributes
	ColorMode ColorMode

	MinColors uint32 // min!=max => user-resizable? ie you can flash between 1, 2, 3, etc different colours
	MaxColors uint32

	Colors []*colorful.Color
}

func extractModes(buf []byte, offset *int) (modes []*Mode, activeModeIdx uint32) {
	modeCount := extractUint16(buf, offset)
	activeM := extractUint32(buf, offset)

	ms := make([]*Mode, modeCount)
	for i := uint16(0); i < modeCount; i++ {
		ms[i] = extractMode(buf, offset, i)
	}

	return ms, activeM
}

func extractMode(buf []byte, offset *int, idx uint16) *Mode {
	m := &Mode{Index: idx}

	m.Name = extractString(buf, offset)

	m.Value = extractUint32(buf, offset)
	m.Flags = ModeFlags(extractUint32(buf, offset))
	m.MinSpeed = extractUint32(buf, offset)
	m.MaxSpeed = extractUint32(buf, offset)
	m.MinColors = extractUint32(buf, offset)
	m.MaxColors = extractUint32(buf, offset)
	m.Speed = extractUint32(buf, offset)
	m.Direction = ModeDirection(extractUint32(buf, offset))
	m.ColorMode = ColorMode(extractUint32(buf, offset))

	m.Colors = extractColors(buf, offset)

	return m
}
