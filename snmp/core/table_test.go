package core

import (
	"fmt"
	"testing"
)

// TestTable
// Initial test creates a table based on the UPS-MIB, upsInput table.
func TestTable(t *testing.T) {
	t.Logf("TestTable start")

	// In order to create the table, we need to create an SNMP Server.
	// In order to create the SNMP server, we need to have an SnmpClient.

	// Create SecurityParameters for the config that should connect to the emulator.
	securityParameters, err := NewSecurityParameters(
		"simulator",  // User Name
		SHA,          // Authentication Protocol
		"auctoritas", // Authentication Passphrase
		AES,          // Privacy Protocol
		"privatus")   // Privacy Passphrase
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	t.Logf("securityParameters: %+v", securityParameters)

	// Create a config.
	config, err := NewDeviceConfig(
		"v3",        // SNMP v3
		"127.0.0.1", // Endpoint
		1024,        // Port
		securityParameters,
		"public") //  Context name
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	t.Logf("config: %+v", config)

	// Create a client.
	client, err := NewSnmpClient(config)
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	t.Logf("client: %+v", client)

	// Create SnmpServerBase
	snmpServer, err := NewSnmpServerBase(
		client,
		config,
		"test_rack")
	if err != nil {
		t.Error(err) // Fail the test.
	}

	t.Logf("snmpServer: %+v", snmpServer)

	// Create SnmpTable for the UPS input power.
	testUpsInputTable, err := NewSnmpTable(
		"upsInputTable",         // Table name. Same as OID .1.3.6.1.2.1.33.1.3.3 (Walk OID)
		".1.3.6.1.2.1.33.1.3.3", // Walk OID
		[]string{ // Column names
			"upsInputLineIndex",
			"upsInputFrequency",
			"upsInputVoltage",
			"upsInputCurrent",
			"upsInputTruePower",
		},
		snmpServer, // server
		"1",        // rowBase
		"",         // indexColumn
		"2",        // readableColumn
		false)      // flattened table
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	testUpsInputTable.Dump()

	// Call DeviceEnumerator for testUpsInputTable.
	// This is currently the default which does nothing, but that may change.
	fmt.Printf("Calling Device Enumerate()\n")
	fmt.Printf("testUpsInputTable: %v\n", testUpsInputTable)
	_, err = testUpsInputTable.DevEnumerator.DeviceEnumerator(nil)
	fmt.Printf("Called Device Enumerate()\n")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("TestTable end")
}
