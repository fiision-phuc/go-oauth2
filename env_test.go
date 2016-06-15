package oauth2

import (
	"os"
	"testing"
)

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
