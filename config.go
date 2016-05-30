package oauth2

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/phuc0302/go-oauth2/utils"
)

// Define HTTP Methods.
const (
	COPY    = "COPY"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	LINK    = "LINK"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PURGE   = "PURGE"
	PUT     = "PUT"
	UNLINK  = "UNLINK"
)

// Define configuration file's name.
const (
	debug   = "oauth2.debug.cfg"
	release = "oauth2.release.cfg"
)

// Define OAuth2 flows.
const (
	AuthorizationCodeGrant = "authorization_code" // For apps running on a web server
	ClientCredentialsGrant = "client_credentials" // For application access
	ImplicitGrant          = "implicit"           // For browser-based or mobile apps
	PasswordGrant          = "password"           // For logging in with a username and password
	RefreshTokenGrant      = "refresh_token"      // Should allow refresh token or not

)

// Define OAuth2 tables.
const (
	TableRefreshToken = "oauth_refresh_token"
	TableAccessToken  = "oauth_access_token"
	TableClient       = "oauth_client"
	TableUser         = "oauth_user"
)

// config descripts a configuration  object  that  will  be  used  during application life time.
type config struct {
	Host         string        `json:"host,omitempty"`
	Port         string        `json:"port,omitempty"`
	TLSPort      string        `json:"tls_port,omitempty"`
	HeaderSize   int           `json:"headers_size,omitempty"`  // In KB
	ReadTimeout  time.Duration `json:"timeout_read,omitempty"`  // In seconds
	WriteTimeout time.Duration `json:"timeout_write,omitempty"` // In seconds

	AllowMethods         []string          `json:"allow_methods,omitempty"`
	StaticFolders        map[string]string `json:"static_folders,omitempty"`
	ReverseRedirectPaths map[string]int    `json:"redirect_paths,omitempty"`

	//	GrantTypes                []string      `json:"grant_types,omitempty"`
	//	AllowRefreshToken         bool          `json:"allow_refresh_token,omitempty"`
	//	AccessTokenDuration       time.Duration `json:"access_token_duration,omitempty"`       // In seconds
	//	RefreshTokenDuration      time.Duration `json:"refresh_token_duration,omitempty"`      // In seconds
	//	AuthorizationCodeDuration time.Duration `json:"authorization_code_duration,omitempty"` // In seconds

	// Validation
	RedirectPaths     map[int]string `json:"-"`
	ClientValidation  *regexp.Regexp `json:"-"`
	GrantsValidation  *regexp.Regexp `json:"-"`
	MethodsValidation *regexp.Regexp `json:"-"`
}

// createConfig generates a default configuration file.
func createConfig(configFile string) {
	if utils.FileExisted(configFile) {
		os.Remove(configFile)
	}

	// Create default config
	config := config{
		Host:         "localhost",
		Port:         "8080",
		TLSPort:      "8443",
		HeaderSize:   5,
		ReadTimeout:  15,
		WriteTimeout: 15,

		AllowMethods: []string{COPY, DELETE, GET, HEAD, LINK, OPTIONS, PATCH, POST, PURGE, PUT, UNLINK},
		ReverseRedirectPaths: map[string]int{
			"/login": 401,
		},
		StaticFolders: map[string]string{
			"/assets":    "assets",
			"/resources": "resources",
		},

		//		GrantTypes:                []string{AuthorizationCodeGrant, ClientCredentialsGrant, PasswordGrant, RefreshTokenGrant},
		//		AccessTokenDuration:       3600,
		//		RefreshTokenDuration:      1209600,
		//		AuthorizationCodeDuration: 30,
	}

	// Create new file
	configJSON, _ := json.MarshalIndent(config, "", "  ")
	file, _ := os.Create(configFile)
	file.Write(configJSON)
	file.Close()
}

// loadConfig retrieves previous configuration from file.
func loadConfig(configFile string) *config {
	// Generate config file if neccessary
	if !utils.FileExisted(configFile) {
		createConfig(configFile)
	}

	file, err := os.Open(configFile)
	if err != nil {
		return nil
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}

	config := config{}
	json.Unmarshal(bytes, &config)

	folders := make(map[string]string)
	for path, folder := range config.StaticFolders {
		folders[utils.FormatPath(path)] = folder
	}
	config.StaticFolders = folders

	//	// Convert duration to seconds
	//	config.AccessTokenDuration = config.AccessTokenDuration * time.Second
	//	config.RefreshTokenDuration = config.RefreshTokenDuration * time.Second
	//	config.AuthorizationCodeDuration = config.AuthorizationCodeDuration * time.Second

	//	// Define regular expressions
	//	//	regexp.MustCompile(`:[^/#?()\.\\]+`)
	//	config.grantsValidation = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(config.GrantTypes, "|")))
	//	config.methodsValidation = regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(config.AllowMethods, "|")))

	//	for _, grant := range config.GrantTypes {
	//		if grant == RefreshTokenGrant {
	//			config.AllowRefreshToken = true
	//			break
	//		}
	//	}
	return &config
}
