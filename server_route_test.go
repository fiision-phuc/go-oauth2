package oauth2

import (
	"testing"

	"github.com/phuc0302/go-oauth2/test"
)

func Test_bindHandler(t *testing.T) {
	// [Test 1] Create new route empty string
	route := createRoute("")
	route.bindHandler(Get, func(request *Request, security *Security) {})
	if route.handlers == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if route.path != "" {
			t.Errorf(test.ExpectedStringButFoundString, "", route.path)
		}
		if route.handlers[Get] == nil {
			t.Error(test.ExpectedNotNil)
		}
	}

	// [Test 2] Create new route with non empty string
	route = createRoute("/example/{userID}")
	route.bindHandler(Get, func(request *Request, security *Security) {})
	if route.path != "/example/{userID}" {
		t.Errorf(test.ExpectedStringButFoundString, "/example/{userID}", route.path)
	}
	matched, params := route.match(Get, "/example/1")
	if !matched {
		t.Errorf(test.ExpectedBoolButFoundBool, true, matched)
	}
	if len(params) != 1 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 1, len(params))
	} else {
		if params["userID"] != "1" {
			t.Errorf(test.ExpectedStringButFoundString, "1", params["userID"])
		}
	}
}
func Test_bindHandlerWithPanic(t *testing.T) {
	route := createRoute("")

	defer func() {
		if r := recover(); r != nil {
			/* Expected panic */
		}
	}()
	route.bindHandler(Post, nil)
	t.Errorf(test.ExpectedPanic)
}

func Test_invokeHandler(t *testing.T) {
	route := createRoute("/example/{userID}/profile/{profileID}")
	defer func() {
		if r := recover(); r != nil {
			/* Expected panic */
		}
	}()
	route.bindHandler(Get, func(request *Request, security *Security) {
		panic("Test if func had been invoked or not.")
	})
	route.invokeHandler(nil, nil)
	t.Errorf(test.ExpectedPanic)
}

func Test_match_InvalidHTTPMethod(t *testing.T) {
	route := createRoute("/example/{userID}/profile/{profileID}")
	route.bindHandler(Get, func(request *Request, security *Security) {})

	matched, pathParams := route.match(Post, "/example/1")
	if matched {
		t.Errorf(test.ExpectedBoolButFoundBool, false, matched)
	}
	if pathParams != nil {
		t.Error(test.ExpectedNil)
	}
}

func Test_match_InvalidHTTPMethodAndInvalidPath(t *testing.T) {
	route := createRoute("/example/{userID}/profile/{profileID}")
	route.bindHandler(Get, func(request *Request, security *Security) {})

	matched, pathParams := route.match(Get, "/example/1/profile")
	if matched {
		t.Errorf(test.ExpectedBoolButFoundBool, false, matched)
	}
	if pathParams != nil {
		t.Error(test.ExpectedNil)
	}
}

func Test_match_ValidHTTPMethodAndValidPath(t *testing.T) {
	route := createRoute("/example/{userID}/profile/{profileID}")
	route.bindHandler(Get, func(request *Request, security *Security) {})

	matched, pathParams := route.match(Get, "/example/1/profile/1")
	if !matched {
		t.Errorf(test.ExpectedBoolButFoundBool, true, matched)
	}
	if pathParams == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if pathParams["userID"] != "1" {
			t.Errorf(test.ExpectedStringButFoundString, "1", pathParams["userID"])
		}
		if pathParams["profileID"] != "1" {
			t.Errorf(test.ExpectedStringButFoundString, "1", pathParams["profileID"])
		}
	}
}
