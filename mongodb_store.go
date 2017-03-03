package oauth2

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/phuc0302/go-mongo"
	"github.com/phuc0302/go-oauth2/oauth_table"
	"github.com/phuc0302/go-oauth2/util"
	"gopkg.in/mgo.v2/bson"
)

// MongoDBStore describes a mongodb store.
type MongoDBStore struct {
}

// FindUserWithID returns user with user_id.
func (d *MongoDBStore) FindUserWithID(userID string) User {
	/* Condition validation */
	if len(userID) == 0 || !bson.IsObjectIdHex(userID) {
		return nil
	}

	user := new(MongoDBUser)
	if err := mongo.EntityWithID(oauthTable.User, bson.ObjectIdHex(userID), user); err == nil {
		return user
	}
	return nil
}

// FindUserWithClient returns user associated with client_id and client_secret.
func (d *MongoDBStore) FindUserWithClient(clientID string, clientSecret string) User {
	/* Condition validation */
	if len(clientID) == 0 || len(clientSecret) == 0 {
		return nil
	}

	user := new(MongoDBUser)
	if err := mongo.EntityWithID(oauthTable.User, bson.ObjectIdHex(clientID), user); err == nil && util.ComparePassword(user.Pass, clientSecret) {
		return user
	}
	return nil
}

// FindUserWithCredential returns user associated with username and password.
func (d *MongoDBStore) FindUserWithCredential(username string, password string) User {
	/* Condition validation */
	if len(username) == 0 || len(password) == 0 {
		return nil
	}

	user := new(MongoDBUser)
	if err := mongo.EntityWithCriteria(oauthTable.User, bson.M{"username": username}, user); err == nil && util.ComparePassword(user.Pass, password) {
		return user
	}
	return nil
}

// FindClientWithID returns user associated with client_id.
func (d *MongoDBStore) FindClientWithID(clientID string) Client {
	/* Condition validation */
	if len(clientID) == 0 || !bson.IsObjectIdHex(clientID) {
		return nil
	}

	client := new(MongoDBClient)
	if err := mongo.EntityWithID(oauthTable.Client, bson.ObjectIdHex(clientID), client); err == nil {
		return client
	}
	return nil
}

// FindClientWithCredential returns client with client_id and client_secret.
func (d *MongoDBStore) FindClientWithCredential(clientID string, clientSecret string) Client {
	/* Condition validation */
	if len(clientID) == 0 || len(clientSecret) == 0 || !bson.IsObjectIdHex(clientID) || !bson.IsObjectIdHex(clientSecret) {
		return nil
	}

	client := new(MongoDBClient)
	if err := mongo.EntityWithCriteria(oauthTable.Client, bson.M{"_id": bson.ObjectIdHex(clientID), "client_secret": bson.ObjectIdHex(clientSecret)}, client); err == nil {
		return client
	}
	return nil
}

// FindAccessToken returns access_token.
func (d *MongoDBStore) FindAccessToken(token string) Token {
	return d.parseToken(token)
}

// FindAccessTokenWithCredential returns access_token associated with client_id and user_id.
func (d *MongoDBStore) FindAccessTokenWithCredential(clientID string, userID string) Token {
	return d.queryTokenWithCredential(oauthTable.AccessToken, clientID, userID)
}

// CreateAccessToken returns new access_token.
func (d *MongoDBStore) CreateAccessToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) Token {
	return d.createToken(oauthTable.AccessToken, clientID, userID, createdTime, expiredTime)
}

// DeleteAccessToken deletes access_token.
func (d *MongoDBStore) DeleteAccessToken(token Token) {
	d.deleteToken(oauthTable.AccessToken, token)
}

// FindRefreshToken returns refresh_token.
func (d *MongoDBStore) FindRefreshToken(token string) Token {
	return d.parseToken(token)
}

// FindRefreshTokenWithCredential returns refresh_token associated with client_id and user_id.
func (d *MongoDBStore) FindRefreshTokenWithCredential(clientID string, userID string) Token {
	return d.queryTokenWithCredential(oauthTable.RefreshToken, clientID, userID)
}

// CreateRefreshToken returns new refresh_token.
func (d *MongoDBStore) CreateRefreshToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) Token {
	return d.createToken(oauthTable.RefreshToken, clientID, userID, createdTime, expiredTime)
}

// DeleteRefreshToken deletes refresh_token.
func (d *MongoDBStore) DeleteRefreshToken(token Token) {
	d.deleteToken(oauthTable.RefreshToken, token)
}

//func (d *MongoDBTokenStore) FindAuthorizationCode(authorizationCode string) {
//}
//func (d *MongoDBTokenStore) SaveAuthorizationCode(authorizationCode string, clientID string, expires time.Time) {
//}

// MARK: Struct's private funcs.
// Parse Token
func (d *MongoDBStore) parseToken(token string) Token {
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

		t := &MongoDBToken{
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
func (d *MongoDBStore) queryTokenWithCredential(table string, clientID string, userID string) Token {
	/* Condition validation */
	if len(clientID) == 0 || len(userID) == 0 || !bson.IsObjectIdHex(clientID) || !bson.IsObjectIdHex(userID) {
		return nil
	}

	token := new(MongoDBToken)
	if err := mongo.EntityWithCriteria(table, bson.M{"user_id": bson.ObjectIdHex(userID), "client_id": bson.ObjectIdHex(clientID)}, token); err != nil {
		return nil
	}
	return token
}

// Create token
func (d *MongoDBStore) createToken(table string, clientID string, userID string, createdTime time.Time, expiredTime time.Time) Token {
	/* Condition validation */
	if len(clientID) == 0 || len(userID) == 0 || !bson.IsObjectIdHex(clientID) || !bson.IsObjectIdHex(userID) {
		return nil
	}

	newToken := &MongoDBToken{
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
func (d *MongoDBStore) deleteToken(table string, token Token) {
	/* Condition validation */
	if token == nil || len(token.ClientID()) == 0 || len(token.UserID()) == 0 || !bson.IsObjectIdHex(token.ClientID()) || !bson.IsObjectIdHex(token.UserID()) {
		return
	}

	if defaultToken, ok := token.(*MongoDBToken); ok {
		mongo.DeleteEntity(table, defaultToken.ID)
	} else {
		u := bson.ObjectIdHex(token.UserID())
		c := bson.ObjectIdHex(token.ClientID())
		mongo.DeleteEntityWithCriteria(table, bson.M{"user_id": u, "client_id": c})
	}
}
