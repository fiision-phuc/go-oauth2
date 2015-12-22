package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// BindForm binds data into given form object.
func BindForm(values url.Values, inputForm interface{}) error {
	/* Condition validation */
	if values == nil || inputForm == nil {
		return fmt.Errorf("Input must not be nil.")
	}

	// Dereference pointer
	reflector := reflect.ValueOf(inputForm)
	if reflector.Kind() == reflect.Ptr {
		reflector = reflector.Elem()
	}
	if reflector.Kind() != reflect.Struct {
		return fmt.Errorf("Input must be struct type.")
	}

	inputStruct := reflector.Type()
	for idx := 0; idx < reflector.NumField(); idx++ {
		property := reflector.Field(idx)
		structField := inputStruct.Field(idx)
		propertyTag := structField.Tag.Get("inject")

		if property.CanSet() && len(propertyTag) != 0 {
			propertyType := property.Type()
			data := values.Get(propertyTag)

			dataType := reflect.TypeOf(data)
			value := reflect.ValueOf(data)

			/* Condition validation: Skip if value is not valid */
			if !value.IsValid() {
				continue
			}

			if dataType == propertyType {
				property.Set(value)

			} else if dataType.Kind() == reflect.String {
				input := strings.ToLower(value.String())

				switch propertyType.Kind() {

				case reflect.Bool:
					b, err := strconv.ParseBool(input)
					if err == nil {
						property.Set(reflect.ValueOf(b))
					}
					break

				case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
					f, err := strconv.ParseFloat(input, 128)
					if err == nil {
						property.Set(reflect.ValueOf(f))
					}
					break

				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					integer, err := strconv.ParseInt(input, 10, 64)
					if err == nil {
						property.Set(reflect.ValueOf(integer))
					}
					break

				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					unsignInteger, err := strconv.ParseUint(input, 10, 64)
					if err == nil {
						property.Set(reflect.ValueOf(unsignInteger))
					}
					break
				}
			}
		}
	}
	return nil
}

// ParseForm parses url-encode form into map.
func ParseForm(request *http.Request) url.Values {
	err := request.ParseForm()
	if err != nil {
		panic(err)
	}
	return request.Form
}

// ParseMultipartForm parses multipart form into map.
func ParseMultipartForm(request *http.Request) url.Values {
	err := request.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		panic(err)
	}
	params := url.Values(request.MultipartForm.Value)
	//	for k, v := range request.URL.Query() {
	//		params[k] = v
	//	}
	return params
}

// ParseStatus parses data into status object.
func ParseStatus(response *http.Response) *Status {
	data, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	status := Status{}
	json.Unmarshal(data, &status)

	return &status
}
