package test

import (
	"os"

	"github.com/phuc0302/go-oauth2"
	"github.com/phuc0302/go-oauth2/mongo"
	"github.com/phuc0302/go-oauth2/utils"
	"gopkg.in/mgo.v2/bson"
)

// Template error message
const (
	TemplateError        = "Invalid %s parameter."
	TemplateErrorMessage = "Expected \"Invalid %s parameter.\" but found \"%s\"."
)

// SetupOAuth2Server returns an oauth2 server with default settings for testing purpose.
func SetupOAuth2Server(function func(s *oauth2.Server, tokenStore oauth2.TokenStore, client oauth2.AuthClient, admin oauth2.AuthUser)) {
	defer os.Remove(oauth2.ConfigFile)
	defer os.Remove(mongo.ConfigFile)
	mongo.ConnectMongo()

	// Clean up date after used.
	session, database := mongo.GetMonotonicSession()
	defer session.Close()
	defer database.DropDatabase()

	// Create single client for testing
	clientCollection := database.C(oauth2.TableClient)
	appClient := oauth2.AuthClientDefault{
		ClientID:     bson.NewObjectId(),
		ClientSecret: bson.NewObjectId(),
		GrantTypes:   []string{oauth2.PasswordGrant, oauth2.RefreshTokenGrant},
		RedirectURIs: []string{"http://www.sample01.com"},
	}
	clientCollection.Insert(appClient)

	// Create single user for testing
	userCollection := database.C(oauth2.TableUser)
	password, _ := utils.EncryptPassword("P@ssw0rd")
	admin := &oauth2.AuthUserDefault{
		UserID:   bson.NewObjectId(),
		Username: "admin",
		Password: password,
		Roles:    []string{"r_user", "r_admin"},
	}
	userCollection.Insert(admin)

	// Setup server
	tokenStore := &oauth2.MongoDBTokenStore{}
	server := oauth2.DefaultServerWithTokenStore(tokenStore)

	function(server, tokenStore, &appClient, admin)
}
