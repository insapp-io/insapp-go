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
func AuthMiddleware(next http.HandlerFunc, role string) http.HandlerFunc {
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
		authToken, refreshToken, err := CheckAndRefreshStringTokens(authCookie.Value, refreshCookie.Value, role)
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

// LoginUserController logs a user in using CAS.
// If the ticket is valid, auth and refresh tokens are generated.
func LoginUserController(w http.ResponseWriter, r *http.Request) {
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

// LoginAssociationController logs an association in.
func LoginAssociationController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var login AssociationLogin
	decoder.Decode(&login)

	user, err := checkLoginForAssociation(login)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{
			"error": fmt.Sprintf("failed to authenticate: %s", err.Error()),
		})

		return
	}

	var role string
	if user.Master {
		role = "admin"
	} else {
		role = "association"
	}
	authToken, refreshToken := CreateNewTokens(user.ID, role)

	// Set the cookies to these newly created tokens
	setAuthAndRefreshCookies(&w, authToken, refreshToken)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(user)
}

// LogoutUserController logs a user out.
func LogoutUserController(w http.ResponseWriter, r *http.Request) {
	DeleteTokenCookies(&w, r)

	w.WriteHeader(http.StatusOK)
}

// LogoutAssociationController logs an association out.
func LogoutAssociationController(w http.ResponseWriter, r *http.Request) {
	DeleteTokenCookies(&w, r)

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

func checkLoginForAssociation(login AssociationLogin) (*AssociationUser, error) {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("association_user")

	var result AssociationUser
	err1 := db.Find(bson.M{
		"username": login.Username,
	}).One(&result)

	if err1 != nil {
		return nil, errors.New("unknown user")
	}

	err2 := db.Find(bson.M{
		"username": login.Username,
		"password": GetMD5Hash(login.Password),
	}).One(&result)

	if err2 != nil {
		return nil, errors.New("wrong password")
	}

	return &result, nil
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

	// The expiration times are set to the refresh token expiration time

	if config.Environment == "local" {
		http.SetCookie(*w, &http.Cookie{
			Name:     "AuthToken",
			Value:    authStringToken,
			Domain:   config.Domain,
			Path:     "/",
			Secure:   false,
			Expires:  time.Unix(refreshToken.Claims.(TokenClaims).ExpiresAt, 0),
			HttpOnly: true,
		})

		http.SetCookie(*w, &http.Cookie{
			Name:     "RefreshToken",
			Value:    refreshStringToken,
			Domain:   config.Domain,
			Path:     "/",
			Secure:   false,
			Expires:  time.Unix(refreshToken.Claims.(TokenClaims).ExpiresAt, 0),
			HttpOnly: true,
		})
	} else {
		http.SetCookie(*w, &http.Cookie{
			Name:     "AuthToken",
			Value:    authStringToken,
			Domain:   config.Domain,
			Path:     "/",
			Secure:   true,
			Expires:  time.Unix(refreshToken.Claims.(TokenClaims).ExpiresAt, 0),
			HttpOnly: true,
		})

		http.SetCookie(*w, &http.Cookie{
			Name:     "RefreshToken",
			Value:    refreshStringToken,
			Domain:   config.Domain,
			Path:     "/",
			Secure:   true,
			Expires:  time.Unix(refreshToken.Claims.(TokenClaims).ExpiresAt, 0),
			HttpOnly: true,
		})
	}

	return nil
}

func nullifyTokenCookies(w *http.ResponseWriter, r *http.Request) {
	if config.Environment == "local" {
		http.SetCookie(*w, &http.Cookie{
			Name:     "AuthToken",
			Value:    "",
			Domain:   config.Domain,
			Path:     "/",
			Secure:   false,
			Expires:  time.Now().Add(-1000 * time.Hour),
			HttpOnly: true,
		})

		http.SetCookie(*w, &http.Cookie{
			Name:     "RefreshToken",
			Value:    "",
			Domain:   config.Domain,
			Path:     "/",
			Secure:   false,
			Expires:  time.Now().Add(-1000 * time.Hour),
			HttpOnly: true,
		})
	} else {
		http.SetCookie(*w, &http.Cookie{
			Name:     "AuthToken",
			Value:    "",
			Domain:   config.Domain,
			Path:     "/",
			Secure:   true,
			Expires:  time.Now().Add(-1000 * time.Hour),
			HttpOnly: true,
		})

		http.SetCookie(*w, &http.Cookie{
			Name:     "RefreshToken",
			Value:    "",
			Domain:   config.Domain,
			Path:     "/",
			Secure:   true,
			Expires:  time.Now().Add(-1000 * time.Hour),
			HttpOnly: true,
		})
	}

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
