package oauth2

import (
	"os"
	"testing"
)

func Test_AddRoles(t *testing.T) {
	defer os.Remove(ConfigFile)
	server := DefaultServer()

	server.AddRoles("//..//user/**", "r_user")
	if len(server.userRoles) != 1 {
		t.Errorf("Expect user's role validation must be 1 but found %d", len(server.userRoles))
	}
	for rule, _ := range server.userRoles {
		if !rule.MatchString("/user/username/password") {
			t.Error("Expect \"/user/username/password\" require r_user but fount not.")
		}

		if rule.MatchString("/username/password") {
			t.Error("Expect \"/username/password\" require none but fount not.")
		}
	}

	server = DefaultServer()

	server.AddRoles("//..//user/:userId/**", "r_user, r_admin")
	if len(server.userRoles) != 1 {
		t.Errorf("Expect user's role validation must be 1 but found %d", len(server.userRoles))
	}
	for rule, _ := range server.userRoles {
		if !rule.MatchString("/user/123456/username/password") {
			t.Error("Expect \"/user/123456/username/password\" require r_user but fount not.")
		}
	}
}

//func Test_ServeHTTP(t *testing.T) {
//	defer os.Remove(ConfigFile)

//	store := createStore()
//	server := DefaultServerWithTokenStore(store)

//	go server.Run()

//	response, _ := http.PostForm("http://localhost:8080/token", url.Values{
//		"grant_type":    []string{PasswordGrant},
//		"client_id":     []string{store.clients[0].GetClientID()},
//		"client_secret": []string{store.clients[0].GetClientSecret()},
//		"username":      []string{"admin"},
//		"password":      []string{"admin"},
//	})

//	if response.StatusCode != 200 {
//		t.Errorf("Expect http status 200 but found %d", response.StatusCode)
//	}

//	response, _ = http.Get("http://localhost:8080/resources/README")
//	if response.StatusCode != 200 {
//		t.Errorf("Expect http status 200 but found %d", response.StatusCode)
//	}
//}
