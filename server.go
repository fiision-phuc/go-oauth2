package oauth2

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Server object description.
type Server struct {
	*Config

	routes []Route
	groups []string
	logger *log.Logger
}

// DefaultServer create a server object with preset config.
func DefaultServer() *Server {
	// Load configuration file
	config := LoadConfigs()
	if config == nil {
		config = CreateConfigs()
	}

	// Create default server
	server := &Server{
		Config: config,

		routes: make([]Route, 0),
		groups: make([]string, 0),
		logger: log.New(os.Stdout, "[OAuth2] ", 0),
	}

	// Pre-define oauth2 urls
	grantAuthorization := new(GrantAuthorization)
	grantToken := new(GrantToken)

	server.Get("/auth", grantAuthorization.HandleForm)
	server.Post("/token", grantToken.HandleForm)
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
	address := fmt.Sprintf("%s:%s", s.Host, s.Port)
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
