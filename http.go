package oauth2

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/phuc0302/go-oauth2/utils"
)

// ServeHTTP handle HTTP request and HTTP response.
func (s *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	request.URL.Path = utils.FormatPath(request.URL.Path)
	request.Method = strings.ToUpper(request.Method)

	// Create request context
	context := objectFactory.CreateRequestContext(request, response)
	defer RecoveryRequest(context, s.sandbox)

	/* Condition validation: validate request methods */
	if !methodsValidation.MatchString(request.Method) {
		context.OutputError(utils.Status405())
		return
	}

	// Should redirect request to static folder or not?
	if request.Method == GET && len(cfg.StaticFolders) > 0 {
		for prefix, folder := range cfg.StaticFolders {
			if strings.HasPrefix(request.URL.Path, prefix) {
				newPath := strings.Replace(request.URL.Path, prefix, folder, 1)
				request.URL, _ = url.Parse(newPath)

				s.serveResource(context, request, response)
				return
			}
		}
	}

	// Create security context
	security := objectFactory.CreateSecurityContext(context)
	s.serveRequest(context, security)
}

func (s *Server) serveRequest(context *Request, security *Security) {
	// FIX FIX FIX: Add priority here so that we can move the mosted used node to top

	if route, pathParams := s.router.MatchRoute(context, security); route != nil {
		context.PathParams = pathParams
		route.InvokeHandler(context, security)
	} else {
		context.OutputError(utils.Status503())
	}
}

func (s *Server) serveResource(context *Request, request *http.Request, response http.ResponseWriter) {
	/* Condition validation: Check if file exist or not */
	if resourcePath := request.URL.Path; !utils.FileExisted(resourcePath) {
		context.OutputError(utils.Status404())
	} else {
		if file, err := os.Open(resourcePath); err != nil {
			context.OutputError(utils.Status404())
		} else {
			defer file.Close()

			/* Condition validation: Only serve file, not directory */
			if info, _ := file.Stat(); info.IsDir() {
				context.OutputError(utils.Status403())
			} else {
				http.ServeContent(response, request, resourcePath, info.ModTime(), file)
			}
		}
	}
}
