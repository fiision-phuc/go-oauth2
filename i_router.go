package oauth2

// IRouter descripts a router component's characteristic.
type IRouter interface {

	// Bind a factory to a router.
	BindFactory(factory IFactory)

	//	// Group all same url's prefix with user's roles.
	//	GroupRole(s *Server, groupPath string, roles string)

	// Bind an url pattern with user's roles.
	BindRole(httpMethod string, urlPattern string, roles string)

	//	// Group all same url's prefix.
	//	GroupRoute(s *Server, groupPath string, function func(s *Server))

	// Bind an url pattern with a handler.
	BindRoute(httpMethod string, urlPattern string, handler interface{})
}
