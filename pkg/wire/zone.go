package wire

import "fmt"

//go:generate stringer -type=ZoneType
// ZoneType specifies how the zone should be interpreted.
// On the wire, all zones are necessarily 1-D vectors, so ZoneType says how they should be interpreted.
type ZoneType uint32

const (
	// Singular means a point. Doesn't imply precisely 1 LED; I've seen 8-led Singular
	// zones. TODO what are the invariants?
	Singular ZoneType = 0
	// Linear means a line of LEDs.
	Linear ZoneType = 1
	// Planar menas a grid of LEDs. I think this is coupled to a non-empty Matrix (with the
	// matrix width and height giving the shape of the grid, and the matrix
	// itself giving the positions of its LEDs). If this is the case, this
	// field seems redundant?
	Planar ZoneType = 2
)

// Zone represents a sub-section of a device, ie some of its LEDs.
// Zones are only useful with using per-LED-color "direct" modes.
// All Devices have at least one zone; if the physical device doesn't really have separate zones (like a DRAM DIMM), OpenRGB fakes one up.
type Zone struct {
	Index uint16
	Name  string
	Type  ZoneType
	// min!=max => user-resizable (ie tell OpenRGB how many addresses the device plugged into this controller has)
	MinLEDs      uint32
	MaxLEDs      uint32
	TotalLEDs    uint32 // current size
	MatrixWidth  uint32 // zero for non-Planar zones
	MatrixHeight uint32 // zero for non-Planar zones
	// To be interpreted as a grid sized as MatrixWidth x MatrixHeight. This is meant for keyboards and the like, specifying which parts of the grid have bottons/lights ,and which don't. TODO: don't know what the values in this represent.
	Maxtrix [][]uint32
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
