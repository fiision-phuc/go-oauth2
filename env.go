package oauth2

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/phuc0302/go-oauth2/utils"
)

var Development = false

/** Persist all configuration to file. */
func CreateConfigs() bool {
	// Remove old file is existed
	if utils.FileExisted("oauth2.cnf") {
		os.Remove("oauth2.cnf")
	}

	host := GetEnv(ENV_HOST)
	port := GetEnv(ENV_PORT)
	if len(host) == 0 {
		host = "localhost"
	}
	if len(port) == 0 {
		port = "8080"
	}

	// Create default config
	config := map[string]interface{}{
		ENV_HOST: host,
		ENV_PORT: port,

		ENV_HEADERS_SIZE:  5,
		ENV_TIMEOUT_READ:  15,
		ENV_TIMEOUT_WRITE: 15,

		ENV_ALLOW_METHODS:  []string{COPY, DELETE, GET, HEAD, LINK, OPTIONS, PATCH, POST, PURGE, PUT, UNLINK},
		ENV_STATIC_FOLDERS: map[string]string{"/resources": "resources"},
	}

	// Create new file
	configJson, _ := json.MarshalIndent(config, "", "  ")
	file, _ := os.Create("oauth2.cnf")
	file.Write(configJson)
	file.Close()

	return true
}

/** Retrieve preset configuration file. */
func LoadConfigs() *Config {
	if !utils.FileExisted("oauth2.cnf") {
		return nil
	}

	file, err := os.Open("oauth2.cnf")
	if err != nil {
		return nil
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}

	config := Config{}
	json.Unmarshal(bytes, &config)

	folders := make(map[string]string)
	for path, folder := range config.StaticFolders {
		folders[utils.FormatPath(path)] = folder
	}
	config.StaticFolders = folders
	return &config
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
