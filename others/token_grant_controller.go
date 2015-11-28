package oauth2

import (
	"net/url"
	"time"

	"github.com/phuc0302/go-cocktail"
)

type TokenGrantController struct {
	config *Config
	store  TokenStore

	values    url.Values
	grantType string
}

// MARK: Struct's constructors
func CreateTokenGrantController(oauth2Config *Config) *TokenGrantController {
	return &TokenGrantController{
		config: oauth2Config,
		store:  oauth2Config.Store,
	}
}

// MARK: Struct's public functions
func (t *TokenGrantController) HandleForm(c *cocktail.Context) {
	info := CreateOAuth2Client(c)
	values := c.Queries

	/* Condition validation: Validate request */
	if values == nil {
		err := common.Status{Status: 400, Title: "Invalid Request", Description: "Invalid Request."}
		c.OutputError(&err)
		return
	}

	/* Condition validation: Validate grant_type */
	grantType := values.Get(GrantType)
	grantsValidation := t.config.grantsValidation

	if len(grantType) == 0 {
		err := common.Status{Status: 400, Title: "Invalid Request", Description: "Missing grant_type parameter."}
		c.OutputError(&err)
		return
	} else if !grantsValidation.MatchString(grantType) {
		err := common.Status{Status: 400, Title: "Invalid Request", Description: "Invalid grant_type parameter."}
		c.OutputError(&err)
		return
	}

	/* Condition validation: Validate client_id */
	if len(info.ClientId) == 0 {
		err := common.Status{Status: 400, Title: "Invalid Request", Description: "Missing client_id parameter."}
		c.OutputError(&err)
		return
	}

	/* Condition validation: Validate client_secret */
	if len(info.ClientSecret) == 0 {
		err := common.Status{Status: 400, Title: "Invalid Request", Description: "Missing client_secret parameter."}
		c.OutputError(&err)
		return
	}

	/* Condition validation: Validate client */
	client := t.store.FindClient(info.ClientId, info.ClientSecret, "")
	if client == nil {
		err := common.Status{Status: 400, Title: "Invalid Request", Description: "Invalid client_id or client_secret parameter."}
		c.OutputError(&err)
		return
	}

	switch grantType {
	case AUTHORIZATION_CODE:
		t.handleAuthorizationCodeGrant(c, values, client)

	case CLIENT_CREDENTIALS:
		t.handleClientCredentialsGrant()

	case PASSWORD:
		t.handlePasswordGrant(c, values, client)

	case REFRESH_TOKEN:
		t.handleRefreshTokenGrant(c, values, client)

	default:
		err := common.Status{Status: 400, Title: "Invalid Request", Description: "Invalid grant_type parameter."}
		c.OutputError(&err)
		return
	}
}

// MARK: Struct's private functions
func (t *TokenGrantController) handleAuthorizationCodeGrant(c *cocktail.Context, values url.Values, client *OAuth2Client) {
	/* Condition validation: Validate code */
	authorizationCode := values.Get(Code)
	if len(authorizationCode) == 0 {
		err := common.Status{Status: 400, Title: "Invalid Request", Description: "Missing code parameter."}
		c.OutputError(&err)
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
func (t *TokenGrantController) handlePasswordGrant(c *cocktail.Context, values url.Values, client *OAuth2Client) {
	user := values.Get(Username)
	pass := values.Get(Password)

	/* Condition validation: Validate username and password parameters */
	if len(user) == 0 || len(pass) == 0 {
		err := common.Status{Status: 400, Title: "Invalid Request", Description: "Missing or invalid username and password parameters."}
		c.OutputError(&err)
		return
	}

	/* Condition validation: Validate user's credentials */
	oauth2User := t.store.FindUser(user, pass)
	if oauth2User == nil {
		err := common.Status{Status: 400, Title: "invalid Grant", Description: "Invalid user's credentials."}
		c.OutputError(&err)
		return
	}
}
func (t *TokenGrantController) handleRefreshTokenGrant(c *cocktail.Context, values url.Values, client *OAuth2Client) {
	/* Condition validation: Validate refresh_token parameter */
	token := t.values.Get(RefreshToken)
	if len(token) == 0 {
		err := common.Status{Status: 400, Title: "Invalid Request", Description: "Missing refresh_token parameter."}
		c.OutputError(&err)
		return
	}

	/* Condition validation: Validate refresh_token */
	refreshToken := t.store.FindRefreshToken(token)

	if refreshToken == nil || refreshToken.ClientId != client.ClientId {
		err := common.Status{Status: 400, Title: "Invalid Request", Description: "Invalid refresh_token parameter."}
		c.OutputError(&err)
		return
	} else if refreshToken.ExpiredTime.Unix() < time.Now().Unix() {
		err := common.Status{Status: 400, Title: "Invalid Grant", Description: "Refresh token is expired."}
		c.OutputError(&err)
		return
	}

	// if (!refreshToken.user && !refreshToken.userId) {
	//   return done(error('server_error', false,
	//     'No user/userId parameter returned from getRefreshToken'));
	// }

	// self.user = refreshToken.user || { id: refreshToken.userId };

	// if (self.model.revokeRefreshToken) {
	//   return self.model.revokeRefreshToken(token, function (err) {
	//     if (err) return done(error('server_error', false, err));
	//     done();
	//   });
	// }
}
func (t *TokenGrantController) handleClientCredentialsGrant() {
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

/**
 * Check the grant type is allowed for this client
 *
 * @param  {Function} done
 * @this   OAuth
 */
func (t *TokenGrantController) CheckGrantTypeAllowed() {
	//	clientId := g.values.Get("client_id")
	//	allowed := g.Model.ValidateGrantType(clientId, g.grantType)

	//	if !allowed {
	//		panic("invalid_client - The grant type is unauthorised for this client_id")
	//	}
}

/**
 * Generate an access token
 *
 * @param  {Function} done
 * @this   OAuth
 */
func (g *TokenGrantController) GenerateAccessToken() string {
	return GenerateToken()
}

/**
 * Save access token with model
 *
 * @param  {Function} done
 * @this   OAuth
 */
func (g *TokenGrantController) SaveAccessToken() {
	// var accessToken = this.accessToken;

	// // Object idicates a reissue
	// if (typeof accessToken === 'object' && accessToken.accessToken) {
	//   this.accessToken = accessToken.accessToken;
	//   return done();
	// }

	// var expires = null;
	// if (this.config.accessTokenLifetime !== null) {
	//   expires = new Date(this.now);
	//   expires.setSeconds(expires.getSeconds() + this.config.accessTokenLifetime);
	// }

	// this.model.saveAccessToken(accessToken, this.client.clientId, expires,
	//     this.user, function (err) {
	//   if (err) return done(error('server_error', false, err));
	//   done();
	// });
}

/**
 * Generate a refresh token
 *
 * @param  {Function} done
 * @this   OAuth
 */
func (g *TokenGrantController) GenerateRefreshToken() string {
	// if (this.config.grants.indexOf('refresh_token') === -1) return done();
	return GenerateToken()
}

/**
 * Save refresh token with model
 *
 * @param  {Function} done
 * @this   OAuth
 */
func (g *TokenGrantController) SaveRefreshToken() {
	// var refreshToken = this.refreshToken;

	// if (!refreshToken) return done();

	// // Object idicates a reissue
	// if (typeof refreshToken === 'object' && refreshToken.refreshToken) {
	//   this.refreshToken = refreshToken.refreshToken;
	//   return done();
	// }

	// var expires = null;
	// if (this.config.refreshTokenLifetime !== null) {
	//   expires = new Date(this.now);
	//   expires.setSeconds(expires.getSeconds() + this.config.refreshTokenLifetime);
	// }

	// this.model.saveRefreshToken(refreshToken, this.client.clientId, expires,
	//     this.user, function (err) {
	//   if (err) return done(error('server_error', false, err));
	//   done();
	// });
}

/**
 * Create an access token and save it with the model
 *
 * @param  {Function} done
 * @this   OAuth
 */
func (g *TokenGrantController) SendResponse() {
	// var response = {
	//   token_type: 'bearer',
	//   access_token: this.accessToken
	// };

	// if (this.config.accessTokenLifetime !== null) {
	//   response.expires_in = this.config.accessTokenLifetime;
	// }

	// if (this.refreshToken) response.refresh_token = this.refreshToken;

	// this.res.set({'Cache-Control': 'no-store', 'Pragma': 'no-cache'});
	// this.res.jsonp(response);

	// if (this.config.continueAfterResponse)
	//   done();
}
