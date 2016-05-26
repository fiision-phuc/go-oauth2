package context

import (
	"regexp"
	"strings"
)

var bearerRegex = regexp.MustCompile("^(B|b)earer\\s\\w+$")

// SecurityContext descripts a user's security scope.
type SecurityContext struct {
	AuthUser         AuthUser
	AuthClient       AuthClient
	AuthAccessToken  Token
	AuthRefreshToken Token
}

func CreateSecurityContext(requestContext *Request, tokenStore TokenStore) *SecurityContext {
	headerToken := strings.Trim(requestContext.Header["authorization"], " ")
	isBearer := bearerRegex.MatchString(headerToken)

	/* Condition validation: Validate existing of authorization header */
	if !isBearer {
		return nil
	}

	headerToken = headerToken[7:]
	accessToken := tokenStore.FindAccessToken(headerToken)

	/* Condition validation: Validate expiration time */
	if accessToken == nil || accessToken.IsExpired() {
		return nil
	}

	client := tokenStore.FindClientWithID(accessToken.GetClientID())
	user := tokenStore.FindUserWithID(accessToken.GetUserID())
	securityContext := &SecurityContext{
		AuthClient:      client,
		AuthUser:        user,
		AuthAccessToken: accessToken,
	}
	return securityContext
}
