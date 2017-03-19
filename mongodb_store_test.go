package oauth2

import (
	"reflect"
	"testing"
	"time"

	"github.com/phuc0302/go-server"
	"github.com/phuc0302/go-server/expected_format"
)

func Test_MongoDBStore_CreateMongoDBStore(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	if mongoStore, ok := Store.(*MongoDBStore); ok {
		if mongoStore.privateKey == nil {
			t.Error(expectedFormat.NotNil)
		}

		if _, ok := server.Cfg.GetExtension("jwt_key").(string); ok {
			// Everything is fine.
		} else {
			t.Error(expectedFormat.NotNil)
		}
	} else {
		t.Error("Expected Store should be instance of M")
	}
}

func Test_MongoDBStore_FindUserWithID(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	user := Store.FindUserWithID(u.UserID.Hex())
	if user == nil {
		t.Errorf(expectedFormat.StringButFoundString, u.User1, user)
	} else {
		if user.UserID() != u.UserID.Hex() {
			t.Errorf(expectedFormat.StringButFoundString, u.UserID.Hex(), user.UserID())
		}
		if user.Username() != "admin" {
			t.Errorf(expectedFormat.StringButFoundString, "admin", user.Username())
		}
		if !reflect.DeepEqual(user.UserRoles(), []string{"r_user", "r_admin"}) {
			t.Errorf(expectedFormat.StringButFoundString, []string{"r_user", "r_admin"}, user.UserRoles())
		}
	}
}

func Test_MongoDBStore_FindUserWithClient(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	user := Store.FindUserWithClient(u.ClientID.Hex(), u.ClientSecret.Hex())
	if user == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if user.UserID() != u.ClientID.Hex() {
			t.Errorf(expectedFormat.StringButFoundString, u.ClientID.Hex(), user.UserID())
		}
		if user.Username() != u.ClientID.Hex() {
			t.Errorf(expectedFormat.StringButFoundString, u.ClientID.Hex(), user.Username())
		}
		if !reflect.DeepEqual(user.UserRoles(), []string{"r_device"}) {
			t.Errorf(expectedFormat.StringButFoundString, []string{"r_device"}, user.UserRoles())
		}
	}
}

func Test_MongoDBStore_FindUserWithCredential(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	user := Store.FindUserWithCredential("admin", "Password")
	if user == nil {
		t.Errorf(expectedFormat.StringButFoundString, u.User1, user)
	} else {
		if user.UserID() != u.UserID.Hex() {
			t.Errorf(expectedFormat.StringButFoundString, u.UserID.Hex(), user.UserID())
		}
		if user.Username() != "admin" {
			t.Errorf(expectedFormat.StringButFoundString, "admin", user.Username())
		}
		if !reflect.DeepEqual(user.UserRoles(), []string{"r_user", "r_admin"}) {
			t.Errorf(expectedFormat.StringButFoundString, []string{"r_user", "r_admin"}, user.UserRoles())
		}
	}
}

func Test_MongoDBStore_FindClientWithID(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	client := Store.FindClientWithID(u.ClientID.Hex())
	if client == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if client.ClientID() != u.ClientID.Hex() {
			t.Errorf(expectedFormat.StringButFoundString, u.ClientID.Hex(), client.ClientID())
		}
		if client.ClientSecret() != u.ClientSecret.Hex() {
			t.Errorf(expectedFormat.StringButFoundString, u.ClientSecret.Hex(), client.ClientSecret())
		}
		if !reflect.DeepEqual(client.GrantTypes(), []string{AuthorizationCodeGrant, PasswordGrant, RefreshTokenGrant}) {
			t.Errorf(expectedFormat.StringButFoundString, []string{AuthorizationCodeGrant, PasswordGrant, RefreshTokenGrant}, client.GrantTypes())
		}
		if !reflect.DeepEqual(client.RedirectURIs(), []string{"http://www.sample01.com", "http://www.sample02.com"}) {
			t.Errorf(expectedFormat.StringButFoundString, []string{"http://www.sample01.com", "http://www.sample02.com"}, client.RedirectURIs())
		}
	}
}

func Test_MongoDBStore_FindClientWithCredential(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	client := Store.FindClientWithCredential(u.ClientID.Hex(), u.ClientSecret.Hex())
	if client == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if client.ClientID() != u.ClientID.Hex() {
			t.Errorf(expectedFormat.StringButFoundString, u.ClientID.Hex(), client.ClientID())
		}
		if client.ClientSecret() != u.ClientSecret.Hex() {
			t.Errorf(expectedFormat.StringButFoundString, u.ClientSecret.Hex(), client.ClientSecret())
		}
		if !reflect.DeepEqual(client.GrantTypes(), []string{AuthorizationCodeGrant, PasswordGrant, RefreshTokenGrant}) {
			t.Errorf(expectedFormat.StringButFoundString, []string{AuthorizationCodeGrant, PasswordGrant, RefreshTokenGrant}, client.GrantTypes())
		}
		if !reflect.DeepEqual(client.RedirectURIs(), []string{"http://www.sample01.com", "http://www.sample02.com"}) {
			t.Errorf(expectedFormat.StringButFoundString, []string{"http://www.sample01.com", "http://www.sample02.com"}, client.RedirectURIs())
		}
	}
}

func Test_MongoDBStore_CreateAccessToken(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	token := Store.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	if token == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if token.ClientID() != u.ClientID.Hex() {
			t.Errorf(expectedFormat.StringButFoundString, u.ClientID.Hex(), token.ClientID())
		}
		if token.UserID() != u.UserID.Hex() {
			t.Errorf(expectedFormat.StringButFoundString, u.UserID.Hex(), token.ClientID())
		}
		if token.IsExpired() {
			t.Errorf(expectedFormat.BoolButFoundBool, false, token.IsExpired())
		}
	}
}

func Test_MongoDBStore_FindAccessToken(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	token1 := Store.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	token2 := Store.FindAccessToken(token1.Token())
	if token2 == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if token2.ClientID() != token1.ClientID() {
			t.Errorf(expectedFormat.StringButFoundString, token1.ClientID(), token2.ClientID())
		}
		if token2.UserID() != token1.UserID() {
			t.Errorf(expectedFormat.StringButFoundString, token1.UserID(), token2.UserID())
		}
		if token2.IsExpired() {
			t.Errorf(expectedFormat.BoolButFoundBool, false, token2.IsExpired())
		}
	}
}

func Test_MongoDBStore_FindAccessTokenWithCredential(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	token1 := Store.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	token2 := Store.FindAccessTokenWithCredential(u.ClientID.Hex(), u.UserID.Hex())
	if token2 == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if token2.ClientID() != token1.ClientID() {
			t.Errorf(expectedFormat.StringButFoundString, token1.ClientID(), token2.ClientID())
		}
		if token2.UserID() != token1.UserID() {
			t.Errorf(expectedFormat.StringButFoundString, token1.UserID(), token2.UserID())
		}
		if token2.IsExpired() {
			t.Errorf(expectedFormat.BoolButFoundBool, false, token2.IsExpired())
		}
	}
}

func Test_MongoDBStore_DeleteAccessToken(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	token1 := Store.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	Store.DeleteAccessToken(token1)

	token2 := Store.FindAccessTokenWithCredential(token1.ClientID(), token1.UserID())
	if token2 != nil {
		t.Errorf(expectedFormat.Nil)
	}
}

func Test_MongoDBStore_CreateRefreshToken(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	token := Store.CreateRefreshToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	if token == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if token.ClientID() != u.ClientID.Hex() {
			t.Errorf(expectedFormat.StringButFoundString, u.ClientID.Hex(), token.ClientID())
		}
		if token.UserID() != u.UserID.Hex() {
			t.Errorf(expectedFormat.StringButFoundString, u.UserID.Hex(), token.ClientID())
		}
		if token.IsExpired() {
			t.Errorf(expectedFormat.BoolButFoundBool, false, token.IsExpired())
		}
	}
}

func Test_MongoDBStore_FindRefreshToken(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	token1 := Store.CreateRefreshToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	token2 := Store.FindRefreshToken(token1.Token())
	if token2 == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if token2.ClientID() != token1.ClientID() {
			t.Errorf(expectedFormat.StringButFoundString, token1.ClientID(), token2.ClientID())
		}
		if token2.UserID() != token1.UserID() {
			t.Errorf(expectedFormat.StringButFoundString, token1.UserID(), token2.UserID())
		}
		if token2.IsExpired() {
			t.Errorf(expectedFormat.BoolButFoundBool, false, token2.IsExpired())
		}
	}
}

func Test_MongoDBStore_FindRefreshTokenWithCredential(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	token1 := Store.CreateRefreshToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	token2 := Store.FindRefreshTokenWithCredential(token1.ClientID(), token1.UserID())
	if token2 == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if token2.ClientID() != token1.ClientID() {
			t.Errorf(expectedFormat.StringButFoundString, token1.ClientID(), token2.ClientID())
		}
		if token2.UserID() != token1.UserID() {
			t.Errorf(expectedFormat.StringButFoundString, token1.UserID(), token2.UserID())
		}
		if token2.IsExpired() {
			t.Errorf(expectedFormat.BoolButFoundBool, false, token2.IsExpired())
		}
	}
}

func Test_MongoDBStore_DeleteRefreshToken(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	token1 := Store.CreateRefreshToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	Store.DeleteRefreshToken(token1)

	token2 := Store.FindRefreshTokenWithCredential(token1.ClientID(), token1.UserID())
	if token2 != nil {
		t.Errorf(expectedFormat.Nil)
	}
}
