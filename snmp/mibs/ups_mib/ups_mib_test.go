package mibs

import (
	"fmt"
	"testing"

	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// TestUpsMib
// Initial test creates all tables based on the UPS-MIB.
func TestUpsMib(t *testing.T) { // nolint: gocyclo
	// In order to create the table, we need to create an SNMP Server.
	// In order to create the SNMP server, we need to have an SnmpClient.

	// Create SecurityParameters for the config that should connect to the emulator.
	securityParameters, err := core.NewSecurityParameters(
		"simulator",  // User Name
		core.SHA,     // Authentication Protocol
		"auctoritas", // Authentication Passphrase
		core.AES,     // Privacy Protocol
		"privatus")   // Privacy Passphrase
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	t.Logf("securityParameters: %+v", securityParameters)

	// Create a config.
	config, err := core.NewDeviceConfig(
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
	client, err := core.NewSnmpClient(config)
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	t.Logf("client: %+v", client)

	// Create SnmpServerBase
	snmpServer, err := core.NewSnmpServerBase(
		client,
		config,
		"test_rack")
	if err != nil {
		t.Fatal(err) // Fail the test.
	}
	t.Logf("snmpServer: %+v", snmpServer)

	// Create the UpsMib and dump it.
	fmt.Printf("TestUpsMib creating mib\n")
	testUpsMib, err := NewUpsMib(snmpServer)
	if err != nil {
		t.Fatal(err) // Fail the test.
	}
	fmt.Printf("TestUpsMib created mib\n")
	//fmt.Printf("testUpsMib: %+v\n", testUpsMib)
	testUpsMib.Dump()

	// We should have 19 tables.
	tableCount := len(testUpsMib.Tables)
	if tableCount != 19 {
		t.Fatalf("testUpsMib: Expected 19 tables, got %d", tableCount)
	}

	// Get the ups identity data from the test MIB.
	upsIdentity := testUpsMib.UpsIdentityTable.UpsIdentity
	fmt.Printf("upsIdentity:  %+v\n", upsIdentity)
	fmt.Printf("Manufacturer:         %v\n", upsIdentity.Manufacturer)
	fmt.Printf("Model:                %v\n", upsIdentity.Model)
	fmt.Printf("UpsSoftwareVersion:   %v\n", upsIdentity.UpsSoftwareVersion)
	fmt.Printf("AgentSoftwareVersion: %v\n", upsIdentity.AgentSoftwareVersion)
	fmt.Printf("Name:                 %v\n", upsIdentity.Name)
	fmt.Printf("AttachedDevices:      %v\n", upsIdentity.AttachedDevices)

	// Verify expected ups identity data from the test MIB.
	if upsIdentity == nil {
		t.Fatal("upsIdentity is nil")
	}

	if upsIdentity.Manufacturer != "Eaton Corporation" {
		t.Fatalf("Expected upsIdentity.Manufacturer [Eaton Corporation], got [%v]", upsIdentity.Manufacturer)
	}

	if upsIdentity.Model != "PXGMS UPS + EATON 93PM" {
		t.Fatalf("Expected upsIdentity.Model [PXGMS UPS + EATON 93PM], got [%v]", upsIdentity.Model)
	}

	if upsIdentity.UpsSoftwareVersion != "INV: 1.44.0000" {
		t.Fatalf("Expected upsIdentity.UpsSoftwareVersion [INV: 1.44.0000], got [%v]", upsIdentity.UpsSoftwareVersion)
	}

	if upsIdentity.AgentSoftwareVersion != "2.3.7" {
		t.Fatalf("Expected upsIdentity.AgentSoftwareVersion [2.3.7], got [%v]", upsIdentity.AgentSoftwareVersion)
	}

	if upsIdentity.Name != "ID: EM111UXX06, Msg: 9PL15N0000E40R2" {
		t.Fatalf("Expected upsIdentity.Name [ID: EM111UXX06, Msg: 9PL15N0000E40R2], got [%v]", upsIdentity.Name)
	}

	if upsIdentity.AttachedDevices != "Attached Devices not set" {
		t.Fatalf("Expected upsIdentity.AttachedDevices [Attached Devices not set], got [%v]", upsIdentity.AttachedDevices)
	}

	// Call the ups battery table device enumerator.
	upsBatteryTable := testUpsMib.UpsBatteryTable
	fmt.Printf("ZZZ (Test): Calling UpsBatteryTable device enumerator (Third try)\n")
	fmt.Printf("This is the one that works!!!\n")
	devices, err := upsBatteryTable.SnmpTable.DevEnumerator.DeviceEnumerator(nil)
	fmt.Printf("ZZZ (Test): Called UpsBatteryTable device enumerator. devices %+v, err %x\n", devices, err)
	// Ensure devices and no error.
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) == 0 { // len(nil) is defined as zero in golang.
		// It's really easy to mess up the enumeration call and call the default,
		// so test that we don't do that.
		t.Fatalf("Expected devices, got none.\n")
	}

	// Enumerate UpsInputTable devices.
	upsInputTable := testUpsMib.UpsInputTable
	devices, err = upsInputTable.SnmpTable.DevEnumerator.DeviceEnumerator(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) == 0 {
		t.Fatalf("Expected devices, got none.\n")
	}

	// Enumerate the mib.
	devices, err = testUpsMib.EnumerateDevices(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) == 0 {
		t.Fatalf("Expected devices, got none.\n")
	}
	if len(devices) != 40 {
		t.Fatalf("Expected 40 devices, got %d.\n", len(devices))
	}

	fmt.Printf("Dumping devices enumerated from UPS-MIB\n")
	for i := 0; i < len(devices); i++ {
		fmt.Printf("UPS-MIB device[%d]: %v %v %v %v %v row:%v column:%v\n", i,
			devices[i].Data["table_name"],
			devices[i].Type,
			devices[i].Data["info"],
			devices[i].Data["oid"],
			devices[i].Data["base_oid"],
			devices[i].Data["row"],
			devices[i].Data["column"])
	}

	fmt.Printf("Dumping full first device data: %v\n", devices[0].Data)
	fmt.Printf("Dumping second first device data: %v\n", devices[1].Data)

	t.Logf("TestUpsMib end")
}
