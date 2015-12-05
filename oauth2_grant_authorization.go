package oauth2

import (
	"fmt"
	"time"

	"github.com/phuc0302/go-oauth2/utils"
)

type GrantAuthorization struct {
	config *Config
	store  Store
}

// MARK: Struct's constructors
func CreateAuthorizationGrantController(oauth2Config *Config) *GrantAuthorization {
	return &GrantAuthorization{
		config: oauth2Config,
		//		store:  oauth2Config.Store,
	}
}

/**
 * Check Request Params
 *
 * @param  {Function} done
 * @this   OAuth
 */
func (a *GrantAuthorization) HandleForm(c *RequestContext) {
	client := createClient(c)
	values := c.Queries
	/* Condition validation: Validate client_id */
	if len(client.ClientID) == 0 {
		err := utils.Status400()
		err.Description = "Missing client_id parameter."
		c.OutputError(err)
		return
	}

	/* Condition validation: Validate response_type */
	responseType := values.Get("response_type")
	if len(responseType) == 0 {
		err := utils.Status400()
		err.Description = "Missing response_type parameter."
		c.OutputError(err)
		return
	} else if responseType != "code" {
		err := utils.Status400()
		err.Description = "Invalid response_type parameter (must be \"code\")."
		c.OutputError(err)
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
	client = a.store.FindClientWithCredential(client.ClientID, client.ClientSecret)
	if client == nil {
		err := utils.Status400()
		err.Description = "Invalid client credentials."
		c.OutputError(err)
		return
	}

	// Generate authorization code
	authorizationCode := utils.GenerateToken()
	a.store.SaveAuthorizationCode(authorizationCode, client.ClientID, time.Now().Add(a.config.DurationAuthorizationCode))

	// Redirect
	state := values.Get("state")
	if len(state) == 0 {
		c.OutputRedirect(utils.Status302(), fmt.Sprintf("%s?code=%s", client.RedirectURI, authorizationCode))
	} else {
		c.OutputRedirect(utils.Status302(), fmt.Sprintf("%s?code=%s&state=%s", client.RedirectURI, authorizationCode, state))
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
