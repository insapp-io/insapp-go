package insapp

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// TokenJTI models JTI keeping track of tokens.
type TokenJTI struct {
	ID  bson.ObjectId `bson:"_id,omitempty"`
	JTI string        `json:"jti"`
}

// Login is the data provided by the user to authenticate.
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AssociationUser is the data provided by an association to authenticate.
type AssociationUser struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Username    string        `json:"username"`
	Association bson.ObjectId `json:"association" bson:"association"`
	Password    string        `json:"password"`
	Master      bool          `json:"master"`
	Owner       bson.ObjectId `json:"owner" bson:"owner,omitempty"`
}

// AuthMiddleware makes sure the user is authenticated before handling the request.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthCookie, authErr := r.Cookie("AuthToken")

		// Unauthorized attempt: no auth cookie
		if authErr == http.ErrNoCookie {
			nullifyTokenCookies(&w, r)
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Internal error
		if authErr != nil {
			nullifyTokenCookies(&w, r)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		RefreshCookie, refreshErr := r.Cookie("RefreshToken")

		// Unauthorized attempt: no refresh cookie
		if refreshErr == http.ErrNoCookie {
			nullifyTokenCookies(&w, r)
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Internal error
		if refreshErr != nil {
			nullifyTokenCookies(&w, r)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		// Check the JWT for validity
		authToken, refreshToken, err := CheckAndRefreshTokens(AuthCookie.Value, RefreshCookie.Value)
		if err != nil {
			// Unauthorized attempt: JWT is not valid
			if err.Error() == "Unauthorized" {
				nullifyTokenCookies(&w, r)
				http.Error(w, http.StatusText(401), 401)
				return
			}

			nullifyTokenCookies(&w, r)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		setAuthAndRefreshCookies(&w, authToken, refreshToken)

		next.ServeHTTP(w, r)
	})
}

// LogInUserController logs the user using CAS.
// If the credentials are correct, a JWT access token is generated.
func LogInUserController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	ticket := mux.Vars(r)["ticket"]
	username, err := isTicketValid(ticket)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{
			"error": err,
		})
	}

	var login Login
	decoder.Decode(&login)
	login.Username = strings.ToLower(username)

	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("user")

	count, _ := db.Find(bson.M{
		"username": login.Username,
	}).Count()

	var user User
	if count == 0 {
		user = AddUser(NewUser(login.Username))
	} else {
		db.Find(bson.M{
			"username": login.Username,
		}).One(&user)
	}

	authToken, refreshToken, err := CreateNewTokens(user.Username, "user")
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(bson.M{
			"error": err,
		})
	} else {
		// Set the cookies to these newly created tokens
		setAuthAndRefreshCookies(&w, authToken, refreshToken)
		w.WriteHeader(http.StatusOK)
	}
}

// isTicketValid checks the validity of the given ticket with the CAS
func isTicketValid(ticket string) (string, error) {
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

func setAuthAndRefreshCookies(w *http.ResponseWriter, authToken string, refreshToken string) {
	http.SetCookie(*w, &http.Cookie{
		Name:     "AuthToken",
		Value:    authToken,
		HttpOnly: true,
	})

	http.SetCookie(*w, &http.Cookie{
		Name:     "RefreshToken",
		Value:    refreshToken,
		HttpOnly: true,
	})
}

func nullifyTokenCookies(w *http.ResponseWriter, r *http.Request) {
	authCookie := http.Cookie{
		Name:     "AuthToken",
		Value:    "",
		Expires:  time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(*w, &authCookie)

	refreshCookie := http.Cookie{
		Name:     "RefreshToken",
		Value:    "",
		Expires:  time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(*w, &refreshCookie)

	// If present, revoke the refresh cookie from the database
	RefreshCookie, refreshErr := r.Cookie("RefreshToken")
	if refreshErr == http.ErrNoCookie {
		// Do nothing, there is no refresh cookie present
		return
	}

	if refreshErr != nil {
		http.Error(*w, http.StatusText(500), 500)
	}

	RevokeRefreshToken(RefreshCookie.Value)
}

func checkLoginForAssociation(login Login) (bson.ObjectId, bool, error) {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("association_user")

	var result []AssociationUser
	db.Find(bson.M{
		"username": login.Username,
		"password": GetMD5Hash(login.Password),
	}).All(&result)

	if len(result) > 0 {
		return result[0].Association, result[0].Master, nil
	}

	return bson.ObjectId(""), false, errors.New("failed to authenticate")
}

func CheckRefreshToken(jti string) bool {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("tokens")

	count, err := db.Find(bson.M{
		"jti": jti,
	}).Count()

	return err != nil && count > 0
}

func StoreRefreshToken() TokenJTI {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("tokens")

	jti, _ := GenerateRandomString(32)
	for CheckRefreshToken(jti) {
		jti, _ = GenerateRandomString(32)
	}

	var token TokenJTI
	token.JTI = jti
	db.Insert(token)

	return token
}

func DeleteRefreshToken(jti string) {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("tokens")

	db.Remove(bson.M{"jti": jti})
}
