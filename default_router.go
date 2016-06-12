package oauth2

import (
	"bytes"
	"regexp"

	"github.com/Sirupsen/logrus"
)

// DefaultRouter descripts a default router component implementation.
type DefaultRouter struct {
	routes    []IRoute
	groups    []string
	userRoles map[*regexp.Regexp][]string
}

// GroupRole groups all same url's prefix with user's roles.
func (r *DefaultRouter) GroupRole(s *Server, groupPath string, roles string) {
}

// BindRole binds an url pattern with user's roles.
func (r *DefaultRouter) BindRole(httpMethod string, urlPattern string, roles string) {
	//	/* Condition validation: Ignore role validation if there is no token store */
	//	if s.tokenStore == nil {
	//		return
	//	}

	//	pattern = utils.FormatPath(pattern)

	//	pattern = pathParamRegex.ReplaceAllStringFunc(pattern, func(m string) string {
	//		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	//	})
	//	pattern = globsRegex.ReplaceAllStringFunc(pattern, func(m string) string {
	//		return fmt.Sprintf(`(?P<_%d>[^#?]*)`, 0)
	//	})
	//	pattern += `\/?`

	//	userRoles := strings.Split(roles, ",")

	//	if s.userRoles == nil {
	//		s.userRoles = make(map[*regexp.Regexp][]string, 1)
	//	}
	//	s.userRoles[regexp.MustCompile(pattern)] = userRoles
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
func (r *DefaultRouter) MatchRoute(httpMethod string, path string) (IRoute, map[string]string) {
	for _, route := range r.routes {
		ok, pathParams := route.MatchURLPattern(httpMethod, path)
		if !ok {
			continue
		}

		//		// Validate authentication & roles if neccessary
		//		if tokenStore != nil {
		//			securityContext := objectFactory.CreateSecurityContext(context)

		//			for rule, roles := range r.userRoles {
		//				if rule.MatchString(context.URLPath) {
		//					regexRoles := regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(roles, "|")))

		//					if securityContext != nil && securityContext.AuthUser != nil {
		//						for _, role := range securityContext.AuthUser.UserRoles() {
		//							if regexRoles.MatchString(role) {
		//								route.InvokeHandler(context, securityContext)
		//								return
		//							}
		//						}
		//					}
		//					context.OutputError(utils.Status401())
		//					return
		//				}
		//			}
		//		}

		return route, pathParams
	}
	return nil, nil
}
