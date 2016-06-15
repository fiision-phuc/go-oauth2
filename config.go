package oauth2

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/phuc0302/go-oauth2/utils"
)

// config describes a configuration  object  that  will  be  used  during application life time.
type config struct {

	// Server
	Host    string `json:"host,omitempty"`
	Port    int    `json:"port,omitempty"`
	TLSPort int    `json:"tls_port,omitempty"`

	// Header
	HeaderSize    int           `json:"header_size,omitempty"`    // In KB
	MultipartSize int64         `json:"multipart_size,omitempty"` // In MB
	ReadTimeout   time.Duration `json:"timeout_read,omitempty"`   // In seconds
	WriteTimeout  time.Duration `json:"timeout_write,omitempty"`  // In seconds

	// HTTP Method
	AllowMethods  []string          `json:"allow_methods,omitempty"`
	RedirectPaths map[string]int    `json:"redirect_paths,omitempty"`
	StaticFolders map[string]string `json:"static_folders,omitempty"`

	// Log
	LogLevel     string `json:"log_level,omitempty"`
	SlackURL     string `json:"slack_url,omitempty"`
	SlackIcon    string `json:"slack_icon,omitempty"`
	SlackUser    string `json:"slack_user,omitempty"`
	SlackChannel string `json:"slack_channel,omitempty"`

	// Jwt
	PrivateKey []byte `json:"private_key,omitempty"`

	// OAuth2.0
	GrantTypes                []string      `json:"grant_types,omitempty"`
	AllowRefreshToken         bool          `json:"allow_refresh_token,omitempty"`
	AccessTokenDuration       time.Duration `json:"access_token_duration,omitempty"`       // In seconds
	RefreshTokenDuration      time.Duration `json:"refresh_token_duration,omitempty"`      // In seconds
	AuthorizationCodeDuration time.Duration `json:"authorization_code_duration,omitempty"` // In seconds
}

// createConfig generates a default configuration file.
func createConfig(configFile string) {
	if utils.FileExisted(configFile) {
		os.Remove(configFile)
	}

	// Create default config
	config := config{
		Host:    "localhost",
		Port:    8080,
		TLSPort: 8443,

		HeaderSize:    5,
		MultipartSize: 1,
		ReadTimeout:   15,
		WriteTimeout:  15,

		AllowMethods: []string{COPY, DELETE, GET, HEAD, LINK, OPTIONS, PATCH, POST, PURGE, PUT, UNLINK},
		RedirectPaths: map[string]int{
			"/login": 401,
		},
		StaticFolders: map[string]string{
			"/assets":    "assets",
			"/resources": "resources",
		},

		LogLevel:     "debug",
		SlackURL:     "https://hooks.slack.com/services/",
		SlackIcon:    ":ghost:",
		SlackUser:    "OAuth2",
		SlackChannel: "#channel",

		GrantTypes:                []string{AuthorizationCodeGrant, ClientCredentialsGrant, PasswordGrant, RefreshTokenGrant},
		AllowRefreshToken:         true,
		AccessTokenDuration:       259200,
		RefreshTokenDuration:      7776000,
		AuthorizationCodeDuration: 300,
	}

	// Generate jwt key
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	config.PrivateKey = privateKeyDer

	// Create new file
	configJSON, _ := json.MarshalIndent(config, "", "  ")
	file, _ := os.Create(configFile)
	file.Write(configJSON)
	file.Close()
}

// loadConfig retrieves previous configuration from file.
func loadConfig(configFile string) config {
	// Generate config file if neccessary
	if !utils.FileExisted(configFile) {
		createConfig(configFile)
	}

	// Load config file
	config := config{}
	file, _ := os.Open(configFile)
	bytes, _ := ioutil.ReadAll(file)

	if err := json.Unmarshal(bytes, &config); err == nil {
		// Convert duration to seconds
		config.HeaderSize <<= 10
		config.MultipartSize <<= 20
		config.ReadTimeout *= time.Second
		config.WriteTimeout *= time.Second
		config.AccessTokenDuration *= time.Second
		config.RefreshTokenDuration *= time.Second
		config.AuthorizationCodeDuration *= time.Second

		// Define redirectPaths
		redirectPaths = make(map[int]string, len(config.RedirectPaths))
		for path, status := range config.RedirectPaths {
			redirectPaths[status] = path
		}

		// Define jwt
		privateKey, _ = x509.ParsePKCS1PrivateKey(config.PrivateKey)

		// Define regular expressions
		//	regexp.MustCompile(`:[^/#?()\.\\]+`)
		grantsValidation = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(config.GrantTypes, "|")))
		methodsValidation = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(config.AllowMethods, "|")))
	}
	return config
}
