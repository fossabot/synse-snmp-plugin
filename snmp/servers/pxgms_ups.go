package servers

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/mibs/ups_mib"
)

// This file supports the PXGMS UPS + EATON 93PM SNMP Server.

//// ParseDeviceConfigs is a wrapper around config.ParseDeviceConfig() that
//// takes a directory parameter for sanity.
//func ParseDeviceConfigs(deviceDirectory string) (
//	deviceConfigs []*config.DeviceConfig, err error) {
//
//	// Set EnvDevicePath.
//	err = os.Setenv(config.EnvDevicePath, deviceDirectory)
//	if err != nil {
//		return nil, err
//	}
//	// Unset env on exit.
//	defer func() {
//		_ = os.Unsetenv(config.EnvDevicePath)
//	}()
//
//	// Parse the Device configuration.
//	deviceConfigs, err = config.ParseDeviceConfig()
//	if err != nil {
//		return nil, err
//	}
//	for i := 0; i < len(deviceConfigs); i++ {
//		fmt.Printf("deviceConfigs[%d]: %+v\n", i, deviceConfigs[i])
//	}
//	return deviceConfigs, nil
//}

// PxgmsUps represents the PXGMS UPS + EATON 93PM SNMP Server.
type PxgmsUps struct {
	*core.SnmpServerBase                        // base class.
	UpsMib               *mibs.UpsMib           // Supported Mibs.
	DeviceConfigs        []*config.DeviceConfig // Enumerated device configs.
}

// NewPxgmsUps creates the PxgmsUps structure.
// func NewPxgmsUps(configFilePath string) (ups *PxgmsUps, err error) {
// TODO: Name sucks. Should initialize all. (Upses?)
func NewPxgmsUps(data map[string]interface{}) (ups *PxgmsUps, err error) { // nolint: gocyclo

	// FIXME (etd): for now, commenting this out -- with the auto-enum config passed
	// in as a param here, we don't need to parse the device config to get that info

	//// Load the device configs.
	//deviceConfigs, err := ParseDeviceConfigs("../../config/device")
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Dump them for now.
	//fmt.Printf("Found device %d configs.\n", len(deviceConfigs))
	//for i := 0; i < len(deviceConfigs); i++ {
	//	fmt.Printf("deviceConfig[%d] %T: %+v\n", i, deviceConfigs[i], deviceConfigs[i])
	//}
	//
	//// Find the right one. For now there is only one.
	//upsDeviceConfigs, err := devices.FindDeviceConfigsByModel(deviceConfigs, "PXGMS_UPS") // TODO: String should be PXGMS UPS, not PXGMS_UPS if possible.
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Dump the upses.
	//fmt.Printf("Found %d upsDeviceConfigs.\n", len(upsDeviceConfigs))
	//for i := 0; i < len(upsDeviceConfigs); i++ {
	//	fmt.Printf("upses[%d] %T: %+v\n", i, upsDeviceConfigs[i], upsDeviceConfigs[i])
	//}
	//
	//// For each ups:
	//// TODO: Only support one UPS for now. We can change this in the future.
	//if len(upsDeviceConfigs) != 1 {
	//	// If not one (what we currently support) then there should be none
	//	// configured and this container should not be running.
	//	return nil, fmt.Errorf("No ups device config found")
	//}
	//
	//// Get our UPS device config.
	//upsDeviceConfig := upsDeviceConfigs[0]
	//fmt.Printf("upsDeviceConfig: %+v\n", upsDeviceConfig)

	// FIXME (etd): sorta a hack just to get things moving, but adding in a check against
	// the model here. There could probably be something at a higher level that checks this
	// and initializes the right stuff based on the specified model.
	model := data["model"].(string)
	if model != "PXGMS_UPS" {
		return nil, fmt.Errorf("only PXGMS_UPS models are currently supported")
	}

	// the autoenum config is map[string]interface{}, but GetDeviceConfig requires map[string]string
	// so we need to convert the autoenum type
	tmpMap := map[string]string{}
	for k, v := range data {
		tmpMap[k] = fmt.Sprint(v)
	}

	// Create the SNMP DeviceConfig,
	snmpDeviceConfig, err := core.GetDeviceConfig(tmpMap)
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

	// Output enumerated devices.
	for i := 0; i < len(snmpDevices); i++ {
		fmt.Printf("snmpDevice[%d]: %+v\n", i, snmpDevices[i])
	}

	// Set up the object.
	return &PxgmsUps{
		SnmpServerBase: snmpServerBase,
		UpsMib:         upsMib,
		DeviceConfigs:  snmpDevices,
	}, nil
}
