package validation

import (
	"log/slog"
	"reflect"
	"strings"
)

type Form[T any] struct {
	params T
	val    reflect.Value
	Error  *Error
}

func NewForm[T any](params T) *Form[T] {
	return &Form[T]{
		params: params,
		val:    reflect.ValueOf(params),
		Error:  NewError(),
	}
}

func (f *Form[T]) Required(fields ...string) {
	for _, field := range fields {
		if strings.TrimSpace(f.val.FieldByName(field).String()) == "" {
			jsonTag := f.getJSONTag(field)
			f.Error.Add(jsonTag, "This field is required.")
		}
	}
}

func (f *Form[T]) PasswordsMatch(password string, passwordConfirmation string) {
	p := f.val.FieldByName(password).String()
	pc := f.val.FieldByName(passwordConfirmation).String()

	slog.Debug("passwords match", "password", p, "password_confirmation", pc)

	if p != "" && pc != "" && p != pc {
		jsonTag := f.getJSONTag(password)
		f.Error.Add(jsonTag, "Passwords do not match.")
	}
}

func (f *Form[T]) IsEmail(field string) {
	email := f.val.FieldByName(field).String()
	if !IsValidEmail(email) {
		jsonTag := f.getJSONTag(field)
		f.Error.Add(jsonTag, "Email is not a valid email address.")
	}
}

func (f *Form[T]) IsValid() bool {
	return f.Error.Count() == 0
}

func (f *Form[T]) getJSONTag(field string) string {
	jsonTag, ok := GetJSONTag(f.params, field)
	if !ok {
		slog.Error("cannot find json tag", "field", field)
		return field
	}

	return jsonTag
}

// GetJSONTag retrieves the JSON tag for a field with the specified name
func GetJSONTag(structure any, fieldName string) (string, bool) {
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
