package oauth2

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	COPY    = "COPY"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	LINK    = "LINK"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PURGE   = "PURGE"
	PUT     = "PUT"
	UNLINK  = "UNLINK"
)

//var (
//	groups  []string = make([]string, 0)
//	methods []string = []string{COPY, DELETE, GET, HEAD, LINK, OPTIONS, PATCH, POST, PURGE, PUT, UNLINK}
//)

type Server struct {
	AllowMethods  []string
	StaticFolders []string

	routes []Route
	groups []string
	logger *log.Logger
}

// MARK: Struct's constructors
func DefaultServer() *Server {
	return &Server{
		AllowMethods:  []string{COPY, DELETE, GET, HEAD, LINK, OPTIONS, PATCH, POST, PURGE, PUT, UNLINK},
		StaticFolders: make([]string, 0),

		routes: make([]Route, 0),
		groups: make([]string, 0),
		logger: log.New(os.Stdout, "[OAuth2] ", 0),
	}
}

// MARK: Struct's public functions
func (s *Server) Run() {
	headersSize := GetEnv(ENV_HEADERS_SIZE)
	readTimeout := GetEnv(ENV_TIMEOUT_READ)
	writeTimeout := GetEnv(ENV_TIMEOUT_WRITE)

	address := fmt.Sprintf("%s:%s", GetEnv(ENV_HOST), GetEnv(ENV_PORT))
	t1, _ := strconv.ParseInt(headersSize, 10, 64)
	t2, _ := strconv.ParseInt(readTimeout, 10, 64)
	t3, _ := strconv.ParseInt(writeTimeout, 10, 64)

	server := &http.Server{
		Addr:           address,
		ReadTimeout:    time.Duration(t2) * time.Second,
		WriteTimeout:   time.Duration(t3) * time.Second,
		MaxHeaderBytes: int(t1) << 10, // 512kb
		Handler:        s,
	}

	s.logger.Printf("listening on %s\n", address)
	s.logger.Fatalln(server.ListenAndServe())
}
func (s *Server) RunTLS(certFile string, keyFile string) {
	headersSize := GetEnv(ENV_HEADERS_SIZE)
	readTimeout := GetEnv(ENV_TIMEOUT_READ)
	writeTimeout := GetEnv(ENV_TIMEOUT_WRITE)

	address := fmt.Sprintf("%s:%s", GetEnv(ENV_HOST), GetEnv(ENV_PORT))
	t1, _ := strconv.ParseInt(headersSize, 10, 64)
	t2, _ := strconv.ParseInt(readTimeout, 10, 64)
	t3, _ := strconv.ParseInt(writeTimeout, 10, 64)

	server := &http.Server{
		Addr:           address,
		ReadTimeout:    time.Duration(t2) * time.Second,
		WriteTimeout:   time.Duration(t3) * time.Second,
		MaxHeaderBytes: int(t1) << 10, // 512kb
		Handler:        s,
	}

	s.logger.Printf("listening on %s\n", address)
	s.logger.Fatalln(server.ListenAndServeTLS(certFile, keyFile))
}
