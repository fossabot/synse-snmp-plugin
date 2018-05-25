package core

import (
	"fmt"
)

// This file contains utility functions. In the future we could put them in
// a "library" repo.

// CopyMapStringString returns a copy of the map passed in without error.
func CopyMapStringString(m map[string]string) map[string]string {
	target := make(map[string]string)
	for k, v := range m {
		target[k] = v
	}
	return target
}

// GetRackAndBoard pulls the rack and board ids out of data.
// Used for pulling them out of the data in a DeviceEnumerator.
func GetRackAndBoard(data map[string]interface{}) (rack string, board string, err error) {
	// Parameter check.
	if data == nil {
		return "", "", fmt.Errorf("data is nil")
	}

	// Get the rack id.
	value, ok := data["rack"]
	if !ok {
		return "", "", fmt.Errorf("rack is not in data")
	}
	rack, ok = value.(string)
	if !ok {
		return "", "", fmt.Errorf("rack is not a string, %T", value)
	}

	// Get the board id.
	value, ok = data["board"]
	if !ok {
		return "", "", fmt.Errorf("board is not in data")
	}
	board, ok = value.(string)
	if !ok {
		return "", "", fmt.Errorf("board is not a string, %T", value)
	}
	return rack, board, nil
}

// MergeMapStringString returns a new map with the contents of both maps passed
// in. Errors out on duplicate keys.
func MergeMapStringString(a map[string]string, b map[string]string) (map[string]string, error) {
	merged := CopyMapStringString(a)
	for k, v := range b {
		_, inMap := merged[k]
		if inMap {
			return nil, fmt.Errorf("Key %v already in merged map: %v", k, merged)
		}
		merged[k] = v
	}
	return merged, nil
}

// TranslatePrintableASCII translates byte arrays from gosnmp to a printable
// string if possible. If this call fails, the caller should normally just keep
// the raw byte array. This call makes no attempt to support extended (8bit)
// ASCII. We need this function since there is no differentiation between
// strings and byte arrays in the SNMP protocol.
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
