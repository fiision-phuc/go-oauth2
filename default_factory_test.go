package oauth2

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/phuc0302/go-oauth2/test"
	"github.com/phuc0302/go-oauth2/utils"
)

func Test_CreateRequestContext(t *testing.T) {
	cfg = loadConfig(debug)
	defer os.Remove(debug)

	objectFactory = &DefaultFactory{}
	testCase := ""

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)

		switch testCase {

		case "Test 1":
			if context.Path != utils.FormatPath(r.URL.Path) {
				t.Errorf(test.ExpectedStringButFoundString, utils.FormatPath(r.URL.Path), context.Path)
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
			break

		case "Test 2":
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
			break

		case "Test 3":
			if context.Header["content-type"] != "application/x-www-form-urlencoded" {
				t.Errorf(test.ExpectedStringButFoundString, "application/x-www-form-urlencoded", context.Header["content-type"])
			}
			if context.QueryParams != nil {
				t.Error(test.ExpectedNil)
			}
			break

		case "Test 4":
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
			break

		case "Test 5":
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
			break

		default:
			break
		}
	}))
	defer ts.Close()

	// [Test 1] Send Get request
	testCase = "Test 1"
	http.Get(ts.URL)

	// [Test 2] Send Get request with query params
	testCase = "Test 2"
	http.Get(fmt.Sprintf("%s?userID=1&profileID=2", ts.URL))

	// [Test 3] Send Post request
	testCase = "Test 3"
	http.Post(ts.URL, strings.ToUpper("application/x-www-form-urlencoded"), nil)

	// [Test 4] Send Post request with data
	testCase = "Test 4"
	http.Post(ts.URL, strings.ToUpper("application/x-www-form-urlencoded"), strings.NewReader("userID=1&profileID=2"))

	// [Test 5] Send Post request with multipart data
	testCase = "Test 5"
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

func Test_CreateSecurityContext(t *testing.T) {
	t.Error("There is no test case available yet!")
}

func Test_CreateRoute(t *testing.T) {
	objectFactory = &DefaultFactory{}
	route := generateRoute("/example/{userID}/profile/{profileID}")
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
