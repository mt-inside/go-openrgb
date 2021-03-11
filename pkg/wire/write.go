package wire

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
)

// UpdateLEDs is the struct used on the wire for the payload to an UpdateLEDs
// command
type UpdateLEDs struct {
	DeviceID uint32
	Colors   []colorful.Color
}

// SendUpdateLEDs sends a (whole-device) UpdateLEDs command.
func SendUpdateLEDs(c *Client, wd *UpdateLEDs) error {
	colorsLen := len(wd.Colors)
	bufLen := 4 + 2 + (colorsLen * 4) // TODO this using sizeof, but not reflect?
	buf := make([]byte, bufLen)

	offset := 0

	insertUint32(buf, &offset, uint32(bufLen))
	insertUint16(buf, &offset, uint16(colorsLen))
	for _, color := range wd.Colors {
		insertColor(buf, &offset, color)
	}

	// We allocated the buffer up-front for speed (because Buffer and a Writer will be slow), but assert we got that calculation correct
	if offset != bufLen {
		panic(fmt.Errorf("Assertion failed: mismatch between supposed length of buffer: %d, and amount used: %d", bufLen, offset))
	}

	err := c.sendCommand(wd.DeviceID, cmdUpdateLEDs, buf)
	if err != nil {
		return fmt.Errorf("Can't set device-level LEDs: %w", err)
	}

	return nil
}

// Updates just one zone (or, rather, legacy code that renders the command to do so).
// Use this to implement a zone writer, also do a singular writer, then have Thither() take/work on a patch, and be optimal
// - if the patch has >1 device use WriteDevices. >1 zone, use WriteDevice, >1 LED use this, 1 LED use WriteLED
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

	// We allocated the buffer up-front for speed (because Buffer and a Writer will be slow), but assert we got that calculation correct
	if offset != bufLen {
		panic(fmt.Errorf("Assertion failed: mismatch between supposed length of buffer: %d, and amount used: %d", bufLen, offset))
	}

	return buf
}

// TODO: because this is basically a serialisation of the wire object, move it into mode.go next to extractMode
// TODO There should be one func in this file per cmd verb
// most, like this one, should defer to "insertFoo". Write*LEDs are special, because those commands don't have "Request" mirrors, and so they need their own structs

/* NB:
* - This updates the definition of a Mode.
*   - It does this by using cmdUpdateMode, which calls RGBController::SetModeDescrition() (NetworkServer.cpp:671)
*   - You can't make a new one (I think that's what cmdSetCustomMode is for?)
*   - This looks at index, finds that entry in the modes array, and replaces its values with the ones here
* - SetModeDescription() ALSO sets the device's active mode to this mode's index. WTF.
*   - It DOESN'T call UpdateMode(), but the very next time in NetworkServer.cpp does
* - Internally, OpenRGB has an RGBController::SetMode() function: https://gitlab.com/CalcProgrammer1/OpenRGB/-/blob/master/RGBController/RGBController.cpp#L1345
*   - This just takes an int, not a whole mode description, and then calls UpdateMode()
*   - SetMode() is called by the qt GUI in a few places (makes sense, will be on a random thread)
* - UpdateMode() seems to set a flag which makes some other thread call DeviceUpdateMode() (virtual on RGBController; implemented by the subclasses)
* - RGBController::DeviceUpdateMode()
*   - is called from the CLI: it arranges for device->active_mode to be changed (applying flags, loading profiles), then calls device->DeviceUpdateMode() direct, presumably knowing it's on the right thread
*   - is implemented by the actual controller subclasses and actually changes the mode
* So...
* - To change the mode of a device, you change device->active_mode, then arrange for DeviceUpdateMode() to get called
*   - The CLI changes active_mode, then calls DeviceUpdateMode() directly, as it's only got one thread
*   - The GUI calls SetMode(), which changes active_mode, and then calls UpdateMode(), to get DeviceUpdateMode() called on the right thread
*   - The NetworkServer only exposes SetModeDescription(). This does change active_mode, and then call UpdateMode(), but you have to play the whole Mode description back at it.
*     - NetworkProtocol.h is downright misleading, cause the comment says it calls UpdateMode() when actually it calls SetModeDescription()
 */

// SendUpdateMode sends a command to update the details of an existing mode; overwriting the mode at the given index.
// On the network API as it stands (v2), this is the only way to change the active mode - reassert a mode exactly as-was, but the act of asserting it causes it to activate.
func SendUpdateMode(c *Client, deviceID uint32, mode *Mode) error {
	bufLen := 4 + (2 + len(mode.Name) + 1) + 9*4 + 2 + len(mode.Colors)*4 // TODO this using sizeof, but not reflect?
	buf := make([]byte, bufLen)

	offset := 0
	insertUint32(buf, &offset, mode.Index)
	insertString(buf, &offset, mode.Name)
	insertUint32(buf, &offset, mode.Value)
	insertUint32(buf, &offset, uint32(mode.Flags))
	insertUint32(buf, &offset, mode.MinSpeed)
	insertUint32(buf, &offset, mode.MaxSpeed)
	insertUint32(buf, &offset, mode.MinColors)
	insertUint32(buf, &offset, mode.MaxColors)
	insertUint32(buf, &offset, mode.Speed)
	insertUint32(buf, &offset, uint32(mode.Direction))
	insertUint32(buf, &offset, uint32(mode.ColorMode))
	// TODO: extract to insertColors, to match extracting.
	// Because of this, pre-calculating buf length will be impossible, so just use a damn Buffer and Writer (effectively ByteBuilder)
	insertUint16(buf, &offset, uint16(len(mode.Colors)))
	for _, c := range mode.Colors {
		insertColor(buf, &offset, c)
	}

	// We allocated the buffer up-front for speed (because Buffer and a Writer will be slow), but assert we got that calculation correct
	if offset != bufLen {
		panic(fmt.Errorf("Assertion failed: mismatch between supposed length of buffer: %d, and amount used: %d", bufLen, offset))
	}

	err := c.sendCommand(deviceID, cmdUpdateMode, buf)
	if err != nil {
		return fmt.Errorf("Can't set active Mode LEDs: %w", err)
	}

	return nil
}
