package oauth2

import (
	"time"

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
	/* Condition validation */
	if len(token) == 0 || !bson.IsObjectIdHex(token) {
		return nil
	}
	accessToken := DefaultToken{}

	err := mongo.EntityWithID(TableAccessToken, bson.ObjectIdHex(token), &accessToken)
	if err != nil {
		return nil
	}
	return &accessToken
}

// FindAccessTokenWithCredential returns access_token associated with client_id and user_id.
func (m *DefaultMongoStore) FindAccessTokenWithCredential(clientID string, userID string) IToken {
	/* Condition validation */
	if len(clientID) == 0 || len(userID) == 0 || !bson.IsObjectIdHex(clientID) || !bson.IsObjectIdHex(userID) {
		return nil
	}
	accessToken := DefaultToken{}

	err := mongo.EntityWithCriteria(TableAccessToken, bson.M{"user_id": bson.ObjectIdHex(userID), "client_id": bson.ObjectIdHex(clientID)}, &accessToken)
	if err != nil {
		return nil
	}
	return &accessToken
}

// CreateAccessToken returns new access_token.
func (m *DefaultMongoStore) CreateAccessToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) IToken {
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

	err := mongo.SaveEntity(TableAccessToken, newToken.ID, newToken)
	if err != nil {
		return nil
	}
	return newToken
}

// DeleteAccessToken deletes access_token.
func (m *DefaultMongoStore) DeleteAccessToken(token IToken) {
	/* Condition validation */
	if token == nil {
		return
	}
	mongo.DeleteEntity(TableAccessToken, bson.ObjectIdHex(token.Token()))
}

// FindRefreshToken returns refresh_token.
func (m *DefaultMongoStore) FindRefreshToken(token string) IToken {
	/* Condition validation */
	if len(token) == 0 || !bson.IsObjectIdHex(token) {
		return nil
	}
	refreshToken := DefaultToken{}

	err := mongo.EntityWithID(TableRefreshToken, bson.ObjectIdHex(token), &refreshToken)
	if err != nil {
		return nil
	}
	return &refreshToken
}

// FindRefreshTokenWithCredential returns refresh_token associated with client_id and user_id.
func (m *DefaultMongoStore) FindRefreshTokenWithCredential(clientID string, userID string) IToken {
	/* Condition validation */
	if len(clientID) == 0 || len(userID) == 0 || !bson.IsObjectIdHex(clientID) || !bson.IsObjectIdHex(userID) {
		return nil
	}
	refreshToken := DefaultToken{}

	err := mongo.EntityWithCriteria(TableRefreshToken, bson.M{"user_id": bson.ObjectIdHex(userID), "client_id": bson.ObjectIdHex(clientID)}, &refreshToken)
	if err != nil {
		return nil
	}
	return &refreshToken
}

// CreateRefreshToken returns new refresh_token.
func (m *DefaultMongoStore) CreateRefreshToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) IToken {
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

	err := mongo.SaveEntity(TableRefreshToken, newToken.ID, newToken)
	if err != nil {
		return nil
	}
	return newToken
}

// DeleteRefreshToken deletes refresh_token.
func (m *DefaultMongoStore) DeleteRefreshToken(token IToken) {
	/* Condition validation */
	if token == nil {
		return
	}
	mongo.DeleteEntity(TableRefreshToken, bson.ObjectIdHex(token.Token()))
}

//func (m *MongoDBTokenStore) FindAuthorizationCode(authorizationCode string) {
//}
//func (m *MongoDBTokenStore) SaveAuthorizationCode(authorizationCode string, clientID string, expires time.Time) {
//}
