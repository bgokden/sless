package main

import (
	"fmt"
	"net/http"
)

type server string

func (s server) Serve(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "server example\n")
}

// exported
var Server server
