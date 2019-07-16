// prUserApp
// Manage pavedroad users
// jscharber

package main

import (
	"context"
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
  "github.com/google/uuid"
)

// Initialize setups database connection object and the http server
//
func (a *prUserApp) Initialize() {

	// Override defaults
	a.initializeEnvironment()

	// Build connection strings
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s",
		dbconf.username,
		dbconf.password,
		dbconf.database,
		dbconf.sslMode,
		dbconf.ip,
		dbconf.port)

	httpconf.listenString = fmt.Sprintf("%s:%s", httpconf.ip, httpconf.port)

	var err error
	a.DB, err = sql.Open(dbconf.dbDriver, connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Start the server
func (a *prUserApp) Run(addr string) {

	log.Println("Listing at: " + addr)
	srv := &http.Server{
		Handler:      a.Router,
		Addr:         addr,
		WriteTimeout: httpconf.writeTimeout * time.Second,
		ReadTimeout:  httpconf.readTimeout * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// Listen for SIGHUP
	c := make(chan os.Signal, 1)
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), httpconf.shutdownTimeout)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

// Get for ennvironment variable overrides
func (a *prUserApp) initializeEnvironment() {
	var envVar = ""

	//look for environment variables overrides
	envVar = os.Getenv("APP_DB_USERNAME")
	if envVar != "" {
		dbconf.username = envVar
	}

	envVar = os.Getenv("APP_DB_PASSWORD")
	if envVar != "" {
		dbconf.password = envVar
	}

	envVar = os.Getenv("APP_DB_NAME")
	if envVar != "" {
		dbconf.database = envVar
	}
	envVar = os.Getenv("APP_DB_SSL_MODE")
	if envVar != "" {
		dbconf.sslMode = envVar
	}

	envVar = os.Getenv("APP_DB_SQL_DRIVER")
	if envVar != "" {
		dbconf.dbDriver = envVar
	}

	envVar = os.Getenv("APP_DB_IP")
	if envVar != "" {
		dbconf.ip = envVar
	}

	envVar = os.Getenv("APP_DB_PORT")
	if envVar != "" {
		dbconf.port = envVar
	}

	envVar = os.Getenv("HTTP_IP_ADDR")
	if envVar != "" {
		httpconf.ip = envVar
	}

	envVar = os.Getenv("HTTP_IP_PORT")
	if envVar != "" {
		httpconf.port = envVar
	}

	envVar = os.Getenv("HTTP_READ_TIMEOUT")
	if envVar != "" {
		to, err := strconv.Atoi(envVar)
		if err == nil {
			log.Printf("failed to convert HTTP_READ_TIMEOUT: %s to int", envVar)
		} else {
			httpconf.readTimeout = time.Duration(to) * time.Second
		}
		log.Printf("Read timeout: %d", httpconf.readTimeout)
	}

	envVar = os.Getenv("HTTP_WRITE_TIMEOUT")
	if envVar != "" {
		to, err := strconv.Atoi(envVar)
		if err == nil {
			log.Printf("failed to convert HTTP_READ_TIMEOUT: %s to int", envVar)
		} else {
			httpconf.writeTimeout = time.Duration(to) * time.Second
		}
		log.Printf("Write timeout: %d", httpconf.writeTimeout)
	}

	envVar = os.Getenv("HTTP_SHUTDOWN_TIMEOUT")
	if envVar != "" {
		if envVar != "" {
			to, err := strconv.Atoi(envVar)
			if err != nil {
				httpconf.shutdownTimeout = time.Second * time.Duration(to)
			} else {
				httpconf.shutdownTimeout = time.Second * httpconf.shutdownTimeout
			}
			log.Println("Shutdown timeout", httpconf.shutdownTimeout)
		}
	}

	envVar = os.Getenv("HTTP_LOG")
	if envVar != "" {
		httpconf.logPath = envVar
	}

}

// endpoints for prUser microservice
// k8s stype URI
// api/v1/namespace/{namespace}/resourcetype
// default namespace is pavedroad.io
//
// List mappngs
// GET /api/v1/namespace/pavedroad.io/prUsersLIST
//
// Map a user cred into uuid
// GET /api/v1/namespace/pavedroad.io/prUsers/{key}
//
// create a new mapping
// POST /api/v1/namespace/pavedroad.io/prUsers
//
// update a mapping
// PUT /api/v1/namespace/pavedroad.io/prUsers/{key}
//
// partial update of a mapping
// PATCH /api/v1/namespace/pavedroad.io/prUsers/{key}
//
// Delete a mapping
// DELETE /api/v1/namespace/pavedroad.io/prUsers/{key}
func (a *prUserApp) initializeRoutes() {
	uri := prUserAPIVersion + "/" + prUserNamespaceID + "/{namespace}/" +
		prUserResourceType + "LIST"
	a.Router.HandleFunc(uri, a.getUsers).Methods("GET")
	fmt.Println("LIST" + uri)

	uri = prUserAPIVersion + "/" + prUserNamespaceID + "/{namespace}/" +
		prUserResourceType + "/{key}"
	a.Router.HandleFunc(uri, a.getUser).Methods("GET")
	fmt.Println("GET" + uri)

	uri = prUserAPIVersion + "/" + prUserNamespaceID + "/{namespace}/" + prUserResourceType
	a.Router.HandleFunc(uri, a.createUser).Methods("POST")
	fmt.Println("POST" + uri)

	uri = prUserAPIVersion + "/" + prUserNamespaceID + "/{namespace}/" +
		prUserResourceType + prUserUUID
	a.Router.HandleFunc(uri, a.updateUser).Methods("PUT")
	fmt.Println("PUT" + uri)

	uri = prUserAPIVersion + "/" + prUserNamespaceID + "/{namespace}/" +
		prUserResourceType + prUserUUID
	a.Router.HandleFunc(uri, a.deleteUser).Methods("DELETE")
	fmt.Println("DELETE" + uri)
}

// getUsers
// return a list of all users
func (a *prUserApp) getUsers(w http.ResponseWriter, r *http.Request) {
	User := prUser{}

	//vars := mux.Vars(r)
	//fmt.Println("list users: ", vars)

	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	users, err := User.getUsers(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

// getUser: return a user given a key
func (a *prUserApp) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	User := prUser{}
  var err error

	switch option := r.FormValue(prLookUp); option {
	case "name":
	  err = User.getUser(a.DB, vars["key"], NAME)
		break

  default:
     // Postgres doesn't handle invalid UUIDs
     if _, er := uuid.Parse(vars["key"]); er != nil {
		    respondWithError(w, http.StatusBadRequest, er.Error())
		    return
     }
	   err = User.getUser(a.DB, vars["key"], UUID)
	}

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, User)
}

// createUser
// Use POST to create a user
func (a *prUserApp) createUser(w http.ResponseWriter, r *http.Request) {
	// New map structure
	user := prUser{}

	htmlData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal(htmlData, &user)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	//dumpUser(user)

	//var out bytes.Buffer
	//json.Indent(&out, htmlData, "", "\t")
	//out.WriteTo(os.Stdout)

	ct := time.Now().UTC()
	user.Created = ct.Format(time.RFC3339)
	user.Updated = ct.Format(time.RFC3339)

	// Save into backend storage
	if _, err := user.createUser(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

// updateUser Use PUT to update a user
func (a *prUserApp) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := prUser{}

	htmlData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(htmlData, &user)
	if err != nil {
		log.Println(err)
		return
	}

	ct := time.Now().UTC()
	user.Metadata.UUID = vars["ke"]
	user.Updated = ct.Format(time.RFC3339)

	if err := user.updateUser(a.DB, vars["key"]); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	respondWithJSON(w, http.StatusOK, user)

}

// deleteUser
// Use DELETE to delete a user
func (a *prUserApp) deleteUser(w http.ResponseWriter, r *http.Request) {
	User := prUser{}
	vars := mux.Vars(r)

	err := User.deleteUser(a.DB, vars["key"])
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

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func openLogFile(logfile string) {
	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

		if err != nil {
			log.Fatal("OpenLogfile: os.OpenFile:", err)
		}

		log.SetOutput(lf)
	}
}

/*
func dumpUser(m prUser) {
	fmt.Println("Dump prUserIdMapp")
	fmt.Printf("apiVersion: %s\n", m.APIVersion)
	fmt.Printf("objVersion: %s\n", m.ObjVersion)
	fmt.Printf("kind: %s\n", m.Kind)
	fmt.Printf("credential: %s\n", m.Credential)
	fmt.Printf("userUUID: %s\n", m.UserUUID)
	fmt.Printf("loginCount: %d\n", m.LoginCount)
	fmt.Printf("created: %s\n", m.Created)
	fmt.Printf("updated: %s\n", m.Updated)
	fmt.Printf("active: %s\n", m.Active)
}
*/
