package oauth2

import "github.com/phuc0302/go-oauth2/context"

// IRoute descripts a route component's characteristic.
type IRoute interface {

	// Bind handler with specific http method.
	BindHandler(method string, handler interface{})
	// Invoke handler.
	InvokeHandler(c *context.Request, s *context.Security)

	// Return registered url pattern.
	URLPattern() string
	// Match url pattern.
	MatchURLPattern(method string, urlPattern string) (bool, map[string]string)
}
