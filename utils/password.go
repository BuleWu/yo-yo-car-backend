package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash), err
}

func VerifyPassword(hashed string, raw_pwd []byte) bool {
	bytehash := []byte(hashed)
	err := bcrypt.CompareHashAndPassword(bytehash, raw_pwd)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
