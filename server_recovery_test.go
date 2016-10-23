package oauth2

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/phuc0302/go-oauth2/test"
	"github.com/phuc0302/go-oauth2/util"
)

func Test_recovery_Panic(t *testing.T) {
	// Setup server & test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := createRequestContext(r, w)
		defer recovery(request, true)
		panic(util.Status404())
	}))
	defer ts.Close()

	response, _ := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader("key=value"))
	if response.StatusCode != 404 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 404, response.StatusCode)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		response.Body.Close()
		fmt.Println(string(data))
	}
}
