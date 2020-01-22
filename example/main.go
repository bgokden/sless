package main

import (
	"fmt"
	"net/http"

	goappenv "github.com/bgokden/go-app-env"
)

type ExampleGoApp struct{}

// GetName is used for service registery
func (e *ExampleGoApp) GetName() string {
	return "example-go-app"
}

// GetTag will be used for versioning
func (e *ExampleGoApp) GetTag() string {
	return "v0.1.0"
}

// GetDependencies will be used for initializing other services
func (e *ExampleGoApp) GetDependencies() []string {
	return []string{}
}

// RunWithEnv is the main loop that will be initialized.
func (e *ExampleGoApp) RunWithEnv(goappenv goappenv.GoAppEnv) error {
	logger := goappenv.GetLogger()
	logger.Printf("I am running in env %v\n", goappenv.GetName())
	mux := goappenv.GetHttpServer()
	mux.HandleFunc("/example0", serve)
	return nil
}

const text = "server example from github\n"

func serve(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, text)
}

func main() {
	exampleGoApp := &ExampleGoApp{}
	err := exampleGoApp.RunWithEnv(goappenv.Base())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	select {}
}

var GoApp ExampleGoApp
