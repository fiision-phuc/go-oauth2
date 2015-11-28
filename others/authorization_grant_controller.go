package oauth2

import (
	"fmt"
	"time"

	"github.com/phuc0302/go-cocktail"
)

type AuthorizationGrantController struct {
	config *Config
	store  TokenStore
}

// MARK: Struct's constructors
func CreateAuthorizationGrantController(oauth2Config *Config) *AuthorizationGrantController {
	return &AuthorizationGrantController{
		config: oauth2Config,
		store:  oauth2Config.Store,
	}
}

/**
 * Check Request Params
 *
 * @param  {Function} done
 * @this   OAuth
 */
func (a *AuthorizationGrantController) HandleForm(c *cocktail.Context) {
	info := CreateOAuth2Client(c)
	values := c.Queries

	/* Condition validation: Validate request */
	if values == nil {
		err := Status{Status: 400, Error: "Invalid Request", ErrorDescription: "Invalid Request."}
		c.OutputError(&err)
		return
	}

	/* Condition validation: Validate client_id */
	if len(info.ClientId) == 0 {
		err := Status{Status: 400, Error: "Invalid Request", ErrorDescription: "Missing client_id parameter."}
		c.OutputError(&err)
		return
	}

	/* Condition validation: Validate response_type */
	responseType := values.Get(ResponseType)
	if len(responseType) == 0 {
		err := common.Status{Status: 400, Error: "Invalid Request", ErrorDescription: "Missing response_type parameter."}
		c.OutputError(&err)
		return
	} else if responseType != "code" {
		err := common.Status{Status: 400, Error: "Invalid Request", ErrorDescription: "Invalid response_type parameter (must be \"code\")."}
		c.OutputError(&err)
		return
	}

	// Note: Allow missing redirect_uri parameter
	//	/* Condition validation: Validate redirect_uri */
	//	redirectUri := values.Get(RedirectUri)
	//	if len(RedirectUri) == 0 {
	//		err := common.Status{Status: 400, Error: "Invalid Request", ErrorDescription: "Missing redirect_uri parameter."}
	//		c.OutputError(&err)
	//		return
	//	}

	/* Condition validation: Validate client credentials */
	client := a.store.FindClient(info.ClientId, info.ClientSecret, info.RedirectUri)
	if client == nil {
		err := common.Status{Status: 400, Error: "Invalid Request", ErrorDescription: "Invalid client credentials."}
		c.OutputError(&err)
		return
	}

	// Generate authorization code
	authorizationCode := GenerateToken()
	a.store.SaveAuthorizationCode(authorizationCode, info.ClientId, time.Now().Add(a.config.Lifetime.AuthorizationCode))

	// Redirect
	state := values.Get(State)
	if len(state) == 0 {
		c.OutputRedirect(common.Status302(), fmt.Sprintf("%s?code=%s", info.RedirectUri, authorizationCode))
	} else {
		c.OutputRedirect(common.Status302(), fmt.Sprintf("%s?code=%s&state=%s", info.RedirectUri, authorizationCode, state))
	}
}

///**
// * Check client against model
// *
// * @param  {Function} done
// * @this   OAuth
// */
//func (a *AuthorizationController) CheckUserApproved() {
//	// var self = this;
//	// this.check(this.req, function (err, allowed, user) {
//	//   if (err) return done(error('server_error', false, err));
//	//   if (!allowed) {
//	//     return done(error('access_denied',
//	//       'The user denied access to your application'));
//	//   }
//	//   self.user = user;
//	//   done();
//	// });
//}
