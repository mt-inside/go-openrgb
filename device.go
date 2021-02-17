package main

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/lucasb-eyer/go-colorful"
)

func extractUint8(buf []byte, offset *int) uint8 {
	value := buf[*offset] // ignores bit-endinaness, as if "endinaness" only means byte-order. Seems to be ok, but might blow up when we talk to an ARM-build server
	*offset += int(reflect.TypeOf(value).Size())
	return value
}
func extractUint16(buf []byte, offset *int) uint16 {
	value := binary.LittleEndian.Uint16(buf[*offset:])
	*offset += int(reflect.TypeOf(value).Size())
	return value
}
func extractUint32(buf []byte, offset *int) uint32 {
	value := binary.LittleEndian.Uint32(buf[*offset:])
	*offset += int(reflect.TypeOf(value).Size())
	return value
}
func extractString(buf []byte, offset *int) (value string) {
	strLen := extractUint16(buf, offset)
	value = string(buf[*offset : *offset+int(strLen)-1]) // The strings, despite having length headers, also contain a null terminator, which we don't need
	*offset += int(strLen)
	return
}

// TODO think about the public API of all this.
// Ideall hide eg send on client (keep it all in the same package and you can do that, despite separate classes)
func FetchDevices(c *Client) ([]*Device, error) {
	deviceCount, err := fetchDeviceCount(c)
	if err != nil {
		return []*Device{}, fmt.Errorf("Couldn't fetch Device count: %w", err)
	}

	ds := make([]*Device, deviceCount)
	for i := 0; i < deviceCount; i++ {
		ds[i], err = fetchDevice(c, i)
		if err != nil {
			return []*Device{}, fmt.Errorf("Couldn't fetch Device %d: %w", i, err)
		}
	}

	return ds, nil
}

/* API is so shit.
* Everything has a length header.
* Devices' length is fetched by a separate command.
* Then each device by a command.
* Within a device, zones etc aren't API commands, they're packed into the binary blob, with thier length preceeding them */
func fetchDeviceCount(c *Client) (int, error) {
	if err := c.sendCommand(0, cmdGetDevCnt, []byte{}); err != nil {
		return 0, fmt.Errorf("Couldn't fetch Device count: %w", err)
	}

	body, err := c.readMessage()
	if err != nil {
		return 0, fmt.Errorf("Couldn't fetch Device count: %w", err)
	}

	deviceCount := int(binary.LittleEndian.Uint32(body))

	return deviceCount, nil

}
func fetchDevice(c *Client, i int) (*Device, error) {
	if err := c.sendCommand(uint32(i), cmdGetDevData, []byte{}); err != nil {
		return &Device{}, fmt.Errorf("Couldn't fetch Device %d: %w", i, err)
	}

	body, err := c.readMessage()
	if err != nil {
		return &Device{}, fmt.Errorf("Couldn't fetch Device %d: %w", i, err)
	}

	device := extractDevice(body)

	return device, nil
}

//go:generate stringer -type=DeviceType
type DeviceType int

const (
	Motherboard DeviceType = 0
	DIMM        DeviceType = 1
)

type Device struct {
	Type          DeviceType
	Name          string
	Description   string
	Version       string
	Serial        string
	Location      string
	ActiveModeIdx uint32            // TODO change to *Mode
	Modes         []*Mode           // TODO: what mean?
	Zones         []*Zone           // TODO: what mean?
	LEDs          []*LED            // TODO: what mean?
	Colors        []*colorful.Color // TODO: what mean?
}

// FIXME binary.Read() !!! Document - problem is the embedded variable-length strings. Use for modes / colours / LEDs ?

func extractDevice(buf []byte) *Device {
	d := &Device{}
	offset := 0

	foo := extractUint32(buf, &offset) // TODO what?
	fmt.Println("TODO foo: ", foo)

	d.Type = DeviceType(extractUint32(buf, &offset))

	d.Name = extractString(buf, &offset)
	d.Description = extractString(buf, &offset)
	d.Version = extractString(buf, &offset)
	d.Serial = extractString(buf, &offset)
	d.Location = extractString(buf, &offset)

	//TODO
	d.Modes, d.ActiveModeIdx = extractModes(buf, &offset)

	//TODO
	d.Zones = extractZones(buf, &offset)

	//TODO
	d.LEDs = extractLEDs(buf, &offset)

	//TODO
	d.Colors = extractColors(buf, &offset)

	return d
}

//go:generate stringer -type=ColorMode
type ColorMode int

const (
	None         ColorMode = 0
	PerLED       ColorMode = 1
	ModeSpecific ColorMode = 2
	Random       ColorMode = 3
)

type Mode struct {
	Name      string
	Value     uint32
	Flags     uint32
	MinSpeed  uint32
	MaxSpeed  uint32
	MinColors uint32
	MaxColors uint32
	Speed     uint32
	Direction uint32
	ColorMode ColorMode
	Colors    []*colorful.Color
}

func extractModes(buf []byte, offset *int) (modes []*Mode, activeModeIdx uint32) {
	modeCount := int(extractUint16(buf, offset))
	activeM := extractUint32(buf, offset)

	ms := make([]*Mode, modeCount)
	for i := 0; i < modeCount; i++ {
		ms[i] = extractMode(buf, offset)
	}

	return ms, activeM
}

func extractMode(buf []byte, offset *int) *Mode {
	m := &Mode{}

	m.Name = extractString(buf, offset)

	m.Value = extractUint32(buf, offset)
	m.Flags = extractUint32(buf, offset)
	m.MinSpeed = extractUint32(buf, offset)
	m.MaxSpeed = extractUint32(buf, offset)
	m.MinColors = extractUint32(buf, offset)
	m.MaxColors = extractUint32(buf, offset)
	m.Speed = extractUint32(buf, offset)
	m.Direction = extractUint32(buf, offset)
	m.ColorMode = ColorMode(extractUint32(buf, offset))

	m.Colors = extractColors(buf, offset) // FIXME

	return m
}

func extractColors(buf []byte, offset *int) []*colorful.Color {
	colorCount := int(extractUint16(buf, offset))

	cs := make([]*colorful.Color, colorCount)
	for i := 0; i < colorCount; i++ {
		cs[i] = extractColor(buf, offset)
	}

	return cs
}

func extractColor(buf []byte, offset *int) *colorful.Color {
	r := extractUint8(buf, offset)
	g := extractUint8(buf, offset)
	b := extractUint8(buf, offset)
	*offset += 1 // Colors are padded to 4 bytes

	return &colorful.Color{float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0}
}

//go:generate stringer -type=ZoneType
type ZoneType uint32

const (
	Zero ZoneType = 0
	One  ZoneType = 1
)

type Zone struct {
	Name      string
	Type      ZoneType
	MinLEDs   uint32
	MaxLEDs   uint32
	TotalLEDs uint32
}

func extractZones(buf []byte, offset *int) []*Zone {
	zoneCount := int(extractUint16(buf, offset))

	zs := make([]*Zone, zoneCount)
	for i := 0; i < zoneCount; i++ {
		zs[i] = extractZone(buf, offset)
	}

	return zs
}

func extractZone(buf []byte, offset *int) *Zone {
	z := &Zone{}

	z.Name = extractString(buf, offset)

	z.Type = ZoneType(extractUint32(buf, offset))
	z.MinLEDs = extractUint32(buf, offset)
	z.MaxLEDs = extractUint32(buf, offset)
	z.TotalLEDs = extractUint32(buf, offset)

	matrixSize := extractUint16(buf, offset)
	*offset += int(matrixSize) // TODO what?
	fmt.Println("TODO matrix size: ", matrixSize)

	return z
}

type LED struct {
	Name  string
	Color *colorful.Color
}

func extractLEDs(buf []byte, offset *int) []*LED {
	ledsCount := int(extractUint16(buf, offset))

	ls := make([]*LED, ledsCount)
	for i := 0; i < ledsCount; i++ {
		ls[i] = extractLED(buf, offset)
	}

	return ls
}

func extractLED(buf []byte, offset *int) *LED {
	l := &LED{}

	l.Name = extractString(buf, offset)
	l.Color = extractColor(buf, offset)

	return l
}
