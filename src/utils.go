package src

import (
	"crypto/rand"
	"encoding/hex"
)

func Map[S []E, E, R any](s S, callback func(int, E) *R) []R {
	ret := []R{}
	for i, e := range s {
		if r := callback(i, e); r != nil {
			ret = append(ret, *r)
		}
	}
	return ret
}

func CreateRandomString() string {
	var bState = make([]byte, 32)
	if _, err := rand.Read(bState); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bState)
}
