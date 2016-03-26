package oauth2

import (
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/nfnt/resize"
	"github.com/phuc0302/go-oauth2/utils"
)

// RequestContext descripts a HTTP URL request scope.
type RequestContext struct {
	Header map[string]string

	URLPath     string
	Queries     url.Values
	PathQueries map[string]string

	request  *http.Request
	response http.ResponseWriter
}

// CreateRequestContext return a default context.
func CreateRequestContext(request *http.Request, response http.ResponseWriter) *RequestContext {
	context := &RequestContext{
		URLPath:  request.URL.Path,
		request:  request,
		response: response,
	}

	// Format request headers
	context.Header = make(map[string]string, len(request.Header))
	for k, v := range request.Header {
		context.Header[strings.ToLower(k)] = strings.ToLower(v[0])
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
		contentType := request.Header.Get("content-type")

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

/////////////////////////////////////////////////////////////////////////////////////////////////

// BasicAuth returns basic authentication info.
func (c *RequestContext) BasicAuth() (username string, password string, ok bool) {
	return c.request.BasicAuth()
}

// Method returns HTTP method.
func (c *RequestContext) Method() string {
	return c.request.Method
}

// Protocol returns HTTP protocol
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

// MultipartFile returns an uploaded file by name.
func (c *RequestContext) MultipartFile(name string) (multipart.File, *multipart.FileHeader, error) {
	return c.request.FormFile(name)
}

// MoveImage moves a multipart image to destination and resize if neccessary.
func (c *RequestContext) MoveImage(name string, destinationPath string, width uint, height uint) error {
	input, imageInfo, err := c.MultipartFile(name)
	if err != nil {
		return err
	}
	defer input.Close()

	// Decode image
	var decodedImage image.Image
	if path.Ext(imageInfo.Filename) == ".jpg" {
		decodedImage, err = jpeg.Decode(input)
	} else if path.Ext(imageInfo.Filename) == ".png" {
		decodedImage, err = png.Decode(input)
	}

	/* Condition validation: Validate decode image process */
	if err != nil {
		return err
	}

	// Create output file
	output, _ := os.Create(destinationPath)
	defer output.Close()

	// Continue if image can be decoded.
	resizedImage := resize.Resize(width, height, decodedImage, resize.NearestNeighbor)
	jpeg.Encode(output, resizedImage, nil)
	return nil
}

// RawData returns a raw body.
func (c *RequestContext) RawData() ([]byte, error) {
	return ioutil.ReadAll(c.request.Body)
}

/////////////////////////////////////////////////////////////////////////////////////////////////

// OutputHeader returns an additional header.
func (c *RequestContext) OutputHeader(headerName string, headerValue string) {
	c.response.Header().Set(headerName, headerValue)
}

// OutputError returns an error JSON.
func (c *RequestContext) OutputError(status *utils.Status) {
	c.response.Header().Set("Content-Type", "application/problem+json")
	c.response.WriteHeader(status.Code)

	cause, _ := json.Marshal(status)
	c.response.Write(cause)
}

// OutputRedirect returns a redirect instruction.
func (c *RequestContext) OutputRedirect(status *utils.Status, url string) {
	http.Redirect(c.response, c.request, url, status.Code)
}

// OutputJSON returns a JSON.
func (c *RequestContext) OutputJSON(status *utils.Status, model interface{}) {
	c.response.Header().Set("Content-Type", "application/json")
	c.response.WriteHeader(status.Code)

	data, _ := json.Marshal(model)
	c.response.Write(data)
}

// OutputHTML will render a HTML page.
func (c *RequestContext) OutputHTML(filePath string, model interface{}) {
	tmpl, error := template.ParseFiles(filePath)
	if error != nil {
		c.OutputError(utils.Status404())
	} else {
		tmpl.Execute(c.response, model)
	}
}

// OutputText returns a string.
func (c *RequestContext) OutputText(status *utils.Status, data string) {
	c.response.Header().Set("Content-Type", "text/plain")
	c.response.WriteHeader(status.Code)
	c.response.Write([]byte(data))
}
