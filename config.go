package oauth2

import "time"

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

const (
	ENV_HOST = "HOST"
	ENV_PORT = "PORT"

	ENV_HEADERS_SIZE  = "headers_size" // In Kb
	ENV_TIMEOUT_READ  = "timeout_read" // In seconds
	ENV_TIMEOUT_WRITE = "timeout_wrte" // In seconds

	ENV_ALLOW_METHODS  = "allow_methods"
	ENV_STATIC_FOLDERS = "static_folders"
)

type Config struct {
	Host string `json:"HOST,omitempty"`
	Port string `json:"PORT,omitempty"`

	HeaderSize   int           `json:"headers_size,omitempty"`
	TimeoutRead  time.Duration `json:"timeout_read,omitempty"`
	TimeoutWrite time.Duration `json:"timeout_wrte,omitempty"`

	AllowMethods  []string          `json:"allow_methods,omitempty"`
	StaticFolders map[string]string `json:"static_folders,omitempty"`
}
