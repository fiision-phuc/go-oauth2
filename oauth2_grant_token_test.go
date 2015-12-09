package oauth2

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

func Test_GeneralValidation(t *testing.T) {
	defer os.Remove(ConfigFile)
	store := createStore()
	config := LoadConfigs()
	controller := CreateTokenGrant(config, store)
	templateError := "Invalid %s parameter."
	templateErrorMessage := "Expected \"Invalid %s parameter.\" but found \"%s\""

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := CreateRequestContext(r, w)
		controller.HandleForm(context)
	}))
	defer ts.Close()

	// Test missing grant_type
	response, _ := http.PostForm(ts.URL, url.Values{})
	status := parseError(response)
	if status == nil {
		t.Error("Expected error return but found nil.")
	}
	if status.Code != 400 {
		t.Errorf("Expected error code 400 but found %d", status.Code)
	}
	if status.Description != fmt.Sprintf(templateError, "grant_type") {
		t.Errorf(templateErrorMessage, "grant_type", status.Description)
	}

	// Test missing client_id
	response, _ = http.PostForm(ts.URL, url.Values{
		"grant_type": []string{AuthorizationCodeGrant},
	})
	status = parseError(response)
	if status.Description != fmt.Sprintf(templateError, "client_id") {
		t.Errorf(templateErrorMessage, "client_id", status.Description)
	}

	// Test missing client_secret
	response, _ = http.PostForm(ts.URL, url.Values{
		"grant_type": []string{AuthorizationCodeGrant},
		"client_id":  []string{store.clients[0].GetClientID()},
	})
	status = parseError(response)
	if status.Description != fmt.Sprintf(templateError, "client_secret") {
		t.Errorf(templateErrorMessage, "client_secret", status.Description)
	}
}

func Test_NotAllowRefreshGrantFlow(t *testing.T) {
	defer os.Remove(ConfigFile)
	store := createStore()
	config := LoadConfigs()

	// Modify config
	config.Grant = []string{AuthorizationCodeGrant, ClientCredentialsGrant, PasswordGrant}
	config.allowRefreshToken = false
	config.grantsValidation = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(config.Grant, "|")))

	controller := CreateTokenGrant(config, store)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := CreateRequestContext(r, w)
		controller.HandleForm(context)
	}))
	defer ts.Close()

	// Test invalid grant_type
	response, _ := http.PostForm(ts.URL, url.Values{
		"grant_type":    []string{RefreshTokenGrant},
		"client_id":     []string{store.clients[0].GetClientID()},
		"client_secret": []string{store.clients[0].GetClientSecret()},
	})
	status := parseError(response)
	if status.Description != fmt.Sprintf("Invalid %s parameter.", "grant_type") {
		t.Errorf("Expected \"Invalid %s parameter.\" but found \"%s\"", "grant_type", status.Description)
	}

	// Test valid request token
	response, _ = http.PostForm(ts.URL, url.Values{
		"grant_type":    []string{PasswordGrant},
		"client_id":     []string{store.clients[0].GetClientID()},
		"client_secret": []string{store.clients[0].GetClientSecret()},
		"username":      []string{"admin"},
		"password":      []string{"admin"},
	})
	token := parseResult(response)
	if token == nil {
		t.Error("Expected not nil but found nil.")
	}
	if token.RefreshToken != "" {
		t.Errorf("Expected nil refresh token but found %s.", token.RefreshToken)
	}
}

func Test_PasswordGrantFlow(t *testing.T) {
	defer os.Remove(ConfigFile)
	store := createStore()
	config := LoadConfigs()
	controller := CreateTokenGrant(config, store)
	templateError := "Invalid %s parameter."
	templateErrorMessage := "Expected \"Invalid %s parameter.\" but found \"%s\""

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := CreateRequestContext(r, w)
		controller.HandleForm(context)
	}))
	defer ts.Close()

	// Test missing username or password
	response, _ := http.PostForm(ts.URL, url.Values{
		"grant_type":    []string{PasswordGrant},
		"client_id":     []string{store.clients[0].GetClientID()},
		"client_secret": []string{store.clients[0].GetClientSecret()},
	})
	status := parseError(response)
	if status.Description != fmt.Sprintf(templateError, "username or password") {
		t.Errorf(templateErrorMessage, "username or password", status.Description)
	}

	// Test invalid username or password
	response, _ = http.PostForm(ts.URL, url.Values{
		"grant_type":    []string{PasswordGrant},
		"client_id":     []string{store.clients[0].GetClientID()},
		"client_secret": []string{store.clients[0].GetClientSecret()},
		"username":      []string{"admin1"},
		"password":      []string{"admin1"},
	})
	status = parseError(response)
	if status.Description != fmt.Sprintf(templateError, "username or password") {
		t.Errorf(templateErrorMessage, "username or password", status.Description)
	}

	// Test valid username and password
	response, _ = http.PostForm(ts.URL, url.Values{
		"grant_type":    []string{PasswordGrant},
		"client_id":     []string{store.clients[0].GetClientID()},
		"client_secret": []string{store.clients[0].GetClientSecret()},
		"username":      []string{"admin"},
		"password":      []string{"admin"},
	})
	token1 := parseResult(response)
	if token1 == nil {
		t.Error("Expected not nil but found nil.")
	}
	if token1.AccessToken != store.accessTokens[1].GetToken() {
		t.Errorf("Expected %s but found %s", store.accessTokens[1].GetToken(), token1.AccessToken)
	}
	if token1.RefreshToken != store.refreshTokens[1].GetToken() {
		t.Errorf("Expected %s but found %s", store.refreshTokens[1].GetToken(), token1.RefreshToken)
	}

	// Test request second token should be the same as the first one
	response, _ = http.PostForm(ts.URL, url.Values{
		"grant_type":    []string{PasswordGrant},
		"client_id":     []string{store.clients[0].GetClientID()},
		"client_secret": []string{store.clients[0].GetClientSecret()},
		"username":      []string{"admin"},
		"password":      []string{"admin"},
	})
	token2 := parseResult(response)
	if token2.AccessToken != token1.AccessToken {
		t.Errorf("Expected %s but found %s", token1.AccessToken, token2.AccessToken)
	}
	if token2.RefreshToken != token1.RefreshToken {
		t.Errorf("Expected %s but found %s", token1.RefreshToken, token2.RefreshToken)
	}

	// Test request existing token should be deleted
	response, _ = http.PostForm(ts.URL, url.Values{
		"grant_type":    []string{PasswordGrant},
		"client_id":     []string{store.clients[0].GetClientID()},
		"client_secret": []string{store.clients[0].GetClientSecret()},
		"username":      []string{"admin2"},
		"password":      []string{"admin2"},
	})
	token3 := parseResult(response)
	if token3.AccessToken == store.accessTokens[0].GetToken() {
		t.Errorf("Expected %s but found %s", store.accessTokens[2].GetToken(), token3.AccessToken)
	}
	if token3.RefreshToken == store.refreshTokens[0].GetToken() {
		t.Errorf("Expected %s but found %s", store.refreshTokens[2].GetToken(), token1.RefreshToken)
	}
	if token3.AccessToken != store.accessTokens[2].GetToken() {
		t.Errorf("Expected %s but found %s", store.accessTokens[2].GetToken(), token3.AccessToken)
	}
	if token3.RefreshToken != store.refreshTokens[2].GetToken() {
		t.Errorf("Expected %s but found %s", store.refreshTokens[2].GetToken(), token1.RefreshToken)
	}
}

func Test_RefreshGrantFlow(t *testing.T) {
	defer os.Remove(ConfigFile)
	store := createStore()
	config := LoadConfigs()
	controller := CreateTokenGrant(config, store)
	templateError := "Invalid %s parameter."
	templateErrorMessage := "Expected \"Invalid %s parameter.\" but found \"%s\""

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := CreateRequestContext(r, w)
		controller.HandleForm(context)
	}))
	defer ts.Close()

	// Send first request to get refresh token
	response, _ := http.PostForm(ts.URL, url.Values{
		"grant_type":    []string{PasswordGrant},
		"client_id":     []string{store.clients[0].GetClientID()},
		"client_secret": []string{store.clients[0].GetClientSecret()},
		"username":      []string{"admin"},
		"password":      []string{"admin"},
	})
	token1 := parseResult(response)

	// Test missing refresh_token
	response, _ = http.PostForm(ts.URL, url.Values{
		"grant_type":    []string{RefreshTokenGrant},
		"client_id":     []string{store.clients[0].GetClientID()},
		"client_secret": []string{store.clients[0].GetClientSecret()},
	})
	status := parseError(response)
	if status.Description != fmt.Sprintf(templateError, "refresh_token") {
		t.Errorf(templateErrorMessage, "refresh_token", status.Description)
	}

	// Send valid request
	response, _ = http.PostForm(ts.URL, url.Values{
		"grant_type":    []string{RefreshTokenGrant},
		"client_id":     []string{store.clients[0].GetClientID()},
		"client_secret": []string{store.clients[0].GetClientSecret()},
		"refresh_token": []string{token1.RefreshToken},
	})
	token2 := parseResult(response)
	if token2.AccessToken == token1.AccessToken {
		t.Errorf("Expect new access token but found %s", token1.AccessToken)
	}
	if token2.RefreshToken == token1.RefreshToken {
		t.Errorf("Expect new refresh token but found %s", token1.RefreshToken)
	}
}
