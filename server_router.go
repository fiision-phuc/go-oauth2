package oauth2

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

// ContextHandler is an alias for function to handle request context & security context.
type ContextHandler func(requestContext *Request, securityContext *Security)

// Router describes a router's implementation.
type Router struct {
	groups []string

	mux   *mux.Router
	roles map[*regexp.Regexp]*regexp.Regexp
}

// CreateRouter creates a router object.
func CreateRouter() *Router {
	router := Router{
		mux: new(mux.Router),
	}
	return &router
}

// GroupRoles groups all same url's prefix with user's roles.
func (r *Router) GroupRoles(groupPath string, roles ...string) {
	/* Condition validation: ignore role validation if there is no token store */
	if TokenStore == nil {
		return
	}

	// Format pattern
	groupPath = pathFinder.ReplaceAllStringFunc(groupPath, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:len(m)-1])
	})
	groupPath = globsFinder.ReplaceAllStringFunc(groupPath, func(m string) string {
		return fmt.Sprintf(`(?P<_%d>[^#?]*)`, 0)
	})

	// Define validator
	if r.roles == nil {
		r.roles = make(map[*regexp.Regexp]*regexp.Regexp)
	}
	r.roles[regexp.MustCompile(groupPath)] = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(roles, "|")))
}

// BindRoles binds an url pattern with user's roles.
func (r *Router) BindRoles(httpMethod string, urlPattern string, roles ...string) {
	/* Condition validation: ignore role validation if there is no token store */
	if TokenStore == nil {
		return
	}
}

// GroupRoute groups all same url's prefix.
func (r *Router) GroupRoute(s *Server, groupPath string, function func(s *Server)) {
	//	subMux := r.mux.PathPrefix(groupPath).Subrouter()
	r.groups = append(r.groups, groupPath)
	function(s)
	r.groups = r.groups[:len(r.groups)-1]
}

// BindRoute binds an url pattern with a handler.
func (r *Router) BindRoute(httpMethod string, urlPattern string, handler ContextHandler) {
	// Format url pattern before assigned to route
	if len(r.groups) > 0 {
		var buffer bytes.Buffer
		for _, path := range r.groups {
			buffer.WriteString(path)
		}

		if len(urlPattern) > 0 {
			buffer.WriteString(urlPattern)
		}
		urlPattern = buffer.String()
	}
	logrus.Infof("%s -> %s", httpMethod, urlPattern)

	// Register route
	r.mux.HandleFunc(urlPattern, func(response http.ResponseWriter, request *http.Request) {
		context := objectFactory.CreateRequestContext(request, response)
		security := objectFactory.CreateSecurityContext(context)

		// Validate user's authorized first
		isAuthorized := true
		for rule, roles := range r.roles {
			if rule.MatchString(context.Path) {
				isAuthorized = false

				if TokenStore != nil && security != nil && security.User != nil {
					for _, role := range security.User.UserRoles() {
						if roles.MatchString(role) {
							isAuthorized = true
							break // If user is authorized, break the loop
						}
					}
				} else {
					//					return nil, nil
				}
			}
		}
		fmt.Println(isAuthorized)
		handler(context, security)
	}).Methods(httpMethod)
}
