package server

import (
	"bytes"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/phuc0302/go-oauth2/util"
)

// Router describes a router component implementation.
type Router struct {
	groups []string
	routes []*Route
}

// DefaultRouter creates new router component.
func DefaultRouter() *Router {
	return new(Router)
}

// GroupRoute generates path's prefix for following urls.
func (r *Router) GroupRoute(s *Server, pathPrefix string, handler HandleGroupFunc) {
	r.groups = append(r.groups, pathPrefix)
	handler(s)
	r.groups = r.groups[:len(r.groups)-1]
}

// BindRoute binds a path with handler.
func (r *Router) BindRoute(method string, path string, handler HandleContextFunc) {
	path = r.mergeGroup(path)
	logrus.Infof("%-6s -> %s", strings.ToUpper(method), path)

	// Define regex pattern
	regexPattern := util.ConvertPath(path)

	// Look for existing one before create new
	for _, route := range r.routes {
		if route.regex.String() == regexPattern {
			route.BindHandler(method, handler)
			return
		}
	}
	newRoute := DefaultRoute(regexPattern)
	newRoute.BindHandler(method, handler)

	// Append to current list
	r.routes = append(r.routes, newRoute)
}

// MatchRoute matches a route with a path.
func (r *Router) MatchRoute(method string, path string) (*Route, map[string]string) {
	// Match route
	for _, route := range r.routes {
		if ok, pathParams := route.Match(method, path); ok {
			return route, pathParams
		}
	}
	return nil, nil
}

// mergeGroup constructs path from multiple parts.
//
// - parameter path: final path
func (r *Router) mergeGroup(path string) string {
	if len(r.groups) > 0 {
		var buffer bytes.Buffer
		for _, path := range r.groups {
			buffer.WriteString(path)
		}

		if len(path) > 0 {
			buffer.WriteString(path)
		}
		path = buffer.String()
	}
	return httprouter.CleanPath(path)
}
