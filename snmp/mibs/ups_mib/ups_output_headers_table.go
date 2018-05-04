package mibs

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsOutputHeadersTable represts SNMP OID .1.3.6.1.2.1.33.1.4
type UpsOutputHeadersTable struct {
	*core.SnmpTable // base class
}

// NewUpsOutputHeadersTable constructs the UpsOutputHeadersTable.
func NewUpsOutputHeadersTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsOutputHeadersTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Output-Headers-Table", // Table Name
		".1.3.6.1.2.1.33.1.4",              // WalkOid
		[]string{ // Column Names
			"upsOutputSource",
			"upsOutputFrequency",
			"upsOutputNumLines",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		true)           // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsOutputHeadersTable{SnmpTable: snmpTable}
	table.DevEnumerator = UpsOutputHeadersTableDeviceEnumerator{table}
	return table, nil
}

// UpsOutputHeadersTableDeviceEnumerator overrides the default SnmpTable device
// enumerator for the output headers table.
type UpsOutputHeadersTableDeviceEnumerator struct {
	Table *UpsOutputHeadersTable // Pointer back to the table.
}

// DeviceEnumerator overrides the default SnmpTable device enumerator.
func (enumerator UpsOutputHeadersTableDeviceEnumerator) DeviceEnumerator(
	data map[string]interface{}) (devices []*config.DeviceConfig, err error) {
	fmt.Printf("ZZZ: Override: UpsBatteryOutputHeadersDeviceEnumerator, enumerator.Table: %+v\n", enumerator.Table)

	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	// This is always a single row table.

	// upsOutputSource
	device := config.DeviceConfig{
		Version: "1",
		Type:    "status", // TODO: This is new for synse.
		Model:   model,
		Location: config.Location{
			Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
			Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
		},
		Data: map[string]string{
			"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
			"info":       "upsOutputSource",
			"base_oid":   table.Rows[0].BaseOid,
			"table_name": table.Name,
			"row":        "0",
			"column":     "1",
		},
	}
	devices = append(devices, &device)

	// upsOutputFrequency
	device2 := config.DeviceConfig{
		Version: "1",
		Type:    "frequency",
		Model:   model,
		Location: config.Location{
			Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
			Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
		},
		Data: map[string]string{
			"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
			"info":       "upsOutputFrequency",
			"base_oid":   table.Rows[0].BaseOid,
			"table_name": table.Name,
			"row":        "0",
			"column":     "2",
		},
	}
	devices = append(devices, &device2)

	// upsOutputNumLines
	device3 := config.DeviceConfig{
		Version: "1",
		Type:    "status", // TODO: This is new for synse.
		Model:   model,
		Location: config.Location{
			Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
			Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
		},
		Data: map[string]string{
			"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
			"info":       "upsOutputNumLines",
			"base_oid":   table.Rows[0].BaseOid,
			"table_name": table.Name,
			"row":        "0",
			"column":     "3",
		},
	}
	devices = append(devices, &device3)

	return devices, err
}
