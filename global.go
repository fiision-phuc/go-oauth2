package oauth2

import "regexp"

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
