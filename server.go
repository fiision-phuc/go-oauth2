package oauth2

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/johntdyer/slackrus"
)

var router = new(mux.Router)

// Server describes server object.
type Server struct {
	sandbox bool

	routerOld *Router
}

// DefaultServer returns a server with build in components.
func DefaultServer(isSandbox bool) *Server {
	factory := &DefaultFactory{}
	return CreateServer(factory, isSandbox)
}

// CreateServer create a server object with preset config and oauth2.0 enabled.
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

	// Enable slack notification if neccessary
	if len(Cfg.SlackURL) > 0 {
		logrus.AddHook(&slackrus.SlackrusHook{
			HookURL:        Cfg.SlackURL,     // "https://hooks.slack.com/services/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"
			Channel:        Cfg.SlackChannel, // "#keywords"
			Username:       Cfg.SlackUser,    // "Server"
			IconEmoji:      Cfg.SlackIcon,    // ":ghost:"
			AcceptedLevels: slackrus.LevelThreshold(level),
		})
	}

	// Register components
	objectFactory = instance
	TokenStore = objectFactory.CreateStore()

	// Create server
	server := Server{
		sandbox:   isSandbox,
		routerOld: objectFactory.CreateRouter(),
	}

	// Pre-define oauth2 urls
	if TokenStore != nil {
		//	grantAuthorization := new(AuthorizationGrant)
		tokenGrant := new(TokenGrant)

		//	server.Get("/authorize", grantAuthorization.HandleForm)
		server.Get("/token", tokenGrant.HandleForm)
		server.Post("/token", tokenGrant.HandleForm)
	}
	return &server
}

// GroupRoles binds user's roles to all url with same prefix.
func (s *Server) GroupRoles(groupPath string, roles ...string) {
	s.routerOld.GroupRoles(groupPath, roles...)
}

// BindRoles an url pattern with user's roles.
func (s *Server) BindRoles(httpMethod string, urlPattern string, roles ...string) {
	s.routerOld.BindRoles(httpMethod, urlPattern, roles...)
}

// GroupRoute routes all url with same prefix.
func (s *Server) GroupRoute(urlGroup string, function func(s *Server)) {
	s.routerOld.GroupRoute(s, urlGroup, function)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Server's listener.

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
// Server's routing apis.

// Copy routes copy request to registered handler.
func (s *Server) Copy(urlPattern string, handler func(requestContext *Request, securityContext *Security)) {
	router.HandleFunc(urlPattern, func(response http.ResponseWriter, request *http.Request) {
		context := objectFactory.CreateRequestContext(request, response)
		security := objectFactory.CreateSecurityContext(context)
		handler(context, security)
	}).Methods(Copy)
}

// Delete routes delete request to registered handler.
func (s *Server) Delete(urlPattern string, handler func(requestContext *Request, securityContext *Security)) {
	router.HandleFunc(urlPattern, func(response http.ResponseWriter, request *http.Request) {
		context := objectFactory.CreateRequestContext(request, response)
		security := objectFactory.CreateSecurityContext(context)
		handler(context, security)
	}).Methods(Delete)
}

// Get routes get request to registered handler.
func (s *Server) Get(urlPattern string, handler func(requestContext *Request, securityContext *Security)) {
	router.HandleFunc(urlPattern, func(response http.ResponseWriter, request *http.Request) {
		context := objectFactory.CreateRequestContext(request, response)
		security := objectFactory.CreateSecurityContext(context)
		handler(context, security)
	}).Methods(Get)
}

// Head routes head request to registered handler.
func (s *Server) Head(urlPattern string, handler func(requestContext *Request, securityContext *Security)) {
	router.HandleFunc(urlPattern, func(response http.ResponseWriter, request *http.Request) {
		context := objectFactory.CreateRequestContext(request, response)
		security := objectFactory.CreateSecurityContext(context)
		handler(context, security)
	}).Methods(Head)
}

// Link routes link request to registered handler.
func (s *Server) Link(urlPattern string, handler func(requestContext *Request, securityContext *Security)) {
	router.HandleFunc(urlPattern, func(response http.ResponseWriter, request *http.Request) {
		context := objectFactory.CreateRequestContext(request, response)
		security := objectFactory.CreateSecurityContext(context)
		handler(context, security)
	}).Methods(Link)
}

// Options routes options request to registered handler.
func (s *Server) Options(urlPattern string, handler func(requestContext *Request, securityContext *Security)) {
	router.HandleFunc(urlPattern, func(response http.ResponseWriter, request *http.Request) {
		context := objectFactory.CreateRequestContext(request, response)
		security := objectFactory.CreateSecurityContext(context)
		handler(context, security)
	}).Methods(Options)
}

// Patch routes patch request to registered handler.
func (s *Server) Patch(urlPattern string, handler func(requestContext *Request, securityContext *Security)) {
	router.HandleFunc(urlPattern, func(response http.ResponseWriter, request *http.Request) {
		context := objectFactory.CreateRequestContext(request, response)
		security := objectFactory.CreateSecurityContext(context)
		handler(context, security)
	}).Methods(Patch)
}

// Post routes post request to registered handler.
func (s *Server) Post(urlPattern string, handler func(requestContext *Request, securityContext *Security)) {
	router.HandleFunc(urlPattern, func(response http.ResponseWriter, request *http.Request) {
		context := objectFactory.CreateRequestContext(request, response)
		security := objectFactory.CreateSecurityContext(context)
		handler(context, security)
	}).Methods(Unlink)
}

// Purge routes purge request to registered handler.
func (s *Server) Purge(urlPattern string, handler func(requestContext *Request, securityContext *Security)) {
	router.HandleFunc(urlPattern, func(response http.ResponseWriter, request *http.Request) {
		context := objectFactory.CreateRequestContext(request, response)
		security := objectFactory.CreateSecurityContext(context)
		handler(context, security)
	}).Methods(Unlink)
}

// Put routes put request to registered handler.
func (s *Server) Put(urlPattern string, handler func(requestContext *Request, securityContext *Security)) {
	router.HandleFunc(urlPattern, func(response http.ResponseWriter, request *http.Request) {
		context := objectFactory.CreateRequestContext(request, response)
		security := objectFactory.CreateSecurityContext(context)
		handler(context, security)
	}).Methods(Unlink)
}

// Unlink routes unlink request to registered handler.
func (s *Server) Unlink(urlPattern string, handler func(requestContext *Request, securityContext *Security)) {
	router.HandleFunc(urlPattern, func(response http.ResponseWriter, request *http.Request) {
		context := objectFactory.CreateRequestContext(request, response)
		security := objectFactory.CreateSecurityContext(context)
		handler(context, security)
	}).Methods(Unlink)
}
