package validation

import "testing"

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
