package oauth2

import (
	"bytes"

	"github.com/phuc0302/go-oauth2/utils"
)

// Group all url with same prefix.
func (s *Server) Group(urlGroup string, function func(s *Server)) {
	s.groups = append(s.groups, urlGroup)
	function(s)

	s.groups = s.groups[:len(s.groups)-1]
}

// Copy method.
func (s *Server) Copy(urlPath string, handler interface{}) {
	defer RecoveryInternal(s.logger)
	s.addRoute(COPY, urlPath, handler)
}

// Delete method.
func (s *Server) Delete(urlPath string, handler interface{}) {
	defer RecoveryInternal(s.logger)
	s.addRoute(DELETE, urlPath, handler)
}

// Get method.
func (s *Server) Get(urlPath string, handler interface{}) {
	defer RecoveryInternal(s.logger)
	s.addRoute(GET, urlPath, handler)
}

// Head method.
func (s *Server) Head(urlPath string, handler interface{}) {
	defer RecoveryInternal(s.logger)
	s.addRoute(HEAD, urlPath, handler)
}

// Link method.
func (s *Server) Link(urlPath string, handler interface{}) {
	defer RecoveryInternal(s.logger)
	s.addRoute(LINK, urlPath, handler)
}

// Options method.
func (s *Server) Options(urlPath string, handler interface{}) {
	defer RecoveryInternal(s.logger)
	s.addRoute(OPTIONS, urlPath, handler)
}

// Patch method.
func (s *Server) Patch(urlPath string, handler interface{}) {
	defer RecoveryInternal(s.logger)
	s.addRoute(PATCH, urlPath, handler)
}

// Post method.
func (s *Server) Post(urlPath string, handler interface{}) {
	defer RecoveryInternal(s.logger)
	s.addRoute(POST, urlPath, handler)
}

// Purge method.
func (s *Server) Purge(urlPath string, handler interface{}) {
	defer RecoveryInternal(s.logger)
	s.addRoute(PURGE, urlPath, handler)
}

// Put method.
func (s *Server) Put(urlPath string, handler interface{}) {
	defer RecoveryInternal(s.logger)
	s.addRoute(PUT, urlPath, handler)
}

// Unlink method.
func (s *Server) Unlink(urlPath string, handler interface{}) {
	defer RecoveryInternal(s.logger)
	s.addRoute(UNLINK, urlPath, handler)
}

// MARK: Struct's private functions
func (s *Server) addRoute(method string, pattern string, handler interface{}) {
	// Format pattern before assigned to route
	if len(s.groups) > 0 {
		var groupPattern bytes.Buffer

		for _, g := range s.groups {
			groupPattern.WriteString(utils.FormatPath(g))
		}

		if len(pattern) > 0 {
			groupPattern.WriteString(utils.FormatPath(pattern))
		}
		pattern = groupPattern.String()
	} else {
		pattern = utils.FormatPath(pattern)
	}

	// Look for existing one before create new
	for _, route := range s.routes {
		if route.GetPattern() == pattern {
			route.AddHandler(method, handler)
			s.logger.Printf("%-6s -> %s\n", method, pattern)
			return
		}
	}

	// Create new route
	newRoute := CreateDefaultRoute(pattern)
	newRoute.AddHandler(method, handler)

	// Add to collection
	s.routes = append(s.routes, newRoute)
	s.logger.Printf("%-6s -> %s\n", method, pattern)
}
