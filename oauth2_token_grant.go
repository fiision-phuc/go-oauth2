package oauth2

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/phuc0302/go-oauth2/utils"
)

// TokenGrant describes a token grant controller.
type TokenGrant struct {
}

// HandleForm validates authentication form.
func (g *TokenGrant) HandleForm(c *Request) {
	security := &Security{}
	err := g.validateForm(c, security)
	if err != nil {
		c.OutputError(err)
	} else {
		g.finalizeToken(c, security)
	}
}

// validateForm validates general information
func (g *TokenGrant) validateForm(c *Request, s *Security) *utils.Status {
	// Bind
	var inputForm struct {
		GrantType    string `grant_type`
		ClientID     string `client_id`
		ClientSecret string `client_secret`
	}
	c.BindForm(&inputForm)

	// If client_id and client_secret are not include, try to look at the authorization header
	if len(inputForm.ClientID) == 0 && len(inputForm.ClientSecret) == 0 {
		inputForm.ClientID, inputForm.ClientSecret, _ = c.request.BasicAuth()
	}

	/* Condition validation: Validate grant_type */
	if !grantsValidation.MatchString(inputForm.GrantType) {
		return utils.Status400WithDescription("Invalid grant_type parameter.")
	}

	/* Condition validation: Validate client_id */
	if len(inputForm.ClientID) == 0 {
		return utils.Status400WithDescription("Invalid client_id parameter.")
	}

	/* Condition validation: Validate client_secret */
	if len(inputForm.ClientSecret) == 0 {
		return utils.Status400WithDescription("Invalid client_secret parameter.")
	}

	/* Condition validation: Check the store */
	recordClient := tokenStore.FindClientWithCredential(inputForm.ClientID, inputForm.ClientSecret)
	if recordClient == nil {
		return utils.Status400WithDescription("Invalid client_id or client_secret parameter.")
	}

	/* Condition validation: Check grant_type for client */
	clientGrantsValidation := regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(recordClient.GrantTypes(), "|")))
	isGranted := clientGrantsValidation.MatchString(inputForm.GrantType)

	if !isGranted {
		return utils.Status400WithDescription("The grant_type is unauthorised for this client_id.")
	}
	s.AuthClient = recordClient

	// Choose authentication flow
	switch inputForm.GrantType {

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
		return g.passwordFlow(c, s)

	case RefreshTokenGrant:
		return g.refreshTokenFlow(c, s)
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

// passwordFlow implements user's authentication with user's credential.
func (g *TokenGrant) passwordFlow(c *Request, s *Security) *utils.Status {
	var passwordForm struct {
		Username string `username`
		Password string `password`
	}
	c.BindForm(&passwordForm)

	/* Condition validation: Validate username and password parameters */
	if len(passwordForm.Username) == 0 || len(passwordForm.Password) == 0 {
		return utils.Status400WithDescription("Invalid username or password parameter.")
	}

	/* Condition validation: Validate user's credentials */
	recordUser := tokenStore.FindUserWithCredential(passwordForm.Username, passwordForm.Password)
	if recordUser == nil {
		return utils.Status400WithDescription("Invalid username or password parameter.")
	}

	s.AuthUser = recordUser
	return nil
}

// useRefreshTokenFlow handle refresh token flow.
func (g *TokenGrant) refreshTokenFlow(c *Request, s *Security) *utils.Status {
	queryToken := c.QueryParams["refresh_token"]

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
			tokenStore.DeleteAccessToken(accessToken)
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
