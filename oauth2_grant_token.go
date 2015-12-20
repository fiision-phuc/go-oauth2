package oauth2

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/phuc0302/go-oauth2/utils"
)

type TokenGrant struct {
	config *Config
	store  TokenStore

	values    url.Values
	grantType string
}

// MARK: Struct's constructors
func CreateTokenGrant(config *Config, store TokenStore) *TokenGrant {
	return &TokenGrant{
		config: config,
		store:  store,
	}
}

// MARK: Struct's public functions
func (g *TokenGrant) HandleGet(c *RequestContext) {
	prefix := ""
	for path, folder := range g.config.StaticFolders {
		if strings.HasSuffix(path, "/templates") {
			prefix = folder
			break
		}
	}

	c.OutputHTML(fmt.Sprintf("%s/login.html", prefix), nil)
}
func (g *TokenGrant) HandleForm(c *RequestContext) {
	securityContext := &SecurityContext{}
	err := g.validateForm(c, securityContext)
	if err != nil {
		c.OutputError(err)
	} else {
		g.finalizeToken(c, securityContext)
	}
}

// MARK: Struct's private functions
func (g *TokenGrant) validateForm(c *RequestContext, s *SecurityContext) *utils.Status {
	grantType := c.Queries.Get("grant_type")
	clientID := c.Queries.Get("client_id")
	clientSecret := c.Queries.Get("client_secret")

	/* Condition validation: Validate grant_type */
	if !(len(grantType) >= 0 && g.config.grantsValidation.MatchString(grantType)) {
		return utils.Status400WithDescription("Invalid grant_type parameter.")
	}

	// If client_id and client_secret are not include, try to look at the authorization header
	if len(clientID) == 0 && len(clientSecret) == 0 {
		clientID, clientSecret, _ = c.BasicAuth()
	}

	/* Condition validation: Validate client_id */
	if len(clientID) == 0 {
		return utils.Status400WithDescription("Invalid client_id parameter.")
	}

	/* Condition validation: Validate client_secret */
	if len(clientSecret) == 0 {
		return utils.Status400WithDescription("Invalid client_secret parameter.")
	}

	/* Condition validation: Check the store */
	recordClient := g.store.FindClientWithCredential(clientID, clientSecret)
	if recordClient == nil {
		return utils.Status400WithDescription("Invalid client_id or client_secret parameter.")
	}

	/* Condition validation: Check grant_type for client */
	isGranted := false
	for _, recordGrant := range recordClient.GetGrantTypes() {
		if recordGrant == grantType {
			isGranted = true
			break
		}
	}
	if !isGranted {
		return utils.Status400WithDescription("The grant_type is unauthorised for this client_id.")
	}
	s.AuthClient = recordClient

	// Choose authentication flow
	switch grantType {

	case AuthorizationCodeGrant:
		// FIX FIX FIX: Going to do soon
		//		g.handleAuthorizationCodeGrant(c, values, queryClient)
		break

		//	case ImplicitGrant:
		//		// FIX FIX FIX: Going to do soon
		//		break

	case ClientCredentialsGrant:
		// FIX FIX FIX: Going to do soon
		//		g.handleClientCredentialsGrant()
		break

	case PasswordGrant:
		return g.usePasswordFlow(c, s)

	case RefreshTokenGrant:
		return g.useRefreshTokenFlow(c, s)
		break
	}
	return nil
}

func (t *TokenGrant) handleAuthorizationCodeGrant(c *RequestContext, values url.Values, client *AuthClientDefault) {
	//	/* Condition validation: Validate redirect_uri */
	//	if len(queryClient.RedirectURI) == 0 {
	//		err := utils.Status400WithDescription("Missing redirect_uri parameter.")
	//		c.OutputError(err)
	//		return
	//	}

	//	/* Condition validation: Check redirect_uri for client */
	//	isAllowRedirectURI := false
	//	for _, redirectURI := range recordClient.RedirectURIs {
	//		if redirectURI == queryClient.RedirectURI {
	//			isAllowRedirectURI = true
	//			break
	//		}
	//	}
	//	if !isAllowRedirectURI {
	//		err := utils.Status400WithDescription("The redirect_uri had not been registered for this client_id.")
	//		c.OutputError(err)
	//		return
	//	}

	/* Condition validation: Validate code */
	authorizationCode := values.Get("code")
	if len(authorizationCode) == 0 {
		err := utils.Status400()
		err.Description = "Missing code parameter."
		c.OutputError(err)
		return
	}

	//	t.store.FindAuthorizationCode(authorizationCode)
	// this.model.getAuthCode(code, function (err, authCode) {

	//   if (!authCode || authCode.clientId !== self.client.clientId) {
	//     return done(error('invalid_grant', 'Invalid code'));
	//   } else if (authCode.expires < self.now) {
	//     return done(error('invalid_grant', 'Code has expired'));
	//   }

	//   self.user = authCode.user || { id: authCode.userId };
	//   if (!self.user.id) {
	//     return done(error('server_error', false,
	//       'No user/userId parameter returned from getauthCode'));
	//   }

	//   done();
	// });
}

func (t *TokenGrant) handleClientCredentialsGrant() {
	// // Client credentials
	// var clientId = this.client.clientId,
	//   clientSecret = this.client.clientSecret;

	// if (!clientId || !clientSecret) {
	//   return done(error('invalid_client',
	//     'Missing parameters. "client_id" and "client_secret" are required'));
	// }

	// var self = this;
	// return this.model.getUserFromClient(clientId, clientSecret,
	//     function (err, user) {
	//   if (err) return done(error('server_error', false, err));
	//   if (!user) {
	//     return done(error('invalid_grant', 'Client credentials are invalid'));
	//   }

	//   self.user = user;
	//   done();
	// });
}

// usePasswordFlow handle password flow.
func (g *TokenGrant) usePasswordFlow(c *RequestContext, s *SecurityContext) *utils.Status {
	username := c.Queries.Get("username")
	password := c.Queries.Get("password")

	/* Condition validation: Validate username and password parameters */
	if len(username) == 0 || len(password) == 0 {
		return utils.Status400WithDescription("Invalid username or password parameter.")
	}

	/* Condition validation: Validate user's credentials */
	recordUser := g.store.FindUserWithCredential(username, password)
	if recordUser == nil {
		return utils.Status400WithDescription("Invalid username or password parameter.")
	}

	s.AuthUser = recordUser
	return nil
}

// useRefreshTokenFlow handle refresh token flow.
func (g *TokenGrant) useRefreshTokenFlow(c *RequestContext, s *SecurityContext) *utils.Status {
	queryToken := c.Queries.Get("refresh_token")

	/* Condition validation: Validate refresh_token parameter */
	if len(queryToken) == 0 {
		return utils.Status400WithDescription("Invalid refresh_token parameter.")
	}

	/* Condition validation: Validate refresh_token */
	recordToken := g.store.FindRefreshToken(queryToken)
	if recordToken == nil || recordToken.GetClientID() != s.AuthClient.GetClientID() {
		return utils.Status400WithDescription("Invalid refresh_token parameter.")

	} else if recordToken.IsExpired() {
		return utils.Status400WithDescription("refresh_token is expired.")
	}

	s.AuthUser = g.store.FindUserWithID(recordToken.GetUserID())
	s.AuthRefreshToken = recordToken
	return nil
}

// finalizeToken summary and return result to client.
func (g *TokenGrant) finalizeToken(c *RequestContext, s *SecurityContext) {
	now := time.Now()

	// Generate access token if neccessary
	if s.AuthAccessToken == nil {
		accessToken := g.store.FindAccessTokenWithCredential(s.AuthClient.GetClientID(), s.AuthUser.GetUserID())
		if accessToken != nil && accessToken.IsExpired() {
			//			g.store.DeleteAccessToken(accessToken) // Note: Let the cron delete, it should be safer.
			accessToken = nil
		}

		if accessToken == nil {
			accessToken = g.store.CreateAccessToken(
				s.AuthClient.GetClientID(),
				s.AuthUser.GetUserID(),
				now,
				now.Add(g.config.DurationAccessToken),
			)
		}
		s.AuthAccessToken = accessToken
	}

	// Generate refresh token if neccessary
	if g.config.allowRefreshToken && s.AuthRefreshToken == nil {
		refreshToken := g.store.FindRefreshTokenWithCredential(s.AuthClient.GetClientID(), s.AuthUser.GetUserID())
		if refreshToken != nil && refreshToken.IsExpired() {
			//			g.store.DeleteRefreshToken(refreshToken) // Note: Let the cron delete, it should be safer.
			refreshToken = nil
		}

		if refreshToken == nil {
			refreshToken = g.store.CreateRefreshToken(
				s.AuthClient.GetClientID(),
				s.AuthUser.GetUserID(),
				now,
				now.Add(g.config.DurationRefreshToken),
			)
		}
		s.AuthRefreshToken = refreshToken
	}

	// Generate response token
	tokenResponse := &TokenResponse{
		TokenType:   "Bearer",
		AccessToken: s.AuthAccessToken.GetToken(),
		ExpiresIn:   s.AuthAccessToken.GetExpiredTime().Unix() - time.Now().Unix(),
	}

	// Only add refresh_token if allowed
	if g.config.allowRefreshToken {
		tokenResponse.RefreshToken = s.AuthRefreshToken.GetToken()
	}
	c.OutputJSON(utils.Status200(), tokenResponse)
}
