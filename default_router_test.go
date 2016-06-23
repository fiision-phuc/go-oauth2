package oauth2

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/phuc0302/go-oauth2/test"
)

func Test_GroupRoles(t *testing.T) {
	u := new(UnitTest)
	defer u.Teardown()
	u.Setup()

	router, _ := objectFactory.CreateRouter().(*DefaultRouter)
	if router == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if router.userRoles != nil {
			t.Error(test.ExpectedNil)
		}

		router.GroupRoles("/private(.htm[l]?)?**", "r_user", "r_admin")

		for rule, roles := range router.userRoles {
			if !rule.MatchString("/private") {
				t.Errorf(test.ExpectedBoolButFoundBool, true, rule.MatchString("/private"))
			}
			if !rule.MatchString("/private.htm") {
				t.Errorf(test.ExpectedBoolButFoundBool, true, rule.MatchString("/private.htm"))
			}
			if !rule.MatchString("/private.html") {
				t.Errorf(test.ExpectedBoolButFoundBool, true, rule.MatchString("/private.html"))
			}
			if !rule.MatchString("/private/1/2/3") {
				t.Errorf(test.ExpectedBoolButFoundBool, true, rule.MatchString("/private.html/1/2/3"))
			}
			if !rule.MatchString("/private.htm/1/2/3") {
				t.Errorf(test.ExpectedBoolButFoundBool, true, rule.MatchString("/private.html/1/2/3"))
			}
			if !rule.MatchString("/private.html/1/2/3") {
				t.Errorf(test.ExpectedBoolButFoundBool, true, rule.MatchString("/private.html/1/2/3"))
			}

			if !roles.MatchString("r_user") {
				t.Errorf(test.ExpectedBoolButFoundBool, true, roles.MatchString("r_user"))
			}
			if !roles.MatchString("r_admin") {
				t.Errorf(test.ExpectedBoolButFoundBool, true, roles.MatchString("r_admin"))
			}
		}
	}
}

func Test_BindRole(t *testing.T) {
	u := new(UnitTest)
	defer u.Teardown()
	u.Setup()

	// Setup router
	router := objectFactory.CreateRouter()

	router.BindRoute(GET, "/", func() {})
	router.GroupRoute(nil, "/user/profile(.htm[l]?)?", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(POST, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoute(nil, "/private", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoles("/private**", "r_admin")

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

	t.Error("Not yet implemented!")
}

func Test_GroupRoute(t *testing.T) {
	u := new(UnitTest)
	defer u.Teardown()
	u.Setup()

	router := new(DefaultRouter)
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
	u := new(UnitTest)
	defer u.Teardown()
	u.Setup()

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
	u := new(UnitTest)
	defer u.Teardown()
	u.Setup()

	// Setup router
	router := objectFactory.CreateRouter()

	router.BindRoute(GET, "/", func() {})
	router.GroupRoute(nil, "/user/profile(.htm[l]?)?", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(POST, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoute(nil, "/private", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoles("/private**", "r_admin")

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
	u := new(UnitTest)
	defer u.Teardown()
	u.Setup()

	// Setup router
	router := objectFactory.CreateRouter()

	router.BindRoute(GET, "/", func() {})
	router.GroupRoute(nil, "/user/profile(.htm[l]?)?", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoute(nil, "/private", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoles("/private**", "r_admin")

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
	u := new(UnitTest)
	defer u.Teardown()
	u.Setup()

	// Setup router
	router := objectFactory.CreateRouter()

	router.BindRoute(GET, "/", func() {})
	router.GroupRoute(nil, "/user/profile(.htm[l]?)?", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(POST, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoute(nil, "/private", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoles("/private**", "r_admin")

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		Security := objectFactory.CreateSecurityContext(context)

		route, _ := router.MatchRoute(context, Security)
		if route == nil {
			t.Error(context.Path)
		}
	}))
	defer ts.Close()
	http.Get(ts.URL)
	http.Get(fmt.Sprintf("%s/user/profile", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile?userID=1", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile/", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile/?userID=1", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile/1", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile/1?userID=1", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile/1/", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile/1/?userID=1", ts.URL))

	http.Get(fmt.Sprintf("%s/user/profile.htm", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile.htm?userID=1", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile.htm/", ts.URL))
	http.Get(fmt.Sprintf("%s/user/profile.html/?userID=1", ts.URL))
}

func Test_MatchRoute_SendRequestToSecureResourceWithoutAccessToken(t *testing.T) {
	u := new(UnitTest)
	defer u.Teardown()
	u.Setup()

	// Setup router
	router := objectFactory.CreateRouter()

	router.BindRoute(GET, "/", func() {})
	router.GroupRoute(nil, "/user/profile(.htm[l]?)?", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(POST, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoute(nil, "/private", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoles("/private**", "r_admin")

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		Security := objectFactory.CreateSecurityContext(context)

		route, pathParams := router.MatchRoute(context, Security)
		if route != nil {
			t.Error(context.Path)
		}
		if pathParams != nil {
			t.Error(context.Path)
		}
	}))
	defer ts.Close()
	http.Get(fmt.Sprintf("%s/private", ts.URL))
	http.Get(fmt.Sprintf("%s/private/", ts.URL))
	http.Get(fmt.Sprintf("%s/private/1", ts.URL))
	http.Get(fmt.Sprintf("%s/private/1/", ts.URL))
}

func Test_MatchRoute_SendRequestToSecureResourceWithAccessToken(t *testing.T) {
	u := new(UnitTest)
	defer u.Teardown()
	u.Setup()

	// Setup router
	router := objectFactory.CreateRouter()

	router.BindRoute(GET, "/", func() {})
	router.GroupRoute(nil, "/user/profile(.htm[l]?)?", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(POST, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoute(nil, "/private", func(s *Server) {
		router.BindRoute(GET, "", func() {})
		router.BindRoute(GET, "/{profileID}", func() {})
	})
	router.GroupRoles("/private**", "r_admin")

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)
		Security := objectFactory.CreateSecurityContext(context)

		route, _ := router.MatchRoute(context, Security)
		if route == nil {
			t.Error(context.Path)
		}
	}))
	defer ts.Close()

	now := time.Now()
	token := TokenStore.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), now, now.Add(Cfg.AccessTokenDuration))

	http.Get(fmt.Sprintf("%s/private?access_token=%s", ts.URL, token.Token()))
	http.Get(fmt.Sprintf("%s/private/1?access_token=%s", ts.URL, token.Token()))
}
