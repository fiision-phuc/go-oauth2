package utils

import (
	"crypto/sha1"
	"fmt"
	"time"
)

// GenerateToken return a random token base on timestamp input.
func GenerateToken() string {
	timestamp := time.Now().Unix()
	data := []byte(fmt.Sprintf("%d", timestamp))

	hash := fmt.Sprintf("%x", sha1.Sum(data))
	return hash
}
