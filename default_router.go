package oauth2

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/phuc0302/go-oauth2/utils"
)

// DefaultRouter descripts a default router component implementation.
type DefaultRouter struct {
	routes    []IRoute
	groups    []string
	userRoles map[*regexp.Regexp][]string
}

// GroupRole groups all same url's prefix with user's roles.
func (r *DefaultRouter) GroupRole(s *Server, groupPath string, roles ...string) {
	/* Condition validation: Ignore role validation if there is no token store */
	if TokenStore == nil {
		return
	}

	// Format path
	groupPath = pathFinder.ReplaceAllStringFunc(groupPath, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:len(m)-1])
	})
	groupPath = globsFinder.ReplaceAllStringFunc(groupPath, func(m string) string {
		return fmt.Sprintf(`(?P<_%d>[^#?]*)`, 0)
	})
	groupPath += `\/?`

	// Define validator
	if r.userRoles == nil {
		r.userRoles = make(map[*regexp.Regexp][]string)
	}
	r.userRoles[regexp.MustCompile(groupPath)] = roles
}

// BindRole binds an url pattern with user's roles.
func (r *DefaultRouter) BindRole(httpMethod string, urlPattern string, roles ...string) {
	/* Condition validation: Ignore role validation if there is no token store */
	if TokenStore == nil {
		return
	}
}

// GroupRoute groups all same url's prefix.
func (r *DefaultRouter) GroupRoute(s *Server, groupPath string, function func(s *Server)) {
	r.groups = append(r.groups, groupPath)
	function(s)
	r.groups = r.groups[:len(r.groups)-1]
}

// BindRoute binds an url pattern with a handler.
func (r *DefaultRouter) BindRoute(httpMethod string, urlPattern string, handler interface{}) {
	// Format url pattern before assigned to route
	if len(r.groups) > 0 {
		var buffer bytes.Buffer

		for _, path := range r.groups {
			buffer.WriteString(path)
		}

		buffer.WriteString(urlPattern)
		urlPattern = buffer.String()
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
func (r *DefaultRouter) MatchRoute(context *Request, security *Security) (IRoute, map[string]string) {
	for _, route := range r.routes {
		if ok, pathParams := route.MatchURLPattern(context.request.Method, context.Path); !ok {
			continue
		} else {
			// Validate authentication & roles if neccessary
			if TokenStore != nil && security != nil && security.AuthUser != nil {
				for rule, roles := range r.userRoles {

					if rule.MatchString(context.Path) {
						regexRoles := regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(roles, "|")))

						for _, role := range security.AuthUser.UserRoles() {
							if regexRoles.MatchString(role) {
								return route, pathParams
							}
						}
						panic(utils.Status401())
					}
				}
			} else {
				// Simply return
				return route, pathParams
			}
		}
	}
	return nil, nil
}
