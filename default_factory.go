package oauth2

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/phuc0302/go-oauth2/utils"
)

var (
	globsRegex     = regexp.MustCompile(`\*\*`)
	pathParamRegex = regexp.MustCompile(`:[^/#?()\.\\]+`)
)

// DefaultFactory descripts a default factory object.
type DefaultFactory struct {
}

// CreateRequestContext creates new request context.
func (d *DefaultFactory) CreateRequestContext(request *http.Request, response http.ResponseWriter) *Request {
	context := &Request{
		URLPath:  request.URL.Path,
		request:  request,
		response: response,
	}

	// Format request headers
	context.Header = make(map[string]string, len(request.Header))
	for k, v := range request.Header {
		context.Header[strings.ToLower(k)] = strings.ToLower(v[0])
	}

	// Parse body context if neccessary
	switch context.Method() {

	case GET:
		params := request.URL.Query()
		if len(params) > 0 {
			context.Queries = params
		}
		break

	case POST, PATCH:
		contentType := request.Header.Get("content-type")

		if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			params := utils.ParseForm(request)
			if len(params) > 0 {
				context.Queries = params
			}
		} else if strings.Contains(contentType, "multipart/form-data") {
			params := utils.ParseMultipartForm(request)

			if len(params) > 0 {
				context.Queries = params
			}
		}
		break

	default:
		break
	}
	return context
}

// CreateSecurityContext creates new security context.
func (d *DefaultFactory) CreateSecurityContext() *Security {
	return nil
}

// CreateRoute creates new route component.
func (d *DefaultFactory) CreateRoute(urlPattern string) IRoute {
	regexPattern := pathParamRegex.ReplaceAllStringFunc(urlPattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	})
	regexPattern += "/?"

	route := DefaultRoute{
		urlPattern: urlPattern,
		handlers:   make(map[string]interface{}),
		regex:      regexp.MustCompile(regexPattern),
	}
	return &route
}

// CreateRouter creates new router component.
func (d *DefaultFactory) CreateRouter() IRouter {
	return nil
}

// CreateStore creates new store component.
func (d *DefaultFactory) CreateStore() IStore {
	return nil
}
