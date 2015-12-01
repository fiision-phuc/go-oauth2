package oauth2

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/phuc0302/go-oauth2/utils"
)

// CreateConfigs persist all configuration to file.
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
	config := Config{
		Development:  true,
		Host:         host,
		Port:         port,
		HeaderSize:   5,
		TimeoutRead:  15,
		TimeoutWrite: 15,

		AllowMethods:  []string{COPY, DELETE, GET, HEAD, LINK, OPTIONS, PATCH, POST, PURGE, PUT, UNLINK},
		StaticFolders: map[string]string{"/resources": "resources"},
	}

	// Create new file
	configJSON, _ := json.MarshalIndent(config, "", "  ")
	file, _ := os.Create("oauth2.cnf")
	file.Write(configJSON)
	file.Close()

	return true
}

// LoadConfigs retrieve preset configuration file.
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

// GetEnv retrieve value from environment.
func GetEnv(key string) string {
	if len(key) == 0 {
		return ""
	}
	return os.Getenv(key)
}

// SetEnv persist key-value to environment.
func SetEnv(key string, value string) {
	if len(key) == 0 || len(value) == 0 {
		return
	}
	os.Setenv(key, value)
}
