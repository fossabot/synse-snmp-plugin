package mibs

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsOutputTable represts SNMP OID .1.3.6.1.2.1.33.1.4.4
type UpsOutputTable struct {
	*core.SnmpTable // base class
}

// NewUpsOutputTable constructs the UpsOutputTable.
func NewUpsOutputTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsOutputTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Output-Table", // Table Name
		".1.3.6.1.2.1.33.1.4.4",    // WalkOid
		[]string{ // Column Names
			"upsOutputLineIndex", // MIB says not accessable. Have seen it in walks.
			"upsOutputVoltage",   // RMS Volts
			"upsOutputCurrent",   // .1 RMS Amp
			"upsOutputPower",     // Watts
			"upsOutputPercentLoad",
		},
		snmpServerBase, // snmpServer
		"1",            // rowBase
		"",             // indexColumn
		"2",            // readableColumn
		false)          // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsOutputTable{SnmpTable: snmpTable}
	table.DevEnumerator = UpsOutputTableDeviceEnumerator{table}
	return table, nil
}

// UpsOutputTableDeviceEnumerator overrides the default SnmpTable device
// enumerator for the output table.
type UpsOutputTableDeviceEnumerator struct {
	Table *UpsOutputTable // Pointer back to the table.
}

// DeviceEnumerator overrides the default SnmpTable device enumerator.
func (enumerator UpsOutputTableDeviceEnumerator) DeviceEnumerator(
	data map[string]interface{}) (devices []*config.DeviceConfig, err error) {
	fmt.Printf("ZZZ: Override: UpsOutputTableDeviceEnumerator, enumerator.Table: %+v\n", enumerator.Table)

	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	for i := 0; i < len(table.Rows); i++ {
		// upsOutputVoltage
		device := config.DeviceConfig{
			Version: "1",
			Type:    "voltage",
			Model:   model,
			Location: config.Location{
				Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
				Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
			},
			Data: map[string]string{
				"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
				"info":       fmt.Sprintf("upsOutputVoltage%d", i),
				"base_oid":   table.Rows[i].BaseOid,
				"table_name": table.Name,
				"row":        fmt.Sprintf("%d", i),
				"column":     "2",
			},
		}
		devices = append(devices, &device)

		// upsOutputCurrent
		device2 := config.DeviceConfig{
			Version: "1",
			Type:    "current",
			Model:   model,
			Location: config.Location{
				Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
				Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
			},
			Data: map[string]string{
				"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
				"info":       fmt.Sprintf("upsOutputCurrent%d", i),
				"base_oid":   table.Rows[i].BaseOid,
				"table_name": table.Name,
				"row":        fmt.Sprintf("%d", i),
				"column":     "3",
			},
		}
		devices = append(devices, &device2)

		// upsOutputPower
		device3 := config.DeviceConfig{
			Version: "1",
			Type:    "power",
			Model:   model,
			Location: config.Location{
				Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
				Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
			},
			Data: map[string]string{
				"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
				"info":       fmt.Sprintf("upsOutputPower%d", i),
				"base_oid":   table.Rows[i].BaseOid,
				"table_name": table.Name,
				"row":        fmt.Sprintf("%d", i),
				"column":     "4",
			},
		}
		devices = append(devices, &device3)

		// upsOutputPercentLoad
		device4 := config.DeviceConfig{
			Version: "1",
			Type:    "status",
			Model:   model,
			Location: config.Location{
				Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
				Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
			},
			Data: map[string]string{
				"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
				"info":       fmt.Sprintf("upsOutputPercentLoad%d", i),
				"base_oid":   table.Rows[i].BaseOid,
				"table_name": table.Name,
				"row":        fmt.Sprintf("%d", i),
				"column":     "5",
			},
		}
		devices = append(devices, &device4)
	}
	return devices, err
}