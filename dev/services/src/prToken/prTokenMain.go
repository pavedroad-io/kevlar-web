// prTokenMain.go
// entry point for prToken Microservice

package main

import "fmt"

// environment variables:
// See prTokenApp.go for defaults
// APP_DB_USERNAME APP_DB_PASSWORD APP_DB_NAME APP_DB_SSL_MODE APP_DB_IP
// APP_DB_PORT
//
// Server addressing environment variables
// IP_ADDR IP_PORT

func main() {
	a := prTokenApp{}
	a.Initialize()
  fmt.Println("Starting server listing on: " + ServerString)
  a.Run(ServerString)
}
