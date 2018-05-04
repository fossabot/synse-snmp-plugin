package mibs

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsInputTable represts SNMP OID .1.3.6.1.2.1.33.1.3.3
type UpsInputTable struct {
	*core.SnmpTable // base class
}

// NewUpsInputTable constructs the UpsInputTable.
func NewUpsInputTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsInputTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Input-Table", // Table Name
		".1.3.6.1.2.1.33.1.3.3",   // WalkOid
		[]string{ // Column Names
			"upsInputLineIndex", // MIB says not accessable. Have seen it in walks.
			"upsInputFrequency", // .1 Hertz
			"upsInputVoltage",   // RMS Volts
			"upsInputCurrent",   // .1 RMS Amp
			"upsInputTruePower", // Down with False Power! (Watts)
		},
		snmpServerBase, // snmpServer
		"1",            // rowBase
		"",             // indexColumn
		"2",            // readableColumn
		false)          // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsInputTable{SnmpTable: snmpTable}
	table.DevEnumerator = UpsInputTableDeviceEnumerator{table}
	return table, nil
}

// UpsInputTableDeviceEnumerator overrides the default SnmpTable device
// enumerator for the input table.
type UpsInputTableDeviceEnumerator struct {
	Table *UpsInputTable // Pointer back to the table.
}

// DeviceEnumerator overrides the default SnmpTable device enumerator.
func (enumerator UpsInputTableDeviceEnumerator) DeviceEnumerator(
	data map[string]interface{}) (devices []*config.DeviceConfig, err error) {
	fmt.Printf("ZZZ: Override: UpsInputTableDeviceEnumerator, enumerator.Table: %+v\n", enumerator.Table)

	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	for i := 0; i < len(table.Rows); i++ {
		// upsInputFrequency
		device := config.DeviceConfig{
			Version: "1",
			Type:    "frequency",
			Model:   model,
			Location: config.Location{
				Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
				Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
			},
			Data: map[string]string{
				"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
				"info":       fmt.Sprintf("upsInputFrequency%d", i),
				"base_oid":   table.Rows[i].BaseOid,
				"table_name": table.Name,
				"row":        fmt.Sprintf("%d", i),
				"column":     "2",
			},
		}
		devices = append(devices, &device)

		// upsInputVoltage
		device2 := config.DeviceConfig{
			Version: "1",
			Type:    "voltage",
			Model:   model,
			Location: config.Location{
				Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
				Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
			},
			Data: map[string]string{
				"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
				"info":       fmt.Sprintf("upsInputVoltage%d", i),
				"base_oid":   table.Rows[i].BaseOid,
				"table_name": table.Name,
				"row":        fmt.Sprintf("%d", i),
				"column":     "3",
			},
		}
		devices = append(devices, &device2)

		// upsInputCurrent
		device3 := config.DeviceConfig{
			Version: "1",
			Type:    "current",
			Model:   model,
			Location: config.Location{
				Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
				Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
			},
			Data: map[string]string{
				"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
				"info":       fmt.Sprintf("upsInputCurrent%d", i),
				"base_oid":   table.Rows[i].BaseOid,
				"table_name": table.Name,
				"row":        fmt.Sprintf("%d", i),
				"column":     "4",
			},
		}
		devices = append(devices, &device3)

		// upsInputTruePower
		device4 := config.DeviceConfig{
			Version: "1",
			Type:    "power",
			Model:   model,
			Location: config.Location{
				Rack:  "TODO", // TODO: Needs to be passed in by the data parameter.
				Board: "TODO", // TODO: Needs to be passed in by whatever doles out the board ids.
			},
			Data: map[string]string{
				"id":         "TODO", // Needs to be passed in by the board (UPS SNMP Server)
				"info":       fmt.Sprintf("upsInputTruePower%d", i),
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