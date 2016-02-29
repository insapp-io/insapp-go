package main

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

func GetEventController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assocationID := vars["id"]
	var res = GetEvent(bson.ObjectIdHex(assocationID))
	json.NewEncoder(w).Encode(res)
}

func GetFutureEventsController(w http.ResponseWriter, r *http.Request) {
	var res = GetFutureEvent()
	json.NewEncoder(w).Encode(res)
}

func AddEventController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var event Event
	decoder.Decode(&event)
	res := AddEvent(event)
	json.NewEncoder(w).Encode(res)
}

func UpdateEventController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var event Event
	decoder.Decode(&event)
	vars := mux.Vars(r)
	eventID := vars["id"]
	res := UpdateEvent(bson.ObjectIdHex(eventID), event)
	json.NewEncoder(w).Encode(res)
}

func DeleteEventController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := GetEvent(bson.ObjectIdHex(vars["id"]))
	res := DeleteEvent(event)
	json.NewEncoder(w).Encode(res)
}

func AddParticipantController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := bson.ObjectIdHex(vars["id"])
	userID := bson.ObjectIdHex(vars["userID"])
	res := AddParticipant(event, userID)
	json.NewEncoder(w).Encode(res)
}

func RemoveParticipantController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := bson.ObjectIdHex(vars["id"])
	userID := bson.ObjectIdHex(vars["userID"])
	res := RemoveParticipant(event, userID)
	json.NewEncoder(w).Encode(res)
}
