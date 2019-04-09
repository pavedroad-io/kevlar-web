// prTokenMain.go
// entry point for prToken Microservice

package main

import "os"

//import "fmt"

// Set database authentication environment variables:
// APP_DB_USERNAME
// APP_DB_PASSWORD
// APP_DB_NAME
// APP_DB_SSL_MODE
// APP_DB_IP
// APP_DB_PORT
// 172.18.0.2
//
// Server addressing environment variables
// IP_ADDR
// IP_PORT

func main() {
	a := prTokenApp{}

	username := os.Getenv("APP_DB_USERNAME")
	if username == "" {
		username = "root"
	}

	password := os.Getenv("APP_DB_PASSWORD")
	if password == "" {
		password = ""
	}

	database := os.Getenv("APP_DB_NAME")
	if database == "" {
		database = "kevlar-web"
	}

	sslMode := os.Getenv("APP_DB_SSL_MODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	dbDriver := os.Getenv("APP_DB_SQL_DRIVER")
	if dbDriver == "" {
		dbDriver = "postgres"
	}

	dbIp := os.Getenv("APP_DB_IP")
	if dbIp == "" {
		dbIp = "127.0.0.1"
	}

	dbPort := os.Getenv("APP_DB_PORT")
	if dbPort == "" {
		dbPort = "26257"
	}

	serverAddr := os.Getenv("IP_ADDR")
	if serverAddr == "" {
		serverAddr = "127.0.0.1"
	}

	serverPort := os.Getenv("IP_PORT")
	if serverPort == "" {
		serverPort = "8081"
	}

	a.Initialize(username,
		password,
		database,
		sslMode,
		dbDriver,
		dbIp,
		dbPort,
		serverAddr,
		serverPort)

	// TODO(jscharber): Update run to use ip/port envvars
	a.Run(":8081")
}
