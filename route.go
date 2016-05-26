package oauth2

// Route descripts a characteristic of routing components.
type Route interface {
	/** Add handler to route. */
	AddHandler(method string, handler interface{})
	/** Invoke handler. */
	InvokeHandler(c *Request, s *SecurityContext)

	/** Look up matched url path. */
	GetPattern() string
	Match(method string, urlPath string) (bool, map[string]string)
}
