package oauth2

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/phuc0302/go-oauth2/test"
)

func Test_TokenGrantGeneralValidation(t *testing.T) {
	defer teardown()
	setup()

	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		controller.HandleForm(context)
	}))
	defer ts.Close()

	// [Test 1] Missing grant_type
	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", nil)
	status := parseError(response)
	if status == nil {
		t.Error(test.ExpectedNotNil)
	}
	if status.Code != 400 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 400, status.Code)
	}
	if status.Description != fmt.Sprintf("Invalid %s parameter.", "grant_type") {
		t.Errorf(test.ExpectedInvalidParameter, "grant_type", status.Description)
	}

	// [Test 2] Missing client_id
	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s", AuthorizationCodeGrant)))
	status = parseError(response)
	if status.Description != fmt.Sprintf("Invalid %s parameter.", "client_id") {
		t.Errorf(test.ExpectedInvalidParameter, "client_id", status.Description)
	}

	// [Test 3] Missing client_secret
	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s", AuthorizationCodeGrant, clientID.Hex())))
	status = parseError(response)
	if status.Description != fmt.Sprintf("Invalid %s parameter.", "client_secret") {
		t.Errorf(test.ExpectedInvalidParameter, "client_secret", status.Description)
	}
}

func Test_PasswordGrantFlow(t *testing.T) {
	defer teardown()
	setup()

	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		controller.HandleForm(context)
	}))
	defer ts.Close()

	// [Test 1] Missing username
	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s",
		PasswordGrant,
		clientID.Hex(),
		clientSecret.Hex(),
	)))
	status := parseError(response)
	if status.Description != fmt.Sprintf("Invalid %s parameter.", "username or password") {
		t.Errorf(test.ExpectedInvalidParameter, "username or password", status.Description)
	}

	// [Test 2] Missing password
	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s",
		PasswordGrant,
		clientID.Hex(),
		clientSecret.Hex(),
		"admin",
	)))
	status = parseError(response)
	if status.Description != fmt.Sprintf("Invalid %s parameter.", "username or password") {
		t.Errorf(test.ExpectedInvalidParameter, "username or password", status.Description)
	}

	// [Test 3] valid username and password
	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
		PasswordGrant,
		clientID.Hex(),
		clientSecret.Hex(),
		"admin",
		"admin",
	)))
	token1 := parseResult(response)

	recordedAccessToken := tokenStore.FindAccessToken(token1.AccessToken)
	recordedRefreshToken := tokenStore.FindRefreshToken(token1.RefreshToken)
	if token1 == nil {
		t.Error(test.ExpectedNotNil)
	}
	if recordedAccessToken == nil {
		t.Error(test.ExpectedNotNil)
	}
	if recordedRefreshToken == nil {
		t.Error(test.ExpectedNotNil)
	}
}

func Test_TokenGrantNotAllowRefreshToken(t *testing.T) {
	defer teardown()
	setup()

	// Modify config
	cfg.GrantTypes = []string{
		AuthorizationCodeGrant,
		ClientCredentialsGrant,
		PasswordGrant,
	}
	cfg.AllowRefreshToken = false
	grantsValidation = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(cfg.GrantTypes, "|")))

	// Setup test server
	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		controller.HandleForm(context)
	}))
	defer ts.Close()

	// [Test 1] Invalid grant_type
	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s", RefreshTokenGrant, clientID.Hex(), clientSecret.Hex())))
	status := parseError(response)
	if status.Description != fmt.Sprintf("Invalid %s parameter.", "grant_type") {
		t.Errorf(test.ExpectedInvalidParameter, "grant_type", status.Description)
	}

	// [Test 2] Valid request token
	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
		PasswordGrant,
		clientID.Hex(),
		clientSecret.Hex(),
		username,
		password,
	)))
	token := parseResult(response)
	if token == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if len(token.RefreshToken) > 0 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 0, len(token.RefreshToken))
		}
	}
}

func Test_RefreshGrantFlow(t *testing.T) {
	defer teardown()
	setup()

	controller := new(TokenGrant)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		controller.HandleForm(context)
	}))
	defer ts.Close()

	// Send first request to get refresh token
	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
		PasswordGrant,
		clientID.Hex(),
		clientSecret.Hex(),
		username,
		password,
	)))
	token1 := parseResult(response)

	// [Test 1] Missing refresh_token
	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s",
		RefreshTokenGrant,
		clientID.Hex(),
		clientSecret.Hex(),
	)))
	status := parseError(response)
	if status.Description != fmt.Sprintf("Invalid %s parameter.", "refresh_token") {
		t.Errorf(test.ExpectedInvalidParameter, "refresh_token", status.Description)
	}

	// Send valid request
	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&refresh_token=%s",
		RefreshTokenGrant,
		clientID.Hex(),
		clientSecret.Hex(),
		token1.RefreshToken,
	)))
	token2 := parseResult(response)
	if token2.AccessToken == token1.AccessToken {
		t.Errorf("Expected new access_token but found \"%s\".", token1.AccessToken)
	}
}
