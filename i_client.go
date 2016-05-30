package oauth2

// IClient descripts a client's characteristic.
type IClient interface {
	GetClientID() string
	GetClientSecret() string
	GetGrantTypes() []string   // Server side only
	GetRedirectURIs() []string // Server side only
}
