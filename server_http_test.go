package oauth2

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/phuc0302/go-oauth2/utils"
)

func Test_ServeHTTP(t *testing.T) {
	defer os.Remove(ConfigFile)

	store := createStore()
	server := DefaultServerWithTokenStore(store)

	request, _ := http.NewRequest("POST", "http://localhost:8080/token", strings.NewReader(""))
	request.Header.Set("content-type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != 400 {
		t.Errorf("Expected http status 400 but found %d", response.Code)
	}

	request, _ = http.NewRequest("GET", "http://localhost:8080/resources/README", nil)
	response = httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Errorf("Expected http status 200 but found %d", response.Code)
	}
}

func Test_serveRequestWithOAuth2Disable(t *testing.T) {
	defer os.Remove(ConfigFile)
	server := DefaultServer()

	server.Get("/sample", func(c *RequestContext) {
		data := map[string]string{"apple": "apple"}
		c.OutputJSON(utils.Status200(), data)
	})

	// Send invalid url request
	request, _ := http.NewRequest("GET", "http://localhost:8080/data", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != 503 {
		t.Errorf("Expected http status 503 but found %d", response.Code)
	}

	// Send valid url request
	request, _ = http.NewRequest("GET", "http://localhost:8080/sample", nil)
	response = httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Errorf("Expected http status 200 but found %d", response.Code)
	}

	data, _ := ioutil.ReadAll(response.Body)
	if string(data) != "{\"apple\":\"apple\"}" {
		t.Errorf("Expected \"%s\" but found \"%s\"", "{\"apple\":\"apple\"}", string(data))
	}
}

func Test_serveRequestWithOAuth2Enable(t *testing.T) {
	defer os.Remove(ConfigFile)
	//	server := DefaultServerWithTokenStore(createStore())
}
