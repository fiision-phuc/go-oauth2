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
func Test_BindHandlerWithPanic(t *testing.T) {
	route := DefaultRoute{}

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

	defer func() {
		if r := recover(); r != nil {
			/* Expected panic */
		}
	}()
	route.BindHandler(GET, func() {
		panic("Test if func had been invoked or not.")
	})
	route.InvokeHandler(nil, nil)
	t.Errorf(test.ExpectedPanic)
}

func Test_URLPattern(t *testing.T) {
	defer teardown()
	setup()

	route := objectFactory.CreateRoute("/example/{userID}")
	if route.URLPattern() != "/example/{userID}" {
		t.Errorf(test.ExpectedStringButFoundString, "/example/{userID}", route.URLPattern())
	}
}

func Test_MatchURLPattern_InvalidHTTPMethod(t *testing.T) {
	defer teardown()
	setup()

	route := objectFactory.CreateRoute("/example/{userID}/profile/{profileID}")
	route.BindHandler(GET, func() {})

	isMatched, pathParams := route.MatchURLPattern(POST, "/example/1")
	if isMatched {
		t.Errorf(test.ExpectedBoolButFoundBool, false, isMatched)
	}
	if pathParams != nil {
		t.Error(test.ExpectedNil)
	}
}

func Test_MatchURLPattern_InvalidHTTPMethodButInvalidPath(t *testing.T) {
	defer teardown()
	setup()

	route := objectFactory.CreateRoute("/example/{userID}/profile/{profileID}")
	route.BindHandler(GET, func() {})

	isMatched, pathParams := route.MatchURLPattern(GET, "/example/1/profile")
	if isMatched {
		t.Errorf(test.ExpectedBoolButFoundBool, false, isMatched)
	}
	if pathParams != nil {
		t.Error(test.ExpectedNil)
	}
}

func Test_MatchURLPattern_ValidHTTPMethodAndValidPath(t *testing.T) {
	defer teardown()
	setup()

	route := objectFactory.CreateRoute("/example/{userID}/profile/{profileID}")
	route.BindHandler(GET, func() {})

	isMatched, pathParams := route.MatchURLPattern(GET, "/example/1/profile/1")
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
