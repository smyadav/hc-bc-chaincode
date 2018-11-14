package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Asset
type Asset struct {
	AssetID      string `json:"AssetID"`                // id of the asset
	Manufacturer string `json:"Manufacturer,omitempty"` // manufacturer
	OnWarranty   bool   `json:"OnWarranty,omitempty"`
	Status       string `json:"Status,omitempty"`
}

// initPrescription: create a new Asset
func (t *Chaincode) insertAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0       		1      2     	3		   4
	// "TailNumber", "AssetID", Manufacturer, "OnWarranty", "Status"
	if len(args) < 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init inserAsset")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non empty boolean")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}

	TailNumber := args[0]
	AssetID := args[1]

	Manufacturer := args[2]
	OnWarranty, err := strconv.ParseBool(args[3])

	if err != nil {
		return shim.Error("3rd arguement must be non empty boolean")
	}

	Status := args[4]

	// get flight Record
	FlightRecord := FlightRecord{}

	// retrieve flight record as bytes
	FlightRecordAsBytes, err := stub.GetState(TailNumber)
	if err != nil {
		return shim.Error("(PS)Tail Number is missing: " + err.Error())
	}

	// return error if the patient record does not exist
	if FlightRecordAsBytes == nil {
		return shim.Error("Flight Record does not exist: " + err.Error())
	}

	// convert flight record as bytes to struct
	if err := json.Unmarshal(FlightRecordAsBytes, &FlightRecord); err != nil {
		return shim.Error(err.Error())
	}

	newAsset := Asset{
		AssetID:      AssetID,
		Manufacturer: Manufacturer,
		OnWarranty:   OnWarranty,
		Status:       Status,
	}

	// see if AssetID already exists in Flight record
	for _, tempAsset := range FlightRecord.Assets {
		if tempAsset.AssetID == newAsset.AssetID {
			return shim.Error("AssetID already exists: " + tempAsset.AssetID)
		}
	}

	// add new Asset to Flight record
	FlightRecord.Assets = append(FlightRecord.Assets, newAsset)

	// convert record to JSON bytes
	FlightRecordAsBytes, err = json.Marshal(FlightRecord)
	if err != nil {
		return shim.Error("Error attempting to marshal Asset: " + err.Error())
	}
	fmt.Printf("Asset as json bytes: %s", string(FlightRecordAsBytes))

	// put record to state ledger
	err = stub.PutState(FlightRecord.TailNumber, FlightRecordAsBytes)
	if err != nil {
		return shim.Error("Error putting Asset to ledger: " + err.Error())
	}
	fmt.Printf("Entered state")

	fmt.Println("- end insertObject (success)")
	return shim.Success(nil)
}


// initPrescription: Update an Asset
func (t *Chaincode) updateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0       		1      2     	3		   4
	// "TailNumber", "AssetID", Manufacturer, "OnWarranty", "Status"
	if len(args) < 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init updateAsset")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non empty boolean")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}

	TailNumber := args[0]
	AssetID := args[1]

	Manufacturer := args[2]
	OnWarranty, err := strconv.ParseBool(args[3])

	if err != nil {
		return shim.Error("3rd arguement must be non empty boolean")
	}

	Status := args[4]

	// get flight Record
	FlightRecord := FlightRecord{}

	// retrieve flight record as bytes
	FlightRecordAsBytes, err := stub.GetState(TailNumber)
	if err != nil {
		return shim.Error("(PS)Tail Number is missing: " + err.Error())
	}

	// return error if the flight record does not exist
	if FlightRecordAsBytes == nil {
		return shim.Error("Flight Record does not exist: " + err.Error())
	}

	// convert flight record as bytes to struct
	if err := json.Unmarshal(FlightRecordAsBytes, &FlightRecord); err != nil {
		return shim.Error(err.Error())
	}

	newAsset := Asset{
		AssetID:      AssetID,
		Manufacturer: Manufacturer,
		OnWarranty:   OnWarranty,
		Status:       Status,
	}

	// see if AssetID already exists in Flight record, update the parameters
	for key, tempAsset := range FlightRecord.Assets {
		if tempAsset.AssetID == newAsset.AssetID {
			fmt.Printf("AssetID already exists: " + tempAsset.AssetID)

			//patientRecord.RxList[key].Pharmacist = pharmacist

			FlightRecord.Assets[key].Manufacturer = Manufacturer
			FlightRecord.Assets[key].OnWarranty = OnWarranty
			FlightRecord.Assets[key].Status = Status
		}
	}


	// convert record to JSON bytes
	FlightRecordAsBytes, err = json.Marshal(FlightRecord)
	if err != nil {
		return shim.Error("Error attempting to marshal Asset: " + err.Error())
	}
	fmt.Printf("Asset as json bytes: %s", string(FlightRecordAsBytes))

	// put record to state ledger
	err = stub.PutState(FlightRecord.TailNumber, FlightRecordAsBytes)
	if err != nil {
		return shim.Error("Error putting Asset to ledger: " + err.Error())
	}
	fmt.Printf("Entered state")

	fmt.Println("- end insertObject (success)")
	return shim.Success(nil)
}




// get Asset with approved attribute
func (t *Chaincode) getAssetForFlight(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//	0			1
	// "TailNumber", "AssetID"

	if len(args) < 2 {
		return shim.Error("Expecting 2 arguements")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st arguement must be a non empty string")
	}

	if len(args[1]) <= 0 {
		return shim.Error("2nd arguement must be a non empty string")
	}

	TailNumber := args[0]
	AssetID := args[1]

	// create empty patient record interface
	FlightRecord := FlightRecord{}

	// get current state of the given patient record
	FlightRecordAsBytes, err := stub.GetState(TailNumber)
	if err != nil {
		return shim.Error("Unable to get record: " + err.Error())
	}

	// convert patient record as bytes to struct
	if err := json.Unmarshal(FlightRecordAsBytes, &FlightRecord); err != nil {
		return shim.Error(err.Error())
	}

	// create custom struct for response of list of prescriptions for a given patient
	response := struct {
		TailNumber string `json:"TailNumber"`
		Asset      Asset  `json:"Asset,omitempty"`
	}{
		TailNumber: FlightRecord.TailNumber,
	}
	for _, cAsset := range FlightRecord.Assets {
		if cAsset.AssetID == AssetID {
			fmt.Println("AssetID " + AssetID)
			response.Asset = cAsset
			break
		}
	}

	// convert reponse to bytes
	responseAsBytes, err := json.Marshal(response)
	if err != nil {
		return shim.Error(err.Error())
	}

	// return results
	return shim.Success(responseAsBytes)
}
// get flight record with approved attribute
func (t *Chaincode) getFlightRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//	0
	// "TailNumber"

	if len(args) < 1 {
		return shim.Error("Expecting 1 arguement")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st arguement must be a non empty string")
	}

	TailNumber := args[0]

	// create empty flight FlightRecord interface
	FlightRecord := FlightRecord{}

	// get current state of the given Flight record
	FlightRecordAsBytes, err := stub.GetState(TailNumber)
	if err != nil {
		return shim.Error("Unable to get record: " + err.Error())
	}

	// convert patient record as bytes to struct
	if err := json.Unmarshal(FlightRecordAsBytes, &FlightRecord); err != nil {
		return shim.Error(err.Error())
	}

	// create custom struct for response of list of prescriptions for a given patient

	// convert reponse to bytes
	responseAsBytes, err := json.Marshal(FlightRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	// return results
	return shim.Success(responseAsBytes)
}
func (t *Chaincode) initFlightRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	TailNumber := args[0]

	newFlightRecord := FlightRecord{
		ObjectType: "record",
		TailNumber: TailNumber,
	}

	FlightRecordAsBytes, err := json.Marshal(newFlightRecord)
	if err != nil {
		return shim.Error("Unbale to marshal flight record:" + err.Error())
	}

	err = stub.PutState(TailNumber, FlightRecordAsBytes)
	if err != nil {
		return shim.Error("Error putting state in to ledger: " + err.Error())
	}

	return shim.Success(nil)
}
