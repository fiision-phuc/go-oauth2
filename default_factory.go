package oauth2

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// DefaultFactory describes a default factory object.
type DefaultFactory struct {
}

// CreateRequestContext creates new request context.
func (d *DefaultFactory) CreateRequestContext(request *http.Request, response http.ResponseWriter) *Request {
	context := &Request{
		Path:     request.URL.Path,
		request:  request,
		response: response,
	}

	// Format request headers
	if len(request.Header) > 0 {
		context.Header = make(map[string]string)

		for k, v := range request.Header {
			context.Header[strings.ToLower(k)] = strings.ToLower(v[0])
		}
	}

	// Parse body context if neccessary
	var params url.Values
	switch context.request.Method {

	case GET:
		params = request.URL.Query()
		break

	case POST, PATCH:
		contentType := context.Header["content-type"]

		if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			err := request.ParseForm()
			if err == nil {
				params = request.Form
			}
		} else if strings.HasPrefix(contentType, "multipart/form-data; boundary") {
			err := request.ParseMultipartForm(cfg.MultipartSize)
			if err == nil {
				params = request.MultipartForm.Value
			}
		}
		break

	default:
		break
	}

	// Process params
	if len(params) > 0 {
		context.QueryParams = make(map[string]string)

		for k, v := range params {
			context.QueryParams[k] = v[0]
		}
	}
	return context
}

// CreateSecurityContext creates new security context.
func (d *DefaultFactory) CreateSecurityContext(requestContext *Request) *Security {
	headerToken := requestContext.Header["authorization"]
	isBearer := bearerRegex.MatchString(headerToken)

	/* Condition validation: Validate existing of authorization header */
	if !isBearer {
		return nil
	}

	headerToken = headerToken[7:]
	accessToken := tokenStore.FindAccessToken(headerToken)

	/* Condition validation: Validate expiration time */
	if accessToken == nil || accessToken.IsExpired() {
		return nil
	}

	client := tokenStore.FindClientWithID(accessToken.ClientID())
	user := tokenStore.FindUserWithID(accessToken.UserID())
	securityContext := &Security{
		AuthClient:      client,
		AuthUser:        user,
		AuthAccessToken: accessToken,
	}
	return securityContext
}

// CreateRoute creates new route component.
func (d *DefaultFactory) CreateRoute(urlPattern string) IRoute {
	regexPattern := pathParamRegex.ReplaceAllStringFunc(urlPattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:len(m)-1])
	})
	regexPattern += "/?"

	route := DefaultRoute{
		path:     urlPattern,
		handlers: map[string]interface{}{},
		regex:    regexp.MustCompile(regexPattern),
	}
	return &route
}

// CreateRouter creates new router component.
func (d *DefaultFactory) CreateRouter() IRouter {
	return &DefaultRouter{}
}

// CreateStore creates new store component.
func (d *DefaultFactory) CreateStore() IStore {
	return &DefaultMongoStore{}
}
