package mibs

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsBatteryTable represts SNMP OID .1.3.6.1.2.1.33.1.2
type UpsBatteryTable struct {
	*core.SnmpTable // base class
}

// Battery Status enumeration.
const (
	BatteryStatusUnknown  = 1
	BatteryStatusNormal   = 2
	BatteryStatusLow      = 3
	BatteryStatusDepleted = 4
)

// GetBatteryStatus gets the status of the battery as a string from the SNMP data.
// We may not get the data from the UPS, so status can be nil.
func GetBatteryStatus(status *int) string {
	if status == nil {
		return "undefined"
	}
	theStatus := *status
	if theStatus == 1 {
		return "unknown"
	}
	if theStatus == 2 {
		return "normal"
	}
	if theStatus == 3 {
		return "depleted"
	}
	return "undefined"
}

// NewUpsBatteryTable constructs the UpsBatteryTable.
func NewUpsBatteryTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsBatteryTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Battery-Table", // Table Name
		".1.3.6.1.2.1.33.1.2",       // WalkOid
		[]string{ // Column Names
			"upsBatteryStatus",
			"upsSecondsOnBattery", // Zero if not on battery power.
			"upsEstimatedMinutesRemaining",
			"upsEstimatedChargeRemaining", // Percentage
			"upsBatteryVoltage",           // Units .1 VDC.
			"upsBatteryCurrent",           // Units .1 Amp DC.
			"upsBacontteryTemperature",    // Units degrees C.
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		true)           // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsBatteryTable{SnmpTable: snmpTable}
	// Override the default Device Enumerator
	table.DevEnumerator = UpsBatteryTableDeviceEnumerator{table}
	return table, nil
}

// UpsBatteryTableDeviceEnumerator overrides the default SnmpTable device
// enumerator for the battery table.
type UpsBatteryTableDeviceEnumerator struct {
	Table *UpsBatteryTable // Pointer back to the table.
}

// DeviceEnumerator overrides the default SnmpTable device enumerator.
func (enumerator UpsBatteryTableDeviceEnumerator) DeviceEnumerator(
	data map[string]interface{}) (devices []*config.DeviceConfig, err error) {

	// Pull out the table, mib, device model, SNMP DeviceConfig
	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	snmpDeviceConfigMap, err := table.SnmpServerBase.DeviceConfig.ToMap()
	if err != nil {
		return nil, err
	}

	// This is always a single row table.

	// upsBatteryStatus
	// deviceData gets shimmed into the DeviceConfig for each synse device.
	// It varies slightly for each device below.
	deviceData := map[string]string{
		"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
		"info":       "upsBatteryStatus",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "1",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 1), // base_oid and integer column.
	}
	deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device := config.DeviceConfig{
		Version: "1",
		Type:    "status", // TODO: This is new for synse.
		Model:   model,
		Location: config.Location{
			Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
			Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
		},
		Data: deviceData,
	}
	devices = append(devices, &device)

	// upsSecondsOnBattery --------------------------------------------------------
	deviceData = map[string]string{
		"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
		"info":       "upsSecondsOnBattery",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "2",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 2), // base_oid and integer column.
	}
	deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device2 := config.DeviceConfig{
		Version: "1",
		Type:    "status", // TODO: This is new for synse.
		Model:   model,
		Location: config.Location{
			Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
			Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
		},
		Data: deviceData,
	}
	devices = append(devices, &device2)

	// upsEstimatedMinutesRemaining -----------------------------------------------
	deviceData = map[string]string{
		"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
		"info":       "upsEstimatedMinutesRemaining",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "3",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 3), // base_oid and integer column.
	}
	deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device3 := config.DeviceConfig{
		Version: "1",
		Type:    "status", // TODO: This is new for synse.
		Model:   model,
		Location: config.Location{
			Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
			Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
		},
		Data: deviceData,
	}
	devices = append(devices, &device3)

	// upsEstimatedChargeRemaining ------------------------------------------------
	deviceData = map[string]string{
		"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
		"info":       "upsEstimatedChargeRemaining",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "4",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 4), // base_oid and integer column.
	}
	deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device4 := config.DeviceConfig{
		Version: "1",
		Type:    "status", // TODO: This is new for synse.
		Model:   model,
		Location: config.Location{
			Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
			Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
		},
		Data: deviceData,
	}
	devices = append(devices, &device4)

	// upsBatteryVoltage ----------------------------------------------------------
	deviceData = map[string]string{
		"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
		"info":       "upsBatteryVoltage",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "5",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 5), // base_oid and integer column.
	}
	deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device5 := config.DeviceConfig{
		Version: "1",
		Type:    "voltage",
		Model:   model,
		Location: config.Location{
			Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
			Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
		},
		Data: deviceData,
	}
	devices = append(devices, &device5)

	// upsBatteryCurrent ---------------------------------------------------------
	deviceData = map[string]string{
		"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
		"info":       "upsBatteryCurrent",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "6",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 6), // base_oid and integer column.
	}
	deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device6 := config.DeviceConfig{
		Version: "1",
		Type:    "current",
		Model:   model,
		Location: config.Location{
			Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
			Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
		},
		Data: deviceData,
	}
	devices = append(devices, &device6)

	// upsBatteryTemperature  -----------------------------------------------------
	deviceData = map[string]string{
		"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
		"info":       "upsBatteryTemperature",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "7",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 7), // base_oid and integer column.
	}
	deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device7 := config.DeviceConfig{
		Version: "1",
		Type:    "temperature",
		Model:   model,
		Location: config.Location{
			Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
			Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
		},
		Data: deviceData,
	}
	devices = append(devices, &device7)

	//for i := 0; i < len(devices); i++ {
	//	fmt.Printf("device[%d]: %+v\n", i, devices[i])
	//}
	return devices, err
}
