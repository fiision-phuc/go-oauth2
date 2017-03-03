package oauth2

import (
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/phuc0302/go-oauth2/test"

	"gopkg.in/mgo.v2/bson"
)

func Test_MongoDBToken(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	token := MongoDBToken{
		ID:      bson.NewObjectId(),
		User:    bson.NewObjectId(),
		Client:  bson.NewObjectId(),
		Created: time.Now(),
	}
	token.Expired = token.Created.Add(cfg.RefreshTokenDuration)

	// Test token
	tokenString := token.Token()
	if len(tokenString) == 0 {
		t.Error(test.ExpectedNotNil)
	} else {
		jwtToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return &privateKey.PublicKey, nil
		})

		if err != nil || !jwtToken.Valid {
			t.Error("Expected token string should be able to decoded.")
		}
	}
}
