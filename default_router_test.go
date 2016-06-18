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

	if router, ok := objectFactory.CreateRouter().(*DefaultRouter); !ok {
		t.Errorf(test.ExpectedBoolButFoundBool, true, ok)
	} else {
		// [Test 1] First bind
		router.BindRoute(GET, "/", func() {})
		if len(router.routes) != 1 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 1, len(router.routes))
		}

		// [Test 2] Second bind
		router.BindRoute(GET, "/sample", func() {})
		if len(router.routes) != 2 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 2, len(router.routes))
		}
	}
}

func Test_MatchRoute_InvalidPath(t *testing.T) {
	defer teardown()
	setup()

	// Setup router
	router := objectFactory.CreateRouter()
	router.GroupRoute(nil, "/user/profile", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoute(nil, "/private", func(s *Server) {
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRole(nil, "/private**", "r_admin")

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		Security := objectFactory.CreateSecurityContext(context)

		route, pathParams := router.MatchRoute(context, Security)
		if route != nil {
			t.Error(test.ExpectedNil)
		}
		if pathParams != nil {
			t.Error(test.ExpectedNil)
		}
	}))
	defer ts.Close()
	http.Get(fmt.Sprintf("%s/user", ts.URL))
}

func Test_MatchRoute_InvalidHTTPMethod(t *testing.T) {
	defer teardown()
	setup()

	// Setup router
	router := objectFactory.CreateRouter()
	router.GroupRoute(nil, "/user/profile", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoute(nil, "/private", func(s *Server) {
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRole(nil, "/private**", "r_admin")

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		Security := objectFactory.CreateSecurityContext(context)

		route, pathParams := router.MatchRoute(context, Security)
		if route != nil {
			t.Error(test.ExpectedNil)
		}
		if pathParams != nil {
			t.Error(test.ExpectedNil)
		}
	}))
	defer ts.Close()
	http.Post(fmt.Sprintf("%s/user/profile", ts.URL), "application/x-www-form-urlencoded", nil)
}

func Test_MatchRoute_ValidHTTPMethodAndPath(t *testing.T) {
	defer teardown()
	setup()

	// Setup router
	router := objectFactory.CreateRouter()
	router.BindRoute(GET, "/", func() {})
	router.GroupRoute(nil, "/user/profile", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoute(nil, "/private", func(s *Server) {
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRole(nil, "/private**", "r_admin")

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		Security := objectFactory.CreateSecurityContext(context)

		route, pathParams := router.MatchRoute(context, Security)
		if route == nil {
			t.Error(test.ExpectedNotNil)
		}
		if pathParams != nil {
			t.Error(test.ExpectedNil)
		}
	}))
	defer ts.Close()
	fmt.Println("----------", fmt.Sprintf("%s/user/profile", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile", ts.URL))
}

func Test_MatchRoute_SendRequestToSecureResourceWithoutAccessToken(t *testing.T) {
	defer teardown()
	setup()

	// Setup router
	router := objectFactory.CreateRouter()
	router.GroupRoute(nil, "/user/profile", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoute(nil, "/private", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRole(nil, "/private**", "r_admin")

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		Security := objectFactory.CreateSecurityContext(context)

		route, pathParams := router.MatchRoute(context, Security)
		if route != nil {
			t.Error(test.ExpectedNil)
		}
		if pathParams != nil {
			t.Error(test.ExpectedNil)
		}
	}))
	defer ts.Close()
	http.Get(fmt.Sprintf("%s/private", ts.URL))
	//	http.Get(fmt.Sprintf("%s/private/", ts.URL))
	//	http.Get(fmt.Sprintf("%s/private/1", ts.URL))
}

func Test_MatchRoute_SendRequestToSecureResourceWithAccessToken(t *testing.T) {
	defer teardown()
	setup()

	// Setup router
	router := objectFactory.CreateRouter()
	router.GroupRoute(nil, "/user/profile", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(POST, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoute(nil, "/private", func(s *Server) {
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRole(nil, "/private**", "r_admin")

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		Security := objectFactory.CreateSecurityContext(context)

		route, pathParams := router.MatchRoute(context, Security)
		if route == nil {
			t.Error(test.ExpectedNotNil)
		}
		if pathParams == nil {
			t.Error(test.ExpectedNil)
		}
	}))
	defer ts.Close()

	now := time.Now()
	token := TokenStore.CreateAccessToken(clientID.Hex(), userID.Hex(), now, now.Add(Cfg.AccessTokenDuration))

	http.Get(fmt.Sprintf("%s/private?access_token=%s", ts.URL, token.Token()))
	http.Get(fmt.Sprintf("%s/private/?access_token=%s", ts.URL, token.Token()))
	http.Get(fmt.Sprintf("%s/private/1?access_token=%s", ts.URL, token.Token()))
}
