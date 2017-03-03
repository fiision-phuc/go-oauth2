package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"
	"github.com/julienschmidt/httprouter"
	"github.com/phuc0302/go-oauth2/util"
)

// Server describes server object.
type Server struct {
	router *Router
}

// CreateServer returns a server with custom components.
//
// - parameter sandboxMode: instruct which config file should be loaded
func CreateServer(sandboxMode bool) *Server {
	// Load config file
	if sandboxMode {
		fmt.Println("Server is in sandboxMode.")
		Cfg = LoadConfig(debug)
	} else {
		fmt.Println("Server is in productionMode.")
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
			HookURL:        Cfg.SlackURL,
			Channel:        Cfg.SlackChannel,
			Username:       Cfg.SlackUser,
			IconEmoji:      Cfg.SlackIcon,
			AcceptedLevels: slackrus.LevelThreshold(level),
		})
	}

	// Create server
	server := Server{router: new(Router)}
	return &server
}

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

// ServeHTTP handle HTTP request and HTTP response.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer Recovery(w, r)
	method := strings.ToLower(r.Method)
	path := httprouter.CleanPath(r.URL.Path)

	/* Condition validation: validate request method */
	if !methodsValidation.MatchString(method) {
		panic(util.Status405())
	}

	// Find route to handle request
	if route, pathParams := s.router.MatchRoute(method, path); route != nil {
		context := CreateContext(w, r)
		if pathParams != nil {
			context.PathParams = pathParams
		}
		route.InvokeHandlers(context)
	} else {
		if len(Cfg.StaticFolders) > 0 && method == Get {
			for prefix, folder := range Cfg.StaticFolders {

				if strings.HasPrefix(path, prefix) {
					path = strings.Replace(path, prefix, folder, 1)

					if file, err := os.Open(path); err == nil {
						defer file.Close()

						if info, _ := file.Stat(); !info.IsDir() {
							http.ServeContent(w, r, path, info.ModTime(), file)
							return
						}
					}

					panic(util.Status404())
					return
				}
			}
		}
		panic(util.Status503())
	}
}

// MARK: Server's routing
// GroupRoute routes all url with same prefix.
func (s *Server) GroupRoute(urlPrefix string, handler HandleGroupFunc) {
	s.router.GroupRoute(s, urlPrefix, handler)
}

// Copy routes copy request to registered handler.
func (s *Server) Copy(urlPattern string, handler HandleContextFunc) {
	s.router.BindRoute(Copy, urlPattern, handler)
}

// Delete routes delete request to registered handler.
func (s *Server) Delete(urlPattern string, handler HandleContextFunc) {
	s.router.BindRoute(Delete, urlPattern, handler)
}

// Get routes get request to registered handler.
func (s *Server) Get(urlPattern string, handler HandleContextFunc) {
	s.router.BindRoute(Get, urlPattern, handler)
}

// Head routes head request to registered handler.
func (s *Server) Head(urlPattern string, handler HandleContextFunc) {
	s.router.BindRoute(Head, urlPattern, handler)
}

// Link routes link request to registered handler.
func (s *Server) Link(urlPattern string, handler HandleContextFunc) {
	s.router.BindRoute(Link, urlPattern, handler)
}

// Options routes options request to registered handler.
func (s *Server) Options(urlPattern string, handler HandleContextFunc) {
	s.router.BindRoute(Options, urlPattern, handler)
}

// Patch routes patch request to registered handler.
func (s *Server) Patch(urlPattern string, handler HandleContextFunc) {
	s.router.BindRoute(Patch, urlPattern, handler)
}

// Post routes post request to registered handler.
func (s *Server) Post(urlPattern string, handler HandleContextFunc) {
	s.router.BindRoute(Post, urlPattern, handler)
}

// Purge routes purge request to registered handler.
func (s *Server) Purge(urlPattern string, handler HandleContextFunc) {
	s.router.BindRoute(Purge, urlPattern, handler)
}

// Put routes put request to registered handler.
func (s *Server) Put(urlPattern string, handler HandleContextFunc) {
	s.router.BindRoute(Put, urlPattern, handler)
}

// Unlink routes unlink request to registered handler.
func (s *Server) Unlink(urlPattern string, handler HandleContextFunc) {
	s.router.BindRoute(Unlink, urlPattern, handler)
}
