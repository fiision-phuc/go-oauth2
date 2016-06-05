package oauth2

import (
	"testing"

	"github.com/phuc0302/go-oauth2/test"
)

func Test_GroupRole(t *testing.T) {
}

func Test_BindRole(t *testing.T) {
}

func Test_GroupRoute(t *testing.T) {
	objectFactory = &DefaultFactory{}
	router := DefaultRouter{}

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
	objectFactory = &DefaultFactory{}
	router := DefaultRouter{}

	// Test first bind
	router.BindRoute(GET, "/", func() {})
	if router.routes == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if len(router.routes) != 1 {
			t.Errorf(test.ExpectedNumberButFoundNumber, 1, len(router.routes))
		}
	}

	// Test second bind
	router.BindRoute(GET, "/sample", func() {})
	if len(router.routes) != 2 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 2, len(router.routes))
	}
}

func Test_MatchRoute(t *testing.T) {
	objectFactory = &DefaultFactory{}
	router := DefaultRouter{}

	router.GroupRoute(nil, "/user/profile", func(s *Server) {
		router.BindRoute(GET, "", func() {})

		router.BindRoute(GET, "/{profileID}", func() {})
		router.BindRoute(POST, "/{profileID}", func() {})
	})

	// [Test 1] Invalid path
	route, pathParams := router.MatchRoute(GET, "/user")
	if route != nil {
		t.Error(test.ExpectedNil)
	}
	if pathParams != nil {
		t.Error(test.ExpectedNil)
	}

	// [Test 2] Invalid HTTP method
	route, pathParams = router.MatchRoute(POST, "/user/profile")
	if route != nil {
		t.Error(test.ExpectedNil)
	}
	if pathParams != nil {
		t.Error(test.ExpectedNil)
	}

	// [Test 3] Valid HTTP method & path
	route, pathParams = router.MatchRoute(GET, "/user/profile")
	if route == nil {
		t.Error(test.ExpectedNotNil)
	}
	if pathParams != nil {
		t.Error(test.ExpectedNil)
	}

	// [Test 4] Valid HTTP method & path
	route, pathParams = router.MatchRoute(GET, "/user/profile/1")
	if route == nil {
		t.Error(test.ExpectedNotNil)
	}
	if pathParams == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if pathParams["profileID"] != "1" {
			t.Errorf(test.ExpectedStringButFoundString, "1", pathParams["profileID"])
		}
	}
}
