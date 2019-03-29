// prTokenMain.go
// entry point for prToken Microservice

package main

import "os"

// Set database authentication environment variables:
// APP_DB_USERNAME
// APP_DB_PASSWORD
// APP_DB_NAME
//
// Server addressing environment variables
// IP_ADDR
// IP_PORT

func main() {
	a := prTokenApp{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("IP_ADDR"),
		os.Getenv("IP_PORT"))
    // TODO(jscharber): Update run to use ip/port envvars
	a.Run(":8080")
}
