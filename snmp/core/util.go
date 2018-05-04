package core

import (
	"fmt"
)

// This file contains utility functions. In the future we could put them in
// a "library" repo.

// TranslatePrintableASCII translates byte arrays from gosnmp to a printable string if possible.
// If this call fails, the caller should normally just keep the raw byte array.
// This call makes no attemp to support extended (8bit) ASCII.
func TranslatePrintableASCII(x interface{}) (string, error) {
	bytes, ok := x.([]uint8)
	if !ok {
		return "", fmt.Errorf("Failure converting type: %T, data: %v to byte array", x, x)
	}

	for i := 0; i < len(bytes); i++ {
		if !(bytes[i] < 0x80 && bytes[i] > 0x1f) {
			return "", fmt.Errorf("Unable to convert %x byte %x at %d to printable Ascii", bytes, bytes[i], i)
		}
	}
	return string(bytes), nil
}
