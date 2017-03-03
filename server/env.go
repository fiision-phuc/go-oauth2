package server

import "os"

// GetEnv retrieves value from environment.
func GetEnv(key string) string {
	if len(key) == 0 {
		return ""
	}
	return os.Getenv(key)
}

// SetEnv persists key-value to environment.
func SetEnv(key string, value string) {
	if len(key) == 0 || len(value) == 0 {
		return
	}
	os.Setenv(key, value)
}
