package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
	"strconv"
	"encoding/json"
	"go/doc"
)

var myLogger = logging.MustGetLogger("digital_im")

type SyndicatedLoanManagement struct {
}

type Contract struct {
	Borrower                    string `json:"borrower"`
	BorrowerContact             string `json:"borrowercontact"`
	Purpose                     string `json:"purpose"`
	LoanAmount                  int64 `json:"loanamount"`
	InterestRate                int64 `json:"interestrate"`
	LeadBankCharge              int64 `json:"leadbankcharge"`
	ProcessingFees              int64 `json:"processingfees"`

	Tranch1Date                 int64 `json:"tranch1date"`
	Tranch1Amount               int64 `json:"tranch1amount"`
	Tranch1Comment              string `json:"tranch1comment"`

	Tranch2Date                 int64 `json:"tranch2date"`
	Tranch2Amount               int64 `json:"tranch2amount"`
	Tranch2Comment              string `json:"tranch2comment"`

	PrepaymentCharges           int64 `json:"prepaymentcharges"`
	PhysicalInspection          string `json:"physicalinspection"`

	Lender1Name                 string `json:"lender1name"`
	Lender1Share                int64 `json:"lender1share"`
	Lender1SharePercent         int64 `json:"lender1sharepercent"`
	Lender1EditAcceptanceStatus string `json:"lender1editacceptancestatus"`

	Lender2Name                 string `json:"lender2name"`
	Lender2Share                int64 `json:"lender2share"`
	Lender2SharePercent         int64 `json:"lender2sharepercent"`
	Lender2EditAcceptanceStatus string `json:"lender2editacceptancestatus"`

	Lender3Name                 string `json:"lender3name"`
	Lender3Share                int64 `json:"lender3share"`
	Lender3SharePercent         int64 `json:"lender3sharepercent"`
	Lender3EditAcceptanceStatus string `json:"lender3editacceptancestatus"`

	Lender4Name                 string `json:"lender4name"`
	Lender4Share                int64 `json:"lender4share"`
	Lender4SharePercent         int64 `json:"lender4sharepercent"`
	Lender4EditAcceptanceStatus string `json:"lender4editacceptancestatus"`

	Penalty30Days               int64 `json:"penalty30days"`
	Penalty45Days               int64 `json:"penalty45days"`
	Penalty90Days               int64 `json:"penalty90days"`

	PaySchedule1RecordDate      int64 `json:"payschedule1recorddate"`
	PaySchedule1Amount          int64 `json:"payschedule1amount"`
	PaySchedule1Status          string `json:"payschedule1status"`

	PaySchedule2RecordDate      int64 `json:"payschedule2recorddate"`
	PaySchedule2Amount          int64 `json:"payschedule2amount"`
	PaySchedule2Status          string `json:"payschedule2status"`

	PaySchedule3RecordDate      int64 `json:"payschedule3recorddate"`
	PaySchedule3Amount          int64 `json:"payschedule3amount"`
	PaySchedule3Status          string `json:"payschedule3status"`

	PaySchedule4RecordDate      int64 `json:"payschedule4recorddate"`
	PaySchedule4Amount          int64 `json:"payschedule4amount"`
	PaySchedule4Status          string `json:"payschedule4status"`

	PaySchedule5RecordDate      int64 `json:"payschedule5recorddate"`
	PaySchedule5Amount          int64 `json:"payschedule5amount"`
	PaySchedule5Status          string `json:"payschedule5status"`

	AgreementFreeFlow           string `json:"agreementfreeflow"`

	ContractStatus              string `json:"contractstatus"`

	Doc1                        string `json:"doc1"`
	Doc2                        string `json:"doc2"`
	Doc3                        string `json:"doc3"`
}

func createContractTable(stub *shim.ChaincodeStub) error {
	err := stub.CreateTable("Contract", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Borrower", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "BorrowerContact", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Purpose", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "LoanAmount", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "InterestRate", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "LeadBankCharge", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "ProcessingFees", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Tranch1Date", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Tranch1Amount", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Tranch1Comment", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Tranch2Date", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Tranch2Amount", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Tranch2Comment", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "PrepaymentCharges", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "PhysicalInspection", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Lender1Name", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Lender1Share", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Lender1SharePercent", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Lender1EditAcceptanceStatus", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Lender2Name", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Lender2Share", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Lender2SharePercent", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Lender2EditAcceptanceStatus", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Lender3Name", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Lender3Share", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Lender3SharePercent", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Lender3EditAcceptanceStatus", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Lender4Name", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Lender4Share", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Lender4SharePercent", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Lender4EditAcceptanceStatus", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Penalty30Days", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Penalty45Days", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Penalty90Days", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "PaySchedule1RecordDate", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "PaySchedule1Amount", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "PaySchedule1Status", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "PaySchedule2RecordDate", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "PaySchedule2Amount", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "PaySchedule2Status", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "PaySchedule3RecordDate", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "PaySchedule3Amount", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "PaySchedule3Status", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "PaySchedule4RecordDate", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "PaySchedule4Amount", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "PaySchedule4Status", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "PaySchedule5RecordDate", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "PaySchedule5Amount", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "PaySchedule5Status", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "AgreementFreeFlow", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ContractStatus", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Doc1", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Doc2", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Doc3", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	return err
}

type AuditTrail struct {
	EventType string `json:"eventtype"`
	EventName string `json:"eventname"`
	Date      int64  `json:"date"`
	ActionBy  string  `json:"actionby"`
	NewData   string  `json:"newdata"`
}

func createAuditTrailTable(stub *shim.ChaincodeStub) error {
	err := stub.CreateTable("AuditTrail", []*shim.ColumnDefinition{
		// EventType Contact/Re-Sale/Payment
		&shim.ColumnDefinition{Name: "EventType", Type: shim.ColumnDefinition_STRING, Key: true},
		// Contract Draft Created / Legal Council Report Published
		&shim.ColumnDefinition{Name: "EventName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Date", Type: shim.ColumnDefinition_INT64, Key: true},
		&shim.ColumnDefinition{Name: "ActionBy", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "NewData", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	return err
}

type Suggestions struct {
	LenderName string `json:"lendername"`
	WhichField string  `json:"whichfield"`
	WhatValue  string  `json:"whatvalue"`
	Action     string  `json:"action"`
	Comment    string `json:"comment"`
	Date       int64  `json:"date"`
}

func createSugesstionsTable(stub *shim.ChaincodeStub) error {
	// Create SugesstionsTable
	err := stub.CreateTable("Suggestions", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "LenderName", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "WhichField", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "WhatValue", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Action", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Comment", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Date", Type: shim.ColumnDefinition_INT64, Key: false},
	})
	return err
}

type Payments struct {
	Borrower	string `json:"borrower"`
	PaymentDate	int64  `json:"paymentDate"`
	Interest	int64  `json:"interest"`
	Principal	int64  `json:"principal"`
	Penalty		int64  `json:"penalty"`
	Fees		int64  `json:"fees"`
}

func createPaymentsTable(stub *shim.ChaincodeStub) error {
	// Create SugesstionsTable
	err := stub.CreateTable("Payments", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Borrower", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "PaymentDate", Type: shim.ColumnDefinition_INT64, Key: true},
		&shim.ColumnDefinition{Name: "Interest", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Principal", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Penalty", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "Fees", Type: shim.ColumnDefinition_INT64, Key: false},
	})
	return err
}

type Positions struct {
	Borrower        string  `json:"borrower"`
	RecordDate      int64   `json:"recorddate"`
	Lender1Name     string  `json:"lender1name"`
	Lender1Position int64        `json:"lender1position"`
	Lender2Name     string  `json:"lender2name"`
	Lender2Position int64        `json:"lender2position"`
	Lender3Name     string  `json:"lender3name"`
	Lender3Position int64        `json:"lender3position"`
	Lender4Name     string  `json:"lender4name"`
	Lender4Position int64        `json:"lender4position"`
	ContractState   string  `json:"contractstate"`
}

func createPositionsTable(stub *shim.ChaincodeStub) error {
	// Create SuggesstionsTable
	err := stub.CreateTable("Positions", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Borrower", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "RecordDate", Type: shim.ColumnDefinition_INT64, Key: true},

		&shim.ColumnDefinition{Name: "Lender1Name", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Lender1Position", Type: shim.ColumnDefinition_INT64, Key: false},

		&shim.ColumnDefinition{Name: "Lender2Name", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Lender2Position", Type: shim.ColumnDefinition_INT64, Key: false},

		&shim.ColumnDefinition{Name: "Lender3Name", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Lender3Position", Type: shim.ColumnDefinition_INT64, Key: false},

		&shim.ColumnDefinition{Name: "Lender4Name", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Lender4Position", Type: shim.ColumnDefinition_INT64, Key: false},

		&shim.ColumnDefinition{Name: "ContractState", Type: shim.ColumnDefinition_STRING, Key: false},

	})

	return err
}

func (t *SyndicatedLoanManagement) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	myLogger.Debug("Init Chaincode...")
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}
	// Create Contract table
	err := createContractTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating Contract table.")
	}
	// Create AuditTrail table
	err1 := createAuditTrailTable(stub)
	if err1 != nil {
		return nil, errors.New("Failed creating AuditTrail table.")
	}

	// Create Sugesstions table
	err2 := createSugesstionsTable(stub)
	if err2 != nil {
		return nil, errors.New("Failed creating Sugesstions table.")
	}

	// Create Positions table
	err3 := createPositionsTable(stub)
	if err3 != nil {
		return nil, errors.New("Failed creating Positions table.")
	}

	// Create Positions table
	err4 := createPaymentsTable(stub)
	if err4 != nil {
		return nil, errors.New("Failed creating Payments table.")
	}

	myLogger.Debug("Init Chaincode...done")
	return nil, nil
}

func (t *SyndicatedLoanManagement) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	switch function {
	case "createContract":
		return t.createContract(stub, args)
	case "addSuggesstion":
		return t.addSuggestion(stub, args)
	default:
		return nil, errors.New("Unknown operation")
	}
	return nil, errors.New("Received unknown function invocation")
}

func (t *SyndicatedLoanManagement) createContract(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	myLogger.Info("In create contract -------->")
	if len(args) < 54 {
		return nil, errors.New("create contract failed. Must include 54 column values")
	}
	Borrower := args[0]
	BorrowerContact := args[1]
	Purpose := args[2]
	LoanAmount, err := strconv.ParseInt(args[3],10, 64)
	InterestRate, err := strconv.ParseInt(args[4],10, 64)
	LeadBankCharge, err := strconv.ParseInt(args[5],10, 64)
	ProcessingFees, err := strconv.ParseInt(args[6],10, 64)
	Tranch1Date, err := strconv.ParseInt(args[7],10, 64)
	Tranch1Amount, err := strconv.ParseInt(args[8],10, 64)
	Tranch1Comment := args[9]
	Tranch2Date, err := strconv.ParseInt(args[10],10, 64)
	Tranch2Amount, err := strconv.ParseInt(args[12],10, 64)
	Tranch2Comment := args[12]
	PrepaymentCharges, err := strconv.ParseInt(args[13],10, 64)
	PhysicalInspection := args[14]
	Lender1Name := args[15]
	Lender1Share, err := strconv.ParseInt(args[16],10, 64)
	Lender1SharePercent, err := strconv.ParseInt(args[17],10, 64)
	Lender1EditAcceptanceStatus := args[18]
	Lender2Name := args[19]
	Lender2Share, err := strconv.ParseInt(args[20],10, 64)
	Lender2SharePercent, err := strconv.ParseInt(args[21],10, 64)
	Lender2EditAcceptanceStatus := args[22]
	Lender3Name := args[23]
	Lender3Share, err := strconv.ParseInt(args[24],10, 64)
	Lender3SharePercent, err := strconv.ParseInt(args[25],10, 64)
	Lender3EditAcceptanceStatus := args[26]
	Lender4Name := args[27]
	Lender4Share, err := strconv.ParseInt(args[28],10, 64)
	Lender4SharePercent, err := strconv.ParseInt(args[29],10, 64)
	Lender4EditAcceptanceStatus := args[30]
	Penalty30Days, err := strconv.ParseInt(args[31],10, 64)
	Penalty45Days, err := strconv.ParseInt(args[32],10, 64)
	Penalty90Days, err := strconv.ParseInt(args[33],10, 64)
	PaySchedule1RecordDate, err := strconv.ParseInt(args[34],10, 64)
	PaySchedule1Amount, err := strconv.ParseInt(args[35],10, 64)
	PaySchedule1Status := args[36]
	PaySchedule2RecordDate, err := strconv.ParseInt(args[37],10, 64)
	PaySchedule2Amount, err := strconv.ParseInt(args[38],10, 64)
	PaySchedule2Status := args[39]
	PaySchedule3RecordDate, err := strconv.ParseInt(args[40],10, 64)
	PaySchedule3Amount, err := strconv.ParseInt(args[41],10, 64)
	PaySchedule3Status := args[42]
	PaySchedule4RecordDate, err := strconv.ParseInt(args[43],10, 64)
	PaySchedule4Amount, err := strconv.ParseInt(args[44],10, 64)
	PaySchedule4Status := args[45]
	PaySchedule5RecordDate, err := strconv.ParseInt(args[46],10, 64)
	PaySchedule5Amount, err := strconv.ParseInt(args[47],10, 64)
	PaySchedule5Status := args[48]
	AgreementFreeFlow := args[49]
	ContractStatus := args[50]
	Doc1 := args[51]
	Doc2 := args[52]
	Doc3 := args[53]

	var columns []*shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: Borrower}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: BorrowerContact}}
	col3 := shim.Column{Value: &shim.Column_String_{String_: Purpose}}
	col4 := shim.Column{Value: &shim.Column_Int64{Int64: LoanAmount                  }}
	col5 := shim.Column{Value: &shim.Column_Int64{Int64: InterestRate                }}
	col6 := shim.Column{Value: &shim.Column_Int64{Int64: LeadBankCharge              }}
	col7 := shim.Column{Value: &shim.Column_Int64{Int64: ProcessingFees              }}
	col8 := shim.Column{Value: &shim.Column_Int64{Int64: Tranch1Date                 }}
	col9 := shim.Column{Value: &shim.Column_Int64{Int64: Tranch1Amount               }}
	col10 := shim.Column{Value: &shim.Column_String_{String_: Tranch1Comment              }}
	col11 := shim.Column{Value: &shim.Column_Int64{Int64: Tranch2Date                 }}
	col12 := shim.Column{Value: &shim.Column_Int64{Int64: Tranch2Amount               }}
	col13 := shim.Column{Value: &shim.Column_String_{String_: Tranch2Comment              }}
	col14 := shim.Column{Value: &shim.Column_Int64{Int64: PrepaymentCharges           }}
	col15 := shim.Column{Value: &shim.Column_String_{String_: PhysicalInspection          }}
	col16 := shim.Column{Value: &shim.Column_String_{String_: Lender1Name                 }}
	col17 := shim.Column{Value: &shim.Column_Int64{Int64: Lender1Share                }}
	col18 := shim.Column{Value: &shim.Column_Int64{Int64: Lender1SharePercent         }}
	col19 := shim.Column{Value: &shim.Column_String_{String_: Lender1EditAcceptanceStatus }}
	col20 := shim.Column{Value: &shim.Column_String_{String_: Lender2Name                 }}
	col21 := shim.Column{Value: &shim.Column_Int64{Int64: Lender2Share                }}
	col22 := shim.Column{Value: &shim.Column_Int64{Int64: Lender2SharePercent         }}
	col23 := shim.Column{Value: &shim.Column_String_{String_: Lender2EditAcceptanceStatus }}
	col24 := shim.Column{Value: &shim.Column_String_{String_: Lender3Name                 }}
	col25 := shim.Column{Value: &shim.Column_Int64{Int64: Lender3Share                }}
	col26 := shim.Column{Value: &shim.Column_Int64{Int64: Lender3SharePercent         }}
	col27 := shim.Column{Value: &shim.Column_String_{String_: Lender3EditAcceptanceStatus }}
	col28 := shim.Column{Value: &shim.Column_String_{String_: Lender4Name                 }}
	col29 := shim.Column{Value: &shim.Column_Int64{Int64: Lender4Share                }}
	col30 := shim.Column{Value: &shim.Column_Int64{Int64: Lender4SharePercent         }}
	col31 := shim.Column{Value: &shim.Column_String_{String_: Lender4EditAcceptanceStatus }}
	col32 := shim.Column{Value: &shim.Column_Int64{Int64: Penalty30Days               }}
	col33 := shim.Column{Value: &shim.Column_Int64{Int64: Penalty45Days               }}
	col34 := shim.Column{Value: &shim.Column_Int64{Int64: Penalty90Days               }}
	col35 := shim.Column{Value: &shim.Column_Int64{Int64: PaySchedule1RecordDate      }}
	col36 := shim.Column{Value: &shim.Column_Int64{Int64: PaySchedule1Amount          }}
	col37 := shim.Column{Value: &shim.Column_String_{String_: PaySchedule1Status          }}
	col38 := shim.Column{Value: &shim.Column_Int64{Int64: PaySchedule2RecordDate      }}
	col39 := shim.Column{Value: &shim.Column_Int64{Int64: PaySchedule2Amount          }}
	col40 := shim.Column{Value: &shim.Column_String_{String_: PaySchedule2Status          }}
	col41 := shim.Column{Value: &shim.Column_Int64{Int64: PaySchedule3RecordDate      }}
	col42 := shim.Column{Value: &shim.Column_Int64{Int64: PaySchedule3Amount          }}
	col43 := shim.Column{Value: &shim.Column_String_{String_: PaySchedule3Status          }}
	col44 := shim.Column{Value: &shim.Column_Int64{Int64: PaySchedule4RecordDate      }}
	col45 := shim.Column{Value: &shim.Column_Int64{Int64: PaySchedule4Amount          }}
	col46 := shim.Column{Value: &shim.Column_String_{String_: PaySchedule4Status          }}
	col47 := shim.Column{Value: &shim.Column_Int64{Int64: PaySchedule5RecordDate      }}
	col48 := shim.Column{Value: &shim.Column_Int64{Int64: PaySchedule5Amount          }}
	col49 := shim.Column{Value: &shim.Column_String_{String_: PaySchedule5Status          }}
	col50 := shim.Column{Value: &shim.Column_String_{String_: AgreementFreeFlow           }}
	col51 := shim.Column{Value: &shim.Column_String_{String_: ContractStatus              }}
	col52 := shim.Column{Value: &shim.Column_String_{String_: Doc1                        }}
	col53 := shim.Column{Value: &shim.Column_String_{String_: Doc2                        }}
	col54 := shim.Column{Value: &shim.Column_String_{String_: Doc3                        }}

	columns = append(columns, &col1)
	columns = append(columns, &col2)
	columns = append(columns, &col3)
	columns = append(columns, &col4)
	columns = append(columns, &col5)
	columns = append(columns, &col6)
	columns = append(columns, &col7)
	columns = append(columns, &col8)
	columns = append(columns, &col9)
	columns = append(columns, &col10)
	columns = append(columns, &col11)
	columns = append(columns, &col12)
	columns = append(columns, &col13)
	columns = append(columns, &col14)
	columns = append(columns, &col15)
	columns = append(columns, &col16)
	columns = append(columns, &col17)
	columns = append(columns, &col18)
	columns = append(columns, &col19)
	columns = append(columns, &col20)
	columns = append(columns, &col21)
	columns = append(columns, &col22)
	columns = append(columns, &col23)
	columns = append(columns, &col24)
	columns = append(columns, &col25)
	columns = append(columns, &col26)
	columns = append(columns, &col27)
	columns = append(columns, &col28)
	columns = append(columns, &col29)
	columns = append(columns, &col30)
	columns = append(columns, &col31)
	columns = append(columns, &col32)
	columns = append(columns, &col33)
	columns = append(columns, &col34)
	columns = append(columns, &col35)
	columns = append(columns, &col36)
	columns = append(columns, &col37)
	columns = append(columns, &col38)
	columns = append(columns, &col39)
	columns = append(columns, &col40)
	columns = append(columns, &col41)
	columns = append(columns, &col42)
	columns = append(columns, &col43)
	columns = append(columns, &col44)
	columns = append(columns, &col45)
	columns = append(columns, &col46)
	columns = append(columns, &col47)
	columns = append(columns, &col48)
	columns = append(columns, &col49)
	columns = append(columns, &col50)
	columns = append(columns, &col51)
	columns = append(columns, &col52)
	columns = append(columns, &col53)
	columns = append(columns, &col54)

	row := shim.Row{Columns: columns}
	ok, err := stub.InsertRow("Contract", row)
	if err != nil {
		return nil, fmt.Errorf("Contract Creation operation failed. %s", err)
	}
	if !ok {
		return nil, errors.New("Contract Creation operation failed. Row with given key already exists")
	}

	myLogger.Info("Contract created. -------->")
	return nil, nil
}

func (t *SyndicatedLoanManagement) addSuggestion(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	myLogger.Info("In add Suggestion -------->")
	if len(args) < 6 {
		return nil, errors.New("add Suggestion failed. Must include 6 column values")
	}
	LenderName := args[0]
	WhichField := args[1]
	WhatValue := args[2]
	Action := args[3]
	Comment := args[4]
	Date, err := strconv.ParseInt(args[5],10, 64)

	var columns []*shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: LenderName}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: WhichField}}
	col3 := shim.Column{Value: &shim.Column_String_{String_: WhatValue}}
	col4 := shim.Column{Value: &shim.Column_String_{String_: Action }}
	col5 := shim.Column{Value: &shim.Column_String_{String_: Comment  }}
	col6 := shim.Column{Value: &shim.Column_Int64{Int64: Date}}

	columns = append(columns, &col1)
	columns = append(columns, &col2)
	columns = append(columns, &col3)
	columns = append(columns, &col4)
	columns = append(columns, &col5)
	columns = append(columns, &col6)


	row := shim.Row{Columns: columns}
	ok, err := stub.InsertRow("Suggestions", row)
	if err != nil {
		return nil, fmt.Errorf("Add Suggestion operation failed. %s", err)
	}
	if !ok {
		return nil, errors.New("Add Suggestion operation failed. Row with given key already exists")
	}

	myLogger.Info("Suggestion created. -------->")
	return nil, nil
}

func (t *SyndicatedLoanManagement) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	switch function {
	case "getContractDetails":
		return t.getContractDetails(stub, args)
	case "getSuggestions":
		return t.getSuggestions(stub, args)
	default:
		return nil, errors.New("Unknown operation")
	}
	return nil, errors.New("Received unknown function invocation")
}

func (t *SyndicatedLoanManagement) getContractDetails(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	myLogger.Info("In ContractDetails --------->")
	if len(args) < 1 {
		return nil, errors.New("getContractDetails failed. Must include 1 key value")
	}

	borrower := args[0]

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: borrower}}
	columns = append(columns, col1)
	myLogger.Info("Querying  Contract table-------->")
	row, err := stub.GetRow("Contract", columns)
	if err != nil {
		return nil, fmt.Errorf("getContractDetails operation failed. %s", err)
	}

	vals := row.GetColumns()
	var contract Contract
	contract.Borrower = vals[0].GetString_()
	contract.BorrowerContact = vals[1].GetString_()
	contract.Purpose = vals[2].GetString_()
	contract.LoanAmount = vals[3].GetInt64()
	contract.InterestRate = vals[4].GetInt64()
	contract.LeadBankCharge = vals[5].GetInt64()
	contract.ProcessingFees = vals[6].GetInt64()
	contract.Tranch1Date = vals[7].GetInt64()
	contract.Tranch1Amount = vals[8].GetInt64()
	contract.Tranch1Comment = vals[9].GetString_()
	contract.Tranch2Date = vals[10].GetInt64()
	contract.Tranch2Amount = vals[11].GetInt64()
	contract.Tranch2Comment = vals[12].GetString_()
	contract.PrepaymentCharges = vals[13].GetInt64()
	contract.PhysicalInspection = vals[14].GetString_()
	contract.Lender1Name = vals[15].GetString_()
	contract.Lender1Share = vals[16].GetInt64()
	contract.Lender1SharePercent = vals[17].GetInt64()
	contract.Lender1EditAcceptanceStatus = vals[18].GetString_()
	contract.Lender2Name = vals[19].GetString_()
	contract.Lender2Share = vals[20].GetInt64()
	contract.Lender2SharePercent = vals[21].GetInt64()
	contract.Lender2EditAcceptanceStatus = vals[22].GetString_()
	contract.Lender3Name = vals[23].GetString_()
	contract.Lender3Share = vals[24].GetInt64()
	contract.Lender3SharePercent = vals[25].GetInt64()
	contract.Lender3EditAcceptanceStatus = vals[26].GetString_()
	contract.Lender4Name = vals[27].GetString_()
	contract.Lender4Share = vals[28].GetInt64()
	contract.Lender4SharePercent = vals[29].GetInt64()
	contract.Lender4EditAcceptanceStatus = vals[30].GetString_()
	contract.Penalty30Days = vals[31].GetInt64()
	contract.Penalty45Days = vals[32].GetInt64()
	contract.Penalty90Days = vals[33].GetInt64()
	contract.PaySchedule1RecordDate = vals[34].GetInt64()
	contract.PaySchedule1Amount = vals[35].GetInt64()
	contract.PaySchedule1Status = vals[36].GetString_()
	contract.PaySchedule2RecordDate = vals[37].GetInt64()
	contract.PaySchedule2Amount = vals[38].GetInt64()
	contract.PaySchedule2Status = vals[39].GetString_()
	contract.PaySchedule3RecordDate = vals[40].GetInt64()
	contract.PaySchedule3Amount = vals[41].GetInt64()
	contract.PaySchedule3Status = vals[42].GetString_()
	contract.PaySchedule4RecordDate = vals[43].GetInt64()
	contract.PaySchedule4Amount = vals[44].GetInt64()
	contract.PaySchedule4Status = vals[45].GetString_()
	contract.PaySchedule5RecordDate = vals[46].GetInt64()
	contract.PaySchedule5Amount = vals[47].GetInt64()
	contract.PaySchedule5Status = vals[48].GetString_()
	contract.AgreementFreeFlow = vals[49].GetString_()
	contract.ContractStatus = vals[50].GetString_()
	contract.Doc1 = vals[51].GetString_()
	contract.Doc2 = vals[52].GetString_()
	contract.Doc3 = vals[53].GetString_()

	jsonRows, err := json.Marshal(contract)
	if err != nil {
		return nil, fmt.Errorf("getContractDetails operation failed. Error marshaling JSON: %s", err)
	}
	return jsonRows,nil
}

func (t *SyndicatedLoanManagement) getSuggestions(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	myLogger.Info("In getSuggesstions --------->")
	if len(args) < 1 {
		return nil, errors.New("getSuggesstions failed. Must include 1 key value")
	}

	borrower := args[0]

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: borrower}}
	columns = append(columns, col1)
	myLogger.Info("Querying  Suggesstions table-------->")
	rowChannel, err := stub.GetRows("Suggestions", columns)
	if err != nil {
		return nil, fmt.Errorf("getSuggesstions operation failed. %s", err)
	}

	var rows []Suggestions

	for {
		select {
		case row, ok := <-rowChannel:
			if !ok {
				row = nil
			} else {
				var suggestion Suggestions
				cols := row.GetColumns()
				suggestion.LenderName = cols[0].GetString_()
				suggestion.WhichField = cols[1].GetString_()
				suggestion.WhatValue = cols[2].GetString_()
				suggestion.Action = cols[3].GetString_()
				suggestion.Comment = cols[4].GetString_()
				suggestion.Date = cols[5].GetInt64()

				rows = append(rows, suggestion);
			}
		}
		if rowChannel == nil {
			break
		}
	}
	jsonRows, err := json.Marshal(rows)
	if err != nil {
		return nil, fmt.Errorf("getSuggestions operation failed. Error marshaling JSON: %s", err)
	}
	return jsonRows,nil
}

func main() {
	err := shim.Start(new(SyndicatedLoanManagement))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}




