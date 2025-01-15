package user

import "fmt"

type EmailExistsError struct {
	Email string
}

func (u *EmailExistsError) Error() string {
	return fmt.Sprintf("User with email: %s already exists.", u.Email)
}
