package oauth2

import (
	"encoding/base64"
	"strings"

	"github.com/phuc0302/go-cocktail"
)

type BasicAuthentication struct {
	Username string
	Password string
}

// MARK: Struct's constructors
func CreateBasicAuthentication(c *cocktail.Context) *BasicAuthentication {
	/* Condition validation: Validate authorization header */
	AuthorizationHeader := c.Headers.Get("Authorization")
	if len(AuthorizationHeader) == 0 {
		return nil
	}

	/* Condition validation: Validate basic authentication format */
	tokens := strings.SplitN(AuthorizationHeader, " ", 2)
	if len(tokens) != 2 || tokens[0] != "Basic" {
		return nil
	}

	// Decode basic authentication
	bytes, err := base64.StdEncoding.DecodeString(tokens[1])
	if err != nil {
		return nil
	}

	pair := strings.SplitN(string(bytes), ":", 2)
	if len(pair) != 2 {
		return nil
	}

	return &BasicAuthentication{
		Username: pair[0],
		Password: pair[1],
	}
}
