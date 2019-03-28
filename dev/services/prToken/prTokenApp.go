// prTokenApp
// Business logic specific to managing GitHub tokens for kevlar web users
// jscharber

package main

import (
	_ "encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_ "strconv"
)

const (
	PrTokenAPIVersion       string = "/api/v1"
	PrTokenNamespaceID      string = "namespace"
	PrTokenDefaultNamespace string = "pavedroad.io"
	PrTokenResourceType     string = "prToken"
)

type prTokenApp struct {
	Router *mux.Router
}

func (a *prTokenApp) Initialize(user, password, dbname string) {
	//connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	//var err error
	//a.DB, err = sql.Open("mysql", connectionString)
	//if err != nil {
	//log.Fatal(err)
	//}

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
// GET /api/v1/namespace/pavedroad.io/prTokenLIST
//
// Get a specific token
// GET /api/v1/namespace/pavedroad.io/prToken/{usertoken}
//
// create a new token
// POST /api/v1/namespace/pavedroad.io/prToken/{usertoken}
//
// update a specific toke
// PUT /api/v1/namespace/pavedroad.io/prToken/{usertoken}
//
// partial update of a specific toke
// PATCH /api/v1/namespace/pavedroad.io/prToken/{usertoken}
//
// Delete a specific toke
// DELETE /api/v1/namespace/pavedroad.io/prToken/{usertoken}

func (a *prTokenApp) initializeRoutes() {
    //Get list of tokens
	uri := PrTokenAPIVersion + "/" + PrTokenNamespaceID + "/{namespace}/" +
           PrTokenResourceType+"LIST"
	a.Router.HandleFunc(uri, a.getTokens).Methods("GET")
	fmt.Println("LIST"+uri)

    //Get a token
	uri = PrTokenAPIVersion + "/" + PrTokenNamespaceID + "/{namespace}/" +
          PrTokenResourceType + "/{usertoken}"
	a.Router.HandleFunc(uri, a.getToken).Methods("GET")
	fmt.Println("GET"+uri)

    //Create a token
	uri = PrTokenAPIVersion + "/" + PrTokenNamespaceID + "/{namespace}/" +
          PrTokenResourceType + "/{usertoken}"
	a.Router.HandleFunc(uri, a.createToken).Methods("POST")
	fmt.Println("POST"+uri)

    //update a token
	uri = PrTokenAPIVersion + "/" + PrTokenNamespaceID + "/{namespace}/" +
          PrTokenResourceType + "/{usertoken}"
	a.Router.HandleFunc(uri, a.updateToken).Methods("PUT")
	fmt.Println("PUT"+uri)

    //delete a token
	uri = PrTokenAPIVersion + "/" + PrTokenNamespaceID + "/{namespace}/" +
          PrTokenResourceType + "/{usertoken}"
	a.Router.HandleFunc(uri, a.deleteToken).Methods("DELETE")
	fmt.Println("DELETE"+uri)

}

func (a *prTokenApp) getTokens(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    fmt.Println(vars)
	// id, err := strconv.Atoi(vars["id"])
	//	if err != nil {
	//		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
	//		return
	//	}
	//
	//	u := user{ID: id}
	//	if err := u.getUser(a.DB); err != nil {
	//		switch err {
	//		case sql.ErrNoRows:
	//			respondWithError(w, http.StatusNotFound, "User not found")
	//		default:
	//			respondWithError(w, http.StatusInternalServerError, err.Error())
	//		}
	//		return
	//	}

	//respondWithJSON(w, http.StatusOK, u)
}

func (a *prTokenApp) getToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    fmt.Println(vars)
}

func (a *prTokenApp) createToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    fmt.Println(vars)
}

func (a *prTokenApp) updateToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    fmt.Println(vars)
}

func (a *prTokenApp) deleteToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    fmt.Println(vars)
}



