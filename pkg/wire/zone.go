package wire

import "fmt"

//go:generate stringer -type=ZoneType
type ZoneType uint32

const (
	Singular ZoneType = 0 // Doesn't imply precisely 1 LED. TODO what are the invariants.
	Linear   ZoneType = 1
	Planar   ZoneType = 2 // Still a linear array of LEDs on the wire, but also with a Matrix to tell you how to place them
)

type Zone struct {
	Index        uint16
	Name         string
	Type         ZoneType
	MinLEDs      uint32 // min!=max => user-resizable (depending on what's plugged in)
	MaxLEDs      uint32
	TotalLEDs    uint32 // current size?
	MatrixWidth  uint32
	MatrixHeight uint32
	Maxtrix      [][]uint32
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

	matrixSize := extractUint16(buf, offset) // This seems to exist so you can skip a matrix if you don't like it?
	if z.Type != Planar {                    // TODO: the server code switches on whether maxtrixSize == 0
		if matrixSize != 0 {
			panic("Assertion failed: for non-planar Zones, matrix size should be 0")
		}
	} else {
		z.MatrixHeight = extractUint32(buf, offset)
		z.MatrixWidth = extractUint32(buf, offset)
		if uint16(4+4+(z.MatrixHeight*z.MatrixWidth*4)) != matrixSize {
			panic("Assertion failed: planar Zone matrix sizes don't add up")
		}
		fmt.Printf("TODO matrix: %d x %d\n", z.MatrixWidth, z.MatrixHeight)

		*offset += int(matrixSize) - 4 - 4
	}

	return z
}
