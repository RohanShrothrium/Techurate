package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
var (
	fileName = "chaincode"
)
// Defining the structure of the chaincode
type TechChaincode struct{
}

// Define structure of a User

type User struct{
	username  string 
	aadhar    string 
	DOB       string 
	email     string 
	phone_num string 
	KYC_flag  string 
}

// ============ INIT ===============
func (t *TechChaincode) Init(stub shim.ChaincodeStubInterface)pb.Response{
	// Whatever variable initialisation you want can be done here //
	return shim.Success(nil)
}

// =========== INVOKE ==============
func  (t *TechChaincode) Invoke(stub shim.ChaincodeStubInterface)pb.Response{
	fmt.Println("Entering Invoke")

	// IF-ELSE-IF all the functions 
	function, args := stub.GetFunctionAndParameters()
	if function == "CreateUser" {
		return t.CreateUser(stub, args)
	}else if function == "ChangeAadhar" {
		return t.ChangeAadhar(stub, args)
	}else if function == "AproveKYC" {
		return t.ApproveKYC(stub, args)
	}else if function == "QueryFlag" {
		return t.QueryFlag(stub, args)
	}else if function == "init" {
		return t.Init(stub)
	}

	fmt.Println("invoke did not find func : " + function) //error
	return shim.Error("Received unknown function invocation")
	// end of all functions
}

/////////////////////////////////////////////////////////
////////            CREATE A USER           /////////////
////////////////////////////////////////////////////////
func  (t *TechChaincode) CreateUser(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	// username ==== aadhar ==== DOB ==== email ==== phone_num //
	if len(args) != 5 {
		fmt.Println("username ==== aadhar ==== DOB ==== email ==== phone_num")
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	// Assigning values to variables
	var username = args[0]
	var aadhar = args[1]
	var DOB = args[2]
	var email = args[3]
	var phone_num = args[4]
	var KYC_flag = "False"

	userAsByte, err := stub.GetState(username)

	// ======== Check user with username already exists =========
	if err != nil {
		return shim.Error("Error encountered: " + err.Error())
	}else if userAsByte != nil {
		fmt.Println(username + " already exists")
		return shim.Error(username + " already exists")
	}

	// ============= create User with a user name- username ==============//
	var User = &User{username: username, aadhar: aadhar, DOB: DOB, email: email, phone_num: phone_num, KYC_flag: KYC_flag}
	UserJSONasBytes, err := json.Marshal(User)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ============ User is being put to blockchain ============//
	err = stub.PutState(username, UserJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("User Successfully created")
	return shim.Success(nil)
}

///////////////////////////////////////////////////////
//////             CHANGE AADHAR          /////////////
//////////////////////////////////////////////////////
func (t* TechChaincode) ChangeAadhar(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	// username ==== aadhar
	if len(args) != 2 {
		fmt.Println("Give arguments as: username , aadhar number")
		return shim.Error("Incorrect number of arguments")
	}

	// ======== Assigning values ==========
	var username = args[0]
	var aadhar = args[1]

	// Check if User exists
	UserAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error("Failed to get User:" + err.Error())
	}else if UserAsBytes == nil {
		return shim.Error("User does not exist")
	}
	var UserChanged = User{}
	err = json.Unmarshal(UserAsBytes, &UserChanged)
	if err != nil {
		return shim.Error(err.Error())
	}
	UserChanged.aadhar = aadhar
	UserAsBytes, err = json.Marshal(UserChanged)
	err = stub.PutState(username, UserAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("Aadhar number Succesfully Updated")
	return shim.Success(nil)
}

/////////////////////////////////////////////////////////
////////             APPROVE KYC           /////////////
////////////////////////////////////////////////////////
func (t* TechChaincode) ApproveKYC(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	// username =====
	if len(args) != 1 {
		fmt.Println("Give arguments as: username ")
		return shim.Error("Incorrect number of arguments")
	}

	// ======== Assigning values ==========
	var username = args[0]

	// Check if User exists
	UserAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error("Failed to get User:" + err.Error())
	}else if UserAsBytes == nil {
		return shim.Error("User does not exist")
	}
	var UserChanged = User{}
	err = json.Unmarshal(UserAsBytes, &UserChanged)
	if err != nil {
		return shim.Error(err.Error())
	}
	UserChanged.KYC_flag = "True"
	UserAsBytes, err = json.Marshal(UserChanged)
	err = stub.PutState(username, UserAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("KYC Succesfully Done")
	return shim.Success(nil)
}

/////////////////////////////////////////////////////////
////////            QUERY KYC_FLAG          /////////////
////////////////////////////////////////////////////////
func (t* TechChaincode) QueryFlag(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	// username====
	if len(args) != 1{
		fmt.Println("Give User Name as first argument")
		return shim.Error("Incorrect number of arguments")
	}
	var username = args[0]
	UserAsBytes, err := stub.GetState(username)
	if err != nil {
		fmt.Println("Invalid username")
		return shim.Error(err.Error())
	}
	return shim.Success(UserAsBytes)
}

///////////////////////////////////////////////////////
//////     THIS IS THE MAIN FUNCTION     /////////////
//////////////////////////////////////////////////////
func  main() {
	err := shim.Start(new(TechChaincode))
	if err != nil {
		fmt.Printf("Error starting Chaincode: %s", err)
	}
}