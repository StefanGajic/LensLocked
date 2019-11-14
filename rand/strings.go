package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

//generate n random bytes or return an error if theere was one
//use crypto/rand pcg so its safe to use with things like remember token
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//generate a byte slice of size nBytes and than return
//a string that is base64 URL encoded verson of that byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// helper func designed to generate remember tokens of a
//predeterminated byte size
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
