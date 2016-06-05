package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/phuc0302/go-oauth2/test"
	"github.com/phuc0302/go-oauth2/utils"

	"gopkg.in/mgo.v2/bson"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

var userID = bson.NewObjectId()
var username = "admin"
var password = "admin"

var clientID = bson.NewObjectId()
var clientSecret = bson.NewObjectId()

var createdTime, _ = time.Parse(time.RFC822, "02 Jan 06 15:04 MST")

////////////////////////////////////////////////////////////////////////////////////////////////////

func createStore() *DefaultInMemoryStore {
	return &DefaultInMemoryStore{
		clients: []IClient{
			&DefaultClient{
				ID:        clientID,
				Secret:    clientSecret,
				Grants:    []string{PasswordGrant, RefreshTokenGrant},
				Redirects: []string{"http://sample01.com", "http://sample02.com"},
			},
		},
		users: []IUser{
			&DefaultUser{
				ID:    userID,
				User:  "admin",
				Pass:  "admin",
				Roles: []string{"r_user", "r_admin"},
			},
			&DefaultUser{
				ID:   bson.NewObjectId(),
				User: "admin2",
				Pass: "admin2",
			},
		},
		accessTokens: []IToken{
			&DefaultToken{
				ID:      bson.NewObjectId(),
				User:    userID,
				Client:  clientID,
				Created: createdTime,
				Expired: createdTime.Add(3600 * time.Second),
			},
		},
		refreshTokens: []IToken{
			&DefaultToken{
				ID:      bson.NewObjectId(),
				User:    userID,
				Client:  clientID,
				Created: createdTime,
				Expired: createdTime.Add(1209600 * time.Second),
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
	defer os.Remove(debug)
	server := DefaultServer(true)

	if cfg == nil {
		t.Error(test.ExpectedNotNil)
	}
	if objectFactory == nil {
		t.Error(test.ExpectedNotNil)
	}
	if server.router == nil {
		t.Error(test.ExpectedNotNil)
	}
}

func Test_Run(t *testing.T) {
	defer os.Remove(debug)

	server := DefaultServer(true)
	go server.Run()

	response, err := http.PostForm("http://localhost:8080/token", url.Values{
		"grant_type":    []string{PasswordGrant},
		"client_id":     []string{bson.NewObjectId().Hex()},
		"client_secret": []string{bson.NewObjectId().Hex()},
		"username":      []string{"admin"},
		"password":      []string{"admin"},
	})
	//	status := parseError(response)

	if response == nil {
		fmt.Println(err)
		t.Error(test.ExpectedNotNil)
	} else {
		if response.StatusCode != 405 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 405, response.StatusCode)
		}
	}

}

func Test_AddRolesWillBeDisabled(t *testing.T) {
	//	defer os.Remove(ConfigFile)
	//	server := DefaultServer()

	//	server.AddRoles("//..//user/**", "r_user")
	//	if len(server.userRoles) != 0 {
	//		t.Errorf("Expect user's role validation must be ignored but found %d", len(server.userRoles))
	//	}
}

func Test_AddRolesWillBeEnabled(t *testing.T) {
	//	defer os.Remove(ConfigFile)
	//	server := DefaultServerWithTokenStore(createStore())

	//	server.AddRoles("//..//user/**", "r_user, r_admin")
	//	if len(server.userRoles) != 1 {
	//		t.Errorf("Expect user's role validation must be 1 but found %d", len(server.userRoles))
	//	}
	//	for rule, v := range server.userRoles {
	//		if len(v) != 2 {
	//			t.Errorf("Expect roles validation must be 2 but found %d", len(v))
	//		}
	//		if !rule.MatchString("/user/username/password") {
	//			t.Error("Expect \"/user/username/password\" require r_user but fount not.")
	//		}

	//		if rule.MatchString("/username/password") {
	//			t.Error("Expect \"/username/password\" require none but fount not.")
	//		}
	//	}

	//	server = DefaultServerWithTokenStore(createStore())
	//	server.AddRoles("//..//user/:userId/**", "r_user, r_admin")
	//	for rule, _ := range server.userRoles {
	//		if !rule.MatchString("/user/123456/username/password") {
	//			t.Error("Expect \"/user/123456/username/password\" require r_user but fount not.")
	//		}
	//	}
}
