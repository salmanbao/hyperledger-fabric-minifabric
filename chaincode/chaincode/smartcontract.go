package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// Smart contract
type SmartContract struct {
	contractapi.Contract
}

// Device represents an IoT device
type Device struct {
	ID       string `json:"id"`
	Owner    string `json:"owner"`
	Location string `json:"location"`
	Status   string `json:"status"` // e.g., active, inactive
}

// DataRecord represents IoT data from a device
type DataRecord struct {
	DeviceID   string `json:"deviceID"`
	Timestamp  string `json:"timestamp"`
	Data       string `json:"data"`
	Status     string `json:"status"` // e.g., pending, verified, rejected
	VerifierID string `json:"verifierID,omitempty"`
}

// RegisterDevice registers a new IoT device on the ledger
func (s *SmartContract) RegisterDevice(ctx contractapi.TransactionContextInterface, deviceID string, owner string, location string) error {
	// Check if device already exists
	exists, err := s.DeviceExists(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to check device existence: %v", err)
	}
	if exists {
		return fmt.Errorf("device %s already registered", deviceID)
	}

	// Create new device
	device := Device{
		ID:       deviceID,
		Owner:    owner,
		Location: location,
		Status:   "active",
	}

	// Marshal device data and save to ledger
	deviceJSON, err := json.Marshal(device)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(deviceID, deviceJSON)
}

// DeviceExists checks if a device exists on the ledger
func (s *SmartContract) DeviceExists(ctx contractapi.TransactionContextInterface, deviceID string) (bool, error) {
	deviceJSON, err := ctx.GetStub().GetState(deviceID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return deviceJSON != nil, nil
}

// SubmitData allows an IoT device to submit data to the ledger
func (s *SmartContract) SubmitData(ctx contractapi.TransactionContextInterface, deviceID string, timestamp string, data string) error {
	// Ensure device is registered
	exists, err := s.DeviceExists(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to check device existence: %v", err)
	}
	if !exists {
		return fmt.Errorf("device %s not registered", deviceID)
	}

	// Create data record
	dataRecord := DataRecord{
		DeviceID:  deviceID,
		Timestamp: timestamp,
		Data:      data,
		Status:    "pending",
	}

	// Create composite key for data record
	dataKey, err := ctx.GetStub().CreateCompositeKey("DataRecord", []string{deviceID, timestamp})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	// Marshal data record and save to ledger
	dataJSON, err := json.Marshal(dataRecord)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(dataKey, dataJSON)
}

// VerifyData verifies a submitted data record and updates its status
func (s *SmartContract) VerifyData(ctx contractapi.TransactionContextInterface, deviceID string, timestamp string, verifierID string, isValid bool) error {
	// Retrieve data record from ledger
	dataKey, err := ctx.GetStub().CreateCompositeKey("DataRecord", []string{deviceID, timestamp})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	dataJSON, err := ctx.GetStub().GetState(dataKey)
	if err != nil {
		return fmt.Errorf("failed to get data record: %v", err)
	}
	if dataJSON == nil {
		return fmt.Errorf("data record for device %s at %s does not exist", deviceID, timestamp)
	}

	// Unmarshal data record
	var dataRecord DataRecord
	err = json.Unmarshal(dataJSON, &dataRecord)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data record: %v", err)
	}

	// Update verification status and verifierID
	if isValid {
		dataRecord.Status = "verified"
	} else {
		dataRecord.Status = "rejected"
	}
	dataRecord.VerifierID = verifierID

	// Marshal updated data record and save to ledger
	updatedDataJSON, err := json.Marshal(dataRecord)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(dataKey, updatedDataJSON)
}

// GetDevice retrieves a device by its ID
func (s *SmartContract) GetDevice(ctx contractapi.TransactionContextInterface, deviceID string) (*Device, error) {
	deviceJSON, err := ctx.GetStub().GetState(deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to read device: %v", err)
	}
	if deviceJSON == nil {
		return nil, fmt.Errorf("device %s does not exist", deviceID)
	}

	var device Device
	err = json.Unmarshal(deviceJSON, &device)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal device JSON: %v", err)
	}
	return &device, nil
}

// GetDataRecord retrieves a data record by deviceID and timestamp
func (s *SmartContract) GetDataRecord(ctx contractapi.TransactionContextInterface, deviceID string, timestamp string) (*DataRecord, error) {
	dataKey, err := ctx.GetStub().CreateCompositeKey("DataRecord", []string{deviceID, timestamp})
	if err != nil {
		return nil, fmt.Errorf("failed to create composite key: %v", err)
	}

	dataJSON, err := ctx.GetStub().GetState(dataKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get data record: %v", err)
	}
	if dataJSON == nil {
		return nil, fmt.Errorf("data record not found")
	}

	var dataRecord DataRecord
	err = json.Unmarshal(dataJSON, &dataRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data record: %v", err)
	}
	return &dataRecord, nil
}
