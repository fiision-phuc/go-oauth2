package oauth2

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/phuc0302/go-oauth2/util"
)

// TokenGrant describes a token grant controller.
type TokenGrant struct {
}

// HandleForm validates authentication form.
func (g *TokenGrant) HandleForm(c *Request, s *Security) {
	security := new(Security)

	if err := g.validateForm(c, security); err == nil {
		g.finalizeToken(c, security)
	} else {
		c.OutputError(err)
	}
}

// validateForm validates general information
func (g *TokenGrant) validateForm(c *Request, s *Security) *util.Status {
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
		return util.Status400WithDescription("Invalid grant_type parameter.")
	}

	/* Condition validation: Validate client_id */
	if len(inputForm.ClientID) == 0 {
		return util.Status400WithDescription("Invalid client_id parameter.")
	}

	/* Condition validation: Validate client_secret */
	if len(inputForm.ClientSecret) == 0 {
		return util.Status400WithDescription("Invalid client_secret parameter.")
	}

	/* Condition validation: Check the store */
	recordClient := TokenStore.FindClientWithCredential(inputForm.ClientID, inputForm.ClientSecret)
	if recordClient == nil {
		return util.Status400WithDescription("Invalid client_id or client_secret parameter.")
	}

	/* Condition validation: Check grant_type for client */
	clientGrantsValidation := regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(recordClient.GrantTypes(), "|")))
	if isGranted := clientGrantsValidation.MatchString(inputForm.GrantType); !isGranted {
		return util.Status400WithDescription("The grant_type is unauthorised for this client_id.")
	}
	s.Client = recordClient

	// Choose authentication flow
	switch inputForm.GrantType {

	case AuthorizationCodeGrant:
		// TODO: Going to do soon
		//		g.handleAuthorizationCodeGrant(c, values, queryClient)
		break

		//	case ImplicitGrant:
		// TODO: Going to do soon
		//		break

	case ClientCredentialsGrant:
		// TODO: Going to do soon
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
	//		err := util.Status400WithDescription("Missing redirect_uri parameter.")
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
	//		err := util.Status400WithDescription("The redirect_uri had not been registered for this client_id.")
	//		c.OutputError(err)
	//		return
	//	}

	/* Condition validation: Validate code */
	authorizationCode := values.Get("code")
	if len(authorizationCode) == 0 {
		err := util.Status400()
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
func (g *TokenGrant) passwordFlow(c *Request, s *Security) *util.Status {
	var passwordForm struct {
		Username string `username`
		Password string `password`
	}
	c.BindForm(&passwordForm)

	/* Condition validation: Validate username and password parameters */
	if len(passwordForm.Username) == 0 || len(passwordForm.Password) == 0 {
		return util.Status400WithDescription("Invalid username or password parameter.")
	}

	/* Condition validation: Validate user's credentials */
	if recordUser := TokenStore.FindUserWithCredential(passwordForm.Username, passwordForm.Password); recordUser != nil {
		s.User = recordUser
		return nil
	}
	return util.Status400WithDescription("Invalid username or password parameter.")

}

// useRefreshTokenFlow handle refresh token flow.
func (g *TokenGrant) refreshTokenFlow(c *Request, s *Security) *util.Status {
	/* Condition validation: Validate refresh_token parameter */
	if queryToken := c.QueryParams["refresh_token"]; len(queryToken) > 0 {

		/* Condition validation: Validate refresh_token */
		refreshToken := TokenStore.FindRefreshToken(queryToken)

		if refreshToken == nil || refreshToken.ClientID() != s.Client.ClientID() {
			return util.Status400WithDescription("Invalid refresh_token parameter.")
		} else if refreshToken.IsExpired() {
			return util.Status400WithDescription("refresh_token is expired.")
		}
		s.User = TokenStore.FindUserWithID(refreshToken.UserID())

		// Delete current access token
		accessToken := TokenStore.FindAccessTokenWithCredential(refreshToken.ClientID(), refreshToken.UserID())
		TokenStore.DeleteAccessToken(accessToken)

		// Delete current refresh token
		TokenStore.DeleteRefreshToken(refreshToken)
		refreshToken = nil

		// Update security context
		s.RefreshToken = nil
		s.AccessToken = nil

		// Delete current refresh token
		return nil
	}
	return util.Status400WithDescription("Invalid refresh_token parameter.")
}

// finalizeToken summary and return result to client.
func (g *TokenGrant) finalizeToken(c *Request, s *Security) {
	now := time.Now()

	// Generate access token if neccessary
	if s.AccessToken == nil {
		accessToken := TokenStore.FindAccessTokenWithCredential(s.Client.ClientID(), s.User.UserID())
		if accessToken != nil && accessToken.IsExpired() {
			TokenStore.DeleteAccessToken(accessToken)
			accessToken = nil
		}

		if accessToken == nil {
			accessToken = TokenStore.CreateAccessToken(
				s.Client.ClientID(),
				s.User.UserID(),
				now,
				now.Add(Cfg.AccessTokenDuration),
			)
		}
		s.AccessToken = accessToken
	}

	// Generate refresh token if neccessary
	if Cfg.AllowRefreshToken && s.RefreshToken == nil {
		refreshToken := TokenStore.FindRefreshTokenWithCredential(s.Client.ClientID(), s.User.UserID())
		if refreshToken != nil && refreshToken.IsExpired() {
			TokenStore.DeleteRefreshToken(refreshToken)
			refreshToken = nil
		}

		if refreshToken == nil {
			refreshToken = TokenStore.CreateRefreshToken(
				s.Client.ClientID(),
				s.User.UserID(),
				now,
				now.Add(Cfg.RefreshTokenDuration),
			)
		}
		s.RefreshToken = refreshToken
	}

	// Generate response token
	tokenResponse := &TokenResponse{
		TokenType:   "Bearer",
		AccessToken: s.AccessToken.Token(),
		ExpiresIn:   s.AccessToken.ExpiredTime().Unix() - time.Now().UTC().Unix(),
		Roles:       s.User.UserRoles(),
	}

	// Only add refresh_token if allowed
	if Cfg.AllowRefreshToken {
		tokenResponse.RefreshToken = s.RefreshToken.Token()
	}
	c.OutputJSON(util.Status200(), tokenResponse)
}
