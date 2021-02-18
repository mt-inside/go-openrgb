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
func extractColor(buf []byte, offset *int) *colorful.Color {
	r := extractUint8(buf, offset)
	g := extractUint8(buf, offset)
	b := extractUint8(buf, offset)
	*offset += 1 // Colors are padded to 4 bytes

	return &colorful.Color{R: float64(r) / 255.0, G: float64(g) / 255.0, B: float64(b) / 255.0}
}

func extractString(buf []byte, offset *int) (value string) {
	strLen := extractUint16(buf, offset)
	value = string(buf[*offset : *offset+int(strLen)-1]) // The strings, despite having length headers, also contain a null terminator, which we don't need
	*offset += int(strLen)
	return
}

func insertUint8(buf []byte, offset *int, value uint8) {
	buf[*offset] = value // similar endinaness concerns to extract8
	*offset += int(reflect.TypeOf(value).Size())
}
func insertUint16(buf []byte, offset *int, value uint16) {
	binary.LittleEndian.PutUint16(buf[*offset:], value)
	*offset += int(reflect.TypeOf(value).Size())
}
func insertUint32(buf []byte, offset *int, value uint32) {
	binary.LittleEndian.PutUint32(buf[*offset:], value)
	*offset += int(reflect.TypeOf(value).Size())
}
func insertColor(buf []byte, offset *int, value colorful.Color) {
	insertUint8(buf, offset, uint8(value.R*255.0))
	insertUint8(buf, offset, uint8(value.G*255.0))
	insertUint8(buf, offset, uint8(value.B*255.0))
	*offset += 1 // Colors are padded to 4 bytes
}

// TODO think about the public API of all this.
// Ideall hide eg send on client (keep it all in the same package and you can do that, despite separate classes)
func FetchDevices(c *Client) ([]*Device, error) {
	deviceCount, err := fetchDeviceCount(c)
	if err != nil {
		return []*Device{}, fmt.Errorf("Couldn't fetch Device count: %w", err)
	}

	ds := make([]*Device, deviceCount)
	for i := uint32(0); i < deviceCount; i++ {
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
func fetchDeviceCount(c *Client) (uint32, error) {
	if err := c.sendCommand(0, cmdGetDevCnt, []byte{}); err != nil {
		return 0, fmt.Errorf("Couldn't fetch Device count: %w", err)
	}

	body, err := c.readMessage()
	if err != nil {
		return 0, fmt.Errorf("Couldn't fetch Device count: %w", err)
	}

	deviceCount := binary.LittleEndian.Uint32(body)

	return deviceCount, nil

}
func fetchDevice(c *Client, i uint32) (*Device, error) {
	if err := c.sendCommand(uint32(i), cmdGetDevData, []byte{}); err != nil {
		return &Device{}, fmt.Errorf("Couldn't fetch Device %d: %w", i, err)
	}

	body, err := c.readMessage()
	if err != nil {
		return &Device{}, fmt.Errorf("Couldn't fetch Device %d: %w", i, err)
	}

	device := extractDevice(body, i)

	return device, nil
}

//go:generate stringer -type=DeviceType
type DeviceType int

const (
	Motherboard DeviceType = 0
	DIMM        DeviceType = 1
)

type Device struct {
	Index         uint32
	Type          DeviceType
	Name          string
	Description   string
	Version       string
	Serial        string
	Location      string
	ActiveModeIdx uint32 // TODO hide this, add functions for SetActiceMode(*Mode - got with ByName or whatever), GetActiceMode() *Mode
	Modes         []*Mode
	Zones         []*Zone           // TODO move under all the directs.
	LEDs          []*LED            // TODO move under Zone
	Colors        []*colorful.Color // TODO move under LED
}

// binary.Read() is neat, but every type (except for Color) has a headed string at the front, so that wouldn't work. Also requires construction of a Reader, which might be slow.

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

//go:generate stringer -type=ColorMode
type ColorMode int

const (
	None         ColorMode = 0
	PerLED       ColorMode = 1
	ModeSpecific ColorMode = 2
	Random       ColorMode = 3
)

type Mode struct {
	Name string

	Value uint32 // driver-internal
	Flags uint32 // useful for the app, but idk what's in it

	MinSpeed uint32
	Speed    uint32 // speed of the effect
	MaxSpeed uint32

	Direction uint32 // direction of the effect

	ColorMode ColorMode

	MinColors uint32 // min!=max => user-resizable? ie you can flash between 1, 2, 3, etc different colours
	MaxColors uint32

	Colors []*colorful.Color
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

	m.Colors = extractColors(buf, offset)

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

//go:generate stringer -type=ZoneType
type ZoneType uint32

const (
	Zero ZoneType = 0
	One  ZoneType = 1
)

type Zone struct {
	Name      string
	Type      ZoneType
	MinLEDs   uint32 // min!=max => user-resizable (depending on what's plugged in)
	MaxLEDs   uint32
	TotalLEDs uint32 // current size?
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
	Value uint32 // driver-internal, eg LED mapping
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
	l.Value = extractUint32(buf, offset)

	return l
}

// TODO refactor - should probably be on client?
func getCommandLEDs(colors []colorful.Color) []byte {
	colorsLen := len(colors)
	bufLen := 4 + 2 + (colorsLen * 4) // TODO this using sizeof, but not reflect?
	buf := make([]byte, bufLen)

	offset := 0

	insertUint32(buf, &offset, uint32(bufLen))
	insertUint16(buf, &offset, uint16(colorsLen))
	for _, color := range colors {
		insertColor(buf, &offset, color)
	}

	return buf
}
