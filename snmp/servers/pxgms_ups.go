package servers

import (
	"fmt"
	"os"

	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/devices"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
	mibs "github.com/vapor-ware/synse-snmp-plugin/snmp/mibs/ups_mib"
	//"github.com/vapor-ware/synse-snmp-plugin/snmp/mibs/ups_mib"
)

// This file supports the PXGMS UPS + EATON 93PM SNMP Server.

// ParseDeviceConfigs is a wrapper around config.ParseDeviceConfig() that
// takes a directory parameter for sanity.
func ParseDeviceConfigs(deviceDirectory string) (
	deviceConfigs []*config.DeviceConfig, err error) {

	// Set EnvDevicePath.
	err = os.Setenv(config.EnvDevicePath, deviceDirectory)
	if err != nil {
		return nil, err
	}
	// Unset env on exit.
	defer func() {
		_ = os.Unsetenv(config.EnvDevicePath)
	}()

	// Parse the Device configuration.
	deviceConfigs, err = config.ParseDeviceConfig()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(deviceConfigs); i++ {
		fmt.Printf("deviceConfigs[%d]: %+v\n", i, deviceConfigs[i])
	}
	return deviceConfigs, nil
}

// PxgmsUps represents the PXGMS UPS + EATON 93PM SNMP Server.
type PxgmsUps struct {
	*core.SnmpServerBase // base class.
	// TODO: Mibs *[]core.Mib, // Pointers to the mibs that we support for this device. // TODO: Want the derrived class here.
	UpsMib *mibs.UpsMib // Supported Mibs.
	// TODO: Devices []*DeviceConfig // Supported devices.
	//Devices []*DeviceConfig // Supported devices.
	//DeviceConfigs []*core.DeviceConfig
	//DeviceConfigs []*sdk.DeviceConfig
	DeviceConfigs []*config.DeviceConfig
}

// NewPxgmsUps creates the PxgmsUps structure.
// func NewPxgmsUps(configFilePath string) (ups *PxgmsUps, err error) {
// TODO: Name sucks. Should initialize all. (Upses?)
func NewPxgmsUps() (ups *PxgmsUps, err error) { // nolint: gocyclo

	// Load the device configs.
	deviceConfigs, err := ParseDeviceConfigs("../../config/device")
	if err != nil {
		return nil, err
	}

	// Dump them for now.
	fmt.Printf("Found device %d configs.\n", len(deviceConfigs))
	for i := 0; i < len(deviceConfigs); i++ {
		fmt.Printf("deviceConfig[%d] %T: %+v\n", i, deviceConfigs[i], deviceConfigs[i])
	}

	// Find the right one. For now there is only one.
	//upses, err := FindDeviceConfigsByModel(deviceConfigs, "PXGMS_UPS") // TODO: String should be PXGMS UPS, not PXGMS_UPS if possible.
	upsDeviceConfigs, err := devices.FindDeviceConfigsByModel(deviceConfigs, "PXGMS_UPS") // TODO: String should be PXGMS UPS, not PXGMS_UPS if possible.
	if err != nil {
		return nil, err
	}

	// Dump the upses.
	fmt.Printf("Found %d upsDeviceConfigs.\n", len(upsDeviceConfigs))
	for i := 0; i < len(upsDeviceConfigs); i++ {
		fmt.Printf("upses[%d] %T: %+v\n", i, upsDeviceConfigs[i], upsDeviceConfigs[i])
	}

	// For each ups:
	// TODO: Only support one UPS for now. We can change this in the future.
	if len(upsDeviceConfigs) != 1 {
		// If not one (what we currently support) then there should be none
		// configured and this container should not be running.
		return nil, fmt.Errorf("No ups device config found")
	}

	// Get our UPS device config.
	upsDeviceConfig := upsDeviceConfigs[0]
	fmt.Printf("upsDeviceConfig: %+v\n", upsDeviceConfig)

	// Create the SNMP DeviceConfig,
	snmpDeviceConfig, err := core.GetDeviceConfig(upsDeviceConfig.Data)
	if err != nil {
		return nil, err
	}
	fmt.Printf("snmpDeviceConfig: %+v\n", snmpDeviceConfig)

	// Create SNMP client.
	snmpClient, err := core.NewSnmpClient(snmpDeviceConfig)
	if err != nil {
		return nil, err
	}
	fmt.Printf("snmpClient: %+v\n", snmpClient)

	// Create SnmpServerBase.
	snmpServerBase, err := core.NewSnmpServerBase(snmpClient, snmpDeviceConfig)
	if err != nil {
		return nil, err
	}
	fmt.Printf("snmpServerBase: %+v\n", snmpServerBase)

	// Create the UpsMib.
	upsMib, err := mibs.NewUpsMib(snmpServerBase)
	if err != nil {
		return nil, err
	}
	fmt.Printf("upsMib: %+v\n", upsMib)

	// Enumerate the mib.
	snmpDevices, err := upsMib.EnumerateDevices(
		map[string]interface{}{"rack": "test_rack", "board": "test_board"})
	if err != nil {
		return nil, err
	}

	//fmt.Printf("snmpDevices: %+v\n", snmpDevices)
	for i := 0; i < len(snmpDevices); i++ {
		fmt.Printf("snmpDevice[%d]: %+v\n", i, snmpDevices[i])
	}

	// TODO: Set up the object.
	return &PxgmsUps{
		SnmpServerBase: snmpServerBase,
		UpsMib:         upsMib,
		DeviceConfigs:  snmpDevices,
	}, nil

	//return nil, nil // TODO
}
