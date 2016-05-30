package oauth2

import (
	"net/http"

	"github.com/phuc0302/go-oauth2/context"
)

// IFactory descripts a factory component's characteristic.
type IFactory interface {

	// Create new request context.
	CreateRequestContext(request *http.Request, response http.ResponseWriter) *context.Request

	// Create new security context.
	CreateSecurityContext() *context.Security

	// Create new route component.
	CreateRoute() IRoute

	// Create new router component.
	CreateRouter() IRouter
}
