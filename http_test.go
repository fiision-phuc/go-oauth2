package oauth2

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/phuc0302/go-oauth2/test"
	"github.com/phuc0302/go-oauth2/utils"
)

func Test_ServeHTTP(t *testing.T) {
	defer os.Remove(debug)
	defer os.RemoveAll("resources")

	server := DefaultServer(true)
	utils.CreateDir("resources", (os.ModeDir | os.ModePerm))

	// Generate resources file
	input, _ := os.Open("LICENSE")
	output, _ := os.Create("resources/LICENSE")
	defer input.Close()
	defer output.Close()
	io.Copy(output, input)

	// Update allow methods
	Cfg.AllowMethods = []string{GET, POST, PATCH, DELETE}
	methodsValidation = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(Cfg.AllowMethods, "|")))

	// [Test 1] Invalid resource request
	request, _ := http.NewRequest("GET", "http://localhost:8080/resources/README", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != 404 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 404, response.Code)
	}

	// [Test 2] Valid resource request
	request, _ = http.NewRequest("GET", "http://localhost:8080/resources/LICENSE", nil)
	response = httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 200, response.Code)
	}

	// [Test 3] Invalid HTTP method
	request, _ = http.NewRequest("LINK", "http://localhost:8080/token", nil)
	response = httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != 405 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 405, response.Code)
	}

	// [Test 4] Invalid url
	request, _ = http.NewRequest("POST", "http://localhost:8080", strings.NewReader(""))
	request.Header.Set("content-type", "application/x-www-form-urlencoded")
	response = httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != 503 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 503, response.Code)
	}

	// [Test 5] Valid url
	server.Get("/sample", func(c *Request) {
		data := map[string]string{"apple": "apple"}
		c.OutputJSON(utils.Status200(), data)
	})

	request, _ = http.NewRequest("GET", "http://localhost:8080/sample", nil)
	response = httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 200, response.Code)
	} else {
		bytes, _ := ioutil.ReadAll(response.Body)

		if string(bytes) != "{\"apple\":\"apple\"}" {
			t.Errorf(test.ExpectedStringButFoundString, "{\"apple\":\"apple\"}", string(bytes))
		}
	}
}

func Test_serveRequestWithOAuth2Disable(t *testing.T) {
	//	defer os.Remove(ConfigFile)
	//	server := DefaultServer()

	//		server.Get("/sample", func(c *Request) {
	//			data := map[string]string{"apple": "apple"}
	//			c.OutputJSON(utils.Status200(), data)
	//		})

	//	// Send invalid url request
	//	request, _ := http.NewRequest("GET", "http://localhost:8080/data", nil)
	//	response := httptest.NewRecorder()
	//	server.ServeHTTP(response, request)
	//	if response.Code != 503 {
	//		t.Errorf("Expected http status 503 but found %d", response.Code)
	//	}

	//	// Send valid url request
	//	request, _ = http.NewRequest("GET", "http://localhost:8080/sample", nil)
	//	response = httptest.NewRecorder()
	//	server.ServeHTTP(response, request)
	//	if response.Code != 200 {
	//		t.Errorf("Expected http status 200 but found %d", response.Code)
	//	}

	//	data, _ := ioutil.ReadAll(response.Body)
	//	if string(data) != "{\"apple\":\"apple\"}" {
	//		t.Errorf("Expected \"%s\" but found \"%s\"", "{\"apple\":\"apple\"}", string(data))
	//	}
}

func Test_serveRequestWithOAuth2Enable(t *testing.T) {
	//	defer os.Remove(ConfigFile)
	//	s := DefaultServerWithTokenStore(createStore())

	//	s.Get("/user", func(c *Request) {
	//		data := map[string]string{"user": "r_user"}
	//		c.OutputJSON(utils.Status200(), data)
	//	})
	//	s.Get("/admin", func(c *Request) {
	//		data := map[string]string{"user": "r_admin"}
	//		c.OutputJSON(utils.Status200(), data)
	//	})
	//	s.Get("/manager", func(c *Request) {
	//		data := map[string]string{"user": "r_manager"}
	//		c.OutputJSON(utils.Status200(), data)
	//	})

	//	s.AddRoles("/user", "r_user")
	//	s.AddRoles("/admin", "r_admin")
	//	s.AddRoles("/manager", "r_manager")

	//	// Get token
	//	request, _ := http.NewRequest("POST", "http://localhost:8080/token", strings.NewReader(fmt.Sprintf(
	//		"grant_type=%s&client_id=%s&client_secret=%s&username=%s&password=%s",
	//		PasswordGrant,
	//		clientID.Hex(),
	//		clientSecret.Hex(),
	//		username,
	//		password,
	//	)))
	//	request.Header.Set("content-type", "application/x-www-form-urlencoded")
	//	response := httptest.NewRecorder()
	//	s.ServeHTTP(response, request)

	//	if response.Code != 200 {
	//		t.Errorf("Expected http status 200 but found %d", response.Code)
	//	}

	//	token := TokenResponse{}
	//	json.Unmarshal(response.Body.Bytes(), &token)

	//	// Test unauthorized access
	//	request, _ = http.NewRequest("GET", "http://localhost:8080/user", nil)
	//	response = httptest.NewRecorder()
	//	s.ServeHTTP(response, request)
	//	if response.Code != 401 {
	//		t.Errorf("Expected http status 401 but found %d", response.Code)
	//	}

	//	// Test authorized access
	//	request.Header.Set("authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	//	response = httptest.NewRecorder()
	//	s.ServeHTTP(response, request)
	//	if response.Code != 200 {
	//		t.Errorf("Expected http status 200 but found %d", response.Code)
	//	}
	//	if string(response.Body.Bytes()) != "{\"user\":\"r_user\"}" {
	//		t.Errorf("Expected \"%s\" but found \"%s\"", "{\"user\":\"r_user\"}", string(response.Body.Bytes()))
	//	}

	//	// Text invalid role
	//	request, _ = http.NewRequest("GET", "http://localhost:8080/manager", nil)
	//	response = httptest.NewRecorder()
	//	s.ServeHTTP(response, request)
	//	if response.Code != 401 {
	//		t.Errorf("Expected http status 401 but found %d", response.Code)
	//	}
}
