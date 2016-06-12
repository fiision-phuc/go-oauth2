package oauth2

import (
	"crypto/rsa"
	"regexp"
)

// Define HTTP Methods.
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

// Define configuration file's name.
const (
	debug   = "oauth2.debug.cfg"
	release = "oauth2.release.cfg"
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
	TableRefreshToken = "oauth_refresh_token"
	TableAccessToken  = "oauth_access_token"
	TableClient       = "oauth_client"
	TableUser         = "oauth_user"
)

// Define global variables.
var (
	// Global config.
	cfg *config

	// Global objects.
	tokenStore    IStore
	objectFactory IFactory

	//Global jwt
	privateKey *rsa.PrivateKey

	// Global validation
	redirectPaths map[int]string
	//	clientValidation  *regexp.Regexp
	grantsValidation  *regexp.Regexp
	methodsValidation *regexp.Regexp

	// Define finder
	bearerRegex    = regexp.MustCompile("^(B|b)earer\\s.+$")
	globsRegex     = regexp.MustCompile(`\*\*`)
	pathParamRegex = regexp.MustCompile(`{[^/#?()\.\\]+}`)
)
