package wire

// LED represents an actual LED in a device, which can emit a color. On the
// wire protocol however, the LED objects don't contain the colors they're
// emitting.
// Note that in some devices, the addresses on the bus don't correspond
// one-to-one with physical LEDs (probably SMT devices). Sometimes an address
// addresses more than one LED. This object represents addresses, ie what can
// be addressed by the controller.
type LED struct {
	Index uint16
	Name  string
	Value uint32 // driver-internal, eg LED mapping
}

func extractLEDs(buf []byte, offset *int) []*LED {
	ledsCount := extractUint16(buf, offset)

	ls := make([]*LED, ledsCount)
	for i := uint16(0); i < ledsCount; i++ {
		ls[i] = extractLED(buf, offset, i)
	}

	return ls
}

func extractLED(buf []byte, offset *int, idx uint16) *LED {
	l := &LED{Index: idx}

	l.Name = extractString(buf, offset)
	l.Value = extractUint32(buf, offset)

	return l
}
