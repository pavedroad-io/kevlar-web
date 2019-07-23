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
	"reflect"
	_ "strconv"
	"strings"
	"testing"
	"time"
)

const (
	prUpdated         string = "updated"
	prCreated         string = "created"
	prActive          string = "active"
	prUserURL string = "/api/v1/namespace/pavedroad.io/prUsers/%s"
)

var newUserJSON=`
{
  "apiVersion" :"core.pavedroad.io/v1alpha1",
  "objVersion": "v1beta1",
  "kind" : "prUser",
  "metadata": {
    "name" : "test",
    "namespace" : "pavedroad",
    "uuid" : "", 
    "profiles": {"key": "value"}
  },
  "labels": {"key": "value"},
  "created": "",
  "updated": "",
  "active": "true"
}`

var a prUserApp

func TestMain(m *testing.M) {
	a = prUserApp{}
	a.Initialize()

	clearDB()
	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	//fmt.Println(tableCreationQuery)
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		fmt.Println("Table check failed:", err)
		log.Fatal(err)
	}

	if _, err := a.DB.Exec(indexCreate); err != nil {
		fmt.Println("Table check failed:", err)
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM pavedroad.prUser")
}

func clearDB() {
	a.DB.Exec("DROP DATABASE IF EXISTS pavedroad")
	a.DB.Exec("CREATE DATABASE pavedroad")
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS pavedroad.prUser (
    uuid UUID DEFAULT uuid_v4()::UUID PRIMARY KEY,
    prUser JSONB
);`

const indexCreate = `
CREATE INDEX IF NOT EXISTS prUserIdx ON pavedroad.prUser USING GIN (prUser);`


func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/v1/namespace/pavedroad.io/prUsersLIST", nil)
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

func TestGetWithBadUserUUID(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET",
		"/api/v1/namespace/pavedroad.io/prUsers/43ae99c9", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "invalid UUID length: 8" {
		t.Errorf("Expected the 'error' key of the response to be set to 'invalid UUID length: 8'. Got '%s'", m["error"])
	}
}

func TestGetWrongUUID(t *testing.T) {
	clearTable()
  nt := NewPrUser()
  addPrUser(nt)
  badUid := "00000000-d01d-4c09-a4e7-59026d143b89"

	statement := fmt.Sprintf(prUserURL, badUid)

	req, _ := http.NewRequest("GET", statement, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestCreatePrUser(t *testing.T) {
	clearTable()

	payload := []byte(`{
  "apiVersion" :"core.pavedroad.io/v1alpha1",
  "objVersion": "v1beta1",
  "kind" : "prUser",
  "metadata": {
    "name" : "test",
    "namespace" : "pavedroad",
    "uuid" : "",
    "profiles" : {"key": "value"}
  },
  "labels": {"key": "value"},
  "created": "",
  "updated": "",
  "active": "true"
}`)

	req, _ := http.NewRequest("POST", "/api/v1/namespace/pavedroad.io/prUsers", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	//var md map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["apiVersion"] != "core.pavedroad.io/v1alpha1" {
		t.Errorf("Expected apiVersion to be 'core.pavedroad.io/v1alpha1'. Got '%v'", m["apiVersion"])
	}

	if m["kind"] != "prUser" {
		t.Errorf("Expected kind to be 'prUser'. Got '%v'", m["kind"])
	}

	if m["objVersion"] != "v1beta1" {
		t.Errorf("Expected objVersion to be 'v1beta1'. Got '%v'", m["objVersion"])
	}

  tname := m["Metadata"].(map[string]interface{})["name"]
	if tname != "test" {
		t.Errorf("Expected login to be 'test', Got '%v'", tname)
	}

  tuid := m["Metadata"].(map[string]interface{})["uuid"]
	if tuid == "" {
		t.Errorf("Expected uuid to be uuid Got '%v'", tuid)
	}

	//Test we can decode the data
	cs, ok := m["created"].(string)
	if ok {
		c, err := time.Parse(time.RFC3339, cs)
		if err != nil {
			t.Errorf("Parse failed on parse updated time Got '%v'", c)
		}
	} else {
		t.Errorf("Expected updated of string type Got '%v'", reflect.TypeOf(m["updated"]))
	}

	us, ok := m["updated"].(string)
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

func TestMarshallprUsre(t *testing.T) {
	nt := NewPrUser()
	_, err := json.Marshal(nt)
	if err != nil {
		t.Errorf("Marshal of prUser failed: Got '%v'", err)
	}
}

func addPrUser(t *prUser) (string) {

  statement := fmt.Sprintf("INSERT INTO pavedroad.pruser(pruser) VALUES('%s') RETURNING uuid", newUserJSON)
  rows, er1 := a.DB.Query(statement)

	if er1 != nil {
		log.Printf("Insert failed error %s", er1)
		return ""
	}

	defer rows.Close()

  for rows.Next() {
    err := rows.Scan(&t.Metadata.UUID)
    if err != nil {
      return ""
    }
  }

  return t.Metadata.UUID
}

//create a token, set default values
func NewPrUser() (t *prUser) {
	n := prUser{
		APIVersion: "core.pavedroad.io/v1alpha1",
		Kind:       "prUser",
		ObjVersion: "v1beta1",
    Metadata: userMetadata{
      Name: "test",
      Namespace: "pavedroad",
      UUID: "",
    },
		Created:    "2002-10-02T15:00:00.05Z",
		Updated:    "2002-10-02T15:00:00.05Z",
		Active:     "true"}

    profiles := make(map[string](string))
    labels := make(map[string](string))
    profiles["name"] = "fool"
    labels["status"] = "terminated"

    n.Metadata.Profiles = profiles
    n.Labels = labels

	return &n
}

//test getting a mapper
func TestGetPrUser(t *testing.T) {
	clearTable()
	nt := NewPrUser()
	uid := addPrUser(nt)
	statement := fmt.Sprintf(prUserURL, uid)
  fmt.Println(statement)

	req, err := http.NewRequest("GET", statement, nil)
  if err != nil {
		panic(err)
  }

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateUser(t *testing.T) {
	clearTable()
	nt := NewPrUser()
	uid := addPrUser(nt)

	statement := fmt.Sprintf(prUserURL, uid)
	req, _ := http.NewRequest("GET", statement, nil)
	response := executeRequest(req)

	json.Unmarshal(response.Body.Bytes(), &nt)

	ut := nt

	//Update the new struct
	ut.Active = "eslaf"

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
}

func TestDeletePrUser(t *testing.T) {
	clearTable()
	nt := NewPrUser()
	uid := addPrUser(nt)

	statement := fmt.Sprintf(prUserURL, uid)
	req, _ := http.NewRequest("DELETE", statement, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", statement, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestDumpPrUser(t *testing.T) {
	nt := NewPrUser()

  err := dumpUser(*nt)

	if err != nil {
		t.Errorf("Expected erro to be nill. Got %v", err)
	}
}


