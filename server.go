package oauth2

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"
)

// Server describes server object.
type Server struct {
	sandbox bool
	router  *router
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// CreateServer returns a server with custom components.
func CreateServer(instance IFactory, isSandbox bool) *Server {
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
	objectFactory = instance
	TokenStore = instance.CreateStore()

	// Create server
	server := Server{
		sandbox: isSandbox,
		router:  new(router),
	}

	// Setup OAuth2.0
	if TokenStore != nil {
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
	factory := new(DefaultFactory)
	return CreateServer(factory, isSandbox)
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
	s.router.BindRoute(Copy, urlPattern, handler)
}

// Delete routes delete request to registered handler.
func (s *Server) Delete(urlPattern string, handler ContextHandler) {
	s.router.BindRoute(Delete, urlPattern, handler)
}

// Get routes get request to registered handler.
func (s *Server) Get(urlPattern string, handler ContextHandler) {
	s.router.BindRoute(Get, urlPattern, handler)
}

// Head routes head request to registered handler.
func (s *Server) Head(urlPattern string, handler ContextHandler) {
	s.router.BindRoute(Head, urlPattern, handler)
}

// Link routes link request to registered handler.
func (s *Server) Link(urlPattern string, handler ContextHandler) {
	s.router.BindRoute(Link, urlPattern, handler)
}

// Options routes options request to registered handler.
func (s *Server) Options(urlPattern string, handler ContextHandler) {
	s.router.BindRoute(Options, urlPattern, handler)
}

// Patch routes patch request to registered handler.
func (s *Server) Patch(urlPattern string, handler ContextHandler) {
	s.router.BindRoute(Patch, urlPattern, handler)
}

// Post routes post request to registered handler.
func (s *Server) Post(urlPattern string, handler ContextHandler) {
	s.router.BindRoute(Post, urlPattern, handler)
}

// Purge routes purge request to registered handler.
func (s *Server) Purge(urlPattern string, handler ContextHandler) {
	s.router.BindRoute(Purge, urlPattern, handler)
}

// Put routes put request to registered handler.
func (s *Server) Put(urlPattern string, handler ContextHandler) {
	s.router.BindRoute(Put, urlPattern, handler)
}

// Unlink routes unlink request to registered handler.
func (s *Server) Unlink(urlPattern string, handler ContextHandler) {
	s.router.BindRoute(Unlink, urlPattern, handler)
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
