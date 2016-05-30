package oauth2

import (
	"fmt"
	"reflect"
	"regexp"
)

// DefaultRoute describes default implementation for route.
type DefaultRoute struct {
	urlPattern string
	regex      *regexp.Regexp
	handlers   map[string]interface{}
}

// CreateDefaultRoute returns a default route object.
func CreateDefaultRoute(pattern string) IRoute {
	regexPattern := pathParamRegex.ReplaceAllStringFunc(pattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	})
	regexPattern += "/?"

	route := DefaultRoute{pattern, regexp.MustCompile(regexPattern), make(map[string]interface{})}
	return &route
}

// MARK: Route's members
func (r *DefaultRoute) BindHandler(method string, handler interface{}) {
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		panic("Request handler must be a function type.")
	}
	r.handlers[method] = handler
}
func (r *DefaultRoute) InvokeHandler(c *Request, s *Security) {
	invoker := CreateInvoker()
	handler := r.handlers[c.Method()]

	// Call handler
	invoker.Map(c)
	invoker.Map(s)
	_, err := invoker.Invoke(handler)

	// Condition validation: Validate error
	if err != nil {
		panic(err)
	}
}

func (r *DefaultRoute) URLPattern() string {
	return r.urlPattern
}
func (r *DefaultRoute) MatchURLPattern(method string, urlPath string) (bool, map[string]string) {
	// Condition validation: Match request url
	matches := r.regex.FindStringSubmatch(urlPath)
	if len(matches) == 0 || matches[0] != urlPath {
		return false, nil
	}

	// Condition validation: Match request method
	handler := r.handlers[method]
	if handler == nil {
		return false, nil
	}

	// Extract path params
	params := make(map[string]string)
	for i, name := range r.regex.SubexpNames() {
		if len(name) > 0 {
			params[name] = matches[i]
		}
	}
	return true, params
}
