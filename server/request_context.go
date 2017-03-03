package server

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/julienschmidt/httprouter"
	"github.com/phuc0302/go-oauth2/util"
)

// RequestContext describes a HTTP URL request scope.
type RequestContext struct {
	Method      string
	Path        string
	Header      map[string]string
	PathParams  map[string]string
	QueryParams map[string]string

	request  *http.Request
	response http.ResponseWriter
	extra    map[string]interface{}
}

// CreateContext creates new request context.
func CreateContext(response http.ResponseWriter, request *http.Request) *RequestContext {
	context := &RequestContext{
		Path:   httprouter.CleanPath(request.URL.Path),
		Method: strings.ToLower(request.Method),

		request:  request,
		response: response,
		extra:    make(map[string]interface{}),
	}

	// Format request headers
	if len(request.Header) > 0 {
		context.Header = make(map[string]string)

		for k, v := range request.Header {
			if header := strings.ToLower(k); header == "authorization" {
				context.Header[header] = v[0]
			} else {
				context.Header[header] = strings.ToLower(v[0])
			}
		}
	}

	// Parse body context if neccessary
	var params url.Values
	switch context.Method {

	case Get:
		params = request.URL.Query()
		break

	case Patch, Post:
		if contentType := context.Header["content-type"]; contentType == "application/x-www-form-urlencoded" {
			if err := request.ParseForm(); err == nil {
				params = request.Form
			}
		} else if strings.HasPrefix(contentType, "multipart/form-data; boundary") {
			if err := request.ParseMultipartForm(Cfg.MultipartSize); err == nil {
				params = request.MultipartForm.Value
			}
		}
		break

	default:
		break
	}

	// Process params
	if len(params) > 0 {
		context.QueryParams = make(map[string]string)

		for k, v := range params {
			context.QueryParams[k] = v[0]
		}
	}
	return context
}

// BasicAuth returns username & password.
func (c *RequestContext) BasicAuth() (username string, password string, ok bool) {
	username, password, ok = c.request.BasicAuth()
	return
}

// BindForm converts urlencode/multipart form to object.
func (c *RequestContext) BindForm(inputForm interface{}) error {
	return util.BindForm(c.QueryParams, inputForm)
}

// BindJSON converts json data to object.
func (c *RequestContext) BindJSON(jsonObject interface{}) error {
	bytes, err := ioutil.ReadAll(c.request.Body)

	/* Condition validation: Validate parse process */
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, jsonObject)
	return err
}

// MultipartFile returns an uploaded file by name.
func (c *RequestContext) MultipartFile(name string) (multipart.File, *multipart.FileHeader, error) {
	return c.request.FormFile(name)
}

//// MoveImage moves a multipart image to destination and resize if neccessary.
//func (c *Request) MoveImage(name string, destinationPath string, width uint, height uint) error {
//	input, imageInfo, err := c.MultipartFile(name)
//	if err != nil {
//		return err
//	}
//	defer input.Close()

//	// Decode image
//	var decodedImage image.Image
//	if path.Ext(imageInfo.Filename) == ".jpg" {
//		decodedImage, err = jpeg.Decode(input)
//	} else if path.Ext(imageInfo.Filename) == ".png" {
//		decodedImage, err = png.Decode(input)
//	}

//	/* Condition validation: Validate decode image process */
//	if err != nil {
//		return err
//	}

//	// Create output file
//	output, _ := os.Create(destinationPath)
//	defer output.Close()

//	// Continue if image can be decoded.
//	resizedImage := resize.Resize(width, height, decodedImage, resize.NearestNeighbor)
//	jpeg.Encode(output, resizedImage, nil)
//	return nil
//}

// OutputHeader returns an additional header.
func (c *RequestContext) OutputHeader(headerName string, headerValue string) {
	c.response.Header().Set(headerName, headerValue)
}

// OutputError returns an error JSON.
func (c *RequestContext) OutputError(status *util.Status) {
	if redirectURL := redirectPaths[status.Code]; len(redirectURL) > 0 {
		c.OutputRedirect(status, redirectURL)
	} else {
		c.response.Header().Set("Content-Type", "application/problem+json")
		c.response.WriteHeader(status.Code)
		cause, _ := json.Marshal(status)
		c.response.Write(cause)
	}
}

// OutputRedirect returns a redirect instruction.
func (c *RequestContext) OutputRedirect(status *util.Status, url string) {
	http.Redirect(c.response, c.request, url, status.Code)
}

// OutputJSON returns a JSON.
func (c *RequestContext) OutputJSON(status *util.Status, model interface{}) {
	c.response.Header().Set("Content-Type", "application/json")
	c.response.WriteHeader(status.Code)
	data, _ := json.Marshal(model)
	c.response.Write(data)
}

// OutputHTML returns a HTML page.
func (c *RequestContext) OutputHTML(filePath string, model interface{}) {
	if tmpl, err := template.ParseFiles(filePath); err == nil {
		tmpl.Execute(c.response, model)
	} else {
		c.OutputError(util.Status404())
	}
}

// OutputText returns a string.
func (c *RequestContext) OutputText(status *util.Status, data string) {
	c.response.Header().Set("Content-Type", "text/plain")
	c.response.WriteHeader(status.Code)
	c.response.Write([]byte(data))
}

// GetExtra returns extra data that had been associated with key if there is any.
func (c *RequestContext) GetExtra(key string) interface{} {
	return c.extra[key]
}

// SetExtra associates value with key.
func (c *RequestContext) SetExtra(key string, value interface{}) {
	c.extra[key] = value
}
