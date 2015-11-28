package oauth2

import (
	"net/http"
	"os"

	"github.com/phuc0302/go-oauth2/utils"
)

// MARK: http.Handler's members
func (s *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	//	request.URL.Path = utils.FormatPath(request.URL.Path)
	//	request.Method = strings.ToUpper(request.Method)

	//	// Create context
	//	context := CreateContext(request, response)
	//	defer RecoveryRequest(context)

	//	if len(s.StaticFolder) > 0 && strings.HasPrefix(request.URL.Path, s.StaticFolder) {
	//		s.serveResource(context, request, response)
	//	} else {
	//		s.serveRequest(context)
	//	}
}

// MARK: Struct's private functions
func (s *Server) serveRequest(context *Context) {
	//	isAllowed := s.Router.ShouldAllow(context.Method)
	//	if !isAllowed {
	//		context.OutputError(Status405())
	//		return
	//	}

	//	// Let delegate decide if the request should be handled or not
	//	if s.Delegate != nil && !s.Delegate.ShouldServeHTTP(context) {
	//		context.OutputError(Status404())
	//		return
	//	} else if s.Delegate != nil {
	//		s.Delegate.WillServeHTTP(context)
	//	}

	//	isHandled := s.Router.ServeRequest(context)
	//	if !isHandled {
	//		context.OutputError(Status404())
	//	}

	//	if s.Delegate != nil {
	//		s.Delegate.DidServeHTTP(context)
	//	}
}

func (s *Server) serveResource(context *Context, request *http.Request, response http.ResponseWriter) {
	/* Condition validation: Only GET is accepted when request a static resources */
	if request.Method != GET {
		context.OutputError(Status405())
		return
	}
	resourcePath := request.URL.Path[1:]

	/* Condition validation: Check if file exist or not */
	if !utils.FileExisted(resourcePath) {
		context.OutputError(Status404())
		return
	}

	// Open file as read only
	f, err := os.OpenFile(resourcePath, os.O_RDONLY, 0)
	defer f.Close()

	if err != nil {
		context.OutputError(Status404())
		return
	}

	/* Condition validation: Only serve file, not directory */
	fi, _ := f.Stat()
	if fi.IsDir() {
		context.OutputError(Status403())
		return
	}

	http.ServeContent(response, request, resourcePath, fi.ModTime(), f)
}
