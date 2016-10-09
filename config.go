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

	"github.com/phuc0302/go-oauth2/util"
)

// Define configuration file's name.
const (
	debug   = "oauth2.debug.cfg"
	release = "oauth2.release.cfg"
)

// Define HTTP Methods.
const (
	Copy    = "copy"
	Delete  = "delete"
	Get     = "get"
	Head    = "head"
	Link    = "link"
	Options = "options"
	Patch   = "patch"
	Post    = "post"
	Purge   = "purge"
	Put     = "put"
	Unlink  = "unlink"
)

// Define OAuth2 flows.
const (
	AuthorizationCodeGrant = "authorization_code" // For apps running on a web server
	ClientCredentialsGrant = "client_credentials" // For application access
	ImplicitGrant          = "implicit"           // For browser-based or mobile apps
	PasswordGrant          = "password"           // For logging in with a username and password
	RefreshTokenGrant      = "refresh_token"      // Should allow refresh token or not

)

// Define OAuth2 tables.
const (
	TableAccessToken  = "oauth_access_token"
	TableClient       = "oauth_client"
	TableRefreshToken = "oauth_refresh_token"
	TableUser         = "oauth_user"
)

// Config describes a configuration object that will be used during application life time.
type Config struct {

	// Server
	Host    string `json:"host"`
	Port    int    `json:"port"`
	TLSPort int    `json:"tls_port"`

	// Header
	HeaderSize    int           `json:"header_size"`    // In KB
	MultipartSize int64         `json:"multipart_size"` // In MB
	ReadTimeout   time.Duration `json:"timeout_read"`   // In seconds
	WriteTimeout  time.Duration `json:"timeout_write"`  // In seconds

	// HTTP Method
	AllowMethods  []string          `json:"allow_methods"`
	RedirectPaths map[string]int    `json:"redirect_paths"`
	StaticFolders map[string]string `json:"static_folders"`

	// Log
	LogLevel     string `json:"log_level"`
	SlackURL     string `json:"slack_url"`
	SlackIcon    string `json:"slack_icon"`
	SlackUser    string `json:"slack_user"`
	SlackChannel string `json:"slack_channel"`

	// Jwt
	PrivateKey []byte `json:"private_key"`

	// OAuth2.0
	GrantTypes                []string      `json:"grant_types"`
	AllowRefreshToken         bool          `json:"allow_refresh_token"`
	AccessTokenDuration       time.Duration `json:"access_token_duration"`       // In seconds
	RefreshTokenDuration      time.Duration `json:"refresh_token_duration"`      // In seconds
	AuthorizationCodeDuration time.Duration `json:"authorization_code_duration"` // In seconds
}

// CreateConfig generates a default configuration file.
func CreateConfig(configFile string) {
	if util.FileExisted(configFile) {
		os.Remove(configFile)
	}

	// Create default config
	config := Config{
		Host:    "localhost",
		Port:    8080,
		TLSPort: 8443,

		HeaderSize:    5,
		MultipartSize: 1,
		ReadTimeout:   15,
		WriteTimeout:  15,

		AllowMethods: []string{Copy, Delete, Get, Head, Link, Options, Patch, Post, Purge, Put, Unlink},
		RedirectPaths: map[string]int{
			"/login": 401,
		},
		StaticFolders: map[string]string{
			"/assets":    "assets",
			"/resources": "resources",
		},

		LogLevel:     "debug",
		SlackURL:     "",
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

// LoadConfig retrieves previous configuration from file.
func LoadConfig(configFile string) Config {
	// Generate config file if neccessary
	if !util.FileExisted(configFile) {
		CreateConfig(configFile)
	}

	// Load config file
	config := Config{}
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
