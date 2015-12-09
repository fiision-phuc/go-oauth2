package oauth2

import (
	"regexp"
	"strings"

	"github.com/phuc0302/go-oauth2/utils"
)

// Pre-compile regex
var bearerRegex = regexp.MustCompile("^Bearer\\s(\\w+)$")

type SecurityContext struct {
	AuthUser         AuthUser
	AuthClient       AuthClient
	AuthAccessToken  Token
	AuthRefreshToken Token
}

func CreateSecurityContextWithRequestContext(requestContext *RequestContext, tokenStore TokenStore) (*SecurityContext, *utils.Status) {
	headerToken := strings.Trim(requestContext.Header["authorization"], " ")
	isMatched := bearerRegex.MatchString(headerToken)

	if !isMatched {
		return nil, utils.Status401WithDescription("Malformed auth header")
	} else {
		headerToken = headerToken[7:]
		accessToken := tokenStore.FindAccessToken(headerToken)

		if accessToken == nil || accessToken.IsExpired() {
			return nil, utils.Status401WithDescription("Malformed auth header")
		} else {
			client := tokenStore.FindClientWithID(accessToken.GetClientID())
			user := tokenStore.FindUserWithID(accessToken.GetUserID())

			securityContext := &SecurityContext{
				AuthClient:      client,
				AuthUser:        user,
				AuthAccessToken: accessToken,
			}
			return securityContext, nil
		}
	}
}
