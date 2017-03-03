package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/phuc0302/go-oauth2/test"
)

func Test_GroupRoute(t *testing.T) {
	router := DefaultRouter()
	router.GroupRoute(nil, "/user/profile", func(s *Server) {
		router.BindRoute(Get, "", func(request *RequestContext) {})
		router.BindRoute(Get, "/{profileID}", func(request *RequestContext) {})
		router.BindRoute(Post, "/{profileID}", func(request *RequestContext) {})
	})

	if router.routes == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if len(router.routes) != 2 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 2, len(router.routes))
		} else {
			route0 := router.routes[0]
			if route0.regex.String() != "^/user/profile/?$" {
				t.Errorf(test.ExpectedStringButFoundString, "^/user/profile/?$", route0.regex.String())
			}
			if route0.handlers[Get] == nil {
				t.Error(test.ExpectedNotNil)
			}

			route1 := router.routes[1]
			if route1.regex.String() != "^/user/profile/(?P<profileID>[^/#?]+)/?$" {
				t.Errorf(test.ExpectedStringButFoundString, "^/user/profile/(?P<profileID>[^/#?]+)/?$", route1.regex.String())
			}
			if route1.handlers[Get] == nil {
				t.Error(test.ExpectedNotNil)
			}
			if route1.handlers[Post] == nil {
				t.Error(test.ExpectedNotNil)
			}
		}
	}
}

func Test_BindRoute(t *testing.T) {
	router := DefaultRouter()

	// [Test 1] First bind
	router.BindRoute(Get, "/", func(c *RequestContext) {})
	if len(router.routes) != 1 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 1, len(router.routes))
	}

	// [Test 2] Second bind
	router.BindRoute(Get, "/sample", func(c *RequestContext) {})
	if len(router.routes) != 2 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 2, len(router.routes))
	}
}

func Test_MatchRoute_InvalidPath(t *testing.T) {
	// Setup router
	router := DefaultRouter()
	router.BindRoute(Get, "/", func(request *RequestContext) {})
	router.GroupRoute(nil, "/user/profile(.htm[l]?)?", func(s *Server) {
		router.BindRoute(Get, "", func(request *RequestContext) {})
		router.BindRoute(Post, "", func(request *RequestContext) {})
		router.BindRoute(Get, "/{profileID}", func(request *RequestContext) {})
	})
	router.GroupRoute(nil, "/private", func(s *Server) {
		router.BindRoute(Get, "", func(request *RequestContext) {})
		router.BindRoute(Get, "/{profileID}", func(request *RequestContext) {})
	})

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := httprouter.CleanPath(r.URL.Path)
		method := strings.ToLower(r.Method)

		route, pathParams := router.MatchRoute(method, path)
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
	// Setup router
	router := DefaultRouter()
	router.BindRoute(Get, "/", func(request *RequestContext) {})
	router.GroupRoute(nil, "/user/profile(.htm[l]?)?", func(s *Server) {
		router.BindRoute(Get, "", func(request *RequestContext) {})
		router.BindRoute(Get, "/{profileID}", func(request *RequestContext) {})
	})
	router.GroupRoute(nil, "/private", func(s *Server) {
		router.BindRoute(Get, "", func(request *RequestContext) {})
		router.BindRoute(Get, "/{profileID}", func(request *RequestContext) {})
	})

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := httprouter.CleanPath(r.URL.Path)
		method := strings.ToLower(r.Method)

		route, pathParams := router.MatchRoute(method, path)
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
	// Setup router
	router := DefaultRouter()
	router.BindRoute(Get, "/", func(request *RequestContext) {})
	router.GroupRoute(nil, "/user/profile(.htm[l]?)?", func(s *Server) {
		router.BindRoute(Get, "", func(request *RequestContext) {})
		router.BindRoute(Post, "", func(request *RequestContext) {})
		router.BindRoute(Get, "/{profileID}", func(request *RequestContext) {})
	})
	router.GroupRoute(nil, "/private", func(s *Server) {
		router.BindRoute(Get, "", func(request *RequestContext) {})
		router.BindRoute(Get, "/{profileID}", func(request *RequestContext) {})
	})

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := httprouter.CleanPath(r.URL.Path)
		method := strings.ToLower(r.Method)

		route, _ := router.MatchRoute(method, path)
		if route == nil {
			t.Error(path)
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
