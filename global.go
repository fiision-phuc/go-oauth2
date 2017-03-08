package oauth2

import (
	"crypto/rsa"
	"fmt"
	"regexp"
	"strings"

	"github.com/phuc0302/go-mongo"
	"github.com/phuc0302/go-oauth2/oauth_key"
	"github.com/phuc0302/go-server"
	"github.com/phuc0302/go-server/util"
)

// Global variables.
var (
	// Global public config's instance.
	Cfg Config

	// Global public token store's instance.
	Store TokenStore

	// Global internal private key.
	privateKey *rsa.PrivateKey
)

// Global regex.
var (
	// Bearer regex.
	bearerFinder = regexp.MustCompile("^(B|b)earer\\s.+$")

	// OAuth2 grant regex.
	grantsValidation *regexp.Regexp
)

// ValidateToken returns a wrapper oauth token for HandleContextFunc.
func ValidateToken() server.Adapter {
	return func(f server.HandleContextFunc) server.HandleContextFunc {
		return func(c *server.RequestContext) {
			tokenString := c.Header["authorization"]

			/* Condition validation: Validate existing of authorization header */
			if isBearer := bearerFinder.MatchString(tokenString); isBearer {
				tokenString = tokenString[7:]
			} else {
				if tokenString = c.QueryParams["access_token"]; len(tokenString) > 0 {
					delete(c.QueryParams, "access_token")
				}
			}

			/* Condition validation: validate token */
			if accessToken := Store.FindAccessToken(tokenString); accessToken != nil && !accessToken.IsExpired() {
				client := Store.FindClientWithID(accessToken.ClientID())
				user := Store.FindUserWithID(accessToken.UserID())

				oauthContext := &OAuthContext{
					Client:      client,
					User:        user,
					AccessToken: accessToken,
				}
				c.SetExtra(oauthKey.OAuthContext, oauthContext)

			} else if username, password, ok := c.BasicAuth(); ok {
				client := Store.FindClientWithCredential(username, password)
				user := Store.FindUserWithClient(username, password)

				if client != nil && user != nil {
					oauthContext := &OAuthContext{
						Client:      client,
						User:        user,
						AccessToken: accessToken,
					}
					c.SetExtra(oauthKey.OAuthContext, oauthContext)
				}
			} else {
				panic(util.Status401())
			}
			f(c)
		}
	}
}

// ValidateToken returns a wrapper oauth token for HandleContextFunc.
func ValidateRoles(roles ...string) server.Adapter {
	return func(f server.HandleContextFunc) server.HandleContextFunc {
		/* Condition validation: validate role input */
		if roles == nil || len(roles) <= 0 {
			return f
		}

		return func(c *server.RequestContext) {
			if oauthContext, ok := c.GetExtra(oauthKey.OAuthContext).(*OAuthContext); ok && Store != nil && oauthContext.User != nil {
				roleValidator := regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(roles, "|")))
				isAuthorized := false

				for _, role := range oauthContext.User.UserRoles() {
					if roleValidator.MatchString(role) {
						isAuthorized = true
						break
					}
				}

				// If user is not authorized, break
				if !isAuthorized {
					panic(util.Status401())
				}
			} else {
				panic(util.Status401())
			}
			f(c)
		}
	}
}

// CreateServer returns a server with custom components.
func CreateServer(tokenStore TokenStore, sandboxMode bool) *server.Server {
	LoadConfig()

	// Register global components
	Store = tokenStore

	// Create server
	server := server.CreateServer(sandboxMode)

	// Setup OAuth2.0
	if Store != nil {
		//	grantAuthorization := new(AuthorizationGrant)
		tokenGrant := new(TokenGrant)

		//	server.Get("/authorize", grantAuthorization.HandleForm)
		server.Get("/token", tokenGrant.HandleForm)
		server.Post("/token", tokenGrant.HandleForm)
	}
	return server
}

// DefaultServer returns a server with build in components.
func DefaultServer(sandboxMode bool) *server.Server {
	mongo.ConnectMongo()
	tokenStore := new(MongoDBStore)
	return CreateServer(tokenStore, sandboxMode)
}
