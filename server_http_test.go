package oauth2

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func Test_ServeHTTP(t *testing.T) {
	defer os.Remove(ConfigFile)

	store := createStore()
	server := DefaultServerWithTokenStore(store)

	request, _ := http.NewRequest("POST", "http://localhost:8080/token", strings.NewReader(fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s", PasswordGrant, store.clients[0].GetClientID(), store.clients[0].GetClientSecret(), "admin", "admin")))
	request.Header.Set("content-type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Errorf("Expect http status 200 but found %d", response.Code)
	}

	request, _ = http.NewRequest("GET", "http://localhost:8080/resources/README", nil)
	response = httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Errorf("Expect http status 200 but found %d", response.Code)
	}
}
