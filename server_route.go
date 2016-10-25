package oauth2

import "regexp"

// ServerRoute describes a route component.
type ServerRoute struct {
	path  string
	regex *regexp.Regexp

	handlers map[string]ContextHandler
}

// bindHandler binds handler with specific http method.
func (r *ServerRoute) bindHandler(method string, handler ContextHandler) {
	/* Condition validation: only accept function */
	if handler == nil {
		panic("Request handler must not be nil.")
	}

	// Create handlers if neccessary
	if r.handlers == nil {
		r.handlers = make(map[string]ContextHandler)
	}

	// Bind handler
	r.handlers[method] = handler
}

// invokeHandler invokes handler.
func (r *ServerRoute) invokeHandler(c *RequestContext, s *OAuthContext) {
	handler := r.handlers[c.request.Method]
	handler(c, s)
}

// match matchs path.
func (r *ServerRoute) match(method string, path string) (bool, map[string]string) {
	if matches := r.regex.FindStringSubmatch(path); len(matches) > 0 && matches[0] == path {
		if handler := r.handlers[method]; handler != nil {
			// Find path params if there is any
			var params map[string]string
			if names := r.regex.SubexpNames(); len(names) > 1 {

				params = make(map[string]string)
				for i, name := range names {
					if len(name) > 0 {
						params[name] = matches[i]
					}
				}
			}

			// Return result
			return true, params
		}
	}
	return false, nil
}
