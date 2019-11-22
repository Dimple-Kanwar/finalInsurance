package main

import (
	"bytes"
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
	UserID         string    `json:"userID"`
	UserType       string    `json:"userType"`
	FullName       string    `json:"fullname"`
	Farms          []Farm    `json:"farms"`
	AccountDetails []Account `json:"accountDetails"`
	HomeAddress    string    `json:"homeAddress"`
	Phone          int       `json:"phone"`
	Email          string    `json:"email"`
}

//==============================================================================================================================
//	Define Account structure
//==============================================================================================================================
type Account struct {
	AccountNumber int     `json:"accountNumber"`
	Balance       float64 `json:"balance"`
	BankName      string  `json:"bankName"`
}

//==============================================================================================================================
//	Define Farm structure
//==============================================================================================================================
type Farm struct {
	FarmID      string `json:"farmID"`
	Address     string `json:"address"`
	Coordinates string `json:"coordinates"`
	CropDetails []Crop `json:"cropDetails"`
}

type Crop struct {
	CropName  string `json:"cropName"`
	CropType  string `json:"cropType"`
	Season    string `json:"season"`
	CropState string `json:"cropState"`
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
//{userID, userType, fullname, farms, account_details, homeAddress, phone,}
//==============================================================================================================================
func (s *Users) registerUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("============= START : Registering a user =============")
	fmt.Println("args: ", args)
	var userID, fullname, farmID, address, coordinates, cropName, cropType, season, cropState, bankName, homeAddress, email string
	var accountNumber, phone int
	var balance float64
	userType := strings.ToLower(args[1])
	var err error
	if userType == "farmer" {
		if len(args) != 15 {
			msg := "Incorrect number of arguments. Expected 15 arguments."
			fmt.Println(msg)
			return shim.Error(msg)
		}
		userID = args[0]
		fullname = args[2]
		farmID = args[3]
		address = args[4]
		coordinates = args[5]
		cropName = args[6]
		cropType = args[7]
		season = args[8]
		cropState = args[9]
		accountNumber, err = strconv.Atoi(args[10])
		if err != nil {
			return shim.Error("accountNumber must be a numeric string")
		}
		balance = 0.0
		bankName = args[11]
		homeAddress = args[12]
		phone, err = strconv.Atoi(args[13])
		if err != nil {
			return shim.Error("phone must be a numeric string")
		}
		email = args[14]
	} else if userType == "insurer" {
		if len(args) != 8 {
			msg := "Incorrect number of arguments. Expected 8 arguments."
			fmt.Println(msg)
			return shim.Error(msg)
		}
		userID = args[0]
		fullname = args[2]
		farmID = ""
		address = ""
		coordinates = ""
		cropName = ""
		cropType = ""
		season = ""
		cropState = ""
		accountNumber, err = strconv.Atoi(args[3])
		if err != nil {
			return shim.Error("accountNumber must be a numeric string")
		}
		balance = 0.0
		bankName = args[4]
		homeAddress = args[5]
		phone, err = strconv.Atoi(args[6])
		if err != nil {
			return shim.Error("phone must be a numeric string")
		}
		email = args[7]
	} else {
		msg := "Invalid User Type."
		fmt.Println(msg)
		return shim.Error(msg)
	}

	// ==== Check if user already exists ====
	userAsBytes, err := APIstub.GetState(userID)
	if err != nil {
		return shim.Error("Failed to get user: " + err.Error())
	} else if userAsBytes != nil {
		msg := "This user already exists: " + userID
		fmt.Println(msg)
		return shim.Error(msg)
	}
	// ==== Create user object and marshal to JSON ====
	//accounts := &Account{accountNumber, balance, bankName}
	//crops :=  []Crop{{cropName, cropType, season, cropState}}
	farms := []Farm{{FarmID: farmID, Address: address, Coordinates: coordinates, CropDetails: []Crop{{cropName, cropType, season, cropState}}}}
	user := User{UserID: userID, UserType: userType, FullName: fullname, Farms: farms, AccountDetails: []Account{{accountNumber, balance, bankName}}, HomeAddress: homeAddress, Phone: phone, Email: email}
	userJSONasBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(userID, userJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// maintain the index
	indexName := "userID~fullname"
	idNameIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{userID, fullname})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the user.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	APIstub.PutState(idNameIndexKey, value)

	fmt.Println("============= END : User registration done =============")
	return shim.Success(nil)
}

//==============================================================================================================================
//	fetchUserDataByUserID - Fetch user details by user id
//==============================================================================================================================
func (s *Users) fetchUserDataByUserID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("============= START : Fetching user data by user id =============")
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting user id of the user to query")
	}
	userID := args[0]
	valAsbytes, err := APIstub.GetState(userID) //get the user from chaincode state
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + userID + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp := "{\"Error\":\"user does not exist: " + userID + "\"}"
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

	userType := strings.ToLower(args[0])

	queryString := fmt.Sprintf("{\"selector\":{\"userType\":\"%s\"}}", userType)
	fmt.Println("fetchUserByType:: queryString: ", queryString)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("fetchUserByType:: queryResults: ", queryResults)
	fmt.Println("============= END : Fetched user data by user type =============")
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

// ===========================================================================================
// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

//==============================================================================================================================
//	fetchBalanceByAccount - Fetch user's account balance by account number
//==============================================================================================================================
// func (s *Users) fetchBalanceByAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
// 	fmt.Println("============= START : Fetching user's account balance by account number =============")
// 	if len(args) < 2 {
// 		return shim.Error("Incorrect number of arguments. Expecting 2")
// 	}

// 	userID := args[0]
// 	account_number := args[1]

// 	queryString := fmt.Sprintf("{\"selector\":{\"userID\":\"%s\",\"account_number\":\"%d\"}}", userID, account_number)
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

	userID := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"userID\":\"%s\"}}", userID)
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

	userID := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"userID\":\"%s\"}}", userID)
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
