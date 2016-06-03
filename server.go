package oauth2

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
)

var (
	cfg     *config
	factory IFactory
	store   IStore
)

// Server object description.
type Server struct {
	sandbox bool
	router  IRouter

	logger    *log.Logger
	userRoles map[*regexp.Regexp][]string
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
		cfg = loadConfig(debug)
	} else {
		cfg = loadConfig(release)
	}

	// Register components
	factory = instance
	store = factory.CreateStore()

	// Create server
	server := Server{
		router: factory.CreateRouter(),
		logger: log.New(os.Stdout, "[OAuth2.0] ", 0),
	}

	// Pre-define oauth2 urls
	if store != nil {
		//	grantAuthorization := new(AuthorizationGrant)
		tokenGrant := CreateTokenGrant(cfg, store)

		//	server.Get("/authorize", grantAuthorization.HandleForm)
		server.Post("/token", tokenGrant.HandleForm)
	}
	return &server
}

// Run will start server on http port.
func (s *Server) Run() {
	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	server := &http.Server{
		Addr:           address,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		MaxHeaderBytes: cfg.HeaderSize,
		Handler:        s,
	}

	s.logger.Printf("listening on %s\n", address)
	s.logger.Fatalln(server.ListenAndServe())
}

// RunTLS will start server on https port.
func (s *Server) RunTLS(certFile string, keyFile string) {
	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.TLSPort)
	server := &http.Server{
		Addr:           address,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		MaxHeaderBytes: cfg.HeaderSize,
		Handler:        s,
	}

	s.logger.Printf("listening on %s\n", address)
	s.logger.Fatalln(server.ListenAndServeTLS(certFile, keyFile))
}
