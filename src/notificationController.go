package main

import (
	"encoding/json"
	"net/http"
	"gopkg.in/mgo.v2/bson"
)


// AddUserController will answer a JSON of the
// brand new created user (from the JSON Body)
func UpdateNotificationUserController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user NotificationUser
	decoder.Decode(&user)
	CreateOrUpdateNotificationUser(user)
	json.NewEncoder(w).Encode(bson.M{"status": "ok"})
}
