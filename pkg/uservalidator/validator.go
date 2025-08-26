package uservalidator

import (
	"errors"
	"unicode"
)

var (
	ErrPasswordTooShort  = errors.New("password too short")
	ErrPasswordNoDigit   = errors.New("password must contain at least one digit")
	ErrPasswordNoUpper   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoSpecial = errors.New("password must contain at least one special character")
)

const (
	minPasswordLength = 8
	maxEmailLength    = 256
)

type UserValidator interface {
	ValidatePassword(password string) error
	ValidateEmail(email string) error
}

type userValidator struct{}

func New() UserValidator {
	return &userValidator{}
}

func (v *userValidator) ValidatePassword(password string) error {
	if len(password) < minPasswordLength {
		return ErrPasswordTooShort
	}

	hasDigit := false
	hasUpper := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsUpper(char):
			hasUpper = true
		case !unicode.IsLetter(char) && !unicode.IsDigit(char):
			hasSpecial = true
		}
	}

	if !hasDigit {
		return ErrPasswordNoDigit
	}
	if !hasUpper {
		return ErrPasswordNoUpper
	}
	if !hasSpecial {
		return ErrPasswordNoSpecial
	}

	return nil
}

func (v *userValidator) ValidateEmail(email string) error {
	return nil
}
