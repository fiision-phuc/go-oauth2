package oauth2

import (
	"os"
	"testing"
	"time"

	"github.com/phuc0302/go-oauth2/mongo"
	"github.com/phuc0302/go-oauth2/utils"
)

//	// Access Token
//	FindAccessToken(token string) Token
//	FindAccessTokenWithCredential(clientID string, userID string) Token
//	CreateAccessToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) Token
//	DeleteAccessToken(token Token)

//	// Refresh Token
//	FindRefreshToken(token string) Token
//	FindRefreshTokenWithCredential(clientID string, userID string) Token
//	CreateRefreshToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) Token
//	DeleteRefreshToken(token Token)

func Test_MongoDBTokenStore(t *testing.T) {
	defer os.Remove(mongo.ConfigFile)
	mongo.ConnectMongo()

	// Prepare data
	session, database := mongo.GetMonotonicSession()
	defer session.Close()
	database.C("client").DropCollection()
	database.C("user").DropCollection()

	password1, _ := utils.EncryptPassword("admin")
	user1 := &AuthUserDefault{
		UserID:   userID,
		Username: "admin",
		Password: password1,
		Roles:    []string{"r_user", "r_admin"},
	}

	password2, _ := utils.EncryptPassword(clientSecret.Hex())
	user2 := &AuthUserDefault{
		UserID:   clientID,
		Username: clientID.Hex(),
		Password: password2,
		Roles:    []string{"r_device"},
	}
	database.C("user").Insert(user1, user2)

	client1 := &AuthClientDefault{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantTypes:   []string{PasswordGrant, RefreshTokenGrant},

		RedirectURIs: []string{"http://www.sample01.com", "http://www.sample02.com"},
	}
	database.C("client").Insert(client1)

	// Create store
	store := MongoDBTokenStore{}

	// Testing process
	recordUser1 := store.FindUserWithID(userID.Hex())
	if recordUser1 == nil || recordUser1.GetUserID() != user1.UserID.Hex() || recordUser1.GetUsername() != user1.Username || recordUser1.GetPassword() != user1.Password {
		t.Errorf("Expected %s but found %s", user1, recordUser1)
	}
	recordUser2 := store.FindUserWithClient(clientID.Hex(), clientSecret.Hex())
	if recordUser2 == nil || recordUser2.GetUserID() != user2.UserID.Hex() || recordUser2.GetUsername() != user2.Username || recordUser2.GetPassword() != user2.Password {
		t.Errorf("Expected %s but found %s", user2, recordUser2)
	}
	recordUser3 := store.FindUserWithCredential("admin", "admin")
	if recordUser3 == nil || recordUser3.GetUserID() != user1.UserID.Hex() || recordUser3.GetUsername() != user1.Username || recordUser3.GetPassword() != user1.Password {
		t.Errorf("Expected %s but found %s", user1, recordUser1)
	}

	// Client
	recordClient1 := store.FindClientWithID(clientID.Hex())
	if recordClient1 == nil || recordClient1.GetClientID() != client1.ClientID.Hex() || recordClient1.GetClientSecret() != client1.ClientSecret.Hex() {
		t.Errorf("Expected %s but found %s", client1, recordClient1)
	}
	recordClient2 := store.FindClientWithCredential(clientID.Hex(), clientSecret.Hex())
	if recordClient2 == nil || recordClient2.GetClientID() != client1.ClientID.Hex() || recordClient2.GetClientSecret() != client1.ClientSecret.Hex() {
		t.Errorf("Expected %s but found %s", client1, recordClient2)
	}

	// Access token
	token1 := store.CreateAccessToken(clientID.Hex(), userID.Hex(), time.Now(), time.Now().Add(3600))
	if token1 == nil {
		t.Errorf("Expected not nil but found nil.")
	}
	if token1.GetClientID() != clientID.Hex() {
		t.Errorf("Expected %s but found %s", clientID.Hex(), token1.GetClientID())
	}
	if token1.GetUserID() != userID.Hex() {
		t.Errorf("Expected %s but found %s", userID.Hex(), token1.GetClientID())
	}
	token2 := store.FindAccessToken(token1.GetToken())
	if token2 == nil {
		t.Errorf("Expected not nil but found nil.")
	}
	token3 := store.FindAccessTokenWithCredential(token1.GetClientID(), token1.GetUserID())
	if token3 == nil {
		t.Errorf("Expected not nil but found nil.")
	}
	store.DeleteAccessToken(token1)
	token4 := store.FindAccessTokenWithCredential(token1.GetClientID(), token1.GetUserID())
	if token4 != nil {
		t.Errorf("Expected nil but found not nil.")
	}

	// Refresh token
	refreshToken1 := store.CreateRefreshToken(clientID.Hex(), userID.Hex(), time.Now(), time.Now().Add(3600))
	if refreshToken1 == nil {
		t.Errorf("Expected not nil but found nil.")
	}
	if refreshToken1.GetClientID() != clientID.Hex() {
		t.Errorf("Expected %s but found %s", clientID.Hex(), refreshToken1.GetClientID())
	}
	if refreshToken1.GetUserID() != userID.Hex() {
		t.Errorf("Expected %s but found %s", userID.Hex(), refreshToken1.GetClientID())
	}
	refreshToken2 := store.FindRefreshToken(refreshToken1.GetToken())
	if refreshToken2 == nil {
		t.Errorf("Expected not nil but found nil.")
	}
	refreshToken3 := store.FindRefreshTokenWithCredential(refreshToken1.GetClientID(), refreshToken1.GetUserID())
	if refreshToken3 == nil {
		t.Errorf("Expected not nil but found nil.")
	}
	store.DeleteRefreshToken(refreshToken1)
	refreshToken4 := store.FindRefreshTokenWithCredential(refreshToken1.GetClientID(), refreshToken1.GetUserID())
	if refreshToken4 != nil {
		t.Errorf("Expected nil but found not nil.")
	}
}
