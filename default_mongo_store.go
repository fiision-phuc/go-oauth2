package oauth2

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/phuc0302/go-oauth2/mongo"
	"github.com/phuc0302/go-oauth2/utils"
	"gopkg.in/mgo.v2/bson"
)

// DefaultMongoStore describes a mongodb store.
type DefaultMongoStore struct {
}

// FindUserWithID returns user with user_id.
func (m *DefaultMongoStore) FindUserWithID(userID string) IUser {
	/* Condition validation */
	if len(userID) == 0 || !bson.IsObjectIdHex(userID) {
		return nil
	}
	user := DefaultUser{}

	err := mongo.EntityWithID(TableUser, bson.ObjectIdHex(userID), &user)
	if err != nil {
		return nil
	}
	return &user
}

// FindUserWithClient returns user associated with client_id and client_secret.
func (m *DefaultMongoStore) FindUserWithClient(clientID string, clientSecret string) IUser {
	/* Condition validation */
	if len(clientID) == 0 || len(clientSecret) == 0 {
		return nil
	}
	user := DefaultUser{}

	err := mongo.EntityWithID(TableUser, bson.ObjectIdHex(clientID), &user)
	if err != nil {
		return nil
	}

	if !utils.ComparePassword(user.Pass, clientSecret) {
		return nil
	}
	return &user
}

// FindUserWithCredential returns user associated with username and password.
func (m *DefaultMongoStore) FindUserWithCredential(username string, password string) IUser {
	/* Condition validation */
	if len(username) == 0 || len(password) == 0 {
		return nil
	}
	user := DefaultUser{}

	err := mongo.EntityWithCriteria(TableUser, bson.M{"username": username}, &user)
	if err != nil {
		return nil
	}

	if !utils.ComparePassword(user.Pass, password) {
		return nil
	}
	return &user
}

// FindClientWithID returns user associated with client_id.
func (m *DefaultMongoStore) FindClientWithID(clientID string) IClient {
	/* Condition validation */
	if len(clientID) == 0 || !bson.IsObjectIdHex(clientID) {
		return nil
	}
	client := DefaultClient{}

	err := mongo.EntityWithID(TableClient, bson.ObjectIdHex(clientID), &client)
	if err != nil {
		return nil
	}
	return &client
}

// FindClientWithCredential returns client with client_id and client_secret.
func (m *DefaultMongoStore) FindClientWithCredential(clientID string, clientSecret string) IClient {
	/* Condition validation */
	if len(clientID) == 0 || len(clientSecret) == 0 || !bson.IsObjectIdHex(clientID) || !bson.IsObjectIdHex(clientSecret) {
		return nil
	}
	client := DefaultClient{}

	err := mongo.EntityWithCriteria(TableClient, bson.M{"_id": bson.ObjectIdHex(clientID), "client_secret": bson.ObjectIdHex(clientSecret)}, &client)
	if err != nil {
		return nil
	}
	return &client
}

// FindAccessToken returns access_token.
func (m *DefaultMongoStore) FindAccessToken(token string) IToken {
	return m.parseToken(token)
}

// FindAccessTokenWithCredential returns access_token associated with client_id and user_id.
func (m *DefaultMongoStore) FindAccessTokenWithCredential(clientID string, userID string) IToken {
	return m.queryTokenWithCredential(TableAccessToken, clientID, userID)
}

// CreateAccessToken returns new access_token.
func (m *DefaultMongoStore) CreateAccessToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) IToken {
	return m.createToken(TableAccessToken, clientID, userID, createdTime, expiredTime)
}

// DeleteAccessToken deletes access_token.
func (m *DefaultMongoStore) DeleteAccessToken(token IToken) {
	m.deleteToken(TableAccessToken, token)
}

// FindRefreshToken returns refresh_token.
func (m *DefaultMongoStore) FindRefreshToken(token string) IToken {
	return m.parseToken(token)
}

// FindRefreshTokenWithCredential returns refresh_token associated with client_id and user_id.
func (m *DefaultMongoStore) FindRefreshTokenWithCredential(clientID string, userID string) IToken {
	return m.queryTokenWithCredential(TableRefreshToken, clientID, userID)
}

// CreateRefreshToken returns new refresh_token.
func (m *DefaultMongoStore) CreateRefreshToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) IToken {
	return m.createToken(TableRefreshToken, clientID, userID, createdTime, expiredTime)
}

// DeleteRefreshToken deletes refresh_token.
func (m *DefaultMongoStore) DeleteRefreshToken(token IToken) {
	m.deleteToken(TableRefreshToken, token)
}

//func (m *MongoDBTokenStore) FindAuthorizationCode(authorizationCode string) {
//}
//func (m *MongoDBTokenStore) SaveAuthorizationCode(authorizationCode string, clientID string, expires time.Time) {
//}

// Parse Token
func (m *DefaultMongoStore) parseToken(token string) IToken {
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

	tokenID, _ := jwtToken.Claims["_id"].(string)
	userID, _ := jwtToken.Claims["user_id"].(string)
	clientID, _ := jwtToken.Claims["client_id"].(string)
	createdTime, _ := jwtToken.Claims["created_time"].(string)
	expiredTime, _ := jwtToken.Claims["expired_time"].(string)
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

// Find token with credential
func (m *DefaultMongoStore) queryTokenWithCredential(table string, clientID string, userID string) IToken {
	/* Condition validation */
	if len(clientID) == 0 || len(userID) == 0 || !bson.IsObjectIdHex(clientID) || !bson.IsObjectIdHex(userID) {
		return nil
	}

	accessToken := DefaultToken{}
	if err := mongo.EntityWithCriteria(table, bson.M{"user_id": bson.ObjectIdHex(userID), "client_id": bson.ObjectIdHex(clientID)}, &accessToken); err != nil {
		return nil
	}
	return &accessToken
}

// Create token
func (m *DefaultMongoStore) createToken(table string, clientID string, userID string, createdTime time.Time, expiredTime time.Time) IToken {
	/* Condition validation */
	if len(clientID) == 0 || len(userID) == 0 || !bson.IsObjectIdHex(clientID) || !bson.IsObjectIdHex(userID) {
		return nil
	}

	newToken := &DefaultToken{
		ID:      bson.NewObjectId(),
		User:    bson.ObjectIdHex(userID),
		Client:  bson.ObjectIdHex(clientID),
		Created: createdTime,
		Expired: expiredTime,
	}

	if err := mongo.SaveEntity(table, newToken.ID, newToken); err != nil {
		return nil
	}
	return newToken
}

// Delete token
func (m *DefaultMongoStore) deleteToken(table string, token IToken) {
	/* Condition validation */
	if token == nil || len(token.ClientID()) == 0 || len(token.UserID()) == 0 || !bson.IsObjectIdHex(token.ClientID()) || !bson.IsObjectIdHex(token.UserID()) {
		return
	}

	u := bson.ObjectIdHex(token.UserID())
	c := bson.ObjectIdHex(token.ClientID())
	mongo.DeleteEntityWithCriteria(table, bson.M{"user_id": u, "client_id": c})
}
