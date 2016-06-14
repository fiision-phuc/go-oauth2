package oauth2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/phuc0302/go-oauth2/test"
	"github.com/phuc0302/go-oauth2/utils"
)

func Test_BindForm(t *testing.T) {
	defer teardown()
	setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var form struct {
			UserID    string `userID`
			ProfileID int64  `profileID`
		}

		context := objectFactory.CreateRequestContext(r, w)
		context.BindForm(&form)

		if form.UserID != "1" {
			t.Errorf(test.ExpectedStringButFoundString, "1", form.UserID)
		}
		if form.ProfileID != 2 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 2, form.ProfileID)
		}
	}))
	defer ts.Close()
	http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader("userID=1&profileID=2"))
}

func Test_BindJSON(t *testing.T) {
	defer teardown()
	setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)

		status := utils.Status{}
		context.BindJSON(&status)

		if status.Code != 200 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 200, status.Code)
		}
		if status.Description != http.StatusText(200) {
			t.Errorf(test.ExpectedStringButFoundString, http.StatusText(200), status.Description)
		}
	}))
	defer ts.Close()
	b, _ := json.Marshal(utils.Status200())

	request, _ := http.NewRequest("POST", ts.URL, bytes.NewBuffer(b))
	request.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	client.Do(request)
}

func Test_OutputError(t *testing.T) {
	defer teardown()
	setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)

		context.OutputHeader("test-header", "test-header-value")
		context.OutputError(utils.Status400())
	}))
	defer ts.Close()

	// [Test 1] Output header
	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", nil)
	if response.Header.Get("test-header") != "test-header-value" {
		t.Errorf(test.ExpectedStringButFoundString, "test-header-value", response.Header.Get("test-header"))
	}

	// [Test 2] Output error
	bytes, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode != 400 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 400, response.StatusCode)
	}
	if string(bytes) != "{\"status\":400,\"error\":\"Bad Request\",\"error_description\":\"Bad Request\"}" {
		t.Errorf(test.ExpectedStringButFoundString, "{\"status\":400,\"error\":\"Bad Request\",\"error_description\":\"Bad Request\"}", string(bytes))
	}
}

func Test_OutputRedirect(t *testing.T) {
	defer teardown()
	setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		context.OutputRedirect(utils.Status301(), "https://www.google.com")
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", nil)
	if response.StatusCode != 301 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 301, response.StatusCode)
	}
	if response.Header.Get("Location") != "https://www.google.com" {
		t.Errorf(test.ExpectedStringButFoundString, "https://www.google.com", response.Header.Get("Location"))
	}
}

func Test_OutputText(t *testing.T) {
	defer teardown()
	setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		context.OutputText(utils.Status200(), "Sample test!")
	}))
	defer ts.Close()

	response, _ := http.Get(ts.URL)
	bytes, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 200, response.StatusCode)
	}
	if string(bytes) != "Sample test!" {
		t.Errorf(test.ExpectedStringButFoundString, "Sample test!", string(bytes))
	}
}
