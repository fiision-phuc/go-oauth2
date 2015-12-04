package oauth2

import (
	"os"
	"reflect"
	"testing"

	"github.com/phuc0302/go-oauth2/utils"
)

func Test_CreateConfig(t *testing.T) {
	config := CreateConfigs()
	defer os.Remove("oauth2.cnf")

	if config == nil {
		t.Error("Expected config is not nil but found nil.")
	}

	if !utils.FileExisted("oauth2.cnf") {
		t.Error("Expected oauth2.cnf file had been created but found nil.")
	}
}

func Test_LoadConfig(t *testing.T) {
	defer os.Remove("oauth2.cnf")
	config := LoadConfigs()
	if config != nil {
		t.Errorf("Expected return nil when oauth2.cnf is not available.")
	}

	CreateConfigs()
	config = LoadConfigs()

	allowMethods := []string{COPY, DELETE, GET, HEAD, LINK, OPTIONS, PATCH, POST, PURGE, PUT, UNLINK}
	if !reflect.DeepEqual(allowMethods, config.AllowMethods) {
		t.Errorf("Expected %s but found %s", allowMethods, config.AllowMethods)
	}

	staticFolders := map[string]string{"/resources": "resources"}
	if !reflect.DeepEqual(staticFolders, config.StaticFolders) {
		t.Errorf("Expected %s but found %s", staticFolders, config.StaticFolders)
	}

	if config.TimeoutRead != 15 {
		t.Errorf("Expected read timeout is 15 seconds but found %d", config.TimeoutRead)
	}

	if config.TimeoutWrite != 15 {
		t.Errorf("Expected write timeout is 15 seconds but found %d", config.TimeoutRead)
	}
}

func Test_SetEnv(t *testing.T) {
	tests := []struct {
		key   string
		value string
	}{
		{"", "development"},
		{"development", ""},
		{"not_development", "not_development"},
	}

	for _, test := range tests {
		SetEnv(test.key, test.value)
		envValue := os.Getenv(test.key)

		if len(test.key) == 0 && len(envValue) != 0 {
			t.Errorf("Expect ignored environment value is: %s, but found: %s", test.value, envValue)
		} else if len(test.key) != 0 && test.value != envValue {
			t.Errorf("Expect environment value is: %s, but found: %s", test.value, envValue)
		}
	}
}

func Test_GetEnv(t *testing.T) {
	tests := []struct {
		key   string
		value string
	}{
		{"", "development"},
		{"development", ""},
		{"not_development", "not_development"},
	}

	for _, test := range tests {
		SetEnv(test.key, test.value)
		envValue := GetEnv(test.key)

		if len(test.key) == 0 && len(envValue) != 0 {
			t.Errorf("Expect ignored environment value is: %s, but found: %s", test.value, envValue)
		} else if len(test.key) != 0 && test.value != envValue {
			t.Errorf("Expect environment value is: %s, but found: %s", test.value, envValue)
		}
	}
}
