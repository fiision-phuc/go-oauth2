package oauth2

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/phuc0302/go-oauth2/oauth_key"
	"github.com/phuc0302/go-server"
	"github.com/phuc0302/go-server/util"
)

// OAuthContext describes a user's oauth scope.
type OAuthContext struct {

	// Registered user. Always available.
	User User
	// Registered client. Always available.
	Client Client
	// Access token that had been given to user. Always available.
	AccessToken Token
	// Refresh token that had been given to user. Might not be available all the time.
	RefreshToken Token
}

// OAuthResponse describes a granted response that will be returned to client.
type OAuthResponse struct {
	TokenType    string `json:"token_type,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`

	Roles []string `json:"roles,omitempty"`
}

// ValidateToken returns a wrapper oauth token validation func before HandleContextFunc.
//
// @return
// - func {server.Adapter} (a wrapper func around developer's server.HandleContextFunc)
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
				c.SetExtra(oauthKey.Context, oauthContext)

			} else if username, password, ok := c.BasicAuth(); ok {
				client := Store.FindClientWithCredential(username, password)
				user := Store.FindUserWithClient(username, password)

				if client != nil && user != nil {
					oauthContext := &OAuthContext{
						Client:      client,
						User:        user,
						AccessToken: accessToken,
					}
					c.SetExtra(oauthKey.Context, oauthContext)
				}
			} else {
				panic(util.Status401())
			}
			f(c)
		}
	}
}

// ValidateRoles returns a wrapper user's roles validation func before HandleContextFunc.
//
// @param
// - roles {[]string} (a list of acceptable users' roles)
//
// @return
// - func {server.Adapter} (a wrapper func around developer's server.HandleContextFunc)
func ValidateRoles(roles ...string) server.Adapter {
	return func(f server.HandleContextFunc) server.HandleContextFunc {
		/* Condition validation: validate role input */
		if roles == nil || len(roles) <= 0 {
			return f
		}

		return func(c *server.RequestContext) {
			if oauthContext, ok := c.GetExtra(oauthKey.Context).(*OAuthContext); ok && Store != nil && oauthContext.User != nil {
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
