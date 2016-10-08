package oauth2

import (
	"net/http"
)

// IFactory describes a factory component's characteristic.
type IFactory interface {

	// Create new request context.
	CreateRequestContext(request *http.Request, response http.ResponseWriter) *Request

	// Create new security context.
	CreateSecurityContext(requestContext *Request) *Security

	// Create new route component.
	CreateRoute(urlPattern string) *Route

	// Create new router component.
	CreateRouter() *Router

	// Create new store component.
	CreateStore() IStore
}
