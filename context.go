package oauth2

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/phuc0302/go-oauth2/utils"
)

// Context represent a request scope.
type Context struct {
	URLPath     string
	Queries     url.Values
	PathQueries map[string]string

	AuthUser         *User
	AuthClient       *Client
	AuthAccessToken  *AccessToken
	AuthRefreshToken *RefreshToken

	request  *http.Request
	response http.ResponseWriter
}

// CreateContext return a default context.
func CreateContext(request *http.Request, response http.ResponseWriter) *Context {
	context := &Context{
		URLPath:  request.URL.Path,
		request:  request,
		response: response,
	}

	// Parse body context if neccessary
	switch context.Method() {

	case GET:
		params := request.URL.Query()
		if len(params) > 0 {
			context.Queries = params
		}
		break

	case POST, PATCH:
		contentType := strings.ToLower(request.Header.Get("CONTENT-TYPE"))

		if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			params := utils.ParseForm(request)
			if len(params) > 0 {
				context.Queries = params
			}
		} else if strings.Contains(contentType, "multipart/form-data") {
			params := utils.ParseMultipartForm(request)

			if len(params) > 0 {
				context.Queries = params
			}
		}
		break

	default:
		break
	}
	return context
}
