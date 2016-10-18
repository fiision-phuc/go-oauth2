package oauth2

import (
	"fmt"
	"regexp"

	"github.com/phuc0302/go-oauth2/inject"
)

// route describes a route component implementation.
type route struct {
	path  string
	regex *regexp.Regexp

	handlers map[string]ContextHandler
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// createRoute creates new route component.
func createRoute(path string) *route {
	regexPattern := pathFinder.ReplaceAllStringFunc(path, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:len(m)-1])
	})
	regexPattern = globsFinder.ReplaceAllStringFunc(regexPattern, func(m string) string {
		return fmt.Sprintf(`(?P<_%d>[^#?]*)`, 0)
	})

	if len(regexPattern) == 1 && regexPattern == "/" {
		regexPattern = fmt.Sprintf("^%s?$", regexPattern)
	} else {
		regexPattern = fmt.Sprintf("^%s/?$", regexPattern)
	}

	route := &route{
		path:     path,
		handlers: map[string]ContextHandler{},
		regex:    regexp.MustCompile(regexPattern),
	}
	return route
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// bindHandler binds handler with specific http method.
func (r *route) bindHandler(method string, handler ContextHandler) {
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
func (r *route) invokeHandler(c *Request, s *Security) {
	invoker := inject.CreateInvoker()
	handler := r.handlers[c.request.Method]

	// Call handler
	invoker.Map(c)
	invoker.Map(s)

	/* Condition validation: Validate error */
	if _, err := invoker.Invoke(handler); err != nil {
		panic(err)
	}
}

// match matchs path.
func (r *route) match(method string, path string) (bool, map[string]string) {
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
