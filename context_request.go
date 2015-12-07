package oauth2

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/phuc0302/go-oauth2/utils"
)

// RequestContext represent a request scope.
type RequestContext struct {
	URLPath     string
	Queries     url.Values
	PathQueries map[string]string

	AuthUser         AuthUser
	AuthClient       AuthClient
	AuthAccessToken  Token
	AuthRefreshToken Token

	request  *http.Request
	response http.ResponseWriter
}

// CreateContext return a default context.
func CreateRequestContext(request *http.Request, response http.ResponseWriter) *RequestContext {
	context := &RequestContext{
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

////////////////////////////////////////////////////////////////////////////////
// INPUT																	    //
////////////////////////////////////////////////////////////////////////////////

func (c *RequestContext) BasicAuth() (username string, password string, ok bool) {
	return c.request.BasicAuth()
}

func (c *RequestContext) Header(headerName string) string {
	return c.request.Header.Get(headerName)
}

func (c *RequestContext) Method() string {
	return c.request.Method
}

func (c *RequestContext) Protocol() string {
	return c.request.Proto
}

// BindForm converts urlencode/multipart form to object.
func (c *RequestContext) BindForm(inputForm interface{}) error {
	return utils.BindForm(c.Queries, inputForm)
}

// BindJSON converts json data to object.
func (c *RequestContext) BindJSON(jsonObject interface{}) error {
	//	return c.request.FormFile(name)
	return nil
}

// GetMultipartFile return an upload file by name.
func (c *RequestContext) MultipartFile(name string) (multipart.File, *multipart.FileHeader, error) {
	return c.request.FormFile(name)
}

////////////////////////////////////////////////////////////////////////////////
// Output						   											    //
////////////////////////////////////////////////////////////////////////////////

// OutputHeader return an additional header.
func (c *RequestContext) OutputHeader(headerName string, headerValue string) {
	c.response.Header().Set(headerName, headerValue)
}

// OutputError return an error json.
func (c *RequestContext) OutputError(status *utils.Status) {
	c.response.Header().Set("Content-Type", "application/problem+json")
	c.response.WriteHeader(status.Code)

	cause, _ := json.Marshal(status)
	c.response.Write(cause)
}

// OutputRedirect return a redirect instruction.
func (c *RequestContext) OutputRedirect(status *utils.Status, url string) {
	http.Redirect(c.response, c.request, url, status.Code)
}

// OutputJSON return a json.
func (c *RequestContext) OutputJSON(status *utils.Status, model interface{}) {
	c.response.Header().Set("Content-Type", "application/json")
	c.response.WriteHeader(status.Code)

	data, _ := json.Marshal(model)
	c.response.Write(data)
}

// OutputHTML will render a html page.
func (c *RequestContext) OutputHTML(filePath string, model interface{}) {
	tmpl, error := template.ParseFiles(filePath)
	if error != nil {
		c.OutputError(utils.Status404())
	} else {
		tmpl.Execute(c.response, model)
	}
}
