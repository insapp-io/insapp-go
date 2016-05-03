package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type WithCORS struct {
	r *mux.Router
}

func main() {
	log.Println("Starting server on 0.0.0.0:9000")
	log.Fatal(http.ListenAndServe(":9000", &WithCORS{NewRouter()}))
}

// Simple wrapper to Allow CORS.
// See: https://groups.google.com/forum/#!topic/golang-nuts/-Sh616lXNRE
// See: http://stackoverflow.com/a/24818638/1058612.
func (s *WithCORS) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		res.Header().Set("Access-Control-Allow-Origin", origin)
		res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		res.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}

	// Stop here for a Preflighted OPTIONS request.
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(res, req)
}
