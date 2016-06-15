package oauth2

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/phuc0302/go-oauth2/test"
)

func Test_GroupRole(t *testing.T) {
	t.Error("Not yet implemented!")
}

func Test_BindRole(t *testing.T) {
	t.Error("Not yet implemented!")
}

func Test_GroupRoute(t *testing.T) {
	defer teardown()
	setup()

	router, _ := objectFactory.CreateRouter().(*DefaultRouter)
	router.GroupRoute(nil, "/user/profile", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
		router.BindRoute(POST, "/{profileID}", func() {})
	})

	if router.routes == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if len(router.routes) != 2 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 2, len(router.routes))
		} else {
			r := router.routes[1]
			route, _ := r.(*DefaultRoute)

			if route.path != "/user/profile/{profileID}" {
				t.Errorf(test.ExpectedStringButFoundString, "/user/profile/{profileID}", route.path)
			}
			if route.handlers[GET] == nil {
				t.Error(test.ExpectedNotNil)
			}
			if route.handlers[POST] == nil {
				t.Error(test.ExpectedNotNil)
			}
		}
	}
}

func Test_BindRoute(t *testing.T) {
	defer teardown()
	setup()

	router, _ := objectFactory.CreateRouter().(*DefaultRouter)

	// [Test 1] First bind
	router.BindRoute(GET, "/", func() {})
	if router.routes == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if len(router.routes) != 1 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 1, len(router.routes))
		}
	}

	// [Test 2] Second bind
	router.BindRoute(GET, "/sample", func() {})
	if len(router.routes) != 2 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 2, len(router.routes))
	}
}

func Test_MatchRoute(t *testing.T) {
	defer teardown()
	setup()

	// Setup router
	router := objectFactory.CreateRouter()
	router.GroupRoute(nil, "/user/profile", func(s *Server) {
		router.BindRoute(GET, "", func() {})

		router.BindRoute(GET, "/{profileID}", func() {})
		router.BindRoute(POST, "/{profileID}", func() {})
	})
	router.GroupRoute(nil, "/private", func(s *Server) {
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRole(nil, "/private**", "r_admin")

	// Setup test server
	testCase := ""
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		Security := objectFactory.CreateSecurityContext(context)

		if route, pathParams := router.MatchRoute(context, Security); route != nil {
			switch testCase {

			case "Test 5":
				if route != nil {
					t.Error(test.ExpectedNil)
				}
				if pathParams != nil {
					t.Error(test.ExpectedNil)
				}
				break
			}

			context.PathParams = pathParams
			route.InvokeHandler(context, Security)
		} else {
			switch testCase {

			case "Test 1", "Test 2":
				if route != nil {
					t.Error(test.ExpectedNil)
				}
				if pathParams != nil {
					t.Error(test.ExpectedNil)
				}
				break

			case "Test 3":
				if route == nil {
					t.Error(test.ExpectedNotNil)
				}
				if pathParams != nil {
					t.Error(test.ExpectedNil)
				}
				break

			case "Test 4":
				if route == nil {
					t.Error(test.ExpectedNotNil)
				}
				if pathParams == nil {
					t.Error(test.ExpectedNil)
				} else {
					if pathParams["profileID"] != "1" {
						t.Errorf(test.ExpectedStringButFoundString, "1", pathParams["profileID"])
					}
				}
				break

			case "Test 6":
				if route == nil {
					t.Error(test.ExpectedNotNil)
				}
				if pathParams == nil {
					t.Error(test.ExpectedNil)
				}
				break
			}
		}
	}))
	defer ts.Close()

	// [Test 1] Invalid path
	testCase = "Test 1"
	http.Get(fmt.Sprintf("%s/user", ts.URL))

	// [Test 2] Invalid HTTP method
	testCase = "Test 2"
	http.Post(fmt.Sprintf("%s/user/profile", ts.URL), "application/x-www-form-urlencoded", nil)

	// [Test 3] Valid HTTP method & path
	testCase = "Test 3"
	http.Get(fmt.Sprintf("%s/user/profile", ts.URL))

	// [Test 4] Valid HTTP method & path
	testCase = "Test 4"
	http.Get(fmt.Sprintf("%s/user/profile/1", ts.URL))

	// [Test 5] Send request to secure resource without access_token
	testCase = "Test 5"
	http.Get(fmt.Sprintf("%s/private/profile/1", ts.URL))

	// [Test 6] Send request to secure resource with access_token
	testCase = "Test 5"
	now := time.Now()
	token := tokenStore.CreateAccessToken(clientID.Hex(), userID.Hex(), now, now.Add(cfg.AccessTokenDuration))
	http.Get(fmt.Sprintf("%s/private/profile/1?access_token=%s", ts.URL, token.Token()))
}
