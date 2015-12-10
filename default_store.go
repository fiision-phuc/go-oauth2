package oauth2

import (
	"time"

	"github.com/phuc0302/go-oauth2/utils"

	"gopkg.in/mgo.v2/bson"
)

// AuthClientDefault descripts a mongodb AuthClient document.
type AuthClientDefault struct {
	ClientID     bson.ObjectId `bson:"_id,omitempty"`
	ClientSecret bson.ObjectId `bson:"client_secret,omitempty"`
	GrantTypes   []string      `bson:"grant_types,omitempty"`
	RedirectURIs []string      `bson:"redirect_uris,omitempty"`
}

// GetClientID returns client_id.
func (a *AuthClientDefault) GetClientID() string { return a.ClientID.Hex() }

// GetClientSecret returns client_secret.
func (a *AuthClientDefault) GetClientSecret() string { return a.ClientSecret.Hex() }

// GetGrantTypes returns grant_types.
func (a *AuthClientDefault) GetGrantTypes() []string { return a.GrantTypes }

// GetRedirectURIs returns redirect_uris.
func (a *AuthClientDefault) GetRedirectURIs() []string { return a.RedirectURIs }

////////////////////////////////////////////////////////////////////////////////////////////////////

// AuthUserDefault descripts a mongodb AuthUser document.
type AuthUserDefault struct {
	UserID   bson.ObjectId `bson:"_id,omitempty"`
	Username string        `bson:"username,omitempty"`
	Password string        `bson:"password,omitempty"`
	Roles    []string      `bson:"roles,omitempty"`
}

// GetUserID returns user_id.
func (a *AuthUserDefault) GetUserID() string { return a.UserID.Hex() }

// GetUsername returns username.
func (a *AuthUserDefault) GetUsername() string { return a.Username }

// GetPassword returns password.
func (a *AuthUserDefault) GetPassword() string { return a.Password }

// GetUserRoles returns roles.
func (a *AuthUserDefault) GetUserRoles() []string { return a.Roles }

/////////////////////////////////////////////////////////////////////////////////////////////////

// TokenDefault descripts a mongodb Token document.
type TokenDefault struct {
	TokenID     bson.ObjectId `bson:"_id,omitempty"`
	UserID      bson.ObjectId `bson:"user_id,omitempty"`
	ClientID    bson.ObjectId `bson:"client_id,omitempty"`
	Token       string        `bson:"token,omitempty"`
	CreatedTime time.Time     `bson:"created_time,omitempty"`
	ExpiredTime time.Time     `bson:"expired_time,omitempty"`
}

// GetTokenID returns token_id.
func (t *TokenDefault) GetTokenID() string { return t.TokenID.Hex() }

// GetClientID returns client_id.
func (t *TokenDefault) GetClientID() string { return t.ClientID.Hex() }

// GetUserID returns user_id.
func (t *TokenDefault) GetUserID() string { return t.UserID.Hex() }

// GetToken returns token.
func (t *TokenDefault) GetToken() string { return t.Token }

// SetToken updates token.
func (t *TokenDefault) SetToken(token string) { t.Token = token }

// IsExpired validate if this token is expired or not.
func (t *TokenDefault) IsExpired() bool { return time.Now().Unix() >= t.ExpiredTime.Unix() }

// GetCreatedTime returns created_time.
func (t *TokenDefault) GetCreatedTime() time.Time { return t.CreatedTime }

// SetCreatedTime updates created_time.
func (t *TokenDefault) SetCreatedTime(createdTime time.Time) { t.CreatedTime = createdTime }

// GetExpiredTime returns expired_time.
func (t *TokenDefault) GetExpiredTime() time.Time { return t.ExpiredTime }

// SetExpiredTime updates expired_time.
func (t *TokenDefault) SetExpiredTime(expiredTime time.Time) { t.ExpiredTime = expiredTime }

/////////////////////////////////////////////////////////////////////////////////////////////////

// InMemoryTokenStore descripts an in memory token store.
type InMemoryTokenStore struct {
	clients       []AuthClient
	users         []AuthUser
	accessTokens  []Token
	refreshTokens []Token
}

// FindUserWithID returns user with user_id.
func (s *InMemoryTokenStore) FindUserWithID(userID string) AuthUser {
	/* Condition validation */
	if len(userID) == 0 {
		return nil
	}

	for _, user := range s.users {
		if user.GetUserID() == userID {
			return user
		}
	}
	return nil
}

// FindUserWithClient returns user associated with client_id and client_secret.
func (s *InMemoryTokenStore) FindUserWithClient(clientID string, clientSecret string) AuthUser {
	/* Condition validation */
	if len(clientID) == 0 || len(clientSecret) == 0 {
		return nil
	}

	return nil
}

// FindUserWithCredential returns user associated with username and password.
func (s *InMemoryTokenStore) FindUserWithCredential(username string, password string) AuthUser {
	/* Condition validation */
	if len(username) == 0 || len(password) == 0 {
		return nil
	}

	for _, user := range s.users {
		if user.GetUsername() == username && user.GetPassword() == password {
			return user
		}
	}
	return nil
}

// FindClientWithID returns user associated with client_id.
func (s *InMemoryTokenStore) FindClientWithID(clientID string) AuthClient {
	/* Condition validation */
	if len(clientID) == 0 {
		return nil
	}

	for _, client := range s.clients {
		if client.GetClientID() == clientID {
			return client
		}
	}
	return nil
}

// FindClientWithCredential returns client with client_id and client_secret.
func (s *InMemoryTokenStore) FindClientWithCredential(clientID string, clientSecret string) AuthClient {
	for _, client := range s.clients {
		if client.GetClientID() == clientID && client.GetClientSecret() == clientSecret {
			return client
		}
	}
	return nil
}

// FindAccessToken returns access_token.
func (s *InMemoryTokenStore) FindAccessToken(token string) Token {
	for _, recordToken := range s.accessTokens {
		if recordToken.GetToken() == token {
			return recordToken
		}
	}
	return nil
}

// FindAccessTokenWithCredential returns access_token associated with client_id and user_id.
func (s *InMemoryTokenStore) FindAccessTokenWithCredential(clientID string, userID string) Token {
	for _, recordToken := range s.accessTokens {
		if recordToken.GetUserID() == userID && recordToken.GetClientID() == clientID {
			return recordToken
		}
	}
	return nil
}

// CreateAccessToken returns new access_token.
func (s *InMemoryTokenStore) CreateAccessToken(clientID string, userID string, token string, createdTime time.Time, expiredTime time.Time) Token {
	newToken := &TokenDefault{
		TokenID:     bson.NewObjectId(),
		UserID:      bson.ObjectIdHex(userID),
		ClientID:    bson.ObjectIdHex(clientID),
		Token:       utils.GenerateToken(),
		CreatedTime: createdTime,
		ExpiredTime: expiredTime,
	}

	s.accessTokens = append(s.accessTokens, newToken)
	return newToken
}

// DeleteAccessToken deletes access_token.
func (s *InMemoryTokenStore) DeleteAccessToken(token Token) {
	for idx, recordToken := range s.accessTokens {
		if recordToken == token {
			s.accessTokens = append(s.accessTokens[:idx], s.accessTokens[idx+1:]...)
			break
		}
	}
}

// SaveAccessToken saves/updates access_token.
func (s *InMemoryTokenStore) SaveAccessToken(token Token) {
	isUpdated := false
	for _, recordToken := range s.accessTokens {
		if recordToken == token {
			token.SetToken(token.GetToken())
			token.SetCreatedTime(token.GetCreatedTime())
			token.SetExpiredTime(token.GetExpiredTime())
			isUpdated = true
			break
		}
	}

	if !isUpdated {
		s.accessTokens = append(s.accessTokens, token)
	}
}

// FindRefreshToken returns refresh_token.
func (s *InMemoryTokenStore) FindRefreshToken(token string) Token {
	for _, recordToken := range s.refreshTokens {
		if recordToken.GetToken() == token {
			return recordToken
		}
	}
	return nil
}

// FindRefreshTokenWithCredential returns refresh_token associated with client_id and user_id.
func (s *InMemoryTokenStore) FindRefreshTokenWithCredential(clientID string, userID string) Token {
	for _, recordToken := range s.refreshTokens {
		if recordToken.GetUserID() == userID && recordToken.GetClientID() == clientID {
			return recordToken
		}
	}
	return nil
}

// CreateRefreshToken returns new refresh_token.
func (s *InMemoryTokenStore) CreateRefreshToken(clientID string, userID string, token string, createdTime time.Time, expiredTime time.Time) Token {
	newToken := &TokenDefault{
		TokenID:     bson.NewObjectId(),
		UserID:      bson.ObjectIdHex(userID),
		ClientID:    bson.ObjectIdHex(clientID),
		Token:       utils.GenerateToken(),
		CreatedTime: createdTime,
		ExpiredTime: expiredTime,
	}

	s.refreshTokens = append(s.refreshTokens, newToken)
	return newToken
}

// DeleteRefreshToken deletes refresh_token.
func (s *InMemoryTokenStore) DeleteRefreshToken(token Token) {
	for idx, recordToken := range s.refreshTokens {
		if recordToken == token {
			s.refreshTokens = append(s.refreshTokens[:idx], s.refreshTokens[idx+1:]...)
			break
		}
	}
}

// SaveRefreshToken saves/updates refresh_token.
func (s *InMemoryTokenStore) SaveRefreshToken(token Token) {
	isUpdated := false
	for _, recordToken := range s.refreshTokens {
		if recordToken == token {
			recordToken.SetToken(token.GetToken())
			recordToken.SetCreatedTime(token.GetCreatedTime())
			recordToken.SetExpiredTime(token.GetExpiredTime())
			isUpdated = true
			break
		}
	}

	if !isUpdated {
		s.refreshTokens = append(s.refreshTokens, token)
	}
}

//func (s *InMemoryTokenStore) FindAuthorizationCode(authorizationCode string) {
//}
//func (s *InMemoryTokenStore) SaveAuthorizationCode(authorizationCode string, clientID string, expires time.Time) {
//}
