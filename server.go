package oauth2

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/phuc0302/go-oauth2/util"
)

// Server describes server object.
type Server struct {
	sandbox bool
	router  *ServerRouter
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Run will start server on http port.
func (s *Server) Run() {
	address := fmt.Sprintf("%s:%d", Cfg.Host, Cfg.Port)
	server := &http.Server{
		Addr:           address,
		ReadTimeout:    Cfg.ReadTimeout,
		WriteTimeout:   Cfg.WriteTimeout,
		MaxHeaderBytes: Cfg.HeaderSize,
		Handler:        s,
	}
	logrus.Infof("listening on %s", address)
	logrus.Fatal(server.ListenAndServe())
}

// RunTLS will start server on https port.
func (s *Server) RunTLS(certFile string, keyFile string) {
	address := fmt.Sprintf("%s:%d", Cfg.Host, Cfg.TLSPort)
	server := &http.Server{
		Addr:           address,
		ReadTimeout:    Cfg.ReadTimeout,
		WriteTimeout:   Cfg.WriteTimeout,
		MaxHeaderBytes: Cfg.HeaderSize,
		Handler:        s,
	}
	logrus.Infof("listening on %s\n", address)
	logrus.Fatal(server.ListenAndServeTLS(certFile, keyFile))
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// GroupRoles binds user's roles to all url with same prefix.
func (s *Server) GroupRoles(groupPath string, roles ...string) {
	s.router.groupRoles(groupPath, roles...)
}

//// BindRoles an url pattern with user's roles.
//func (s *Server) BindRoles(httpMethod string, urlPattern string, roles ...string) {
//	s.router.BindRoles(httpMethod, urlPattern, roles...)
//}

// GroupRoute routes all url with same prefix.
func (s *Server) GroupRoute(urlPrefix string, handler GroupHandler) {
	s.router.groupRoute(s, urlPrefix, handler)
}

// Copy routes copy request to registered handler.
func (s *Server) Copy(urlPattern string, handler ContextHandler) {
	s.router.bindRoute(Copy, urlPattern, handler)
}

// Delete routes delete request to registered handler.
func (s *Server) Delete(urlPattern string, handler ContextHandler) {
	s.router.bindRoute(Delete, urlPattern, handler)
}

// Get routes get request to registered handler.
func (s *Server) Get(urlPattern string, handler ContextHandler) {
	s.router.bindRoute(Get, urlPattern, handler)
}

// Head routes head request to registered handler.
func (s *Server) Head(urlPattern string, handler ContextHandler) {
	s.router.bindRoute(Head, urlPattern, handler)
}

// Link routes link request to registered handler.
func (s *Server) Link(urlPattern string, handler ContextHandler) {
	s.router.bindRoute(Link, urlPattern, handler)
}

// Options routes options request to registered handler.
func (s *Server) Options(urlPattern string, handler ContextHandler) {
	s.router.bindRoute(Options, urlPattern, handler)
}

// Patch routes patch request to registered handler.
func (s *Server) Patch(urlPattern string, handler ContextHandler) {
	s.router.bindRoute(Patch, urlPattern, handler)
}

// Post routes post request to registered handler.
func (s *Server) Post(urlPattern string, handler ContextHandler) {
	s.router.bindRoute(Post, urlPattern, handler)
}

// Purge routes purge request to registered handler.
func (s *Server) Purge(urlPattern string, handler ContextHandler) {
	s.router.bindRoute(Purge, urlPattern, handler)
}

// Put routes put request to registered handler.
func (s *Server) Put(urlPattern string, handler ContextHandler) {
	s.router.bindRoute(Put, urlPattern, handler)
}

// Unlink routes unlink request to registered handler.
func (s *Server) Unlink(urlPattern string, handler ContextHandler) {
	s.router.bindRoute(Unlink, urlPattern, handler)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
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
	} else {
		panic(util.Status503())
	}
}
