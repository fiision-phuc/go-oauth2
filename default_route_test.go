package oauth2

import (
	"testing"

	"github.com/phuc0302/go-oauth2/test"
)

func Test_BindHandler(t *testing.T) {
	route := DefaultRoute{}
	route.BindHandler(GET, func() {})

	if route.handlers == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if route.handlers[GET] == nil {
			t.Error(test.ExpectedNotNil)
		}
	}

	// [Test 2] Expected panic if not a func
	defer func() {
		if r := recover(); r != nil {
			/* Expected panic */
		}
	}()
	route.BindHandler(POST, "")
	t.Errorf(test.ExpectedPanic)
}

func Test_InvokeHandler(t *testing.T) {
	defer teardown()
	setup()
	route := objectFactory.CreateRoute("/example/{userID}/profile/{profileID}")
	route.BindHandler(GET, func() {
		panic("Test if func had been invoked or not.")
	})

	defer func() {
		if r := recover(); r != nil {
			/* Expected panic */
		}
	}()
	route.InvokeHandler(nil, nil)
	t.Errorf(test.ExpectedPanic)
}

func Test_URLPattern(t *testing.T) {
	route := DefaultRoute{
		path: "/example/{userID}",
	}

	if route.URLPattern() != "/example/{userID}" {
		t.Errorf(test.ExpectedStringButFoundString, "/example/{userID}", route.URLPattern())
	}
}

func Test_MatchURLPattern(t *testing.T) {
	defer teardown()
	setup()
	route := objectFactory.CreateRoute("/example/{userID}/profile/{profileID}")
	route.BindHandler(GET, func() {})

	// [Test 1] Invalid HTTP method
	isMatched, pathParams := route.MatchURLPattern(POST, "/example/1")
	if isMatched {
		t.Errorf(test.ExpectedBoolButFoundBool, false, isMatched)
	}
	if pathParams != nil {
		t.Error(test.ExpectedNil)
	}

	// [Test 2] Valid HTTP method but invalid path
	isMatched, pathParams = route.MatchURLPattern(GET, "/example/1/profile")
	if isMatched {
		t.Errorf(test.ExpectedBoolButFoundBool, false, isMatched)
	}
	if pathParams != nil {
		t.Error(test.ExpectedNil)
	}

	// [Test 3] Valid HTTP method and valid path
	isMatched, pathParams = route.MatchURLPattern(GET, "/example/1/profile/1")
	if !isMatched {
		t.Errorf(test.ExpectedBoolButFoundBool, true, isMatched)
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
