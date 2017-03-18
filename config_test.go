package oauth2

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/phuc0302/go-oauth2/oauth_key"
	"github.com/phuc0302/go-server"
	"github.com/phuc0302/go-server/expected_format"
)

func Test_CreateConfigWithPanic(t *testing.T) {
	defer os.Remove(server.Debug)
	defer func() {
		if err := recover(); err != nil {
			// Panic had been fired.
		}
	}()

	createConfig()
	t.Error(expectedFormat.Panic)
}

func Test_LoadConfigWithPanic(t *testing.T) {
	defer os.Remove(server.Debug)
	defer func() {
		if err := recover(); err != nil {
			// Panic had been fired.
		}
	}()

	loadConfig()
	t.Error(expectedFormat.Panic)
}

func Test_CreateConfig(t *testing.T) {
	defer os.Remove(server.Debug)
	server.Initialize(true)
	createConfig()

	if server.Cfg.GetExtension(oauthKey.Config) == nil {
		t.Error(expectedFormat.NotNil)
	}
}

func Test_LoadConfig(t *testing.T) {
	defer os.Remove(server.Debug)
	server.Initialize(true)
	config := loadConfig()

	if config.AllowRefreshToken != true {
		t.Errorf(expectedFormat.BoolButFoundBool, true, config.AllowRefreshToken)
	}
	if config.AccessTokenDuration != 259200*time.Second {
		t.Errorf(expectedFormat.NumberButFoundNumber, 259200*time.Second, config.AccessTokenDuration)
	}
	if config.RefreshTokenDuration != 7776000*time.Second {
		t.Errorf(expectedFormat.NumberButFoundNumber, 7776000*time.Second, config.RefreshTokenDuration)
	}
	if config.AuthorizationCodeDuration != 300*time.Second {
		t.Errorf(expectedFormat.NumberButFoundNumber, 300*time.Second, config.AuthorizationCodeDuration)
	}

	// Validate grant types
	grantTypes := []string{AuthorizationCodeGrant, ClientCredentialsGrant, PasswordGrant, RefreshTokenGrant}
	if !reflect.DeepEqual(grantTypes, config.GrantTypes) {
		t.Errorf(expectedFormat.StringButFoundString, grantTypes, config.GrantTypes)
	}
	if grantsValidation == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if !grantsValidation.MatchString(AuthorizationCodeGrant) {
			t.Errorf(expectedFormat.BoolButFoundBool, true, grantsValidation.MatchString(AuthorizationCodeGrant))
		}
	}
}
