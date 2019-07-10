// main_test.go


package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	_ "strconv"
	"testing"
  "reflect"
  "time"
)
const (
  prUpdated     string = "updated"
  prCreated     string = "created"
  prActive      string = "active"
  prUserIdMapperURL    string = "/api/v1/namespace/pavedroad.io/prUserIdMappers/%s"
)

var a prUserIdMapperApp

func TestMain(m *testing.M) {
	a = prUserIdMapperApp{}
	a.Initialize()

	clearDB()
	ensureTableExists()

	code := m.Run()

  fmt.Println("cleartable")
	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
  //fmt.Println(tableCreationQuery)
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		fmt.Println("Table check failed:", err)
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM pavedroad.prUserIdMapper")
}

func clearDB() {
	a.DB.Exec("DROP DATABASE IF EXISTS pavedroad")
	a.DB.Exec("CREATE DATABASE pavedroad")
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS pavedroad.prUserIdMapper (
    apiVersion  VARCHAR(254) NOT NULL,
    kind VARCHAR(25) NOT NULL,
    objVersion CHAR(10) NOT NULL,
    credential VARCHAR(254) NOT NULL PRIMARY KEY,
    userUUID VARCHAR(254),
    loginCount INT,
    created TIMESTAMPTZ NOT NULL,
    updated TIMESTAMPTZ,
    active CHAR(5) NOT NULL
);`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/v1/namespace/pavedroad.io/prUserIdMappersLIST", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentToken(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET",
		"/api/v1/namespace/pavedroad.io/prUserIdMappers/43ae99c9", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "credential 43ae99c9 does not exist" {
		t.Errorf("Expected the 'error' key of the response to be set to 'credential 43ae99c9 does not exist'. Got '%s'", m["error"])
	}
}

func TestCreatePrUserIdMapper(t *testing.T) {
	clearTable()

	payload := []byte(`{
  "apiVersion" :"core.pavedroad.io/v1alpha1",
  "objVersion": "v1beta1",
  "kind" : "prUserIdMapper",
  "login": "john@scharber.com",
  "userUUID": "",
  "loginCount": 10,
  "created": "",
  "updated": "",
  "active": "true"}`)


	req, _ := http.NewRequest("POST", "/api/v1/namespace/pavedroad.io/prUserIdMappers", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	//var md map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["apiVersion"] != "core.pavedroad.io/v1alpha1" {
		t.Errorf("Expected apiVersion to be 'core.pavedroad.io/v1alpha1'. Got '%v'", m["apiVersion"])
	}

	if m["kind"] != "prUserIdMapper" {
		t.Errorf("Expected kind to be 'prUserIdMapper'. Got '%v'", m["kind"])
	}

	if m["objVersion"] != "v1beta1" {
		t.Errorf("Expected objVersion to be 'v1beta1'. Got '%v'", m["objVersion"])
	}

	if m["login"] != "john@scharber.com" {
		t.Errorf("Expected login to be 'john@scharber.com', Got '%v'", m["login"])
	}

	if m["userUUID"] != "" {
		t.Errorf("Expected userUUID to be '', Got '%v'", m["userUUID"])
	}

  // numbers are converted to float64
  var ec float64 = 1
	if m["loginCount"] != ec {
		t.Errorf("Expected loginCount to be 1, Got %v", m["loginCount"])
	}

  //Test we can decode the data
  cs, ok :=  m["created"].(string)
  if ok {
    c, err := time.Parse(time.RFC3339, cs)
  	if err != nil {
  		t.Errorf("Parse failed on parse updated time Got '%v'", c)
  	}
  } else {
  		t.Errorf("Expected updated of string type Got '%v'", reflect.TypeOf(m["updated"]))
  }

  us, ok :=  m["updated"].(string)
  if ok {
    u, err := time.Parse(time.RFC3339, us)
  	if err != nil {
  		t.Errorf("Parse failed on parse updated time Got '%v'", u)
  	}
  } else {
  		t.Errorf("Expected updated of string type Got '%v'", reflect.TypeOf(m["updated"]))
  }

	if m["active"] != "true" {
		t.Errorf("Expected active to be 'true', Got '%v'", m["active"])
	}

}

func TestMarshallprUsreIdMapper(t *testing.T) {
  nt := NewPrUserIdMapper()
  _, err := json.Marshal(nt)
  if err != nil {
    t.Errorf("Marshal of prUserIdMapper failed: Got '%v'", err)
  }
}

func addPrUserIdMapper(t *prUserIdMapper) (cred string) {

  statement := fmt.Sprintf("INSERT INTO pavedroad.pruseridmapper(apiVersion, objVersion, kind, credential, userUUID, loginCount, created, updated, active) VALUES('%s', '%s', '%s', '%s', '%s', '%d', '%s', '%s', '%s');",
    t.APIVersion, t.ObjVersion, t.Kind, t.Credential, t.UserUUID, t.LoginCount,
    t.Created, t.Updated, t.Active)
  //fmt.Println(statement)

  rows, er1 := a.DB.Query(statement)

  if er1 != nil {
    log.Printf("Insert failed for: %s", t.Credential)
    log.Printf("SQL Error: %s", er1)
    return ""
  }

  defer rows.Close()

  return t.Credential
}

//create a token, set default values
func NewPrUserIdMapper () (t* prUserIdMapper) {
  n := prUserIdMapper{
  APIVersion: "core.pavedroad.io/v1alpha1",
  Kind: "prUserIdMapper",
  ObjVersion: "v1beta1",
  Credential: "foobar@youspace.com",
  UserUUID: "",
  LoginCount: 0,
  Created: "2002-10-02T15:00:00.05Z",
  Updated: "2002-10-02T15:00:00.05Z",
  Active: "true"}

  return &n
}

//test getting a mapper
func TestGetPrUserIdMapper(t *testing.T) {
	clearTable()
  nt := NewPrUserIdMapper()
  cred := addPrUserIdMapper(nt)

  statement := fmt.Sprintf("/api/v1/namespace/pavedroad.io/prUserIdMappers/%s", cred)

	req, _ := http.NewRequest("GET", statement, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateUser(t *testing.T) {
	clearTable()
  nt := NewPrUserIdMapper()
  cred := addPrUserIdMapper(nt)

  statement := fmt.Sprintf(prUserIdMapperURL, cred)
	req, _ := http.NewRequest("GET", statement, nil)
	response := executeRequest(req)

	//var originalPrUserIdMapper map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &nt)

  ut := nt

  //Update the new struct
	ut.Active = "eslaf"
	ut.UserUUID = "12345"

  jb, err := json.Marshal(ut)
  if err != nil {
    panic(err)
  }

	req, _ = http.NewRequest("PUT", statement, strings.NewReader(string(jb)))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["active"] != "eslaf" {
		t.Errorf("Expected active to be eslaf. Got %v", m["active"])
	}

	if m["userUUID"] != "12345" {
		t.Errorf("Expected UserUUID to be faluse. Got %v", m["userUUID"])
	}

}

func TestDeletePrUserIdMapper(t *testing.T) {
	clearTable()
  nt := NewPrUserIdMapper()
  cred := addPrUserIdMapper(nt)

  statement := fmt.Sprintf(prUserIdMapperURL, cred)
	req, _ := http.NewRequest("DELETE", statement, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", statement, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
