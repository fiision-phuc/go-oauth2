package oauth2

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/phuc0302/go-oauth2/test"
	"github.com/phuc0302/go-oauth2/util"
)

func Test_ServeHTTP_InvalidResource(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup server & test server
	server := DefaultServer(true)
	server.Get("/sample", func(c *RequestContext, s *OAuthContext) {
		c.OutputJSON(util.Status200(), map[string]string{"apple": "apple"})
	})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	}))
	defer ts.Close()

	response, _ := http.Get(fmt.Sprintf("%s/%s", ts.URL, "resources/README"))
	if response.StatusCode != 404 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 404, response.StatusCode)
	}
}

func Test_ServeHTTP_ValidResource(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup server & test server
	server := DefaultServer(true)
	server.Get("/sample", func(c *RequestContext, s *OAuthContext) {
		c.OutputJSON(util.Status200(), map[string]string{"apple": "apple"})
	})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	}))
	defer ts.Close()

	response, _ := http.Get(fmt.Sprintf("%s/%s", ts.URL, "resources/LICENSE"))
	if response.StatusCode != 200 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 200, response.StatusCode)
	}
}

func Test_ServeHTTP_InvalidHTTPMethod(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup server & test server
	server := DefaultServer(true)
	server.Get("/sample", func(c *RequestContext, s *OAuthContext) {
		c.OutputJSON(util.Status200(), map[string]string{"apple": "apple"})
	})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	}))
	defer ts.Close()

	// Update allow methods
	Cfg.AllowMethods = []string{Get, Post, Patch, Delete}
	methodsValidation = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(Cfg.AllowMethods, "|")))

	request, _ := http.NewRequest("LINK", fmt.Sprintf("%s/%s", ts.URL, "token"), nil)
	response, _ := http.DefaultClient.Do(request)
	if response.StatusCode != 405 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 405, response.StatusCode)
	}
}

func Test_ServeHTTP_InvalidURL(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup server & test server
	server := DefaultServer(true)
	server.Get("/sample", func(c *RequestContext, s *OAuthContext) {
		c.OutputJSON(util.Status200(), map[string]string{"apple": "apple"})
	})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	}))
	defer ts.Close()

	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/%s", ts.URL, "sample"), strings.NewReader(""))
	request.Header.Set("content-type", "application/x-www-form-urlencoded")

	response, _ := http.DefaultClient.Do(request)
	if response.StatusCode != 503 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 503, response.StatusCode)
	}
}

func Test_ServeHTTP_ValidURL(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup server & test server
	server := DefaultServer(true)
	server.Get("/sample", func(c *RequestContext, s *OAuthContext) {
		c.OutputJSON(util.Status200(), map[string]string{"apple": "apple"})
	})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	}))
	defer ts.Close()

	response, _ := http.Get(fmt.Sprintf("%s/%s", ts.URL, "sample"))
	if response.StatusCode != 200 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 200, response.StatusCode)
	} else {
		bytes, _ := ioutil.ReadAll(response.Body)

		if string(bytes) != "{\"apple\":\"apple\"}" {
			t.Errorf(test.ExpectedStringButFoundString, "{\"apple\":\"apple\"}", string(bytes))
		}
	}
}
