package context

import "mime/multipart"

// BindForm converts urlencode/multipart form to object.
func (c *Context) BindForm(formObject interface{}) error {
	//	return c.request.FormFile(name)
}

// BindJSON converts json form to object.
func (c *Context) BindJSON(jsonObject interface{}) error {
	//	return c.request.FormFile(name)
}

// GetMultipartFile return an upload file by name.
func (c *Context) GetMultipartFile(name string) (multipart.File, *multipart.FileHeader, error) {
	return c.request.FormFile(name)
}
