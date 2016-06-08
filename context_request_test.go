package oauth2

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phuc0302/go-oauth2/test"
	"github.com/phuc0302/go-oauth2/utils"
)

func Test_OutputError(t *testing.T) {
	objectFactory = &DefaultFactory{}

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := objectFactory.CreateRequestContext(r, w)

		context.OutputHeader("test-header", "test-header-value")
		context.OutputError(utils.Status400())
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", nil)
	if response.Header.Get("test-header") != "test-header-value" {
		t.Errorf(test.ExpectedStringButFoundString, "test-header-value", response.Header.Get("test-header"))
	}
	if response.StatusCode != 400 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 400, response.StatusCode)
	}

	bytes, _ := ioutil.ReadAll(response.Body)
	if string(bytes) != "{\"status\":400,\"error\":\"Bad Request\",\"error_description\":\"Bad Request\"}" {
		t.Errorf(test.ExpectedStringButFoundString, "{\"status\":400,\"error\":\"Bad Request\",\"error_description\":\"Bad Request\"}", string(bytes))
	}
}
