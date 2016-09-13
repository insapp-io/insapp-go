package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
  "os/exec"

	"github.com/freehaha/token-auth/memory"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Login struct {
	Username string
	Password string
}

type Credentials struct {
	ID 				bson.ObjectId		`bson:"_id,omitempty"`
	Username 	string					`json:"username"`
	AuthToken 		string			`json:"authtoken"`
	User 			bson.ObjectId		`json:"user" bson:"user"`
	Token 		*memstore.MemoryToken					`json:"token"`
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
	var credentials Credentials
	decoder.Decode(&credentials)
	cred, err := checkLoginForUser(credentials)
	if err == nil {
		sessionToken := logUser(cred.User)
		cred.Token = sessionToken
		json.NewEncoder(w).Encode(cred)
	} else {
		json.NewEncoder(w).Encode(bson.M{"error": err})
	}
}

func SignInUserController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var login Login
	decoder.Decode(&login)
	isCASValid, err := verifyUserWithCAS(login)
	if isCASValid {
		user := AddUser(User{Username:login.Username})
		token := generateAuthToken()
		credentials := Credentials{AuthToken: token, User: user.ID, Username: user.Username}
		result := addCredentials(credentials)
		json.NewEncoder(w).Encode(result)
	} else {
		json.NewEncoder(w).Encode(bson.M{"error": err})
	}
}

func generateAuthToken() (string){
	out, _ := exec.Command("uuidgen").Output()
	return string(out)
}

func addCredentials(credentials Credentials) (Credentials){
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("credentials")
	db.Insert(credentials)
	var result Credentials
	db.Find(bson.M{"username": credentials.Username}).One(&result)
	return result
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

func verifyUserWithCAS(login Login) (bool, error){
	return true, nil
}

func checkLoginForUser(credentials Credentials) (Credentials, error) {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("credentials")
	var result []Credentials
	db.Find(bson.M{"username": credentials.Username, "authtoken": credentials.AuthToken}).All(&result)
	if len(result) > 0 {
		return result[0], nil
	}
	return Credentials{}, errors.New("No User Found")
}

func logAssociation(id bson.ObjectId, master bool) *memstore.MemoryToken {
	if master {
		memStoreUser.NewToken(string(id))
		memStoreAssociationUser.NewToken(string(id))
		return memStoreSuperUser.NewToken(string(id))
	}
	memStoreUser.NewToken(string(id))
	return memStoreAssociationUser.NewToken(string(id))
}

func logUser(id bson.ObjectId) *memstore.MemoryToken {
	return memStoreUser.NewToken(string(id))
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
