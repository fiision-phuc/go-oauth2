package oauth2

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/phuc0302/go-mongo"
	"github.com/phuc0302/go-oauth2/test"
	"github.com/phuc0302/go-oauth2/util"
)

func Test_TokenGrant_validateForm_MissingGrantType(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := createRequestContext(r, w)
		defer recovery(context, true)

		controller.HandleForm(context, nil)
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", nil)
	status := util.ParseStatus(response)
	if status == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if status.Code != 400 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 400, status.Code)
		}
		if status.Description != fmt.Sprintf(InvalidParameter, "grant_type") {
			t.Errorf(test.ExpectedInvalidParameter, "grant_type", status.Description)
		}
	}
}
func Test_TokenGrant_validateForm_MissingClientID(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := createRequestContext(r, w)
		defer recovery(context, true)

		controller.HandleForm(context, nil)
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s", AuthorizationCodeGrant)))
	status := util.ParseStatus(response)
	if status == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if status.Code != 400 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 400, status.Code)
		}
		if status.Description != fmt.Sprintf(InvalidParameter, "client_id") {
			t.Errorf(test.ExpectedInvalidParameter, "client_id", status.Description)
		}
	}

	// [Test 3] Missing client_secret
	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s", AuthorizationCodeGrant, u.ClientID.Hex())))
	status = util.ParseStatus(response)
	if status.Description != fmt.Sprintf(InvalidParameter, "client_secret") {
		t.Errorf(test.ExpectedInvalidParameter, "client_secret", status.Description)
	}
}
func Test_TokenGrant_validateForm_MissingClientSecret(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := createRequestContext(r, w)
		defer recovery(context, true)

		controller.HandleForm(context, nil)
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s", AuthorizationCodeGrant, u.ClientID.Hex())))
	status := util.ParseStatus(response)
	if status == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if status.Code != 400 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 400, status.Code)
		}
		if status.Description != fmt.Sprintf(InvalidParameter, "client_secret") {
			t.Errorf(test.ExpectedInvalidParameter, "client_secret", status.Description)
		}
	}
}

func Test_TokenGrant_passwordFlow_MissingUsername(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := createRequestContext(r, w)
		defer recovery(context, true)

		controller.HandleForm(context, nil)
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s",
		PasswordGrant,
		u.ClientID.Hex(),
		u.ClientSecret.Hex(),
	)))
	status := util.ParseStatus(response)
	if status == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if status.Code != 400 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 400, status.Code)
		}
		if status.Description != fmt.Sprintf(InvalidParameter, "username or password") {
			t.Errorf(test.ExpectedInvalidParameter, "username or password", status.Description)
		}
	}
}
func Test_TokenGrant_passwordFlow_MissingPassword(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := createRequestContext(r, w)
		defer recovery(context, true)

		controller.HandleForm(context, nil)
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
		t.Error(test.ExpectedNotNil)
	} else {
		if status.Code != 400 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 400, status.Code)
		}
		if status.Description != fmt.Sprintf(InvalidParameter, "username or password") {
			t.Errorf(test.ExpectedInvalidParameter, "username or password", status.Description)
		}
	}
}
func Test_TokenGrant_passwordFlow_ValidParams(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := createRequestContext(r, w)
		defer recovery(context, true)

		controller.HandleForm(context, nil)
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
		PasswordGrant,
		u.ClientID.Hex(),
		u.ClientSecret.Hex(),
		"admin",
		"admin",
	)))

	token1 := parseResult(response)
	if token1 == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		recordedAccessToken := store.FindAccessToken(token1.AccessToken)
		recordedRefreshToken := store.FindRefreshToken(token1.RefreshToken)
		if recordedAccessToken == nil {
			t.Error(test.ExpectedNotNil)
		}
		if recordedRefreshToken == nil {
			t.Error(test.ExpectedNotNil)
		}
	}
}

func Test_TokenGrant_NotAllowRefreshToken(t *testing.T) {
	u := new(TestUnit)
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
		context := createRequestContext(r, w)
		defer recovery(context, true)

		controller.HandleForm(context, nil)
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
		t.Error(test.ExpectedNotNil)
	} else {
		if len(token.RefreshToken) > 0 {
			t.Error(test.ExpectedNil)
		}
	}
}

func Test_TokenGrant_refreshTokenFlow_MissingRefreshToken(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := createRequestContext(r, w)
		defer recovery(context, true)

		controller.HandleForm(context, nil)
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
		t.Error(test.ExpectedNotNil)
	} else {
		if status.Code != 400 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 400, status.Code)
		}
		if status.Description != fmt.Sprintf(InvalidParameter, "refresh_token") {
			t.Errorf(test.ExpectedInvalidParameter, "refresh_token", status.Description)
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
//		context := createRequestContext(r, w)
//		defer recovery(context, true)

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
//		t.Error(test.ExpectedNotNil)
//	} else {
//		if token2.AccessToken == token1.AccessToken {
//			t.Errorf("Expected new access_token but found \"%s\".", token2.AccessToken)
//		}

//		recordedAccessToken1 := TokenStore.FindAccessToken(token1.AccessToken)
//		recordedRefreshToken1 := TokenStore.FindRefreshToken(token1.RefreshToken)
//		if recordedAccessToken1 != nil {
//			t.Error(test.ExpectedNil)
//		}
//		if recordedRefreshToken1 != nil {
//			t.Error(test.ExpectedNil)
//		}

//		recordedAccessToken2 := TokenStore.FindAccessToken(token2.AccessToken)
//		recordedRefreshToken2 := TokenStore.FindRefreshToken(token2.RefreshToken)
//		if recordedAccessToken2 == nil {
//			t.Error(test.ExpectedNotNil)
//		}
//		if recordedRefreshToken2 == nil {
//			t.Error(test.ExpectedNotNil)
//		}
//	}
//}
func Test_TokenGrant_refreshTokenFlow_ValidParams(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := createRequestContext(r, w)
		defer recovery(context, true)

		controller.HandleForm(context, nil)
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
		t.Error(test.ExpectedNotNil)
	} else {
		if token2.AccessToken == token1.AccessToken {
			t.Errorf("Expected new access_token but found \"%s\".", token2.AccessToken)
		}

		accessToken1, _ := store.FindAccessToken(token1.AccessToken).(*DefaultToken)
		if err := mongo.EntityWithID(TableAccessToken, accessToken1.ID, new(DefaultToken)); err == nil {
			t.Error(test.ExpectedNil)
		}
		refreshToken1, _ := store.FindRefreshToken(token1.RefreshToken).(*DefaultToken)
		if err := mongo.EntityWithID(TableRefreshToken, refreshToken1.ID, new(DefaultToken)); err == nil {
			t.Error(test.ExpectedNil)
		}

		accessToken2, _ := store.FindAccessToken(token2.AccessToken).(*DefaultToken)
		if accessToken2 == nil {
			t.Error(test.ExpectedNotNil)
		} else {
			if accessToken2.ID == accessToken1.ID {
				t.Errorf("Expected new access_token but found \"%s\".", token2.AccessToken)
			}
		}
		refreshToken2, _ := store.FindRefreshToken(token2.RefreshToken).(*DefaultToken)
		if refreshToken2 == nil {
			t.Error(test.ExpectedNotNil)
		} else {
			if refreshToken2.ID == refreshToken1.ID {
				t.Errorf("Expected new refresh_token but found \"%s\".", token2.RefreshToken)
			}
		}
	}
}
