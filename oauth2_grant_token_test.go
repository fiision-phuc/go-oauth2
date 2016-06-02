package oauth2

import "testing"

////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	templateError        = "Invalid %s parameter."
	templateErrorMessage = "Expected \"Invalid %s parameter.\" but found \"%s\""
)

func Test_TokenGrantGeneralValidation(t *testing.T) {
	//	defer os.Remove(ConfigFile)
	//	config := loadConfig()
	//	store := createStore()

	//	controller := CreateTokenGrant(config, store)
	//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		context := CreateRequest(r, w)
	//		controller.HandleForm(context)
	//	}))
	//	defer ts.Close()

	//	// Test missing grant_type
	//	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", nil)
	//	status := parseError(response)
	//	if status == nil {
	//		t.Error("Expected error return but found nil.")
	//	}
	//	if status.Code != 400 {
	//		t.Errorf("Expected error code 400 but found %d", status.Code)
	//	}
	//	if status.Description != fmt.Sprintf(templateError, "grant_type") {
	//		t.Errorf(templateErrorMessage, "grant_type", status.Description)
	//	}

	//	// Test missing client_id
	//	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s", AuthorizationCodeGrant)))
	//	status = parseError(response)
	//	if status.Description != fmt.Sprintf(templateError, "client_id") {
	//		t.Errorf(templateErrorMessage, "client_id", status.Description)
	//	}

	//	// Test missing client_secret
	//	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s", AuthorizationCodeGrant, clientID.Hex())))
	//	status = parseError(response)
	//	if status.Description != fmt.Sprintf(templateError, "client_secret") {
	//		t.Errorf(templateErrorMessage, "client_secret", status.Description)
	//	}
}

func Test_TokenGrantNotAllowRefreshToken(t *testing.T) {
	//	defer os.Remove(ConfigFile)
	//	config := loadConfig()
	//	store := createStore()

	//	// Modify config
	//	config.Grant = []string{
	//		AuthorizationCodeGrant,
	//		ClientCredentialsGrant,
	//		PasswordGrant,
	//	}
	//	config.AllowRefreshToken = false
	//	config.grantsValidation = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(config.Grant, "|")))

	//	controller := CreateTokenGrant(config, store)
	//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		context := CreateRequest(r, w)
	//		controller.HandleForm(context)
	//	}))
	//	defer ts.Close()

	//	// Test invalid grant_type
	//	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s", RefreshTokenGrant, clientID.Hex(), clientSecret.Hex())))
	//	status := parseError(response)
	//	if status.Description != fmt.Sprintf(templateError, "grant_type") {
	//		t.Errorf(templateErrorMessage, "grant_type", status.Description)
	//	}

	//	// Test valid request token
	//	response, _ = http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
	//		PasswordGrant,
	//		clientID.Hex(),
	//		clientSecret.Hex(),
	//		username,
	//		password,
	//	)))
	//	token := parseResult(response)
	//	if token == nil {
	//		t.Error("Expected not nil but found nil.")
	//	}
	//	if token.RefreshToken != "" {
	//		t.Errorf("Expected nil refresh token but found %s.", token.RefreshToken)
	//	}
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
