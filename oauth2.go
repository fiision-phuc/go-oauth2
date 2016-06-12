package oauth2

// TokenResponse describes a granted response that will be returned to client.
type TokenResponse struct {
	TokenType    string `json:"token_type,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`

	Roles []string `json:"roles,omitempty"`
}
