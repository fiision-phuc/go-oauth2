package oauth2

import (
	"testing"
	"time"

	"github.com/phuc0302/go-oauth2/test"
)

func Test_DefaultMongoStore(t *testing.T) {
	defer teardown()
	setup()

	// [Test 1] FindUserWithID
	recordUser1 := tokenStore.FindUserWithID(userID.Hex())
	if recordUser1 == nil {
		t.Errorf(test.ExpectedStringButFoundString, user1, recordUser1)
	}

	// [Test 2] FindUserWithClient
	recordUser2 := tokenStore.FindUserWithClient(clientID.Hex(), clientSecret.Hex())
	if recordUser2 == nil {
		t.Errorf(test.ExpectedStringButFoundString, user2, recordUser2)
	}

	// [Test 3] FindUserWithCredential
	recordUser3 := tokenStore.FindUserWithCredential("admin", "admin")
	if recordUser3 == nil {
		t.Errorf(test.ExpectedStringButFoundString, user1, recordUser1)
	}

	// [Test 4] FindClientWithID
	recordClient1 := tokenStore.FindClientWithID(clientID.Hex())
	if recordClient1 == nil {
		t.Errorf(test.ExpectedStringButFoundString, client1, recordClient1)
	}

	// [Test 5] FindClientWithCredential
	recordClient2 := tokenStore.FindClientWithCredential(clientID.Hex(), clientSecret.Hex())
	if recordClient2 == nil {
		t.Errorf(test.ExpectedStringButFoundString, client1, recordClient2)
	}

	// [Test 6] CreateAccessToken
	token1 := tokenStore.CreateAccessToken(clientID.Hex(), userID.Hex(), time.Now(), time.Now().Add(3600))
	if token1 == nil {
		t.Error(test.ExpectedNotNil)
	}
	if token1.ClientID() != clientID.Hex() {
		t.Errorf(test.ExpectedStringButFoundString, clientID.Hex(), token1.ClientID())
	}
	if token1.UserID() != userID.Hex() {
		t.Errorf(test.ExpectedStringButFoundString, userID.Hex(), token1.ClientID())
	}

	// [Test 7] FindAccessToken
	token2 := tokenStore.FindAccessToken(token1.Token())
	if token2 == nil {
		t.Error(test.ExpectedNotNil)
	}

	// [Test 8] FindAccessTokenWithCredential
	token3 := tokenStore.FindAccessTokenWithCredential(token1.ClientID(), token1.UserID())
	if token3 == nil {
		t.Error(test.ExpectedNotNil)
	}

	// [Test 9] DeleteAccessToken
	tokenStore.DeleteAccessToken(token1)
	token4 := tokenStore.FindAccessTokenWithCredential(token1.ClientID(), token1.UserID())
	if token4 != nil {
		t.Errorf(test.ExpectedNil)
	}

	// [Test 10] CreateRefreshToken
	refreshToken1 := tokenStore.CreateRefreshToken(clientID.Hex(), userID.Hex(), time.Now(), time.Now().Add(3600))
	if refreshToken1 == nil {
		t.Errorf(test.ExpectedNotNil)
	}
	if refreshToken1.ClientID() != clientID.Hex() {
		t.Errorf(test.ExpectedStringButFoundString, clientID.Hex(), refreshToken1.ClientID())
	}
	if refreshToken1.UserID() != userID.Hex() {
		t.Errorf(test.ExpectedStringButFoundString, userID.Hex(), refreshToken1.ClientID())
	}

	// [Test 11] FindRefreshToken
	refreshToken2 := tokenStore.FindRefreshToken(refreshToken1.Token())
	if refreshToken2 == nil {
		t.Errorf(test.ExpectedNotNil)
	}

	// [Test 12] FindRefreshTokenWithCredential
	refreshToken3 := tokenStore.FindRefreshTokenWithCredential(refreshToken1.ClientID(), refreshToken1.UserID())
	if refreshToken3 == nil {
		t.Errorf(test.ExpectedNotNil)
	}

	// [Test 13] DeleteRefreshToken
	tokenStore.DeleteRefreshToken(refreshToken1)
	refreshToken4 := tokenStore.FindRefreshTokenWithCredential(refreshToken1.ClientID(), refreshToken1.UserID())
	if refreshToken4 != nil {
		t.Errorf(test.ExpectedNil)
	}
}
