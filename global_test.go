package oauth2

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/phuc0302/go-oauth2/oauth_key"
	"github.com/phuc0302/go-oauth2/server"
	"github.com/phuc0302/go-oauth2/test"
	"github.com/phuc0302/go-oauth2/util"
)

func Test_ValidateToken_NoAccessToken(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				/* Expected panic */
			}
		}()

		context := server.CreateContext(w, r)
		f1 := ValidateToken()

		f2 := f1(func(r *server.RequestContext) {})
		f2(context)

		t.Errorf(test.ExpectedPanic)
	}))
	defer ts.Close()
	http.Get(ts.URL)
}

func Test_ValidateToken_WithBasicAuth(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := server.CreateContext(w, r)
		f1 := ValidateToken()

		f2 := f1(func(r *server.RequestContext) {})
		f2(context)

		if security, ok := context.GetExtra(oauthKey.OAuthContext).(*OAuthContext); ok {
			if security.Client == nil {
				t.Error(test.ExpectedNotNil)
			}
			if security.User == nil {
				t.Error(test.ExpectedNotNil)
			}
		} else {
			t.Error(test.ExpectedNotNil)
		}
	}))
	defer ts.Close()

	// Send token as query param
	request, _ := http.NewRequest("Get", ts.URL, nil)
	request.SetBasicAuth(u.ClientID.Hex(), u.ClientSecret.Hex())

	http.DefaultClient.Do(request)
}

func Test_ValidateToken_WithGetAccessToken(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := server.CreateContext(w, r)
		f1 := ValidateToken()

		f2 := f1(func(r *server.RequestContext) {})
		f2(context)

		if security, ok := context.GetExtra(oauthKey.OAuthContext).(*OAuthContext); ok {
			if security.AccessToken == nil {
				t.Error(test.ExpectedNotNil)
			}
			if security.RefreshToken != nil {
				t.Error(test.ExpectedNil)
			}
			if security.Client == nil {
				t.Error(test.ExpectedNotNil)
			}
			if security.User == nil {
				t.Error(test.ExpectedNotNil)
			}
		} else {
			t.Error(test.ExpectedNotNil)
		}
	}))
	defer ts.Close()

	// Generate token
	now := time.Now()
	token := Store.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), now, now.Add(cfg.AccessTokenDuration))

	// Send token as query param
	http.Get(fmt.Sprintf("%s?access_token=%s", ts.URL, token.Token()))
}

func Test_ValidateToken_WithPostAccessToken(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := server.CreateContext(w, r)
		f1 := ValidateToken()

		f2 := f1(func(r *server.RequestContext) {})
		f2(context)

		if security, ok := context.GetExtra(oauthKey.OAuthContext).(*OAuthContext); ok {
			if security.AccessToken == nil {
				t.Error(test.ExpectedNotNil)
			}
			if security.RefreshToken != nil {
				t.Error(test.ExpectedNil)
			}
			if security.Client == nil {
				t.Error(test.ExpectedNotNil)
			}
			if security.User == nil {
				t.Error(test.ExpectedNotNil)
			}
		} else {
			t.Error(test.ExpectedNotNil)
		}
	}))
	defer ts.Close()

	// Generate token
	now := time.Now().UTC()
	token := Store.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), now, now.Add(cfg.AccessTokenDuration))

	// Send token as authorization header
	request, _ := http.NewRequest("POST", ts.URL, nil)
	client := http.DefaultClient

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token()))
	client.Do(request)
}

func Test_ValidateRoles_InvalidRoles(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err, ok := recover().(*util.Status); ok {
				if err.Code != 401 {
					t.Errorf(test.ExpectedNumberButFoundNumber, 401, err.Code)
				}
			} else {
				t.Error(test.ExpectedNotNil)
			}
		}()

		context := server.CreateContext(w, r)
		f1 := func(c *server.RequestContext) {
			data := map[string]string{"key": "Hello world"}
			c.OutputJSON(util.Status200(), data)
		}

		f1 = server.Adapt(f1, ValidateToken(), ValidateRoles("r_android"))
		f1(context)
	}))
	defer ts.Close()

	// Generate token
	now := time.Now().UTC()
	token := Store.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), now, now.Add(cfg.AccessTokenDuration))

	// Send token as authorization header
	request, _ := http.NewRequest("POST", ts.URL, nil)
	client := http.DefaultClient

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token()))
	client.Do(request)
}

func Test_ValidateRoles_ValidRoles(t *testing.T) {
	u := new(TestEnv)
	defer u.Teardown()
	u.Setup()

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := server.CreateContext(w, r)
		f1 := func(c *server.RequestContext) {
			c.OutputError(util.Status200())
		}

		f1 = server.Adapt(f1, ValidateToken(), ValidateRoles("r_admin"))
		f1(context)
	}))
	defer ts.Close()

	// Generate token
	now := time.Now().UTC()
	token := Store.CreateAccessToken(u.ClientID.Hex(), u.UserID.Hex(), now, now.Add(cfg.AccessTokenDuration))

	// Send token as authorization header
	request, _ := http.NewRequest("POST", ts.URL, nil)
	client := http.DefaultClient

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token()))
	response, _ := client.Do(request)

	status := util.ParseStatus(response)
	if status.Code != 200 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 200, status.Code)
	}
}
