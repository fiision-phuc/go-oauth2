package oauth2

import (
	"net/url"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/phuc0302/go-oauth2/utils"
)

type GrantToken struct {
	config *Config
	store  Store

	values    url.Values
	grantType string
}

// MARK: Struct's constructors
func CreateGrantToken(config *Config, store Store) *GrantToken {
	return &GrantToken{
		config: config,
		store:  store,
	}
}

// MARK: Struct's public functions
func (g *GrantToken) HandleForm(c *RequestContext) {
	err := g.validateForm(c)
	if err != nil {
		c.OutputError(err)
	} else {
		g.finalizeToken(c)
	}
}

// MARK: Struct's private functions
func (g *GrantToken) validateForm(c *RequestContext) *utils.Status {
	queryClient := createClient(c)

	/* Condition validation: Validate grant_type */
	if !(len(queryClient.GrantType) >= 0 && g.config.grantsValidation.MatchString(queryClient.GrantType)) {
		return utils.Status400WithDescription("Invalid grant_type parameter.")
	}

	/* Condition validation: Validate client_id */
	if len(queryClient.ClientID) == 0 {
		return utils.Status400WithDescription("Invalid client_id parameter.")
	}

	/* Condition validation: Validate client_secret */
	if len(queryClient.ClientSecret) == 0 {
		return utils.Status400WithDescription("Invalid client_secret parameter.")
	}

	/* Condition validation: Check the store */
	recordClient := g.store.FindClientWithCredential(queryClient.ClientID, queryClient.ClientSecret)
	if recordClient == nil {
		return utils.Status400WithDescription("Invalid client_id or client_secret parameter.")
	}

	/* Condition validation: Check grant_type for client */
	isGranted := false
	for _, grantType := range recordClient.GrantTypes {
		if grantType == queryClient.GrantType {
			isGranted = true
			break
		}
	}
	if !isGranted {
		return utils.Status400WithDescription("The grant_type is unauthorised for this client_id.")
	}
	c.AuthClient = recordClient

	// Choose authentication flow
	switch queryClient.GrantType {

	case AuthorizationCodeGrant:
		// FIX FIX FIX: Going to do soon
		//		g.handleAuthorizationCodeGrant(c, values, queryClient)
		break

	case ImplicitGrant:
		// FIX FIX FIX: Going to do soon
		break

	case ClientCredentialsGrant:
		// FIX FIX FIX: Going to do soon
		//		g.handleClientCredentialsGrant()
		break

	case PasswordGrant:
		return g.usePasswordFlow(c)

	case RefreshTokenGrant:
		return g.useRefreshTokenFlow(c)
		break
	}
	return nil
}

func (t *GrantToken) handleAuthorizationCodeGrant(c *RequestContext, values url.Values, client *Client) {
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

	t.store.FindAuthorizationCode(authorizationCode)
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

func (t *GrantToken) handleClientCredentialsGrant() {
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

func (g *GrantToken) usePasswordFlow(c *RequestContext) *utils.Status {
	queryUser := User{}
	c.BindForm(&queryUser)

	/* Condition validation: Validate username and password parameters */
	if len(queryUser.Username) == 0 || len(queryUser.Password) == 0 {
		return utils.Status400WithDescription("Invalid username or password parameter.")
	}

	/* Condition validation: Validate user's credentials */
	recordUser := g.store.FindUserWithCredential(queryUser.Username, queryUser.Password)
	if recordUser == nil {
		return utils.Status400WithDescription("Invalid username or password parameter.")
	}

	c.AuthUser = recordUser
	return nil
}

func (g *GrantToken) useRefreshTokenFlow(c *RequestContext) *utils.Status {
	/* Condition validation: Validate refresh_token parameter */
	queryToken := c.Queries.Get("refresh_token")
	if len(queryToken) == 0 {
		return utils.Status400WithDescription("Invalid refresh_token parameter.")
	}

	/* Condition validation: Validate refresh_token */
	recordToken := g.store.FindRefreshToken(queryToken)
	if recordToken == nil || recordToken.ClientID != c.AuthClient.ClientID {
		return utils.Status400WithDescription("Invalid refresh_token parameter.")

	} else if recordToken.ExpiredTime.Unix() < time.Now().Unix() {
		return utils.Status400WithDescription("refresh_token is expired.")
	}

	c.AuthAccessToken = g.store.FindAccessTokenWithCredential(recordToken.ClientID, recordToken.UserID)
	c.AuthUser = g.store.FindUserWithID(recordToken.UserID)
	c.AuthRefreshToken = recordToken
	now := time.Now()

	// Update access token
	c.AuthAccessToken.Token = utils.GenerateToken()
	c.AuthAccessToken.CreatedTime = now
	c.AuthAccessToken.ExpiredTime = now.Add(g.config.DurationAccessToken)
	g.store.SaveAccessToken(c.AuthAccessToken)

	// Update refresh token
	c.AuthRefreshToken.Token = utils.GenerateToken()
	c.AuthRefreshToken.CreatedTime = now
	c.AuthRefreshToken.ExpiredTime = now.Add(g.config.DurationRefreshToken)
	g.store.SaveRefreshToken(c.AuthRefreshToken)

	return nil
}

func (g *GrantToken) finalizeToken(c *RequestContext) {
	now := time.Now()

	// Generate access token if neccessary
	if c.AuthAccessToken == nil {
		accessToken := g.store.FindAccessTokenWithCredential(c.AuthClient.ClientID, c.AuthUser.UserID)
		if accessToken != nil && accessToken.isExpired() {
			g.store.DeleteAccessToken(accessToken)
			accessToken = nil
		}

		if accessToken == nil {
			accessToken = &Token{
				TokenID:     bson.NewObjectId(),
				UserID:      c.AuthUser.UserID,
				ClientID:    c.AuthClient.ClientID,
				Token:       utils.GenerateToken(),
				CreatedTime: now,
			}
			accessToken.ExpiredTime = now.Add(g.config.DurationAccessToken)
			g.store.SaveAccessToken(accessToken)
		}
		c.AuthAccessToken = accessToken
	}

	// Generate refresh token if neccessary
	if g.config.allowRefreshToken && c.AuthRefreshToken == nil {
		refreshToken := g.store.FindRefreshTokenWithCredential(c.AuthClient.ClientID, c.AuthUser.UserID)
		if refreshToken != nil && refreshToken.isExpired() {
			g.store.DeleteRefreshToken(refreshToken)
			refreshToken = nil
		}

		if refreshToken == nil {
			refreshToken = &Token{
				TokenID:     bson.NewObjectId(),
				UserID:      c.AuthUser.UserID,
				ClientID:    c.AuthClient.ClientID,
				Token:       utils.GenerateToken(),
				CreatedTime: now,
			}
			refreshToken.ExpiredTime = now.Add(g.config.DurationRefreshToken)
			g.store.SaveRefreshToken(refreshToken)
		}
		c.AuthRefreshToken = refreshToken
	}

	tokenResponse := &TokenResponse{
		TokenType:   "Bearer",
		AccessToken: c.AuthAccessToken.Token,
		ExpiresIn:   c.AuthAccessToken.ExpiredTime.Unix() - time.Now().Unix(),
	}

	if g.config.allowRefreshToken {
		tokenResponse.RefreshToken = c.AuthRefreshToken.Token
	}
	c.OutputJSON(utils.Status200(), tokenResponse)
}
