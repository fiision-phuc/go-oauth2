package oauth2

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/phuc0302/go-oauth2/test"
)

func Test_CreateRequestContext_GetRequest(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)

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
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)

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
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)

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
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)

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
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)

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
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)

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

func Test_CreateSecurityContext_NoAccessToken(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		security := objectFactory.CreateSecurityContext(context)

		if security != nil {
			t.Error(test.ExpectedNil)
		}
	}))
	defer ts.Close()
	http.Get(ts.URL)
}
func Test_CreateSecurityContext_WithBasicAuth(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		security := objectFactory.CreateSecurityContext(context)

		if security == nil {
			t.Error(test.ExpectedNotNil)
		} else {
			if security.Client == nil {
				t.Error(test.ExpectedNotNil)
			}
			if security.User == nil {
				t.Error(test.ExpectedNotNil)
			}
		}
	}))
	defer ts.Close()

	// Send token as query param
	request, _ := http.NewRequest("Get", ts.URL, nil)
	request.SetBasicAuth(u.ClientID.Hex(), u.ClientSecret.Hex())

	http.DefaultClient.Do(request)
}
func Test_CreateSecurityContext_WithGetAccessToken(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		security := objectFactory.CreateSecurityContext(context)

		if security == nil {
			t.Error(test.ExpectedNotNil)
		} else {
			if security.AccessToken == nil {
				t.Error(test.ExpectedNotNil)
			}
			if security.RefreshToken != nil {
				t.Error(test.ExpectedNil)
			}
			if security.Client == nil {
				t.Error(test.ExpectedNotNil)
			}
			if security.User == nil {
				t.Error(test.ExpectedNotNil)
			}
		}
	}))
	defer ts.Close()

	// Generate token
	now := time.Now()
	token := TokenStore.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), now, now.Add(Cfg.AccessTokenDuration))

	// Send token as query param
	http.Get(fmt.Sprintf("%s?access_token=%s", ts.URL, token.Token()))
}
func Test_CreateSecurityContext_WithPostAccessToken(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		security := objectFactory.CreateSecurityContext(context)

		if security == nil {
			t.Error(test.ExpectedNotNil)
		} else {
			if security.AccessToken == nil {
				t.Error(test.ExpectedNotNil)
			}
			if security.RefreshToken != nil {
				t.Error(test.ExpectedNil)
			}
			if security.Client == nil {
				t.Error(test.ExpectedNotNil)
			}
			if security.User == nil {
				t.Error(test.ExpectedNotNil)
			}
		}
	}))
	defer ts.Close()

	// Generate token
	now := time.Now().UTC()
	token := TokenStore.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), now, now.Add(Cfg.AccessTokenDuration))

	// Send token as authorization header
	request, _ := http.NewRequest("POST", ts.URL, nil)
	client := http.DefaultClient

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token()))
	client.Do(request)
}

func Test_CreateRoute(t *testing.T) {
	objectFactory = &DefaultFactory{}
	route := objectFactory.CreateRoute("/example/{userID}/profile/{profileID}")

	route.BindHandler(GET, func() {})
	if route == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		defaultRoute, ok := route.(*DefaultRoute)
		if ok {
			if defaultRoute.path != "/example/{userID}/profile/{profileID}" {
				t.Errorf(test.ExpectedStringButFoundString, "/example/{userID}/profile/{profileID}", defaultRoute.path)
			}
			if defaultRoute.regex == nil {
				t.Error(test.ExpectedNotNil)
			}
		} else {
			t.Errorf(test.ExpectedBoolButFoundBool, true, ok)
		}
	}
}

func Test_CreateRouter(t *testing.T) {
	objectFactory = &DefaultFactory{}
	router := objectFactory.CreateRouter()

	_, ok := router.(*DefaultRouter)
	if !ok {
		t.Errorf(test.ExpectedBoolButFoundBool, true, ok)
	}
}

func Test_CreateStore(t *testing.T) {
	objectFactory = &DefaultFactory{}
	store := objectFactory.CreateStore()

	_, ok := store.(*DefaultMongoStore)
	if !ok {
		t.Errorf(test.ExpectedBoolButFoundBool, true, ok)
	}
}
