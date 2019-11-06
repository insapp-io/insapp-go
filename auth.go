package insapp

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

// TokenClaims is the JWT encoding format.
type TokenClaims struct {
	ID   bson.ObjectId `json:"id"`
	Role string        `json:"role"`
	jwt.StandardClaims
}

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

const (
	refreshTokenValidTime = time.Hour * 672
	authTokenValidTime    = time.Hour * 24
)

// InitJWT reads the key files before starting http handlers
func InitJWT() error {
	signBytes, err := ioutil.ReadFile(config.PrivateKeyPath)
	if err != nil {
		return err
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return err
	}

	verifyBytes, err := ioutil.ReadFile(config.PublicKeyPath)
	if err != nil {
		return err
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return err
	}

	return nil
}

// CreateNewTokens creates auth and refresh tokens.
func CreateNewTokens(ID bson.ObjectId, role string) (string, string, error) {
	// Generate the auth token
	authTokenString, err := createAuthTokenString(ID, role)
	if err != nil {
		return "", "", err
	}

	// Generate the refresh token
	refreshTokenString, err := createRefreshTokenString(ID, role)
	if err != nil {
		return "", "", err
	}

	return authTokenString, refreshTokenString, nil
}

// CheckAndRefreshTokens renews the auth token, if needed.
func CheckAndRefreshTokens(authTokenString string, refreshTokenString string, role string) (string, string, error) {
	var newAuthTokenString string
	var newRefreshTokenString string

	// Check that it matches with the auth token claims
	authToken, err := jwt.ParseWithClaims(authTokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err != nil {
		return "", "", err
	}

	// The auth token is still valid
	if _, ok := authToken.Claims.(*TokenClaims); ok && authToken.Valid {
		// Check the role
		if authToken.Claims.(*TokenClaims).Role != role {
			return "", "", errors.New("Unauthorized")
		}

		// Update the expiration time of refresh token
		newRefreshTokenString, err = updateRefreshTokenExpiration(refreshTokenString)

		return authTokenString, newRefreshTokenString, nil
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		// The auth token has expired
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			newAuthTokenString, err = updateAuthTokenString(authTokenString, refreshTokenString)
			if err != nil {
				return "", "", err
			}

			// Update the expiration time of refresh token string
			newRefreshTokenString, err = updateRefreshTokenExpiration(refreshTokenString)
			if err != nil {
				return "", "", err
			}

			return newAuthTokenString, newRefreshTokenString, nil
		}
	}

	return "", "", err
}

// RevokeRefreshToken deletes the given token from the database, if valid.
func RevokeRefreshToken(refreshTokenString string) error {
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err != nil {
		return errors.New("Could not parse refresh token with claims")
	}

	refreshTokenClaims, ok := refreshToken.Claims.(*TokenClaims)
	if !ok {
		return errors.New("Could not read refresh token claims")
	}

	deleteRefreshToken(refreshTokenClaims.StandardClaims.Id)

	return nil
}

func GetUserFromRequest(r *http.Request) bson.ObjectId {
	authCookie, _ := r.Cookie("AuthToken")

	// Check that it matches with the auth token claims
	authToken, _ := jwt.ParseWithClaims(authCookie.Value, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	return authToken.Claims.(*TokenClaims).ID
}

// createAuthTokenString creates an auth token
func createAuthTokenString(id bson.ObjectId, role string) (string, error) {
	authTokenExpiration := time.Now().Add(authTokenValidTime).Unix()

	authClaims := TokenClaims{
		ID:   id,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: authTokenExpiration,
		},
	}

	// Create a signer for rsa 256
	authJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), authClaims)

	// Generate the auth token string
	return authJwt.SignedString(signKey)
}

func updateAuthTokenString(authTokenString string, refreshTokenString string) (string, error) {
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	refreshTokenClaims, ok := refreshToken.Claims.(*TokenClaims)
	if !ok {
		return "", err
	}

	// Check that the refresh token has not been revoked
	if checkRefreshToken(refreshTokenClaims.StandardClaims.Id) {
		// Has the refresh token expired?
		if refreshToken.Valid {
			// We can issue a new auth token
			authToken, err := jwt.ParseWithClaims(authTokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
				return verifyKey, nil
			})

			authTokenClaims, ok := authToken.Claims.(*TokenClaims)
			if !ok {
				return "", err
			}

			return createAuthTokenString(authTokenClaims.ID, authTokenClaims.Role)
		}

		// The refresh token has expired: revoke the token
		deleteRefreshToken(refreshTokenClaims.StandardClaims.Id)

		return "", errors.New("Unauthorized")
	}

	// The refresh token has been revoked!
	return "", errors.New("Unauthorized")
}

func updateRefreshTokenExpiration(refreshTokenString string) (string, error) {
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	refreshTokenClaims, ok := refreshToken.Claims.(*TokenClaims)
	if !ok {
		return "", err
	}

	refreshTokenExpiration := time.Now().Add(refreshTokenValidTime).Unix()

	refreshClaims := TokenClaims{
		ID:   refreshTokenClaims.ID,
		Role: refreshTokenClaims.Role,
		StandardClaims: jwt.StandardClaims{
			Id:        refreshTokenClaims.StandardClaims.Id,
			ExpiresAt: refreshTokenExpiration,
		},
	}

	// Create a signer for rsa 256
	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)

	// Generate the refresh token string
	return refreshJwt.SignedString(signKey)
}

// createRefreshTokenString create a refresh token
func createRefreshTokenString(id bson.ObjectId, role string) (string, error) {
	refreshTokenExpiration := time.Now().Add(refreshTokenValidTime).Unix()

	// Store a token in the database
	token := storeRefreshToken()

	refreshClaims := TokenClaims{
		ID:   id,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			Id:        token.JTI,
			ExpiresAt: refreshTokenExpiration,
		},
	}

	// Create a signer for rsa 256
	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)

	// Generate the refresh token string
	return refreshJwt.SignedString(signKey)
}

func checkRefreshToken(jti string) bool {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("tokens")

	count, err := db.Find(bson.M{
		"jti": jti,
	}).Count()

	return err != nil && count > 0
}

func storeRefreshToken() TokenJTI {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("tokens")

	jti, _ := GenerateRandomString(32)
	for checkRefreshToken(jti) {
		jti, _ = GenerateRandomString(32)
	}

	var token TokenJTI
	token.JTI = jti
	db.Insert(token)

	return token
}

func deleteRefreshToken(jti string) {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("tokens")

	db.Remove(bson.M{"jti": jti})
}
