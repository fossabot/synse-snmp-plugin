package devices

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	//"github.com/vapor-ware/synse-sdk/sdk"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/mibs/ups_mib"
)

// FindDevicesConfigsByType returns all elements in a DeviceConfig array where the Type is t.
// TODO: Could some of these be SDK helper functions? No idea if the answer is yes or no.
func FindDeviceConfigsByType(devices []*config.DeviceConfig, t string) (matches []*config.DeviceConfig, err error) {
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
func DumpDeviceConfigs(devices []*config.DeviceConfig, header string) {
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

// ParseProtoypeConfigs is a wrapper around config.ParsePrototyeConfig() that
// takes a directory parameter for sanity.
func ParsePrototypeConfigs(prototypeDirectory string) (prototypeConfigs []*config.PrototypeConfig, err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fmt.Printf("pwd is: %v\n", pwd)
	fmt.Printf("prototypeDirectory is: %v\n", prototypeDirectory)

	// ls in the correct directory.
	fmt.Printf("ls %v:\n", prototypeDirectory)
	files, err := ioutil.ReadDir(prototypeDirectory)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}

	// Set EnvProtoPath.
	err = os.Setenv(config.EnvProtoPath, prototypeDirectory)
	if err != nil {
		return nil, err
	}
	// Unset env on exit.
	defer func() {
		_ = os.Unsetenv(config.EnvProtoPath)
	}()

	// Parse the Protoype configuration.
	prototypeConfigs, err = config.ParsePrototypeConfig()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(prototypeConfigs); i++ {
		fmt.Printf("prototypeConfigs[%d]: %+v\n", i, prototypeConfigs[i])
	}
	return prototypeConfigs, nil
}

// FindPrototypeConfigByType finds a prototype config in the given set where Type matches t or nil if not found.
func FindPrototypeConfigByType(prototypeConfigs []*config.PrototypeConfig, t string) (prototypeConfig *config.PrototypeConfig) {
	if prototypeConfigs == nil {
		return nil
	}
	for i := 0; i < len(prototypeConfigs); i++ {
		if prototypeConfigs[i].Type == t {
			return prototypeConfigs[i]
		}
	}
	return nil
}

// Create Device creates the Device structure in test land for now.
// Make your own Devices! (DeviceConfig is dynamic with SNMP)
// TODO:
func CreateDevice(
	deviceConfig *config.DeviceConfig,
	prototypeConfig *config.PrototypeConfig,
	deviceHandler *sdk.DeviceHandler,
	plugin *sdk.Plugin) (device *sdk.Device, err error) {

	return sdk.NewDevice(
		prototypeConfig,
		deviceConfig,
		deviceHandler,
		plugin)
}

/*
// ParseDeviceConfigs is a wrapper around config.ParseDeviceConfig() that
// takes a directory parameter for sanity.
func ParseDeviceConfigs(deviceDirectory string) (deviceConfigs []*config.DeviceConfig, err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fmt.Printf("pwd is: %v\n", pwd)
	fmt.Printf("deviceDirectory is: %v\n", deviceDirectory)

	// ls in the correct directory.
	fmt.Printf("ls %v:\n", deviceDirectory)
	files, err := ioutil.ReadDir(deviceDirectory)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}

	// Set EnvProtoPath.
	err = os.Setenv(config.EnvDevicePath, deviceDirectory)
	if err != nil {
		return nil, err
	}
	// Unset env on exit.
	defer func() {
		_ = os.Unsetenv(config.EnvProtoPath)
	}()

	// Parse the Protoype configuration.
	deviceConfigs, err = config.ParseDeviceConfig()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(deviceConfigs); i++ {
		fmt.Printf("deviceConfigs[%d]: %+v\n", i, deviceConfigs[i])
	}
	return deviceConfigs, nil
}

// FindDeviceConfigByType finds a device config in the given set where Type matches t or nil if not found.
func FindDeviceConfigByType(deviceConfigs []*config.DeviceConfig, t string) (deviceConfig *config.DeviceConfig) {
	if deviceConfigs == nil {
		return nil
	}
	for i := 0; i < len(deviceConfigs); i++ {
		if deviceConfigs[i].Type == t {
			return deviceConfigs[i]
		}
	}
	return nil
}
*/

// TODO: Explain Why? This is needed, but why?
func testDeviceIdentifier(x map[string]string) string {
	return ""
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

	// Create a snmp config.
	snmpConfig, err := core.NewDeviceConfig(
		"v3",        // SNMP v3
		"127.0.0.1", // Endpoint
		1024,        // Port
		securityParameters,
		"public") //  Context name
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	t.Logf("snmpConfig: %+v", snmpConfig)

	// Create a client.
	client, err := core.NewSnmpClient(snmpConfig)
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	t.Logf("client: %+v", client)

	// Create SnmpServerBase
	snmpServer, err := core.NewSnmpServerBase(
		client,
		snmpConfig,
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

	powerDeviceConfigs, err := FindDeviceConfigsByType(devices, "power")
	if err != nil {
		t.Fatal(err)
	}

	DumpDeviceConfigs(powerDeviceConfigs, "Power device configs")

	// Prototype configs are in ${PWD}/../config/proto
	// In order to parse them, we need to set environment variable EnvProtoPath to the directory which is really funky.
	// Why not just pass in the directory as a parameter?
	prototypeConfigs, err := ParsePrototypeConfigs("../config/proto")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("prototypeConfigs: %+v\n", prototypeConfigs)

	powerPrototypeConfig := FindPrototypeConfigByType(prototypeConfigs, "power")
	fmt.Printf("powerPrototypeConfig: %+v\n", powerPrototypeConfig)

	powerDeviceHandler := &SnmpPower
	fmt.Printf("powerDeviceHandler: %+v\n", powerDeviceHandler)

	// Need handlers to create a plugin.
	handlers, err := sdk.NewHandlers(testDeviceIdentifier, testUpsMib.EnumerateDevices)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("handlers: %+v\n", handlers)

	// Need a plugin config to create a plugin.
	pluginConfig := config.PluginConfig{
		Name:    "test config",
		Version: "test config v1",
		Network: config.NetworkSettings{
			Type:    "tcp",
			Address: "test_config",
		},
		Settings: config.Settings{
			Read:        config.ReadSettings{Buffer: 1024},
			Write:       config.WriteSettings{Buffer: 1024},
			Transaction: config.TransactionSettings{TTL: "2s"},
		},
	}
	fmt.Printf("pluginConfig: %+v\n", pluginConfig)

	// Create a plugin.
	plugin, err := sdk.NewPlugin(handlers, &pluginConfig)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("plugin: %+v\n", plugin)

	// At long last we should be able to create the Device structure.
	powerDevice, err := CreateDevice(powerDeviceConfigs[0], powerPrototypeConfig, powerDeviceHandler, plugin)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("powerDevice: %+v\n", powerDevice)

	readings, err := SnmpPowerRead(powerDevice)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Power Reading: %T, %+v\n", readings, readings)
	for i := 0; i < len(readings); i++ {
		fmt.Printf("Reading[%d]: %T, %+v\n", i, readings[i], readings[i])
	}

}
