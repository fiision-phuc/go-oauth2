package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/phuc0302/go-oauth2/utils"
)

func Test_CreateConfig(t *testing.T) {
	CreateConfigs(Debug)
	defer os.Remove(Debug)

	if !utils.FileExisted(Debug) {
		t.Errorf("Expected %s file had been created but found nil.", Debug)
	}
}

func Test_LoadConfig(t *testing.T) {
	defer os.Remove(Debug)
	defer os.Remove(Release)
	config := LoadConfigs(Debug)

	if config == nil {
		t.Errorf("%s could not be loaded.", Debug)
	}

	allowMethods := []string{COPY, DELETE, GET, HEAD, LINK, OPTIONS, PATCH, POST, PURGE, PUT, UNLINK}
	if !reflect.DeepEqual(allowMethods, config.AllowMethods) {
		t.Errorf("Expected '%s' but found '%s'.", allowMethods, config.AllowMethods)
	}

	staticFolders := map[string]string{
		"/assets":    "assets",
		"/resources": "resources",
	}
	if !reflect.DeepEqual(staticFolders, config.StaticFolders) {
		t.Errorf("Expected '%s' but found '%s'.", staticFolders, config.StaticFolders)
	}

	if config.ReadTimeout != 15 {
		t.Errorf("Expected read timeout is 15 seconds but found %d seconds.", config.ReadTimeout)
	}

	if config.WriteTimeout != 15 {
		t.Errorf("Expected write timeout is 15 seconds but found %d seconds.", config.WriteTimeout)
	}
}
