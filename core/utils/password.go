package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash), err
}

func VerifyPassword(hashed string, rawPwd string) bool {
	bytehash := []byte(hashed)
	byteRawPwd := []byte(rawPwd)
	err := bcrypt.CompareHashAndPassword(bytehash, byteRawPwd)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
