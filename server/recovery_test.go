package server

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

func Test_recovery_development(t *testing.T) {
	// Setup server & test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer Recovery(w, r)
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

func Test_recovery_production(t *testing.T) {
	// Setup server & test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer Recovery(w, r)
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
