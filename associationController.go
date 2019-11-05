package insapp

import (
	"encoding/json"
	"net/http"

	tauth "github.com/freehaha/token-auth"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// GetMyAssociationController will answer a JSON of the associations owned by the applicant master association
func GetMyAssociationController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	associationID := vars["id"]
	var res = GetMyAssociations(bson.ObjectIdHex(associationID))
	_ = json.NewEncoder(w).Encode(res)
}

// GetAssociationController will answer a JSON of the association
// linked to the given id in the URL
func GetAssociationController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	associationID := vars["id"]
	var res = GetAssociation(bson.ObjectIdHex(associationID))
	_ = json.NewEncoder(w).Encode(res)
}

// GetAllAssociationsController will answer a JSON of all associations
func GetAllAssociationsController(w http.ResponseWriter, r *http.Request) {
	var res = GetAllAssociations()
	_ = json.NewEncoder(w).Encode(res)
}

// AddAssociationController will answer a JSON of the
// brand new created association (from the JSON Body)
// Should be protected
func AddAssociationController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var association Association
	_ = decoder.Decode(&association)

	isValidMail := VerifyEmail(association.Email)
	if !isValidMail {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "email already used"})
		return
	}

	res := AddAssociation(association)
	password := GeneratePassword()
	token := tauth.Get(r)
	id := bson.ObjectIdHex(token.Claims("id").(string))

	var user AssociationUser
	user.Association = res.ID
	user.Username = res.Email
	user.Master = false
	user.Owner = id
	user.Password = GetMD5Hash(password)
	AddAssociationUser(user)
	_ = SendAssociationEmailSubscription(user.Username, password)
	_ = json.NewEncoder(w).Encode(res)
}

// UpdateAssociationController will answer the JSON of the
// modified association (from the JSON Body)
// Should be protected
func UpdateAssociationController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var association Association
	_ = decoder.Decode(&association)
	vars := mux.Vars(r)
	associationID := vars["id"]

	res := UpdateAssociation(bson.ObjectIdHex(associationID), association)

	_ = json.NewEncoder(w).Encode(res)
}

// DeleteAssociationController will answer a JSON of an
// empty association if the deletion has succeed
// Should be protected
func DeleteAssociationController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	associationID := vars["id"]

	res := DeleteAssociation(bson.ObjectIdHex(associationID))

	_ = json.NewEncoder(w).Encode(res)
}

// VerifyEmail return true if email is not already used
func VerifyEmail(email string) bool {
	association := GetAssociationFromEmail(email)
	return association.Email == ""
}
