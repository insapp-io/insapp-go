package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/freehaha/token-auth/memory"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Device   string `json:"device"`
}

type Credentials struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Username  string        `json:"username"`
	AuthToken string        `json:"authtoken"`
	User      bson.ObjectId `json:"user" bson:"user"`
	Device    string        `json:"device"`
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
	auth, master, err := checkLoginForAssociation(login)

	if err == nil {
		sessionToken := logAssociation(auth, master)
		json.NewEncoder(w).Encode(bson.M{"token": sessionToken.Token, "master": master, "associationID": auth})
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(bson.M{"error": "failed to authenticate"})
	}
}

func LogUserController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var credentials Credentials
	decoder.Decode(&credentials)
	cred, err := checkLoginForUser(credentials)

	if err == nil {
		sessionToken := logUser(cred.User)
		user := GetUser(cred.User)
		json.NewEncoder(w).Encode(bson.M{"credentials": credentials, "sessionToken": sessionToken, "user": user})
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(bson.M{"error": err})
	}
}

func SignInUserController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var login Login
	decoder.Decode(&login)

	vars := mux.Vars(r)
	ticket := vars["ticket"]

	// w.WriteHeader(http.StatusForbidden)
	// json.NewEncoder(w).Encode(bson.M{"error": "De manière temporaire, les inscriptions sont désactivées. Réessaye Lundi 😊" })
	// return

	username, err := verifyTicket(ticket)
	login.Username = username
	login.Username = strings.ToLower(login.Username)

	if err == nil && len(login.Username) > 0 && len(login.Device) > 0 {
		session := GetMongoSession()
		defer session.Close()
		db := session.DB("insapp").C("user")

		count, _ := db.Find(bson.M{"username": login.Username}).Count()
		var user User
		if count == 0 {
			user = AddUser(User{Name: "", Username: login.Username, Description: "", Email: "", EmailPublic: false, Promotion: "", Events: []bson.ObjectId{}, PostsLiked: []bson.ObjectId{}})
		} else {
			db.Find(bson.M{"username": login.Username}).One(&user)
		}

		token := generateAuthToken()
		credentials := Credentials{AuthToken: token, User: user.ID, Username: user.Username, Device: login.Device}
		result := addCredentials(credentials)
		json.NewEncoder(w).Encode(result)
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(bson.M{"error": err})
	}
}

func generateAuthToken() string {
	out, _ := exec.Command("uuidgen").Output()
	return strings.TrimSpace(string(out))
}

func DeleteCredentialsForUser(id bson.ObjectId) {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("credentials")

	db.Remove(bson.M{"user": id})
}

func addCredentials(credentials Credentials) Credentials {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("credentials")

	var cred Credentials
	db.Find(bson.M{"username": credentials.Username}).One(&cred)
	db.RemoveId(cred.ID)
	db.Insert(credentials)

	var result Credentials
	db.Find(bson.M{"username": credentials.Username}).One(&result)

	return result
}

func checkLoginForAssociation(login Login) (bson.ObjectId, bool, error) {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("association_user")

	var result []AssociationUser
	db.Find(bson.M{"username": login.Username, "password": GetMD5Hash(login.Password)}).All(&result)
	if len(result) > 0 {
		return result[0].Association, result[0].Master, nil
	}

	return bson.ObjectId(""), false, errors.New("failed to authenticate")
}

func verifyTicket(ticket string) (string, error) {
	response, err := http.Get("https://cas.insa-rennes.fr/cas/serviceValidate?service=https%3A%2F%2Finsapp.fr%2F&ticket=" + ticket)
	if err != nil {
		return "", errors.New("unable to verify identity")
	}
	defer response.Body.Close()

	htmlData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("unable to verify identity")
	}
	xml := string(htmlData)
	if !strings.Contains(xml, "<cas:authenticationSuccess>") && !strings.Contains(xml, "<cas:user>") {
		return "", errors.New("unable to verify identity")
	}

	username := strings.Split(xml, "<cas:user>")[1]
	username = strings.Split(username, "</cas:user>")[0]

	if !(len(username) > 2) {
		return "", errors.New("unable to verify identity")
	}
	return username, nil
}

func checkLoginForUser(credentials Credentials) (Credentials, error) {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("credentials")

	var result []Credentials
	db.Find(bson.M{"username": credentials.Username, "authtoken": credentials.AuthToken}).All(&result)
	if len(result) > 0 {
		return result[0], nil
	}

	return Credentials{}, errors.New("wrong credentials")
}

func logAssociation(id bson.ObjectId, master bool) *memstore.MemoryToken {
	if master {
		memStoreUser.NewToken(id.Hex())
		memStoreAssociationUser.NewToken(id.Hex())
		return memStoreSuperUser.NewToken(id.Hex())
	}
	memStoreUser.NewToken(id.Hex())
	return memStoreAssociationUser.NewToken(id.Hex())
}

func logUser(id bson.ObjectId) *memstore.MemoryToken {
	return memStoreUser.NewToken(id.Hex())
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
