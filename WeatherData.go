package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	oraclizeapi "github.com/oraclize/fabric-api"
)

// WeatherData defines the Smart Contract structure
type WeatherData struct {
}
// Setting API KEY for Agro API account 
var API_KEY = "b662848ce4c60f1334a86073c58f27de"

func (s *WeatherData) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *WeatherData) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "createPolygon" {
		return s.createPolygon(APIstub)
	}
	if function == "fetchCurrentWeather" {
		return s.fetchCurrentWeather(APIstub)
	}
	if function == "fetchWeatherData" {
		return s.fetchWeatherData(APIstub)
	}
	if function == "verifyCause" {
		return s.verifyCause(APIstub)
	}
	if function == "fetchPolygon" {
		return s.fetchPolygon(APIstub)
	}
	if function == "listAllPolygons" {
		return s.listAllPolygons(APIstub)
	}
	fmt.Println("function:", function, args[0])
	return shim.Error("Invalid Smart Contract function name.")
}
// ================================================================================================================================================================================
// createPolygon - Create a new polygon for the farm coordinates
// ================================================================================================================================================================================
func (s *WeatherData) createPolygon(APIstub shim.ChaincodeStubInterface, polygon_coordinates []float32, poly_name string ) sc.Response {
	fmt.Println("============= START : Creating a new polygon =============")
	var datasource = "URL" // Setting the Oraclize datasource
	var polygon_name = poly_name // Setting Polygon name
	if len(poly_name) == 0{
		message := "Required Polygon name to create a new polygon";
		fmt.Println(message);
		return shim.Error(message);
	}else if len(polygon_coordinates) != 5{
		message := "Incorrect number of polygon coordinates. Required 5 coordinates in sequence";
		fmt.Println(message);
		return shim.Error(message);
	}

	var geo_json = {
		"type":"Feature",
		"properties":{
		},
		"geometry":{
			"type":"Polygon",
			"coordinates":[
			   [
				  coordinates[0],
				  coordinates[1],
				  coordinates[2],
				  coordinates[3],
				  coordinates[4]	  
			   ]
			]
		}
	}
	/* [-121.1958,37.6683],
	[-121.1779,37.6687],
	[-121.1773,37.6792],
	[-121.1958,37.6792],
	[-121.1958,37.668] */
	var query = "json(http://api.agromonitoring.com/agro/1.0/polygons?appid="+ API_KEY + ")" // Setting the query
	result, proof := oraclizeapi.OraclizeQuery_sync(APIstub, datasource, query, oraclizeapi.TLSNOTARY)
	fmt.Printf("proof: %s", proof)
	fmt.Printf("\nresult: %s\n", result)
	var response = {
		"result": result,
		"proof": proof
	}
	fmt.Println("============= END : Created a new polygon =============")
	return shim.Success(result)
}

// ================================================================================================================================================================================
// fetchCurrentWeather - Fetch current weather data for the polygon id
// ================================================================================================================================================================================
func (s *WeatherData) fetchCurrentWeather(APIstub shim.ChaincodeStubInterface, polyId string ) sc.Response {
	if len(polyId) == 0 {
		message := "Required Polygon id to fetch current weather";
		fmt.Println(message);
		return shim.Error(message);
	}
	fmt.Println("============= START : Fetching current weather data for the polygon =============")
	var datasource = "URL" // Setting the Oraclize datasource
	var query = "json(http://api.agromonitoring.com/agro/1.0/weather?polyid=" + polyId + "&appid=" + API_KEY + ")" // Setting the query
	result, proof := oraclizeapi.OraclizeQuery_sync(APIstub, datasource, query, oraclizeapi.TLSNOTARY)
	fmt.Printf("proof: %s", proof)
	fmt.Printf("\nresult: %s\n", result)
	var response = {
		"result": result,
		"proof": proof
	}
	fmt.Println("============= END : Fetched current weather data for the polygon =============")
	return shim.Success(response)
}

// ================================================================================================================================================================================
// verifyCause - check that the farm coordinates are under drought or excessive rainfall area
// ================================================================================================================================================================================
func (s *WeatherData) verifyCause(APIstub shim.ChaincodeStubInterface, polyId string ) sc.Response {
	
	fmt.Println("============= START : Fetching polygon details by the polygon id =============")
	var datasource = "URL" // Setting the Oraclize datasource
	var query = "json(http://api.agromonitoring.com/agro/1.0/weather?polyid=" + polyId + "&appid=" + API_KEY + ")" // Setting the query
	result, proof := oraclizeapi.OraclizeQuery_sync(APIstub, datasource, query, oraclizeapi.TLSNOTARY)
	fmt.Printf("proof: %s", proof)
	fmt.Printf("\nresult: %s\n", result)
	var response = {
		"result": result,
		"proof": proof
	}
	fmt.Println("============= END : Fetched polygon details by the polygon id =============")
	return shim.Success(response)
	return shim.Success(response);
}

// ================================================================================================================================================================================
// fetchPolygon - Fetch polygon details by polygon id
// ================================================================================================================================================================================
func (s *WeatherData) fetchPolygon(APIstub shim.ChaincodeStubInterface, polyId string ) sc.Response {
	if len(polyId) == 0 {
		message := "Required Polygon id to fetch polygon details";
		fmt.Println(message);
		return shim.Error(message);
	}
	fmt.Println("============= START : Fetching polygon details by the polygon id =============")
	var datasource = "URL" // Setting the Oraclize datasource
	var query = "json(http://api.agromonitoring.com/agro/1.0/polygons/" + polyId + "?appid=" + API_KEY + ")" // Setting the query
	result, proof := oraclizeapi.OraclizeQuery_sync(APIstub, datasource, query, oraclizeapi.TLSNOTARY)
	fmt.Printf("proof: %s", proof)
	fmt.Printf("\nresult: %s\n", result)
	var response = {
		"result": result,
		"proof": proof
	}
	fmt.Println("============= END : Fetched polygon details by the polygon id =============")
	return shim.Success(response)
}

// ================================================================================================================================================================================
// listAllPolygons - Fetch all polygons
// ================================================================================================================================================================================
func (s *WeatherData) listAllPolygons(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("============= START : Fetching all polygons =============")
	var datasource = "URL" // Setting the Oraclize datasource
	var query = "json(http://api.agromonitoring.com/agro/1.0/polygons?appid=" + API_KEY + ")" // Setting the query
	result, proof := oraclizeapi.OraclizeQuery_sync(APIstub, datasource, query, oraclizeapi.TLSNOTARY)
	fmt.Printf("proof: %s", proof)
	fmt.Printf("\nresult: %s\n", result)
	var response = {
		"result": result,
		"proof": proof
	}
	fmt.Println("============= END : Fetched all polygons =============")
	return shim.Success(response) // array of structs
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(new(WeatherData))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
