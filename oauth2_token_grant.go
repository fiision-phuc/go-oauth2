package oauth2

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/phuc0302/go-oauth2/util"
)

// TokenGrant describes a token grant controller.
type TokenGrant struct {
}

// HandleForm validates authentication form.
func (g *TokenGrant) HandleForm(c *RequestContext, s *OAuthContext) {
	s = new(OAuthContext)

	g.handleForm(c, s)
	g.finalizeToken(c, s)
}

// handleForm validates general information
func (g *TokenGrant) handleForm(c *RequestContext, s *OAuthContext) {
	// If client_id and client_secret are not include, try to look at the authorization header
	if c.QueryParams != nil && len(c.QueryParams["client_id"]) == 0 && len(c.QueryParams["client_secret"]) == 0 {
		c.QueryParams["client_id"], c.QueryParams["client_secret"], _ = c.request.BasicAuth()
	}

	// Bind
	var inputForm struct {
		GrantType    string `field:"grant_type"`
		ClientID     string `field:"client_id" validation:"^\\w+$"`
		ClientSecret string `field:"client_secret" validation:"^\\w+$"`
	}
	err := c.BindForm(&inputForm)

	/* Condition validation: Validate grant_type */
	if !grantsValidation.MatchString(inputForm.GrantType) {
		panic(util.Status400WithDescription(fmt.Sprintf(InvalidParameter, "grant_type")))
	}

	/* Condition validation: Validate binding process */
	if err != nil {
		panic(util.Status400WithDescription(err.Error()))
	}

	/* Condition validation: Check the store */
	recordClient := store.FindClientWithCredential(inputForm.ClientID, inputForm.ClientSecret)
	if recordClient == nil {
		panic(util.Status400WithDescription(fmt.Sprintf(InvalidParameter, "client_id or client_secret")))
	}

	/* Condition validation: Check grant_type for client */
	clientGrantsValidation := regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(recordClient.GrantTypes(), "|")))
	if isGranted := clientGrantsValidation.MatchString(inputForm.GrantType); !isGranted {
		panic(util.Status400WithDescription("The \"grant_type\" is unauthorised for this \"client_id\"."))
	}
	s.Client = recordClient

	// Choose authentication flow
	switch inputForm.GrantType {

	case AuthorizationCodeGrant:
		// TODO: Going to do soon
		g.handleAuthorizationCodeGrant(c, s)
		break

		//	case ImplicitGrant:
		// TODO: Going to do soon
		//		break

	case ClientCredentialsGrant:
		g.handleClientCredentialsGrant(inputForm.ClientID, inputForm.ClientSecret, c, s)
		break

	case PasswordGrant:
		g.passwordFlow(c, s)
		break

	case RefreshTokenGrant:
		g.refreshTokenFlow(c, s)
		break
	}
}

func (t *TokenGrant) handleAuthorizationCodeGrant(c *RequestContext, s *OAuthContext) {
	// Bind
	var inputForm struct {
		Code        string `field:"code" validation:"^\\w+$"`
		RedirectURI string `field:"redirect_uri" validation:"^\\w+$"`
	}

	/* Condition validation: Validate binding process */
	err := c.BindForm(&inputForm)
	if err != nil {
		panic(util.Status400WithDescription(err.Error()))
	}

	/* Condition validation: Check redirect_uri for client */
	isAllow := false
	for _, redirectURI := range s.Client.RedirectURIs() {
		if redirectURI == inputForm.RedirectURI {
			isAllow = true
			break
		}
	}
	if !isAllow {
		panic(util.Status400WithDescription("The \"redirect_uri\" had not been registered for this \"client_id\"."))
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

func (t *TokenGrant) handleClientCredentialsGrant(clientID string, clientSecret string, c *RequestContext, s *OAuthContext) {
	if user := store.FindUserWithClient(clientID, clientSecret); user != nil {
		s.User = user
	} else {
		panic(util.Status400WithDescription(fmt.Sprintf(InvalidParameter, "client_id or client_secret")))
	}
}

// passwordFlow implements user's authentication with user's credential.
func (g *TokenGrant) passwordFlow(c *RequestContext, s *OAuthContext) {
	var passwordForm struct {
		Username string `field:"username" validation:"^\\w+$"`
		Password string `field:"password" validation:"^\\w+$"`
	}
	c.BindForm(&passwordForm)

	/* Condition validation: Validate username and password parameters */
	if len(passwordForm.Username) == 0 || len(passwordForm.Password) == 0 {
		panic(util.Status400WithDescription(fmt.Sprintf(InvalidParameter, "username or password")))
	}

	/* Condition validation: Validate user's credentials */
	if recordUser := store.FindUserWithCredential(passwordForm.Username, passwordForm.Password); recordUser != nil {
		s.User = recordUser
	} else {
		panic(util.Status400WithDescription(fmt.Sprintf(InvalidParameter, "username or password")))
	}
}

// useRefreshTokenFlow handle refresh token flow.
func (g *TokenGrant) refreshTokenFlow(c *RequestContext, s *OAuthContext) {
	/* Condition validation: Validate refresh_token parameter */
	if queryToken := c.QueryParams["refresh_token"]; len(queryToken) > 0 {

		/* Condition validation: Validate refresh_token */
		refreshToken := store.FindRefreshToken(queryToken)

		if refreshToken == nil || refreshToken.ClientID() != s.Client.ClientID() {
			panic(util.Status400WithDescription(fmt.Sprintf(InvalidParameter, "refresh_token")))
		}
		if refreshToken.IsExpired() {
			panic(util.Status400WithDescription("\refresh_token\" is expired."))
		}
		s.User = store.FindUserWithID(refreshToken.UserID())

		// Delete current access token
		accessToken := store.FindAccessTokenWithCredential(refreshToken.ClientID(), refreshToken.UserID())
		store.DeleteAccessToken(accessToken)

		// Delete current refresh token
		store.DeleteRefreshToken(refreshToken)
		refreshToken = nil

		// Update security context
		s.RefreshToken = nil
		s.AccessToken = nil
	} else {
		panic(util.Status400WithDescription(fmt.Sprintf(InvalidParameter, "refresh_token")))
	}
}

// finalizeToken summary and return result to client.
func (g *TokenGrant) finalizeToken(c *RequestContext, s *OAuthContext) {
	now := time.Now()

	// Generate access token if neccessary
	if s.AccessToken == nil {
		accessToken := store.FindAccessTokenWithCredential(s.Client.ClientID(), s.User.UserID())
		if accessToken != nil && accessToken.IsExpired() {
			store.DeleteAccessToken(accessToken)
			accessToken = nil
		}

		if accessToken == nil {
			accessToken = store.CreateAccessToken(
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
		refreshToken := store.FindRefreshTokenWithCredential(s.Client.ClientID(), s.User.UserID())
		if refreshToken != nil && refreshToken.IsExpired() {
			store.DeleteRefreshToken(refreshToken)
			refreshToken = nil
		}

		if refreshToken == nil {
			refreshToken = store.CreateRefreshToken(
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
