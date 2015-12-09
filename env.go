package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/phuc0302/go-oauth2/utils"
)

// ConfigFile defines configuration file's name.
const ConfigFile = "oauth2.cfg"

// CreateConfigs generates a default configuration file.
func CreateConfigs() {
	if utils.FileExisted(ConfigFile) {
		os.Remove(ConfigFile)
	}

	host := GetEnv("HOST")
	port := GetEnv("PORT")
	if len(host) == 0 {
		host = "localhost"
	}
	if len(port) == 0 {
		port = "8080"
	}

	// Create default config
	config := Config{
		Development: true,

		Host:          host,
		Port:          port,
		TLSPort:       "8443",
		HeaderSize:    5,
		TimeoutRead:   10,
		TimeoutWrite:  10,
		AllowMethods:  []string{COPY, DELETE, GET, HEAD, LINK, OPTIONS, PATCH, POST, PURGE, PUT, UNLINK},
		StaticFolders: map[string]string{"/resources": "resources"},

		Grant:                     []string{AuthorizationCodeGrant, ClientCredentialsGrant, PasswordGrant, RefreshTokenGrant},
		DurationAccessToken:       3600,
		DurationRefreshToken:      1209600,
		DurationAuthorizationCode: 30,
	}

	// Create new file
	configJSON, _ := json.MarshalIndent(config, "", "  ")
	file, _ := os.Create(ConfigFile)
	file.Write(configJSON)
	file.Close()
}

// LoadConfigs retrieve preset configuration file.
func LoadConfigs() *Config {
	if !utils.FileExisted(ConfigFile) {
		CreateConfigs()
	}

	file, err := os.Open(ConfigFile)
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

	// Convert duration to seconds
	config.DurationAccessToken = config.DurationAccessToken * time.Second
	config.DurationRefreshToken = config.DurationRefreshToken * time.Second
	config.DurationAuthorizationCode = config.DurationAuthorizationCode * time.Second

	// Define regular expressions
	//	regexp.MustCompile(`:[^/#?()\.\\]+`)
	config.grantsValidation = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(config.Grant, "|")))
	config.methodsValidation = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(config.AllowMethods, "|")))

	for _, grant := range config.Grant {
		if grant == RefreshTokenGrant {
			config.allowRefreshToken = true
			break
		}
	}
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
