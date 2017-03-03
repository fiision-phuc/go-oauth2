package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/phuc0302/go-oauth2/test"
	"github.com/phuc0302/go-oauth2/util"
)

func Test_CreateRequestContext_GetRequest(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := CreateContext(w, r)

		if context.Path != httprouter.CleanPath(r.URL.Path) {
			t.Errorf(test.ExpectedStringButFoundString, httprouter.CleanPath(r.URL.Path), context.Path)
		}
		if context.Header == nil {
			t.Error(test.ExpectedNotNil)
		} else {
			if len(context.Header) != 2 {
				t.Errorf(test.ExpectedNumberButFoundNumber, 2, len(context.Header))
			} else {
				if context.Header["user-agent"] != "go-http-client/1.1" {
					t.Errorf(test.ExpectedStringButFoundString, "go-http-client/1.1", context.Header["user-agent"])
				}
				if context.Header["accept-encoding"] != "gzip" {
					t.Errorf(test.ExpectedStringButFoundString, "gzip", context.Header["accept-encoding"])
				}
			}
		}
		if context.PathParams != nil {
			t.Error(test.ExpectedNil)
		}
		if context.QueryParams != nil {
			t.Error(test.ExpectedNil)
		}
	}))
	defer ts.Close()
	http.Get(ts.URL)
}

func Test_CreateRequestContext_GetRequestWithQueryParams(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := CreateContext(w, r)

		if context.QueryParams == nil {
			t.Error(test.ExpectedNotNil)
		} else {
			if len(context.QueryParams) != 2 {
				t.Errorf(test.ExpectedNumberButFoundNumber, 2, len(context.QueryParams))
			} else {
				if context.QueryParams["userID"] != "1" {
					t.Errorf(test.ExpectedStringButFoundString, "1", context.QueryParams["userID"])
				}
				if context.QueryParams["profileID"] != "2" {
					t.Errorf(test.ExpectedStringButFoundString, "2", context.QueryParams["profileID"])
				}
			}
		}
	}))
	defer ts.Close()
	http.Get(fmt.Sprintf("%s?userID=1&profileID=2", ts.URL))
}

func Test_CreateRequestContext_PostFormRequest(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := CreateContext(w, r)

		if context.Header["content-type"] != "application/x-www-form-urlencoded" {
			t.Errorf(test.ExpectedStringButFoundString, "application/x-www-form-urlencoded", context.Header["content-type"])
		}
		if context.QueryParams != nil {
			t.Error(test.ExpectedNil)
		}
	}))
	defer ts.Close()
	http.Post(ts.URL, strings.ToUpper("application/x-www-form-urlencoded"), nil)
}

func Test_CreateRequestContext_PostFormRequestWithData(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := CreateContext(w, r)

		if context.Header["content-type"] != "application/x-www-form-urlencoded" {
			t.Errorf(test.ExpectedStringButFoundString, "application/x-www-form-urlencoded", context.Header["content-type"])
		}
		if context.QueryParams == nil {
			t.Error(test.ExpectedNotNil)
		} else {
			if context.QueryParams["userID"] != "1" {
				t.Errorf(test.ExpectedStringButFoundString, "1", context.QueryParams["userID"])
			}
			if context.QueryParams["profileID"] != "2" {
				t.Errorf(test.ExpectedStringButFoundString, "2", context.QueryParams["profileID"])
			}
		}

	}))
	defer ts.Close()
	http.Post(ts.URL, strings.ToUpper("application/x-www-form-urlencoded"), strings.NewReader("userID=1&profileID=2"))
}

func Test_CreateRequestContext_PostMultipartRequest(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := CreateContext(w, r)

		if context.Header["content-type"] != "multipart/form-data; boundary=gc0p4jq0m2yt08ju534c0p" {
			t.Errorf(test.ExpectedStringButFoundString, "multipart/form-data; boundary=gc0p4jq0m2yt08ju534c0p", context.Header["content-type"])
		}
		if context.QueryParams != nil {
			t.Error(test.ExpectedNil)
		}
	}))
	defer ts.Close()

	request, _ := http.NewRequest("POST", ts.URL, nil)
	request.Header.Set("content-type", "multipart/form-data; boundary=gc0p4Jq0M2Yt08jU534c0p")

	http.DefaultClient.Do(request)
}

func Test_CreateRequestContext_PostMultipartRequestWithData(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := CreateContext(w, r)

		if context.Header["content-type"] != "multipart/form-data; boundary=gc0p4jq0m2yt08ju534c0p" {
			t.Errorf(test.ExpectedStringButFoundString, "multipart/form-data; boundary=gc0p4jq0m2yt08ju534c0p", context.Header["content-type"])
		}
		if context.QueryParams == nil {
			t.Error(test.ExpectedNotNil)
		} else {
			if context.QueryParams["userID"] != "1" {
				t.Errorf(test.ExpectedStringButFoundString, "1", context.QueryParams["userID"])
			}
			if context.QueryParams["profileID"] != "2" {
				t.Errorf(test.ExpectedStringButFoundString, "2", context.QueryParams["profileID"])
			}
		}
	}))
	defer ts.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.SetBoundary("gc0p4Jq0M2Yt08jU534c0p")

	p := map[string]string{
		"userID":    "1",
		"profileID": "2",
	}
	for key, val := range p {
		_ = writer.WriteField(key, val)
	}
	writer.Close()

	request, _ := http.NewRequest("POST", ts.URL, body)
	request.Header.Set("content-type", "multipart/form-data; boundary=gc0p4Jq0M2Yt08jU534c0p")

	http.DefaultClient.Do(request)
}

func Test_BindForm(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var form struct {
			UserID    string `field:"userID"`
			ProfileID int64  `field:"profileID"`
		}

		context := CreateContext(w, r)
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
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := CreateContext(w, r)

		status := new(util.Status)
		context.BindJSON(status)

		if status.Code != 200 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 200, status.Code)
		}
		if status.Description != http.StatusText(200) {
			t.Errorf(test.ExpectedStringButFoundString, http.StatusText(200), status.Description)
		}
	}))
	defer ts.Close()
	b, _ := json.Marshal(util.Status200())

	request, _ := http.NewRequest("POST", ts.URL, bytes.NewBuffer(b))
	request.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	client.Do(request)
}

func Test_OutputHeader(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := CreateContext(w, r)

		context.OutputHeader("test-header", "test-header-value")
		context.OutputError(util.Status200())
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", nil)
	if response.Header.Get("test-header") != "test-header-value" {
		t.Errorf(test.ExpectedStringButFoundString, "test-header-value", response.Header.Get("test-header"))
	}
}

func Test_OutputError(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := CreateContext(w, r)
		context.OutputError(util.Status400())
	}))
	defer ts.Close()

	response, _ := http.Get(ts.URL)
	bytes, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 400 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 400, response.StatusCode)
	} else {
		if string(bytes) != "{\"status\":400,\"error\":\"Bad Request\",\"error_description\":\"Bad Request\"}" {
			t.Errorf(test.ExpectedStringButFoundString, "{\"status\":400,\"error\":\"Bad Request\",\"error_description\":\"Bad Request\"}", string(bytes))
		}
	}
}

func Test_OutputRedirect(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := CreateContext(w, r)
		context.OutputRedirect(util.Status301(), "https://www.google.com")
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
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := CreateContext(w, r)
		context.OutputText(util.Status200(), "Sample test!")
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

func BenchmarkURLEncodeForm(b *testing.B) {
	request, _ := http.NewRequest("POST", "http://localhost:8080/", strings.NewReader("userID=1&profileID=2"))
	request.Header["content-type"] = []string{"application/x-www-form-urlencoded"}

	for n := 0; n < b.N; n++ {
		CreateContext(nil, request)
	}
}
