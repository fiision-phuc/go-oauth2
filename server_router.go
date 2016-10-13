package oauth2

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

/// Router descripts a default router component implementation.
type router struct {
	groups []string
	routes []IRoute
	roles  map[*regexp.Regexp]*regexp.Regexp
}

// groupRoles groups prefix url's with user's roles.
//
// @param groupPath:
func (r *router) groupRoles(groupPath string, roles ...string) {
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
	if r.roles == nil {
		r.roles = make(map[*regexp.Regexp]*regexp.Regexp)
	}
	r.roles[regexp.MustCompile(groupPath)] = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(roles, "|")))
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

// BindRoute binds an url pattern with a handler.
func (r *router) BindRoute(httpMethod string, path string, handler interface{}) {
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
	logrus.Infof("%s -> %s", httpMethod, path)

	// Look for existing one before create new
	for _, route := range r.routes {
		if route.URLPattern() == path {
			route.BindHandler(httpMethod, handler)
			return
		}
	}

	// Create new route
	newRoute := objectFactory.CreateRoute(path)
	newRoute.BindHandler(httpMethod, handler)
	r.routes = append(r.routes, newRoute)
}

// MatchRoute matches a route with an url path.
func (r *router) MatchRoute(context *Request, security *Security) (IRoute, map[string]string) {
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
