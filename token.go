package oauth2

import (
	"crypto/sha1"
	"fmt"
	"time"
)

func GenerateToken() string {
	timestamp := time.Now().Unix()
	data := []byte(fmt.Sprintf("%d", timestamp))

	hash := fmt.Sprintf("%x", sha1.Sum(data))
	return hash
}
