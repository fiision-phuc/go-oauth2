go-oauth2
=========

OAuth2 Provider library written in Go
-------------------------------------

This is a ported project from [OAuth2 Provider][3ee06010] that had been
written in Node.js with some additional features.

-   Pure implementation in GoLang.
-   Current implementation only supports password_flow & refresh_token_flow.
-   Use [JWT][368ba6d5].
-   Default buildin with MongoDB.
-   Allow to customize the server.

### Example Server

````go
package main

import (
	"fmt"

	"github.com/phuc0302/go-oauth2"
	"github.com/phuc0302/go-oauth2/oauth_key"
	"github.com/phuc0302/go-oauth2/oauth_role"
	"github.com/phuc0302/go-server"
	"github.com/phuc0302/go-server/util"
)

func main() {
	// Initialize server with sandbox mode enable and using MongoDB.
	oauth2.InitializeWithMongoDB(true, true)

	// Define handler
	f := func(c *server.RequestContext) {
		if s, ok := c.GetExtra(oauthKey.Context).(oauth2.OAuthContext); ok {
			c.OutputText(util.Status200(), fmt.Sprintf("Hello, your ID is: %s", s.User.UserID()))
		} else {
			panic(util.Status401())
		}
	}

	// Bind handler with HTTP GET
	oauth2.BindGet("/protected", oauthRole.All(), f)

	// Start server
	oauth2.Run()
}
````

### Author

Phuc, Tran Huu
phuc@fiisionstudio.com

[3ee06010]: https://github.com/t1msh/node-oauth20-provider "OAuth2.0 Provider"
[368ba6d5]: https://github.com/dgrijalva/jwt-go "Json Web Token"
