package oauth2

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

	// Setup test server
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

	//	// [Test 1] Invalid path
	//	route, pathParams := router.MatchRoute(GET, "/user")
	//	if route != nil {
	//		t.Error(test.ExpectedNil)
	//	}
	//	if pathParams != nil {
	//		t.Error(test.ExpectedNil)
	//	}

	//	// [Test 2] Invalid HTTP method
	//	route, pathParams = router.MatchRoute(POST, "/user/profile")
	//	if route != nil {
	//		t.Error(test.ExpectedNil)
	//	}
	//	if pathParams != nil {
	//		t.Error(test.ExpectedNil)
	//	}

	//	// [Test 3] Valid HTTP method & path
	//	route, pathParams = router.MatchRoute(GET, "/user/profile")
	//	if route == nil {
	//		t.Error(test.ExpectedNotNil)
	//	}
	//	if pathParams != nil {
	//		t.Error(test.ExpectedNil)
	//	}

	//	// [Test 4] Valid HTTP method & path
	//	route, pathParams = router.MatchRoute(GET, "/user/profile/1")
	//	if route == nil {
	//		t.Error(test.ExpectedNotNil)
	//	}
	//	if pathParams == nil {
	//		t.Error(test.ExpectedNotNil)
	//	} else {
	//		if pathParams["profileID"] != "1" {
	//			t.Errorf(test.ExpectedStringButFoundString, "1", pathParams["profileID"])
	//		}
	//	}
}
