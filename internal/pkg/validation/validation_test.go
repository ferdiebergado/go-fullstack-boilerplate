//go:build !integration

package validation

import (
	"reflect"
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		// Valid emails
		{"example@example.com", true},
		{"user+alias@domain.co.uk", true},
		{"user_name@sub.domain.com", true},
		{"123@domain.com", true},
		{"user@domain.io", true},

		// Invalid emails
		{"plainaddress", false},
		{"@missingusername.com", false},
		{"username@.com", false},
		{"username@domain,com", false},
		{"username@domain..com", false},
		{"username@.domain.com", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := IsValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("IsValidEmail(%q) = %v; want %v", tt.email, result, tt.expected)
			}
		})
	}
}

// Define test structs
type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

type Person struct {
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Address Address `json:"address"`
}

func TestTrimStructFields(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name: "Trim string fields",
			input: &Person{
				Name:  " John Doe ",
				Email: " john.doe@example.com ",
				Address: Address{
					Street: " 123 Main St ",
					City:   " Some City ",
				},
			},
			expected: &Person{
				Name:  "John Doe",
				Email: "john.doe@example.com",
				Address: Address{
					Street: "123 Main St",
					City:   "Some City",
				},
			},
		},
		{
			name: "No trimming needed",
			input: &Person{
				Name:  "John Doe",
				Email: "john.doe@example.com",
				Address: Address{
					Street: "123 Main St",
					City:   "Some City",
				},
			},
			expected: &Person{
				Name:  "John Doe",
				Email: "john.doe@example.com",
				Address: Address{
					Street: "123 Main St",
					City:   "Some City",
				},
			},
		},
		{
			name: "Empty string fields",
			input: &Person{
				Name:  " ",
				Email: " ",
				Address: Address{
					Street: " ",
					City:   " ",
				},
			},
			expected: &Person{
				Name:  "",
				Email: "",
				Address: Address{
					Street: "",
					City:   "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TrimStructFields(tt.input)

			// Compare the result with the expected value
			if !reflect.DeepEqual(tt.input, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, tt.input)
			}
		})
	}
}
