package oauth2

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

// Server object description.
type Server struct {
	*config

	router  IRouter
	factory IFactory

	logger *log.Logger

	tokenStore IStore
	userRoles  map[*regexp.Regexp][]string
}

// DefaultServer create a server object with preset config.
func DefaultServer() *Server {
	factory := &DefaultFactory{}
	server := &Server{
		config: loadConfig(debug),

		factory: factory,
		router:  factory.CreateRouter(),
	}
	return server
}

// DefaultServerWithTokenStore create a server object with preset config and oauth2.0 enabled.
func DefaultServerWithTokenStore(tokenStore IStore) *Server {
	cfg := loadConfig(debug)

	server := &Server{
		config: cfg,
		logger: log.New(os.Stdout, "[OAuth2.0] ", 0),
	}

	if tokenStore != nil {
		server.tokenStore = tokenStore

		// Pre-define oauth2 urls
		//	grantAuthorization := new(AuthorizationGrant)
		tokenGrant := CreateTokenGrant(cfg, tokenStore)

		//	server.Get("/authorize", grantAuthorization.HandleForm)
		server.Post("/token", tokenGrant.HandleForm)
	}
	return server
}

// Run will start server on http port.
func (s *Server) Run() {
	address := fmt.Sprintf("%s:%s", s.Host, s.Port)
	server := &http.Server{
		Addr:           address,
		ReadTimeout:    s.ReadTimeout * time.Second,
		WriteTimeout:   s.WriteTimeout * time.Second,
		MaxHeaderBytes: s.HeaderSize << 10,
		Handler:        s,
	}

	s.logger.Printf("listening on %s\n", address)
	s.logger.Fatalln(server.ListenAndServe())
}

// RunTLS will start server on https port.
func (s *Server) RunTLS(certFile string, keyFile string) {
	address := fmt.Sprintf("%s:%s", s.Host, s.TLSPort)
	server := &http.Server{
		Addr:           address,
		ReadTimeout:    s.ReadTimeout * time.Second,
		WriteTimeout:   s.WriteTimeout * time.Second,
		MaxHeaderBytes: s.HeaderSize << 10,
		Handler:        s,
	}

	s.logger.Printf("listening on %s\n", address)
	s.logger.Fatalln(server.ListenAndServeTLS(certFile, keyFile))
}
