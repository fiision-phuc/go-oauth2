package oauth2

import (
	"regexp"

	"github.com/phuc0302/go-mongo"
	"github.com/phuc0302/go-oauth2/oauth_role"
	"github.com/phuc0302/go-server"
)

// Global variables.
var (
	// Global public config's instance.
	Cfg *Config

	// Global public token store's instance.
	Store TokenStore

	// OAuth2 grant regex.
	grantsValidation *regexp.Regexp

	// Bearer regex.
	bearerFinder = regexp.MustCompile("^(B|b)earer\\s.+$")
)

// InitializeWithMongoDB will init server with MongoDB; either in sandbox mode or production mode,
// will register oauth2 service or not.
//
// @param
// - sandboxMode {bool} (instruction in which config file should be loaded)
// - sandboxMode {bool} (instruction in which should bind authorize & token service or not)
func InitializeWithMongoDB(sandboxMode bool, bindService bool) {
	mongo.ConnectMongo()
	Initialize(nil, sandboxMode, bindService)
}

// Initialize will init server as above func. However, the database will be your choice.
//
// @param
// - tokenStore {TokenStore} (your own token store implementation. If null, default will be used)
// - sandboxMode {bool} (instruction in which config file should be loaded)
// - sandboxMode {bool} (instruction in which should bind authorize & token service or not)
func Initialize(tokenStore TokenStore, sandboxMode bool, bindService bool) {
	server.Initialize(sandboxMode)
	if tokenStore == nil {
		tokenStore = CreateMongoDBStore()
	}

	// Register global components
	Store = tokenStore

	// Setup OAuth2.0
	if bindService {
		//	grantAuthorization := new(AuthorizationGrant)
		tokenGrant := new(TokenGrant)

		//	server.Get("/authorize", grantAuthorization.HandleForm)
		server.BindGet("/token", tokenGrant.HandleForm)
		server.BindPost("/token", tokenGrant.HandleForm)
	}
}

// Run will start HTTP server.
func Run() {
	server.Run()
}

// RunTLS will start HTTPS server.
func RunTLS(certFile string, keyFile string) {
	server.RunTLS(certFile, keyFile)
}

// GroupRoute is a wrapper func for server.GroupRoute func.
//
// @param
// - prefixURI {string} (the prefix for url)
// - handler {server.HandleGroupFunc} (the callback func)
func GroupRoute(prefixURI string, groupHandler server.HandleGroupFunc) {
	server.GroupRoute(prefixURI, groupHandler)
}

// BindCopy is a wrapper func for server.BindCopy func. By default, this will add ValidateToken
// and ValidateRoles around HandleContextFunc. If roles is not defined, by default, a route will
// accept all roles.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - roles {[]string} (a list of acceptable users' roles)
// - handler {server.HandleContextFunc} (the callback func)
func BindCopy(patternURL string, roles []string, handler server.HandleContextFunc) {
	if roles == nil || len(roles) == 0 {
		roles = oauthRole.All()
	}
	server.BindCopy(patternURL, server.Adapt(handler, ValidateToken(), ValidateRoles(roles...)))
}

// BindDelete is a wrapper func for server.BindDelete func. By default, this will add ValidateToken
// and ValidateRoles around HandleContextFunc. If roles is not defined, by default, a route will
// accept all roles.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - roles {[]string} (a list of acceptable users' roles)
// - handler {server.HandleContextFunc} (the callback func)
func BindDelete(patternURL string, roles []string, handler server.HandleContextFunc) {
	if roles == nil || len(roles) == 0 {
		roles = oauthRole.All()
	}
	server.BindDelete(patternURL, server.Adapt(handler, ValidateToken(), ValidateRoles(roles...)))
}

// BindGet is a wrapper func for server.BindGet func. By default, this will add ValidateToken
// and ValidateRoles around HandleContextFunc. If roles is not defined, by default, a route will
// accept all roles.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - roles {[]string} (a list of acceptable users' roles)
// - handler {server.HandleContextFunc} (the callback func)
func BindGet(patternURL string, roles []string, handler server.HandleContextFunc) {
	if roles == nil || len(roles) == 0 {
		roles = oauthRole.All()
	}
	server.BindGet(patternURL, server.Adapt(handler, ValidateToken(), ValidateRoles(roles...)))
}

// BindHead is a wrapper func for server.BindHead func. By default, this will add ValidateToken
// and ValidateRoles around HandleContextFunc. If roles is not defined, by default, a route will
// accept all roles.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - roles {[]string} (a list of acceptable users' roles)
// - handler {server.HandleContextFunc} (the callback func)
func BindHead(patternURL string, roles []string, handler server.HandleContextFunc) {
	if roles == nil || len(roles) == 0 {
		roles = oauthRole.All()
	}
	server.BindHead(patternURL, server.Adapt(handler, ValidateToken(), ValidateRoles(roles...)))
}

// BindLink is a wrapper func for server.BindLink func. By default, this will add ValidateToken
// and ValidateRoles around HandleContextFunc. If roles is not defined, by default, a route will
// accept all roles.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - roles {[]string} (a list of acceptable users' roles)
// - handler {server.HandleContextFunc} (the callback func)
func BindLink(patternURL string, roles []string, handler server.HandleContextFunc) {
	if roles == nil || len(roles) == 0 {
		roles = oauthRole.All()
	}
	server.BindLink(patternURL, server.Adapt(handler, ValidateToken(), ValidateRoles(roles...)))
}

// BindOptions is a wrapper func for server.BindOptions func. By default, this will add ValidateToken
// and ValidateRoles around HandleContextFunc. If roles is not defined, by default, a route will
// accept all roles.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - roles {[]string} (a list of acceptable users' roles)
// - handler {server.HandleContextFunc} (the callback func)
func BindOptions(patternURL string, roles []string, handler server.HandleContextFunc) {
	if roles == nil || len(roles) == 0 {
		roles = oauthRole.All()
	}
	server.BindOptions(patternURL, server.Adapt(handler, ValidateToken(), ValidateRoles(roles...)))
}

// BindPatch is a wrapper func for server.BindPatch func. By default, this will add ValidateToken
// and ValidateRoles around HandleContextFunc. If roles is not defined, by default, a route will
// accept all roles.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - roles {[]string} (a list of acceptable users' roles)
// - handler {server.HandleContextFunc} (the callback func)
func BindPatch(patternURL string, roles []string, handler server.HandleContextFunc) {
	if roles == nil || len(roles) == 0 {
		roles = oauthRole.All()
	}
	server.BindPatch(patternURL, server.Adapt(handler, ValidateToken(), ValidateRoles(roles...)))
}

// BindPost is a wrapper func for server.BindPost func. By default, this will add ValidateToken
// and ValidateRoles around HandleContextFunc. If roles is not defined, by default, a route will
// accept all roles.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - roles {[]string} (a list of acceptable users' roles)
// - handler {server.HandleContextFunc} (the callback func)
func BindPost(patternURL string, roles []string, handler server.HandleContextFunc) {
	if roles == nil || len(roles) == 0 {
		roles = oauthRole.All()
	}
	server.BindPost(patternURL, server.Adapt(handler, ValidateToken(), ValidateRoles(roles...)))
}

// BindPurge is a wrapper func for server.BindPurge func. By default, this will add ValidateToken
// and ValidateRoles around HandleContextFunc. If roles is not defined, by default, a route will
// accept all roles.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - roles {[]string} (a list of acceptable users' roles)
// - handler {server.HandleContextFunc} (the callback func)
func BindPurge(patternURL string, roles []string, handler server.HandleContextFunc) {
	if roles == nil || len(roles) == 0 {
		roles = oauthRole.All()
	}
	server.BindPurge(patternURL, server.Adapt(handler, ValidateToken(), ValidateRoles(roles...)))
}

// BindPut is a wrapper func for server.BindPut func. By default, this will add ValidateToken
// and ValidateRoles around HandleContextFunc. If roles is not defined, by default, a route will
// accept all roles.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - roles {[]string} (a list of acceptable users' roles)
// - handler {server.HandleContextFunc} (the callback func)
func BindPut(patternURL string, roles []string, handler server.HandleContextFunc) {
	if roles == nil || len(roles) == 0 {
		roles = oauthRole.All()
	}
	server.BindPut(patternURL, server.Adapt(handler, ValidateToken(), ValidateRoles(roles...)))
}

// BindUnlink is a wrapper func for server.BindUnlink func. By default, this will add ValidateToken
// and ValidateRoles around HandleContextFunc. If roles is not defined, by default, a route will
// accept all roles.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - roles {[]string} (a list of acceptable users' roles)
// - handler {server.HandleContextFunc} (the callback func)
func BindUnlink(patternURL string, roles []string, handler server.HandleContextFunc) {
	if roles == nil || len(roles) == 0 {
		roles = oauthRole.All()
	}
	server.BindUnlink(patternURL, server.Adapt(handler, ValidateToken(), ValidateRoles(roles...)))
}
