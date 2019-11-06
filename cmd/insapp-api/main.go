package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	insapp "github.com/thomas-bouvier/insapp-go"
)

type withCORS struct {
	r *mux.Router
}

func main() {
	config := insapp.InitConfig()

	log.Println("Starting server on 0.0.0.0:" + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, &withCORS{insapp.NewRouter()}))
}

// Simple wrapper to Allow CORS
func (s *withCORS) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")
	res.Header().Set("Access-Control-Allow-Origin", origin)
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	res.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Origin")
	res.Header().Set("Access-Control-Expose-Headers", "Content-Range")

	// Stop here for a Preflighted OPTIONS request.
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(res, req)
}
