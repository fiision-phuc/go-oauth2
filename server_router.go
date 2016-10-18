package oauth2

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

/// router describes a router component implementation.
type router struct {
	groups []string
	routes []*route
	roles  map[*regexp.Regexp]*regexp.Regexp
}

// groupRoles groups prefix url's with user's roles.
//
// @param pathPrefix:
func (r *router) groupRoles(pathPrefix string, roles ...string) {
	/* Condition validation: Ignore role validation if there is no token store */
	if TokenStore == nil {
		return
	}

	// Format pattern
	pathPrefix = pathFinder.ReplaceAllStringFunc(pathPrefix, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:len(m)-1])
	})
	pathPrefix = globsFinder.ReplaceAllStringFunc(pathPrefix, func(m string) string {
		return fmt.Sprintf(`(?P<_%d>[^#?]*)`, 0)
	})

	// Define validator
	if r.roles == nil {
		r.roles = make(map[*regexp.Regexp]*regexp.Regexp)
	}
	r.roles[regexp.MustCompile(pathPrefix)] = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(roles, "|")))
}

// BindRoles binds an url pattern with user's roles.
//func (r *DefaultRouter) BindRoles(httpMethod string, urlPattern string, roles ...string) {
//	/* Condition validation: Ignore role validation if there is no token store */
//	if TokenStore == nil || len(r.routes) == 0 {
//		return
//	}
//}

// groupRoute generates path's prefix for following urls.
//
// @param s
// @param pathPrefix
// @param handler
func (r *router) groupRoute(s *Server, pathPrefix string, handler GroupHandler) {
	r.groups = append(r.groups, pathPrefix)
	handler(s)
	r.groups = r.groups[:len(r.groups)-1]
}

// bindRoute binds a path with a handler.
func (r *router) bindRoute(method string, path string, handler ContextHandler) {
	// Format url pattern before assigned to route
	if len(r.groups) > 0 {
		var buffer bytes.Buffer
		for _, path := range r.groups {
			buffer.WriteString(path)
		}

		if len(path) > 0 {
			buffer.WriteString(httprouter.CleanPath(path))
		}
		path = buffer.String()
	} else {
		path = httprouter.CleanPath(path)
	}
	logrus.Infof("%07s -> %s", strings.ToUpper(method), path)

	// Look for existing one before create new
	for _, route := range r.routes {
		if route.path == path {
			route.BindHandler(method, handler)
			return
		}
	}

	// Create new route
	newRoute := createRoute(path)
	newRoute.BindHandler(method, handler)
	r.routes = append(r.routes, newRoute)
}

// matchRoute matches a route with a path.
func (r *router) matchRoute(context *Request, security *Security) (*route, map[string]string) {
	// Validate user's authorized first
	isAuthorized := true
	for rule, roles := range r.roles {
		if rule.MatchString(context.Path) {
			isAuthorized = false

			if TokenStore != nil && security != nil && security.User != nil {
				for _, role := range security.User.UserRoles() {
					if roles.MatchString(role) {
						isAuthorized = true
						break // If user is authorized, break the loop
					}
				}
			} else {
				return nil, nil
			}
		}
	}

	/* Condition validation: Validate user's authorized */
	if !isAuthorized {
		return nil, nil
	}

	// Match route
	for _, route := range r.routes {
		if ok, pathParams := route.MatchURLPattern(context.request.Method, context.Path); ok {
			return route, pathParams
		}
	}
	return nil, nil
}
