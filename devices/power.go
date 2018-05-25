package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// SnmpPower is the handler for the snmp-power device.
var SnmpPower = sdk.DeviceHandler{
	Type:  "power",
	Model: "snmp-power",

	Read:     SnmpPowerRead,
	Write:    nil, // NYI for V1
	BulkRead: nil,
}

// SnmpPowerRead is the read handler function for snmp-power devices.
func SnmpPowerRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

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
	if !ok {
		return nil, fmt.Errorf(
			"Expected int power reading, got type: %T, value: %v",
			result.Data, result.Data)
	}
	resultString := fmt.Sprintf("%d", resultInt)

	// Create the reading.
	readings = []*sdk.Reading{
		sdk.NewReading("power", resultString),
	}
	return readings, nil
}
