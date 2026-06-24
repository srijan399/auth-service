package validation

import (
	"errors"
	"net/mail"
)

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email required")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email")
	}

	return nil
}
