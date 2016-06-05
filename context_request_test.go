package oauth2

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_OutputHeader(t *testing.T) {
	objectFactory = &DefaultFactory{}

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)

	}))
	defer ts.Close()
}
