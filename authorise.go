package oauth2

import "net/http"

type Authorise struct {
	Config *Config
	Model  *TokenStore
	Req    *http.Request
}

/**
 * Get bearer token
 *
 * Extract token from request according to RFC6750
 *
 * @param  {Function} done
 * @this   OAuth
 */
func (a *Authorise) getBearerToken() {
	// var headerToken = this.req.get('Authorization'),
	//   getToken =  this.req.query.access_token,
	//   postToken = this.req.body ? this.req.body.access_token : undefined;

	// // Check exactly one method was used
	// var methodsUsed = (headerToken !== undefined) + (getToken !== undefined) +
	//   (postToken !== undefined);

	// if (methodsUsed > 1) {
	//   return done(error('invalid_request',
	//     'Only one method may be used to authenticate at a time (Auth header,  ' +
	//       'GET or POST).'));
	// } else if (methodsUsed === 0) {
	//   return done(error('invalid_request', 'The access token was not found'));
	// }

	// // Header: http://tools.ietf.org/html/rfc6750#section-2.1
	// if (headerToken) {
	//   var matches = headerToken.match(/Bearer\s(\S+)/);

	//   if (!matches) {
	//     return done(error('invalid_request', 'Malformed auth header'));
	//   }

	//   headerToken = matches[1];
	// }

	// // POST: http://tools.ietf.org/html/rfc6750#section-2.2
	// if (postToken) {
	//   if (this.req.method === 'GET') {
	//     return done(error('invalid_request',
	//       'Method cannot be GET When putting the token in the body.'));
	//   }

	//   if (!this.req.is('application/x-www-form-urlencoded')) {
	//     return done(error('invalid_request', 'When putting the token in the ' +
	//       'body, content type must be application/x-www-form-urlencoded.'));
	//   }
	// }

	// this.bearerToken = headerToken || postToken || getToken;
	// done();
}

/**
 * Check token
 *
 * Check it against model, ensure it's not expired
 * @param  {Function} done
 * @this   OAuth
 */
func (a *Authorise) checkToken() {
	// var self = this;
	// this.model.getAccessToken(this.bearerToken, function (err, token) {
	//   if (err) return done(error('server_error', false, err));

	//   if (!token) {
	//     return done(error('invalid_token',
	//       'The access token provided is invalid.'));
	//   }

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
