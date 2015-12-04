package oauth2

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Client Types	(http://tools.ietf.org/html/rfc6749#section-2.1)
const (
	Confidential = "confidential"
	Public       = "public"
)

type Lifetime struct {
	AccessToken       time.Duration
	RefreshToken      time.Duration
	AuthorizationCode time.Duration
}

type Config struct {
	Store    TokenStore
	Grant    []string
	Lifetime Lifetime

	clientValidation *regexp.Regexp
	grantsValidation *regexp.Regexp
}

func DefaultConfig(tokenStore TokenStore, grantTypes []string) *Config {
	//	These endpoints are https://public-api.wordpress.com/oauth2/authorize and https://public-api.wordpress.com/oauth2/token

	//	regexp.MustCompile(`:[^/#?()\.\\]+`)
	//	^(password|refresh_token|authorization_code)$
	//	strings.Join(Grants, "|")

	// this.regex = {
	//   clientId: config.clientIdRegex || /^[a-z0-9-_]{3,40}$/i,
	//   grantType: new RegExp('^(' + this.grants.join('|') + ')$', 'i')
	// };

	return &Config{
		Store: tokenStore,
		Grant: grantTypes,

		Lifetime: Lifetime{
			AccessToken:       3600,
			RefreshToken:      1209600,
			AuthorizationCode: 30,
		},

		clientValidation: nil,
		grantsValidation: regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(grantTypes, "|"))),
	}
}
