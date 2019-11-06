package insapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// TokenJTI models JTI keeping track of tokens.
type TokenJTI struct {
	ID  bson.ObjectId `bson:"_id,omitempty"`
	JTI string        `json:"jti"`
}

// AssociationLogin is the data provided by an association to authenticate.
type AssociationLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthMiddleware makes sure the user is authenticated before handling the request.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(requestDump))
		}

		authCookie, authErr := r.Cookie("AuthToken")

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

		refreshCookie, refreshErr := r.Cookie("RefreshToken")

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
		authToken, refreshToken, err := CheckAndRefreshStringTokens(authCookie.Value, refreshCookie.Value, "user")
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

// LogUserController logs a user in using CAS.
// If the ticket is valid, auth and refresh tokens are generated.
func LogUserController(w http.ResponseWriter, r *http.Request) {
	ticket := mux.Vars(r)["ticket"]

	username, err := isTicketValid(ticket)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{
			"error": "failed to authenticate",
		})
	}

	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("user")

	count, _ := db.Find(bson.M{
		"username": username,
	}).Count()

	var user User
	if count == 0 {
		user = AddUser(NewUser(username))
	} else {
		db.Find(bson.M{
			"username": username,
		}).One(&user)
	}

	authToken, refreshToken := CreateNewTokens(user.ID, "user")

	// Set the cookies to these newly created tokens
	setAuthAndRefreshCookies(&w, authToken, refreshToken)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(user)
}

// LogAssociationController logs an association in.
func LogAssociationController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var login AssociationLogin
	decoder.Decode(&login)

	id, _, err := checkLoginForAssociation(login)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)

		json.NewEncoder(w).Encode(bson.M{
			"error": "failed to authenticate",
		})
		return
	}

	authToken, refreshToken := CreateNewTokens(id, "association")

	// Set the cookies to these newly created tokens
	setAuthAndRefreshCookies(&w, authToken, refreshToken)
	w.WriteHeader(http.StatusOK)
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

	return strings.ToLower(username), nil
}

func checkLoginForAssociation(login AssociationLogin) (bson.ObjectId, bool, error) {
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

func setAuthAndRefreshCookies(w *http.ResponseWriter, authToken *jwt.Token, refreshToken *jwt.Token) error {
	authStringToken, err1 := authToken.SignedString(signKey)
	if err1 != nil {
		return err1
	}

	refreshStringToken, err2 := refreshToken.SignedString(signKey)
	if err2 != nil {
		return err2
	}

	http.SetCookie(*w, &http.Cookie{
		Name:     "AuthToken",
		Value:    authStringToken,
		Domain:   config.Domain,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	http.SetCookie(*w, &http.Cookie{
		Name:     "RefreshToken",
		Value:    refreshStringToken,
		Domain:   config.Domain,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	return nil
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

	RevokeRefreshStringToken(RefreshCookie.Value)
}

// DeleteTokenCookies deletes the cookies
func DeleteTokenCookies(w *http.ResponseWriter, r *http.Request) {
	_, authErr := r.Cookie("AuthToken")

	if authErr == http.ErrNoCookie {
		nullifyTokenCookies(w, r)
		return
	}

	if authErr != nil {
		nullifyTokenCookies(w, r)
		http.Error(*w, http.StatusText(500), 500)
		return
	}

	// Remove this user's ability to make requests
	nullifyTokenCookies(w, r)
}
