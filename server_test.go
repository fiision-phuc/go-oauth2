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

func createStore() *InMemoryTokenStore {
	return &InMemoryTokenStore{
		clients: []AuthClient{
			&AuthClientDefault{
				ClientID:     bson.NewObjectId(),
				ClientSecret: bson.NewObjectId(),
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
				ClientID:    clientID,
				Token:       utils.GenerateToken(),
				CreatedTime: createdTime,
				ExpiredTime: createdTime.Add(3600 * time.Second),
			},
		},
		refreshTokens: []Token{
			&TokenDefault{
				TokenID:     bson.NewObjectId(),
				UserID:      userID,
				ClientID:    clientID,
				Token:       utils.GenerateToken(),
				CreatedTime: createdTime,
				ExpiredTime: createdTime.Add(1209600 * time.Second),
			},
		},
	}
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

func Test_AddRolesWillBeDisabled(t *testing.T) {
	defer os.Remove(ConfigFile)
	server := DefaultServer()

	server.AddRoles("//..//user/**", "r_user")
	if len(server.userRoles) != 0 {
		t.Errorf("Expect user's role validation must be ignored but found %d", len(server.userRoles))
	}
}

func Test_AddRolesWillBeEnabled(t *testing.T) {
	defer os.Remove(ConfigFile)
	server := DefaultServerWithTokenStore(createStore())

	server.AddRoles("//..//user/**", "r_user, r_admin")
	if len(server.userRoles) != 1 {
		t.Errorf("Expect user's role validation must be 1 but found %d", len(server.userRoles))
	}
	for rule, v := range server.userRoles {
		if len(v) != 2 {
			t.Errorf("Expect roles validation must be 2 but found %d", len(v))
		}
		if !rule.MatchString("/user/username/password") {
			t.Error("Expect \"/user/username/password\" require r_user but fount not.")
		}

		if rule.MatchString("/username/password") {
			t.Error("Expect \"/username/password\" require none but fount not.")
		}
	}

	server = DefaultServerWithTokenStore(createStore())
	server.AddRoles("//..//user/:userId/**", "r_user, r_admin")
	for rule, _ := range server.userRoles {
		if !rule.MatchString("/user/123456/username/password") {
			t.Error("Expect \"/user/123456/username/password\" require r_user but fount not.")
		}
	}
}
