package oauth2

// IClient descripts a client's characteristic.
type IClient interface {

	// Return client's ID.
	ClientID() string
	// Return client's secret.
	ClientSecret() string

	// Return client's allowed grant types.
	GrantTypes() []string
	// Return client's registered redirect URIs.
	RedirectURIs() []string
}
