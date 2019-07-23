// prUserMain.go
// entry point for prUser microservice

package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"time"
)

// Contants to build up a k8s style URL
// prUserIdMapperAPIVersion Version API URL
// prUserNamespaceID Prefix for namespaces
// prUserDefaultNamespace Default namespace
// prUserResourceType CRD Type per k8s
// UUID for the user
const (
	prUserAPIVersion       string = "/api/v1"
	prUserNamespaceID      string = "namespace"
	prUserDefaultNamespace string = "pavedroad.io"
	prUserResourceType     string = "prUsers"
	prUserUUID             string = "/{key}"
	prLookUp               string = "query"
)

// Options for looking up a user
const (
	UUID = iota
	NAME
)

// holds pointers to database and http server
type prUserApp struct {
	Router *mux.Router
	DB     *sql.DB
}

// both db and http configuration can be changed using environment varialbes
type databaseConfig struct {
	username string
	password string
	database string
	sslMode  string
	dbDriver string
	ip       string
	port     string
}

// HTTP server configuration
type httpConfig struct {
	ip              string
	port            string
	shutdownTimeout time.Duration
	readTimeout     time.Duration
	writeTimeout    time.Duration
	listenString    string
	logPath         string
}

// Global for use in the module

// Set default database configuration
var dbconf = databaseConfig{username: "root", password: "", database: "pavedroad", sslMode: "disable", dbDriver: "postgres", ip: "127.0.0.1", port: "26257"}

// Set default http configuration
var httpconf = httpConfig{ip: "127.0.0.1", port: "8083", shutdownTimeout: 15, readTimeout: 60, writeTimeout: 60, listenString: "127.0.0.1:8083", logPath: "logs/user.log"}

// shutdownTimeout will be initialized based on the default or HTTP_SHUTDOWN_TIMEOUT
var shutdowTimeout time.Duration

// main entry point for server
func main() {

	// Setup loggin
	openLogFile(httpconf.logPath)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("Logfile opened %s", httpconf.logPath)

	a := prUserApp{}
	a.Initialize()
	a.Run(httpconf.listenString)
}
