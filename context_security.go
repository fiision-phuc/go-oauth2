package oauth2

import "regexp"

var bearerRegex = regexp.MustCompile("^(B|b)earer\\s\\w+$")

// Security descripts a user's security scope.
type Security struct {
	//	AuthUser         AuthUser
	//	AuthClient       AuthClient
	//	AuthAccessToken  Token
	//	AuthRefreshToken Token
}

//func CreateSecurityContext(requestContext *Request, tokenStore oauth.TokenStore) *Security {
//	headerToken := strings.Trim(requestContext.Header["authorization"], " ")
//	isBearer := bearerRegex.MatchString(headerToken)

//	/* Condition validation: Validate existing of authorization header */
//	if !isBearer {
//		return nil
//	}

//	headerToken = headerToken[7:]
//	accessToken := tokenStore.FindAccessToken(headerToken)

//	/* Condition validation: Validate expiration time */
//	if accessToken == nil || accessToken.IsExpired() {
//		return nil
//	}

//	client := tokenStore.FindClientWithID(accessToken.GetClientID())
//	user := tokenStore.FindUserWithID(accessToken.GetUserID())
//	securityContext := &Security{
//		AuthClient:      client,
//		AuthUser:        user,
//		AuthAccessToken: accessToken,
//	}
//	return securityContext
//}
