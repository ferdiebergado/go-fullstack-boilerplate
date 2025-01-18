package validation

import (
	"reflect"
	"regexp"
	"strings"
)

const emailRegex = `^(?i:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)$`

var re = regexp.MustCompile(emailRegex)

// Validates an email address.
func IsValidEmail(email string) bool {
	return re.MatchString(email)
}

// Trims string fields of a struct.
func TrimStructFields[T any](s T) {
	v := reflect.ValueOf(s).Elem()

	// Ensure we're working with a struct
	if v.Kind() != reflect.Struct {
		return
	}

	// Iterate through the struct fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// Check if the field is a string and trim it
		if field.Kind() == reflect.String {
			field.SetString(strings.TrimSpace(field.String()))
		}

		// If the field is another struct, recurse into it
		if field.Kind() == reflect.Ptr && !field.IsNil() && field.Elem().Kind() == reflect.Struct {
			TrimStructFields(field.Interface())
		}
		if field.Kind() == reflect.Struct {
			TrimStructFields(field.Addr().Interface())
		}
	}
}
