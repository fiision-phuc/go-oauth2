package oauth2

//import (
//	"fmt"
//	"reflect"
//	"regexp"

//	"github.com/phuc0302/go-cocktail-di"
//)

//type DefaultRoute struct {
//	Pattern  string
//	regex    *regexp.Regexp
//	handlers map[string]interface{}
//}

//// MARK: Struct's constructors
//func createDefaultRoute(pattern string) *DefaultRoute {
//	regex := regexp.MustCompile(`:[^/#?()\.\\]+`)

//	// Convert param to regular expression
//	regexPattern := regex.ReplaceAllStringFunc(pattern, func(m string) string {
//		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
//	})
//	regexPattern += "/?"

//	route := DefaultRoute{pattern, regexp.MustCompile(regexPattern), make(map[string]interface{}, 4)}
//	return &route
//}

//// MARK: Route's members
//func (r *DefaultRoute) AddHandler(method string, handler interface{}) {
//	if reflect.TypeOf(handler).Kind() != reflect.Func {
//		panic("Request handler must be a function type.")
//	}
//	r.handlers[method] = handler
//}
//func (r *DefaultRoute) InvokeHandler(c *Context) {
//	injector := di.Injector()
//	handler := r.handlers[c.request.Method]

//	// Call handler
//	injector.Map(c)
//	_, err := injector.Invoke(handler)

//	// Condition validation: Validate error
//	if err != nil {
//		panic(err)
//	}
//}

//func (r *DefaultRoute) GetPattern() string {
//	return r.Pattern
//}
//func (r *DefaultRoute) Match(method string, urlPath string) (bool, map[string]string) {
//	// Condition validation: Match request url
//	matches := r.regex.FindStringSubmatch(urlPath)
//	if len(matches) == 0 || matches[0] != urlPath {
//		return false, nil
//	}

//	// Condition validation: Match request method
//	handler := r.handlers[method]
//	if handler == nil {
//		return false, nil
//	}

//	// Extract path params
//	params := make(map[string]string)
//	for i, name := range r.regex.SubexpNames() {
//		if len(name) > 0 {
//			params[name] = matches[i]
//		}
//	}
//	return true, params
//}
