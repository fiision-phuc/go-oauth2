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

// Configuration file's name.
const configFile = "oauth2.cfg"

// OAuth2.0 flows.
const (
	// For apps running on a web server.
	AuthorizationCodeGrant = "authorization_code"

	// For application access.
	ClientCredentialsGrant = "client_credentials"

	// For browser-based or mobile apps.
	ImplicitGrant = "implicit"

	// For logging in with a username and password.
	PasswordGrant = "password"

	// Should allow refresh token or not.
	RefreshTokenGrant = "refresh_token"
)

// Config describes a configuration object that will be used during application life time.
type Config struct {
	GrantTypes []string `json:"grant_types"`
	PrivateKey []byte   `json:"private_key"`

	AllowRefreshToken         bool          `json:"allow_refresh_token"`
	AccessTokenDuration       time.Duration `json:"access_token_duration"`       // In seconds
	RefreshTokenDuration      time.Duration `json:"refresh_token_duration"`      // In seconds
	AuthorizationCodeDuration time.Duration `json:"authorization_code_duration"` // In seconds
}

// CreateConfig generates a default configuration file.
func CreateConfig() {
	if util.FileExisted(configFile) {
		os.Remove(configFile)
	}

	// Create default config
	config := Config{
		GrantTypes: []string{AuthorizationCodeGrant, ClientCredentialsGrant, PasswordGrant, RefreshTokenGrant},

		AllowRefreshToken:         true,
		AuthorizationCodeDuration: 300,
		AccessTokenDuration:       259200,
		RefreshTokenDuration:      7776000,
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
func LoadConfig() Config {
	// Generate config file if neccessary
	if !util.FileExisted(configFile) {
		CreateConfig()
	}

	// Load config file
	var config Config
	file, _ := os.Open(configFile)
	bytes, _ := ioutil.ReadAll(file)

	if err := json.Unmarshal(bytes, &config); err == nil {
		config.AccessTokenDuration *= time.Second
		config.RefreshTokenDuration *= time.Second
		config.AuthorizationCodeDuration *= time.Second

		// Define jwt
		privateKey, _ = x509.ParsePKCS1PrivateKey(config.PrivateKey)

		// Define regular expressions
		grantsValidation = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(config.GrantTypes, "|")))
	}
	return config
}
