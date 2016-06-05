package oauth2

import (
	"net/url"
	"time"

	"github.com/phuc0302/go-oauth2/utils"
)

type TokenGrant struct {
	values    url.Values
	grantType string
}

// MARK: Struct's constructors
func CreateTokenGrant(config *config, store IStore) *TokenGrant {
	return &TokenGrant{
	//		config: config,
	//		store:  store,
	}
}

// MARK: Struct's public functions
func (g *TokenGrant) HandleForm(c *Request) {
	security := &Security{}
	err := g.validateForm(c, security)
	if err != nil {
		c.OutputError(err)
	} else {
		g.finalizeToken(c, security)
	}
}

// MARK: Struct's private functions
func (g *TokenGrant) validateForm(c *Request, s *Security) *utils.Status {
	grantType := c.QueryParams.Get("grant_type")
	clientID := c.QueryParams.Get("client_id")
	clientSecret := c.QueryParams.Get("client_secret")

	/* Condition validation: Validate grant_type */
	if !(len(grantType) >= 0 && grantsValidation.MatchString(grantType)) {
		return utils.Status400WithDescription("Invalid grant_type parameter.")
	}

	// If client_id and client_secret are not include, try to look at the authorization header
	if len(clientID) == 0 && len(clientSecret) == 0 {
		clientID, clientSecret, _ = c.request.BasicAuth()
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
	recordClient := tokenStore.FindClientWithCredential(clientID, clientSecret)
	if recordClient == nil {
		return utils.Status400WithDescription("Invalid client_id or client_secret parameter.")
	}

	/* Condition validation: Check grant_type for client */
	isGranted := false
	for _, recordGrant := range recordClient.GrantTypes() {
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
	}
	return nil
}

func (t *TokenGrant) handleAuthorizationCodeGrant(c *Request, values url.Values, client IClient) {
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
func (g *TokenGrant) usePasswordFlow(c *Request, s *Security) *utils.Status {
	username := c.QueryParams.Get("username")
	password := c.QueryParams.Get("password")

	/* Condition validation: Validate username and password parameters */
	if len(username) == 0 || len(password) == 0 {
		return utils.Status400WithDescription("Invalid username or password parameter.")
	}

	/* Condition validation: Validate user's credentials */
	recordUser := tokenStore.FindUserWithCredential(username, password)
	if recordUser == nil {
		return utils.Status400WithDescription("Invalid username or password parameter.")
	}

	//	s.AuthUser = recordUser
	return nil
}

// useRefreshTokenFlow handle refresh token flow.
func (g *TokenGrant) useRefreshTokenFlow(c *Request, s *Security) *utils.Status {
	queryToken := c.QueryParams.Get("refresh_token")

	/* Condition validation: Validate refresh_token parameter */
	if len(queryToken) == 0 {
		return utils.Status400WithDescription("Invalid refresh_token parameter.")
	}

	/* Condition validation: Validate refresh_token */
	recordToken := tokenStore.FindRefreshToken(queryToken)
	if recordToken == nil || recordToken.ClientID() != s.AuthClient.ClientID() {
		return utils.Status400WithDescription("Invalid refresh_token parameter.")

	} else if recordToken.IsExpired() {
		return utils.Status400WithDescription("refresh_token is expired.")
	}

	s.AuthUser = tokenStore.FindUserWithID(recordToken.UserID())
	s.AuthRefreshToken = recordToken

	// Delete current access token
	accessToken := tokenStore.FindAccessTokenWithCredential(recordToken.ClientID(), recordToken.UserID())
	tokenStore.DeleteAccessToken(accessToken)
	return nil
}

// finalizeToken summary and return result to client.
func (g *TokenGrant) finalizeToken(c *Request, s *Security) {
	now := time.Now()

	// Generate access token if neccessary
	if s.AuthAccessToken == nil {
		accessToken := tokenStore.FindAccessTokenWithCredential(s.AuthClient.ClientID(), s.AuthUser.UserID())
		if accessToken != nil && accessToken.IsExpired() {
			tokenStore.DeleteAccessToken(accessToken) // Note: Let the cron delete, it should be safer.
			accessToken = nil
		}

		if accessToken == nil {
			accessToken = tokenStore.CreateAccessToken(
				s.AuthClient.ClientID(),
				s.AuthUser.UserID(),
				now,
				now.Add(cfg.AccessTokenDuration),
			)
		}
		s.AuthAccessToken = accessToken
	}

	// Generate refresh token if neccessary
	if cfg.AllowRefreshToken && s.AuthRefreshToken == nil {
		refreshToken := tokenStore.FindRefreshTokenWithCredential(s.AuthClient.ClientID(), s.AuthUser.UserID())
		if refreshToken != nil && refreshToken.IsExpired() {
			tokenStore.DeleteRefreshToken(refreshToken) // Note: Let the cron delete, it should be safer.
			refreshToken = nil
		}

		if refreshToken == nil {
			refreshToken = tokenStore.CreateRefreshToken(
				s.AuthClient.ClientID(),
				s.AuthUser.UserID(),
				now,
				now.Add(cfg.RefreshTokenDuration),
			)
		}
		s.AuthRefreshToken = refreshToken
	}

	// Generate response token
	tokenResponse := &TokenResponse{
		TokenType:   "Bearer",
		AccessToken: s.AuthAccessToken.Token(),
		ExpiresIn:   s.AuthAccessToken.ExpiredTime().Unix() - time.Now().Unix(),
		Roles:       s.AuthUser.UserRoles(),
	}

	// Only add refresh_token if allowed
	if cfg.AllowRefreshToken {
		tokenResponse.RefreshToken = s.AuthRefreshToken.Token()
	}
	c.OutputJSON(utils.Status200(), tokenResponse)
}
