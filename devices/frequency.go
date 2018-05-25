package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// SnmpFrequency is the handler for the snmp-frequency device.
var SnmpFrequency = sdk.DeviceHandler{
	Type:  "frequency",
	Model: "snmp-frequency",

	Read:     SnmpFrequencyRead,
	Write:    nil, // NYI for V1
	BulkRead: nil,
}

// SnmpFrequencyRead is the read handler function for snmp-frequency devices.
func SnmpFrequencyRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

	// Arg checks.
	if device == nil {
		return nil, fmt.Errorf("device is nil")
	}

	// Get the SNMP device config from the strings in data.
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
	resultFloat := float32(resultInt) / 10.0
	if !ok {
		return nil, fmt.Errorf(
			"Expected int frequency reading, got type: %T, value: %v",
			result.Data, result.Data)
	}
	resultString := fmt.Sprintf("%.1f", resultFloat)

	// Create the reading.
	readings = []*sdk.Reading{
		sdk.NewReading("frequency", resultString),
	}
	return readings, nil
}
