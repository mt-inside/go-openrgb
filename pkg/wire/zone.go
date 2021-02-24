package wire

import "fmt"

//go:generate stringer -type=ZoneType
type ZoneType uint32

const (
	Singular ZoneType = 0 // Doesn't imply precisely 1 LED. TODO what are the invariants.
	Linear   ZoneType = 1
	Planar   ZoneType = 2
)

type Zone struct {
	Index     uint16
	Name      string
	Type      ZoneType
	MinLEDs   uint32 // min!=max => user-resizable (depending on what's plugged in)
	MaxLEDs   uint32
	TotalLEDs uint32 // current size?
}

func extractZones(buf []byte, offset *int) []*Zone {
	zoneCount := extractUint16(buf, offset)

	zs := make([]*Zone, zoneCount)
	for i := uint16(0); i < zoneCount; i++ {
		zs[i] = extractZone(buf, offset, i)
	}

	return zs
}

func extractZone(buf []byte, offset *int, idx uint16) *Zone {
	z := &Zone{Index: idx}

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
