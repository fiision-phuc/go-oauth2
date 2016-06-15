package oauth2

// IRoute describes a route component's characteristic.
type IRoute interface {

	// Bind handler with specific http method.
	BindHandler(method string, handler interface{})

	// Invoke handler.
	InvokeHandler(c *Request, s *Security)

	// Return registered url pattern.
	URLPattern() string

	// Match url pattern.
	MatchURLPattern(method string, urlPattern string) (bool, map[string]string)
}
