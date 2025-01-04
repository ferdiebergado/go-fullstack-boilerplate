package validation

import (
	"log/slog"
	"reflect"
	"strings"
)

type Form[T any] struct {
	Params T
	Errors errors
}

func NewForm[T any](params T) *Form[T] {
	return &Form[T]{
		Params: params,
		Errors: make(errors),
	}
}

func (f *Form[T]) Required(fields ...string) {
	for _, field := range fields {
		val := reflect.ValueOf(f.Params)
		if strings.TrimSpace(val.FieldByName(field).String()) == "" {
			jsonTag, ok := GetJSONTag(f.Params, field)
			if !ok {
				slog.Error("cannot find json tag", "field", field)
				return
			}
			f.Errors.Add(jsonTag, "This field is required.")
		}
	}
}

func (f *Form[T]) PasswordsMatch(password string, passwordConfirmation string) {
	val := reflect.ValueOf(f.Params)
	p := val.FieldByName(password).String()
	pc := val.FieldByName(passwordConfirmation).String()

	slog.Debug("passwords match", "password", p, "password_confirmation", pc)

	if p != "" && pc != "" && p != pc {
		jsonTag, ok := GetJSONTag(f.Params, password)
		if !ok {
			slog.Error("cannot find json tag", "password", password, "password_confirmation", passwordConfirmation)
			return
		}
		f.Errors.Add(jsonTag, "Passwords do not match.")
	}
}

func (f *Form[T]) IsEmail(field string) {
	val := reflect.ValueOf(f.Params)
	email := val.FieldByName(field).String()
	if !IsValidEmail(email) {
		jsonTag, ok := GetJSONTag(f.Params, field)
		if !ok {
			slog.Error("cannot find json tag", "field", field)
			return
		}
		f.Errors.Add(jsonTag, "Email is not a valid email address.")
	}
}

func (f *Form[T]) IsValid() bool {
	return len(f.Errors) == 0
}

// GetJSONTag retrieves the JSON tag for a field with the specified name
func GetJSONTag(structure interface{}, fieldName string) (string, bool) {
	typ := reflect.TypeOf(structure)

	if typ.Kind() != reflect.Struct {
		return "", false
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == fieldName {
			return strings.TrimSuffix(field.Tag.Get("json"), ",omitempty"), true
		}
	}
	return "", false
}
