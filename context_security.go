package oauth2

import (
	"regexp"
	"strings"
)

var bearerRegex = regexp.MustCompile("^(B|b)earer\\s\\w+$")

// Security describes a user's security scope.
type Security struct {
	AuthUser         IUser
	AuthClient       IClient
	AuthAccessToken  IToken
	AuthRefreshToken IToken
}

func CreateSecurityContext(requestContext *Request, tokenStore IStore) *Security {
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

	client := tokenStore.FindClientWithID(accessToken.ClientID())
	user := tokenStore.FindUserWithID(accessToken.UserID())
	securityContext := &Security{
		AuthClient:      client,
		AuthUser:        user,
		AuthAccessToken: accessToken,
	}
	return securityContext
}
