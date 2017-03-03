package util

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// BindForm binds data into given form object.
func BindForm(values map[string]string, inputForm interface{}) error {
	/* Condition validation */
	if inputForm == nil {
		panic(Status500WithDescription("InputForm must not be nil."))
	}

	/* Condition validation: validate inputForm type */
	reflector := reflect.ValueOf(inputForm)
	if reflector.Kind() == reflect.Ptr {
		reflector = reflector.Elem()
	}
	if reflector.Kind() != reflect.Struct {
		panic(Status500WithDescription("InputForm must be struct."))
	}

	// Binding process
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
				input := strings.TrimSpace(value.String())
				validation := structField.Tag.Get("validation")

				// Validation input value before inject
				if len(validation) > 0 {
					regex := regexp.MustCompile(validation)
					if !regex.MatchString(input) {
						return fmt.Errorf("Invalid '%s' parameter.", field)
					}
				}

				// Convert process
				switch propertyType.Kind() {

				case reflect.Bool:
					input = strings.ToLower(input)

					if b, err := strconv.ParseBool(input); err == nil {
						property.Set(reflect.ValueOf(b))
						break
					}
				case reflect.Float32:
					if f, err := strconv.ParseFloat(input, 32); err == nil {
						property.Set(reflect.ValueOf(float32(f)))
						break
					}
				case reflect.Float64:
					if f, err := strconv.ParseFloat(input, 64); err == nil {
						property.Set(reflect.ValueOf(f))
						break
					}

				case reflect.Int, reflect.Int32:
					if integer, err := strconv.ParseInt(input, 10, 32); err == nil {
						property.Set(reflect.ValueOf(int(integer)))
						break
					}
				case reflect.Int8:
					if integer, err := strconv.ParseInt(input, 10, 0); err == nil {
						property.Set(reflect.ValueOf(int8(integer)))
						break
					}
				case reflect.Int16:
					if integer, err := strconv.ParseInt(input, 10, 16); err == nil {
						property.Set(reflect.ValueOf(int16(integer)))
						break
					}
				case reflect.Int64:
					if integer, err := strconv.ParseInt(input, 10, 64); err == nil {
						property.Set(reflect.ValueOf(integer))
						break
					}

				case reflect.Uint, reflect.Uint32:
					if unsignInteger, err := strconv.ParseUint(input, 10, 32); err == nil {
						property.Set(reflect.ValueOf(uint(unsignInteger)))
						break
					}
				case reflect.Uint8:
					if unsignInteger, err := strconv.ParseUint(input, 10, 8); err == nil {
						property.Set(reflect.ValueOf(uint8(unsignInteger)))
						break
					}
				case reflect.Uint16:
					if unsignInteger, err := strconv.ParseUint(input, 10, 16); err == nil {
						property.Set(reflect.ValueOf(uint16(unsignInteger)))
						break
					}
				case reflect.Uint64:
					if unsignInteger, err := strconv.ParseUint(input, 10, 64); err == nil {
						property.Set(reflect.ValueOf(unsignInteger))
						break
					}

				case reflect.String:
					if len(value.String()) > 0 {
						property.Set(value)
						break
					}

				default:
					return fmt.Errorf("Invalid \"%s\" parameter.", field)
				}
			}
		}
	}
	return nil
}
