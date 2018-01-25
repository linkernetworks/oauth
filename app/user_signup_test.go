package app

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserSignUp(t *testing.T) {
	// EXECUTOR_NUMBER is a jenkins environment variable
	if os.Getenv("EXECUTOR_NUMBER") != "" {
		t.Skip("Fix this for concurrent build")
	}

	phoneNumber, ok := os.LookupEnv("TEST_SMS_PHONENUMBRR")
	if !ok {
		t.Skipf("need env TEST_SMS_PHONENUMBRR, skip user signup test.\n")
	}
	setupAPIServer(t)
	testLocale := testLocaleFunc(DefaultLocaleLang)

	userSinupAPI := apiBaseUrl + "/v1/signup"
	userPost := url.Values{"email": {"linker@linker.com"}, "password": {"123456"}, "cellphone": {phoneNumber}}

	// test user sign up success
	resp1, err := http.PostForm(userSinupAPI, userPost)
	if err != nil {
		t.Errorf("post form data to %s error: %v", userSinupAPI, err)
	}
	defer resp1.Body.Close()

	// test user has signup
	ret2 := PostForm(userSinupAPI, userPost)
	m2 := ret2.(map[string]interface{})
	assert.Equal(t, m2["message"], testLocale(ErrUserExisted))

	t.Logf("TestAppSignup done, clear mongo db")
	clearDatabase()
}
