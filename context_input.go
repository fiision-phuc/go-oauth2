package oauth2

import (
	"mime/multipart"

	"github.com/phuc0302/go-oauth2/utils"
)

func (c *Context) BasicAuth() (username string, password string, ok bool) {
	return c.request.BasicAuth()
}

func (c *Context) Header(headerName string) string {
	return c.request.Header.Get(headerName)
}

func (c *Context) Method() string {
	return c.request.Method
}

func (c *Context) Protocol() string {
	return c.request.Proto
}

// BindForm converts urlencode/multipart form to object.
func (c *Context) BindForm(inputForm interface{}) error {
	return utils.BindForm(c.Queries, inputForm)
}

// BindJSON converts json data to object.
func (c *Context) BindJSON(jsonObject interface{}) error {
	//	return c.request.FormFile(name)
	return nil
}

// GetMultipartFile return an upload file by name.
func (c *Context) MultipartFile(name string) (multipart.File, *multipart.FileHeader, error) {
	return c.request.FormFile(name)
}
