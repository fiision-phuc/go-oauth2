package oauth2

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/phuc0302/go-oauth2/oauth_key"
	"github.com/phuc0302/go-server"
)

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
	AllowRefreshToken bool `json:"allow_refresh_token"`

	GrantTypes                []string      `json:"grant_types"`
	AccessTokenDuration       time.Duration `json:"access_token_duration"`       // In seconds
	RefreshTokenDuration      time.Duration `json:"refresh_token_duration"`      // In seconds
	AuthorizationCodeDuration time.Duration `json:"authorization_code_duration"` // In seconds
}

// createConfig generates a default oauth2 configuration.
//
// @return
// - config {Config} (an instance of oauth2's configuration)
func createConfig() (config *Config) {
	if server.Cfg == nil {
		panic("Server is not yet being initialized! Please run: 'server.Initialize'.")
	}

	// Create default config
	config = &Config{
		GrantTypes: []string{AuthorizationCodeGrant, ClientCredentialsGrant, PasswordGrant, RefreshTokenGrant},

		AllowRefreshToken:         true,
		AuthorizationCodeDuration: 300,
		AccessTokenDuration:       259200,
		RefreshTokenDuration:      7776000,
	}

	server.Cfg.SetExtension(oauthKey.Config, *config)
	server.Cfg.Save()
	return
}

// loadConfig retrieves previous configuration from file.
//
// @return
// - config {Config} (an instance of oauth2's configuration)
func loadConfig() (config *Config) {
	if server.Cfg == nil {
		panic("Server is not yet being initialized! Please run: 'server.Initialize'.")
	}

	// Generate config file if neccessary
	if info := server.Cfg.GetExtension(oauthKey.Config); info != nil {
		if configJSON, err := json.Marshal(info); err == nil {

			config = new(Config)
			if err = json.Unmarshal(configJSON, config); err == nil {
				// Everything is good to go.
			} else {
				config = createConfig()
			}
		} else {
			config = createConfig()
		}
	} else {
		config = createConfig()
	}

	grantsValidation = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(config.GrantTypes, "|")))
	config.AuthorizationCodeDuration *= time.Second
	config.RefreshTokenDuration *= time.Second
	config.AccessTokenDuration *= time.Second
	return
}
