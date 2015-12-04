package oauth2

import (
	"regexp"
	"strings"

	"github.com/phuc0302/go-oauth2/utils"
)

// Pre-compile regex
var bearerRegex = regexp.MustCompile("^Bearer\\s(\\w+)$")

// getBearerToken extract token from request according to RFC6750
func (c *Context) getBearerToken() string {
	headerToken := strings.Trim(c.Header("Authorization"), " ")
	//	getToken := c.Queries.Get(AccessToken)
	//	postToken := this.req.body ? this.req.body.access_token : undefined;

	//	 // Check exactly one method was used
	//	 var methodsUsed = (headerToken !== undefined) + (getToken !== undefined) +
	//	   (postToken !== undefined);
	//	 if (methodsUsed > 1) {
	//	   return done(error('invalid_request',
	//	     'Only one method may be used to authenticate at a time (Auth header,  ' +
	//	       'GET or POST).'));
	//	 } else if (methodsUsed === 0) {
	//	   return done(error('invalid_request', 'The access token was not found'));
	//	 }

	// Header: http://tools.ietf.org/html/rfc6750#section-2.1
	isMatched := bearerRegex.MatchString(headerToken)
	if !isMatched {
		status := utils.Status401()
		status.Description = "Malformed auth header"

		c.OutputError(status)
		return ""
	} else {
		headerToken = headerToken[7:]
	}
	return headerToken
}

/**
 * Check token
 *
 * Check it against model, ensure it's not expired
 * @param  {Function} done
 * @this   OAuth
 */
func (c *Context) checkToken() {
	// var self = this;
	accessToken := c.getBearerToken()
	// this.model.getAccessToken(this.bearerToken, function (err, token) {
	//   if (err) return done(error('server_error', false, err));

	if len(accessToken) == 0 {
		//     return done(error('invalid_token',
		//       'The access token provided is invalid.'));
	}

	//   if (token.expires !== null &&
	//     (!token.expires || token.expires < new Date())) {
	//     return done(error('invalid_token',
	//       'The access token provided has expired.'));
	//   }

	//   // Expose params
	//   self.req.oauth = { bearerToken: token };
	//   self.req.user = token.user ? token.user : { id: token.userId };

	//   done();
	// });
}
