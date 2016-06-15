package oauth2

// IRouter describes a router component's characteristic.
type IRouter interface {

	// Group all same url's prefix with user's roles.
	GroupRole(s *Server, groupPath string, roles ...string)

	// Bind an url pattern with user's roles.
	BindRole(httpMethod string, urlPattern string, roles ...string)

	// Group all same url's prefix.
	GroupRoute(s *Server, groupPath string, function func(s *Server))

	// Bind an url pattern with a handler.
	BindRoute(httpMethod string, urlPattern string, handler interface{})

	// Match a route with an url path.
	MatchRoute(context *Request, security *Security) (route IRoute, pathParams map[string]string)
}
