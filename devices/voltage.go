package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// SnmpVoltage is the handler for the SNMP OIDs that report voltage.
var SnmpVoltage = sdk.DeviceHandler{
	Name: "voltage",
	Read: SnmpVoltageRead,
}

// SnmpVoltageRead is the read handler function for synse SNMP devices that report voltage.
func SnmpVoltageRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

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
	result, err := snmpClient.Get(fmt.Sprint(data["oid"]))
	if err != nil {
		return nil, err
	}

	// Account for a multiplier if any and convert to float.
	var resultFloat float32
	resultFloat, err = MultiplyReading(result, data)
	if err != nil {
		return nil, err
	}

	// Create the reading.
	reading, err := device.GetOutput("voltage").MakeReading(resultFloat)
	if err != nil {
		return nil, err
	}

	readings = []*sdk.Reading{reading}
	return readings, nil
}
