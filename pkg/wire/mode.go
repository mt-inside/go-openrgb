package wire

import (
	"fmt"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

//go:generate stringer -type=ColorMode
// ColorMode specifies how colors are chosen in this mode.
// Modes may support more than one color mode.
type ColorMode int

const (
	// ColorModeNone means a "color" for this mode wouldn't make sense, eg
	// it's a rainbow effect
	ColorModeNone ColorMode = 0
	// PerLED allows user to chose a color for each individual LED, ie is a
	// "direct" mode
	PerLED ColorMode = 1
	// ModeSpecific means the mode has one (or sometimes a few) colors,
	// which the user can chose. Whereas PerLED means addressing individual
	// LEDs, this would be eg a red static or breathing effect across the
	// whole device/zone. Where >1 color is supported, that means eg
	// flashing the whole device between red and green; it's not the same
	// as per-LED color.
	ModeSpecific ColorMode = 2
	// Random means the mode has one (or a few) colors, but user wants the
	// controller to chose them at random. This isn't the same as None -
	// None means specifying a color is meaningless, eg the device is doing
	// a rainbow. Random means there is a color, but it's chosen at random,
	// eg static at a random color, or breathing with a different color
	// each time.
	Random ColorMode = 3
)

//go:generate stringer -type=EffectDirection
// EffectDirection is the direction of travel of the effect across the device.
type EffectDirection int

// Valid effect directions.
// Technically the Direction field should be ignored if none of the
// HasDirection* flags are set, but I've added a DirectionNA value to use in
// that case, in case the user doesn't check the flags (as 0 is used
// legitimately by the protocol). TODO represent this optionality better, eg
// with a getter func or a pointer.
const (
	Left        EffectDirection = 0
	Right       EffectDirection = 1
	Up          EffectDirection = 2
	Down        EffectDirection = 3
	Horizontal  EffectDirection = 4
	Vertical    EffectDirection = 5
	DirectionNA EffectDirection = 65535
)

// ModeFlags says which attributes this mode has, and thus which fields are applicable.
// The OpenRGB server sends arbitrary data (uninitialised memory?) down the
// wire for fields which, according to these flags, aren't applicable.
type ModeFlags int

// Mode flags:
const (
	HasSpeed             ModeFlags = 1 << 0 // Speed field is in use
	HasDirectionLR       ModeFlags = 1 << 1 // Direction field is in use, and can contain Left or Right
	HasDirectionUD       ModeFlags = 1 << 2 // Direction field is in use, and can contain Up or Down
	HasDirectionHV       ModeFlags = 1 << 3 // Direction field is in use, and can contain Horizontal or Vertical
	HasBrightness        ModeFlags = 1 << 4 // TODO doesn't make any sense because there's not a brightness field (nor a brightness control in the GUI)
	HasPerLEDColor       ModeFlags = 1 << 5 // PerLED is an acceptable value for ColorMode
	HasModeSpecificColor ModeFlags = 1 << 6 // ModeSpecific is an acceptable value for ColorMode
	HasRandomColor       ModeFlags = 1 << 7 // Random is an acceptable value for ColorMode
)

var modeFlagNames = []string{
	"HasSpeed",
	"HasDirectionLR",
	"HasDirectionUD",
	"HasDirectionHV",
	"HasBrightness", // But no brightness field?
	"HasPerLEDColor",
	"HasModeSpecificColor",
	"HasRandomColor",
}

// ModeDirectionFlagsClear returns true if none of the three Direction* flags are set in the given flags bitmap.
// This is not a method if the ModeFlags type, as the Flag values are used outside this package
func ModeDirectionFlagsClear(f ModeFlags) bool {
	return f&(HasDirectionLR|HasDirectionUD|HasDirectionHV) == 0
}

// Can't find a library or generator for this
func (f ModeFlags) String() string {
	if uint32(f) == 0 {
		return "<none>"
	}
	flagNames := []string{}

	for i := 0; i < len(modeFlagNames); i++ {
		if (uint32(f) & (1 << i)) != 0 { // Cast is necessary otherwise x&y is always false...
			flagNames = append(flagNames, modeFlagNames[i])
		}
	}

	return strings.Join(flagNames, ",")
}

// Mode is an "effect" that a controller can render to a set of LEDs.
// This is a big topic, but basically the mode will be something like "static"
// or "breathing", and has various attributes like its direction, speed, and
// color(s).
type Mode struct {
	Index uint32
	Name  string

	Value uint32 // driver-internal

	MinSpeed uint32
	Speed    uint32 // speed of the effect
	MaxSpeed uint32

	Direction EffectDirection // direction of the effect

	Flags     ModeFlags // available ColorModes, Directions, and other attributes
	ColorMode ColorMode

	MinColors uint32 // min!=max => user-resizable? ie you can flash between 1, 2, 3, etc different colours
	MaxColors uint32

	Colors []colorful.Color
}

func extractModes(buf []byte, offset *int) (modes []*Mode, activeModeIdx uint32) {
	/* Confirmed from OpenRGB source:
	 * - number of modes is a `short`: https://gitlab.com/CalcProgrammer1/OpenRGB/-/blob/master/RGBController/RGBController.cpp#L510
	 * - modes don't have an index when rendered into a Device; you have to track it
	 * - but when rendered themselves, Modes' indecies are `int`:  https://gitlab.com/CalcProgrammer1/OpenRGB/-/blob/master/RGBController/RGBController.cpp#L792
	 */
	modeCount := uint32(extractUint16(buf, offset))
	activeM := extractUint32(buf, offset)

	ms := make([]*Mode, modeCount)
	for i := uint32(0); i < modeCount; i++ {
		ms[i] = extractMode(buf, offset, i)
	}

	return ms, activeM
}

func extractMode(buf []byte, offset *int, idx uint32) *Mode {
	m := &Mode{Index: idx}

	m.Name = extractString(buf, offset)

	m.Value = extractUint32(buf, offset)
	m.Flags = ModeFlags(extractUint32(buf, offset))
	m.MinSpeed = extractUint32(buf, offset)
	m.MaxSpeed = extractUint32(buf, offset)
	m.MinColors = extractUint32(buf, offset)
	m.MaxColors = extractUint32(buf, offset)
	m.Speed = extractUint32(buf, offset)
	m.Direction = EffectDirection(extractUint32(buf, offset))
	m.ColorMode = ColorMode(extractUint32(buf, offset))

	if m.Flags&HasSpeed != HasSpeed {
		m.MinSpeed = 0 // Uninitialised if flag is clear (nasty cause 0 is a valid Speed)
		m.Speed = 0
		m.MaxSpeed = 0
	}
	if ModeDirectionFlagsClear(m.Flags) {
		m.Direction = DirectionNA
	}
	if m.ColorMode == ColorModeNone {
		m.MinColors = 0 // These seem uninitialised if Color Mode is None
		m.MaxColors = 0
	}

	m.Colors = extractColors(buf, offset)

	if m.ColorMode == ColorModeNone && len(m.Colors) != 0 {
		panic(fmt.Errorf("Assertion failed: Assumed color mode 'none' would have zero colors"))
	}

	return m
}
