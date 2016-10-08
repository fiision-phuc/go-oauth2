package oauth2

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

// Router descripts a default router component implementation.
type Router struct {
	routes    []*Route
	groups    []string
	userRoles map[*regexp.Regexp]*regexp.Regexp
}

// GroupRoles groups all same url's prefix with user's roles.
func (r *Router) GroupRoles(groupPath string, roles ...string) {
	/* Condition validation: Ignore role validation if there is no token store */
	if TokenStore == nil {
		return
	}

	// Format pattern
	groupPath = pathFinder.ReplaceAllStringFunc(groupPath, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:len(m)-1])
	})
	groupPath = globsFinder.ReplaceAllStringFunc(groupPath, func(m string) string {
		return fmt.Sprintf(`(?P<_%d>[^#?]*)`, 0)
	})

	// Define validator
	if r.userRoles == nil {
		r.userRoles = make(map[*regexp.Regexp]*regexp.Regexp)
	}
	r.userRoles[regexp.MustCompile(groupPath)] = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(roles, "|")))
}

// BindRoles binds an url pattern with user's roles.
func (r *Router) BindRoles(httpMethod string, urlPattern string, roles ...string) {
	/* Condition validation: Ignore role validation if there is no token store */
	if TokenStore == nil || len(r.routes) == 0 {
		return
	}
}

// GroupRoute groups all same url's prefix.
func (r *Router) GroupRoute(s *Server, groupPath string, function func(s *Server)) {
	r.groups = append(r.groups, httprouter.CleanPath(groupPath))
	function(s)
	r.groups = r.groups[:len(r.groups)-1]
}

// BindRoute binds an url pattern with a handler.
func (r *Router) BindRoute(httpMethod string, urlPattern string, handler interface{}) {
	// Format url pattern before assigned to route
	if len(r.groups) > 0 {
		var buffer bytes.Buffer
		for _, path := range r.groups {
			buffer.WriteString(path)
		}

		if len(urlPattern) > 0 {
			buffer.WriteString(httprouter.CleanPath(urlPattern))
		}
		urlPattern = buffer.String()
	} else {
		urlPattern = httprouter.CleanPath(urlPattern)
	}
	logrus.Infof("%s -> %s", httpMethod, urlPattern)

	// Look for existing one before create new
	for _, route := range r.routes {
		if route.URLPattern() == urlPattern {
			route.BindHandler(httpMethod, handler)
			return
		}
	}

	// Create new route
	newRoute := objectFactory.CreateRoute(urlPattern)
	newRoute.BindHandler(httpMethod, handler)
	r.routes = append(r.routes, newRoute)
}

// MatchRoute matches a route with an url path.
func (r *Router) MatchRoute(context *Request, security *Security) (*Route, map[string]string) {
	// Validate user's authorized first
	isAuthorized := true
	for rule, roles := range r.userRoles {
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
