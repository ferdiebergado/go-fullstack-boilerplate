package validation

import "regexp"

// Validates an email address.
func IsValidEmail(email string) bool {
	const emailRegex = `^(?i:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)$`

	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}
