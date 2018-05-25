package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// SnmpVoltage is the handler for the snmp-voltage device.
var SnmpVoltage = sdk.DeviceHandler{
	Type:  "voltage",
	Model: "snmp-voltage",

	Read:     SnmpVoltageRead,
	Write:    nil, // NYI for V1
	BulkRead: nil,
}

// SnmpVoltageRead is the read handler function for snmp-voltage devices.
func SnmpVoltageRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

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
	//fmt.Printf("Power get. oid: %v, result: %+v\n", data["oid"], result)

	// Should be an int.
	resultInt, ok := result.Data.(int)
	resultFloat := float32(resultInt) / 10.0
	//fmt.Printf("Power get. resultInt: %d\n", resultInt)
	if !ok {
		return nil, fmt.Errorf(
			"Expected int voltage reading, got type: %T, value: %v",
			result.Data, result.Data)
	}
	resultString := fmt.Sprintf("%.1f", resultFloat)

	// Create the reading.
	readings = []*sdk.Reading{
		sdk.NewReading("voltage", resultString),
	}
	//fmt.Printf("Power readings: %+v\n", readings)
	//fmt.Printf("Power readings[0]: %+v\n", readings[0])
	return readings, nil
}
