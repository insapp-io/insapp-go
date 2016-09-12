package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/freehaha/token-auth/memory"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Login struct {
	Username string
	Password string
}

type AssociationUser struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Username    string        `json:"username"`
	Association bson.ObjectId `json:"association" bson:"association"`
	Password    string        `json:"password"`
	Master      bool          `json:"master"`
	Owner       bson.ObjectId `json:"owner" bson:"owner,omitempty"`
}

func LogAssociationController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var login Login
	decoder.Decode(&login)
	fmt.Println(login)
	auth, master, err := checkLoginForAssociation(login)
	if err == nil {
		sessionToken := logAssociation(auth, master)
		json.NewEncoder(w).Encode(bson.M{"token": sessionToken.Token, "master": master, "associationID": auth})
	} else {
		json.NewEncoder(w).Encode(bson.M{"error": "Failed to authentificate"})
	}
}

func LogUserController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var login Login
	decoder.Decode(&login)
	auth, err := checkLoginForUser(login)
	if err == nil {
		sessionToken := logUser(auth)
		json.NewEncoder(w).Encode(bson.M{"token": sessionToken.Token, "userID": auth})
	} else {
		json.NewEncoder(w).Encode(bson.M{"error": "Failed to authentificate"})
	}
}

func checkLoginForAssociation(login Login) (bson.ObjectId, bool, error) {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association_user")
	var result []AssociationUser
	db.Find(bson.M{"username": login.Username, "password": GetMD5Hash(login.Password)}).All(&result)
	if len(result) > 0 {
		return result[0].Association, result[0].Master, nil
	}
	return bson.ObjectId(""), false, errors.New("Failed to authentificate")
}

func checkLoginForUser(login Login) (bson.ObjectId, error) {

	//CHECK USER ACCESS
	//return bson.ObjectId(""), errors.New("Failed to authentificate")

	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("user")
	var result []AssociationUser
	db.Find(bson.M{"username": login.Username).All(&result)
	if len(result) > 0 {
		return result[0].ID, nil
	}
	return bson.ObjectId(""), errors.New("No User Found")
}

func logAssociation(id bson.ObjectId, master bool) *memstore.MemoryToken {
	if master {
		memStoreUser.NewToken(string(id))
		return memStoreSuperUser.NewToken(string(id))
	}
	return memStoreUser.NewToken(string(id))
}

func logUser(id bson.ObjectId) *memstore.MemoryToken {
	return memStoreUser.NewToken(string(id))
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
