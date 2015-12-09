package oauth2

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

var (
	globsRegex     = regexp.MustCompile(`\*\*`)
	pathParamRegex = regexp.MustCompile(`:[^/#?()\.\\]+`)
)

// Server object description.
type Server struct {
	*Config

	routes []Route
	groups []string
	logger *log.Logger

	tokenStore TokenStore
	userRoles  map[*regexp.Regexp][]string
}

// DefaultServer create a server object with preset config.
func DefaultServer() *Server {
	return DefaultServerWithTokenStore(nil)
}

// DefaultServerWithTokenStore create a server object with preset config and oauth2.0 enabled.
func DefaultServerWithTokenStore(tokenStore TokenStore) *Server {
	config := LoadConfigs()

	server := &Server{
		Config: config,
		logger: log.New(os.Stdout, "[OAuth2.0] ", 0),
	}

	if tokenStore != nil {
		server.tokenStore = tokenStore

		// Pre-define oauth2 urls
		//	grantAuthorization := new(AuthorizationGrant)
		tokenGrant := CreateTokenGrant(config, tokenStore)

		//	server.Get("/auth", grantAuthorization.HandleForm)
		server.Post("/token", tokenGrant.HandleForm)
	}
	return server
}

// Run will start server on http port.
func (s *Server) Run() {
	address := fmt.Sprintf("%s:%s", s.Host, s.Port)
	server := &http.Server{
		Addr:           address,
		ReadTimeout:    s.TimeoutRead * time.Second,
		WriteTimeout:   s.TimeoutWrite * time.Second,
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
		ReadTimeout:    s.TimeoutRead * time.Second,
		WriteTimeout:   s.TimeoutWrite * time.Second,
		MaxHeaderBytes: s.HeaderSize << 10,
		Handler:        s,
	}

	s.logger.Printf("listening on %s\n", address)
	s.logger.Fatalln(server.ListenAndServeTLS(certFile, keyFile))
}
