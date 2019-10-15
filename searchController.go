package insapp

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type Search struct {
	Terms string `json:"terms"`
}

func SearchUserController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var search Search
	decoder.Decode(&search)
	users := SearchUser(search.Terms)
	json.NewEncoder(w).Encode(bson.M{"users": users})
}

func SearchPostController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var search Search
	decoder.Decode(&search)
	posts := SearchPost(search.Terms)
	json.NewEncoder(w).Encode(bson.M{"posts": posts})
}

func SearchEventController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var search Search
	decoder.Decode(&search)
	events := SearchEvent(search.Terms)
	json.NewEncoder(w).Encode(bson.M{"events": events})
}

func SearchAssociationController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var search Search
	decoder.Decode(&search)
	associations := SearchAssociation(search.Terms)
	json.NewEncoder(w).Encode(bson.M{"associations": associations})
}

func SearchUniversalController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var search Search
	decoder.Decode(&search)

	users := SearchUser(search.Terms)
	posts := SearchPost(search.Terms)
	events := SearchEvent(search.Terms)
	associations := SearchAssociation(search.Terms)
	json.NewEncoder(w).Encode(bson.M{"associations": associations, "users": users, "posts": posts, "events": events})
}
