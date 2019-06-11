// prTokenApp
// Business logic specific to managing GitHub tokens for kevlar web users
// jscharber

package main

import (
	_ "bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
  // PrTokenAPIVersion Version API URL
	PrTokenAPIVersion       string = "/api/v1"
  // PrTokenNamespaceID Prefix for namespaces
	PrTokenNamespaceID      string = "namespace"
  // PrTokenDefaultNamespace Default namespace
	PrTokenDefaultNamespace string = "pavedroad.io"
  // PrTokenResourceType CRD Type per k8s
	PrTokenResourceType     string = "prTokens"
)

type prTokenApp struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *prTokenApp) Initialize(user, password, dbname string, sslmode string, sqldriver string,
	dbIP string, dbPort string, ipAddr string, ipPort string) {

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s",
		user, password, dbname, sslmode, dbIP, dbPort)

	//fmt.Println(sqldriver, connectionString)
	var err error
	a.DB, err = sql.Open(sqldriver, connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *prTokenApp) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// endpoints for prToken microservice
// k8s stype URI
// api/v1/namespace/{namespace}/resourcetype
// default namespace is pavedroad.io
//
// Return a list of tokens
// GET /api/v1/namespace/pavedroad.io/prTokensLIST
//
// Get a specific token
// GET /api/v1/namespace/pavedroad.io/prTokens/{uid}
//
// create a new token
// POST /api/v1/namespace/pavedroad.io/prTokens
//
// update a specific toke
// PUT /api/v1/namespace/pavedroad.io/prTokens/{uid}
//
// partial update of a specific toke
// PATCH /api/v1/namespace/pavedroad.io/prTokens/{uid}
//
// Delete a specific toke
// DELETE /api/v1/namespace/pavedroad.io/prTokens/{uid}
func (a *prTokenApp) initializeRoutes() {
	//Get list of tokens
	uri := PrTokenAPIVersion + "/" + PrTokenNamespaceID + "/{namespace}/" +
		PrTokenResourceType + "LIST"
	a.Router.HandleFunc(uri, a.getTokens).Methods("GET")
	fmt.Println("LIST" + uri)

	//Get a token
	uri = PrTokenAPIVersion + "/" + PrTokenNamespaceID + "/{namespace}/" +
		PrTokenResourceType + "/{uid}"
	a.Router.HandleFunc(uri, a.getToken).Methods("GET")
	fmt.Println("GET" + uri)

	//Create a token
	uri = PrTokenAPIVersion + "/" + PrTokenNamespaceID + "/{namespace}/" +
		PrTokenResourceType
	a.Router.HandleFunc(uri, a.createToken).Methods("POST")
	fmt.Println("POST" + uri)

	//update a token
	uri = PrTokenAPIVersion + "/" + PrTokenNamespaceID + "/{namespace}/" +
		PrTokenResourceType + "/{uid}"
	a.Router.HandleFunc(uri, a.updateToken).Methods("PUT")
	fmt.Println("PUT" + uri)

	//delete a token
	uri = PrTokenAPIVersion + "/" + PrTokenNamespaceID + "/{namespace}/" +
		PrTokenResourceType + "/{uid}"
	a.Router.HandleFunc(uri, a.deleteToken).Methods("DELETE")
	fmt.Println("DELETE" + uri)

}

// getTokens
// return a list of all tokens
func (a *prTokenApp) getTokens(w http.ResponseWriter, r *http.Request) {
	Token := PrToken{}

	//vars := mux.Vars(r)
	//fmt.Println("list tokens: ", vars)

	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	tokens, err := Token.getTokens(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tokens)
}

// getToken: return a token given a UID
func (a *prTokenApp) getToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Token := PrToken{}

	err := Token.getToken(a.DB, vars["uid"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Token)
}

// createToken
// Use POST to create a token
func (a *prTokenApp) createToken(w http.ResponseWriter, r *http.Request) {
	// New token structure
	Token := PrToken{}

	// Read URI variables
	//vars := mux.Vars(r)
	//fmt.Println(vars)

	htmlData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal(htmlData, &Token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//var out bytes.Buffer
	//json.Indent(&out, htmlData, "", "\t")
	//out.WriteTo(os.Stdout)

	ct := time.Now().UTC()
	Token.Created = ct.Format(time.RFC3339)
	Token.Updated = ct.Format(time.RFC3339)

	if err := Token.createToken(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	//fmt.Println(Token)
	respondWithJSON(w, http.StatusCreated, Token)
}

// updateToken
// Use PUT to update a token
func (a *prTokenApp) updateToken(w http.ResponseWriter, r *http.Request) {
	// New token structure
	Token := PrToken{}

	// Read URI variables
	vars := mux.Vars(r)

	htmlData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(htmlData, &Token)
	if err != nil {
		fmt.Println(err)
		return
	}

	if vars["uid"] != Token.Metadata.UID {
		//TODO(jscharber): Change to log message
		em := "UID: " + vars["uid"] + "in does not match payload [" + Token.Metadata.UID + "]"
		fmt.Println(em)
		return
	}

	ct := time.Now().UTC()
	Token.Updated = ct.Format(time.RFC3339)

	if err := Token.updateToken(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	//fmt.Println(Token)
	respondWithJSON(w, http.StatusOK, Token)

}

// deleteToken
// Use DELETE to delete a token
func (a *prTokenApp) deleteToken(w http.ResponseWriter, r *http.Request) {
	Token := PrToken{}
	vars := mux.Vars(r)

	err := Token.deleteToken(a.DB, vars["uid"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
