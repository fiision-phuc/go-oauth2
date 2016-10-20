package oauth2

import "github.com/phuc0302/go-mongo"

// DefaultFactory describes a default factory object.
type DefaultFactory struct {
}

// CreateStore creates new store component.
func (d *DefaultFactory) CreateStore() TokenStore {
	mongo.ConnectMongo()
	return &DefaultMongoStore{}
}
