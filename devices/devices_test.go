package devices

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	//"github.com/vapor-ware/synse-sdk/sdk"
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

	powerDevices, err := FindDeviceConfigsByType(devices, "power")
	if err != nil {
		t.Fatal(err)
	}

	DumpDeviceConfigs(powerDevices, "Power device configs")

	// Prototype configs are in ${PWD}/../config/proto
	// In order to parse them, we need to set environment variable EnvProtoPath to the directory which is really funky.
	// Why not just pass in the directory as a parameter?
	prototypeConfigs, err := ParsePrototypeConfigs("../config/proto")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("prototypeConfigs: %+v\n", prototypeConfigs)

	powerPrototype := FindPrototypeConfigByType(prototypeConfigs, "power")
	fmt.Printf("powerPrototype: %+v\n", powerPrototype)

	/*
		deviceConfigs, err := ParseDeviceConfigs(".")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("deviceConfigs: %+v\n", deviceConfigs)

		powerDevice := FindDeviceConfigByType(deviceConfigs, "power")
		fmt.Printf("powerDevice: %+v\n", powerDevice)
	*/

	/*
			// TODO: Get readings for the power devices we've found.
			for i := 0; i < len(powerDevices); i++ {
				//fmt.Printf("Reading power device [%d], Info: %v\n", i, powerDevices[i].Info)
				//fmt.Printf("Reading power device [%d], Data[Info]: %v\n", i, powerDevices[i].Data["Info"])
				fmt.Printf("Reading power device [%d], Data[info]: %v, Type: %v\n", i, powerDevices[i].Data["info"], powerDevices[i].Type)
				// (variable of type *github.com/vapor-ware/synse-snmp-plugin/vendor/github.com/vapor-ware/synse-sdk/sdk/config.DeviceConfig)
				// as *github.com/vapor-ware/synse-snmp-plugin/vendor/github.com/vapor-ware/synse-sdk/sdk.Device value in argument to SnmpPowerRead (varcheck)

				// TODO: KILLER: Convert to synse-sdk.sdk.config.DeviceConfig
				// TODO: Quick and shitty way is to do it yourself.
				// TODO: Better way is to use the sdk, but that may have pitfalls based on experience.

				// To create a device, we need:
				// ProtoConfig
				// DeviceConfig
				// DeviceHandler
				// PlugIn

		    /*
				// Create a prototye config from the appropriate device using the appropriate proto file.
				//func loadPrototypeFile(directory string, deviceType string) (prototypeConfig *config.PrototypeConfig, err error) {
				prototypeConfig, err := loadPrototypeFile("../config/proto/", powerDevices[i].Type)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Printf("prototypeConfig: %+v\n", prototypeConfig)
	*/

	// In order to parse prototype files, we need to set an environment variable which is really really weird.
	// Passing in a parameter is more typical.

	/*
		readings, err := SnmpPowerRead(powerDevices[i])
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("Readings: %+v\n", readings)
	*/
	/*
		}
	*/
	// TODO: Find all voltage devices. Get readings.
	// TODO: Find all current devices. Get readings.
	// TODO: Find all frequency devices. Get readings.
	// TODO: Find all temperature devices. Get readings.
	// TODO: Find all identity devices. Get readings.
	// TODO: Find all status devices. Get readings.

}
