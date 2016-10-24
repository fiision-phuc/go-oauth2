package inject

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/phuc0302/go-oauth2/util"
)

// BindForm binds data into given form object.
func BindForm(values map[string]string, inputForm interface{}) error {
	/* Condition validation */
	if values == nil || inputForm == nil {
		panic(util.Status500WithDescription("InputForm must not be nil."))
	}

	/* Condition validation: validate inputForm type */
	reflector := reflect.ValueOf(inputForm)
	if reflector.Kind() == reflect.Ptr {
		reflector = reflector.Elem()
	}
	if reflector.Kind() != reflect.Struct {
		panic(util.Status500WithDescription("InputForm must be struct."))
	}

	inputStruct := reflector.Type()
	for idx := 0; idx < reflector.NumField(); idx++ {
		property := reflector.Field(idx)
		structField := inputStruct.Field(idx)
		field := structField.Tag.Get("field")

		if property.CanSet() && len(field) > 0 {
			data := values[field]
			propertyType := property.Type()

			value := reflect.ValueOf(data)
			dataType := reflect.TypeOf(data)

			/* Condition validation: Skip if value is not valid */
			if !value.IsValid() {
				continue
			}

			if dataType == propertyType && dataType.Kind() != reflect.String {
				property.Set(value)
			} else if dataType.Kind() == reflect.String {
				validation := structField.Tag.Get("validation")
				input := strings.ToLower(value.String())

				switch propertyType.Kind() {

				case reflect.Bool:
					b, err := strconv.ParseBool(input)
					if err == nil {
						property.Set(reflect.ValueOf(b))
					}
					break

				case reflect.Float32:
					f, err := strconv.ParseFloat(input, 32)
					if err == nil {
						property.Set(reflect.ValueOf(float32(f)))
					}
					break

				case reflect.Float64:
					f, err := strconv.ParseFloat(input, 64)
					if err == nil {
						property.Set(reflect.ValueOf(f))
					}
					break

				case reflect.Int, reflect.Int32:
					integer, err := strconv.ParseInt(input, 10, 32)
					if err == nil {
						property.Set(reflect.ValueOf(int(integer)))
					}
					break

				case reflect.Int8:
					integer, err := strconv.ParseInt(input, 10, 0)
					if err == nil {
						property.Set(reflect.ValueOf(int8(integer)))
					}
					break

				case reflect.Int16:
					integer, err := strconv.ParseInt(input, 10, 16)
					if err == nil {
						property.Set(reflect.ValueOf(int16(integer)))
					}
					break

				case reflect.Int64:
					integer, err := strconv.ParseInt(input, 10, 64)
					if err == nil {
						property.Set(reflect.ValueOf(integer))
					}
					break

				case reflect.String:
					if len(value.String()) > 0 {
						property.Set(value)
					}
					break

				case reflect.Uint, reflect.Uint32:
					unsignInteger, err := strconv.ParseUint(input, 10, 32)
					if err == nil {
						property.Set(reflect.ValueOf(uint(unsignInteger)))
					}
					break

				case reflect.Uint8:
					unsignInteger, err := strconv.ParseUint(input, 10, 8)
					if err == nil {
						property.Set(reflect.ValueOf(uint8(unsignInteger)))
					}
					break

				case reflect.Uint16:
					unsignInteger, err := strconv.ParseUint(input, 10, 16)
					if err == nil {
						property.Set(reflect.ValueOf(uint16(unsignInteger)))
					}
					break

				case reflect.Uint64:
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
