package oauth2

import "regexp"

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

// Define global variables.
var (
	// Global config.
	cfg *config

	// Global objects.
	tokenStore    IStore
	objectFactory IFactory

	// Global validation
	redirectPaths map[int]string
	//	clientValidation  *regexp.Regexp
	grantsValidation  *regexp.Regexp
	methodsValidation *regexp.Regexp

	// Define finder
	bearerRegex    = regexp.MustCompile("^(B|b)earer\\s\\w+$")
	globsRegex     = regexp.MustCompile(`\*\*`)
	pathParamRegex = regexp.MustCompile(`{[^/#?()\.\\]+}`)
)
