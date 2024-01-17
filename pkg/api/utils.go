package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	jwt "github.com/golang-jwt/jwt/v5"
)

// openIDConnectClaims represents an OpenID Connect JWT.
type openIDConnectClaims struct {
	*jwt.RegisteredClaims
	GivenName string `json:"given_name"`
}

// mustGetRandHexStr generates a random 32 byte hex string.
func mustGetRandHexStr() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

// giveToRequest sets the given value to the request's context.
func giveToRequest(r *http.Request, key string, value any) {
	*r = *r.WithContext(context.WithValue(r.Context(), key, value))
}
