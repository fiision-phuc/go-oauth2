package oauth2

import (
	"os"
	"testing"
	"time"

	"github.com/phuc0302/go-oauth2/mongo"
	"github.com/phuc0302/go-oauth2/utils"
)

func Test_DefaultMongoStore(t *testing.T) {
	defer os.Remove(mongo.ConfigFile)
	mongo.ConnectMongo()

	// Prepare data
	session, database := mongo.GetMonotonicSession()
	defer session.Close()
	database.C("client").DropCollection()
	database.C("user").DropCollection()

	password1, _ := utils.EncryptPassword("admin")
	user1 := &DefaultUser{
		ID:    userID,
		User:  "admin",
		Pass:  password1,
		Roles: []string{"r_user", "r_admin"},
	}

	password2, _ := utils.EncryptPassword(clientSecret.Hex())
	user2 := &DefaultUser{
		ID:    clientID,
		User:  clientID.Hex(),
		Pass:  password2,
		Roles: []string{"r_device"},
	}
	database.C("user").Insert(user1, user2)

	client1 := &DefaultClient{
		ID:     clientID,
		Secret: clientSecret,
		Grants: []string{PasswordGrant, RefreshTokenGrant},

		Redirects: []string{"http://www.sample01.com", "http://www.sample02.com"},
	}
	database.C("client").Insert(client1)

	// Create store
	store := DefaultMongoStore{}

	// Testing process
	recordUser1 := store.FindUserWithID(userID.Hex())
	if recordUser1 == nil || recordUser1.UserID() != user1.ID.Hex() || recordUser1.Username() != user1.User || recordUser1.Password() != user1.Pass {
		t.Errorf("Expected %s but found %s", user1, recordUser1)
	}
	recordUser2 := store.FindUserWithClient(clientID.Hex(), clientSecret.Hex())
	if recordUser2 == nil || recordUser2.UserID() != user2.ID.Hex() || recordUser2.Username() != user2.User || recordUser2.Password() != user2.Pass {
		t.Errorf("Expected %s but found %s", user2, recordUser2)
	}
	recordUser3 := store.FindUserWithCredential("admin", "admin")
	if recordUser3 == nil || recordUser3.UserID() != user1.ID.Hex() || recordUser3.Username() != user1.User || recordUser3.Password() != user1.Pass {
		t.Errorf("Expected %s but found %s", user1, recordUser1)
	}

	// Client
	recordClient1 := store.FindClientWithID(clientID.Hex())
	if recordClient1 == nil || recordClient1.ClientID() != client1.ID.Hex() || recordClient1.ClientSecret() != client1.Secret.Hex() {
		t.Errorf("Expected %s but found %s", client1, recordClient1)
	}
	recordClient2 := store.FindClientWithCredential(clientID.Hex(), clientSecret.Hex())
	if recordClient2 == nil || recordClient2.ClientID() != client1.ID.Hex() || recordClient2.ClientSecret() != client1.Secret.Hex() {
		t.Errorf("Expected %s but found %s", client1, recordClient2)
	}

	// Access token
	token1 := store.CreateAccessToken(clientID.Hex(), userID.Hex(), time.Now(), time.Now().Add(3600))
	if token1 == nil {
		t.Errorf("Expected not nil but found nil.")
	}
	if token1.ClientID() != clientID.Hex() {
		t.Errorf("Expected %s but found %s", clientID.Hex(), token1.ClientID())
	}
	if token1.UserID() != userID.Hex() {
		t.Errorf("Expected %s but found %s", userID.Hex(), token1.ClientID())
	}
	token2 := store.FindAccessToken(token1.Token())
	if token2 == nil {
		t.Errorf("Expected not nil but found nil.")
	}
	token3 := store.FindAccessTokenWithCredential(token1.ClientID(), token1.UserID())
	if token3 == nil {
		t.Errorf("Expected not nil but found nil.")
	}
	store.DeleteAccessToken(token1)
	token4 := store.FindAccessTokenWithCredential(token1.ClientID(), token1.UserID())
	if token4 != nil {
		t.Errorf("Expected nil but found not nil.")
	}

	// Refresh token
	refreshToken1 := store.CreateRefreshToken(clientID.Hex(), userID.Hex(), time.Now(), time.Now().Add(3600))
	if refreshToken1 == nil {
		t.Errorf("Expected not nil but found nil.")
	}
	if refreshToken1.ClientID() != clientID.Hex() {
		t.Errorf("Expected %s but found %s", clientID.Hex(), refreshToken1.ClientID())
	}
	if refreshToken1.UserID() != userID.Hex() {
		t.Errorf("Expected %s but found %s", userID.Hex(), refreshToken1.ClientID())
	}
	refreshToken2 := store.FindRefreshToken(refreshToken1.Token())
	if refreshToken2 == nil {
		t.Errorf("Expected not nil but found nil.")
	}
	refreshToken3 := store.FindRefreshTokenWithCredential(refreshToken1.ClientID(), refreshToken1.UserID())
	if refreshToken3 == nil {
		t.Errorf("Expected not nil but found nil.")
	}
	store.DeleteRefreshToken(refreshToken1)
	refreshToken4 := store.FindRefreshTokenWithCredential(refreshToken1.ClientID(), refreshToken1.UserID())
	if refreshToken4 != nil {
		t.Errorf("Expected nil but found not nil.")
	}
}
