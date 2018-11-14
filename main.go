package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode example simple Chaincode implementation
type Chaincode struct {
}

// Flight record containing Asset details, etc
// summary:
//
type FlightRecord struct {
	ObjectType string  `json:"objType"`
	Assets     []Asset `json:"Assets,omitempty"`     // list of prescriptions that the patient has currently
	TailNumber string  `json:"tailNumber,omitempty"` // current insurance

}

// Main
func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting File Trace chaincode: %s", err)
	}
}

// Init initializes chaincode
func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	response := t.initFlightRecord(stub, []string{"TailNumber1"})
	fmt.Println(response.GetMessage())
	response = t.initFlightRecord(stub, []string{"TailNumber2"})
	fmt.Println(response.GetMessage())
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "insertAsset" {
		// TESTED OK
		return t.insertAsset(stub, args)
	} else if function == "getAssets" {
		// TESTED OK
		return t.getAssetForFlight(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// createIndex - create search index for ledger
// currently used to create a composite key
// used in getPeople by adding each person to a people index that we can query against
func (t *Chaincode) createIndex(stub shim.ChaincodeStubInterface, indexName string, attributes []string) error {
	fmt.Println("- start create index")
	var err error
	//  ==== Index the object to enable range queries, e.g. return all parts made by supplier b ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	indexKey, err := stub.CreateCompositeKey(indexName, attributes)
	if err != nil {
		return err
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of object.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(indexKey, value)

	fmt.Println("- end create index")
	return nil
}
