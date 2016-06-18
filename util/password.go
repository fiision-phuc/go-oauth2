package util

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// ComparePassword will compare input password with hash password.
func ComparePassword(hash string, password string) bool {
	password = strings.Trim(password, " ")
	input := []byte(password)

	err := bcrypt.CompareHashAndPassword([]byte(hash), input)
	return (err == nil)
}

// EncryptPassword will cipher user's password using bcrypt.
func EncryptPassword(password string) (string, error) {
	password = strings.Trim(password, " ")
	input := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(input, bcrypt.DefaultCost)
	return string(hash), err
}
