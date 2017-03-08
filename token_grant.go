package oauth2

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/phuc0302/go-server"
	"github.com/phuc0302/go-server/string_format"
	"github.com/phuc0302/go-server/util"
)

// TokenGrant describes a token grant controller.
type TokenGrant struct {
}

// HandleForm validates authentication form.
func (g *TokenGrant) HandleForm(c *server.RequestContext) {
	s := new(OAuthContext)

	g.generalValidation(c, s)
	g.finalizeToken(c, s)
}

// handleForm validates general information
func (g *TokenGrant) generalValidation(c *server.RequestContext, s *OAuthContext) {
	// If client_id and client_secret are not include, try to look at the authorization header
	if c.QueryParams != nil && len(c.QueryParams["client_id"]) == 0 && len(c.QueryParams["client_secret"]) == 0 {
		c.QueryParams["client_id"], c.QueryParams["client_secret"], _ = c.BasicAuth()
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
		panic(util.Status400WithDescription(fmt.Sprintf(stringFormat.InvalidParameter, "grant_type")))
	}

	/* Condition validation: Validate binding process */
	if err != nil {
		panic(util.Status400WithDescription(err.Error()))
	}

	/* Condition validation: Check the store */
	recordClient := Store.FindClientWithCredential(inputForm.ClientID, inputForm.ClientSecret)
	if recordClient == nil {
		panic(util.Status400WithDescription(fmt.Sprintf(stringFormat.InvalidParameter, "client_id or client_secret")))
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

func (t *TokenGrant) handleAuthorizationCodeGrant(c *server.RequestContext, s *OAuthContext) {
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

func (t *TokenGrant) handleClientCredentialsGrant(clientID string, clientSecret string, c *server.RequestContext, s *OAuthContext) {
	if user := Store.FindUserWithClient(clientID, clientSecret); user != nil {
		s.User = user
	} else {
		panic(util.Status400WithDescription(fmt.Sprintf(stringFormat.InvalidParameter, "client_id or client_secret")))
	}
}

// passwordFlow implements user's authentication with user's credential.
func (g *TokenGrant) passwordFlow(c *server.RequestContext, s *OAuthContext) {
	var passwordForm struct {
		Username string `field:"username" validation:"^[^\\s]+$"`
		Password string `field:"password" validation:"^[^\\s]{8,32}$"`
	}
	if err := c.BindForm(&passwordForm); err != nil {
		panic(util.Status400WithDescription("Invalid 'username or password' parameter."))
	}

	/* Condition validation: Validate user's credentials */
	if recordUser := Store.FindUserWithCredential(passwordForm.Username, passwordForm.Password); recordUser != nil {
		s.User = recordUser
	} else {
		panic(util.Status400WithDescription(fmt.Sprintf(stringFormat.InvalidParameter, "username or password")))
	}
}

// useRefreshTokenFlow handle refresh token flow.
func (g *TokenGrant) refreshTokenFlow(c *server.RequestContext, s *OAuthContext) {
	/* Condition validation: Validate refresh_token parameter */
	if queryToken := c.QueryParams["refresh_token"]; len(queryToken) > 0 {

		/* Condition validation: Validate refresh_token */
		refreshToken := Store.FindRefreshToken(queryToken)

		if refreshToken == nil || refreshToken.ClientID() != s.Client.ClientID() {
			panic(util.Status400WithDescription(fmt.Sprintf(stringFormat.InvalidParameter, "refresh_token")))
		}
		if refreshToken.IsExpired() {
			panic(util.Status400WithDescription("\refresh_token\" is expired."))
		}
		s.User = Store.FindUserWithID(refreshToken.UserID())

		// Delete current access token
		accessToken := Store.FindAccessTokenWithCredential(refreshToken.ClientID(), refreshToken.UserID())
		Store.DeleteAccessToken(accessToken)

		// Delete current refresh token
		Store.DeleteRefreshToken(refreshToken)
		refreshToken = nil

		// Update security context
		s.RefreshToken = nil
		s.AccessToken = nil
	} else {
		panic(util.Status400WithDescription(fmt.Sprintf(stringFormat.InvalidParameter, "refresh_token")))
	}
}

// finalizeToken summary and return result to client.
func (g *TokenGrant) finalizeToken(c *server.RequestContext, s *OAuthContext) {
	now := time.Now()

	// Generate access token if neccessary
	if s.AccessToken == nil {
		accessToken := Store.FindAccessTokenWithCredential(s.Client.ClientID(), s.User.UserID())
		if accessToken != nil && accessToken.IsExpired() {
			Store.DeleteAccessToken(accessToken)
			accessToken = nil
		}

		if accessToken == nil {
			accessToken = Store.CreateAccessToken(
				s.Client.ClientID(),
				s.User.UserID(),
				now,
				now.Add(cfg.AccessTokenDuration),
			)
		}
		s.AccessToken = accessToken
	}

	// Generate refresh token if neccessary
	if cfg.AllowRefreshToken && s.RefreshToken == nil {
		refreshToken := Store.FindRefreshTokenWithCredential(s.Client.ClientID(), s.User.UserID())
		if refreshToken != nil && refreshToken.IsExpired() {
			Store.DeleteRefreshToken(refreshToken)
			refreshToken = nil
		}

		if refreshToken == nil {
			refreshToken = Store.CreateRefreshToken(
				s.Client.ClientID(),
				s.User.UserID(),
				now,
				now.Add(cfg.RefreshTokenDuration),
			)
		}
		s.RefreshToken = refreshToken
	}

	// Generate response token
	tokenResponse := &OAuthResponse{
		TokenType:   "Bearer",
		AccessToken: s.AccessToken.Token(),
		ExpiresIn:   s.AccessToken.ExpiredTime().Unix() - time.Now().UTC().Unix(),
		Roles:       s.User.UserRoles(),
	}

	// Only add refresh_token if allowed
	if cfg.AllowRefreshToken {
		tokenResponse.RefreshToken = s.RefreshToken.Token()
	}
	c.OutputJSON(util.Status200(), tokenResponse)
}
