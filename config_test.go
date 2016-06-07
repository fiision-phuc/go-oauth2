package oauth2

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/phuc0302/go-oauth2/test"
	"github.com/phuc0302/go-oauth2/utils"
)

func Test_CreateConfig(t *testing.T) {
	//	logrus.SetFormatter(&logrus.JSONFormatter{})
	//	logrus.SetFormatter(&logrus.TextFormatter{})
	//	logrus.SetOutput(os.Stderr)

	//	logrus.SetLevel(logrus.DebugLevel)
	//	logrus.AddHook(&slackrus.SlackrusHook{
	//		HookURL:        "https://hooks.slack.com/services/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",
	//		AcceptedLevels: slackrus.LevelThreshold(logrus.DebugLevel),
	//		Channel:        "#keywords",
	//		Username:       "Server",
	//		IconEmoji:      ":ghost:",
	//	})
	//	logrus.Warn("warn")
	//	logrus.Info("info")
	//	logrus.Debug("debug")

	createConfig(debug)
	defer os.Remove(debug)

	if !utils.FileExisted(debug) {
		t.Errorf("Expected %s file had been created but found nil.", debug)
	}
}

func Test_LoadConfig(t *testing.T) {
	defer os.Remove(debug)
	config := loadConfig(debug)

	// Validate configFile
	if config == nil {
		t.Errorf("%s could not be loaded.", debug)
	}

	// Validate basic information
	if config.Host != "localhost" {
		t.Errorf(test.ExpectedStringButFoundString, "localhost", config.Host)
	}
	if config.Port != 8080 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 8080, config.Port)
	}
	if config.TLSPort != 8443 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 8443, config.TLSPort)
	}
	if config.HeaderSize != 5120 {
		t.Errorf(test.ExpectedNumberButFoundNumber, 5120, config.HeaderSize)
	}
	if config.ReadTimeout != 15*time.Second {
		t.Errorf(test.ExpectedNumberButFoundNumber, 15*time.Second, config.ReadTimeout)
	}
	if config.WriteTimeout != 15*time.Second {
		t.Errorf(test.ExpectedNumberButFoundNumber, 15*time.Second, config.WriteTimeout)
	}

	if config.AllowRefreshToken != true {
		t.Errorf(test.ExpectedBoolButFoundBool, true, config.AllowRefreshToken)
	}
	if config.AccessTokenDuration != 259200*time.Second {
		t.Errorf(test.ExpectedNumberButFoundNumber, 259200*time.Second, config.AccessTokenDuration)
	}
	if config.RefreshTokenDuration != 7776000*time.Second {
		t.Errorf(test.ExpectedNumberButFoundNumber, 7776000*time.Second, config.RefreshTokenDuration)
	}
	if config.AuthorizationCodeDuration != 300*time.Second {
		t.Errorf(test.ExpectedNumberButFoundNumber, 300*time.Second, config.AuthorizationCodeDuration)
	}

	// Validate allow methods
	allowMethods := []string{COPY, DELETE, GET, HEAD, LINK, OPTIONS, PATCH, POST, PURGE, PUT, UNLINK}
	if !reflect.DeepEqual(allowMethods, config.AllowMethods) {
		t.Errorf("Expected '%s' but found '%s'.", allowMethods, config.AllowMethods)
	}
	if methodsValidation == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if !methodsValidation.MatchString(COPY) {
			t.Errorf(test.ExpectedBoolButFoundBool, true, methodsValidation.MatchString(COPY))
		}
	}

	// Validate redirect paths
	if redirectPaths == nil || len(redirectPaths) != 1 {
		t.Error(test.ExpectedNotNil)
	}
	if redirectPaths[401] != "/login" {
		t.Errorf(test.ExpectedStringButFoundString, "/login", redirectPaths[401])
	}

	// Validate static folders
	staticFolders := map[string]string{
		"/assets":    "assets",
		"/resources": "resources",
	}
	if !reflect.DeepEqual(staticFolders, config.StaticFolders) {
		t.Errorf(test.ExpectedStringButFoundString, staticFolders, config.StaticFolders)
	}

	// Validate grant types
	grantTypes := []string{AuthorizationCodeGrant, ClientCredentialsGrant, PasswordGrant, RefreshTokenGrant}
	if !reflect.DeepEqual(grantTypes, config.GrantTypes) {
		t.Errorf(test.ExpectedStringButFoundString, grantTypes, config.GrantTypes)
	}
	if grantsValidation == nil {
		t.Error(test.ExpectedNotNil)
	} else {
		if !grantsValidation.MatchString(AuthorizationCodeGrant) {
			t.Errorf(test.ExpectedBoolButFoundBool, true, grantsValidation.MatchString(AuthorizationCodeGrant))
		}
	}
}