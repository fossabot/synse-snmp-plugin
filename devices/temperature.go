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

	// Get the device config. (You can't, it's private, but you can get the members.)
	data := device.Data
	// Create the SnmpClient from the strings in data.
	snmpConfig, err := core.GetDeviceConfig(data)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("snmpConfig: %+v\n", snmpConfig) // use it or lose it

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
	//fmt.Printf("Temperature get. oid: %v, result: %+v\n", data["oid"], result)

	// Should be an int.
	// TODO: Multiplier (?) (Not needed here, but may be in other places..)
	resultInt, ok := result.Data.(int)
	//fmt.Printf("Temperature get. resultInt: %d\n", resultInt)
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
	//fmt.Printf("Temperature readings: %+v\n", readings)
	//fmt.Printf("Temperature readings[0]: %+v\n", readings[0])
	return readings, nil
}
