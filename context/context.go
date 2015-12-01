package context

import (
	"net/http"
	"net/url"
)

// Context object.
type Context struct {
	Protocol string
	Method   string

	URLPath     string
	Queries     url.Values
	PathQueries map[string]string

	request  *http.Request
	response http.ResponseWriter
}

// CreateContext return a default context.
func CreateContext(request *http.Request, response http.ResponseWriter) *Context {
	context := &Context{
		Method:   request.Method,
		URLPath:  request.URL.Path,
		Protocol: request.Proto,

		request:  request,
		response: response,
	}

	context.Headers = request.Header
	return context
}
