package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// SnmpIdentity is the handler for the snmp-identity device.
var SnmpIdentity = sdk.DeviceHandler{
	Type:  "identity",
	Model: "snmp-identity",

	Read:     SnmpIdentityRead,
	Write:    nil, // NYI for V1
	BulkRead: nil,
}

// SnmpIdentityRead is the read handler function for snmp-identity devices.
func SnmpIdentityRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

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
	//fmt.Printf("Identity get. oid: %v, result: %+v\n", data["oid"], result)

	// Should be a string.
	resultString, ok := result.Data.(string)
	//fmt.Printf("Identity get. resultInt: %d\n", resultInt)
	if !ok {
		return nil, fmt.Errorf(
			"Expected int identity reading, got type: %T, value: %v",
			result.Data, result.Data)
	}

	// Create the reading.
	readings = []*sdk.Reading{
		sdk.NewReading("identity", resultString),
	}
	//fmt.Printf("Identity readings: %+v\n", readings)
	//fmt.Printf("Identity readings[0]: %+v\n", readings[0])
	return readings, nil
}
