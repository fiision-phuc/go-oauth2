package oauth2

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/phuc0302/go-mongo"
	"github.com/phuc0302/go-oauth2/oauth_table"
	"github.com/phuc0302/go-server"
	"github.com/phuc0302/go-server/expected_format"
	"github.com/phuc0302/go-server/string_format"
	"github.com/phuc0302/go-server/util"
)

func Test_TokenGrant_validateForm_MissingGrantType(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := server.CreateContext(w, r)
		defer server.Recovery(w, r)

		controller.HandleForm(context)
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", nil)
	status := util.ParseStatus(response)
	if status == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if status.Code != 400 {
			t.Errorf(expectedFormat.NumberButFoundNumber, 400, status.Code)
		}
		if status.Description != fmt.Sprintf(stringFormat.InvalidParameter, "grant_type") {
			t.Errorf(expectedFormat.InvalidParameter, "grant_type", status.Description)
		}
	}
}
func Test_TokenGrant_validateForm_MissingClientID(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := server.CreateContext(w, r)
		defer server.Recovery(w, r)

		controller.HandleForm(context)
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s", AuthorizationCodeGrant)))
	status := util.ParseStatus(response)
	if status == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if status.Code != 400 {
			t.Errorf(expectedFormat.NumberButFoundNumber, 400, status.Code)
		}
		if status.Description != fmt.Sprintf(stringFormat.InvalidParameter, "client_id") {
			t.Errorf(expectedFormat.InvalidParameter, "client_id", status.Description)
		}
	}

	// [Test 3] Missing client_secret
	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s", AuthorizationCodeGrant, u.ClientID.Hex())))
	status = util.ParseStatus(response)
	if status.Description != fmt.Sprintf(stringFormat.InvalidParameter, "client_secret") {
		t.Errorf(expectedFormat.InvalidParameter, "client_secret", status.Description)
	}
}
func Test_TokenGrant_validateForm_MissingClientSecret(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := server.CreateContext(w, r)
		defer server.Recovery(w, r)

		controller.HandleForm(context)
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s", AuthorizationCodeGrant, u.ClientID.Hex())))
	status := util.ParseStatus(response)
	if status == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if status.Code != 400 {
			t.Errorf(expectedFormat.NumberButFoundNumber, 400, status.Code)
		}
		if status.Description != fmt.Sprintf(stringFormat.InvalidParameter, "client_secret") {
			t.Errorf(expectedFormat.InvalidParameter, "client_secret", status.Description)
		}
	}
}

func Test_TokenGrant_passwordFlow_MissingUsername(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := server.CreateContext(w, r)
		defer server.Recovery(w, r)

		controller.HandleForm(context)
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s",
		PasswordGrant,
		u.ClientID.Hex(),
		u.ClientSecret.Hex(),
	)))
	status := util.ParseStatus(response)
	if status == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if status.Code != 400 {
			t.Errorf(expectedFormat.NumberButFoundNumber, 400, status.Code)
		}
		if status.Description != fmt.Sprintf(stringFormat.InvalidParameter, "username or password") {
			t.Errorf(expectedFormat.InvalidParameter, "username or password", status.Description)
		}
	}
}
func Test_TokenGrant_passwordFlow_MissingPassword(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := server.CreateContext(w, r)
		defer server.Recovery(w, r)

		controller.HandleForm(context)
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s",
		PasswordGrant,
		u.ClientID.Hex(),
		u.ClientSecret.Hex(),
		"admin",
	)))
	status := util.ParseStatus(response)
	if status == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if status.Code != 400 {
			t.Errorf(expectedFormat.NumberButFoundNumber, 400, status.Code)
		}
		if status.Description != fmt.Sprintf(stringFormat.InvalidParameter, "username or password") {
			t.Errorf(expectedFormat.InvalidParameter, "username or password", status.Description)
		}
	}
}
func Test_TokenGrant_passwordFlow_ValidParams(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := server.CreateContext(w, r)
		defer server.Recovery(w, r)

		controller.HandleForm(context)
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
		PasswordGrant,
		u.ClientID.Hex(),
		u.ClientSecret.Hex(),
		"admin",
		"Password",
	)))

	token1 := parseResult(response)
	if token1 == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		recordedAccessToken := Store.FindAccessToken(token1.AccessToken)
		recordedRefreshToken := Store.FindRefreshToken(token1.RefreshToken)
		if recordedAccessToken == nil {
			t.Error(expectedFormat.NotNil)
		}
		if recordedRefreshToken == nil {
			t.Error(expectedFormat.NotNil)
		}
	}
}

func Test_TokenGrant_NotAllowRefreshToken(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	// Modify config
	Cfg.GrantTypes = []string{
		AuthorizationCodeGrant,
		ClientCredentialsGrant,
		PasswordGrant,
	}
	Cfg.AllowRefreshToken = false
	grantsValidation = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(Cfg.GrantTypes, "|")))

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := server.CreateContext(w, r)
		defer server.Recovery(w, r)

		controller.HandleForm(context)
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
		PasswordGrant,
		u.ClientID.Hex(),
		u.ClientSecret.Hex(),
		u.Username,
		u.Password,
	)))

	token := parseResult(response)
	if token == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if len(token.RefreshToken) > 0 {
			t.Error(expectedFormat.Nil)
		}
	}
}

func Test_TokenGrant_refreshTokenFlow_MissingRefreshToken(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := server.CreateContext(w, r)
		defer server.Recovery(w, r)

		controller.HandleForm(context)
	}))
	defer ts.Close()

	// Send first request to get refresh token
	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
		PasswordGrant,
		u.ClientID.Hex(),
		u.ClientSecret.Hex(),
		u.Username,
		u.Password,
	)))

	// Send second request
	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s",
		RefreshTokenGrant,
		u.ClientID.Hex(),
		u.ClientSecret.Hex(),
	)))
	status := util.ParseStatus(response)
	if status == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if status.Code != 400 {
			t.Errorf(expectedFormat.NumberButFoundNumber, 400, status.Code)
		}
		if status.Description != fmt.Sprintf(stringFormat.InvalidParameter, "refresh_token") {
			t.Errorf(expectedFormat.InvalidParameter, "refresh_token", status.Description)
		}
	}
}

//func Test_TokenGrant_refreshTokenFlow_ExpiredRefreshToken(t *testing.T) {
//	u := new(TestUnit)
//	defer u.Teardown()
//	u.Setup()

//	// Setup server
//	controller := new(TokenGrant)
//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		context := server.CreateContext(w, r)
//		defer server.Recovery(w, r)

//		controller.HandleForm(context)
//	}))
//	defer ts.Close()

//	// Send first request to get refresh token
//	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
//		PasswordGrant,
//		clientID.Hex(),
//		clientSecret.Hex(),
//		username,
//		password,
//	)))
//	token1 := parseResult(response)

//	// Send second request
//	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&refresh_token=%s",
//		RefreshTokenGrant,
//		clientID.Hex(),
//		clientSecret.Hex(),
//		token1.RefreshToken,
//	)))
//	token2 := parseResult(response)
//	if token2 == nil {
//		t.Error(expectedFormat.NotNil)
//	} else {
//		if token2.AccessToken == token1.AccessToken {
//			t.Errorf("Expected new access_token but found \"%s\".", token2.AccessToken)
//		}

//		recordedAccessToken1 := TokenStore.FindAccessToken(token1.AccessToken)
//		recordedRefreshToken1 := TokenStore.FindRefreshToken(token1.RefreshToken)
//		if recordedAccessToken1 != nil {
//			t.Error(expectedFormat.Nil)
//		}
//		if recordedRefreshToken1 != nil {
//			t.Error(expectedFormat.Nil)
//		}

//		recordedAccessToken2 := TokenStore.FindAccessToken(token2.AccessToken)
//		recordedRefreshToken2 := TokenStore.FindRefreshToken(token2.RefreshToken)
//		if recordedAccessToken2 == nil {
//			t.Error(expectedFormat.NotNil)
//		}
//		if recordedRefreshToken2 == nil {
//			t.Error(expectedFormat.NotNil)
//		}
//	}
//}
func Test_TokenGrant_refreshTokenFlow_ValidParams(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := server.CreateContext(w, r)
		defer server.Recovery(w, r)

		controller.HandleForm(context)
	}))
	defer ts.Close()

	// Send first request to get refresh token
	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
		PasswordGrant,
		u.ClientID.Hex(),
		u.ClientSecret.Hex(),
		u.Username,
		u.Password,
	)))
	token1 := parseResult(response)

	// Send second request
	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&refresh_token=%s",
		RefreshTokenGrant,
		u.ClientID.Hex(),
		u.ClientSecret.Hex(),
		token1.RefreshToken,
	)))
	token2 := parseResult(response)
	if token2 == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if token2.AccessToken == token1.AccessToken {
			t.Errorf("Expected new access_token but found \"%s\".", token2.AccessToken)
		}

		accessToken1, _ := Store.FindAccessToken(token1.AccessToken).(*MongoDBToken)
		if err := mongo.EntityWithID(oauthTable.AccessToken, accessToken1.ID, new(MongoDBToken)); err == nil {
			t.Error(expectedFormat.Nil)
		}
		refreshToken1, _ := Store.FindRefreshToken(token1.RefreshToken).(*MongoDBToken)
		if err := mongo.EntityWithID(oauthTable.RefreshToken, refreshToken1.ID, new(MongoDBToken)); err == nil {
			t.Error(expectedFormat.Nil)
		}

		accessToken2, _ := Store.FindAccessToken(token2.AccessToken).(*MongoDBToken)
		if accessToken2 == nil {
			t.Error(expectedFormat.NotNil)
		} else {
			if accessToken2.ID == accessToken1.ID {
				t.Errorf("Expected new access_token but found \"%s\".", token2.AccessToken)
			}
		}
		refreshToken2, _ := Store.FindRefreshToken(token2.RefreshToken).(*MongoDBToken)
		if refreshToken2 == nil {
			t.Error(expectedFormat.NotNil)
		} else {
			if refreshToken2.ID == refreshToken1.ID {
				t.Errorf("Expected new refresh_token but found \"%s\".", token2.RefreshToken)
			}
		}
	}
}
