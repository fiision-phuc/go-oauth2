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
	r := createRouter()
	if r.roles != nil {
		t.Error(test.ExpectedNil)
	}

	r.groupRoles("/private(.htm[l]?)?**", "r_user", "r_admin")
	for rule, roles := range r.roles {
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

func Test_BindRole(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup router
	r := createRouter()
	r.bindRoute(Get, "/", func(request *RequestContext, security *OAuthContext) {})
	r.groupRoute(nil, "/user/profile(.htm[l]?)?", func(s *Server) {
		r.bindRoute(Get, "", func(request *RequestContext, security *OAuthContext) {})
		r.bindRoute(Post, "", func(request *RequestContext, security *OAuthContext) {})
		r.bindRoute(Get, "/{profileID}", func(request *RequestContext, security *OAuthContext) {})
	})
	r.groupRoute(nil, "/private", func(s *Server) {
		r.bindRoute(Get, "", func(request *RequestContext, security *OAuthContext) {})
		r.bindRoute(Get, "/{profileID}", func(request *RequestContext, security *OAuthContext) {})
	})
	r.groupRoles("/private**", "r_admin")

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		context := createRequestContext(request, writer)
		security := createSecurityContext(context)

		route, pathParams := r.matchRoute(context, security)
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
	router := createRouter()
	router.groupRoute(nil, "/user/profile", func(s *Server) {
		router.bindRoute(Get, "", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Get, "/{profileID}", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Post, "/{profileID}", func(request *RequestContext, security *OAuthContext) {})
	})

	if router.routes == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if len(router.routes) != 2 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 2, len(router.routes))
		} else {
			route := router.routes[1]
			if route.path != "/user/profile/{profileID}" {
				t.Errorf(test.ExpectedStringButFoundString, "/user/profile/{profileID}", route.path)
			}
			if route.handlers[Get] == nil {
				t.Error(test.ExpectedNotNil)
			}
			if route.handlers[Post] == nil {
				t.Error(test.ExpectedNotNil)
			}
		}
	}
}

func Test_BindRoute(t *testing.T) {
	//	u := new(TestUnit)
	//	defer u.Teardown()
	//	u.Setup()

	router := createRouter()

	// [Test 1] First bind
	router.bindRoute(Get, "/", func(c *RequestContext, s *OAuthContext) {})
	if len(router.routes) != 1 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 1, len(router.routes))
	}

	// [Test 2] Second bind
	router.bindRoute(Get, "/sample", func(c *RequestContext, s *OAuthContext) {})
	if len(router.routes) != 2 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 2, len(router.routes))
	}
}

func Test_MatchRoute_InvalidPath(t *testing.T) {
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup router
	router := createRouter()
	router.bindRoute(Get, "/", func(request *RequestContext, security *OAuthContext) {})
	router.groupRoute(nil, "/user/profile(.htm[l]?)?", func(s *Server) {
		router.bindRoute(Get, "", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Post, "", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Get, "/{profileID}", func(request *RequestContext, security *OAuthContext) {})
	})
	router.groupRoute(nil, "/private", func(s *Server) {
		router.bindRoute(Get, "", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Get, "/{profileID}", func(request *RequestContext, security *OAuthContext) {})
	})
	router.groupRoles("/private**", "r_admin")

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := createRequestContext(r, w)
		Security := createSecurityContext(context)

		route, pathParams := router.matchRoute(context, Security)
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
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup router
	router := createRouter()
	router.bindRoute(Get, "/", func(request *RequestContext, security *OAuthContext) {})
	router.groupRoute(nil, "/user/profile(.htm[l]?)?", func(s *Server) {
		router.bindRoute(Get, "", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Get, "/{profileID}", func(request *RequestContext, security *OAuthContext) {})
	})
	router.groupRoute(nil, "/private", func(s *Server) {
		router.bindRoute(Get, "", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Get, "/{profileID}", func(request *RequestContext, security *OAuthContext) {})
	})
	router.groupRoles("/private**", "r_admin")

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := createRequestContext(r, w)
		Security := createSecurityContext(context)

		route, pathParams := router.matchRoute(context, Security)
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
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup router
	router := createRouter()
	router.bindRoute(Get, "/", func(request *RequestContext, security *OAuthContext) {})
	router.groupRoute(nil, "/user/profile(.htm[l]?)?", func(s *Server) {
		router.bindRoute(Get, "", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Post, "", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Get, "/{profileID}", func(request *RequestContext, security *OAuthContext) {})
	})
	router.groupRoute(nil, "/private", func(s *Server) {
		router.bindRoute(Get, "", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Get, "/{profileID}", func(request *RequestContext, security *OAuthContext) {})
	})
	router.groupRoles("/private**", "r_admin")

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := createRequestContext(r, w)
		Security := createSecurityContext(context)

		route, _ := router.matchRoute(context, Security)
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
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup router
	router := createRouter()
	router.bindRoute(Get, "/", func(request *RequestContext, security *OAuthContext) {})
	router.groupRoute(nil, "/user/profile(.htm[l]?)?", func(s *Server) {
		router.bindRoute(Get, "", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Post, "", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Get, "/{profileID}", func(request *RequestContext, security *OAuthContext) {})
	})
	router.groupRoute(nil, "/private", func(s *Server) {
		router.bindRoute(Get, "", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Get, "/{profileID}", func(request *RequestContext, security *OAuthContext) {})
	})
	router.groupRoles("/private**", "r_admin")

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := createRequestContext(r, w)
		Security := createSecurityContext(context)

		route, pathParams := router.matchRoute(context, Security)
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
	u := new(TestUnit)
	defer u.Teardown()
	u.Setup()

	// Setup router
	router := createRouter()
	router.bindRoute(Get, "/", func(request *RequestContext, security *OAuthContext) {})
	router.groupRoute(nil, "/user/profile(.htm[l]?)?", func(s *Server) {
		router.bindRoute(Get, "", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Post, "", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Get, "/{profileID}", func(request *RequestContext, security *OAuthContext) {})
	})
	router.groupRoute(nil, "/private", func(s *Server) {
		router.bindRoute(Get, "", func(request *RequestContext, security *OAuthContext) {})
		router.bindRoute(Get, "/{profileID}", func(request *RequestContext, security *OAuthContext) {})
	})
	router.groupRoles("/private**", "r_admin")

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := createRequestContext(r, w)
		Security := createSecurityContext(context)

		route, _ := router.matchRoute(context, Security)
		if route == nil {
			t.Error(context.Path)
		}
	}))
	defer ts.Close()

	now := time.Now()
	token := store.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), now, now.Add(Cfg.AccessTokenDuration))

	http.Get(fmt.Sprintf("%s/private?access_token=%s", ts.URL, token.Token()))
	http.Get(fmt.Sprintf("%s/private/1?access_token=%s", ts.URL, token.Token()))
}
