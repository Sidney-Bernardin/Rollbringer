package handler

import (
	"crypto/rand"
	"encoding/hex"

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
