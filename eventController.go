package insapp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// GetEventController will answer a JSON of the event
// from the given "id" in the URL. (cf Routes in routes.go)
func GetEventController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	associationID := vars["id"]
	var res = GetEvent(bson.ObjectIdHex(associationID))
	_ = json.NewEncoder(w).Encode(res)
}

// GetFutureEventsController will answer a JSON
// containing all future events from "NOW"
func GetFutureEventsController(w http.ResponseWriter, r *http.Request) {
	id, err := GetUserFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(bson.M{"error": "could not get user ID"})
		return
	}

	user := GetUser(id)
	os := GetNotificationUserForUser(id).Os
	events := GetFutureEvents()
	res := Events{}
	if user.ID != "" {
		for _, event := range events {
			if contains(strings.ToUpper(user.Promotion), event.Promotions) && (contains(os, event.Plateforms) || os == "") || len(event.Promotions) == 0 || len(event.Plateforms) == 0 {
				res = append(res, event)
			}
		}
	} else {
		res = events
	}
	json.NewEncoder(w).Encode(res)
}

func GetEventsForAssociationController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	associationID := vars["id"]
	events := GetEventsForAssociation(bson.ObjectIdHex(associationID))
	_ = json.NewEncoder(w).Encode(events)
}

// AddEventController will answer the JSON
// of the brand new created Event from the JSON body
// Should be protected
func AddEventController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var event Event
	_ = decoder.Decode(&event)

	res := AddEvent(event)
	association := GetAssociation(event.Association)
	_ = json.NewEncoder(w).Encode(res)
	go TriggerNotificationForEvent(event, association.ID, res.ID, "@"+strings.ToLower(association.Name)+" t'invite à "+res.Name+" 📅")
}

// UpdateEventController will answer the JSON
// of the modified Event from the JSON body
// Should be protected
func UpdateEventController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var event Event
	_ = decoder.Decode(&event)
	vars := mux.Vars(r)
	eventID := vars["id"]

	res := UpdateEvent(bson.ObjectIdHex(eventID), event)
	_ = json.NewEncoder(w).Encode(res)
}

// DeleteEventController will answer an empty JSON
// if the deletion has succeed
// Should be protected
func DeleteEventController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := GetEvent(bson.ObjectIdHex(vars["id"]))

	res := DeleteEvent(event)
	_ = json.NewEncoder(w).Encode(res)
}

func AddAttendeeController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := bson.ObjectIdHex(vars["id"])
	userID := bson.ObjectIdHex(vars["userID"])

	event, user := AddAttendeeToGoingList(eventID, userID)

	_ = json.NewEncoder(w).Encode(bson.M{"event": event, "user": user})
}

// AddAttendeeController will answer the JSON
// of the event with the given attendee added
func ChangeAttendeeStatusController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := bson.ObjectIdHex(vars["id"])
	userID := bson.ObjectIdHex(vars["userID"])
	status := vars["status"]

	if status == "going" {
		event, user := AddAttendeeToGoingList(eventID, userID)
		_ = json.NewEncoder(w).Encode(bson.M{"event": event, "user": user})
	} else if status == "maybe" {
		event, user := AddAttendeeToMaybeList(eventID, userID)
		_ = json.NewEncoder(w).Encode(bson.M{"event": event, "user": user})
	} else if status == "notgoing" {
		event, user := AddAttendeeToNotGoingList(eventID, userID)
		_ = json.NewEncoder(w).Encode(bson.M{"event": event, "user": user})
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "bad status"})
	}
}

// RemoveAttendeeController will answer the JSON
// of the event without the given attendee added
func RemoveAttendeeController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := bson.ObjectIdHex(vars["id"])
	userID := bson.ObjectIdHex(vars["userID"])

	RemoveAttendee(eventID, userID, "participants")
	RemoveAttendee(eventID, userID, "notgoing")
	event, user := RemoveAttendee(eventID, userID, "maybe")

	_ = json.NewEncoder(w).Encode(bson.M{"event": event, "user": user})
}

// CommentEventController will answer a JSON of the event
func CommentEventController(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Unable to read the request body"})
	}
	var comment Comment
	if err := json.Unmarshal([]byte(string(body)), &comment); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "wrong format"})
		return
	}

	comment.ID = bson.NewObjectId()
	comment.Date = time.Now()

	vars := mux.Vars(r)
	eventID := vars["id"]

	event := CommentEvent(bson.ObjectIdHex(eventID), comment)
	association := GetAssociation(event.Association)
	user := GetUser(comment.User)

	_ = json.NewEncoder(w).Encode(event)

	if !event.NoNotification {
		_ = SendAssociationEmailForCommentOnEvent(association.Email, event, comment, user)
	}

	for _, tag := range comment.Tags {
		go TriggerNotificationForUserFromEvent(comment.User, bson.ObjectIdHex(tag.User), event.ID, "@"+GetUser(comment.User).Username+" t'a taggé sur '"+event.Name+"'", comment, "eventTag")
	}
}

// UncommentEventController will answer a JSON of the event
func UncommentEventController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["id"]
	commentID := vars["commentID"]

	res := UncommentEvent(bson.ObjectIdHex(eventID), bson.ObjectIdHex(commentID))

	_ = json.NewEncoder(w).Encode(res)
}
