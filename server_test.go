package oauth2

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/phuc0302/go-oauth2/utils"

	"gopkg.in/mgo.v2/bson"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// In Memory Store
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

func (s *InMemoryStore) FindClientWithID(clientID string) AuthClient {
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

////////////////////////////////////////////////////////////////////////////////////////////////////
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

////////////////////////////////////////////////////////////////////////////////////////////////////

func Test_DefaultServer(t *testing.T) {
	defer os.Remove(ConfigFile)
	server := DefaultServer()

	if server.Config == nil {
		t.Errorf("Expected config file should be loaded at creation time but found nil.")
	}

	if server.routes != nil {
		t.Error("Expected routes nil but found not nil.")
	}

	if server.groups != nil {
		t.Error("Expected groups will be nil but found not nil.")
	}

	if server.logger == nil {
		t.Error("Expected logger not nil but found nil.")
	}
}

func Test_DefaultServerWithTokenStore(t *testing.T) {
	defer os.Remove(ConfigFile)
	server := DefaultServerWithTokenStore(createStore())

	if server.Config == nil {
		t.Errorf("Expected config file should be loaded at creation time but found nil.")
	}

	if server.routes == nil {
		t.Error("Expected routes not nil but found nil.")
	}

	if server.groups != nil {
		t.Error("Expected groups will be nil but found not nil.")
	}

	if server.logger == nil {
		t.Error("Expected logger not nil but found nil.")
	}
}

//func Test_Run(t *testing.T) {
//	defer os.Remove(ConfigFile)

//	config := LoadConfigs()
//	config.AllowMethods = []string{GET}

//	configJSON, _ := json.MarshalIndent(config, "", "  ")
//	file, _ := os.Create(ConfigFile)
//	file.Write(configJSON)
//	file.Close()

//	config = LoadConfigs()

//	server := DefaultServerWithTokenStore(createStore())
//	go server.Run()

//	response, _ := http.PostForm("http://localhost:8080/token", url.Values{
//		"grant_type":    []string{PasswordGrant},
//		"client_id":     []string{bson.NewObjectId().Hex()},
//		"client_secret": []string{bson.NewObjectId().Hex()},
//		"username":      []string{"admin"},
//		"password":      []string{"admin"},
//	})
//	status := parseError(response)

//	if status.Code != 405 {
//		t.Errorf("Expect http status 405 but found %d", status.Code)
//	}
//}
