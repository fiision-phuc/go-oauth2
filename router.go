package oauth2

// GroupRole binds user's roles to all url with same prefix.
func (s *Server) GroupRole(groupPath string, roles string) {
	s.router.GroupRole(s, groupPath, roles)
}

// Bind an url pattern with user's roles.
func (s *Server) BindRole(httpMethod string, urlPattern string, roles string) {
	s.router.BindRole(httpMethod, urlPattern, roles)
}

// GroupRoute routes all url with same prefix.
func (s *Server) GroupRoute(urlGroup string, function func(s *Server)) {
	s.router.GroupRoute(s, urlGroup, function)
}

// Copy routes copy request to registered handler.
func (s *Server) Copy(urlPattern string, handler interface{}) {
	s.router.BindRoute(COPY, urlPattern, handler)
}

// Delete routes delete request to registered handler.
func (s *Server) Delete(urlPattern string, handler interface{}) {
	s.router.BindRoute(DELETE, urlPattern, handler)
}

// Get routes get request to registered handler.
func (s *Server) Get(urlPattern string, handler interface{}) {
	s.router.BindRoute(GET, urlPattern, handler)
}

// Head routes head request to registered handler.
func (s *Server) Head(urlPattern string, handler interface{}) {
	s.router.BindRoute(HEAD, urlPattern, handler)
}

// Link routes link request to registered handler.
func (s *Server) Link(urlPattern string, handler interface{}) {
	s.router.BindRoute(LINK, urlPattern, handler)
}

// Options routes options request to registered handler.
func (s *Server) Options(urlPattern string, handler interface{}) {
	s.router.BindRoute(OPTIONS, urlPattern, handler)
}

// Patch routes patch request to registered handler.
func (s *Server) Patch(urlPattern string, handler interface{}) {
	s.router.BindRoute(PATCH, urlPattern, handler)
}

// Post routes post request to registered handler.
func (s *Server) Post(urlPattern string, handler interface{}) {
	s.router.BindRoute(POST, urlPattern, handler)
}

// Purge routes purge request to registered handler.
func (s *Server) Purge(urlPattern string, handler interface{}) {
	s.router.BindRoute(PURGE, urlPattern, handler)
}

// Put routes put request to registered handler.
func (s *Server) Put(urlPattern string, handler interface{}) {
	s.router.BindRoute(PUT, urlPattern, handler)
}

// Unlink routes unlink request to registered handler.
func (s *Server) Unlink(urlPattern string, handler interface{}) {
	s.router.BindRoute(UNLINK, urlPattern, handler)
}
