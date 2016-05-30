package oauth2

import "time"

// IToken descripts a token's characteristic, it can be either access token or refresh token.
type IToken interface {
	GetClientID() string
	GetUserID() string
	GetToken() string

	IsExpired() bool
	GetCreatedTime() time.Time
	GetExpiredTime() time.Time
}
