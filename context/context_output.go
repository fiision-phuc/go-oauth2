package context

import (
	"encoding/json"
	"net/http"
	"text/template"
)

// OutputHeader return an additional header.
func (c *Context) OutputHeader(headerName string, headerValue string) {
	c.response.Header().Set(headerName, headerValue)
}

// OutputError return an error json.
func (c *Context) OutputError(status *Status) {
	c.response.Header().Set("Content-Type", "application/problem+json")
	c.response.WriteHeader(status.Status)

	cause, _ := json.Marshal(status)
	c.response.Write(cause)
}

// OutputRedirect return a redirect instruction.
func (c *Context) OutputRedirect(status *Status, url string) {
	http.Redirect(c.response, c.request, url, status.Status)
}

// OutputJSON return a json.
func (c *Context) OutputJSON(status *Status, model interface{}) {
	c.response.Header().Set("Content-Type", "application/json")
	c.response.WriteHeader(status.Status)

	data, _ := json.Marshal(model)
	c.response.Write(data)
}

// OutputHTML will render a html page.
func (c *Context) OutputHTML(filePath string, model interface{}) {
	tmpl, error := template.ParseFiles(filePath)
	if error != nil {
		c.OutputError(Status404())
	} else {
		tmpl.Execute(c.response, model)
	}
}
