package wire

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
)

type WriteDevice struct {
	Index  uint32
	Colors []colorful.Color
}

func WriteDevices(c *Client, wds []*WriteDevice) error {
	for _, wd := range wds {
		err := writeDevice(c, wd)
		if err != nil {
			return fmt.Errorf("Couldn't sync devices to server: %w", err)
		}
	}

	return nil
}

func writeDevice(c *Client, wd *WriteDevice) error {
	colorsLen := len(wd.Colors)
	bufLen := 4 + 2 + (colorsLen * 4) // TODO this using sizeof, but not reflect?
	buf := make([]byte, bufLen)

	offset := 0

	insertUint32(buf, &offset, uint32(bufLen))
	insertUint16(buf, &offset, uint16(colorsLen))
	for _, color := range wd.Colors {
		insertColor(buf, &offset, color)
	}

	err := c.sendCommand(wd.Index, cmdUpdateLEDs, buf)
	if err != nil {
		return fmt.Errorf("Can't set device-level LEDs: %w", err)
	}

	return nil
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
