package oauth2

import (
	"net/http"
	"net/url"
	"strings"
)

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
			if header := strings.ToLower(k); header == "authorization" {
				context.Header[header] = v[0]
			} else {
				context.Header[header] = strings.ToLower(v[0])
			}
		}
	}

	// Parse body context if neccessary
	var params url.Values
	switch context.request.Method {

	case Get:
		params = request.URL.Query()
		break

	case Post, Patch:
		if contentType := context.Header["content-type"]; contentType == "application/x-www-form-urlencoded" {
			if err := request.ParseForm(); err == nil {
				params = request.Form
			}
		} else if strings.HasPrefix(contentType, "multipart/form-data; boundary") {
			if err := request.ParseMultipartForm(Cfg.MultipartSize); err == nil {
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
func (d *DefaultFactory) CreateSecurityContext(c *Request) *Security {
	tokenString := c.Header["authorization"]

	/* Condition validation: Validate existing of authorization header */
	if isBearer := bearerFinder.MatchString(tokenString); isBearer {
		tokenString = tokenString[7:]
	} else {
		if tokenString = c.QueryParams["access_token"]; len(tokenString) > 0 {
			delete(c.QueryParams, "access_token")
		}
	}

	/* Condition validation: Validate expiration time */
	if accessToken := store.FindAccessToken(tokenString); accessToken != nil && !accessToken.IsExpired() {
		client := store.FindClientWithID(accessToken.ClientID())
		user := store.FindUserWithID(accessToken.UserID())
		return &Security{
			Client:      client,
			User:        user,
			AccessToken: accessToken,
		}
	}

	/* Condition validation: If everything is not work out, try to look for basic auth */
	if username, password, ok := c.BasicAuth(); ok {
		client := store.FindClientWithCredential(username, password)
		user := store.FindUserWithClient(username, password)

		if client != nil && user != nil {
			return &Security{
				Client: client,
				User:   user,
			}
		}
	}
	return nil
}

//// CreateRoute creates new route component.
//func (d *DefaultFactory) CreateRoute(urlPattern string) IRoute {
//	return nil
//	//	regexPattern := pathFinder.ReplaceAllStringFunc(urlPattern, func(m string) string {
//	//		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:len(m)-1])
//	//	})
//	//	regexPattern = globsFinder.ReplaceAllStringFunc(regexPattern, func(m string) string {
//	//		return fmt.Sprintf(`(?P<_%d>[^#?]*)`, 0)
//	//	})
//	//	if len(regexPattern) == 1 && regexPattern == "/" {
//	//		regexPattern = fmt.Sprintf("^%s?$", regexPattern)
//	//	} else {
//	//		regexPattern = fmt.Sprintf("^%s/?$", regexPattern)
//	//	}

//	//	route := route{
//	//		path:     urlPattern,
//	//		handlers: map[string]ContextHandler{},
//	//		regex:    regexp.MustCompile(regexPattern),
//	//	}
//	//	return &route
//}

//// CreateRouter creates new router component.
//func (d *DefaultFactory) CreateRouter() ServerRouter {
//	return nil
//	//	return ServerRouter{}
//}
