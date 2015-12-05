package utils

import (
	"crypto/sha1"
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// GenerateToken return a random token base on timestamp input.
func GenerateToken() string {
	timestamp := time.Now().Unix()
	data := []byte(fmt.Sprintf("%s%d", bson.NewObjectId().Hex(), timestamp))

	hash := fmt.Sprintf("%x", sha1.Sum(data))
	return hash
}
