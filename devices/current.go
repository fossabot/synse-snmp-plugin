package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// SnmpCurrent is the handler for the snmp-current device.
var SnmpCurrent = sdk.DeviceHandler{
	Type:  "current",
	Model: "snmp-current",

	Read:     SnmpCurrentRead,
	Write:    nil, // NYI for V1
	BulkRead: nil,
}

// SnmpCurrentRead is the read handler function for snmp-current devices.
func SnmpCurrentRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

	// Arg checks.
	if device == nil {
		return nil, fmt.Errorf("device is nil")
	}

	// Get the SNMP device config from the strings in the data.
	data := device.Data
	snmpConfig, err := core.GetDeviceConfig(data)
	if err != nil {
		return nil, err
	}

	// Create SnmpClient.
	snmpClient, err := core.NewSnmpClient(snmpConfig)
	if err != nil {
		return nil, err
	}

	// Read the SNMP OID in the device config.
	result, err := snmpClient.Get(data["oid"])
	if err != nil {
		return nil, err
	}

	// Should be an int.
	resultInt, ok := result.Data.(int)
	// All current readings are in .1 Amps at this time.
	// If this changes, the constant multiplier below will not always be true.
	resultFloat := float32(resultInt) / 10.0
	if !ok {
		return nil, fmt.Errorf(
			"Expected int current reading, got type: %T, value: %v",
			result.Data, result.Data)
	}
	resultString := fmt.Sprintf("%.1f", resultFloat)

	// Create the reading.
	readings = []*sdk.Reading{
		sdk.NewReading("current", resultString),
	}
	return readings, nil
}