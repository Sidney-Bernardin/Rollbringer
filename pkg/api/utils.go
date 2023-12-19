package api

import (
	"crypto/rand"
	"encoding/hex"

	jwt "github.com/golang-jwt/jwt/v5"
)

type openIDConnectClaims struct {
	*jwt.RegisteredClaims
	GivenName string `json:"given_name"`
}

func mustGetRandHexStr() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
