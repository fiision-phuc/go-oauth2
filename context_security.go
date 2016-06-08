package oauth2

// Security describes a user's security scope.
type Security struct {
	AuthUser         IUser
	AuthClient       IClient
	AuthAccessToken  IToken
	AuthRefreshToken IToken
}
