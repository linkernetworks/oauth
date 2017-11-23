package validator

import "testing"

func TestMatchRegexpEamil(t *testing.T) {
	errorEmailArray := []string{
		"@email.com",
		"a@.com",
		"aemail.com",
		"a@email.",
		"a@emailcom",
		"a@email",
		"##@email.com",
		"a#email.com",
		"a@#com",
		"a@email#com",
		"a@email.#",
	}

	correctEmailArray := []string{
		"a@email.com",
		// "a.b@email.com",
		"a-b@email.com",
		"ab@email.org",
	}

	for _, e := range errorEmailArray {
		if matchRegexpEamil(e) {
			t.Errorf("email format %s error", e)
		}
	}

	for _, c := range correctEmailArray {
		if !matchRegexpEamil(c) {
			t.Errorf("email format %s should be correct", c)
		}
	}
}

func TestMatchRegexpCellphone(t *testing.T) {
	errorCellphone := []string{
		"13815380336",
		"110",
		"128",
		"+86114",
	}

	correctCellphone := []string{
		"+8615651002356",
		"+8812565116558",
	}

	for _, e := range errorCellphone {
		if matchRegexpCellphone(e) {
			t.Errorf("cellphone format %s error", e)
		}
	}

	for _, c := range correctCellphone {
		if matchRegexpCellphone(c) {
			t.Errorf("cellphone format %s should be correct", c)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	ret, err := ValidateEmail("")
	if err == nil {
		t.Error()
	}

	if !ret.Error {
		t.Error()
	}

	if ret.Message != "email is required." {
		t.Error()
	}
}

func TestValidatePassword(t *testing.T) {
	ret1, err1 := ValidatePassword("")
	if err1 == nil {
		t.Error()
	}
	if ret1.Message != "password is required." {
		t.Error()
	}

	ret2, err2 := ValidatePassword("123")
	if err2 == nil {
		t.Error()
	}
	if ret2.Message != "password length must be equal or longer than six." {
		t.Error()
	}

	// correct
	_, err3 := ValidatePassword("123456")
	if err3 != nil {
		t.Error()
	}
}

func TestValidateCellphone(t *testing.T) {
	ret, err := ValidateCellphone("")
	if err == nil {
		t.Error()
	}
	if ret.Message != "cellphone is required." {
		t.Error()
	}
}
