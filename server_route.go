package oauth2

type Route interface {
	/** Add handler to route. */
	AddHandler(method string, handler interface{})
	/** Invoke handler. */
	InvokeHandler(c *RequestContext)

	/** Look up matched url path. */
	GetPattern() string
	Match(method string, urlPath string) (bool, map[string]string)
}
