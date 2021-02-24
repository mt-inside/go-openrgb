package wire

import "github.com/lucasb-eyer/go-colorful"

// TODO refactor - should probably be on client?
//nolint:deadcode,unused
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

// This is for the subset of device LEDs in a zone, ie for direct mode.
//nolint:deadcode,unused
func getCommandZoneLEDs(zoneID uint32, colors []colorful.Color) []byte {
	colorsLen := len(colors)
	bufLen := 4 + 4 + 2 + (colorsLen * 4) // TODO this using sizeof, but not reflect?
	buf := make([]byte, bufLen)

	offset := 0

	insertUint32(buf, &offset, uint32(bufLen))
	insertUint32(buf, &offset, zoneID)
	insertUint16(buf, &offset, uint16(colorsLen))
	for _, color := range colors {
		insertColor(buf, &offset, color)
	}

	return buf
}
