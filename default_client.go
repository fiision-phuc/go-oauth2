package oauth2

import "gopkg.in/mgo.v2/bson"

// defaultClient descripts a mongodb client.
type DefaultClient struct {
	ClientID     bson.ObjectId `bson:"_id,omitempty"`
	ClientSecret bson.ObjectId `bson:"client_secret,omitempty"`
	GrantTypes   []string      `bson:"grant_types,omitempty"`
	RedirectURIs []string      `bson:"redirect_uris,omitempty"`
}

// GetClientID returns client_id.
func (a *DefaultClient) GetClientID() string { return a.ClientID.Hex() }

// GetClientSecret returns client_secret.
func (a *DefaultClient) GetClientSecret() string { return a.ClientSecret.Hex() }

// GetGrantTypes returns grant_types.
func (a *DefaultClient) GetGrantTypes() []string { return a.GrantTypes }

// GetRedirectURIs returns redirect_uris.
func (a *DefaultClient) GetRedirectURIs() []string { return a.RedirectURIs }
