package oauth2

import "github.com/phuc0302/go-oauth2/context"

type OAuth2Client struct {
	ClientId     string
	ClientSecret string
	RedirectUri  string
}

// MARK: Struct's constructors
func CreateOAuth2Client(c *context.Context) *OAuth2Client {
	basicAuth := CreateBasicAuthentication(c)

	if basicAuth != nil {
		return &OAuth2Client{
			ClientId:     basicAuth.Username,
			ClientSecret: basicAuth.Password,
			RedirectUri:  c.Queries.Get(RedirectUri),
		}
	} else {
		return &OAuth2Client{
			ClientId:     c.Queries.Get(ClientId),
			ClientSecret: c.Queries.Get(ClientSecret),
			RedirectUri:  c.Queries.Get(RedirectUri),
		}
	}
}
