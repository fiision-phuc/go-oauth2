package oauth2

import (
	"net/http"
)

// IFactory descripts a factory component's characteristic.
type IFactory interface {

	// Create new request context.
	CreateRequestContext(request *http.Request, response http.ResponseWriter) *Request

	// Create new security context.
	CreateSecurityContext() *Security

	// Create new route component.
	CreateRoute() IRoute

	// Create new router component.
	CreateRouter() IRouter
}
