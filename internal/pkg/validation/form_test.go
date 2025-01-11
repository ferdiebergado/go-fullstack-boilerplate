//go:build !integration

package validation

import (
	"testing"
)

type TestParams struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func TestNewForm(t *testing.T) {
	params := TestParams{
		Name:                 "John Doe",
		Email:                "john.doe@example.com",
		Password:             "password123",
		PasswordConfirmation: "password123",
	}
	form := NewForm(params)

	if form.params != params {
		t.Errorf("expected Params to be %+v, got %+v", params, form.params)
	}

	if form.Error.Count() != 0 {
		t.Errorf("expected Errors to be empty, got %+v", form.Error)
	}
}

func TestFormRequired(t *testing.T) {
	params := TestParams{
		Name:  "",
		Email: "john.doe@example.com",
	}
	form := NewForm(params)

	form.Required("Name", "Email")

	if form.Error.Count() != 1 {
		t.Errorf("expected 1 error, got %d", form.Error.Count())
	}

	if len(form.Error.Get("name")) == 0 {
		t.Errorf("expected error for field 'name' not found")
	}

	if form.Error.Get("name")[0] != "This field is required." {
		t.Errorf("expected error message 'This field is required.', got '%s'", form.Error.Get("name")[0])
	}
}

func TestFormPasswordsMatch(t *testing.T) {
	t.Run("Should pass when the passwords are the same", func(t *testing.T) {
		params := TestParams{
			Password:             "password123",
			PasswordConfirmation: "password123",
		}
		form := NewForm(params)

		form.PasswordsMatch("Password", "PasswordConfirmation")

		if form.Error.Count() != 0 {
			t.Errorf("expected no error, got %d", form.Error.Count())
		}

		if len(form.Error.Get("password")) != 0 {
			t.Errorf("expected no error for field 'password' but found errors")
		}
	})

	t.Run("Should fail when the passwords are not the same", func(t *testing.T) {
		params := TestParams{
			Password:             "password123",
			PasswordConfirmation: "password321",
		}
		form := NewForm(params)

		form.PasswordsMatch("Password", "PasswordConfirmation")

		if form.Error.Count() != 1 {
			t.Errorf("expected 1 error, got %d", form.Error.Count())
		}

		if len(form.Error.Get("password")) == 0 {
			t.Errorf("expected error for field 'password' not found")
		}

		if form.Error.Get("password")[0] != "Passwords do not match." {
			t.Errorf("expected error message 'Passwords do not match.', got '%s'", form.Error.Get("password")[0])
		}
	})
}

func TestFormIsEmail(t *testing.T) {
	params := TestParams{
		Email: "invalid-email",
	}
	form := NewForm(params)

	form.IsEmail("Email")

	if form.Error.Count() != 1 {
		t.Errorf("expected 1 error, got %d", form.Error.Count())
	}

	if len(form.Error.Get("email")) == 0 {
		t.Errorf("expected error for field 'email' not found")
	}

	if form.Error.Get("email")[0] != "Email is not a valid email address." {
		t.Errorf("expected error message 'Email is not a valid email address.', got '%s'", form.Error.Get("email")[0])
	}
}

func TestFormIsValid(t *testing.T) {
	params := TestParams{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}
	form := NewForm(params)

	form.Required("Name", "Email")
	form.IsEmail("Email")

	if !form.IsValid() {
		t.Errorf("expected form to be valid, but it is not")
	}

	// Add invalid email and test
	params.Email = "invalid-email"
	form = NewForm(params)
	form.Required("Name", "Email")
	form.IsEmail("Email")

	if form.IsValid() {
		t.Errorf("expected form to be invalid, but it is valid")
	}
}
