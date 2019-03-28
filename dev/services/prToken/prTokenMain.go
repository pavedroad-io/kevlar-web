// prTokenMain.go
// entry point for prToken Microservice

package main

import "os"

func main() {
	a := prTokenApp{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	a.Run(":8080")
}
