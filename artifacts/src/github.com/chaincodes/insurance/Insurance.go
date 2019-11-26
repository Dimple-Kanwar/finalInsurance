package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Insurance defines the Smart Contract structure
type Insurance struct {
}

// Define the Policy structure.  Structure tags are used by encoding/json library
type Policy struct {
	policyId      string       `json:"policyId"`
	FarmerId      string       `json:"farmerId"`
	Status        string       `json:"status"`
	StartDate     string       `json:"startDate"`
	ExpiryDate    string       `json:"expiryDate"`
	FarmData      Farm_Details `json:"formData"`
	CropData      Crop_Details `json:"cropData"`
	AmountInsured int          `json:"amountInsured"`
}

//==============================================================================================================================
//	Farm_Details - Defines the structure for a Farm_Details object.
//==============================================================================================================================
type Farm_Details struct {
	Address     string `json:"address"`
	Coordinates string `json:"geo_coordinates"`
}

//==============================================================================================================================
//	Crop_Details - Defines the structure for a Crop_Details object.
//==============================================================================================================================
type Crop_Details struct {
	Crop_name   string `json:"Crop_name"`
	Crop_type   string `json:"Crop_type"`
	Crop_season string `json:"Crop_season"`
}

func (s *Insurance) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Insurance chaincode is initialized successfully.")
	return shim.Success(nil)
}

//==============================================================================================================================
//	Invoke - List all the methods to invoke/query
//==============================================================================================================================
func (s *Insurance) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "newPolicy" {
		return s.newPolicy(APIstub, args)
	}
	if function == "updateInsuranceStatus" {
		return s.updateInsuranceStatus(APIstub, args)
	}
	if function == "fetchInsuranceByStatus" {
		return s.fetchInsuranceByStatus(APIstub, args)
	}
	if function == "fetchAllInsurance" {
		return s.fetchAllInsurance(APIstub, args)
	}
	if function == "ClaimInsurance" {
		return s.ClaimInsurance(APIstub, args)
	}

	fmt.Println("function:", function, args[0])
	return shim.Error("Invalid Smart Contract function name.")
}

//==============================================================================================================================
//	newPolicy - Create a new insurance policy.
//==============================================================================================================================
func (s *Insurance) newPolicy(APIstub shim.ChaincodeStubInterface, policyData []string) sc.Response {
	fmt.Println("============= Creating a policy for the insurance =============")
	if len(args) != 10 {
		return nil, errors.New("Incorrect number of arguments. Expecting 10 arguments")
	}
	var policy Policy

	policy.FarmerId = policyData[0]
	policy.PolicyId = policyData[1]
	policy.AmountInsured = policyData[2]
	policy.Status = "New"
	policy.StartDate = policyData[3]
	policy.ExpiryDate = policyData[4]
	policy.FarmData.Address = policyData[5]
	policy.FarmData.Coordinates = policyData[6]
	policy.CropData.Name = policyData[7]
	policy.CropData.Type = policyData[8]
	policy.CropData.Season = policyData[9]

	fmt.Println("Policy: ", policy)

	// ==== Check if policy already exists ====
	policyAsBytes, err := stub.GetState(policy.PolicyId)
	if err != nil {
		return shim.Error("Failed to get policy: " + err.Error())
	} else if policyAsBytes != nil {
		fmt.Println("This policy already exists: " + policy.PolicyId)
		return shim.Error("This policy already exists: " + policy.PolicyId)
	}

	err = APIstub.PutState(policy.PolicyId, policy)
	if err != nil {
		fmt.Printf("\nError when storing insurance policy: %s", err)
		return shim.Error(err.Error())
	}

	// maintain the index
	indexName := "policyId~farmerId"
	PFIDIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{policy.PolicyId, policy.FarmerId})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the policy.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	APIstub.PutState(PFIDIndexKey, value)

	fmt.Println("============= END : Created policy successfully =============")
	return shim.Success(true)
}

//==============================================================================================================================
//	updateInsuranceStatus - Update Insurance status(ACTIVE/EXPIRED/CLAIMED).
//==============================================================================================================================
func (s *Insurance) updateInsuranceStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("============= START : Updating Insurance status =============")
	policy := fetchInsuranceByPolicyID(APIstub, args[0])
	if !policy {
		return shim.Error("updateInsuranceStatus:: Insurance policy not found.") // need to check if condition
	}
	policy.Status = args[1]
	fmt.Println("updateInsuranceStatus:: Policy: ", policy)
	err = APIstub.PutState(policy.PolicyId, policy)
	if err != nil {
		fmt.Printf("\nupdateInsuranceStatus:: Error when updating insurance policy: %s", err)
		return shim.Error("updateInsuranceStatus:: Error when updating insurance policy.Error: ", err)
	}
	fmt.Println("============= END : Updated Insurance status =============")
	return shim.Success(result)
}

//==============================================================================================================================
//	fetchInsuranceByStatus - Return insurance policy by policy status.
//==============================================================================================================================
func (s *Insurance) fetchInsuranceByStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("============= START : Fetching Insurance by status =============")
	if len(args) == 0 {
		return shim.Error("fetchInsuranceByStatus:: Incorrect number of arguments. Expecting 1 (status)")
	}
	status := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"status\":\"%s\"}}", status)
	fmt.Println("fetchInsuranceByStatus:: queryString: ", queryString)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("fetchInsuranceByStatus:: queryResults: ", queryResults)
	fmt.Println("============= END : Fetched Insurance details by status  =============")
	return shim.Success(queryResults)
}

//==============================================================================================================================
//	fetchInsuranceByPolicyID - Return insurance policy by policy id.
//==============================================================================================================================
func (s *Insurance) fetchInsuranceByPolicyID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("============= START : Fetching Insurance by policy id =============")
	if len(args) == 0 {
		return shim.Error("fetchInsuranceByPolicyID:: Incorrect number of arguments. Expecting 1 (policyID)")
	}
	policyID := args[0]
	policyAsbytes, err := APIstub.GetState(policyID) //get the policy from chaincode state
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + policyID + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp := "{\"Error\":\"Policy does not exist: " + policyID + "\"}"
		return shim.Error(jsonResp)
	}
	fmt.Println("============= END : Fetched Insurance details by policy id  =============")
	return shim.Success(policyAsbytes)
}

//==============================================================================================================================
//	fetchInsuranceByFarmerID - Return insurance policy by farmer id.
//==============================================================================================================================
func (s *Insurance) fetchInsuranceByFarmerID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("============= START : Fetching Insurance by farmer id =============")
	if len(args) == 0 {
		return shim.Error("fetchInsuranceByFarmerId:: Incorrect number of arguments. Expecting 1 (farmerID)")
	}
	farmerID := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"farmerID\":\"%s\"}}", farmerID)
	fmt.Println("fetchInsuranceByFarmerID:: queryString: ", queryString)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("fetchInsuranceByFarmerID:: queryResults: ", queryResults)
	fmt.Println("============= END : Fetched Insurance details by farmer id  =============")
	return shim.Success(queryResults)
}

//==============================================================================================================================
//	fetchAllInsurance - Return all the insurance policies.
//==============================================================================================================================
func (s *Insurance) fetchAllInsurance(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("============= START : Fatching all insurances =============")
	queryString := fmt.Sprintf("{\"selector\":{}}")
	fmt.Println("fetchAllInsurance:: queryString: ", queryString)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("fetchAllInsurance:: queryResults: ", queryResults)
	fmt.Println("============= END : Fetched all insurances =============")
	return shim.Success(queryResults)
}

//==============================================================================================================================
//	claimInsurance - Claim Insurance incase of draught or excess rainfall
//==============================================================================================================================
func (s *Insurance) claimInsurance(APIstub shim.ChaincodeStubInterface, args string) sc.Response {
	fmt.Println("============= START : Claming Insurance  =============")
	if len(args) != 1 {
		return shim.Error("claimInsurance:: Incorrect number of arguments. Expecting 1 (policyId)")
	}
	policyID := args[0]
	var farmer, insurer User
	insuranceAsBytes := fetchInsuranceByPolicyID(policyID) // fetch policy data
	// convert []bytes to insurance struct
	insuranceData, err := json.Marshal(insuranceAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	//consult oracle here
	//check if the geo-coordinates are within the areas marked as under drought by the Oracle service
	function := "fetchUserDataByUserID"
	QueryArgs := util.ToChaincodeArgs(function, userID)
	farmerData, err := stub.QueryChaincode("userCC", QueryArgs, "")
	if err != nil {
		errStr := fmt.Sprintf("Error in fetching farmer details from 'Users' chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return shim.Error(err.Error(errStr))
	}
	json.Unmarshal(farmerData, &farmer)

	function := "fetchUserDataByUserID"
	QueryArgs := util.ToChaincodeArgs(function, userID)
	insurerData, err := stub.QueryChaincode("userCC", QueryArgs, "")
	if err != nil {
		errStr := fmt.Sprintf("Error in fetching insurer details from 'Users' chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return shim.Error(err.Error(errStr))
	}
	json.Unmarshal(insurerData, &insurer)

	//check insurance policy is still active by start date and expiry date
	// Start Date = "27/04/2017"
	// Expiry date = "02/01/2006"
	// dd/mm/yyyy
	result, err := DateWithinRange(insuranceData.StartDate, insuranceData.ExpiryDate)
	if err != nil {
		fmt.Println("Error :", err)
	}
	// %t to print boolean value in fmt.Printf
	fmt.Printf("Given date with range : [%t]\n", result)

	if insuranceData.Status == "Active" {
		//claim insured amount
		farmer.AccountDetails[0].Balance += insuranceData.AmountInsured
		insurer.AccountDetails[0].Balance -= insurerData.AmountInsured
		// update its status to "CLAIMED"
		insuranceData.Status = "Claimed"
		// UPDATE THE LEDGER
	} else if insuranceData.Status == "Claimed" {
		errStr := fmt.Println("Insurance already claimed.")
		fmt.Printf(errStr)
		return shim.Error(err.Error(errStr))
	}

	fmt.Println("============= END : Claimed Insurance =============")
	return shim.Success(result)
}

// To check whether policy is expired or still within the start date and expiry date
//https://www.socketloop.com/tutorials/golang-how-to-check-if-a-date-is-within-certain-range
func DateWithinRange(dateString string, dateFormat string) (bool, error) {

	dateStamp, err := time.Parse(dateFormat, dateString)

	if err != nil {
		return false, err
	}

	today := time.Now()

	twoMonthsAgo := today.AddDate(0, -2, 0)  // minus 2 months
	twoMonthsLater := today.AddDate(0, 2, 0) // plus 2 months

	fmt.Println("Given : ", dateStamp.Format("02/01/2006"))
	fmt.Println("2 months ago : ", twoMonthsAgo.Format("02/01/2006"))
	fmt.Println("2 months later : ", twoMonthsLater.Format("02/01/2006"))

	if dateStamp.Before(twoMonthsLater) && dateStamp.After(twoMonthsAgo) {
		return true, nil
	} else {
		return false, nil
	}

	// default
	return false, nil
}

//==============================================================================================================================
// The main function is only relevant in unit test mode. Only included here for completeness.
//==============================================================================================================================
func main() {
	// Create a new Insurance chaincode
	err := shim.Start(new(Insurance))
	if err != nil {
		fmt.Printf("Error creating new Insurance chaincode: %s", err)
	}
}
