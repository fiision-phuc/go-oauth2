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
	//	defer RecoveryRequest(request, response, s.sandbox)
	request.URL.Path = httprouter.CleanPath(request.URL.Path)
	request.Method = strings.ToLower(request.Method)

	/* Condition validation: validate request method */
	if !methodsValidation.MatchString(request.Method) {
		panic(util.Status405())
		return
		//		context.OutputError(util.Status405())
	}

	// Should redirect request to static folder or not?
	if request.Method == Get && len(Cfg.StaticFolders) > 0 {
		for prefix, folder := range Cfg.StaticFolders {
			if path := request.URL.Path; strings.HasPrefix(path, prefix) {
				path = strings.Replace(path, prefix, folder, 1)

				if file, err := os.Open(path); err == nil {
					defer file.Close()

					if info, _ := file.Stat(); !info.IsDir() {
						http.ServeContent(response, request, path, info.ModTime(), file)
						return
					}
				}
				panic(util.Status404())
				return
			}
		}
	}

	//	// Handle request
	//	var match mux.RouteMatch
	//	if matched := router.Match(request, &match); matched {
	//		match.Handler.ServeHTTP(response, request)
	//		return
	//	}
	panic(util.Status503())
}
