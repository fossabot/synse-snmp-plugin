package core

// TODO: Fix up comment. Run under docker.
// To run these tests, startup the emulator by:
// snmpsimd.py --data-dir=/Users/mhink/src/vapor/opendcre-core/opendcre_core/opendcre_southbound/emulator/snmp/data/rittal/initial/ups --agent-udpv4-endpoint=0.0.0.0:1024

import (
	"fmt"
	"testing"
	"time"
)

// TestClient is the initial positive test against the emulator.
func TestClient(t *testing.T) {
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

	// Walk OID "1.3.6.1" and print results.
	results, err := client.Walk("1.3.6.1")
	if err != nil {
		t.Error(err) // Fail the test.
	}

	// Log output.
	for i, result := range results {
		t.Logf("%d: OID: %v, Data: %v",
			i, result.Oid, result.Data)
	}
}

// getExpectedConfigShaAes gets an expected valid device configuration.
func getExpectedConfigShaAes() *DeviceConfig {
	securityParameters := &SecurityParameters{
		UserName:                 "simulator",
		AuthenticationProtocol:   SHA,
		AuthenticationPassphrase: "auctorias",
		PrivacyProtocol:          AES,
		PrivacyPassphrase:        "privatus",
	}

	return &DeviceConfig{
		Version:            "V3",
		Endpoint:           "127.0.0.1",
		Port:               1024,
		ContextName:        "public",
		Timeout:            time.Duration(30) * time.Second,
		SecurityParameters: securityParameters,
	}
}

// getExpectedConfigShaAes gets an expected valid device configuration.
func getExpectedConfigMd5Des() *DeviceConfig {
	securityParameters := &SecurityParameters{
		UserName:                 "simulator",
		AuthenticationProtocol:   MD5,
		AuthenticationPassphrase: "auctorias",
		PrivacyProtocol:          DES,
		PrivacyPassphrase:        "privatus",
	}

	return &DeviceConfig{
		Version:            "V3",
		Endpoint:           "127.0.0.1",
		Port:               1024,
		ContextName:        "public",
		Timeout:            time.Duration(30) * time.Second,
		SecurityParameters: securityParameters,
	}
}

// verifyConfig checks excpected DeviceConfig fields versus the actual.
func verifyConfig(expected *DeviceConfig, actual *DeviceConfig) (err error) {

	if expected.Version != actual.Version {
		return fmt.Errorf("Fail version: expected: [%v], actual [%v]",
			expected.Version, actual.Version)
	}

	if expected.Endpoint != actual.Endpoint {
		return fmt.Errorf("Fail endpoint: expected: [%v], actual [%v]",
			expected.Endpoint, actual.Endpoint)
	}

	if expected.Port != actual.Port {
		return fmt.Errorf("Fail port: expected: [%v], actual: [%v]",
			expected.Port, actual.Port)
	}

	if expected.SecurityParameters.UserName != actual.SecurityParameters.UserName {
		return fmt.Errorf("Fail userName: expected: [%v], actual [%v]",
			expected.SecurityParameters.UserName, actual.SecurityParameters.UserName)
	}

	if expected.SecurityParameters.AuthenticationProtocol != actual.SecurityParameters.AuthenticationProtocol {
		return fmt.Errorf("Fail AuthenticationProtocol: expected: [%v], actual [%v]",
			expected.SecurityParameters.AuthenticationProtocol,
			actual.SecurityParameters.AuthenticationProtocol)
	}

	if expected.SecurityParameters.AuthenticationPassphrase != actual.SecurityParameters.AuthenticationPassphrase {
		return fmt.Errorf("Fail AuthenticationPassphrase: expected: [%v], actual [%v]",
			expected.SecurityParameters.AuthenticationPassphrase,
			actual.SecurityParameters.AuthenticationPassphrase)
	}

	if expected.SecurityParameters.PrivacyProtocol != actual.SecurityParameters.PrivacyProtocol {
		return fmt.Errorf("Fail PrivacyProtocol: expected: [%v], actual [%v]",
			expected.SecurityParameters.PrivacyProtocol,
			actual.SecurityParameters.PrivacyProtocol)
	}

	if expected.SecurityParameters.PrivacyPassphrase != actual.SecurityParameters.PrivacyPassphrase {
		return fmt.Errorf("Fail PrivacyPassphrase: expected: [%v], actual [%v]",
			expected.SecurityParameters.PrivacyPassphrase,
			actual.SecurityParameters.PrivacyPassphrase)
	}

	if expected.ContextName != actual.ContextName {
		return fmt.Errorf("Fail ContextName: expected: [%v], actual [%v]",
			expected.ContextName, actual.ContextName)
	}
	return nil // Verification passed.
}

// Configuration tests.

// Test a valid configuration.
func TestValidConfigMapShaAes(t *testing.T) {
	yamlConfig := map[string]string{
		"version":                  "v3",
		"endpoint":                 "127.0.0.1",
		"port":                     "1024",
		"userName":                 "simulator",
		"authenticationProtocol":   "SHA",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "AES",
		"privacyPassphrase":        "privatus",
		"contextName":              "public",
	}
	// actual is a SNMP DeviceConfig created by the constructor.
	actual, err := GetDeviceConfig(yamlConfig)
	if err != nil {
		t.Fatal(err) // Fail test.
	}
	t.Logf("deviceConfig: %+v", actual)

	// Test each field against expected.
	expected := getExpectedConfigShaAes()
	err = verifyConfig(expected, actual)
	if err != nil {
		t.Fatal(err) // Fail test.
	}
}

func TestValidConfigMapMd5Des(t *testing.T) {
	yamlConfig := map[string]string{
		"version":                  "v3",
		"endpoint":                 "127.0.0.1",
		"port":                     "1024",
		"userName":                 "simulator",
		"authenticationProtocol":   "MD5",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "DES",
		"privacyPassphrase":        "privatus",
		"contextName":              "public",
	}
	// actual is a SNMP DeviceConfig created by the constructor.
	actual, err := GetDeviceConfig(yamlConfig)
	if err != nil {
		t.Fatal(err) // Fail test.
	}
	t.Logf("deviceConfig: %+v", actual)

	// Test each field against expected.
	expected := getExpectedConfigMd5Des()
	err = verifyConfig(expected, actual)
	if err != nil {
		t.Fatal(err) // Fail test.
	}
}

// Test invalid configurations, one for each required field missing or invalid.

func TestConfigMapInvalidVersion(t *testing.T) {
	yamlConfig := map[string]string{
		"version":                  "v2c", // v2c is not currently supported.
		"endpoint":                 "127.0.0.1",
		"port":                     "1024",
		"userName":                 "simulator",
		"authenticationProtocol":   "SHA",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "AES",
		"privacyPassphrase":        "privatus",
		"contextName":              "public",
	}
	_, err := GetDeviceConfig(yamlConfig)
	if err != nil {
		// We expect an error here: Version [v2c] unsupported
		t.Logf("got error: %v", err)
		expectedError := "Version [v2c] unsupported"
		if err.Error() != expectedError {
			t.Fatalf("Expected error %v, got %v", expectedError, err.Error())
		}
	} else {
		t.Fatal("Got nil error, expected non-nil error")
	}
}

func TestConfigMapForgotVersion(t *testing.T) {
	yamlConfig := map[string]string{
		//"version":                  "v3",
		"endpoint":                 "127.0.0.1",
		"port":                     "1024",
		"userName":                 "simulator",
		"authenticationProtocol":   "SHA",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "AES",
		"privacyPassphrase":        "privatus",
		"contextName":              "public",
	}
	_, err := GetDeviceConfig(yamlConfig)
	if err != nil {
		// We expect an error here: Version [] unsupported
		t.Logf("got error: %v", err)
		expectedError := "Version [] unsupported"
		if err.Error() != expectedError {
			t.Fatalf("Expected error %v, got %v", expectedError, err.Error())
		}
	} else {
		t.Fatal("Got nil error, expected non-nil error")
	}
}

func TestConfigMapForgotEndpoint(t *testing.T) {
	yamlConfig := map[string]string{
		"version": "v3",
		//"endpoint":                 "127.0.0.1",
		"port":                     "1024",
		"userName":                 "simulator",
		"authenticationProtocol":   "SHA",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "AES",
		"privacyPassphrase":        "privatus",
		"contextName":              "public",
	}
	_, err := GetDeviceConfig(yamlConfig)
	if err != nil {
		t.Logf("got error: %v", err)
		// We expect an error here:
		expectedError := "endpoint is an empty string, but should not be"
		if err.Error() != expectedError {
			t.Fatalf("Expected error %v, got %v", expectedError, err.Error())
		}
	} else {
		t.Fatalf("Got nil error, expected non-nil error")
	}
}

func TestConfigMapNonNumericPort(t *testing.T) {
	yamlConfig := map[string]string{
		"version":                  "v3",
		"endpoint":                 "127.0.0.1",
		"port":                     "Z",
		"userName":                 "simulator",
		"authenticationProtocol":   "SHA",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "AES",
		"privacyPassphrase":        "privatus",
		"contextName":              "public",
	}
	_, err := GetDeviceConfig(yamlConfig)
	if err != nil {
		t.Logf("got error: %v", err)
		// We expect an error here:
		expectedError := "strconv.Atoi: parsing \"Z\": invalid syntax"
		if err.Error() != expectedError {
			t.Fatalf("Expected error %v, got %v", expectedError, err.Error())
		}
	} else {
		t.Fatal("Got nil error, expected non-nil error")
	}
}

func TestConfigMapForgotPort(t *testing.T) {
	yamlConfig := map[string]string{
		"version":  "v3",
		"endpoint": "127.0.0.1",
		//"port":                     "Z",
		"userName":                 "simulator",
		"authenticationProtocol":   "SHA",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "AES",
		"privacyPassphrase":        "privatus",
		"contextName":              "public",
	}
	_, err := GetDeviceConfig(yamlConfig)
	if err != nil {
		t.Logf("got error: %v", err)
		// We expect an error here:
		expectedError := "strconv.Atoi: parsing \"\": invalid syntax"
		if err.Error() != expectedError {
			t.Fatalf("Expected error %v, got %v", expectedError, err.Error())
		}
	} else {
		t.Fatal("Got nil error, expected non-nil error")
	}
}

func TestConfigMapForgotContextName(t *testing.T) {
	// This one should be okay. Empty string is valid for context name and is in
	// fact the SNMP default.
	yamlConfig := map[string]string{
		"version":                  "v3",
		"endpoint":                 "127.0.0.1",
		"port":                     "1024",
		"userName":                 "simulator",
		"authenticationProtocol":   "SHA",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "AES",
		"privacyPassphrase":        "privatus",
	}
	_, err := GetDeviceConfig(yamlConfig)
	if err != nil {
		t.Logf("got error: %v", err)
		// We should not expect an error here:
		t.Fatalf("Expected no error, got %v", err.Error())
	}
}