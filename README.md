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
import "github.com/phuc0302/go-oauth2"

// Create server with sandbox mode enable.
server := oauth2.DefaultServer(true)

// Define routing server.
server.Get("/protected", func(c *oauth2.Request, s *oauth2.Security) {
  c.OutputText(utils.Status200(), "This is a protected resources.")
})

// Define who is able to access protected resources.
server.GroupRole("/protected**", "r_user")

// Start server.
server.Run()
````

### Author

Phuc, Tran Huu
phuc@fiisionstudio.com

[3ee06010]: https://github.com/t1msh/node-oauth20-provider "OAuth2.0 Provider"
[368ba6d5]: https://github.com/dgrijalva/jwt-go "Json Web Token"
