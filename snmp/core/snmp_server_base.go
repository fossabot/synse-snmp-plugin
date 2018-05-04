package core

import (
	"fmt"
)

// SnmpServerBase is a base class for specific SnmpServer implementations.
type SnmpServerBase struct {
	SnmpClient   *SnmpClient
	DeviceConfig *DeviceConfig
	RackID       string
	// TODO: What type is this? (Probably map[string]interface{} ???
	// ScanResultsInternal

	// TODO: This is going to be pretty rough.

	// TODO: Slice of SnmpTable as members?
}

// NewSnmpServerBase constructs common code for all SNMP Servers.
func NewSnmpServerBase(
	client *SnmpClient,
	deviceConfig *DeviceConfig,
	rackID string) (*SnmpServerBase, error) { // Dear golang. Your parser is the worst. Completely awful at best.
	// Parameter checks.
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	if deviceConfig == nil {
		return nil, fmt.Errorf("deviceConfig is nil")
	}

	if rackID == "" {
		return nil, fmt.Errorf("rackID is empty")
	}

	// Construct the struct.
	return &SnmpServerBase{
		SnmpClient:   client,
		DeviceConfig: deviceConfig,
		RackID:       rackID,
	}, nil
}

// TODO: Lots

// TODO: Not needed: func ConverSnmpResultSet(results []ReadResult) {}

// readRawRow
// TODO: func (snmpServerBase *SnmpServerBase) readRawRow(rawRow *SnmpRow)

// sortRawRow
// We need to sort the raw SNMP read data by OID in order to line up
// the columns with those in the table. ASCII sort is insufficient.
// table: The table we are operating on.
// snmp_row: The initial row from the scan results that we are reading.
// raw_row: The raw SNMP read from _read_raw_row.
// TODO: Does this work?
//func sortRawRow(table *SnmpTable, row *SnmpRow, rawRow *SnmpRow) (*SnmpRow) {
//func sortRawRow(table *SnmpTable, row *SnmpRow, rawRow *SnmpRow) interface{} {
//func sortRawRow(table *SnmpTable, row *SnmpRow, rawRow *SnmpRow) string {
/*
func sortRawRow(table *SnmpTable, row *SnmpRow, rawRow *SnmpRow) []string {

	logger.Debugf("sorting read_row")
	columnIndex := 1
	//var rowData []string //__ TODO: SnmpRow?? NO
	//var rowData interface{} // Whatever the data are.
	var rowData []string // SNMP OIDS

	baseOid := row.BaseOid

	for i := 0; i < len(table.ColumnList); i++ { // for each column
		dataOid := fmt.Sprintf(baseOid, columnIndex) // Get the OID for the cell TODO: Does this still work? Need to check the formatter.
		rowData = append(rowData, dataOid)
	}

	logger.Debugf("sortRawRow returning:")
	logger.Debugf("%v", rowData)
	return rowData
}
*/
