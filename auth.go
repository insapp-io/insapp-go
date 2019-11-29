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
func CreateNewTokens(ID bson.ObjectId, role string) (*jwt.Token, *jwt.Token) {
	return createAuthToken(ID, role), createRefreshToken(ID, role)
}

// CheckAndRefreshStringTokens renews the auth token, if needed.
func CheckAndRefreshStringTokens(authStringToken string, refreshStringToken string, role string) (*jwt.Token, *jwt.Token, error) {
	refreshToken, err := parseRefreshStringToken(refreshStringToken)
	if err != nil {
		return nil, nil, err
	}

	// Don't use parseAuthStringToken:
	// if err2 is not nil, it could be a validation error handled later in this function
	authToken, err2 := jwt.ParseWithClaims(authStringToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	_, ok := authToken.Claims.(*TokenClaims)
	if !ok {
		return nil, nil, errors.New("Auth token parse error")
	}

	// The auth token is still valid
	if authToken.Valid {
		// Check the role
		roles := map[string]int{
			"user":        0,
			"association": 1,
			"admin":       2,
		}
		var level int
		var requiredLevel int
		if level, ok = roles[authToken.Claims.(*TokenClaims).Role]; !ok {
			return nil, nil, errors.New("Unauthorized")
		}
		if requiredLevel, ok = roles[role]; !ok {
			return nil, nil, errors.New("Unauthorized")
		}
		if level < requiredLevel {
			return nil, nil, errors.New("Unauthorized")
		}

		// Update the expiration time of refresh token
		newRefreshToken := updateRefreshTokenExpiration(refreshToken)

		return authToken, newRefreshToken, nil
	} else if ve, ok := err2.(*jwt.ValidationError); ok {
		// The auth token has expired
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			newAuthToken, err := updateAuthToken(authToken, refreshToken)
			if err != nil {
				return nil, nil, err
			}

			// Update the expiration time of refresh token string
			newRefreshToken := updateRefreshTokenExpiration(refreshToken)

			return newAuthToken, newRefreshToken, nil
		}
	}

	return nil, nil, err2
}

// RevokeRefreshStringToken deletes the given token from the database, if valid.
func RevokeRefreshStringToken(refreshStringToken string) error {
	refreshToken, err := parseRefreshStringToken(refreshStringToken)
	if err != nil {
		return err
	}

	deleteRefreshToken(refreshToken.Claims.(*TokenClaims).StandardClaims.Id)

	return nil
}

// GetUserFromRequest returns the User or AssociationUser ID from the auth cookie.
func GetUserFromRequest(r *http.Request) (bson.ObjectId, error) {
	authCookie, err1 := r.Cookie("AuthToken")
	if err1 != nil {
		return bson.ObjectId(""), err1
	}

	authToken, err2 := parseAuthStringToken(authCookie.Value)
	if err2 != nil {
		return bson.ObjectId(""), err2
	}

	return authToken.Claims.(*TokenClaims).ID, nil
}

func parseAuthStringToken(authStringToken string) (*jwt.Token, error) {
	authToken, err := jwt.ParseWithClaims(authStringToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		return nil, err
	}
	_, ok := authToken.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("Auth token parse error")
	}

	return authToken, nil
}

func parseRefreshStringToken(refreshStringToken string) (*jwt.Token, error) {
	refreshToken, err := jwt.ParseWithClaims(refreshStringToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		return nil, err
	}
	_, ok := refreshToken.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("Refresh token parse error")
	}

	return refreshToken, nil
}

// createAuthToken creates an auth token
func createAuthToken(id bson.ObjectId, role string) *jwt.Token {
	authTokenExpiration := time.Now().Add(authTokenValidTime).Unix()

	authClaims := TokenClaims{
		ID:   id,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: authTokenExpiration,
		},
	}

	// Create a signer for rsa 256
	return jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), authClaims)
}

// createRefreshToken create a refresh token
func createRefreshToken(id bson.ObjectId, role string) *jwt.Token {
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
	return jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)
}

func updateAuthToken(authToken *jwt.Token, refreshToken *jwt.Token) (*jwt.Token, error) {
	// Check that the refresh token has not been revoked
	if checkRefreshToken(refreshToken.Claims.(*TokenClaims).StandardClaims.Id) {
		// The refresh token has not expired: issue a new auth token
		if refreshToken.Valid {
			return createAuthToken(authToken.Claims.(*TokenClaims).ID, authToken.Claims.(*TokenClaims).Role), nil
		}

		// The refresh token has expired: revoke the token
		deleteRefreshToken(refreshToken.Claims.(*TokenClaims).StandardClaims.Id)

		return nil, errors.New("Unauthorized")
	}

	// The refresh token has been revoked!
	return nil, errors.New("Unauthorized")
}

func updateRefreshTokenExpiration(refreshToken *jwt.Token) *jwt.Token {
	refreshTokenExpiration := time.Now().Add(refreshTokenValidTime).Unix()

	refreshClaims := TokenClaims{
		ID:   refreshToken.Claims.(*TokenClaims).ID,
		Role: refreshToken.Claims.(*TokenClaims).Role,
		StandardClaims: jwt.StandardClaims{
			Id:        refreshToken.Claims.(*TokenClaims).StandardClaims.Id,
			ExpiresAt: refreshTokenExpiration,
		},
	}

	// Create a signer for rsa 256
	return jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)
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
