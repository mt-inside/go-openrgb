package wire

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/lucasb-eyer/go-colorful"
)

const knownProtoVer = 1

//go:generate stringer -type=Command
type Command int

const (
	// Devices
	cmdGetDevCnt  Command = 0
	cmdGetDevData Command = 1
	// Protocol
	cmdGetProtocolVersion Command = 40 //nolint:varcheck,deadcode,unused
	cmdSetClientName      Command = 50
	// Streaming support?
	queDevListUpdated Command = 100 //nolint:varcheck,deadcode,unused
	// Manipulation of profiles
	cmdGetProfileList Command = 150 //nolint:varcheck,deadcode,unused
	cmdSaveProfile    Command = 151 //nolint:varcheck,deadcode,unused
	cmdLoadProfile    Command = 152 //nolint:varcheck,deadcode,unused
	cmdDeleteProfile  Command = 153 //nolint:varcheck,deadcode,unused

	// Setting colors
	cmdResizeZone        Command = 1000 //nolint:varcheck,deadcode,unused
	cmdUpdateLEDs        Command = 1050 //nolint:varcheck,deadcode,unused
	cmdUpdateZoneLEDs    Command = 1051 //nolint:varcheck,deadcode,unused
	cmdUpdateSingularLED Command = 1052 //nolint:varcheck,deadcode,unused
	cmdSetCustomMode     Command = 1100 //nolint:varcheck,deadcode,unused
	cmdUpdateMode        Command = 1101 //nolint:varcheck,deadcode,unused
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

func extractColors(buf []byte, offset *int) []colorful.Color {
	colorCount := extractUint16(buf, offset)

	cs := make([]colorful.Color, colorCount)
	for i := uint16(0); i < colorCount; i++ {
		cs[i] = extractColor(buf, offset)
	}

	return cs
}

func extractColor(buf []byte, offset *int) colorful.Color {
	r := extractUint8(buf, offset)
	g := extractUint8(buf, offset)
	b := extractUint8(buf, offset)
	*offset += 1 // Colors are padded to 4 bytes

	return colorful.Color{R: float64(r) / 255.0, G: float64(g) / 255.0, B: float64(b) / 255.0}
}

func extractString(buf []byte, offset *int) (value string) {
	strLen := extractUint16(buf, offset)
	value = string(buf[*offset : *offset+int(strLen)-1]) // The strings, despite having length headers, also contain a null terminator, which we don't need
	*offset += int(strLen)
	return
}

//nolint:deadcode,unused
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
	insertUint8(buf, offset, uint8(0)) // Colors are padded to 4 bytes
}

func insertString(buf []byte, offset *int, value string) {
	origOffset := *offset
	insertUint16(buf, offset, uint16(len(value)))
	copy(buf[*offset:*offset+len(value)], value) // Assumes UTF8
	*offset += len(value)                        // len(str) is byte len, not rune len
	insertUint8(buf, offset, 0)                  // The strings, despite having length headers, also contain a null terminator, which we don't need
	if *offset != origOffset+2+len(value)+1 {
		panic(fmt.Errorf("Assertion failed: offset=%d, calculated=%d", *offset, 2+len(value)+1))
	}
}
