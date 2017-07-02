package main

import (
	"errors"
	"fmt"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//var logger = shim.NewLogger("mylogger")

type SampleChaincode struct {
}

//custom data models
type IdentificationDocs struct{
	Aadhar		int `json:"aadhar"`
	Passport	int `json:"passport"`
	PAN		int `json:"pan"`
	DrivingLicense	int `json:"drivingLicense"`
	voterID		int `json:"voterId"`
}

type PersonalInfo struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Address	  string `json:""address`
	DOB       string `json:"DOB"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
	Identification IdentificationDocs `json:"identification"`
}

type FinancialInfo struct {
	MonthlySalary      int `json:"monthlySalary"`
	MonthlyRent        int `json:"monthlyRent"`
	OtherExpenditure   int `json:"otherExpenditure"`
	MonthlyLoanPayment int `json:"monthlyLoanPayment"`
}

type History struct{
	PoliceCase	string	`json:"policeCase"`
	LawCase		string	`json:"lawCase"`
	FraudCase	string	`json:"fraudCase"`
	TerrorismCase	string	`json:"terrorismCase"`
}

type Customer struct{
	ID			string	      `json:"customerId"`
	NumberOfLoans		int		
	NumberOfPendingLoan	int
	NumberOfCompletedLoans	int
	PersonalInfo		PersonalInfo  `json:"personalInfo"`
	FinancialInfo		FinancialInfo `json:"financialInfo"`
	CustomerHistory		History	      `json:"customerHistory"`
}


/*type LoanApplication struct {
	ID                     string        `json:"id"`
	PropertyId             string        `json:"propertyId"`
	LandId                 string        `json:"landId"`
	PermitId               string        `json:"permitId"`
	BuyerId                string        `json:"buyerId"`
	AppraisalApplicationId string        `json:"appraiserApplicationId"`
	SalesContractId        string        `json:"salesContractId"`
	PersonalInfo           PersonalInfo  `json:"personalInfo"`
	FinancialInfo          FinancialInfo `json:"financialInfo"`
	Status                 string        `json:"status"`
	RequestedAmount        int           `json:"requestedAmount"`
	FairMarketValue        int           `json:"fairMarketValue"`
	ApprovedAmount         int           `json:"approvedAmount"`
	ReviewerId             string        `json:"reviewerId"`
	LastModifiedDate       string        `json:"lastModifiedDate"`
}*/


func GetCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//logger.Debug("Entering GetCustomer")

	if len(args) < 1 {
		//logger.Error("Invalid number of arguments")
		return nil, errors.New("Missing Customer ID")
	}

	var customerId = args[0]
	bytes, err := stub.GetState(customerId)
	if err != nil {
		//logger.Error("Could not fetch loan application with id "+loanApplicationId+" from ledger", err)
		return nil, err
	}
	return bytes, nil
}

func CreateNewCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//logger.Debug("Entering CreateNewCustomer")

	if len(args) < 2 {
		//logger.Error("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for new customer creation")
	}

	var customerId = args[0]
	var customerData = args[1]

	err := stub.PutState(customerId, []byte(customerData))
	if err != nil {
		//logger.Error("Could not save loan application to ledger", err)
		return nil, err
	}

	var customEvent = "{eventType: 'newCustomerCreation', description:" + customerId + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	//logger.Info("Successfully saved loan application")
	return nil, nil

}

/**
Updates the status of the loan application
**/
func UpdateCustomerInformation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//logger.Debug("Entering UpdateCustomer")

	if len(args) < 2 {
		//logger.Error("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for customers' update")
	}

	var customerId = args[0]
	var updatedCustomer = args[1]
	var customer Customer
	err := json.Unmarshal([]byte(updatedCustomer), customer)
	customer.ID=customerId
	
	laBytes, err := json.Marshal(&customer)
	if err != nil {
		//logger.Error("Could not marshal customer update", err)
		return nil, err
	}

	err = stub.PutState(customerId, laBytes)
	if err != nil {
		//logger.Error("Could not save customer post update", err)
		return nil, err
	}

	var customEvent = "{eventType: 'customerInformationUpdate', description:" + customerId + "' Successfully updated information'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	//logger.Info("Successfully updated customer info")
	return nil, nil

}

func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "GetCustomer" {
		return GetCustomer(stub, args)
	}
	return nil, nil
}

func GetCertAttribute(stub shim.ChaincodeStubInterface, attributeName string) (string, error) {
	//logger.Debug("Entering GetCertAttribute")
	attr, err := stub.ReadCertAttribute(attributeName)
	if err != nil {
		return "", errors.New("Couldn't get attribute " + attributeName + ". Error: " + err.Error())
	}
	attrString := string(attr)
	return attrString, nil
}

func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "CreateNewCustomer" {
		username, _ := GetCertAttribute(stub, "username")
		role, _ := GetCertAttribute(stub, "role")
		if role == "Bank_Home_Loan_Admin" {
			return CreateNewCustomer(stub, args)
		} else {
			return nil, errors.New(username + " with role " + role + " does not have access to create a loan application")
		}

	}
	return nil, nil
}

type customEvent struct {
	Type       string `json:"type"`
	Decription string `json:"description"`
}

func main() {

	lld, _ := shim.LogLevel("DEBUG")
	fmt.Println(lld)

	//logger.SetLevel(lld)
	//fmt.Println(logger.IsEnabledFor(lld))

	err := shim.Start(new(SampleChaincode))
	if err != nil {
		//logger.Error("Could not start SampleChaincode")
	} else {
		//logger.Info("SampleChaincode successfully started")
	}

}
