package validator

import "github.com/go-playground/validator/v10"

const PasswordValidator = "password"

const ErrPasswordValidation = "password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one digit and one special character"

func PasswordValidate(fl validator.FieldLevel) bool {
	t := fl.Field().String()

	if len(t) < 8 {
		return false
	}

	var (
		hasUpper, hasLower, hasDigit, hasSpecial bool
	)

	for _, c := range t {
		switch {
		case 'A' <= c && c <= 'Z':
			hasUpper = true
		case 'a' <= c && c <= 'z':
			hasLower = true
		case '0' <= c && c <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	if hasUpper && hasLower && hasDigit && hasSpecial {
		return true
	}

	return false
}
