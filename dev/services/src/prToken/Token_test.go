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
	_ "reflect"
	"testing"
)
const (
  prUpdated     string = "updated"
  prCreated     string = "created"
  prActive      string = "active"
  prMeta        string = "metadata"
  prTokenURL    string = "/api/v1/namespace/pavedroad.io/prTokens/%s"
)

var a prTokenApp

func TestMain(m *testing.M) {
	a = prTokenApp{}
	a.Initialize()

	clearDB()
	ensureTableExists()

	code := m.Run()

  fmt.Println("cleartable")
	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
  fmt.Println(tableCreationQuery)
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		fmt.Println("Table check failed:", err)
		log.Fatal(err)
	}

	if _, err := a.DB.Exec(indexCreateToken); err != nil {
		fmt.Println("Index check failed")
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM kevlarweb.prToken")
        a.DB.Exec("ALTER TABLE kevlarweb.prToken AUTO_INCREMENT = 1")

}

func clearDB() {
	a.DB.Exec("DROP DATABASE IF EXISTS kevlarweb")
	a.DB.Exec("CREATE DATABASE kevlarweb")
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS kevlarweb.prToken (
    uid UUID DEFAULT uuid_v4()::UUID PRIMARY KEY,
    prToken JSONB
);`

const indexCreateToken = `
CREATE INDEX IF NOT EXISTS prTokenIdx ON kevlarweb.prToken USING GIN (prToken)`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/v1/namespace/pavedroad.io/prTokensLIST", nil)
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
		"/api/v1/namespace/pavedroad.io/prTokens/43ae99c9-5ca9-4691-b824-fbd2673946a7", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "UID: 43ae99c9-5ca9-4691-b824-fbd2673946a7 does not exist" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Token not found'. Got '%s'", m["error"])
	}
}

func TestCreateToken(t *testing.T) {
	clearTable()

	payload := []byte(`{"apiVersion" :"core.pavedroad.io/v1alpha1", "kind" : "PrToken",
		"metadata" : { "name" : "testuser-github-token", "namespace" : "jscharber",
		"uid" : "43ae99c9-5ca9-4691-b824-fbd2673946a7", "site" : "github",
    "endPoint": "https://api.github.com", "token" : "mytoken", 
    "scope" : ["user", "repo"] }, "created" : "2012-04-21T18:25:43-05:00",
    "updated" : "2012-04-21T18:25:43-05:00", "active"  : true }`)

	req, _ := http.NewRequest("POST", "/api/v1/namespace/pavedroad.io/prTokens", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	//var md map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["apiVersion"] != "core.pavedroad.io/v1alpha1" {
		t.Errorf("Expected apiVersion to be 'core.pavedroad.io/v1alpha1'. Got '%v'", m["apiVersion"])
	}

	if m["kind"] != "PrToken" {
		t.Errorf("Expected apiVersion to be 'PrToken'. Got '%v'", m["kind"])
	}

  // go see this as a date and the conversation time will not be the same as the
  // original data/time
  if m[prCreated] == "" {
    t.Errorf("Expected created to be set. Got '%v'", m[prCreated])
	}

  if m[prUpdated] == "" {
    t.Errorf("Expected updated to be set. Got '%v'", m[prCreated])
	}

	if m[prActive] != true {
		t.Errorf("Expected active to be 'true'. Got '%v'", m[prActive])
	}

	testName := m[prMeta].(map[string]interface{})["name"]
	if testName != "testuser-github-token" {
		t.Errorf("Expected user name to be 'testuser-github-token'. Got '%v'", testName)
	}

	//testName = m["metadata"].(map[string]interface{})["uid"]
  //fmt.Println(m["metadata"])
  //fmt.Println(testName)
	//if testName != "4601415d-7d00-4707-8fba-151eb9a0a75e" {
		//t.Errorf("Expected user uid to be '4601415d-7d00-4707-8fba-151eb9a0a75e'.\nGot %#v", testName)
	//}

	testName = m[prMeta].(map[string]interface{})["namespace"]
	if testName != "jscharber" {
		t.Errorf("Expected namespace to be 'jscharber'. Got '%v'", testName)
	}

	testName = m[prMeta].(map[string]interface{})["site"]
	if testName != "github" {
		t.Errorf("Expected site to be 'github'. Got '%v'", testName)
	}

	testName = m[prMeta].(map[string]interface{})["endPoint"]
	if testName != "https://api.github.com" {
		t.Errorf("Expected endPoint to be 'https://api.github.com'. Got '%v'", testName)
	}

	testName = m[prMeta].(map[string]interface{})["token"]
	if testName != "mytoken" {
		t.Errorf("Expected token to be 'mytoken'. Got '%v'", testName)
	}

	testName = m[prMeta].(map[string]interface{})["scope"]
	if testName == "" {
		t.Errorf("Expected scope to be set")
	}

}

func addToken(t *PrToken) (uid string) {
  jb, err := json.Marshal(t)
  if err != nil {
    panic(err)
  }

  statement := fmt.Sprintf("INSERT INTO kevlarweb.prtoken(prtoken) VALUES('%s') RETURNING uid", jb)
  rows, er1 := a.DB.Query(statement)

  if er1 != nil {
    fmt.Println("Insert failed")
    return ""
  }

  defer rows.Close()

  for rows.Next() {
    err := rows.Scan(&t.Metadata.UID)
    if err != nil {
      return ""
    }
  }
  return t.Metadata.UID
}

//create a token, set default values
func NewToken () (t* PrToken) {
  Token := PrToken{}

  Token.APIVersion = "core.pavedroad.io/v1alpha1"
  Token.Kind = "PrToken"
  Token.Metadata.Name = "testoken"
  Token.Metadata.UID = ""
  Token.Metadata.Site = "github"
  Token.Metadata.EndPoint = "https://api.github.com"
  Token.Metadata.Token = "#####################"
  Token.Metadata.Scope =  append(Token.Metadata.Scope, "user")
  Token.Metadata.Scope =  append(Token.Metadata.Scope, "repo")
  Token.Created = ""
  Token.Updated = ""
  Token.Active = true

  return &Token
}

//test getting a token
func TestGetToken(t *testing.T) {
	clearTable()
  nt := NewToken()
  uid := addToken(nt)

  statement := fmt.Sprintf("/api/v1/namespace/pavedroad.io/prTokens/%s", uid)

	req, _ := http.NewRequest("GET", statement, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateUser(t *testing.T) {
	clearTable()
  nt := NewToken()
  uid := addToken(nt)

  statement := fmt.Sprintf(prTokenURL, uid)
	req, _ := http.NewRequest("GET", statement, nil)
	response := executeRequest(req)

	//var originalToken map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &nt)

  ut := nt

  //Update the new struct
	ut.Active = false
	ut.Metadata.Name  = "updatedname"

  jb, err := json.Marshal(ut)
  if err != nil {
    panic(err)
  }

	req, _ = http.NewRequest("PUT", statement, strings.NewReader(string(jb)))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["active"] != false {
		t.Errorf("Expected active to be faluse. Got %v", m["id"])
	}

  testName := m[prMeta].(map[string]interface{})["name"]
  if testName != "updatedname" {
    t.Errorf("Expected name to be 'updatedname'. Got '%v'", testName)
  }

}

func TestDeleteToken(t *testing.T) {
	clearTable()
  nt := NewToken()
  uid := addToken(nt)

  statement := fmt.Sprintf(prTokenURL, uid)
	req, _ := http.NewRequest("DELETE", statement, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", statement, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
