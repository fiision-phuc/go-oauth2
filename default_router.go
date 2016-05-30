package oauth2

import "github.com/phuc0302/go-oauth2/i"

// defaultRouter object description.
type defaultRouter struct {
	routes  []i.IRoute
	groups  []string
	factory i.IFactory
}

// MARK: IRouter's members
func (r *defaultRouter) BindFactory(factory i.IFactory) {
	r.factory = factory
}

//func (r *defaultRouter) GroupRole(s *Server, groupPath string, roles string) {
//}

func (r *defaultRouter) BindRole(httpMethod string, urlPattern string, roles string) {
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

//func (r *defaultRouter) GroupRoute(s *Server, groupPath string, function func(s *Server)) {
//	r.groups = append(r.groups, groupPath)
//	function(s)
//	r.groups = r.groups[:len(r.groups)-1]
//}

func (r *defaultRouter) BindRoute(httpMethod string, urlPattern string, handler interface{}) {
	//	defer RecoveryInternal(s.logger)

	//	// Format url pattern before assigned to route
	//	if len(r.groups) > 0 {
	//		var groupPattern bytes.Buffer

	//		for _, g := range r.groups {
	//			groupPattern.WriteString(utils.FormatPath(g))
	//		}

	//		if len(urlPattern) > 0 {
	//			groupPattern.WriteString(utils.FormatPath(urlPattern))
	//		}
	//		urlPattern = groupPattern.String()
	//	} else {
	//		urlPattern = utils.FormatPath(urlPattern)
	//	}

	// Look for existing one before create new
	for _, route := range r.routes {
		if route.URLPattern() == urlPattern {
			route.BindHandler(httpMethod, handler)
			//			s.logger.Printf("%-6s -> %s\n", method, urlPattern)
			return
		}
	}

	//	// Create new route
	//	newRoute := CreateDefaultRoute(urlPattern)
	//	newRoute.BindHandler(httpMethod, handler)

	// Add to collection
	//	r.routes = append(r.routes, newRoute)
	//	s.logger.Printf("%-6s -> %s\n", method, urlPattern)
}
