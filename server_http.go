package oauth2

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/phuc0302/go-oauth2/utils"
)

func (s *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	request.URL.Path = utils.FormatPath(request.URL.Path)
	request.Method = strings.ToUpper(request.Method)

	// Create context
	context := CreateContext(request, response)
	defer RecoveryRequest(context)

	// Validate http request methods
	isAllowed := false
	for _, allowMethod := range s.AllowMethods {
		if request.Method == allowMethod {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		context.OutputError(Status405())
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
func (s *Server) serveRequest(context *Context) {
	// FIX FIX FIX: Add priority here so that we can move the mosted used node to top
	isHandled := false
	for _, route := range s.routes {
		ok, pathParams := route.Match(context.Method, context.UrlPath)
		if !ok {
			continue
		}

		context.Params.PathQueries = pathParams
		route.InvokeHandler(context)
		isHandled = true
		break
	}

	if !isHandled {
		context.OutputError(Status503())
	}
}

func (s *Server) serveResource(context *Context, request *http.Request, response http.ResponseWriter) {
	resourcePath := request.URL.Path

	/* Condition validation: Check if file exist or not */
	if !utils.FileExisted(resourcePath) {
		context.OutputError(Status404())
		return
	}

	// Open file as read only
	file, err := os.Open(resourcePath)
	defer file.Close()

	if err != nil {
		context.OutputError(Status404())
		return
	}

	/* Condition validation: Only serve file, not directory */
	info, _ := file.Stat()
	if info.IsDir() {
		context.OutputError(Status403())
		return
	}
	http.ServeContent(response, request, resourcePath, info.ModTime(), file)
}
