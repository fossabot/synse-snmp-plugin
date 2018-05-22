package devices

import (
	"fmt"
	"testing"

	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/mibs/ups_mib"
)

// FindDevicesByType returns all elements in a DeviceConfig array where the Type is t.
// TODO: Could be an SDK helper function?
func FindDevicesByType(devices []*config.DeviceConfig, t string) (matches []*config.DeviceConfig, err error) {
	if devices == nil {
		return nil, fmt.Errorf("devices is nil")
	}

	for i := 0; i < len(devices); i++ {
		if devices[i].Type == t {
			matches = append(matches, devices[i])
		}
	}
	return matches, err
}

// DumpDevices utility function.
func DumpDevices(devices []*config.DeviceConfig, header string) {
	fmt.Printf("Dumping Devices. ")
	fmt.Print(header)

	if devices == nil {
		fmt.Printf(" <nil>\n")
		return
	}

	fmt.Printf(". Count: %d\n", len(devices))

	for i := 0; i < len(devices); i++ {
		fmt.Printf("device[%d]: %v %v %v %v %v row:%v column:%v\n", i,
			devices[i].Data["table_name"],
			devices[i].Type,
			devices[i].Data["info"],
			devices[i].Data["oid"],
			devices[i].Data["base_oid"],
			devices[i].Data["row"],
			devices[i].Data["column"])
	}
}

// Initial device test. Ensure we can register each type the ups mib supports
// and get a reading from each.
func TestDevices(t *testing.T) { // nolint: gocyclo
	t.Logf("TestDevices")

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
	testUpsMib, err := mibs.NewUpsMib(snmpServer)
	if err != nil {
		t.Fatal(err) // Fail the test.
	}
	fmt.Printf("TestUpsMib created mib\n")

	// Enumerate the mib.
	devices, err := testUpsMib.EnumerateDevices(nil)
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
	fmt.Printf("\n")

	// TODO: Find all power devices. Get readings.
	// func FindDevicesByType(devices []*config.DeviceConfig, t string) (matches []*config.DeviceConfig, err error) {
	powerDevices, err := FindDevicesByType(devices, "power")
	if err != nil {
		t.Fatal(err)
	}

	DumpDevices(powerDevices, "Power devices")

	// TODO: Get readings for the power devices we've found.

	// TODO: Find all voltage devices. Get readings.
	// TODO: Find all current devices. Get readings.
	// TODO: Find all frequency devices. Get readings.
	// TODO: Find all temperature devices. Get readings.
	// TODO: Find all identity devices. Get readings.
	// TODO: Find all status devices. Get readings.

}
