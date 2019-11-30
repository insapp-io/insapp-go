package insapp

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func GetAssociationUserController(w http.ResponseWriter, r *http.Request) {
	authCookie, err1 := r.Cookie("AuthToken")
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(bson.M{
			"error": "failed to retrieve current user",
		})

		return
	}

	authToken, err2 := parseAuthStringToken(authCookie.Value)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(bson.M{
			"error": "failed to retrieve current user",
		})

		return
	}

	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("association_user")

	var user AssociationUser
	err := db.FindId(
		authToken.Claims.(*TokenClaims).ID,
	).One(&user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(bson.M{
			"error": "failed to retrieve current user: not found",
		})

		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

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
	id := vars["id"]

	userID, err := GetUserFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(bson.M{"error": "could not get user ID"})
		return
	}

	ReportUser(bson.ObjectIdHex(id), userID)
	json.NewEncoder(w).Encode(bson.M{})
}
