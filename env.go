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

// Define keywords for HTTP methods.
const (
	COPY    = "COPY"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	LINK    = "LINK"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PURGE   = "PURGE"
	PUT     = "PUT"
	UNLINK  = "UNLINK"
)

// Define keywords for OAuth2.0 flows.
const (
	AuthorizationCodeGrant = "authorization_code" // For apps running on a web server
	ClientCredentialsGrant = "client_credentials" // For application access
	PasswordGrant          = "password"           // For logging in with a username and password
	RefreshTokenGrant      = "refresh_token"      // Should allow refresh token or not

	//	ImplicitGrant          = "implicit"           // For browser-based or mobile apps
)

// Define table names.
const (
	TableAccessToken  = "access_token"
	TableClient       = "client"
	TableRefreshToken = "refresh_token"
	TableUser         = "user"
)

// ConfigFile defines configuration file's name.
const ConfigFile = "oauth2.cfg"

////////////////////////////////////////////////////////////////////////////////////////////////////

// Config descripts a configuration  object  that  will  be  used  during application life time.
type Config struct {
	Development bool `json:"development,omitempty"`

	Host          string            `json:"host,omitempty"`
	Port          string            `json:"port,omitempty"`
	TLSPort       string            `json:"tls_port,omitempty"`
	HeaderSize    int               `json:"headers_size,omitempty"`  // In KB
	TimeoutRead   time.Duration     `json:"timeout_read,omitempty"`  // In seconds
	TimeoutWrite  time.Duration     `json:"timeout_write,omitempty"` // In seconds
	AllowMethods  []string          `json:"allow_methods,omitempty"`
	StaticFolders map[string]string `json:"static_folders,omitempty"`

	Grant                     []string      `json:"grant_types,omitempty"`
	DurationAccessToken       time.Duration `json:"duration_access_token,omitempty"`       // In seconds
	DurationRefreshToken      time.Duration `json:"duration_refresh_token,omitempty"`      // In seconds
	DurationAuthorizationCode time.Duration `json:"duration_authorization_code,omitempty"` // In seconds

	allowRefreshToken bool           `json:"-"`
	clientValidation  *regexp.Regexp `json:"-"`
	grantsValidation  *regexp.Regexp `json:"-"`
	methodsValidation *regexp.Regexp `json:"-"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////

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

		Host:         host,
		Port:         port,
		TLSPort:      "8443",
		HeaderSize:   5,
		TimeoutRead:  10,
		TimeoutWrite: 10,
		AllowMethods: []string{COPY, DELETE, GET, HEAD, LINK, OPTIONS, PATCH, POST, PURGE, PUT, UNLINK},

		StaticFolders: map[string]string{
			"/oauth2/resources": "github.com/phuc0302/go-oauth2/resources",
			"/oauth2/templates": "github.com/phuc0302/go-oauth2/templates",
		},

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

// LoadConfigs retrieves previous configuration from file.
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

////////////////////////////////////////////////////////////////////////////////////////////////////

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
