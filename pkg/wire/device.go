package wire

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
)

//go:generate stringer -type=DeviceType
// DeviceType models the various different types of product that an RGB Controller can be present in.
type DeviceType int

// Types are self-explainatory.
const (
	Motherboard  DeviceType = 0
	DIMM         DeviceType = 1
	GPU          DeviceType = 2
	Cooler       DeviceType = 3
	LEDStrip     DeviceType = 4
	Keyboard     DeviceType = 5
	Mouse        DeviceType = 6
	MouseMat     DeviceType = 7
	Headset      DeviceType = 8
	HeadsetStand DeviceType = 9
	Gamepad      DeviceType = 10
	Light        DeviceType = 11
	Speaker      DeviceType = 12
	Unknown      DeviceType = 13
)

// Device represents the natural idea of a device in a system, eg a RAM stick,
// graphics card, etc with RGB lighting. Technically it maps to a controller
// chip which can be sent commands over a bus.
type Device struct {
	Index         uint32
	Type          DeviceType
	Name          string
	Description   string
	Version       string
	Serial        string
	Location      string // Bus address eg /dev/i2c-0, address 0x3C
	ActiveModeIdx uint32
	Modes         []*Mode
	Zones         []*Zone
	LEDs          []*LED
	Colors        []colorful.Color
}

func extractDevice(buf []byte, idx uint32) *Device {
	d := &Device{Index: idx}
	offset := 0

	totalLen := extractUint32(buf, &offset)
	if len(buf) != int(totalLen) {
		panic(fmt.Sprintf("Assertion failed: msg len %d does not match length header %d", len(buf), totalLen))
	}

	d.Type = DeviceType(extractUint32(buf, &offset))

	d.Name = extractString(buf, &offset)
	d.Description = extractString(buf, &offset)
	d.Version = extractString(buf, &offset)
	d.Serial = extractString(buf, &offset)
	d.Location = extractString(buf, &offset)

	d.Modes, d.ActiveModeIdx = extractModes(buf, &offset)

	d.Zones = extractZones(buf, &offset)

	d.LEDs = extractLEDs(buf, &offset)

	d.Colors = extractColors(buf, &offset)

	return d
}
