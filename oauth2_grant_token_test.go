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
	users         []AuthUserDefault
	clients       []AuthClientDefault
	accessTokens  []Token
	refreshTokens []TokenDefault
}

func createStore() *InMemoryStore {
	return &InMemoryStore{
		users: []AuthUserDefault{
			AuthUserDefault{
				UserID:   bson.NewObjectId(),
				Username: "admin",
				Password: "admin",
			},
			AuthUserDefault{
				UserID:   userID,
				Username: "admin2",
				Password: "admin2",
			},
		},
		clients: []AuthClientDefault{
			AuthClientDefault{
				ClientID:     bson.NewObjectId().Hex(),
				ClientSecret: bson.NewObjectId().Hex(),
				GrantTypes:   []string{PasswordGrant, RefreshTokenGrant},
				RedirectURIs: []string{"http://sample01.com", "http://sample02.com"},
			},
		},
		accessTokens: []TokenDefault{
			TokenDefault{
				TokenID:     bson.NewObjectId(),
				UserID:      userID,
				ClientID:    clientID.Hex(),
				Token:       utils.GenerateToken(),
				CreatedTime: createdTime,
				ExpiredTime: createdTime.Add(3600 * time.Second),
			},
		},
		refreshTokens: []TokenDefault{
			TokenDefault{
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

func (s *InMemoryStore) FindUserWithID(userID bson.ObjectId) *AuthUserDefault {
	for _, user := range s.users {
		if user.UserID == userID {
			return &user
		}
	}
	return nil
}
func (s *InMemoryStore) FindUserWithClient(clientID string, clientSecret string) *AuthUserDefault {
	return nil
}
func (s *InMemoryStore) FindUserWithCredential(username string, password string) *AuthUserDefault {
	for _, user := range s.users {
		if user.Username == username && user.Password == password {
			return &user
		}
	}
	return nil
}

func (s *InMemoryStore) FindClientWithCredential(clientID string, clientSecret string) *AuthClientDefault {
	for _, client := range s.clients {
		if client.ClientID == clientID && client.ClientSecret == clientSecret {
			return &client
		}
	}
	return nil
}

func (s *InMemoryStore) FindToken(accessToken string) Token {
	for _, token := range s.accessTokens {
		if token.Token == accessToken {
			return &token
		}
	}
	return nil
}
func (s *InMemoryStore) FindTokenWithCredential(clientID string, userID bson.ObjectId) Token {
	for _, token := range s.accessTokens {
		if token.UserID == userID && token.ClientID == clientID {
			return &token
		}
	}
	return nil
}
func (s *InMemoryStore) CreateToken(clientID string, userID string, token string, createdTime time.Time, expiredTime time.Time) Token {
	return nil
}
func (s *InMemoryStore) DeleteToken(token Token) {
	//	for idx, token := range s.accessTokens {
	//		if token == *accessToken {
	//			s.accessTokens = append(s.accessTokens[:idx], s.accessTokens[idx+1:]...)
	//			break
	//		}
	//	}
}
func (s *InMemoryStore) SaveToken(queryToken Token) {
	isUpdated := false
	for _, token := range s.accessTokens {
		//		if token.Token == queryToken.GetToken() {
		//			token.Token = queryToken.GetToken()
		//			token.CreatedTime = queryToken.CreatedTime
		//			token.ExpiredTime = queryToken.ExpiredTime
		//			isUpdated = true
		//			break
		//		}
	}

	if !isUpdated {
		s.accessTokens = append(s.accessTokens, *queryToken)
	}
}

func (s *InMemoryStore) FindAuthorizationCode(authorizationCode string) {

}
func (s *InMemoryStore) SaveAuthorizationCode(authorizationCode string, clientID string, expires time.Time) {

}

////////////////////////////////////////////////////////////////////////////////
// Helper																	  //
////////////////////////////////////////////////////////////////////////////////
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

////////////////////////////////////////////////////////////////////////////////
// Test																		  //
////////////////////////////////////////////////////////////////////////////////
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
		"client_id":  []string{store.clients[0].ClientID},
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
		"client_id":     []string{store.clients[0].ClientID},
		"client_secret": []string{store.clients[0].ClientSecret},
	})
	status := parseError(response)
	if status.Description != fmt.Sprintf("Invalid %s parameter.", "grant_type") {
		t.Errorf("Expected \"Invalid %s parameter.\" but found \"%s\"", "grant_type", status.Description)
	}

	// Test valid request token
	response, _ = http.PostForm(ts.URL, url.Values{
		"grant_type":    []string{PasswordGrant},
		"client_id":     []string{store.clients[0].ClientID},
		"client_secret": []string{store.clients[0].ClientSecret},
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
		"client_id":     []string{store.clients[0].ClientID},
		"client_secret": []string{store.clients[0].ClientSecret},
	})
	status := parseError(response)
	if status.Description != fmt.Sprintf(templateError, "username or password") {
		t.Errorf(templateErrorMessage, "username or password", status.Description)
	}

	// Test invalid username or password
	response, _ = http.PostForm(ts.URL, url.Values{
		"grant_type":    []string{PasswordGrant},
		"client_id":     []string{store.clients[0].ClientID},
		"client_secret": []string{store.clients[0].ClientSecret},
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
		"client_id":     []string{store.clients[0].ClientID},
		"client_secret": []string{store.clients[0].ClientSecret},
		"username":      []string{"admin"},
		"password":      []string{"admin"},
	})
	token1 := parseResult(response)
	if token1 == nil {
		t.Error("Expected not nil but found nil.")
	}
	if token1.AccessToken != store.accessTokens[1].Token {
		t.Errorf("Expected %s but found %s", store.accessTokens[1].Token, token1.AccessToken)
	}
	if token1.RefreshToken != store.refreshTokens[1].Token {
		t.Errorf("Expected %s but found %s", store.refreshTokens[1].Token, token1.RefreshToken)
	}

	// Test request second token should be the same as the first one
	response, _ = http.PostForm(ts.URL, url.Values{
		"grant_type":    []string{PasswordGrant},
		"client_id":     []string{store.clients[0].ClientID},
		"client_secret": []string{store.clients[0].ClientSecret},
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
		"client_id":     []string{store.clients[0].ClientID},
		"client_secret": []string{store.clients[0].ClientSecret},
		"username":      []string{"admin2"},
		"password":      []string{"admin2"},
	})
	token3 := parseResult(response)
	if token3.AccessToken == store.accessTokens[0].Token {
		t.Errorf("Expected %s but found %s", store.accessTokens[2].Token, token3.AccessToken)
	}
	if token3.RefreshToken == store.refreshTokens[0].Token {
		t.Errorf("Expected %s but found %s", store.refreshTokens[2].Token, token1.RefreshToken)
	}
	if token3.AccessToken != store.accessTokens[2].Token {
		t.Errorf("Expected %s but found %s", store.accessTokens[2].Token, token3.AccessToken)
	}
	if token3.RefreshToken != store.refreshTokens[2].Token {
		t.Errorf("Expected %s but found %s", store.refreshTokens[2].Token, token1.RefreshToken)
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
		"client_id":     []string{store.clients[0].ClientID},
		"client_secret": []string{store.clients[0].ClientSecret},
		"username":      []string{"admin"},
		"password":      []string{"admin"},
	})
	token1 := parseResult(response)

	// Test missing refresh_token
	response, _ = http.PostForm(ts.URL, url.Values{
		"grant_type":    []string{RefreshTokenGrant},
		"client_id":     []string{store.clients[0].ClientID},
		"client_secret": []string{store.clients[0].ClientSecret},
	})
	status := parseError(response)
	if status.Description != fmt.Sprintf(templateError, "refresh_token") {
		t.Errorf(templateErrorMessage, "refresh_token", status.Description)
	}

	// Send valid request
	response, _ = http.PostForm(ts.URL, url.Values{
		"grant_type":    []string{RefreshTokenGrant},
		"client_id":     []string{store.clients[0].ClientID},
		"client_secret": []string{store.clients[0].ClientSecret},
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
