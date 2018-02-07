package main

import (

	"github.com/kidinamoto01/sdk-test/app"
	"fmt"
)

// Entry point of the Go app

func main() {
	app := app.NewSearchApp()
	fmt.Println("Search app started. Running forever on :46658")
	app.RunForever()
}