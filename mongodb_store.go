package oauth2

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/phuc0302/go-mongo"
	"github.com/phuc0302/go-oauth2/oauth_table"
	"github.com/phuc0302/go-server"
	"github.com/phuc0302/go-server/util"
	"gopkg.in/mgo.v2/bson"
)

// MongoDBStore describes a mongodb token store.
type MongoDBStore struct {
	privateKey *rsa.PrivateKey
}

// CreateMongoDBStore return a default MongoDBStore's instance.
//
// @return
// - tokenStore {TokenStore} (a mongoDB token store's instance)
func CreateMongoDBStore() (tokenStore TokenStore) {
	if server.Cfg == nil {
		panic("Please call server.Initialize before create store.")
	}

	var privateKey *rsa.PrivateKey
	if base64Encoded, ok := server.Cfg.GetExtension("jwt_key").(string); ok {
		if keyDER, err := base64.StdEncoding.DecodeString(base64Encoded); err == nil {
			privateKey, _ = x509.ParsePKCS1PrivateKey(keyDER)
		}
	}

	// Generate JWT key if necessary
	if privateKey == nil {
		privateKey, _ = rsa.GenerateKey(rand.Reader, 1024)
		keyDER := x509.MarshalPKCS1PrivateKey(privateKey)

		// Save to config
		base64Encoded := base64.StdEncoding.EncodeToString(keyDER)
		server.Cfg.SetExtension("jwt_key", base64Encoded)
		server.Cfg.Save()
	}

	return &MongoDBStore{
		privateKey: privateKey,
	}
}

// FindUserWithID returns an user entity according to userID or null. A user entity can either
// human or machine.
//
// @param
// - userID {string} (userID that associated with user's entity)
//
// @return
// - user {User} (an user entity or null)
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

// FindUserWithClient returns a machine user entity.
//
// @param
// - clientID {string} (client's client_id)
// - clientSecret {string} (client's client_secret)
//
// @return
// - user {User} (a machine user entity or null)
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

// FindUserWithCredential returns a human user entity.
//
// @param
// - username {string} (user's username)
// - password {string} (user's password)
//
// @return
// - user {User} (a human user entity or null)
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

// FindClientWithID returns a client entity according to clientID or null.
//
// @param
// - clientID {string} (client's client_id)
//
// @return
// - client {Client} (a client entity or null)
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

// FindClientWithCredential returns a client entity according to clientID and clientSecret or
// null.
//
// @param
// - clientID {string} (client's client_id)
// - clientSecret {string} (client's client_secret)
//
// @return
// - client {Client} (a client entity or null)
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

// FindAccessToken returns an access token entity according to token string or null.
//
// @param
// - token {string} (user's access token in string form)
//
// @return
// - token {Token} (a token's instance or null)
func (d *MongoDBStore) FindAccessToken(token string) Token {
	return d.parseToken(token)
}

// FindAccessTokenWithCredential returns an access token entity according to clientID and
// userID or null.
//
// @param
// - clientID {string} (client's client_id)
// - userID {string} (userID that associated with user's entity)
//
// @return
// - token {Token} (a token's instance or null)
func (d *MongoDBStore) FindAccessTokenWithCredential(clientID string, userID string) Token {
	return d.queryTokenWithCredential(oauthTable.AccessToken, clientID, userID)
}

// CreateAccessToken creates a token's instance.
//
// @param
// - clientID {string} (client's client_id)
// - userID {string} (userID that associated with user's entity)
// - createdTime {time.Time} (token's issued time)
// - expiredTime {time.Time} (token's expired time)
//
// @return
// - token {Token} (an access token's instance)
func (d *MongoDBStore) CreateAccessToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) Token {
	return d.createToken(oauthTable.AccessToken, clientID, userID, createdTime, expiredTime)
}

// DeleteAccessToken deletes an access token from database.
//
// @param
// - token {Token} (an access token's instance)
func (d *MongoDBStore) DeleteAccessToken(token Token) {
	d.deleteToken(oauthTable.AccessToken, token)
}

// FindRefreshToken returns a refresh token entity according to token string or null.
//
// @param
// - token {string} (user's refresh token in string form)
//
// @return
// - token {Token} (a token's instance or null)
func (d *MongoDBStore) FindRefreshToken(token string) Token {
	return d.parseToken(token)
}

// FindRefreshTokenWithCredential returns a refresh token entity according to clientID and
// userID or null.
//
// @param
// - clientID {string} (client's client_id)
// - userID {string} (userID that associated with user's entity)
//
// @return
// - token {Token} (a token's instance or null)
func (d *MongoDBStore) FindRefreshTokenWithCredential(clientID string, userID string) Token {
	return d.queryTokenWithCredential(oauthTable.RefreshToken, clientID, userID)
}

// CreateRefreshToken creates a token's instance.
//
// @param
// - clientID {string} (client's client_id)
// - userID {string} (userID that associated with user's entity)
// - createdTime {time.Time} (token's issued time)
// - expiredTime {time.Time} (token's expired time)
//
// @return
// - token {Token} (a refresh token's instance)
func (d *MongoDBStore) CreateRefreshToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) Token {
	return d.createToken(oauthTable.RefreshToken, clientID, userID, createdTime, expiredTime)
}

// DeleteRefreshToken deletes a refresh token from database.
//
// @param
// - token {Token} (a refresh token's instance)
func (d *MongoDBStore) DeleteRefreshToken(token Token) {
	d.deleteToken(oauthTable.RefreshToken, token)
}

//func (d *MongoDBTokenStore) FindAuthorizationCode(authorizationCode string) {
//}
//func (d *MongoDBTokenStore) SaveAuthorizationCode(authorizationCode string, clientID string, expires time.Time) {
//}

// parseToken convert JWT token to token's instance.
//
// @param
// - token {string} (user's token in string form)
//
// @return
// - token {Token} (a token's instance or null)
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
		return &d.privateKey.PublicKey, nil
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

			privateKey: d.privateKey,
		}
		return t
	}
	return nil
}

// queryTokenWithCredential returns a token entity base on search criteria.
//
// @param
// - table {string} (access token table or refresh token table)
// - clientID {string} (client's client_id)
// - userID {string} (userID that associated with user's entity)
//
// @return
// - token {Token} (a token's instance or null)
func (d *MongoDBStore) queryTokenWithCredential(table string, clientID string, userID string) Token {
	/* Condition validation */
	if len(clientID) == 0 || len(userID) == 0 || !bson.IsObjectIdHex(clientID) || !bson.IsObjectIdHex(userID) {
		return nil
	}

	var token MongoDBToken
	if err := mongo.EntityWithCriteria(table, bson.M{"user_id": bson.ObjectIdHex(userID), "client_id": bson.ObjectIdHex(clientID)}, &token); err != nil {
		return nil
	}

	token.privateKey = d.privateKey
	return &token
}

// createToken creates new token's instance.
//
// @param
// - table {string} (access token table or refresh token table)
// - clientID {string} (client's client_id)
// - userID {string} (userID that associated with user's entity)
// - createdTime {time.Time} (token's issued time)
// - expiredTime {time.Time} (token's expired time)
//
// @return
// - token {Token} (a token's instance)
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

		privateKey: d.privateKey,
	}

	if err := mongo.SaveEntity(table, newToken.ID, newToken); err == nil {
		return newToken
	}
	return nil
}

// deleteToken deletes a token from database.
//
// @param
// - table {string} (access token table or refresh token table)
// - token {Token} (a token's instance)
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
