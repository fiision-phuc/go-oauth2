package oauth2

import (
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/phuc0302/go-oauth2/util"
)

// ServeHTTP handle HTTP request and HTTP response.
func (s *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	request.URL.Path = httprouter.CleanPath(request.URL.Path)
	request.Method = strings.ToLower(request.Method)

	// Create request context
	context := objectFactory.CreateRequestContext(request, response)

	/* Condition validation: validate HTTP method */
	if !methodsValidation.MatchString(request.Method) {
		context.OutputError(util.Status405())
		return
	}
	defer RecoveryRequest(context, s.sandbox)

	// Should redirect request to static folder or not?
	if request.Method == Get && len(Cfg.StaticFolders) > 0 {
		for prefix, folder := range Cfg.StaticFolders {

			if strings.HasPrefix(context.Path, prefix) {
				context.Path = strings.Replace(context.Path, prefix, folder, 1)
				//				newPath := strings.Replace(context.Path, prefix, folder, 1)
				//				request.URL, _ = url.Parse(newPath)

				s.serveResource(context)
				return
			}
		}
	}

	// Create security context
	security := objectFactory.CreateSecurityContext(context)
	s.serveRequest(context, security)
}

func (s *Server) serveRequest(context *Request, security *Security) {
	if route, pathParams := s.router.MatchRoute(context, security); route != nil {
		context.PathParams = pathParams

		route.InvokeHandler(context, security)
		return
	}
	context.OutputError(util.Status503())
}

func (s *Server) serveResource(context *Request) {
	if file, err := os.Open(context.Path); err == nil {
		defer file.Close()

		if info, _ := file.Stat(); !info.IsDir() {
			http.ServeContent(context.response, context.request, context.Path, info.ModTime(), file)
			return
		}
	}
	context.OutputError(util.Status404())
}
