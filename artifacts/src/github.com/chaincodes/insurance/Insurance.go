package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	oraclizeapi "github.com/oraclize/fabric-api"
)

// Insurance defines the Smart Contract structure
type Insurance struct {
}

// Define the Policy structure.  Structure tags are used by encoding/json library
type Policy struct {
	Policy_id   string       `json:"policy_id"`
	Farmer_id   string       `json:"farmer_id"`
	Status      string       `json:"status"`
	Start_date  string       `json:"start_date"`
	Expiry_date string       `json:"expiry_dates"`
	Farm_data   Farm_Details `json:"farm_data"`
	Crop_data   Crop_Details `json:"crop_data"`
	Amount_insured int `json:"amount_insured"`
}

//==============================================================================================================================
//	Farm_Details - Defines the structure for a Farm_Details object.
//==============================================================================================================================
type Farm_Details struct {
	Address         string      `json:"address"`
	Geo_coordinates Coordinates `json:"geo_coordinates"`
}

//==============================================================================================================================
//	Policy_Index - Defines array to store policy id for indexing purpose.
//==============================================================================================================================
type Policy_Index []string

//==============================================================================================================================
//	Coordinates - Defines the structure for a Coordinates object.
//==============================================================================================================================
type Coordinates struct {
	Longitude float32 `json:"longitude"`
	Latitute  float32 `json:"latitute"`
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
	return shim.Success(nil)
}

func (s *Insurance) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "newPolicy" {
		return s.newPolicy (APIstub)
	}
	if function == "updateInsurance" {
		return s.updateInsurance(APIstub)
	}
	if function == "fetchInsuranceByStatus" {
		return s.fetchInsuranceByStatus(APIstub)
	}
	if function == "fetchAllInsurance" {
		return s.fetchAllInsurance(APIstub)
	}
	if function == "ClaimInsurance" {
		return s.ClaimInsurance(APIstub)
	}

	fmt.Println("function:", function, args[0])
	return shim.Error("Invalid Smart Contract function name.")
}

//==============================================================================================================================
//	newPolicy - Create a new insurance policy.
//==============================================================================================================================
func (s *Insurance) newPolicy(APIstub shim.ChaincodeStubInterface, policy_data []string, crop_data []string, farm_data []string) sc.Response {
	fmt.Println("============= Creating a policy for the insurance =============")
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3 (Policy_data,Farm_data,Crop_data)")
	}
	var policy Policy

	policy.Farmer_id = policy_data.farmer_id
	policy.Policy_id = policy_data.policy_id
	policy.Amount_insured = policy_data.amount_insured
	policy.Status = "New"
	policy.Start_date = policy_data.start_date
	policy.Expiry_date = policy_data.Expiry_date
	policy.Farm_data.Address = farm_data.address
	policy.Farm_data.Geo_coordinates.Longitude = farm_data.longitude
	policy.Farm_data.Geo_coordinates.Latitute = farm_data.latitute
	policy.Crop_data.Crop_name = crop_data.crop_name
	policy.Crop_data.Crop_type = crop_data.crop_type
	policy.Crop_data.Crop_season = crop_data.crop_season
	err = APIstub.PutState(policy.Policy_id, policy)
	if (err != nil) {
		fmt.Printf("\nError when storing insurance policy: %s", err);
	 	return nil, err
	}

	fmt.Println(" Adding policy id for indexing");
	Policy_Index.push(policy_data.policy_id);
	err = APIstub.PutState("Policy_Index", Policy_Index)
	if (err != nil) {
		fmt.Printf("\nError when creating new insurance policy index: %s", err);
	 	return nil, err
	}
	fmt.Println(" Added policy id for indexing. Policy Index: ", Policy_Index);
	
	fmt.Println("============= END : Created policy successfully =============")
	return shim.Success(true)
}

//==============================================================================================================================
//	updateInsurance - Update Insurance status(ACTIVE/EXPIRED/CLAIMED).
//==============================================================================================================================
func (s *Insurance) updateInsurance(APIstub shim.ChaincodeStubInterface, policyId string, status string) sc.Response {
	fmt.Println("============= START : Updating Insurance  =============")
	policy := fetchInsuranceByPolicyId(stub, policyId)
	if(!policy) return shim.Error("updateInsurance:: Insurance policy not found.") // need to check if condition
	policy.Status = status
	err = APIstub.PutState(policy.Policy_id, policy)
	if (err != nil) {
		fmt.Printf("\nupdateInsurance:: Error when updating insurance policy: %s", err);
	 	return shim.Error("updateInsurance:: Error when updating insurance policy.Error: ", err)
	}
	fmt.Println("============= END : Updated Insurance =============")
	return shim.Success(result)
}

//==============================================================================================================================
//	fetchInsuranceByStatus - Return insurance policy by policy status.
//==============================================================================================================================
func (s *Insurance) fetchInsuranceByStatus(APIstub shim.ChaincodeStubInterface, status string) sc.Response {
	fmt.Println("============= START : Fetching Insurance by status =============")
	if len(status) == 0 {
		return shim.Error("fetchInsuranceByStatus:: Incorrect number of arguments. Expecting 1 (status)")
	}
	var policies []Policy
	for i := 1; i <= Policy_Index; i++ {
		policyId := Policy_Index(i)
		policy := fetchInsuranceByPolicyId(stub, policyId)
		if (policy.Status == status){
			fmt.Println("fetchInsuranceByStatus:: Found record.")
			policies = append(policies, policy)
		}
	}
	fmt.Println("============= END : Fetched Insurance details by status  =============")
	return shim.Success(policies);
}

//==============================================================================================================================
//	fetchInsuranceByPolicyId - Return insurance policy by policy id.
//==============================================================================================================================
func (s *Insurance) fetchInsuranceByPolicyId(APIstub shim.ChaincodeStubInterface, policy_id string) (sc.Response) {
	fmt.Println("============= START : Fetching Insurance by policy id =============")
	if len(policy_id) == 0 {
		return shim.Error("fetchInsuranceByPolicyId:: Incorrect number of arguments. Expecting 1 (policy_id)")
	}
	policy_data, err := stub.GetState(policy_id)
	if (err != nil) {
		fmt.Println("fetchInsuranceByPolicyId:: Policy not found.Error: ",err);
		return shim.Error(err)
	}
	var policy Policy
	err := json.Unmarshal(policy_data, &policy);
	if (err != nil) return shim.Error(err);
	fmt.Println("============= END : Fetched Insurance details by policy id  =============")
	return shim.Success(policy_data);
}

//==============================================================================================================================
//	fetchInsuranceByFarmerId - Return insurance policy by farmer id.
//==============================================================================================================================
func (s *Insurance) fetchInsuranceByFarmerId(APIstub shim.ChaincodeStubInterface, farmer_id string) sc.Response {
	fmt.Println("============= START : Fetching Insurance by farmer id =============")
	if len(farmer_id) == 0 {
		return shim.Error("fetchInsuranceByFarmerId:: Incorrect number of arguments. Expecting 1 (farmer_id)")
	}
	var policies []Policy
	for i := 1; i <= Policy_Index; i++ {
		policyId := Policy_Index(i)
		policy := fetchInsuranceByPolicyId(stub, policyId)
		if (policy.Farmer_id == farmer_id){
			fmt.Println("fetchInsuranceByFarmerId:: Found record.")
			policies = append(policies, policy)
		}
	}
	fmt.Println("============= END : Fetched Insurance details by farmer id  =============")
	return shim.Success(policies);
}

//==============================================================================================================================
//	fetchAllInsurance - Return all the insurance policies.
//==============================================================================================================================
func (s *Insurance) fetchAllInsurance(APIstub shim.ChaincodeStubInterface) ([]Policy) {
	fmt.Println("============= START : Fatching all insurances =============")
	var policies []Policy
	for i := 1; i <= Policy_Index; i++ {
		policyId := Policy_Index(i)
		policy := fetchInsuranceByPolicyId(stub, policyId)
		policies = append(policies, policy)
	}
	fmt.Println("============= END : Fetched all insurances =============")
	return policies
}

//==============================================================================================================================
//	claimInsurance - Claim Insurance incase of draught or excess rainfall 
//==============================================================================================================================
func (s *Insurance) claimInsurance(APIstub shim.ChaincodeStubInterface. status string) sc.Response {
	fmt.Println("============= START : Claming Insurance  =============")
	insurance_data := fetchInsuranceByPolicyId(policy_id); // fetch policy data
	//consult oracle here
	//check if the geo-coordinates are within the areas marked as under drought by the Oracle service
	//check insurance policy is still active by start date and expiry date
	//claim insured amount
	// amount transfer transaction begins here
	// update its status to "CLAIMED"
	// UPDATE THE LEDGER
	fmt.Println("============= END : Claimed Insurance =============")
	return shim.Success(result)
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
