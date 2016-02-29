package main

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

func GetAssociationController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assocationID := vars["id"]
	var res = GetAssociation(bson.ObjectIdHex(assocationID))
	json.NewEncoder(w).Encode(res)
}

func GetAllAssociationsController(w http.ResponseWriter, r *http.Request) {
	var res = GetAllAssociation()
	json.NewEncoder(w).Encode(res)
}

func AddAssociationController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var association Association
	decoder.Decode(&association)
	res := AddAssociation(association)
	json.NewEncoder(w).Encode(res)
}

func UpdateAssociationController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var association Association
	decoder.Decode(&association)
	vars := mux.Vars(r)
	assocationID := vars["id"]
	res := UpdateAssociation(bson.ObjectIdHex(assocationID), association)
	json.NewEncoder(w).Encode(res)
}

func DeleteAssociationController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res := DeleteAssociation(bson.ObjectIdHex(vars["id"]))
	json.NewEncoder(w).Encode(res)
}
