package oauth2

import (
	"regexp"
	"time"
)

// Define keywords for http methods.
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

// Define keywords for oauth2.0 flows.
const (
	AuthorizationCodeGrant = "authorization_code" // For apps running on a web server
	ClientCredentialsGrant = "client_credentials" // For application access
	PasswordGrant          = "password"           // For logging in with a username and password
	RefreshTokenGrant      = "refresh_token"      // Should allow refresh token or not

	//	ImplicitGrant          = "implicit"           // For browser-based or mobile apps
)

// Config struct descripts a configuration  object  that  will  be  used  during
// application life time.
type Config struct {
	Development bool `json:"development,omitempty"`

	Host              string            `json:"host,omitempty"`
	Port              string            `json:"port,omitempty"`
	TLSPort           string            `json:"tls_port,omitempty"`
	HeaderSize        int               `json:"headers_size,omitempty"`  // In KB
	TimeoutRead       time.Duration     `json:"timeout_read,omitempty"`  // In seconds
	TimeoutWrite      time.Duration     `json:"timeout_write,omitempty"` // In seconds
	AllowMethods      []string          `json:"allow_methods,omitempty"`
	StaticFolders     map[string]string `json:"static_folders,omitempty"`
	methodsValidation *regexp.Regexp    `json:"-"`

	Grant                     []string       `json:"grant_types,omitempty"`
	DurationAccessToken       time.Duration  `json:"duration_access_token,omitempty"`       // In seconds
	DurationRefreshToken      time.Duration  `json:"duration_refresh_token,omitempty"`      // In seconds
	DurationAuthorizationCode time.Duration  `json:"duration_authorization_code,omitempty"` // In seconds
	allowRefreshToken         bool           `json:"-"`
	clientValidation          *regexp.Regexp `json:"-"`
	grantsValidation          *regexp.Regexp `json:"-"`
}
