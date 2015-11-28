package oauth2

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/url"
	"text/template"
)

type Params struct {
	Headers http.Header
	Queries url.Values

	PathQueries map[string]string
}
type Context struct {
	Params

	Method   string
	UrlPath  string
	Protocol string

	request  *http.Request
	response http.ResponseWriter
}

// MARK: Struct's constructors
func CreateContext(request *http.Request, response http.ResponseWriter) *Context {
	context := &Context{
		Method:   request.Method,
		UrlPath:  request.URL.Path,
		Protocol: request.Proto,

		request:  request,
		response: response,
	}

	context.Headers = request.Header
	return context
}

/** Get an upload file by name. */
func (c *Context) GetMultipartFile(name string) (multipart.File, *multipart.FileHeader, error) {
	return c.request.FormFile(name)
}

/** Set additional response header. */
func (c *Context) SetResponseHeader(headerName string, headerValue string) {
	c.response.Header().Set(headerName, headerValue)
}

/** Return result to client. */
func (c *Context) OutputError(status *Status) {
	c.response.Header().Set("Content-Type", "application/problem+json")
	c.response.WriteHeader(status.Status)

	cause, _ := json.Marshal(status)
	c.response.Write(cause)
}

func (c *Context) OutputRedirect(status *Status, url string) {
	http.Redirect(c.response, c.request, url, status.Status)
}

func (c *Context) OutputJson(status *Status, model interface{}) {
	c.response.Header().Set("Content-Type", "application/json")
	c.response.WriteHeader(status.Status)

	data, _ := json.Marshal(model)
	c.response.Write(data)
}

func (c *Context) OutputHtml(filePath string, model interface{}) {
	tmpl, error := template.ParseFiles(filePath)
	if error != nil {
		c.OutputError(Status404())
	} else {
		tmpl.Execute(c.response, model)
	}
}
