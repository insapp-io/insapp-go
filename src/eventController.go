package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

// GetEventController will answer a JSON of the event
// from the given "id" in the URL. (cf Routes in routes.go)
func GetEventController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assocationID := vars["id"]
	var res = GetEvent(bson.ObjectIdHex(assocationID))
	json.NewEncoder(w).Encode(res)
}

// GetFutureEventsController will answer a JSON
// containing all future events from "NOW"
func GetFutureEventsController(w http.ResponseWriter, r *http.Request) {
	var res = GetFutureEvents()
	json.NewEncoder(w).Encode(res)
}

// AddEventController will answer the JSON
// of the brand new created event from the JSON body
func AddEventController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var event Event
	decoder.Decode(&event)
	res := AddEvent(event)
	asso := GetAssociation(event.Association)
	TriggerNotification(asso.Name + " vient de poster un nouvel évènement.")
	json.NewEncoder(w).Encode(res)
}

// UpdateEventController will answer the JSON
// of the brand new modified event from the JSON body
func UpdateEventController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var event Event
	decoder.Decode(&event)
	vars := mux.Vars(r)
	eventID := vars["id"]
	res := UpdateEvent(bson.ObjectIdHex(eventID), event)
	json.NewEncoder(w).Encode(res)
}

// DeleteEventController will answer an empty JSON
// if the deletation has succeed
func DeleteEventController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := GetEvent(bson.ObjectIdHex(vars["id"]))
	res := DeleteEvent(event)
	json.NewEncoder(w).Encode(res)
}

// AddParticipantController will answer the JSON
// of the event with the given partipant added
func AddParticipantController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := bson.ObjectIdHex(vars["id"])
	userID := bson.ObjectIdHex(vars["userID"])
	event, user := AddParticipant(eventID, userID)
	json.NewEncoder(w).Encode(bson.M{"event": event, "user": user})
}

// RemoveParticipantController will answer the JSON
// of the event without the given partipant added
func RemoveParticipantController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := bson.ObjectIdHex(vars["id"])
	userID := bson.ObjectIdHex(vars["userID"])
	event, user := RemoveParticipant(eventID, userID)
	json.NewEncoder(w).Encode(bson.M{"event": event, "user": user})
}

// AddImageEventController will set the image of the event and return the event
func AddImageEventController(w http.ResponseWriter, r *http.Request) {
	fileName := UploadImage(r)
	if fileName == "error" {
		w.Header().Set("status", "400")
		fmt.Fprintln(w, "{}")
	} else {
		vars := mux.Vars(r)
		res := SetImageEvent(bson.ObjectIdHex(vars["id"]), fileName)
		json.NewEncoder(w).Encode(res)
	}
}
