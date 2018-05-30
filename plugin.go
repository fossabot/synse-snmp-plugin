package main

import (
	"fmt"
	"log"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/devices"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/servers"
)

// Build time variables for setting the version info of a Plugin.
var (
	BuildDate     string
	GitCommit     string
	GitTag        string
	GoVersion     string
	VersionString string
)

// DeviceIdentifier defines the SNMP-specific way of uniquely identifying a device
// through its device configuration.
//
// FIXME - this is just a stub for framing up the plugin
// TODO: This will work for the initial cut. This may change later if/when
// we need to support the entity mib and entity sensor mib.
func DeviceIdentifier(data map[string]string) string {
	return data["oid"]
}

func main() {

	handlers, err := sdk.NewHandlers(DeviceIdentifier, nil)
	if err != nil {
		log.Fatal(err)
	}

	plugin, err := sdk.NewPlugin(handlers, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Load the MIB from the configuration still.
	pxgmsUps, err := servers.NewPxgmsUps()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Initialized PxgmsUps: %+v\n", pxgmsUps)

	// Register Device Handlers for all supported devices we interact with over SNMP.
	plugin.RegisterDeviceHandlers(
		&devices.SnmpCurrent,
		&devices.SnmpFrequency,
		&devices.SnmpIdentity,
		&devices.SnmpPower,
		&devices.SnmpStatus,
		&devices.SnmpTemperature,
		&devices.SnmpVoltage,
	)

	// Set build-time version info.
	plugin.SetVersion(sdk.VersionInfo{
		BuildDate:     BuildDate,
		GitCommit:     GitCommit,
		GitTag:        GitTag,
		GoVersion:     GoVersion,
		VersionString: VersionString,
	})

	// Run the plugin.
	err = plugin.Run()
	if err != nil {
		log.Fatal(err)
	}
}
