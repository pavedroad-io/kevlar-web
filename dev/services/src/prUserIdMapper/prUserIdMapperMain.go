// prUserIdMapperMain.go
// entry point for prUserIdMapper microservice

package main

import(
  "github.com/gorilla/mux"
  "database/sql"
  "log"
  "time"
)

// Contants to build up a k8s style URL
const (
  // prUserIdMapperAPIVersion Version API URL
  prUserIdMapperAPIVersion       string = "/api/v1"
  // prUserIdMapperNamespaceID Prefix for namespaces
  prUserIdMapperNamespaceID      string = "namespace"
  // prUserIdMapperDefaultNamespace Default namespace
  prUserIdMapperDefaultNamespace string = "pavedroad.io"
  // prUserIdMapperResourceType CRD Type per k8s
  prUserIdMapperResourceType     string = "prUserIdMappers"
  // The email or account login used by 3rd parth provider
  prUserCred                     string ="/{key}"
)

// holds pointers to database and http server
type prUserIdMapperApp struct {
  Router *mux.Router
  DB     *sql.DB
}

// both db and http configuration can be changed using environment varialbes
type  databaseConfig struct {
  username string
  password string
  database string
  sslMode  string
  dbDriver string
  ip string
  port string
}

// HTTP server configuration
type httpConfig struct {
  ip string
  port string
  shutdownTimeout time.Duration
  readTimeout time.Duration
  writeTimeout time.Duration
  listenString string
  logPath string
}

// Global for use in the module

// Set default database configuration
var dbconf = databaseConfig{username: "root", password: "", database: "pavedroad", sslMode: "disable", dbDriver: "postgres", ip: "127.0.0.1", port: "26257"}

// Set default http configuration
var httpconf = httpConfig {ip: "127.0.0.1", port: "8082", shutdownTimeout: 15, readTimeout: 60, writeTimeout: 60, listenString: "127.0.0.1:8082", logPath: "logs/userIdMapper.log"}

// shutdownTimeout will be initialized based on the default or HTTP_SHUTDOWN_TIMEOUT
var shutdowTimeout time.Duration

// main entry point for server
func main() {

  // Setup loggin
  openLogFile(httpconf.logPath)
  log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
  log.Printf("Logfile opened %s", httpconf.logPath)

	a := prUserIdMapperApp{}
	a.Initialize()
  a.Run(httpconf.listenString)
}
