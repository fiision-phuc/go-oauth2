package oauth2

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/phuc0302/go-oauth2/utils"
)

//route := route{method, nil, handlers, pattern, ""}

//	route.regex = regexp.MustCompile(pattern)

// AddRole apply role to specific pattern
func (s *Server) AddRoles(pattern string, roles string) {
	pattern = utils.FormatPath(pattern)

	pattern = pathParamRegex.ReplaceAllStringFunc(pattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	})
	pattern = globsRegex.ReplaceAllStringFunc(pattern, func(m string) string {
		return fmt.Sprintf(`(?P<_%d>[^#?]*)`, 0)
	})
	pattern += `\/?`

	userRoles := strings.Split(roles, ",")

	if s.userRoles == nil {
		s.userRoles = make(map[*regexp.Regexp][]string, 1)
	}
	s.userRoles[regexp.MustCompile(pattern)] = userRoles
}

// ServeHTTP handle HTTP request and HTTP response
func (s *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	request.URL.Path = utils.FormatPath(request.URL.Path)
	request.Method = strings.ToUpper(request.Method)

	// Create context
	context := CreateRequestContext(request, response)
	defer RecoveryRequest(context, s.Development)

	// Validate http request methods
	if !s.methodsValidation.MatchString(request.Method) {
		context.OutputError(utils.Status405())
		return
	}

	// Should redirect request to static folder or not?
	isStaticRequest := false
	if request.Method == GET && len(s.StaticFolders) > 0 {
		for prefix, folder := range s.StaticFolders {
			if strings.HasPrefix(request.URL.Path, prefix) {
				newPath := strings.Replace(request.URL.Path, prefix, folder, 1)
				request.URL, _ = url.Parse(newPath)
				isStaticRequest = true
				break
			}
		}
	}

	if !isStaticRequest {
		s.serveRequest(context)
	} else {
		s.serveResource(context, request, response)
	}
}

// MARK: Struct's private functions
func (s *Server) serveRequest(context *RequestContext) {
	// FIX FIX FIX: Add priority here so that we can move the mosted used node to top
	isHandled := false

	for _, route := range s.routes {
		ok, pathQueries := route.Match(context.Method(), context.URLPath)
		if !ok {
			continue
		}

		context.PathQueries = pathQueries

		// Validate authentication & roles if neccessary
		securityContext, status := CreateSecurityContextWithRequestContext(context, s.tokenStore)
		for rule, _ := range s.userRoles {
			if rule.MatchString(context.URLPath) {
				if securityContext.AuthUser != nil {

				} else {
					context.OutputError(status)
					return
				}
				break
			}
		}

		route.InvokeHandler(context)
		isHandled = true
		break
	}

	if !isHandled {
		context.OutputError(utils.Status503())
	}
}

func (s *Server) serveResource(context *RequestContext, request *http.Request, response http.ResponseWriter) {
	resourcePath := request.URL.Path

	/* Condition validation: Check if file exist or not */
	if !utils.FileExisted(resourcePath) {
		context.OutputError(utils.Status404())
		return
	}

	// Open file as read only
	file, err := os.Open(resourcePath)
	defer file.Close()

	if err != nil {
		context.OutputError(utils.Status404())
		return
	}

	/* Condition validation: Only serve file, not directory */
	info, _ := file.Stat()
	if info.IsDir() {
		context.OutputError(utils.Status403())
		return
	}
	http.ServeContent(response, request, resourcePath, info.ModTime(), file)
}
