package utils

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// EncryptPassword will cipher user's password using bcrypt
func EncryptPassword(password string) {
	password = strings.Trim(password, " ")
	input := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(input, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hash))

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(hash, input)
	fmt.Println(err) // nil means it is a match
}
