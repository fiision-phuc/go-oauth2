package oauth2

import (
	"crypto/rsa"
	"regexp"
)

// Configuration file's name.
const (
	debug   = "oauth2.debug.cfg"
	release = "oauth2.release.cfg"
)

// HTTP Methods.
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

// OAuth2.0 flows.
const (
	// For apps running on a web server
	AuthorizationCodeGrant = "authorization_code"

	// For application access
	ClientCredentialsGrant = "client_credentials"

	// For browser-based or mobile apps
	ImplicitGrant = "implicit"

	// For logging in with a username and password
	PasswordGrant = "password"

	// Should allow refresh token or not
	RefreshTokenGrant = "refresh_token"
)

// OAuth2 tables.
const (
	TableAccessToken  = "oauth_access_token"
	TableClient       = "oauth_client"
	TableRefreshToken = "oauth_refresh_token"
	TableUser         = "oauth_user"
)

// Global variables.
var (
	// Global public config's instance.
	Cfg Config

	// Global public token store's instance.
	TokenStore IStore

	// Global internal object factory's instance.
	objectFactory IFactory

	// Global internal private key.
	privateKey *rsa.PrivateKey

	// Global internal redirect map.
	redirectPaths map[int]string
)

// Global regex.
var (
	// OAuth2.0 regex
	bearerFinder      = regexp.MustCompile("^(B|b)earer\\s.+$")
	grantsValidation  *regexp.Regexp
	methodsValidation *regexp.Regexp

	// Globs & Path regex
	globsFinder = regexp.MustCompile(`\*\*`)
	pathFinder  = regexp.MustCompile(`{[^/#?()\.\\]+}`)
)

// Type alias
type ContextHandler func(request *Request, security *Security)
type GroupHandler func(server *Server)
