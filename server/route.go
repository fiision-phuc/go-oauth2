package server

import "regexp"

// Route describes a route component.
type Route struct {
	regex    *regexp.Regexp
	handlers map[string]HandleContextFunc
}

// DefaultRoute creates new route component.
func DefaultRoute(regexPattern string) *Route {
	route := &Route{
		regex:    regexp.MustCompile(regexPattern),
		handlers: make(map[string]HandleContextFunc),
	}
	return route
}

// BindHandler binds handler with specific http method.
func (r *Route) BindHandler(method string, handler HandleContextFunc) {
	/* Condition validation: only accept function */
	if handler == nil {
		panic("Request handlers must not be nil.")
	}
	r.handlers[method] = handler
}

// InvokeHandlers invokes handlers.
func (r *Route) InvokeHandlers(c *RequestContext) {
	handler := r.handlers[c.Method]
	handler(c)
}

// Match matchs path.
func (r *Route) Match(method string, path string) (bool, map[string]string) {
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
