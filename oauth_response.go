package oauth2

import "github.com/phuc0302/go-oauth2/server"

// OAuthResponse describes a granted response that will be returned to client.
type OAuthResponse struct {
	TokenType    string `json:"token_type,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`

	Roles []string `json:"roles,omitempty"`
}

// createOAuthContext creates new security context.
func createOAuthContext(c *server.RequestContext) *OAuthContext {
	tokenString := c.Header["authorization"]

	/* Condition validation: Validate existing of authorization header */
	if isBearer := bearerFinder.MatchString(tokenString); isBearer {
		tokenString = tokenString[7:]
	} else {
		if tokenString = c.QueryParams["access_token"]; len(tokenString) > 0 {
			delete(c.QueryParams, "access_token")
		}
	}

	/* Condition validation: Validate expiration time */
	if accessToken := Store.FindAccessToken(tokenString); accessToken != nil && !accessToken.IsExpired() {
		client := Store.FindClientWithID(accessToken.ClientID())
		user := Store.FindUserWithID(accessToken.UserID())
		return &OAuthContext{
			Client:      client,
			User:        user,
			AccessToken: accessToken,
		}
	}

	/* Condition validation: If everything is not work out, try to look for basic auth */
	if username, password, ok := c.BasicAuth(); ok {
		client := Store.FindClientWithCredential(username, password)
		user := Store.FindUserWithClient(username, password)

		if client != nil && user != nil {
			return &OAuthContext{
				Client: client,
				User:   user,
			}
		}
	}
	return nil
}
