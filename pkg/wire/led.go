package wire

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
