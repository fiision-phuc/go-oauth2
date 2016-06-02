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
	createConfig(debug)
	defer os.Remove(debug)

	if !utils.FileExisted(debug) {
		t.Errorf("Expected %s file had been created but found nil.", debug)
	}
}

func Test_LoadConfig(t *testing.T) {
	//	defer os.Remove(debug)
	config := loadConfig(debug)

	if config == nil {
		t.Errorf("%s could not be loaded.", debug)
	}

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

	staticFolders := map[string]string{
		"/assets":    "assets",
		"/resources": "resources",
	}
	if !reflect.DeepEqual(staticFolders, config.StaticFolders) {
		t.Errorf("Expected '%s' but found '%s'.", staticFolders, config.StaticFolders)
	}
}
