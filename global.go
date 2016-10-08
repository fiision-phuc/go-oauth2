package oauth2

import (
	"crypto/rsa"
	"regexp"
)

// Define global variables.
var (
	Cfg        Config
	TokenStore IStore

	// Factory
	objectFactory IFactory
	// Global jwt
	privateKey *rsa.PrivateKey

	// Global validation
	redirectPaths     map[int]string
	grantsValidation  *regexp.Regexp
	methodsValidation *regexp.Regexp

	// Define finder
	bearerFinder = regexp.MustCompile("^(B|b)earer\\s.+$")
	globsFinder  = regexp.MustCompile(`\*\*`)
	pathFinder   = regexp.MustCompile(`{[^/#?()\.\\]+}`)
)
