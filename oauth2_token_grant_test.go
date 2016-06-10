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
	//	defer os.Remove(ConfigFile)
	//	config := loadConfig()
	//	store := createStore()

	//	controller := CreateTokenGrant(config, store)
	//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		context := CreateRequest(r, w)
	//		controller.HandleForm(context)
	//	}))
	//	defer ts.Close()

	//	// Test missing username or password
	//	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s",
	//		PasswordGrant,
	//		clientID.Hex(),
	//		clientSecret.Hex(),
	//	)))
	//	status := parseError(response)
	//	if status.Description != fmt.Sprintf(templateError, "username or password") {
	//		t.Errorf(templateErrorMessage, "username or password", status.Description)
	//	}

	//	// Test invalid username or password
	//	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
	//		PasswordGrant,
	//		clientID.Hex(),
	//		clientSecret.Hex(),
	//		"admin1",
	//		"admin1",
	//	)))
	//	status = parseError(response)
	//	if status.Description != fmt.Sprintf(templateError, "username or password") {
	//		t.Errorf(templateErrorMessage, "username or password", status.Description)
	//	}

	//	// Test valid username and password
	//	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
	//		PasswordGrant,
	//		clientID.Hex(),
	//		clientSecret.Hex(),
	//		"admin2",
	//		"admin2",
	//	)))
	//	token1 := parseResult(response)
	//	if token1 == nil {
	//		t.Error("Expected not nil but found nil.")
	//	}
	//	if token1.AccessToken != store.accessTokens[1].Token() {
	//		t.Errorf("Expected %s but found %s", store.accessTokens[1].Token(), token1.AccessToken)
	//	}
	//	if token1.RefreshToken != store.refreshTokens[1].Token() {
	//		t.Errorf("Expected %s but found %s", store.refreshTokens[1].Token(), token1.RefreshToken)
	//	}

	//	// Test request second token should be the same as the first one
	//	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
	//		PasswordGrant,
	//		clientID.Hex(),
	//		clientSecret.Hex(),
	//		"admin2",
	//		"admin2",
	//	)))
	//	token2 := parseResult(response)
	//	if token2.AccessToken != token1.AccessToken {
	//		t.Errorf("Expected %s but found %s", token1.AccessToken, token2.AccessToken)
	//	}
	//	if token2.RefreshToken != token1.RefreshToken {
	//		t.Errorf("Expected %s but found %s", token1.RefreshToken, token2.RefreshToken)
	//	}

	//	// Test request expired token
	//	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
	//		PasswordGrant,
	//		clientID.Hex(),
	//		clientSecret.Hex(),
	//		username,
	//		password,
	//	)))
	//	token3 := parseResult(response)
	//	if len(store.accessTokens) != 2 {
	//		t.Errorf("Expected %d but found %d", 2, len(store.accessTokens))
	//	}
	//	if len(store.refreshTokens) != 2 {
	//		t.Errorf("Expected %d but found %d", 2, len(store.refreshTokens))
	//	}
	//	if token3.AccessToken != store.accessTokens[1].Token() {
	//		t.Errorf("Expected \"%s\" but found \"%s\".", store.accessTokens[1].Token(), token3.AccessToken)
	//	}
	//	if token3.RefreshToken != store.refreshTokens[1].Token() {
	//		t.Errorf("Expected \"%s\" but found \"%s\".", store.refreshTokens[1].Token(), token1.RefreshToken)
	//	}
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
	//	defer os.Remove(ConfigFile)
	//	config := loadConfig()
	//	store := createStore()

	//	controller := CreateTokenGrant(config, store)
	//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		context := CreateRequest(r, w)
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

	//	// Test missing refresh_token
	//	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s",
	//		RefreshTokenGrant,
	//		clientID.Hex(),
	//		clientSecret.Hex(),
	//	)))
	//	status := parseError(response)
	//	if status.Description != fmt.Sprintf(templateError, "refresh_token") {
	//		t.Errorf(templateErrorMessage, "refresh_token", status.Description)
	//	}

	//	// Send valid request
	//	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&refresh_token=%s",
	//		RefreshTokenGrant,
	//		clientID.Hex(),
	//		clientSecret.Hex(),
	//		token1.RefreshToken,
	//	)))
	//	token2 := parseResult(response)
	//	if token2.AccessToken == token1.AccessToken {
	//		t.Errorf("Expected new access token but found %s", token1.AccessToken)
	//	}
	//	if len(store.accessTokens) != 1 {
	//		t.Errorf("Expected %d but found %d", 1, len(store.accessTokens))
	//	}
}
func Test_RefreshGrantFlowWithExpiredRefreshToken(t *testing.T) {
	//	defer os.Remove(ConfigFile)
	//	config := loadConfig()
	//	store := createStore()

	//	controller := CreateTokenGrant(config, store)
	//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		context := CreateRequest(r, w)
	//		controller.HandleForm(context)
	//	}))
	//	defer ts.Close()

	//	// Send valid request
	//	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&refresh_token=%s",
	//		RefreshTokenGrant,
	//		clientID.Hex(),
	//		clientSecret.Hex(),
	//		store.refreshTokens[0].Token(),
	//	)))
	//	status := parseError(response)
	//	if status.Description != fmt.Sprintf("%s is expired.", "refresh_token") {
	//		t.Errorf("Expected \"%s is expired.\" but found \"%s\"", "refresh_token", status.Description)
	//	}
}
