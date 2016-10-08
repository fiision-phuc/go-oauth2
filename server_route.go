package oauth2

import (
	"reflect"
	"regexp"

	"github.com/phuc0302/go-oauth2/inject"
)

// Route describes a default route component implementation.
type Route struct {
	path     string
	regex    *regexp.Regexp
	handlers map[string]interface{}
}

// BindHandler binds handler with specific http method.
func (r *Route) BindHandler(method string, handler interface{}) {
	/* Condition validation: only accept function */
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		panic("Request handler must be a function type.")
	}

	// Create handlers if neccessary
	if r.handlers == nil {
		r.handlers = map[string]interface{}{}
	}

	// Bind handler
	r.handlers[method] = handler
}

// InvokeHandler invokes handler.
func (r *Route) InvokeHandler(c *Request, s *Security) {
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

// URLPattern returns registered url pattern.
func (r *Route) URLPattern() string {
	return r.path
}

// MatchURLPattern matchs url pattern.
func (r *Route) MatchURLPattern(method string, urlPath string) (bool, map[string]string) {
	if matches := r.regex.FindStringSubmatch(urlPath); len(matches) > 0 && matches[0] == urlPath {
		if handler := r.handlers[method]; handler != nil {
			var params map[string]string

			if names := r.regex.SubexpNames(); len(names) > 1 {
				params = map[string]string{}

				for i, name := range names {
					if len(name) > 0 {
						params[name] = matches[i]
					}
				}
			}
			return true, params
		}
	}
	return false, nil
}
