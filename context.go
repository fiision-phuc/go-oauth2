package oauth2

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"text/template"

	"github.com/phuc0302/go-oauth2/inject"
	"github.com/phuc0302/go-oauth2/util"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// OAuthContext describes a user's oauth scope.
type OAuthContext struct {

	// Registered user. Always available.
	User User
	// Registered client. Always available.
	Client Client
	// Access token that had been given to user. Always available.
	AccessToken Token
	// Refresh token that had been given to user. Might not be available all the time.
	RefreshToken Token
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// RequestContext describes a HTTP URL request scope.
type RequestContext struct {
	Method      string
	Path        string
	Header      map[string]string
	PathParams  map[string]string
	QueryParams map[string]string

	request  *http.Request
	response http.ResponseWriter
}

// BasicAuth returns username & password.
func (c *RequestContext) BasicAuth() (username string, password string, ok bool) {
	username, password, ok = c.request.BasicAuth()
	return
}

// BindForm converts urlencode/multipart form to object.
func (c *RequestContext) BindForm(inputForm interface{}) error {
	return inject.BindForm(c.QueryParams, inputForm)
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
