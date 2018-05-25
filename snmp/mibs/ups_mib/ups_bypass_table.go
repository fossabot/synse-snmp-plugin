package mibs

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsBypassTable represts SNMP OID .1.3.6.1.2.1.33.1.5.3
type UpsBypassTable struct {
	*core.SnmpTable // base class
}

// NewUpsBypassTable constructs the UpsBypassTable.
func NewUpsBypassTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsBypassTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Bypass-Table", // Table Name
		".1.3.6.1.2.1.33.1.2",      // WalkOid
		[]string{ // Column Names
			"upsBypassLineIndex",
			"upsBypassVoltage",
			"upsBypassCurrent",
			"upsBypassPower",
		},
		snmpServerBase, // snmpServer
		"1",            // rowBase
		"",             // indexColumn
		"2",            // readableColumn
		false)          // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsBypassTable{SnmpTable: snmpTable}
	table.DevEnumerator = UpsBypassTableDeviceEnumerator{table}
	return table, nil
}

// UpsBypassTableDeviceEnumerator overrides the default SnmpTable device
// enumerator for the bypass table.
type UpsBypassTableDeviceEnumerator struct {
	Table *UpsBypassTable // Pointer back to the table.
}

// DeviceEnumerator overrides the default SnmpTable device enumerator.
func (enumerator UpsBypassTableDeviceEnumerator) DeviceEnumerator(
	data map[string]interface{}) (devices []*config.DeviceConfig, err error) {

	// Pull out the table, mib, device model, SNMP DeviceConfig.
	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	snmpDeviceConfigMap, err := table.SnmpServerBase.DeviceConfig.ToMap()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(table.Rows); i++ {
		// upsBypassVoltage
		// deviceData gets shimmed into the DeviceConfig for each synse device.
		// It varies slightly for each device below.
		deviceData := map[string]string{
			"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
			"info":       fmt.Sprintf("upsBypassVoltage%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "2",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 2), // base_oid and integer column.
		}
		deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device := config.DeviceConfig{
			Version: "1",
			Type:    "voltage",
			Model:   model,
			Location: config.Location{
				Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
				Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
			},
			Data: deviceData,
		}
		devices = append(devices, &device)

		// upsBypassCurrent ---------------------------------------------------------
		deviceData = map[string]string{
			"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
			"info":       fmt.Sprintf("upsBypassCurrent%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "3",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 3), // base_oid and integer column.
		}
		deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device2 := config.DeviceConfig{
			Version: "1",
			Type:    "current",
			Model:   model,
			Location: config.Location{
				Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
				Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
			},
			Data: deviceData,
		}
		devices = append(devices, &device2)

		// upsBypassPower --------------------------------------------------------------
		deviceData = map[string]string{
			"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
			"info":       fmt.Sprintf("upsBypassPower%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "4",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 4), // base_oid and integer column.
		}
		deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device3 := config.DeviceConfig{
			Version: "1",
			Type:    "power",
			Model:   model,
			Location: config.Location{
				Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
				Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
			},
			Data: deviceData,
		}
		devices = append(devices, &device3)
	}
	return devices, err
}
