package oauth2

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"text/template"

	"github.com/phuc0302/go-oauth2/utils"
)

// Request describes a HTTP URL request scope.
type Request struct {
	Path        string
	Header      map[string]string
	PathParams  map[string]string
	QueryParams map[string]string

	request  *http.Request
	response http.ResponseWriter
}

// BindForm converts urlencode/multipart form to object.
func (c *Request) BindForm(inputForm interface{}) error {
	return utils.BindForm(c.QueryParams, inputForm)
}

// BindJSON converts json data to object.
func (c *Request) BindJSON(jsonObject interface{}) error {
	bytes, err := ioutil.ReadAll(c.request.Body)

	/* Condition validation: Validate parse process */
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, jsonObject)
	return err
}

// MultipartFile returns an uploaded file by name.
func (c *Request) MultipartFile(name string) (multipart.File, *multipart.FileHeader, error) {
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

////////////////////////////////////////////////////////////////////////////////////////////////////

// OutputHeader returns an additional header.
func (c *Request) OutputHeader(headerName string, headerValue string) {
	c.response.Header().Set(headerName, headerValue)
}

// OutputError returns an error JSON.
func (c *Request) OutputError(status *utils.Status) {
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
func (c *Request) OutputRedirect(status *utils.Status, url string) {
	http.Redirect(c.response, c.request, url, status.Code)
}

// OutputJSON returns a JSON.
func (c *Request) OutputJSON(status *utils.Status, model interface{}) {
	c.response.Header().Set("Content-Type", "application/json")
	c.response.WriteHeader(status.Code)

	data, _ := json.Marshal(model)
	c.response.Write(data)
}

// OutputHTML will render a HTML page.
func (c *Request) OutputHTML(filePath string, model interface{}) {
	tmpl, error := template.ParseFiles(filePath)
	if error != nil {
		c.OutputError(utils.Status404())
	} else {
		tmpl.Execute(c.response, model)
	}
}

// OutputText returns a string.
func (c *Request) OutputText(status *utils.Status, data string) {
	c.response.Header().Set("Content-Type", "text/plain")
	c.response.WriteHeader(status.Code)
	c.response.Write([]byte(data))
}
