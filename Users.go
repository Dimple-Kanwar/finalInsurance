package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

//==============================================================================================================================
//	Users defines users chaincode structure
//==============================================================================================================================
type Users struct {
}

//==============================================================================================================================
//	Define the User structure.  Structure tags are used by encoding/json library
//==============================================================================================================================
type User struct {
	User_Id         string    `json:"user_id"`
	User_Type       string    `json:"user_type"`
	FullName        string    `json:"fullname"`
	Farms           []Farm    `json:"farms"`
	Account_Details []Account `json:"account_details"`
	HomeAddress     string    `json:"home_address"`
	Phone           int       `json:"phone"`
	Email           string    `json:"email"`
}

//==============================================================================================================================
//	Define Account structure
//==============================================================================================================================
type Account struct {
	Account_Number int     `json:"account_number"`
	Balance        float32 `json:"balance"`
	Bank_Name      string  `json:"bank_name"`
}

//==============================================================================================================================
//	Define Farm structure
//==============================================================================================================================
type Farm struct {
	Farm_Id      string    `json:"farm_id"`
	Address      string    `json:"address"`
	Coordinates  []float32 `json:"coordinates"`
	Crop_Details []Crop    `json:"crop_details"`
}

type Crop struct {
	Crop_name  string    `json:"crop_name"`
	Crop_type  string    `json:"crop_type"`
	Season     []float32 `json:"season"`
	Crop_state []Crop    `json:"crop_state"`
}

//==============================================================================================================================
//	Init - Init chaincode
//==============================================================================================================================
func (s *Users) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

//==============================================================================================================================
//	Invoke - List all the methods to invoke/query
//==============================================================================================================================
func (s *Users) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested chaincode function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	switch function {
	case "registerUser":
		return s.registerUser(APIstub, args)
	case "fetchUserDataByUserID":
		return s.fetchUserDataByUserID(APIstub, args)
	// case "fetchAllUsers":
	// 	return s.fetchAllUsers(APIstub, args)
	case "fetchUserByType":
		return s.fetchUserByType(APIstub, args)
	// case "fetchBalanceByAccount":
	// 	return s.fetchBalanceByAccount(APIstub, args)
	case "fetchFarmsByUserId":
		return s.fetchFarmsByUserId(APIstub, args)
	case "fetchAccountsByUserId":
		return s.fetchAccountsByUserId(APIstub, args)
	default:
		fmt.Println("invoke did not find func: " + function)
		return shim.Error("Received unknown function invocation")
	}
}

//==============================================================================================================================
//	registerUser - Register a new user
// Input -
//{user_id, user_type, fullname, farms, account_details, home_address, phone,}
//==============================================================================================================================
func (s *Users) registerUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("============= START : Registering a user =============")
	fmt.Println("args: ", args)
	if args.length != 8 {
		msg := "Incorrect number of arguments. Expected 8 arguments."
		fmt.Println(msg)
		return shim.Error(msg)
	}
	user_id := args[0]
	user_type := strings.ToLower(args[1])
	fullname := args[2]
	farms := args[3]
	account_details := args[4]
	home_address := args[5]
	phone, err := strconv.Atoi(args[6])
	if err != nil {
		return shim.Error("phone must be a numeric string")
	}
	email := args[7]

	// ==== Check if user already exists ====
	userAsBytes, err := APIstub.GetState(user_id)
	if err != nil {
		return shim.Error("Failed to get user: " + err.Error())
	} else if userAsBytes != nil {
		msg := "This user already exists: " + user_id
		fmt.Println(msg)
		return shim.Error(msg)
	}
	// ==== Create user object and marshal to JSON ====
	objectType := "user"
	user := &User{user_id, user_type, fullname, farms, account_details, home_address, phone, email}
	userJSONasBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(user_id, userJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// maintain the index
	indexName := "user_id~fullname"
	idNameIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{User.User_Id, User.FullName})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the user.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	APIstub.PutState(idNameIndexKey, value)

	fmt.Println("============= END : User registration done =============")
	return shim.Success(result)
}

//==============================================================================================================================
//	fetchUserDataByUserID - Fetch user details by user id
//==============================================================================================================================
func (s *Users) fetchUserDataByUserID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("============= START : Fetching user data by user id =============")
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting user id of the user to query")
	}
	user_id = args[0]
	valAsbytes, err := APIstub.GetState(user_id) //get the user from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + user_id + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"user does not exist: " + user_id + "\"}"
		return shim.Error(jsonResp)
	}
	fmt.Println("============= END : Fetched user data by user id =============")
	return shim.Success(valAsbytes)
}

//==============================================================================================================================
//	fetchAllUsers - Fetch all users
//==============================================================================================================================
// func (s *Users) fetchAllUsers(APIstub shim.ChaincodeStubInterface) sc.Response {
// 	fmt.Println("============= START : Calling the oraclize chaincode =============")
// 	var datasource = "URL"                                                                  // Setting the Oraclize datasource
// 	var query = "json(https://min-api.cryptocompare.com/data/price?fsym=EUR&tsyms=USD).USD" // Setting the query
// 	result, proof := oraclizeapi.OraclizeQuery_sync(APIstub, datasource, query, oraclizeapi.TLSNOTARY)
// 	fmt.Printf("proof: %s", proof)
// 	fmt.Printf("\nresult: %s\n", result)
// 	fmt.Println("Do something with the result...")
// 	fmt.Println("============= END : Calling the oraclize chaincode =============")
// 	return shim.Success(result)
// }

//==============================================================================================================================
//	fetchUserByType - Fetch users by user type
//==============================================================================================================================
func (s *Users) fetchUserByType(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("============= START : Fetching user data by user type =============")
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	user_type := strings.ToLower(args[0])

	queryString := fmt.Sprintf("{\"selector\":{\"user_type\":\"%s\"}}", user_type)
	fmt.Println("fetchUserByType:: queryString: ", queryString)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("fetchUserByType:: queryResults: ", queryResults)
	fmt.Println("============= END : Fetched user data by user type =============")
	return shim.Success(queryResults)
}

//==============================================================================================================================
//	fetchBalanceByAccount - Fetch user's account balance by account number
//==============================================================================================================================
// func (s *Users) fetchBalanceByAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
// 	fmt.Println("============= START : Fetching user's account balance by account number =============")
// 	if len(args) < 2 {
// 		return shim.Error("Incorrect number of arguments. Expecting 2")
// 	}

// 	user_id := args[0]
// 	account_number := args[1]

// 	queryString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\",\"account_number\":\"%d\"}}", user_id, account_number)
// 	fmt.Println("fetchBalanceByAccount:: queryString: ", queryString)
// 	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	fmt.Println("fetchBalanceByAccount:: queryResults: ", queryResults)
// 	fmt.Println("============= END : Fetched user's account balance by account number =============")
// 	return shim.Success(queryResults)
// }

//==============================================================================================================================
//	fetchFarmsByUserId - Fetch all the farms of a farmer
//==============================================================================================================================
func (s *Users) fetchFarmsByUserId(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("============= START : Fetching farmer's farms by user id =============")
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	user_id := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\"}}", user_id)
	fmt.Println("fetchFarmsByUserId:: queryString: ", queryString)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("fetchFarmsByUserId:: queryResults: ", queryResults)
	fmt.Println("============= END : Fetched farmer's farms by user id =============")
	return shim.Success(queryResults)
}

//==============================================================================================================================
//	fetchAccountsByUserId - Fetch user's account details by user id
//==============================================================================================================================
func (s *Users) fetchAccountsByUserId(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("============= START : Fetching user's accounts by user id =============")
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	user_id := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\"}}", user_id)
	fmt.Println("fetchAccountsByUserId:: queryString: ", queryString)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("fetchAccountsByUserId:: queryResults: ", queryResults)
	fmt.Println("============= END : Fetched user's accounts by user id =============")
	return shim.Success(queryResults)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(new(Users))
	if err != nil {
		fmt.Printf("Error creating new User chaincode: %s", err)
	}
}
