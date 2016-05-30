package oauth2

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/phuc0302/go-oauth2/config"
	"github.com/phuc0302/go-oauth2/d"
	"github.com/phuc0302/go-oauth2/i"
	"github.com/phuc0302/go-oauth2/oauth"
)

// Server object description.
type Server struct {
	*config.Config

	router  i.IRouter
	factory i.IFactory

	logger *log.Logger

	tokenStore oauth.TokenStore
	userRoles  map[*regexp.Regexp][]string
}

// DefaultServer create a server object with preset config.
func DefaultServer() *Server {
	factory := &d.DefaultFactory{}
	server := &Server{
		Config: config.LoadConfigs(config.Debug),

		factory: factory,
		router:  factory.CreateRouter(),
	}
	return server
}

// DefaultServerWithTokenStore create a server object with preset config and oauth2.0 enabled.
func DefaultServerWithTokenStore(tokenStore oauth.TokenStore) *Server {
	cfg := config.LoadConfigs(config.Debug)

	server := &Server{
		Config: cfg,
		logger: log.New(os.Stdout, "[OAuth2.0] ", 0),
	}

	if tokenStore != nil {
		server.tokenStore = tokenStore

		// Pre-define oauth2 urls
		//	grantAuthorization := new(AuthorizationGrant)
		tokenGrant := oauth.CreateTokenGrant(cfg, tokenStore)

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
