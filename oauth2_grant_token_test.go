package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/phuc0302/go-oauth2/utils"

	"gopkg.in/mgo.v2/bson"
)

////////////////////////////////////////////////////////////////////////////////
// In Memory Store															    //
////////////////////////////////////////////////////////////////////////////////
var userID = bson.NewObjectId()
var clientID = bson.NewObjectId()
var createdTime, _ = time.Parse(time.RFC822, "02 Jan 06 15:04 MST")

type InMemoryStore struct {
	clients []AuthClient
	users   []AuthUser

	accessTokens  []Token
	refreshTokens []Token
}

func createStore() *InMemoryStore {
	return &InMemoryStore{
		clients: []AuthClient{
			&AuthClientDefault{
				ClientID:     bson.NewObjectId().Hex(),
				ClientSecret: bson.NewObjectId().Hex(),
				GrantTypes:   []string{PasswordGrant, RefreshTokenGrant},
				RedirectURIs: []string{"http://sample01.com", "http://sample02.com"},
			},
		},
		users: []AuthUser{
			&AuthUserDefault{
				UserID:   bson.NewObjectId(),
				Username: "admin",
				Password: "admin",
			},
			&AuthUserDefault{
				UserID:   userID,
				Username: "admin2",
				Password: "admin2",
			},
		},
		accessTokens: []Token{
			&TokenDefault{
				TokenID:     bson.NewObjectId(),
				UserID:      userID,
				ClientID:    clientID.Hex(),
				Token:       utils.GenerateToken(),
				CreatedTime: createdTime,
				ExpiredTime: createdTime.Add(3600 * time.Second),
			},
		},
		refreshTokens: []Token{
			&TokenDefault{
				TokenID:     bson.NewObjectId(),
				UserID:      userID,
				ClientID:    clientID.Hex(),
				Token:       utils.GenerateToken(),
				CreatedTime: createdTime,
				ExpiredTime: createdTime.Add(1209600 * time.Second),
			},
		},
	}
}

func (s *InMemoryStore) FindUserWithID(userID string) AuthUser {
	for _, user := range s.users {
		if user.GetUserID() == userID {
			return user
		}
	}
	return nil
}
func (s *InMemoryStore) FindUserWithClient(clientID string, clientSecret string) AuthUser {
	return nil
}
func (s *InMemoryStore) FindUserWithCredential(username string, password string) AuthUser {
	for _, user := range s.users {
		if user.GetUsername() == username && user.GetPassword() == password {
			return user
		}
	}
	return nil
}

func (s *InMemoryStore) FindClientWithCredential(clientID string, clientSecret string) AuthClient {
	for _, client := range s.clients {
		if client.GetClientID() == clientID && client.GetClientSecret() == clientSecret {
			return client
		}
	}
	return nil
}

func (s *InMemoryStore) FindAccessToken(token string) Token {
	for _, recordToken := range s.accessTokens {
		if recordToken.GetToken() == token {
			return recordToken
		}
	}
	return nil
}
func (s *InMemoryStore) FindAccessTokenWithCredential(clientID string, userID string) Token {
	for _, recordToken := range s.accessTokens {
		if recordToken.GetUserID() == userID && recordToken.GetClientID() == clientID {
			return recordToken
		}
	}
	return nil
}
func (s *InMemoryStore) CreateAccessToken(clientID string, userID string, token string, createdTime time.Time, expiredTime time.Time) Token {
	newToken := &TokenDefault{
		TokenID:     bson.NewObjectId(),
		UserID:      bson.ObjectIdHex(userID),
		ClientID:    clientID,
		Token:       utils.GenerateToken(),
		CreatedTime: createdTime,
		ExpiredTime: expiredTime,
	}

	s.accessTokens = append(s.accessTokens, newToken)
	return newToken
}
func (s *InMemoryStore) DeleteAccessToken(token Token) {
	for idx, recordToken := range s.accessTokens {
		if recordToken == token {
			s.accessTokens = append(s.accessTokens[:idx], s.accessTokens[idx+1:]...)
			break
		}
	}
}
func (s *InMemoryStore) SaveAccessToken(token Token) {
	isUpdated := false
	for _, recordToken := range s.accessTokens {
		if recordToken == token {
			token.SetToken(token.GetToken())
			token.SetCreatedTime(token.GetCreatedTime())
			token.SetExpiredTime(token.GetExpiredTime())
			isUpdated = true
			break
		}
	}

	if !isUpdated {
		s.accessTokens = append(s.accessTokens, token)
	}
}

func (s *InMemoryStore) FindRefreshToken(token string) Token {
	for _, recordToken := range s.refreshTokens {
		if recordToken.GetToken() == token {
			return recordToken
		}
	}
	return nil
}
func (s *InMemoryStore) FindRefreshTokenWithCredential(clientID string, userID string) Token {
	for _, recordToken := range s.refreshTokens {
		if recordToken.GetUserID() == userID && recordToken.GetClientID() == clientID {
			return recordToken
		}
	}
	return nil
}
func (s *InMemoryStore) CreateRefreshToken(clientID string, userID string, token string, createdTime time.Time, expiredTime time.Time) Token {
	newToken := &TokenDefault{
		TokenID:     bson.NewObjectId(),
		UserID:      bson.ObjectIdHex(userID),
		ClientID:    clientID,
		Token:       utils.GenerateToken(),
		CreatedTime: createdTime,
		ExpiredTime: expiredTime,
	}

	s.refreshTokens = append(s.refreshTokens, newToken)
	return newToken
}
func (s *InMemoryStore) DeleteRefreshToken(token Token) {
	for idx, recordToken := range s.refreshTokens {
		if recordToken == token {
			s.refreshTokens = append(s.refreshTokens[:idx], s.refreshTokens[idx+1:]...)
			break
		}
	}
}
func (s *InMemoryStore) SaveRefreshToken(token Token) {
	isUpdated := false
	for _, recordToken := range s.refreshTokens {
		if recordToken == token {
			recordToken.SetToken(token.GetToken())
			recordToken.SetCreatedTime(token.GetCreatedTime())
			recordToken.SetExpiredTime(token.GetExpiredTime())
			isUpdated = true
			break
		}
	}

	if !isUpdated {
		s.refreshTokens = append(s.refreshTokens, token)
	}
}

func (s *InMemoryStore) FindAuthorizationCode(authorizationCode string) {

}
func (s *InMemoryStore) SaveAuthorizationCode(authorizationCode string, clientID string, expires time.Time) {

}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Helper
func parseError(response *http.Response) *utils.Status {
	data, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	status := utils.Status{}
	json.Unmarshal(data, &status)

	return &status
}
func parseResult(response *http.Response) *TokenResponse {
	data, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	token := TokenResponse{}
	json.Unmarshal(data, &token)

	return &token
}

/////////////////////////////////////////////////////////////////////////////////////////////////

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
