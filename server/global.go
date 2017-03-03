package server

import "regexp"

// Error messages.
const (
	InvalidParameter = "Invalid '%s' parameter."
)

// Global variables
var (
	// Global public config's instance.
	Cfg Config

	// Global internal redirect map.
	redirectPaths map[int]string

	// HTTP method regex
	methodsValidation *regexp.Regexp
)

// HandleGroupFunc defines type alias for group func handler.
type HandleGroupFunc func(*Server)

// HandleContextFunc defines type alias for request context func handler.
type HandleContextFunc func(*RequestContext)

// Adapter defines type alias HandleContextFunc func decorator.
type Adapter func(HandleContextFunc) HandleContextFunc

// Adapt generates decorator for HandleContextFunc func.
func Adapt(f HandleContextFunc, adapters ...Adapter) HandleContextFunc {
	for i := len(adapters) - 1; i >= 0; i-- {
		adapter := adapters[i]
		f = adapter(f)
	}
	return f
}
