package oauth2

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"
	"github.com/julienschmidt/httprouter"
	"github.com/phuc0302/go-mongo"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// CreateServer returns a server with custom components.
func CreateServer(tokenStore TokenStore, isSandbox bool) *Server {
	// Load config file
	if isSandbox {
		Cfg = LoadConfig(debug)
	} else {
		Cfg = LoadConfig(release)
	}

	// Setup logger
	level, err := logrus.ParseLevel(Cfg.LogLevel)
	if err != nil {
		level = logrus.DebugLevel
	}
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(level)

	// Setup slack notification if neccessary
	if len(Cfg.SlackURL) > 0 {
		logrus.AddHook(&slackrus.SlackrusHook{
			HookURL:        Cfg.SlackURL,     // "https://hooks.slack.com/services/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"
			Channel:        Cfg.SlackChannel, // "#Oauth2.0"
			Username:       Cfg.SlackUser,    // "Server"
			IconEmoji:      Cfg.SlackIcon,    // ":ghost:"
			AcceptedLevels: slackrus.LevelThreshold(level),
		})
	}

	// Register global components
	Store = tokenStore

	// Create server
	server := Server{
		sandbox: isSandbox,
		router:  new(ServerRouter),
	}

	// Setup OAuth2.0
	if Store != nil {
		//	grantAuthorization := new(AuthorizationGrant)
		tokenGrant := new(TokenGrant)

		//	server.Get("/authorize", grantAuthorization.HandleForm)
		server.Get("/token", tokenGrant.HandleForm)
		server.Post("/token", tokenGrant.HandleForm)
	}
	return &server
}

// DefaultServer returns a server with build in components.
func DefaultServer(isSandbox bool) *Server {
	mongo.ConnectMongo()
	tokenStore := new(DefaultMongoStore)
	return CreateServer(tokenStore, isSandbox)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// createOAuthContext creates new security context.
func createOAuthContext(c *RequestContext) *OAuthContext {
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
	if accessToken := Store.FindAccessToken(tokenString); accessToken != nil && !accessToken.IsExpired() {
		client := Store.FindClientWithID(accessToken.ClientID())
		user := Store.FindUserWithID(accessToken.UserID())
		return &OAuthContext{
			Client:      client,
			User:        user,
			AccessToken: accessToken,
		}
	}

	/* Condition validation: If everything is not work out, try to look for basic auth */
	if username, password, ok := c.BasicAuth(); ok {
		client := Store.FindClientWithCredential(username, password)
		user := Store.FindUserWithClient(username, password)

		if client != nil && user != nil {
			return &OAuthContext{
				Client: client,
				User:   user,
			}
		}
	}
	return nil
}

// createRequestContext creates new request context.
func createRequestContext(request *http.Request, response http.ResponseWriter) *RequestContext {
	context := &RequestContext{
		Path:     httprouter.CleanPath(request.URL.Path),
		Method:   strings.ToLower(request.Method),
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
	switch context.Method {

	case Get:
		params = request.URL.Query()
		break

	case Patch, Post:
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

////////////////////////////////////////////////////////////////////////////////////////////////////
// createRoute creates new route component.
func createRoute(path string) *ServerRoute {
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

	route := &ServerRoute{
		path:     path,
		handlers: make(map[string]ContextHandler),
		regex:    regexp.MustCompile(regexPattern),
	}
	return route
}

// createRouter creates new router component.
func createRouter() *ServerRouter {
	return new(ServerRouter)
}
