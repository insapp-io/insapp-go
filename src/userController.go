package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

// GetUserController will answer a JSON of the user
// linked to the given id in the URL
func GetUserController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	var res = GetUser(bson.ObjectIdHex(userID))
	json.NewEncoder(w).Encode(res)
}

// AddUserController will answer a JSON of the
// brand new created user (from the JSON Body)
func AddUserController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user User
	decoder.Decode(&user)
	res := AddUser(user)
	json.NewEncoder(w).Encode(res)
}

// UpdateUserController will answer the JSON of the
// modified user (from the JSON Body)
func UpdateUserController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user User
	decoder.Decode(&user)
	vars := mux.Vars(r)
	userID := vars["id"]
	res := UpdateUser(bson.ObjectIdHex(userID), user)
	json.NewEncoder(w).Encode(res)
}

// DeleteUserController will answer a JSON of an
// empty user if the deletation has succeed
func DeleteUserController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := GetUser(bson.ObjectIdHex(vars["id"]))
	res := DeleteUser(user)
	json.NewEncoder(w).Encode(res)
}

func SearchUserController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	users := SearchUser(vars["username"])
	json.NewEncoder(w).Encode(bson.M{"users": users})
}

// AddImageUserController will set the image of the user and return the user
func AddImageUserController(w http.ResponseWriter, r *http.Request) {
	fileName := UploadImage(r)
	if fileName == "error" {
		w.Header().Set("status", "400")
		fmt.Fprintln(w, "{}")
	} else {
		vars := mux.Vars(r)
		res := SetImageUser(bson.ObjectIdHex(vars["id"]), fileName)
		json.NewEncoder(w).Encode(res)
	}
}
