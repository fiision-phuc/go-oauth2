package oauth2

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/phuc0302/go-server/expected_format"
	"github.com/phuc0302/go-server/util"
)

func Test_CreateConfig(t *testing.T) {
	defer os.Remove(configFile)
	CreateConfig()

	if !util.FileExisted(configFile) {
		t.Errorf("Expected %s file had been created but found nil.", configFile)
	}
}

func Test_LoadConfig(t *testing.T) {
	defer os.Remove(configFile)
	config := LoadConfig()

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

	// Validate private key
	if privateKey == nil {
		t.Errorf(expectedFormat.NotNil)
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
