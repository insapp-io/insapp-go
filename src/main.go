package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting server on 0.0.0.0:9000")
	log.Fatal(http.ListenAndServe(":9000", NewRouter()))
}
