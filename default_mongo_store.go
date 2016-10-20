package oauth2

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/phuc0302/go-mongo"
	"github.com/phuc0302/go-oauth2/util"
	"gopkg.in/mgo.v2/bson"
)

// DefaultMongoStore describes a mongodb store.
type DefaultMongoStore struct {
}

// FindUserWithID returns user with user_id.
func (d *DefaultMongoStore) FindUserWithID(userID string) User {
	/* Condition validation */
	if len(userID) == 0 || !bson.IsObjectIdHex(userID) {
		return nil
	}

	user := new(DefaultUser)
	if err := mongo.EntityWithID(TableUser, bson.ObjectIdHex(userID), user); err == nil {
		return user
	}
	return nil
}

// FindUserWithClient returns user associated with client_id and client_secret.
func (d *DefaultMongoStore) FindUserWithClient(clientID string, clientSecret string) User {
	/* Condition validation */
	if len(clientID) == 0 || len(clientSecret) == 0 {
		return nil
	}

	user := new(DefaultUser)
	if err := mongo.EntityWithID(TableUser, bson.ObjectIdHex(clientID), user); err == nil && util.ComparePassword(user.Pass, clientSecret) {
		return user
	}
	return nil
}

// FindUserWithCredential returns user associated with username and password.
func (d *DefaultMongoStore) FindUserWithCredential(username string, password string) User {
	/* Condition validation */
	if len(username) == 0 || len(password) == 0 {
		return nil
	}

	user := new(DefaultUser)
	if err := mongo.EntityWithCriteria(TableUser, bson.M{"username": username}, user); err == nil && util.ComparePassword(user.Pass, password) {
		return user
	}
	return nil
}

// FindClientWithID returns user associated with client_id.
func (d *DefaultMongoStore) FindClientWithID(clientID string) Client {
	/* Condition validation */
	if len(clientID) == 0 || !bson.IsObjectIdHex(clientID) {
		return nil
	}

	client := new(DefaultClient)
	if err := mongo.EntityWithID(TableClient, bson.ObjectIdHex(clientID), client); err == nil {
		return client
	}
	return nil
}

// FindClientWithCredential returns client with client_id and client_secret.
func (d *DefaultMongoStore) FindClientWithCredential(clientID string, clientSecret string) Client {
	/* Condition validation */
	if len(clientID) == 0 || len(clientSecret) == 0 || !bson.IsObjectIdHex(clientID) || !bson.IsObjectIdHex(clientSecret) {
		return nil
	}

	client := new(DefaultClient)
	if err := mongo.EntityWithCriteria(TableClient, bson.M{"_id": bson.ObjectIdHex(clientID), "client_secret": bson.ObjectIdHex(clientSecret)}, client); err == nil {
		return client
	}
	return nil
}

// FindAccessToken returns access_token.
func (d *DefaultMongoStore) FindAccessToken(token string) Token {
	return d.parseToken(token)
}

// FindAccessTokenWithCredential returns access_token associated with client_id and user_id.
func (d *DefaultMongoStore) FindAccessTokenWithCredential(clientID string, userID string) Token {
	return d.queryTokenWithCredential(TableAccessToken, clientID, userID)
}

// CreateAccessToken returns new access_token.
func (d *DefaultMongoStore) CreateAccessToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) Token {
	return d.createToken(TableAccessToken, clientID, userID, createdTime, expiredTime)
}

// DeleteAccessToken deletes access_token.
func (d *DefaultMongoStore) DeleteAccessToken(token Token) {
	d.deleteToken(TableAccessToken, token)
}

// FindRefreshToken returns refresh_token.
func (d *DefaultMongoStore) FindRefreshToken(token string) Token {
	return d.parseToken(token)
}

// FindRefreshTokenWithCredential returns refresh_token associated with client_id and user_id.
func (d *DefaultMongoStore) FindRefreshTokenWithCredential(clientID string, userID string) Token {
	return d.queryTokenWithCredential(TableRefreshToken, clientID, userID)
}

// CreateRefreshToken returns new refresh_token.
func (d *DefaultMongoStore) CreateRefreshToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) Token {
	return d.createToken(TableRefreshToken, clientID, userID, createdTime, expiredTime)
}

// DeleteRefreshToken deletes refresh_token.
func (d *DefaultMongoStore) DeleteRefreshToken(token Token) {
	d.deleteToken(TableRefreshToken, token)
}

//func (d *MongoDBTokenStore) FindAuthorizationCode(authorizationCode string) {
//}
//func (d *MongoDBTokenStore) SaveAuthorizationCode(authorizationCode string, clientID string, expires time.Time) {
//}

// Parse Token
func (d *DefaultMongoStore) parseToken(token string) Token {
	/* Condition validation */
	if len(token) == 0 {
		return nil
	}

	// Parse token
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		/* Condition validation: jwt method should be instance of RSA signing method */
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Invalid signing method: %v", t.Header["alg"])
		}
		return &privateKey.PublicKey, nil
	})

	/* Condition validation: validate parse process */
	if err != nil || !jwtToken.Valid {
		return nil
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		tokenID, _ := claims["_id"].(string)
		userID, _ := claims["user_id"].(string)
		clientID, _ := claims["client_id"].(string)
		createdTime, _ := claims["created_time"].(string)
		expiredTime, _ := claims["expired_time"].(string)
		created, _ := time.Parse(time.RFC3339, createdTime)
		expired, _ := time.Parse(time.RFC3339, expiredTime)

		t := &DefaultToken{
			ID:      bson.ObjectIdHex(tokenID),
			User:    bson.ObjectIdHex(userID),
			Client:  bson.ObjectIdHex(clientID),
			Created: created,
			Expired: expired,
		}
		return t
	}
	return nil
}

// Find token with credential
func (d *DefaultMongoStore) queryTokenWithCredential(table string, clientID string, userID string) Token {
	/* Condition validation */
	if len(clientID) == 0 || len(userID) == 0 || !bson.IsObjectIdHex(clientID) || !bson.IsObjectIdHex(userID) {
		return nil
	}

	token := new(DefaultToken)
	if err := mongo.EntityWithCriteria(table, bson.M{"user_id": bson.ObjectIdHex(userID), "client_id": bson.ObjectIdHex(clientID)}, token); err != nil {
		return nil
	}
	return token
}

// Create token
func (d *DefaultMongoStore) createToken(table string, clientID string, userID string, createdTime time.Time, expiredTime time.Time) Token {
	/* Condition validation */
	if len(clientID) == 0 || len(userID) == 0 || !bson.IsObjectIdHex(clientID) || !bson.IsObjectIdHex(userID) {
		return nil
	}

	newToken := &DefaultToken{
		ID:      bson.NewObjectId(),
		User:    bson.ObjectIdHex(userID),
		Client:  bson.ObjectIdHex(clientID),
		Created: createdTime.UTC(),
		Expired: expiredTime.UTC(),
	}

	if err := mongo.SaveEntity(table, newToken.ID, newToken); err == nil {
		return newToken
	}
	return nil
}

// Delete token
func (d *DefaultMongoStore) deleteToken(table string, token Token) {
	/* Condition validation */
	if token == nil || len(token.ClientID()) == 0 || len(token.UserID()) == 0 || !bson.IsObjectIdHex(token.ClientID()) || !bson.IsObjectIdHex(token.UserID()) {
		return
	}

	u := bson.ObjectIdHex(token.UserID())
	c := bson.ObjectIdHex(token.ClientID())
	mongo.DeleteEntityWithCriteria(table, bson.M{"user_id": u, "client_id": c})
}
