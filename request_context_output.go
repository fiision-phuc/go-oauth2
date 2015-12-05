package oauth2

import (
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/phuc0302/go-oauth2/utils"
)

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
