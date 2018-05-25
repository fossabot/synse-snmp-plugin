package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// SnmpTemperature is the handler for the snmp-temperature device.
var SnmpTemperature = sdk.DeviceHandler{
	Type:  "temperature",
	Model: "snmp-temperature",

	Read:     SnmpTemperatureRead,
	Write:    nil, // NYI for V1
	BulkRead: nil,
}

// SnmpTemperatureRead is the read handler function for snmp-temperature devices.
func SnmpTemperatureRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

	// Arg checks.
	if device == nil {
		return nil, fmt.Errorf("device is nil")
	}

	// Get the device config from the strings in data.
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
	// Currrently all temperatures are degrees C. If this changes
	// we may consider introducing a multipier in the device config data.
	if !ok {
		return nil, fmt.Errorf(
			"Expected int temperature reading, got type: %T, value: %v",
			result.Data, result.Data)
	}
	resultString := fmt.Sprintf("%d", resultInt)

	// Create the reading.
	readings = []*sdk.Reading{
		sdk.NewReading("temperature", resultString),
	}
	return readings, nil
}
