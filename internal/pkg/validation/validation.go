package validation

import "regexp"

const emailRegex = `^(?i:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)$`

var re = regexp.MustCompile(emailRegex)

// Validates an email address.
func IsValidEmail(email string) bool {
	return re.MatchString(email)
}
