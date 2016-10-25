package oauth2

import (
	"net/http"
	"os"
	"strings"

	"github.com/phuc0302/go-oauth2/util"
)

// ServeHTTP handle HTTP request and HTTP response.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := createRequestContext(r, w)
	oauth := createOAuthContext(request)

	// Handle error
	defer recovery(request, s.sandbox)

	/* Condition validation: validate request method */
	if !methodsValidation.MatchString(request.Method) {
		panic(util.Status405())
	}

	// Should redirect request to static folder or not?
	if request.Method == Get && len(Cfg.StaticFolders) > 0 {
		for prefix, folder := range Cfg.StaticFolders {
			if path := request.Path; strings.HasPrefix(path, prefix) {
				path = strings.Replace(path, prefix, folder, 1)

				if file, err := os.Open(path); err == nil {
					defer file.Close()

					if info, _ := file.Stat(); !info.IsDir() {
						http.ServeContent(w, r, path, info.ModTime(), file)
						return
					}
				}
				panic(util.Status404())
			}
		}
	}

	// Find route to handle request
	if route, pathParams := s.router.matchRoute(request, oauth); route != nil {
		if pathParams != nil {
			request.PathParams = pathParams
		}
		route.invokeHandler(request, oauth)
		return
	}
	panic(util.Status503())
}
