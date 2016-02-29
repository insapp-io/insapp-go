package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

var routes = Routes{
	Route{"Index", "GET", "/", Index},

	Route{"GetAssociation", "GET", "/association", GetAllAssociationsController},
	Route{"GetAssociation", "GET", "/association/{id}", GetAssociationController},
	Route{"AddAssociation", "POST", "/association", AddAssociationController},
	Route{"UpdateAssociation", "PUT", "/association/{id}", UpdateAssociationController},
	//Route{"fdsq", "POST", "/association/{id}/{idEvent}", AddEventTest},
	//Route{"fdsq", "DELETE", "/association/{id}/{idEvent}", RemoveEventTest},
	Route{"GetFutureEvents", "GET", "/event", GetFutureEventsController},
	Route{"GetEvent", "GET", "/event/{id}", GetEventController},
	Route{"AddEvent", "POST", "/event", AddEventController},
	Route{"UpdateEvent", "PUT", "/event/{id}", UpdateEventController},
	Route{"DeleteEvent", "DELETE", "/event/{id}", DeleteEventController},
	Route{"AddParticipant", "POST", "/event/{id}/participant/{userID}", AddParticipantController},
	Route{"RemoveParticipant", "DELETE", "/event/{id}/participant/{userID}", RemoveParticipantController},
}
