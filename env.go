package oauth2

import (
	"os"

	"github.com/phuc0302/go-oauth2/utils"
)

var Development = false

const (
	ENV_HOST = "HOST"
	ENV_PORT = "PORT"

	ENV_HEADERS_SIZE  = "headers_size" // In Kb
	ENV_TIMEOUT_READ  = "timeout_read" // In seconds
	ENV_TIMEOUT_WRITE = "timeout_wrte" // In seconds
)

/** Retrieve preset configuration file. */
func LoadConfigs() bool {
	if !utils.FileExisted("oauth2.cnf") {
		return false
	}

	return true
}

/** Persist all configuration to file. */
func SaveConfigs() bool {

	return true
}

/** Retrieve value from environment. */
func GetEnv(key string) string {
	if len(key) == 0 {
		return ""
	}
	return os.Getenv(key)
}

/** Persist key-value to environment. */
func SetEnv(key string, value string) {
	if len(key) == 0 || len(value) == 0 {
		return
	}
	os.Setenv(key, value)
}
