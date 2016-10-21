package oauth2

// IFactory describes a factory component's characteristic.
type IFactory interface {

	// Create new request context.
	//	CreateRequestContext(request *http.Request, response http.ResponseWriter) *Request

	// Create new security context.
	//	CreateSecurityContext(requestContext *Request) *Security

	//	// Create new route component.
	//	CreateRoute(urlPattern string) IRoute

	//	// Create new router component.
	//	CreateRouter() ServerRouter

	// Create new store component.
	CreateStore() TokenStore
}
