package oauth2

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// DefaultInMemoryStore descripts an in memory store.
type DefaultInMemoryStore struct {
	clients       []IClient
	users         []IUser
	accessTokens  []IToken
	refreshTokens []IToken
}

// FindUserWithID returns user with user_id.
func (s *DefaultInMemoryStore) FindUserWithID(userID string) IUser {
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
func (s *DefaultInMemoryStore) FindUserWithClient(clientID string, clientSecret string) IUser {
	/* Condition validation */
	if len(clientID) == 0 || len(clientSecret) == 0 {
		return nil
	}
	return nil
}

// FindUserWithCredential returns user associated with username and password.
func (s *DefaultInMemoryStore) FindUserWithCredential(username string, password string) IUser {
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
func (s *DefaultInMemoryStore) FindClientWithID(clientID string) IClient {
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
func (s *DefaultInMemoryStore) FindClientWithCredential(clientID string, clientSecret string) IClient {
	/* Condition validation */
	if len(clientID) == 0 || len(clientSecret) == 0 {
		return nil
	}

	for _, client := range s.clients {
		if client.GetClientID() == clientID && client.GetClientSecret() == clientSecret {
			return client
		}
	}
	return nil
}

// FindAccessToken returns access_token.
func (s *DefaultInMemoryStore) FindAccessToken(token string) IToken {
	for _, recordToken := range s.accessTokens {
		if recordToken.GetToken() == token {
			return recordToken
		}
	}
	return nil
}

// FindAccessTokenWithCredential returns access_token associated with client_id and user_id.
func (s *DefaultInMemoryStore) FindAccessTokenWithCredential(clientID string, userID string) IToken {
	for _, recordToken := range s.accessTokens {
		if recordToken.GetUserID() == userID && recordToken.GetClientID() == clientID {
			return recordToken
		}
	}
	return nil
}

// CreateAccessToken returns new access_token.
func (s *DefaultInMemoryStore) CreateAccessToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) IToken {
	newToken := &DefaultToken{
		TokenID:     bson.NewObjectId(),
		UserID:      bson.ObjectIdHex(userID),
		ClientID:    bson.ObjectIdHex(clientID),
		CreatedTime: createdTime,
		ExpiredTime: expiredTime,
	}

	s.accessTokens = append(s.accessTokens, newToken)
	return newToken
}

// DeleteAccessToken deletes access_token.
func (s *DefaultInMemoryStore) DeleteAccessToken(token IToken) {
	for idx, recordToken := range s.accessTokens {
		if recordToken == token {
			s.accessTokens = append(s.accessTokens[:idx], s.accessTokens[idx+1:]...)
			break
		}
	}
}

// FindRefreshToken returns refresh_token.
func (s *DefaultInMemoryStore) FindRefreshToken(token string) IToken {
	for _, recordToken := range s.refreshTokens {
		if recordToken.GetToken() == token {
			return recordToken
		}
	}
	return nil
}

// FindRefreshTokenWithCredential returns refresh_token associated with client_id and user_id.
func (s *DefaultInMemoryStore) FindRefreshTokenWithCredential(clientID string, userID string) IToken {
	for _, recordToken := range s.refreshTokens {
		if recordToken.GetUserID() == userID && recordToken.GetClientID() == clientID {
			return recordToken
		}
	}
	return nil
}

// CreateRefreshToken returns new refresh_token.
func (s *DefaultInMemoryStore) CreateRefreshToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) IToken {
	newToken := &DefaultToken{
		TokenID:     bson.NewObjectId(),
		UserID:      bson.ObjectIdHex(userID),
		ClientID:    bson.ObjectIdHex(clientID),
		CreatedTime: createdTime,
		ExpiredTime: expiredTime,
	}

	s.refreshTokens = append(s.refreshTokens, newToken)
	return newToken
}

// DeleteRefreshToken deletes refresh_token.
func (s *DefaultInMemoryStore) DeleteRefreshToken(token IToken) {
	for idx, recordToken := range s.refreshTokens {
		if recordToken == token {
			s.refreshTokens = append(s.refreshTokens[:idx], s.refreshTokens[idx+1:]...)
			break
		}
	}
}

//func (s *InMemoryTokenStore) FindAuthorizationCode(authorizationCode string) {
//}
//func (s *InMemoryTokenStore) SaveAuthorizationCode(authorizationCode string, clientID string, expires time.Time) {
//}
