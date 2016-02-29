package main

import (
	"fmt"
	"net/http"
)

// Index is just a test actually
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Insapp REST API - v.0.1")
}
