package validator

import (
	"errors"
	"regexp"
)

// ValidateEmail
func ValidateEmail(email string) (FieldValidation, error) {
	var emailValidation = FieldValidation{}

	if email == "" {
		emailValidation.Field = "email"
		emailValidation.Error = true
		emailValidation.Message = "email is required."

		return emailValidation, errors.New("email is required.")
	}

	if !matchRegexpEamil(email) {
		emailValidation.Field = "email"
		emailValidation.Error = true
		emailValidation.Message = "email is not validate."

		return emailValidation, errors.New("email is not validate.")
	}

	return emailValidation, nil
}

// matchRegexpEamil validate email by regexp
func matchRegexpEamil(email string) bool {
	emailReg := regexp.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
	return emailReg.MatchString(email)
}

// ValidatePassword
func ValidatePassword(password string) (FieldValidation, error) {
	var passwordValidation = FieldValidation{}

	if password == "" {
		passwordValidation.Field = "password"
		passwordValidation.Error = true
		passwordValidation.Message = "password is required."
		return passwordValidation, errors.New("password is required.")
	}

	if len(password) < 6 {
		passwordValidation.Field = "password"
		passwordValidation.Error = true
		passwordValidation.Message = "password length must be equal or longer than six."

		return passwordValidation, errors.New("password length must be equal or longer than six.")
	}

	return passwordValidation, nil
}

func ValidateCellphone(cellphone string) (FieldValidation, error) {
	var phoneValidation = FieldValidation{}
	if cellphone == "" {
		phoneValidation.Field = "cellphone"
		phoneValidation.Error = true
		phoneValidation.Message = "cellphone is required."

		return phoneValidation, errors.New("cellphone is required.")
	}

	// if !matchRegexpCellphone(cellphone) {
	// 	phoneValidation.Field = "password"
	// 	phoneValidation.Error = true
	// 	phoneValidation.Message = "cellphone is not validate."

	// 	return phoneValidation, errors.New("cellphone is not validate.")
	// }

	return phoneValidation, nil
}

func matchRegexpCellphone(cellphone string) bool {
	// refer to
	// https://stackoverflow.com/questions/2113908/what-regular-expression-will-match-valid-international-phone-numbers
	regExp := `\+(9[976]\d|8[987530]\d|6[987]\d|5[90]\d|42\d|3[875]\d|
2[98654321]\d|9[8543210]|8[6421]|6[6543210]|5[87654321]|
4[987654310]|3[9643210]|2[70]|7|1)
\W*\d\W*\d\W*\d\W*\d\W*\d\W*\d\W*\d\W*\d\W*(\d{1,2})$`
	phoneReg := regexp.MustCompile(regExp)
	return phoneReg.MatchString(cellphone)
}
