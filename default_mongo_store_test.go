package oauth2

import (
	"reflect"
	"testing"
	"time"

	"github.com/phuc0302/go-oauth2/test"
)

func Test_DefaultMongoStore_FindUserWithID(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	user := TokenStore.FindUserWithID(u.UserID.Hex())
	if user == nil {
		t.Errorf(test.ExpectedStringButFoundString, u.User1, user)
	} else {
		if user.UserID() != u.UserID.Hex() {
			t.Errorf(test.ExpectedStringButFoundString, u.UserID.Hex(), user.UserID())
		}
		if user.Username() != "admin" {
			t.Errorf(test.ExpectedStringButFoundString, "admin", user.Username())
		}
		if !reflect.DeepEqual(user.UserRoles(), []string{"r_user", "r_admin"}) {
			t.Errorf(test.ExpectedStringButFoundString, []string{"r_user", "r_admin"}, user.UserRoles())
		}
	}
}

func Test_DefaultMongoStore_FindUserWithClient(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	user := TokenStore.FindUserWithClient(u.ClientID.Hex(), u.ClientSecret.Hex())
	if user == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if user.UserID() != u.ClientID.Hex() {
			t.Errorf(test.ExpectedStringButFoundString, u.ClientID.Hex(), user.UserID())
		}
		if user.Username() != u.ClientID.Hex() {
			t.Errorf(test.ExpectedStringButFoundString, u.ClientID.Hex(), user.Username())
		}
		if !reflect.DeepEqual(user.UserRoles(), []string{"r_device"}) {
			t.Errorf(test.ExpectedStringButFoundString, []string{"r_device"}, user.UserRoles())
		}
	}
}

func Test_DefaultMongoStore_FindUserWithCredential(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	user := TokenStore.FindUserWithCredential("admin", "admin")
	if user == nil {
		t.Errorf(test.ExpectedStringButFoundString, u.User1, user)
	} else {
		if user.UserID() != u.UserID.Hex() {
			t.Errorf(test.ExpectedStringButFoundString, u.UserID.Hex(), user.UserID())
		}
		if user.Username() != "admin" {
			t.Errorf(test.ExpectedStringButFoundString, "admin", user.Username())
		}
		if !reflect.DeepEqual(user.UserRoles(), []string{"r_user", "r_admin"}) {
			t.Errorf(test.ExpectedStringButFoundString, []string{"r_user", "r_admin"}, user.UserRoles())
		}
	}
}

func Test_DefaultMongoStore_FindClientWithID(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	client := TokenStore.FindClientWithID(u.ClientID.Hex())
	if client == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if client.ClientID() != u.ClientID.Hex() {
			t.Errorf(test.ExpectedStringButFoundString, u.ClientID.Hex(), client.ClientID())
		}
		if client.ClientSecret() != u.ClientSecret.Hex() {
			t.Errorf(test.ExpectedStringButFoundString, u.ClientSecret.Hex(), client.ClientSecret())
		}
		if !reflect.DeepEqual(client.GrantTypes(), []string{AuthorizationCodeGrant, PasswordGrant, RefreshTokenGrant}) {
			t.Errorf(test.ExpectedStringButFoundString, []string{AuthorizationCodeGrant, PasswordGrant, RefreshTokenGrant}, client.GrantTypes())
		}
		if !reflect.DeepEqual(client.RedirectURIs(), []string{"http://www.sample01.com", "http://www.sample02.com"}) {
			t.Errorf(test.ExpectedStringButFoundString, []string{"http://www.sample01.com", "http://www.sample02.com"}, client.RedirectURIs())
		}
	}
}

func Test_DefaultMongoStore_FindClientWithCredential(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	client := TokenStore.FindClientWithCredential(u.ClientID.Hex(), u.ClientSecret.Hex())
	if client == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if client.ClientID() != u.ClientID.Hex() {
			t.Errorf(test.ExpectedStringButFoundString, u.ClientID.Hex(), client.ClientID())
		}
		if client.ClientSecret() != u.ClientSecret.Hex() {
			t.Errorf(test.ExpectedStringButFoundString, u.ClientSecret.Hex(), client.ClientSecret())
		}
		if !reflect.DeepEqual(client.GrantTypes(), []string{AuthorizationCodeGrant, PasswordGrant, RefreshTokenGrant}) {
			t.Errorf(test.ExpectedStringButFoundString, []string{AuthorizationCodeGrant, PasswordGrant, RefreshTokenGrant}, client.GrantTypes())
		}
		if !reflect.DeepEqual(client.RedirectURIs(), []string{"http://www.sample01.com", "http://www.sample02.com"}) {
			t.Errorf(test.ExpectedStringButFoundString, []string{"http://www.sample01.com", "http://www.sample02.com"}, client.RedirectURIs())
		}
	}
}

func Test_DefaultMongoStore_CreateAccessToken(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	token := TokenStore.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	if token == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if token.ClientID() != u.ClientID.Hex() {
			t.Errorf(test.ExpectedStringButFoundString, u.ClientID.Hex(), token.ClientID())
		}
		if token.UserID() != u.UserID.Hex() {
			t.Errorf(test.ExpectedStringButFoundString, u.UserID.Hex(), token.ClientID())
		}
		if token.IsExpired() {
			t.Errorf(test.ExpectedBoolButFoundBool, false, token.IsExpired())
		}
	}
}

func Test_DefaultMongoStore_FindAccessToken(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	token1 := TokenStore.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	token2 := TokenStore.FindAccessToken(token1.Token())
	if token2 == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if token2.ClientID() != token1.ClientID() {
			t.Errorf(test.ExpectedStringButFoundString, token1.ClientID(), token2.ClientID())
		}
		if token2.UserID() != token1.UserID() {
			t.Errorf(test.ExpectedStringButFoundString, token1.UserID(), token2.UserID())
		}
		if token2.IsExpired() {
			t.Errorf(test.ExpectedBoolButFoundBool, false, token2.IsExpired())
		}
	}
}

func Test_DefaultMongoStore_FindAccessTokenWithCredential(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	token1 := TokenStore.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	token2 := TokenStore.FindAccessTokenWithCredential(u.ClientID.Hex(), u.UserID.Hex())
	if token2 == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if token2.ClientID() != token1.ClientID() {
			t.Errorf(test.ExpectedStringButFoundString, token1.ClientID(), token2.ClientID())
		}
		if token2.UserID() != token1.UserID() {
			t.Errorf(test.ExpectedStringButFoundString, token1.UserID(), token2.UserID())
		}
		if token2.IsExpired() {
			t.Errorf(test.ExpectedBoolButFoundBool, false, token2.IsExpired())
		}
	}
}

func Test_DefaultMongoStore_DeleteAccessToken(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	token1 := TokenStore.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	TokenStore.DeleteAccessToken(token1)

	token2 := TokenStore.FindAccessTokenWithCredential(token1.ClientID(), token1.UserID())
	if token2 != nil {
		t.Errorf(test.ExpectedNil)
	}
}

func Test_DefaultMongoStore_CreateRefreshToken(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	token := TokenStore.CreateRefreshToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	if token == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if token.ClientID() != u.ClientID.Hex() {
			t.Errorf(test.ExpectedStringButFoundString, u.ClientID.Hex(), token.ClientID())
		}
		if token.UserID() != u.UserID.Hex() {
			t.Errorf(test.ExpectedStringButFoundString, u.UserID.Hex(), token.ClientID())
		}
		if token.IsExpired() {
			t.Errorf(test.ExpectedBoolButFoundBool, false, token.IsExpired())
		}
	}
}

func Test_DefaultMongoStore_FindRefreshToken(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	token1 := TokenStore.CreateRefreshToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	token2 := TokenStore.FindRefreshToken(token1.Token())
	if token2 == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if token2.ClientID() != token1.ClientID() {
			t.Errorf(test.ExpectedStringButFoundString, token1.ClientID(), token2.ClientID())
		}
		if token2.UserID() != token1.UserID() {
			t.Errorf(test.ExpectedStringButFoundString, token1.UserID(), token2.UserID())
		}
		if token2.IsExpired() {
			t.Errorf(test.ExpectedBoolButFoundBool, false, token2.IsExpired())
		}
	}
}

func Test_DefaultMongoStore_FindRefreshTokenWithCredential(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	token1 := TokenStore.CreateRefreshToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	token2 := TokenStore.FindRefreshTokenWithCredential(token1.ClientID(), token1.UserID())
	if token2 == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if token2.ClientID() != token1.ClientID() {
			t.Errorf(test.ExpectedStringButFoundString, token1.ClientID(), token2.ClientID())
		}
		if token2.UserID() != token1.UserID() {
			t.Errorf(test.ExpectedStringButFoundString, token1.UserID(), token2.UserID())
		}
		if token2.IsExpired() {
			t.Errorf(test.ExpectedBoolButFoundBool, false, token2.IsExpired())
		}
	}
}

func Test_DefaultMongoStore_DeleteRefreshToken(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	token1 := TokenStore.CreateRefreshToken(u.ClientID.Hex(), u.UserID.Hex(), time.Now(), time.Now().Add(Cfg.AccessTokenDuration))
	TokenStore.DeleteRefreshToken(token1)

	token2 := TokenStore.FindRefreshTokenWithCredential(token1.ClientID(), token1.UserID())
	if token2 != nil {
		t.Errorf(test.ExpectedNil)
	}
}
