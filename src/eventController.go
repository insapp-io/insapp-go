package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// GetEventController will answer a JSON of the event
// from the given "id" in the URL. (cf Routes in routes.go)
func GetEventController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	associationID := vars["id"]
	var res = GetEvent(bson.ObjectIdHex(associationID))
	json.NewEncoder(w).Encode(res)
}

// GetFutureEventsController will answer a JSON
// containing all future events from "NOW"
func GetFutureEventsController(w http.ResponseWriter, r *http.Request) {
	userId := GetUserFromRequest(r)
	user := GetUser(bson.ObjectIdHex(userId))
	os := GetNotificationUserForUser(bson.ObjectIdHex(userId)).Os
	events := GetFutureEvents()
	res := Events{}
	if user.ID != "" {
		for _, event := range events {
			if Contains(strings.ToUpper(user.Promotion), event.Promotions) && (Contains(os, event.Plateforms) || os == "") || len(event.Promotions) == 0 || len(event.Plateforms) == 0 {
				res = append(res, event)
			}
		}
	} else {
		res = events
	}
	json.NewEncoder(w).Encode(res)
}

// AddEventController will answer the JSON
// of the brand new created event from the JSON body
func AddEventController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var event Event
	decoder.Decode(&event)

	isValid := VerifyAssociationRequest(r, event.Association)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "protected content"})
		return
	}

	res := AddEvent(event)
	asso := GetAssociation(event.Association)
	json.NewEncoder(w).Encode(res)
	go TriggerNotificationForEvent(event, asso.ID, res.ID, "@"+strings.ToLower(asso.Name)+" t'invite à "+res.Name+" 📅")
}

// UpdateEventController will answer the JSON
// of the brand new modified event from the JSON body
func UpdateEventController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var event Event
	decoder.Decode(&event)
	vars := mux.Vars(r)
	eventID := vars["id"]

	isValid := VerifyAssociationRequest(r, event.Association)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "protected content"})
		return
	}

	res := UpdateEvent(bson.ObjectIdHex(eventID), event)
	json.NewEncoder(w).Encode(res)
}

// DeleteEventController will answer an empty JSON
// if the deletation has succeed
func DeleteEventController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := GetEvent(bson.ObjectIdHex(vars["id"]))

	isValid := VerifyAssociationRequest(r, event.Association)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "protected content"})
		return
	}

	res := DeleteEvent(event)
	json.NewEncoder(w).Encode(res)
}

func AddParticipantController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := bson.ObjectIdHex(vars["id"])
	userID := bson.ObjectIdHex(vars["userID"])
	isValid := VerifyUserRequest(r, userID)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "protected content"})
		return
	}
	event, user := AddParticipantToGoingList(eventID, userID)
	json.NewEncoder(w).Encode(bson.M{"event": event, "user": user})
}

// AddParticipantController will answer the JSON
// of the event with the given attendee added
func ChangeAttendeeStatusController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := bson.ObjectIdHex(vars["id"])
	userID := bson.ObjectIdHex(vars["userID"])
	status := vars["status"]
	isValid := VerifyUserRequest(r, userID)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "protected content"})
		return
	}
	if status == "going" {
		event, user := AddParticipantToGoingList(eventID, userID)
		json.NewEncoder(w).Encode(bson.M{"event": event, "user": user})
	} else if status == "maybe" {
		event, user := AddParticipantToMaybeList(eventID, userID)
		json.NewEncoder(w).Encode(bson.M{"event": event, "user": user})
	} else if status == "notgoing" {
		event, user := AddParticipantToNotGoingList(eventID, userID)
		json.NewEncoder(w).Encode(bson.M{"event": event, "user": user})
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(bson.M{"error": "bad status"})
	}
}

// RemoveParticipantController will answer the JSON
// of the event without the given partipant added
func RemoveParticipantController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := bson.ObjectIdHex(vars["id"])
	userID := bson.ObjectIdHex(vars["userID"])
	isValid := VerifyUserRequest(r, userID)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "protected content"})
		return
	}
	RemoveParticipant(eventID, userID, "participants")
	RemoveParticipant(eventID, userID, "notgoing")
	event, user := RemoveParticipant(eventID, userID, "maybe")
	json.NewEncoder(w).Encode(bson.M{"event": event, "user": user})
}

// CommentPostController will answer a JSON of the post
func CommentEventController(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bson.M{"error": "Unable to read the request body"})
	}
	var comment Comment
	if err := json.Unmarshal([]byte(string(body)), &comment); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bson.M{"error": "wrong format"})
		return
	}

	isValid := VerifyUserRequest(r, comment.User)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "protected content"})
		return
	}

	comment.ID = bson.NewObjectId()
	comment.Date = time.Now()

	vars := mux.Vars(r)
	eventID := vars["id"]

	event := CommentEvent(bson.ObjectIdHex(eventID), comment)
	association := GetAssociation(event.Association)
	user := GetUser(comment.User)

	json.NewEncoder(w).Encode(event)

	if !event.NoNotification {
		SendAssociationEmailForCommentOnEvent(association.Email, event, comment, user)
	}

	for _, tag := range comment.Tags {
		go TriggerNotificationForUser(comment.User, bson.ObjectIdHex(tag.User), event.ID, "@"+GetUser(comment.User).Username+" t'a taggé sur \""+event.Name+"\"", comment, "eventTag")
	}
}

// UncommentPostController will answer a JSON of the post
func UncommentEventController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["id"]
	commentID := vars["commentID"]
	comment, err := GetCommentForEvent(bson.ObjectIdHex(eventID), bson.ObjectIdHex(commentID))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(bson.M{"error": "content not found"})
		return
	}
	event := GetEvent(bson.ObjectIdHex(eventID))
	isUserValid := VerifyUserRequest(r, comment.User)
	isAssociationValid := VerifyAssociationRequest(r, event.Association)
	if !isUserValid && !isAssociationValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "protected content"})
		return
	}
	res := UncommentEvent(bson.ObjectIdHex(eventID), bson.ObjectIdHex(commentID))
	json.NewEncoder(w).Encode(res)
}
