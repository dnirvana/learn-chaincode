package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"math/rand"
	"strconv"
	"time"
	"encoding/json"
	"crypto/x509"
	"net/http"
	"io/ioutil"
	"net/url"
	"encoding/pem"
	"strings"
)

//var myLogger = logging.MustGetLogger("digital_im")

type DigitalIdentityManagement struct {
}

type Share struct {
	UsedId          string `json:"userid"`
	SharedWith      string `json:"sharedwith"`
	DocName         string `json:"docname"`
	ExpiryStart     int64  `json:"expirystart"`
	ExpiryEnd       int64  `json:"expiryend"`
}

type ECertResponse struct {
	OK string `json:"OK"`
}

type User struct {
	UserId   string `json:"userid"`
	UserName string `json:"username"`
	DOB      int64 `json:"dob"`
	AddressLine string `json:"address"`
	State string `json:"state"`
	City string `json:"city"`
	ZipCode string `json:"zipcode"`
	Mobile string `json:"mobile"`
	Email string `json:"email"`
	CreationDate int64 `json:"creation"`
	ActicationDate int64 `json:"activation"`
}

type Document struct {
	OwnerId string `json:"ownerid"`
	DocName string `json:"docname"`
	DocType string `json:"doctype"`
	DocInfo string `json:"docinfo"`
	DocData []byte `json:"docdata"`
}

type DocMetaData struct {
	UserId string `json:"userId"`
	Name string `json:"docName"`
	Doctype string `json:"docType"`
	Info string   `json:"info"`
	DocVerifiDet DocVerificationDet `json:"docVerification"`
}

type DocVerificationDet struct {
	VerificationClientId string `json:"verificationClientId"`
	VerificationDt int64 `json:"verificationDt"`
	Comment   string `json:"comment"`
	ExpiryDt  int64 `json:"expiryDt"`
}

func (t *DigitalIdentityManagement) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	//myLogger.Debug("Init Chaincode...")
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	// Create ownership table
	err := createUserTable(stub)
	if err != nil {
		return nil, errors.New("Failed creating USER table.")
	}

	err2 := createDocumentTable(stub)
	if err2 != nil {
		return nil, errors.New("Failed creating DOCUMENTS table.")
	}

	err3 := createDocAuthorizationTable(stub)
	if err3 != nil {
		return nil, errors.New("Failed creating SHARE_AUTHORIZATION table.")
	}

	err4 := createDocumentVerficationTable(stub)
	if err4 != nil {
		return nil, errors.New("Failed creating DOC_VERIFICATION table.")
	}

	// Set the admin
	// The metadata will contain the certificate of the administrator
	/*adminCert, err := stub.GetCallerMetadata()
	if err != nil {
		myLogger.Debug("Failed getting metadata")
		return nil, errors.New("Failed getting metadata.")
	}
	if len(adminCert) == 0 {
		myLogger.Debug("Invalid admin certificate. Empty.")
		return nil, errors.New("Invalid admin certificate. Empty.")
	}

	myLogger.Debug("The administrator is [%x]", adminCert)

	stub.PutState("admin", adminCert)*/

	//myLogger.Debug("Init Chaincode...done")

	return nil, nil
}

func (t *DigitalIdentityManagement) createUser(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) < 9 {
		return nil, errors.New("create user failed. Must include 11 column values")
	}

	userId := args[0] //user id
	//passwd := generatePassword();
	name := args[1] // user name
	dob, err := strconv.ParseInt(args[2],10, 64)
	if err != nil {
		return nil, errors.New("Wrong Data of Birth..")
	}
	addressLine := args[3] // address line
	state := args[4] //state
	city := args[5] // city
	zipCode := args[6] //zip code
	mobile := args[7] //mobile
	email := args[8] //email
	creationDate := time.Now().Unix()

	actionvateDate := int64(0);

	var columns []*shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: name}}
	col3 := shim.Column{Value: &shim.Column_Int64{Int64: dob}}
	col4 := shim.Column{Value: &shim.Column_String_{String_: addressLine}}
	col5 := shim.Column{Value: &shim.Column_String_{String_: state}}
	col6 := shim.Column{Value: &shim.Column_String_{String_: city}}
	col7 := shim.Column{Value: &shim.Column_String_{String_: zipCode}}
	col8 := shim.Column{Value: &shim.Column_String_{String_: mobile}}
	col9 := shim.Column{Value: &shim.Column_String_{String_: email}}
	col11 := shim.Column{Value: &shim.Column_Int64{Int64: creationDate}}
	col12 := shim.Column{Value: &shim.Column_Int64{Int64: actionvateDate}}
	columns = append(columns, &col1)
	columns = append(columns, &col2)
	columns = append(columns, &col3)
	columns = append(columns, &col4)
	columns = append(columns, &col5)
	columns = append(columns, &col6)
	columns = append(columns, &col7)
	columns = append(columns, &col8)
	columns = append(columns, &col9)
	columns = append(columns, &col11)
	columns = append(columns, &col12)

	row := shim.Row{Columns: columns}
	ok, err := stub.InsertRow("USER", row)
	if err != nil {
		return nil, fmt.Errorf("User Creation operation failed. %s", err)
	}
	if !ok {
		return nil, errors.New("User Creation operation failed. Row with given key already exists")
	}
	return nil, nil
}

func (t *DigitalIdentityManagement) uploadDocument(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) < 7 {
		return nil, errors.New("uploadDocument failed. Must include 7 column in args")
	}



	userId := args[0]
	name := args[1]
	doctype := args[2]
	info := args[3]
	verificationComment := args[4]
	expiryDate,err := strconv.ParseInt(args[5],10, 64)
	doc := []byte(args[6])
	ClientId,_ := get_username(stub)
	//ClientId := args[7]

	var columns []*shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: name}}
	col3 := shim.Column{Value: &shim.Column_String_{String_: doctype}}
	col4 := shim.Column{Value: &shim.Column_String_{String_: info}}
	col5 := shim.Column{Value: &shim.Column_Bytes{Bytes: doc}}

	columns = append(columns, &col1)
	columns = append(columns, &col2)
	columns = append(columns, &col3)
	columns = append(columns, &col4)
	columns = append(columns, &col5)

	row := shim.Row{Columns: columns}
	ok, err := stub.InsertRow("DOC_STORE", row)
	if err != nil {
		return nil, fmt.Errorf("uploadDocument operation failed. %s", err)
	}
	if !ok {
		return nil, errors.New("uploadDocument operation failed. Row with given key already exists")
	}



	//Add verification Entry
	var columnVerification []*shim.Column

	col1V := shim.Column{Value: &shim.Column_String_{String_: userId}}
	col2V := shim.Column{Value: &shim.Column_String_{String_: name}}
	col3V := shim.Column{Value: &shim.Column_String_{String_: ClientId}}
	col4V := shim.Column{Value: &shim.Column_Int64{Int64: time.Now().Unix()}}
	col5V := shim.Column{Value: &shim.Column_String_{String_: verificationComment}}
	col6V := shim.Column{Value: &shim.Column_Int64{Int64: expiryDate}}

	columnVerification = append(columnVerification, &col1V)
	columnVerification = append(columnVerification, &col2V)
	columnVerification = append(columnVerification, &col3V)
	columnVerification = append(columnVerification, &col4V)
	columnVerification = append(columnVerification, &col5V)
	columnVerification = append(columnVerification, &col6V)

	verificationRow := shim.Row{Columns: columnVerification}
	okverified, error := stub.InsertRow("DOC_VERIFICATION", verificationRow)
	if error != nil {
		deleteDocStoreEntry(stub, args)
		return nil, fmt.Errorf("Document verification operation failed. %s", err)
	}
	if !okverified {
		return nil, errors.New("Document verification operation failed. Row with given key already exists")
	}


	shareDocArgs := make([]string, 10, 15)
	shareDocArgs[0] = userId  // owner id
	shareDocArgs[1] = ClientId // shared with id
	shareDocArgs[2] = name // document name
	shareDocArgs[3] = "0"
	shareDocArgs[4] = "0"

	doc,error1 := t.shareDocument(stub, shareDocArgs);
	if error1 != nil {
		deleteDocStoreEntry(stub, args)
		deleteDocVerificationEntry(stub, args)
		return nil,errors.New("Error while sharing document.")
	}
	return nil,nil
}

func (t *DigitalIdentityManagement) updateDocument(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) < 9 {
		return nil, errors.New("updateDocument failed. Must include 9 column in args")
	}


	userId := args[0]
	name := args[1]
	doctype := args[2]
	info := args[3]
	doc := []byte(args[6])

	var columns []*shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: name}}
	col3 := shim.Column{Value: &shim.Column_String_{String_: doctype}}
	col4 := shim.Column{Value: &shim.Column_String_{String_: info}}
	col5 := shim.Column{Value: &shim.Column_Bytes{Bytes: doc}}

	columns = append(columns, &col1)
	columns = append(columns, &col2)
	columns = append(columns, &col3)
	columns = append(columns, &col4)
	columns = append(columns, &col5)

	row := shim.Row{Columns: columns}
	ok, err := stub.ReplaceRow("DOC_STORE", row)
	if err != nil {
		return nil, fmt.Errorf("updateDocument operation failed. %s", err)
	}
	if !ok {
		return nil, errors.New("updateDocument operation failed. Row with given key already exists")
	}
	return nil,nil
}

func (t *DigitalIdentityManagement) fetchAllAuthorization(stub *shim.ChaincodeStub, args []string) ([]byte,error) {
	if len(args) < 1 {
		return nil, errors.New("fetchAllAuthorization failed. Must include at least key values")
	}

	uname,_ := get_username(stub)
	affiliation, _ := check_affiliation(stub, uname)

	if affiliation != 2 {
		return nil,errors.New("fetchAllAuthorization is not allowed for user " + uname)
	}
	var columns []shim.Column

	userId := args[0]
	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	columns = append(columns, col1)


	rowChannel, err := stub.GetRows("SHARE_AUTHORIZATIONS", columns)
	if err != nil {
		return nil, errors.New("fetchAllAuthorization failed. Must include at least key values")
	}

	var rows []Share
	for {
		select {
		case row, ok := <-rowChannel:
			if !ok {
				rowChannel = nil
			} else {
				var share Share
				cols := row.GetColumns()
				share.UsedId = cols[0].GetString_()
				share.SharedWith = cols[1].GetString_()
				share.DocName = cols[2].GetString_()
				share.ExpiryStart = cols[3].GetInt64()
				share.ExpiryEnd = cols[4].GetInt64()
				rows = append(rows, share)
			}
		}
		if rowChannel == nil {
			break
		}
	}
	jsonRows, err := json.Marshal(rows)
	if err != nil {
		return nil, fmt.Errorf("fetchAllAuthorization operation failed. Error marshaling JSON: %s", err)
	}
	return jsonRows,nil
}

func deleteDocStoreEntry(stub *shim.ChaincodeStub, args []string) error {
	if len(args) < 1 {
		return errors.New("Delete Doc Store Entry failed. Must include 1 key value")
	}

	userId := args[0]
	docName := args[1]
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: docName}}
	columns = append(columns, col1)
	columns = append(columns, col2)

	err := stub.DeleteRow("DOC_STORE", columns)
	if err != nil {
		return fmt.Errorf("DOC_STORE delete operation failed. %s", err)
	}
	return err
}
func deleteDocVerificationEntry(stub *shim.ChaincodeStub, args []string) error {
	if len(args) < 1 {
		return  errors.New("Delete Doc Store Entry failed. Must include 1 key value")
	}

	userId := args[0]
	docName := args[1]
	verifyingUser := args[2]
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: docName}}
	col3 := shim.Column{Value: &shim.Column_String_{String_: verifyingUser}}
	columns = append(columns, col1)
	columns = append(columns, col2)
	columns = append(columns, col3)

	err := stub.DeleteRow("DOC_VERIFICATION", columns)
	if err != nil {
		return  fmt.Errorf("DOC_STORE delete operation failed. %s", err)
	}
	return err
}




func (t *DigitalIdentityManagement) fetchDocumentList(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("fetchDocumentList failed. Must include 1 aregument value")
	}
	/*var affiliation, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, fmt.Errorf("fetchAuthorizedDocument operation failed. %s", err)
	}
	var ClientId = args[2]*/
	ClientId,_ := get_username(stub)
	affiliation, _ := check_affiliation(stub, ClientId)
	if affiliation == 1 {
		return t.getAuthorizedDocList(stub, args, ClientId)
	} else if affiliation == 2 {
		return t.getAllDocList(stub, args)
	} else {
		return nil, errors.New("Unsupported affiliation"+strconv.Itoa(affiliation))
	}
}

func (t *DigitalIdentityManagement) getAuthorizedDocList(stub *shim.ChaincodeStub, args []string, clientId string) ([]byte, error ) {

	userId := args[0]
	var columns []shim.Column


	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	columns = append(columns, col1)

	var buildQueryColForAuthorized []shim.Column
	userIdCol := shim.Column{Value: &shim.Column_String_{String_: userId}}
	buildQueryColForAuthorized = append(buildQueryColForAuthorized, userIdCol)

	//Get Document details
	rowChannelDoc, errDoc := stub.GetRows("DOC_STORE", buildQueryColForAuthorized)
	if errDoc != nil {
		return nil, fmt.Errorf("get document failed. %s", errDoc)
	}

	rowChannelDocDet, errDoc := stub.GetRows("SHARE_AUTHORIZATIONS", columns)
	if errDoc != nil {
		return nil, fmt.Errorf("get document share auth failed. %s", errDoc)
	}

	var docNames []string
	for {
		select {
		case row, ok := <-rowChannelDocDet:
			if !ok {
				rowChannelDocDet = nil
			} else {
				if !isExpired(row.Columns[3].GetInt64(), row.Columns[4].GetInt64()) {
					if row.Columns[1].GetString_() == clientId {
						var docNameCol = shim.Column{Value: &shim.Column_String_{String_: row.Columns[2].GetString_()}}
						buildQueryColForAuthorized = append(buildQueryColForAuthorized, docNameCol)
						docNames = append(docNames ,row.Columns[2].GetString_())
					}
				}
			}
		}
		if rowChannelDocDet == nil {
			break
		}
	}
	//Get Verification details
	rowChannelDocDet1, errDoc := stub.GetRows("DOC_VERIFICATION", buildQueryColForAuthorized)
	if errDoc != nil {
		return nil, fmt.Errorf("get document verification failed. %s", errDoc)
	}

	var docDet map[string]DocVerificationDet
	docDet = make(map[string]DocVerificationDet)
	for {
		select {
		case row, ok := <-rowChannelDocDet1:
			if !ok {
				rowChannelDocDet1 = nil
			} else {
				var docVerDet = DocVerificationDet{row.Columns[2].GetString_(),
					row.Columns[3].GetInt64(),
					row.Columns[4].GetString_(),
					row.Columns[5].GetInt64()}

				docDet[row.Columns[1].GetString_()] = docVerDet
			}
		}
		if rowChannelDocDet1 == nil {
			break
		}
	}



	var docMetaData []DocMetaData
	//var rows []shim.Row
	for {
		select {
		case row, ok := <-rowChannelDoc:
			if !ok {
				rowChannelDoc = nil
			} else {

				for _, docName := range docNames {

					if (docName == row.Columns[1].GetString_()) {
						var doc DocMetaData
						doc.UserId = row.Columns[0].GetString_()
						doc.Name = row.Columns[1].GetString_()
						doc.Doctype = row.Columns[2].GetString_()
						doc.Info = row.Columns[4].GetString_()
						doc.DocVerifiDet = docDet[row.Columns[1].GetString_()]
						//doc.DocVerifiDet = DocVerificationDet{}
						docMetaData = append(docMetaData, doc);
						//rows = append(rows, row)
						break;
					}
				}
			}
		}
		if rowChannelDoc == nil {
			break
		}
	}

	jsondata, jsonError := json.Marshal(docMetaData)
	if jsonError != nil {
		return nil, fmt.Errorf("get document failed. %s", jsonError)
	}
	return jsondata, nil
}

func (t *DigitalIdentityManagement) getAllDocList(stub *shim.ChaincodeStub, args []string) ([]byte, error ) {
	if len(args) < 1 {
		return nil, errors.New("getAllDocList failed. Must include at least key values")
	}

	userId := args[0]
	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	columns = append(columns, col1)

	rowChannelDocDet, errDoc := stub.GetRows("DOC_VERIFICATION", columns)
	if errDoc != nil {
		return nil, fmt.Errorf("get document verification failed. %s", errDoc)
	}

	var docDet map[string]DocVerificationDet
	docDet = make(map[string]DocVerificationDet)
	for {
		select {
		case row, ok := <-rowChannelDocDet:
			if !ok {
				rowChannelDocDet = nil
			} else {
				var docVerDet = DocVerificationDet{row.Columns[2].GetString_(),
					row.Columns[3].GetInt64(),
					row.Columns[4].GetString_(),
					row.Columns[5].GetInt64()}

				docDet[row.Columns[1].GetString_()] = docVerDet
			}
		}
		if rowChannelDocDet == nil {
			break
		}
	}

	rowChannel, err := stub.GetRows("DOC_STORE", columns)
	if err != nil {
		return nil, errors.New("fetchAllAuthorization failed. Must include at least key values")
	}

	var docMetaData []DocMetaData

	for {
		select {
		case row, ok := <-rowChannel:
			if !ok {
				rowChannel = nil
			} else {

				if docDet[row.Columns[1].GetString_()] == (DocVerificationDet{}) {
					docDet[row.Columns[1].GetString_()] = DocVerificationDet{}
				}

				var doc DocMetaData
				doc.UserId = row.Columns[0].GetString_()
				doc.Name = row.Columns[1].GetString_()
				doc.Doctype = row.Columns[2].GetString_()
				doc.Info = row.Columns[4].GetString_()

				doc.DocVerifiDet = docDet[row.Columns[1].GetString_()]
				//doc.DocVerifiDet = DocVerificationDet{}

				docMetaData = append(docMetaData, doc);
			}
		}
		if rowChannel == nil {
			break
		}
	}
	jsonRows, err := json.Marshal(docMetaData)
	if err != nil {
		return nil, fmt.Errorf("fetchAllAuthorization operation failed. Error marshaling JSON: %s", err)
	}
	return jsonRows, nil

}





func (t *DigitalIdentityManagement) getUserDetails(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("getUserDetails failed. Must include 1 key value")
	}

	userId := args[0]

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	columns = append(columns, col1)

	row, err := stub.GetRow("USER", columns)
	if err != nil {
		return nil, fmt.Errorf("getUserDetails operation failed. %s", err)
	}

	vals := row.GetColumns()
	var user User
	user.UserId = vals[0].GetString_()
	user.UserName = vals[1].GetString_()
	user.DOB = vals[2].GetInt64()
	user.AddressLine = vals[3].GetString_()
	user.State = vals[4].GetString_()
	user.City = vals[5].GetString_()
	user.ZipCode = vals[6].GetString_()
	user.Mobile = vals[7].GetString_()
	user.Email = vals[8].GetString_()
	user.CreationDate = vals[9].GetInt64()
	user.ActicationDate = vals[10].GetInt64()
	jsonRows, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("fetchAllAuthorization operation failed. Error marshaling JSON: %s", err)
	}
	return jsonRows,nil
}

func (t *DigitalIdentityManagement) updateUser(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) < 8 {
		return nil, errors.New("update user failed. Must include 8 column values")
	}

	userId := args[0]

	var cols []shim.Column
	cols1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	cols = append(cols, cols1)

	existingRow, err := stub.GetRow("USER", cols)
	if err != nil {
		return nil, fmt.Errorf("getUserDetails operation failed. %s", err)
	}
	existing := existingRow.GetColumns()


	//passwd := generatePassword();
	name := args[1] // user name
	dob := existing[2].GetInt64()
	addressLine := args[2] // address line
	state := args[3] //state
	city := args[4] // city
	zipCode := args[5] //zip code
	mobile := args[6] //mobile
	email := args[7] //email
	creationDate := time.Now().Unix()

	actionvateDate := existing[10].GetInt64()

	var columns []*shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: name}}
	col3 := shim.Column{Value: &shim.Column_Int64{Int64: dob}}
	col4 := shim.Column{Value: &shim.Column_String_{String_: addressLine}}
	col5 := shim.Column{Value: &shim.Column_String_{String_: state}}
	col6 := shim.Column{Value: &shim.Column_String_{String_: city}}
	col7 := shim.Column{Value: &shim.Column_String_{String_: zipCode}}
	col8 := shim.Column{Value: &shim.Column_String_{String_: mobile}}
	col9 := shim.Column{Value: &shim.Column_String_{String_: email}}
	col11 := shim.Column{Value: &shim.Column_Int64{Int64: creationDate}}
	col12 := shim.Column{Value: &shim.Column_Int64{Int64: actionvateDate}}
	columns = append(columns, &col1)
	columns = append(columns, &col2)
	columns = append(columns, &col3)
	columns = append(columns, &col4)
	columns = append(columns, &col5)
	columns = append(columns, &col6)
	columns = append(columns, &col7)
	columns = append(columns, &col8)
	columns = append(columns, &col9)
	columns = append(columns, &col11)
	columns = append(columns, &col12)

	row := shim.Row{Columns: columns}
	ok, err := stub.ReplaceRow("USER", row)
	if err != nil {
		return nil, fmt.Errorf("User Update operation failed. %s", err)
	}
	if !ok {
		return nil, errors.New("User Update operation failed. Row with given key already exists")
	}
	return nil, nil
}

func isValidUser(stub *shim.ChaincodeStub) (bool, error) {
	adminCert, err := stub.GetCallerMetadata()
	if err != nil {
	//myLogger.Debug("Failed getting metadata")
	return false, errors.New("Failed getting metadata.")
	}
	if len(adminCert) == 0 {
	//myLogger.Debug("Invalid admin certificate. Empty.")
	return false, errors.New("Invalid admin certificate. Empty.")
	}
	return true, nil
}

func (t *DigitalIdentityManagement)  shareDocument(stub *shim.ChaincodeStub, args []string) ([]byte, error){
	if len(args) < 5 {
		return nil, errors.New("shareDocument failed. Must include 3 key value")
	}
	userId := args[0] // owner id
	shareWith := args[1] // shared with id
	docName := args[2] // document name
	fromDate, err := strconv.ParseInt(args[3], 10, 64) // from date
	if err != nil {
		return nil, errors.New("Wrong From Date.")
	}
	toDate, err := strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		return nil,errors.New("Wrong To Date.")
	}
	var columns []*shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: shareWith}}
	col3 := shim.Column{Value: &shim.Column_String_{String_: docName}}
	col4 := shim.Column{Value: &shim.Column_Int64{Int64: fromDate}}
	col5 := shim.Column{Value: &shim.Column_Int64{Int64: toDate}}
	columns = append(columns, &col1)
	columns = append(columns, &col2)
	columns = append(columns, &col3)
	columns = append(columns, &col4)
	columns = append(columns, &col5)

	row := shim.Row{Columns: columns}
	ok, err := stub.InsertRow("SHARE_AUTHORIZATIONS", row)
	if err != nil {
		return nil, fmt.Errorf("shareDocument operation failed. %s", err)
	}
	if !ok {
		return nil, errors.New("shareDocument operation failed. Row with given key already exists")
	}
	return nil, nil

}

func (t *DigitalIdentityManagement) downloadDocument(stub *shim.ChaincodeStub, args []string) ([]byte,error){
	if len(args) < 2 {
		return nil, errors.New("downloadDocument failed. Must include 2 key value")
	}
	userId := args[0]
	docName := args[1]

	var columns []shim.Column
	col1 := shim.Column{Value:&shim.Column_String_{String_:userId}}
	col2 := shim.Column{Value:&shim.Column_String_{String_:docName}}
	columns = append(columns, col1)
	columns = append(columns, col2)

	row, err := stub.GetRow("DOC_STORE", columns)
	if err != nil {
		return nil, fmt.Errorf("downloadDocument operation failed. %s", err)
	}
	cols := row.GetColumns()
	return cols[4].GetBytes(), nil
}

func isUserAuthorized(stub *shim.ChaincodeStub, args []string) (bool, error) {
	if len(args) < 3 {
		return false, errors.New("isUserAuthorized failed. Must include 1 key value")
	}

	userId := args[0] // ownwe id
	sharedWith := args[1] // shared with Id
	docName := args[2]

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: sharedWith}}
	col3 := shim.Column{Value: &shim.Column_String_{String_: docName}}
	columns = append(columns, col1)
	columns = append(columns, col2)
	columns = append(columns, col3)

	row, err := stub.GetRow("SHARE_AUTHORIZATIONS", columns)
	if err != nil {
		return false, fmt.Errorf("isUserAuthorized operation failed. %s", err)
	}
	rowCols := row.GetColumns()
	shareStartDate := rowCols[3].GetInt64()
	shareEndDate := rowCols[4].GetInt64()
	return isExpired(shareStartDate, shareEndDate), nil
}

func  get_username(stub *shim.ChaincodeStub) (string, error) {

	bytes, err := stub.GetCallerCertificate();
	if err != nil { return "", errors.New("Couldn't retrieve caller certificate") }
	x509Cert, err := x509.ParseCertificate(bytes);				// Extract Certificate from result of GetCallerCertificate
	if err != nil { return "", errors.New("Couldn't parse certificate")	}

	return x509Cert.Subject.CommonName, nil
}

func  check_affiliation(stub *shim.ChaincodeStub, username string) (int, error) {
	cert1,_ := get_ecert(stub, username)
	cert := string(cert1)
	decodedCert, err := url.QueryUnescape(cert);    				// make % etc normal //

	if err != nil { return -1, errors.New("Could not decode certificate") }

	pem, _ := pem.Decode([]byte(decodedCert))           				// Make Plain text   //

	x509Cert, err := x509.ParseCertificate(pem.Bytes);				// Extract Certificate from argument //

	if err != nil { return -1, errors.New("Couldn't parse certificate")	}

	cn := x509Cert.Subject.CommonName

	res := strings.Split(cn,"\\")

	affiliation, _ := strconv.Atoi(res[2])

	return affiliation, nil
}

func  get_ecert(stub *shim.ChaincodeStub, name string) ([]byte, error) {

	var cert ECertResponse


	response, err := http.Get("http://"+string("localhost:5000")+"/registrar/"+name+"/ecert") 	// Calls out to the HyperLedger REST API to get the ecert of the user with that name

	if err != nil { return nil, errors.New("Error calling ecert API") }

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)					// Read the response from the http callout into the variable contents

	if err != nil { return nil, errors.New("Could not read body") }

	err = json.Unmarshal(contents, &cert)

	if err != nil { return nil, errors.New("Could not retrieve ecert for user: "+name) }

	return []byte(string(cert.OK)), nil
}

func isExpired(startTime int64, endTime int64) bool {
	currentTime := time.Now().Unix()
	if startTime == 0 && endTime == 0 {
		return false
	}
	if currentTime >= startTime && currentTime <= endTime {
		return true
	}
	return false
}

func checkClientIdTypeAuthorities(stub *shim.ChaincodeStub, clientId string) bool {
	affiliation,_ := check_affiliation(stub, clientId)
	switch affiliation {
	case 1:
		return true
	case 2:
		return false
		default:
		return false
	}
}

func (t *DigitalIdentityManagement) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	// Handle different functions
	/*ok,_ := isValidUser(stub)
	if !ok {
		return nil, errors.New("Invalid User")
	}*/
	switch function {
	case "createUser":
		//create user.
		return t.createUser(stub, args)
	case "updateUser":
		return t.updateUser(stub, args)
	case "uploadDocument":
		//Upload DOcument.
		return t.uploadDocument(stub, args)
	case "shareDocument":
		return t.shareDocument(stub, args)
	case "updateDocument":
		return t.updateDocument(stub, args)
	default:
		return nil, errors.New("Unknown operation")

	}

	return nil, errors.New("Received unknown function invocation")
}

func (t *DigitalIdentityManagement) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	// Handle different functions
	/*ok,_ := isValidUser(stub)
	if !ok {
		return nil, errors.New("Invalid User")
	}*/
	switch function {
	case "getUserDetails":
		// get user details.
		return t.getUserDetails(stub, args)
	case "fetchDocumentList":
		// get name list.
		return t.fetchDocumentList(stub, args)
	case "downloadDocument":
		return t.downloadDocument(stub, args)
	case "fetchAllAuthorization":
		return t.fetchAllAuthorization(stub, args)
	default:
		return nil, errors.New("Unknown operation")

	}

	return nil, errors.New("Received unknown function invocation")
}

func createUserTable(stub *shim.ChaincodeStub) error {
	err := stub.CreateTable("USER", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "USER_ID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "NAME", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "DOB", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "ADDRESS_LINE", Type:shim.ColumnDefinition_STRING, Key:false},
		&shim.ColumnDefinition{Name: "STATE", Type:shim.ColumnDefinition_STRING, Key:false},
		&shim.ColumnDefinition{Name: "CITY", Type:shim.ColumnDefinition_STRING, Key:false},
		&shim.ColumnDefinition{Name: "ZIP_CODE", Type:shim.ColumnDefinition_STRING, Key:false},
		&shim.ColumnDefinition{Name: "EMAIL_ID", Type:shim.ColumnDefinition_STRING, Key:false},
		&shim.ColumnDefinition{Name: "MOBILE", Type:shim.ColumnDefinition_STRING, Key:false},
		&shim.ColumnDefinition{Name: "CREATION_DATE", Type:shim.ColumnDefinition_INT64, Key:false},
		&shim.ColumnDefinition{Name: "ACTIVATION_DATE", Type:shim.ColumnDefinition_INT64, Key:false},
	})
	return err
}


func createDocumentTable(stub *shim.ChaincodeStub) error {
	err := stub.CreateTable("DOC_STORE", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "OWNER_ID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "DOCNAME", Type: shim.ColumnDefinition_STRING, Key:true},
		&shim.ColumnDefinition{Name: "DOCUMENT_TYPE", Type:shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "INFO", Type:shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "DOCUMENT_DATA", Type:shim.ColumnDefinition_BYTES, Key: false},
	})
	return err;
}

func createDocumentVerficationTable(stub *shim.ChaincodeStub) error {
	err := stub.CreateTable("DOC_VERIFICATION", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "OWNER_ID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "DOCNAME", Type: shim.ColumnDefinition_STRING, Key:true},
		&shim.ColumnDefinition{Name: "VERIFING_CLIENT_ID",Type:  shim.ColumnDefinition_STRING, Key:true},
		&shim.ColumnDefinition{Name: "VERFICATION_DATE",Type:  shim.ColumnDefinition_INT64, Key:false},
		&shim.ColumnDefinition{Name: "COMMENT", Type:shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "EXPIRY_DATE",Type:  shim.ColumnDefinition_INT64, Key:false},
	})
	return err;
}

func createDocAuthorizationTable(stub *shim.ChaincodeStub) error {
	err := stub.CreateTable("SHARE_AUTHORIZATIONS", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "OWNER_ID", Type: shim.ColumnDefinition_STRING, Key:true},
		&shim.ColumnDefinition{Name: "SHARED_WITH", Type:shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "NAME",Type:  shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "SHARE_START_DATE",Type:  shim.ColumnDefinition_INT64, Key:false},
		&shim.ColumnDefinition{Name: "SHARE_END_DATE",Type:  shim.ColumnDefinition_INT64, Key:false},
	})
	return err;
}


func generatePassword() string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 10)
	for i := 0; i < 10; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}



func main() {
	err := shim.Start(new(DigitalIdentityManagement))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
