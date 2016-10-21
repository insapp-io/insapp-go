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

	conf, err := Configuration()

	if err != nil{
		log.Println(err)
		log.Fatal("[error] Error when parsing config file. Make sure the config file is valid : ")
		return
	}

	log.Println("Starting server on 0.0.0.0:" + conf.Port)
	log.Fatal(http.ListenAndServe(":" + conf.Port, &WithCORS{NewRouter()}))
}

// Simple wrapper to Allow CORS.
// See: https://groups.google.com/forum/#!topic/golang-nuts/-Sh616lXNRE
// See: http://stackoverflow.com/a/24818638/1058612.
func (s *WithCORS) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")
	res.Header().Set("Access-Control-Allow-Origin", origin)
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	res.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Origin")

	// Stop here for a Preflighted OPTIONS request.
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(res, req)
}
