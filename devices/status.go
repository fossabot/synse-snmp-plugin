package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// SnmpStatus is the handler for the snmp-status device.
var SnmpStatus = sdk.DeviceHandler{
	Type:  "status",
	Model: "snmp-status",

	Read:     SnmpStatusRead,
	Write:    nil, // NYI for V1
	BulkRead: nil,
}

// SnmpStatusRead is the read handler function for snmp-status devices.
func SnmpStatusRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

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
	//fmt.Printf("Status get. oid: %v, result: %+v\n", data["oid"], result)

	// Should be a string.
	resultString := ""
	//fmt.Printf("status reading: %T, %v\n", result.Data, result.Data)
	if result.Data != nil {
		var ok bool
		resultString, ok = result.Data.(string)
		//fmt.Printf("Status get. resultInt: %d\n", resultInt)
		if !ok {
			// Could be an int as well.
			var resultInt int
			resultInt, ok = result.Data.(int)
			if !ok {
				return nil, fmt.Errorf(
					"Expected string or int status reading, got type: %T, value: %v",
					result.Data, result.Data)
			}
			resultString = fmt.Sprintf("%d", resultInt)
		}
	}
	// Create the reading.
	readings = []*sdk.Reading{
		sdk.NewReading("status", resultString),
	}
	//fmt.Printf("Status readings: %+v\n", readings)
	//fmt.Printf("Status readings[0]: %+v\n", readings[0])
	return readings, nil
}
