package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/soniah/gosnmp"
)

// AuthenticationProtocol enumeration for authentication algorithms.
type AuthenticationProtocol uint8

// AuthProtocol enumeration for authentication algorithms.
const (
	NoAuthentication AuthenticationProtocol = 1
	MD5              AuthenticationProtocol = 2
	SHA              AuthenticationProtocol = 3
)

// PrivacyProtocol enumeration for encryption algorithms.
type PrivacyProtocol uint8

// PrivacyProtocol enumeration for encryption algorithms.
const (
	NoPrivacy PrivacyProtocol = 1
	DES       PrivacyProtocol = 2
	AES       PrivacyProtocol = 3
)

// SecurityParameters is a subset of SNMP USM parameters.
type SecurityParameters struct {
	AuthenticationProtocol   AuthenticationProtocol
	PrivacyProtocol          PrivacyProtocol
	UserName                 string // SNMP user name.
	AuthenticationPassphrase string
	PrivacyPassphrase        string
}

// NewSecurityParameters constructs a SecurityParameters.
func NewSecurityParameters(
	userName string,
	authenticationProtocol AuthenticationProtocol,
	authenticationPassphrase string,
	privacyProtocol PrivacyProtocol,
	privacyPassphrase string) (*SecurityParameters, error) {

	// For now, require  authorization and privacy.
	// Empty user/passwords are okay.
	if !(authenticationProtocol == MD5 || authenticationProtocol == SHA) {
		return nil, fmt.Errorf("Unsupported authentication protocol [%v]",
			authenticationProtocol)
	}

	if !(privacyProtocol == DES || privacyProtocol == AES) {
		return nil, fmt.Errorf("Unsupported privacy protocol [%v]",
			privacyProtocol)
	}

	return &SecurityParameters{
		UserName:                 userName,
		AuthenticationProtocol:   authenticationProtocol,
		AuthenticationPassphrase: authenticationPassphrase,
		PrivacyProtocol:          privacyProtocol,
		PrivacyPassphrase:        privacyPassphrase,
	}, nil
}

// DeviceConfig is a thin wrapper around the configuration for gosnmp using SNMP V3.
type DeviceConfig struct {
	Version            string              // SNMP protocol version. Currently only SNMP V3 is supported.
	Endpoint           string              // Endpoint of the SNMP server to connect to.
	ContextName        string              // Context name for SNMP V3 messages.
	Timeout            time.Duration       // Timeout for the SNMP query.
	SecurityParameters *SecurityParameters // SNMP V3 security parameters.
	//SecurityLevel      SecurityLevel       // SecurityLevel for authentication and privacy.
	Port uint16 // UDP port to connect to.
}

// checkForEmptyString checks for an empty string variable and fails with an
// attempt of a reasonable error message on failure.
func checkForEmptyString(variable string, variableName string) (err error) {
	if variable == "" {
		return fmt.Errorf("%v is an empty string, but should not be", variableName)
	}
	return nil
}

// NewDeviceConfig creates an DeviceConfig.
func NewDeviceConfig(
	version string,
	endpoint string,
	port uint16,
	//securityLevel SecurityLevel,
	securityParameters *SecurityParameters,
	contextName string) (*DeviceConfig, error) {

	// Check parameters.
	versionUpper := strings.ToUpper(version)
	if versionUpper != "V3" {
		return nil, fmt.Errorf("Version [%v] unsupported", version)
	}

	if securityParameters == nil {
		return nil, fmt.Errorf("securityParameters is nil")
	}

	// Check strings for emptyness. Version is already checked.
	err := checkForEmptyString(endpoint, "endpoint")
	if err != nil {
		return nil, err
	}

	return &DeviceConfig{
		Version:            versionUpper,
		Endpoint:           endpoint,
		Port:               port,
		SecurityParameters: securityParameters,
		ContextName:        contextName,
		Timeout:            time.Duration(30) * time.Second,
	}, nil
}

// GetDeviceConfig takes the instance configuration for an SNMP device and
// parses it into a DeviceConfig struct, filling in default values for anything
// that is missing and has a default value defined.
func GetDeviceConfig(instanceData map[string]string) (*DeviceConfig, error) {

	// Parse out each field. The contructor call will check the parameters.
	version := instanceData["version"]
	endpoint := instanceData["endpoint"]

	prt, err := strconv.Atoi(instanceData["port"])
	if err != nil {
		return nil, err
	}
	port := uint16(prt)

	userName := instanceData["userName"]

	authenticationProtocolString := strings.ToUpper(instanceData["authenticationProtocol"])
	// Only MD5 and SHA are currently supported.
	var authenticationProtocol AuthenticationProtocol
	if authenticationProtocolString == "MD5" {
		authenticationProtocol = MD5
	} else if authenticationProtocolString == "SHA" {
		authenticationProtocol = SHA
	} else {
		return nil, fmt.Errorf("Unsupported authentication protocol [%v]", authenticationProtocolString)
	}

	authenticationPassphrase := instanceData["authenticationPassphrase"]

	privacyProtocolString := strings.ToUpper(instanceData["privacyProtocol"])
	// Only DES and AES are currently supported.
	var privacyProtocol PrivacyProtocol
	if privacyProtocolString == "DES" {
		privacyProtocol = DES
	} else if privacyProtocolString == "AES" {
		privacyProtocol = AES
	} else {
		return nil, fmt.Errorf("Unsupported privacy protocol [%v]", privacyProtocolString)
	}

	privacyPassphrase := instanceData["privacyPassphrase"]
	contextName := instanceData["contextName"]

	// Create security parameters
	securityParameters, err := NewSecurityParameters(
		userName,
		authenticationProtocol,
		authenticationPassphrase,
		privacyProtocol,
		privacyPassphrase)
	if err != nil {
		return nil, err
	}

	// Create the config.
	return NewDeviceConfig(
		version,
		endpoint,
		port,
		securityParameters,
		contextName)
}

// SnmpClient is a thin wrapper around gosnmp.
type SnmpClient struct {
	Config *DeviceConfig
}

// NewSnmpClient constructs SnmpClient.
func NewSnmpClient(config *DeviceConfig) (*SnmpClient, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	return &SnmpClient{
		Config: config,
	}, nil
}

// ReadResult is the structure for any SNMP read.
type ReadResult struct {
	Oid  string      // The SNMP OID read.
	Data interface{} // The data for the OID. See gosnmp decodeValue() https://github.com/soniah/gosnmp/blob/master/helper.go#L67
}

// Walk performs an SNMP bulk walk.
func (client *SnmpClient) Walk(rootOid string) (results []ReadResult, err error) {

	goSnmp, err := client.createGoSNMP()
	if err != nil {
		return nil, err
	}

	resultSet, err := goSnmp.BulkWalkAll(rootOid)
	err2 := goSnmp.Conn.Close() // Do not leak connection.

	// Return first error.
	if err != nil {
		return nil, err
	}
	if err2 != nil {
		return nil, err2
	}

	// Package results.
	for _, snmpPdu := range resultSet {

		// If it looks like an ASCII string, try to translate it.
		if snmpPdu.Type == gosnmp.OctetString {
			ascii, err := TranslatePrintableASCII(snmpPdu.Value)
			if err == nil {
				snmpPdu.Value = ascii
			}
			// err above is deliberately ignored here. SNMP does not differentiate
			// between ASCII strings and byte array.
		}

		results = append(results, ReadResult{
			Oid:  snmpPdu.Name,
			Data: snmpPdu.Value,
		})
	}
	return results, nil
}

// createGoSNMP is a helper to create gosnmp.GoSNMP from SnmpClient.
// On success, the connection is open.
func (client *SnmpClient) createGoSNMP() (*gosnmp.GoSNMP, error) {

	// Argument checks
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	// Map DeviceConfig parameters to gosnmp parameters.
	securityParameters := client.Config.SecurityParameters
	var authProtocol gosnmp.SnmpV3AuthProtocol
	var privProtocol gosnmp.SnmpV3PrivProtocol

	if securityParameters.AuthenticationProtocol == MD5 {
		authProtocol = gosnmp.MD5
	} else if securityParameters.AuthenticationProtocol == SHA {
		authProtocol = gosnmp.SHA
	} else {
		return nil, fmt.Errorf("Unsupported authentication protocol [%v]", securityParameters.AuthenticationProtocol)
	}

	if securityParameters.PrivacyProtocol == DES {
		privProtocol = gosnmp.DES
	} else if securityParameters.PrivacyProtocol == AES {
		privProtocol = gosnmp.AES
	} else {
		return nil, fmt.Errorf("Unsupported privacy protocol [%v]", securityParameters.PrivacyProtocol)
	}

	goSnmp := &gosnmp.GoSNMP{
		Target:        client.Config.Endpoint,
		Port:          client.Config.Port,
		Version:       gosnmp.Version3,
		Timeout:       client.Config.Timeout,
		SecurityModel: gosnmp.UserSecurityModel,
		MsgFlags:      gosnmp.AuthPriv,
		SecurityParameters: &gosnmp.UsmSecurityParameters{
			UserName:                 client.Config.SecurityParameters.UserName,
			AuthenticationProtocol:   authProtocol,
			AuthenticationPassphrase: client.Config.SecurityParameters.AuthenticationPassphrase,
			PrivacyProtocol:          privProtocol,
			PrivacyPassphrase:        client.Config.SecurityParameters.PrivacyPassphrase,
		},
		ContextName: client.Config.ContextName,
	}

	// Connect
	err := goSnmp.Connect()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect: %+v", goSnmp)
	}
	return goSnmp, err
}