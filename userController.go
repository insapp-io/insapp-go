package insapp

import (
	"encoding/json"
	"net/http"

	tauth "github.com/freehaha/token-auth"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// GetUserController will answer a JSON of the user
// linked to the given id in the URL
func GetUserController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	var res = GetUser(bson.ObjectIdHex(userID))
	json.NewEncoder(w).Encode(res)
}

func GetAllUserController(w http.ResponseWriter, r *http.Request) {
	var res = GetAllUser()
	json.NewEncoder(w).Encode(res)
}

// AddUserController will answer a JSON of the
// brand new created user (from the JSON Body)
func AddUserController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user User
	decoder.Decode(&user)
	res := AddUser(&user)
	json.NewEncoder(w).Encode(res)
}

// UpdateUserController will answer the JSON of the
// modified user (from the JSON Body)
// Should be protected
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
// empty user if the deletion succeeded.
// Should be protected
func DeleteUserController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	DeleteTokenCookies(&w, r)

	user := GetUser(bson.ObjectIdHex(userID))
	res := DeleteUser(user)

	json.NewEncoder(w).Encode(res)
}

func ReportUserController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	token := tauth.Get(r)
	reporterID := token.Claims("id").(string)
	ReportUser(bson.ObjectIdHex(userID), bson.ObjectIdHex(reporterID))
	json.NewEncoder(w).Encode(bson.M{})
}

func GetUserFromRequest(r *http.Request) string {
	token := tauth.Get(r)
	id := token.Claims("id").(string)
	return id
}

func Contains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
